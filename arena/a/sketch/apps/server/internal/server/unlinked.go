package server

import "github.com/gin-gonic/gin"

// Ticket .7 (ADR 0002). Queries: db/queries/unlinked.sql.
func (s *Server) mountUnlinked(r routes) {
	r.admin.POST("/members", s.createUnlinked)             // api.ProfileInput, nickname only
	r.admin.POST("/members/:memberId/link", s.linkApplicant) // api.LinkInput
}

// createUnlinked: approved Member with line_user_id NULL.
func (s *Server) createUnlinked(c *gin.Context) { panic("not implemented") }

// linkApplicant, one transaction: target must be Unlinked (line_user_id NULL, not deleted),
// source must be an Applicant; delete the Applicant first (frees the UNIQUE userId), then copy its
// userId and profile onto the target. The person's next request resolves to the target Member.
func (s *Server) linkApplicant(c *gin.Context) { panic("not implemented") }
