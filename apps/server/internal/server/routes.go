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

// mount is the complete route table. Ticket .1 writes every mount line and an empty mountX in
// each feature file; .2 adds the session middleware to the group lines. Wave tickets never edit this file.
func (s *Server) mount(api *gin.RouterGroup) {
	r := routes{
		public: api,
		self:   api.Group(""),
		member: api.Group(""),
		admin:  api.Group("/admin"),
	}
	s.mountConfig(r)      // .1  server.go
	s.mountAuth(r)        // .2  auth.go
	s.mountMembers(r)     // .2  members.go
	s.mountCampaigns(r)   // .3  campaigns.go
	s.mountUploads(r)     // .3  uploads.go
	s.mountOrders(r)      // .4  orders.go
	s.mountSponsors(r)    // .10 sponsors.go
	s.mountPayments(r)    // .5  payments.go
	s.mountAdminOrders(r) // .7  admin_orders.go
	s.mountUnlinked(r)    // .7  unlinked.go
	s.mountProfile(r)     // .11 profile.go
	s.mountLedger(r)      // .6  ledger.go
	s.mountExports(r)     // .8  exports.go
	s.mountRoster(r)      // .9  roster.go
}
