package server

import "github.com/gin-gonic/gin"

// Ticket .10. Queries: db/queries/sponsors.sql.
func (s *Server) mountSponsors(r routes) {
	r.member.GET("/campaigns/:campaignId/sponsors", s.listSponsors) // loadCampaign gates drafts
	r.admin.POST("/campaigns/:campaignId/sponsors", s.createSponsor)
	r.admin.PUT("/sponsors/:sponsorId", s.updateSponsor) // includes sortOrder (reorder)
	r.admin.DELETE("/sponsors/:sponsorId", s.deleteSponsor)
}

func (s *Server) listSponsors(c *gin.Context)  { panic("not implemented") }
func (s *Server) createSponsor(c *gin.Context) { panic("not implemented") }
func (s *Server) updateSponsor(c *gin.Context) { panic("not implemented") }
func (s *Server) deleteSponsor(c *gin.Context) { panic("not implemented") }
