package server

import (
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

func (s *Server) mountAuth(r routes) {
	r.public.POST("/auth/line", s.lineLogin)
	r.public.POST("/auth/dev", s.devLogin) // 404 unless DEV_LOGIN
	r.public.POST("/auth/admin", s.adminLogin)
	r.public.POST("/auth/logout", s.logout)
	r.self.GET("/me", s.me)
	r.self.POST("/me", s.onboard)
}

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
	if lineUserID == domain.SuperAdminLineUserID {
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
	if id == domain.SuperAdminLineUserID {
		fail(c, &domain.Error{Kind: domain.Invalid, Msg: "ใช้ LINE userId นี้ไม่ได้"})
		return
	}
	s.login(c, id)
}

const (
	adminLoginWindow   = time.Minute
	adminLoginMaxFails = 5
)

// adminLogin signs in the Super Admin with the env username and password. The row is upserted
// so a changed env or a wiped database needs no manual step; no other row is touched.
// shortcut: the failure counter is process-wide and in memory, so a restart resets it; upgrade if the app runs more than one instance.
func (s *Server) adminLogin(c *gin.Context) {
	if s.cfg.SuperAdminPassword == "" {
		fail(c, &domain.Error{Kind: domain.NotFound, Msg: "ไม่พบข้อมูล"})
		return
	}
	in, ok := bind[api.AdminLoginInput](c)
	if !ok {
		return
	}
	if !s.adminAttemptAllowed() {
		abort(c, http.StatusTooManyRequests, "ลองใหม่อีกครั้งในอีกสักครู่")
		return
	}
	// & (not &&) so both fields are always compared.
	if sameSecret(in.Username, s.cfg.SuperAdminUsername)&sameSecret(in.Password, s.cfg.SuperAdminPassword) != 1 {
		s.recordAdminFailure()
		abort(c, http.StatusUnauthorized, "ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง")
		return
	}
	id := domain.SuperAdminLineUserID
	err := s.q.UpsertSuperAdmin(c.Request.Context(), db.UpsertSuperAdminParams{
		LineUserID: &id,
		Nickname:   "ผู้ดูแลระบบ",
		Now:        rfc3339(s.now()),
	})
	if err != nil {
		fail(c, err)
		return
	}
	s.setSession(c, id)
	c.Status(http.StatusNoContent)
}

// sameSecret hashes first so the compare takes the same time whatever the lengths.
func sameSecret(a, b string) int {
	ha, hb := sha256.Sum256([]byte(a)), sha256.Sum256([]byte(b))
	return subtle.ConstantTimeCompare(ha[:], hb[:])
}

// adminAttemptAllowed drops failures older than the window, then allows another attempt
// while fewer than adminLoginMaxFails remain.
func (s *Server) adminAttemptAllowed() bool {
	s.adminLoginMu.Lock()
	defer s.adminLoginMu.Unlock()
	cutoff := s.now().Add(-adminLoginWindow)
	kept := s.adminLoginFails[:0]
	for _, t := range s.adminLoginFails {
		if t.After(cutoff) {
			kept = append(kept, t)
		}
	}
	s.adminLoginFails = kept
	return len(kept) < adminLoginMaxFails
}

func (s *Server) recordAdminFailure() {
	s.adminLoginMu.Lock()
	defer s.adminLoginMu.Unlock()
	s.adminLoginFails = append(s.adminLoginFails, s.now())
}

// login sets the cookie for a verified LINE userId.
func (s *Server) login(c *gin.Context, lineUserID string) {
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

// onboard creates the Applicant. 409 if this LINE userId already has a Member.
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
			Status:     string(api.MemberApplicant),
			Role:       string(api.RoleMember),
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
