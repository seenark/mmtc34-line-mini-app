package server

import "github.com/gin-gonic/gin"

// Ticket .6. Queries: db/queries/ledger.sql.
func (s *Server) mountLedger(r routes) {
	r.admin.GET("/campaigns/:campaignId/orders", s.campaignOrders) // ?status=<effective status>
	r.admin.POST("/orders/:orderId/cash", s.recordCash)           // api.AmountInput
	r.admin.POST("/orders/:orderId/adjustments", s.addAdjustment) // api.AdjustmentInput
}

// campaignOrders: loadOrders(CampaignID) filtered by o.Status() when ?status is given.
func (s *Server) campaignOrders(c *gin.Context) { panic("not implemented") }

// recordCash: apply(EvRecordCash) with fx inserting a confirmed Payment (slip NULL, reviewed_by = actor).
func (s *Server) recordCash(c *gin.Context) { panic("not implemented") }

// addAdjustment: amount != 0, note required; apply(EvAdjust) with fx inserting the Adjustment.
func (s *Server) addAdjustment(c *gin.Context) { panic("not implemented") }
