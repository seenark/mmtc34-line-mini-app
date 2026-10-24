package server

import (
	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

// Ticket .5. Queries: db/queries/payments.sql.
func (s *Server) mountPayments(r routes) {
	r.member.POST("/orders/:orderId/slip", s.uploadSlip) // multipart "file"
	r.member.GET("/slips/:name", s.readSlip)
	r.admin.GET("/payments/pending", s.pendingSlips)
	// By Order, not by Payment: payments_one_pending guarantees one pending Payment per Order.
	r.admin.POST("/orders/:orderId/confirm", s.confirmSlip) // api.AmountInput
	r.admin.POST("/orders/:orderId/reject", s.rejectSlip)   // api.RejectInput
}

// uploadSlip: s.apply(c, EvUploadSlip, fx) where fx saves the file (saveUpload kindSlips)
// and inserts a pending Payment. Saving inside fx means nobody without permission writes to disk.
func (s *Server) uploadSlip(c *gin.Context) {
	s.apply(c, domain.EvUploadSlip, func(q *db.Queries, o domain.Order) error {
		panic("not implemented")
	})
}

// readSlip: Slip owner (via payments -> orders.member_id) or Admin; anyone else 404.
func (s *Server) readSlip(c *gin.Context) { panic("not implemented") }

// pendingSlips: loadOrders(all) filtered to Status() == pending_review, oldest first,
// as []api.Order (each carries Slip name, Total, nickname, Campaign title). The Admin
// badge count is len() of the same response.
func (s *Server) pendingSlips(c *gin.Context) { panic("not implemented") }

// confirmSlip: amount > 0 required (the web prefills Order Total); fx marks the pending
// Payment confirmed with amount, reviewed_by, reviewed_at = now.
func (s *Server) confirmSlip(c *gin.Context) { panic("not implemented") }

// rejectSlip: trimmed reason required -> 400 "กรุณาระบุเหตุผล"; fx marks the pending Payment rejected.
func (s *Server) rejectSlip(c *gin.Context) { panic("not implemented") }
