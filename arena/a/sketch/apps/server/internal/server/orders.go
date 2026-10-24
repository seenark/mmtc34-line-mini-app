package server

import (
	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

// Ticket .4. Queries: db/queries/orders.sql.
func (s *Server) mountOrders(r routes) {
	r.member.POST("/campaigns/:campaignId/orders", s.createOrder)
	r.member.GET("/orders", s.myOrders)
	r.member.GET("/orders/:orderId", s.getOrder)
	// One cancel endpoint for Member and Admin: the state machine decides by actor.
	r.member.POST("/orders/:orderId/cancel", s.cancelOrder)
}

// createOrder: bind api.OrderInput; s.checkout(campaignID, actor.MemberID, nil, in).
func (s *Server) createOrder(c *gin.Context) { panic("not implemented") }

// myOrders: loadOrders(MemberID = me), newest first.
func (s *Server) myOrders(c *gin.Context) { panic("not implemented") }

// getOrder: s.apply(c, domain.EvView, nil) — EvView is how "only owner or Admin" is enforced.
func (s *Server) getOrder(c *gin.Context) { s.apply(c, domain.EvView, nil) }

// cancelOrder: Admin cancel also rejects a pending Slip ("ยกเลิกคำสั่งซื้อ") so the queue never
// shows a Slip for a cancelled Order.
func (s *Server) cancelOrder(c *gin.Context) { panic("not implemented") }
