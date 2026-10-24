package server

import "github.com/gin-gonic/gin"

// Ticket .8. No query file: reads go through loadOrders.
func (s *Server) mountExports(r routes) {
	r.admin.GET("/campaigns/:campaignId/exports/:kind", s.export) // kind: factory.csv|shipping.csv|pickup.csv
}

// export: loadOrders(CampaignID), keep Status() == paid, plus pending_review when
// ?includePending=1. Writes UTF-8 BOM + encoding/csv with
// Content-Disposition: attachment; filename*=UTF-8''<title>-<kind>.
//   factory:  Option, size, print name, number, quantity per Item; blank row; count per Option+size.
//   shipping: ship Orders: name, address, phone, Items summary.
//   pickup:   pickup Orders: nickname, real name, Items summary.
func (s *Server) export(c *gin.Context) { panic("not implemented") }
