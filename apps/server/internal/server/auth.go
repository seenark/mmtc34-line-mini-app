package server

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

func (s *Server) mountAuth(r routes) {
	r.public.POST("/auth/line", s.lineLogin)
	r.public.POST("/auth/dev", s.devLogin) // 404 unless DEV_LOGIN
	r.public.POST("/auth/logout", s.logout)
	r.self.GET("/me", s.me)
	r.self.POST("/me", s.onboard)
}

// lineLogin verifies the ID token with LINE and starts a session for its userId.
func (s *Server) lineLogin(c *gin.Context) {
	in, ok := bind[api.LineLoginInput](c)
	if !ok {
		return
	}
	lineUserID, err := s.verify.Verify(c.Request.Context(), in.IDToken)
	if err != nil {
		slog.Warn("line verify failed", "err", err)
		abort(c, http.StatusUnauthorized, "ยืนยันตัวตนกับ LINE ไม่สำเร็จ")
		return
	}
	s.login(c, lineUserID)
}

// devLogin is lineLogin with the userId taken on trust. The route does not exist unless DEV_LOGIN is on.
func (s *Server) devLogin(c *gin.Context) {
	if !s.cfg.DevLogin {
		fail(c, &domain.Error{Kind: domain.NotFound, Msg: "ไม่พบข้อมูล"})
		return
	}
	in, ok := bind[api.DevLoginInput](c)
	if !ok {
		return
	}
	id := strings.TrimSpace(in.LineUserID)
	if id == "" {
		fail(c, &domain.Error{Kind: domain.Invalid, Msg: "กรุณากรอก LINE userId"})
		return
	}
	s.login(c, id)
}

// login promotes a bootstrap Admin's existing row (a Member created before the id was listed)
// and sets the cookie.
func (s *Server) login(c *gin.Context, lineUserID string) {
	if domain.IsBootstrapAdmin(lineUserID, s.cfg.AdminLineUserIDs) {
		if err := s.q.PromoteBootstrapAdmin(c.Request.Context(), &lineUserID); err != nil {
			fail(c, err)
			return
		}
	}
	s.setSession(c, lineUserID)
	c.Status(http.StatusNoContent)
}

func (s *Server) logout(c *gin.Context) {
	s.clearSession(c)
	c.Status(http.StatusNoContent)
}

// me returns api.Me; Member null means "show onboarding".
func (s *Server) me(c *gin.Context) {
	v := viewerOf(c)
	me := api.Me{LineUserID: v.LineUserID}
	if v.Member != nil {
		m := memberDTO(*v.Member)
		me.Member = &m
	}
	c.JSON(http.StatusOK, me)
}

// onboard creates the Applicant (or an approved Admin for a bootstrap id). 409 if this LINE
// userId already has a Member.
func (s *Server) onboard(c *gin.Context) {
	in, ok := bind[api.ProfileInput](c)
	if !ok {
		return
	}
	in, err := domain.ValidateProfile(in, domain.RequireIdentity)
	if err != nil {
		fail(c, err)
		return
	}
	lineUserID := viewerOf(c).LineUserID
	status, role := api.MemberApplicant, api.RoleMember
	if domain.IsBootstrapAdmin(lineUserID, s.cfg.AdminLineUserIDs) {
		status, role = api.MemberApproved, api.RoleAdmin
	}

	var m db.Member
	err = s.tx(c.Request.Context(), func(q *db.Queries) error {
		_, err := q.GetMemberByLineUserID(c.Request.Context(), &lineUserID)
		if err == nil {
			return &domain.Error{Kind: domain.Conflict, Msg: "ลงทะเบียนแล้ว"}
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
		m, err = q.CreateMember(c.Request.Context(), db.CreateMemberParams{
			LineUserID: &lineUserID,
			Status:     string(status),
			Role:       string(role),
			Cohort:     in.Cohort,
			Nickname:   in.Nickname,
			RealName:   in.RealName,
			Address:    in.Address,
			Phone:      in.Phone,
			Now:        rfc3339(s.now()),
		})
		return err
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, memberDTO(m))
}
