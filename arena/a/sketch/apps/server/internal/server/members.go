package server

import "github.com/gin-gonic/gin"

// Ticket .2. Queries: db/queries/members.sql.
func (s *Server) mountMembers(r routes) {
	r.admin.GET("/members", s.listMembers) // ?status=applicant|approved
	r.admin.GET("/members/:memberId", s.getMember)
	r.admin.POST("/members/:memberId/approve", s.approveApplicant)
	r.admin.POST("/members/:memberId/reject", s.rejectApplicant)
}

func (s *Server) listMembers(c *gin.Context) { panic("not implemented") }
func (s *Server) getMember(c *gin.Context)   { panic("not implemented") }

// approveApplicant: applicant -> approved. Approving an approved Member is a no-op 200.
func (s *Server) approveApplicant(c *gin.Context) { panic("not implemented") }

// rejectApplicant hard-deletes the Applicant row (Applicants own nothing). 409 for non-applicants.
func (s *Server) rejectApplicant(c *gin.Context) { panic("not implemented") }
