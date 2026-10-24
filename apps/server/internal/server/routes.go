package server

import "github.com/gin-gonic/gin"

// routes are the four access levels (spec "Auth"). A feature file only attaches handlers
// to the group whose gate it needs; gating itself lives here and in session.go.
type routes struct {
	public *gin.RouterGroup // anyone: config, login
	self   *gin.RouterGroup // valid session, Member may be nil or an Applicant: me, onboarding, profile
	member *gin.RouterGroup // approved, not deleted
	admin  *gin.RouterGroup // approved Admin; mounted at /api/admin
}

// mount is the complete route table.
func (s *Server) mount(api *gin.RouterGroup) {
	r := routes{
		public: api,
		self:   api.Group("", s.requireSession),
		member: api.Group("", s.requireSession, s.requireApproved),
		admin:  api.Group("/admin", s.requireSession, s.requireAdmin),
	}
	s.mountConfig(r)
	s.mountAuth(r)
	s.mountMembers(r)
	s.mountCampaigns(r)
	s.mountUploads(r)
	s.mountOrders(r)
	s.mountSponsors(r)
	s.mountPayments(r)
	s.mountAdminOrders(r)
	s.mountUnlinked(r)
	s.mountProfile(r)
	s.mountLedger(r)
	s.mountExports(r)
	s.mountRoster(r)
}
