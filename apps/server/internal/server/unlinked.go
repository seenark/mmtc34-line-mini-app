package server

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

func (s *Server) mountUnlinked(r routes) {
	r.admin.POST("/members", s.createUnlinked)
	r.admin.POST("/members/:memberId/link", s.linkMember)
}

// createUnlinked records a classmate who has no LINE account yet: approved, no LINE userId.
func (s *Server) createUnlinked(c *gin.Context) {
	in, ok := bind[api.ProfileInput](c)
	if !ok {
		return
	}
	in, err := domain.ValidateProfile(in, domain.RequireNicknameOnly)
	if err != nil {
		fail(c, err)
		return
	}
	m, err := s.q.CreateMember(c.Request.Context(), db.CreateMemberParams{
		Status:   string(api.MemberApproved),
		Role:     string(api.RoleMember),
		Cohort:   in.Cohort,
		Nickname: in.Nickname,
		RealName: in.RealName,
		Address:  in.Address,
		Phone:    in.Phone,
		Now:      rfc3339(s.now()),
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, memberDTO(m))
}

// linkMember moves an Applicant's LINE userId onto an Unlinked Member (ADR 0002). The Applicant
// row goes first because it holds the unique userId; sessions follow the userId, so the LINE
// user's next request already sees the Unlinked Member's Orders.
func (s *Server) linkMember(c *gin.Context) {
	id, ok := idParam(c, "memberId")
	if !ok {
		return
	}
	in, ok := bind[api.LinkInput](c)
	if !ok {
		return
	}
	ctx := c.Request.Context()
	var linked db.Member
	err := s.tx(ctx, func(q *db.Queries) error {
		target, err := q.GetMember(ctx, id)
		if err != nil {
			return err
		}
		if target.LineUserID != nil || target.DeletedAt != nil {
			return &domain.Error{Kind: domain.Conflict, Msg: "สมาชิกคนนี้เชื่อมกับบัญชี LINE แล้วหรือถูกลบไปแล้ว"}
		}
		app, err := q.GetMember(ctx, in.ApplicantID)
		if errors.Is(err, sql.ErrNoRows) || err == nil && app.Status != string(api.MemberApplicant) {
			return &domain.Error{Kind: domain.Invalid, Msg: "ไม่พบผู้สมัครที่รออนุมัติ"}
		}
		if err != nil {
			return err
		}
		if err := q.DeleteApplicant(ctx, app.ID); err != nil {
			return err
		}
		linked, err = q.LinkMember(ctx, db.LinkMemberParams{
			ID: id, LineUserID: app.LineUserID,
			Cohort: app.Cohort, RealName: app.RealName, Address: app.Address, Phone: app.Phone,
		})
		return err
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, memberDTO(linked))
}
