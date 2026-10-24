package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

func (s *Server) mountMembers(r routes) {
	r.admin.GET("/members", s.listMembers) // ?status=applicant|approved, absent = all
	r.admin.GET("/members/:memberId", s.getMember)
	r.admin.POST("/members/:memberId/approve", s.approveApplicant)
	r.admin.POST("/members/:memberId/reject", s.rejectApplicant)
}

// memberDTO is the only db.Member -> api.Member mapping.
func memberDTO(m db.Member) api.Member {
	return api.Member{
		ID:       m.ID,
		Linked:   m.LineUserID != nil,
		Status:   api.MemberStatus(m.Status),
		Role:     api.Role(m.Role),
		Cohort:   m.Cohort,
		Nickname: m.Nickname,
		RealName: m.RealName,
		Address:  m.Address,
		Phone:    m.Phone,
		Deleted:  m.DeletedAt != nil,
	}
}

func (s *Server) listMembers(c *gin.Context) {
	status := c.Query("status")
	if status != "" && status != string(api.MemberApplicant) && status != string(api.MemberApproved) {
		fail(c, &domain.Error{Kind: domain.Invalid, Msg: "สถานะไม่ถูกต้อง"})
		return
	}
	rows, err := s.q.ListMembers(c.Request.Context(), status)
	if err != nil {
		fail(c, err)
		return
	}
	out := make([]api.Member, len(rows))
	for i, m := range rows {
		out[i] = memberDTO(m)
	}
	c.JSON(http.StatusOK, out)
}

func (s *Server) getMember(c *gin.Context) {
	id, ok := idParam(c, "memberId")
	if !ok {
		return
	}
	m, err := s.q.GetMember(c.Request.Context(), id)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, memberDTO(m))
}

// approveApplicant: applicant -> approved. Approving an approved Member is a no-op 200.
func (s *Server) approveApplicant(c *gin.Context) {
	id, ok := idParam(c, "memberId")
	if !ok {
		return
	}
	m, err := s.q.ApproveMember(c.Request.Context(), id)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, memberDTO(m))
}

// rejectApplicant hard-deletes the Applicant row (Applicants own nothing). 409 for non-applicants.
func (s *Server) rejectApplicant(c *gin.Context) {
	id, ok := idParam(c, "memberId")
	if !ok {
		return
	}
	ctx := c.Request.Context()
	err := s.tx(ctx, func(q *db.Queries) error {
		m, err := q.GetMember(ctx, id)
		if err != nil {
			return err
		}
		if m.Status != string(api.MemberApplicant) {
			return &domain.Error{Kind: domain.Conflict, Msg: "สมาชิกคนนี้อนุมัติแล้ว ปฏิเสธไม่ได้"}
		}
		return q.DeleteApplicant(ctx, id)
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
