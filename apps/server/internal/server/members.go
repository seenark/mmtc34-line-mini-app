package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

func (s *Server) mountMembers(r routes) {
	r.admin.GET("/members", s.listMembers)
	r.admin.GET("/members/:memberId", s.getMember)
}

// memberDTO is the only db.Member -> api.Member mapping.
func memberDTO(m db.Member) api.Member {
	return api.Member{
		ID:         m.ID,
		Linked:     m.LineUserID != nil,
		Status:     api.MemberStatus(m.Status),
		Role:       api.Role(m.Role),
		Cohort:     m.Cohort,
		Nickname:   m.Nickname,
		RealName:   m.RealName,
		Address:    m.Address,
		Phone:      m.Phone,
		Deleted:    m.DeletedAt != nil,
		SuperAdmin: isSuperAdmin(m),
	}
}

func isSuperAdmin(m db.Member) bool {
	return m.LineUserID != nil && *m.LineUserID == domain.SuperAdminLineUserID
}

func (s *Server) listMembers(c *gin.Context) {
	rows, err := s.q.ListMembers(c.Request.Context())
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
