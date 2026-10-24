package server

import "github.com/gin-gonic/gin"

// Ticket .11. Queries: db/queries/profile.sql.
func (s *Server) mountProfile(r routes) {
	r.self.PUT("/me", s.editMe)
	r.self.DELETE("/me", s.deleteMe)
	r.admin.PUT("/members/:memberId", s.adminEditMember)
	r.admin.PUT("/members/:memberId/role", s.setRole)
}

// editMe: profile only; Orders keep their shipping snapshot.
func (s *Server) editMe(c *gin.Context) { panic("not implemented") }

// deleteMe, one transaction: an Applicant row is hard-deleted (it owns nothing); an approved
// Member is erased: cohort/real_name/address/phone = '', line_user_id = NULL, deleted_at = now,
// and ship_* = '' on all their Orders. Nickname, Orders, Payments, Adjustments stay. Clears the cookie.
func (s *Server) deleteMe(c *gin.Context) { panic("not implemented") }

func (s *Server) adminEditMember(c *gin.Context) { panic("not implemented") }

// setRole: an Admin may not demote themself (409), so the app never loses its last active Admin
// through this route.
func (s *Server) setRole(c *gin.Context) { panic("not implemented") }
