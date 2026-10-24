package server

import (
	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
)

// Ticket .3. Queries: db/queries/campaigns.sql, db/queries/options.sql.
func (s *Server) mountCampaigns(r routes) {
	r.member.GET("/campaigns", s.listCampaigns)             // non-draft, effective status
	r.member.GET("/campaigns/:campaignId", s.campaignDetail) // draft -> 404 for Members
	r.admin.GET("/campaigns", s.adminListCampaigns)          // includes drafts
	r.admin.POST("/campaigns", s.createCampaign)
	r.admin.PUT("/campaigns/:campaignId", s.updateCampaign) // also open / close / reopen
	r.admin.POST("/campaigns/:campaignId/options", s.createOption)
	r.admin.PUT("/options/:optionId", s.updateOption) // includes sortOrder (reorder)
	r.admin.DELETE("/options/:optionId", s.deleteOption) // FK RESTRICT -> 409 via fail()
}

func (s *Server) listCampaigns(c *gin.Context)      { panic("not implemented") }
func (s *Server) campaignDetail(c *gin.Context)     { panic("not implemented") }
func (s *Server) adminListCampaigns(c *gin.Context) { panic("not implemented") }
func (s *Server) createCampaign(c *gin.Context)     { panic("not implemented") }
func (s *Server) updateCampaign(c *gin.Context)     { panic("not implemented") }
func (s *Server) createOption(c *gin.Context)       { panic("not implemented") }
func (s *Server) updateOption(c *gin.Context)       { panic("not implemented") }
func (s *Server) deleteOption(c *gin.Context)       { panic("not implemented") }

// toAPICampaign is the only db.Campaign -> api.Campaign mapping; it writes the effective
// status via domain.Campaign.Status(s.now()).
func (s *Server) toAPICampaign(row db.Campaign) api.Campaign { panic("not implemented") }

// loadCampaign returns the Campaign if the actor may see it (domain.Campaign.VisibleTo) else 404.
// Used by orders (.4), sponsors (.10), roster (.9), exports (.8).
func (s *Server) loadCampaign(c *gin.Context, q *db.Queries, id int64) (db.Campaign, error) {
	panic("not implemented")
}
