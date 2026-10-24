package server

import "github.com/gin-gonic/gin"

// Ticket .9. No query file: reads go through loadOrders.
func (s *Server) mountRoster(r routes) {
	r.member.GET("/campaigns/:campaignId/roster", s.roster)
}

// roster: loadOrders(CampaignID), keep Status() in {pending_review, paid}, group by Member in
// first-order sequence, emit []api.RosterEntry. The DTO has no field for contact data or status,
// so leaking them would need a type change, not a forgotten filter.
func (s *Server) roster(c *gin.Context) { panic("not implemented") }
