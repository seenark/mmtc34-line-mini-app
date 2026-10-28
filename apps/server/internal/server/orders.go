package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

func (s *Server) mountOrders(r routes) {
	r.member.POST("/campaigns/:campaignId/orders", s.createOrder)
	r.member.GET("/orders", s.myOrders)
	r.member.GET("/orders/:orderId", s.getOrder)
	// One cancel endpoint for Member and Admin: the state machine decides by actor.
	r.member.POST("/orders/:orderId/cancel", s.cancelOrder)
	r.member.PUT("/orders/:orderId/details", s.editOrderDetails)
}

func (s *Server) createOrder(c *gin.Context) {
	campaignID, ok := idParam(c, "campaignId")
	if !ok {
		return
	}
	in, ok := bind[api.OrderInput](c)
	if !ok {
		return
	}
	s.checkout(c, campaignID, actorOf(c).MemberID, false, in)
}

func (s *Server) myOrders(c *gin.Context) {
	a := actorOf(c)
	orders, err := s.loadOrders(c.Request.Context(), s.q, orderScope{MemberID: &a.MemberID})
	if err != nil {
		fail(c, err)
		return
	}
	out := make([]api.Order, len(orders))
	for i, o := range orders {
		out[i] = toAPIOrder(o, a)
	}
	c.JSON(http.StatusOK, out)
}

// getOrder: EvView is how "only the owner or an Admin" is enforced (others get 404).
func (s *Server) getOrder(c *gin.Context) { s.apply(c, domain.EvView, nil) }

// cancelOrder: an Admin cancel also rejects a pending Slip, so the queue never shows a Slip
// for a cancelled Order.
func (s *Server) cancelOrder(c *gin.Context) {
	me := actorOf(c)
	s.apply(c, domain.EvCancel, func(q *db.Queries, o domain.Order) error {
		now := rfc3339(s.now())
		return q.RejectPendingPayment(c.Request.Context(), db.RejectPendingPaymentParams{
			Reason: "ยกเลิกคำสั่งซื้อ", ReviewedBy: &me.MemberID, Now: &now, OrderID: o.ID,
		})
	})
}

// editOrderDetails lets the owner correct size, print name, number and shipping while the
// Campaign is open, whatever the payment state. Money fields cannot change (ValidateDetails).
func (s *Server) editOrderDetails(c *gin.Context) {
	in, ok := bind[api.OrderDetailsInput](c)
	if !ok {
		return
	}
	s.apply(c, domain.EvEditDetails, func(q *db.Queries, o domain.Order) error {
		ctx := c.Request.Context()
		options, err := s.campaignOptions(ctx, q, o.CampaignID)
		if err != nil {
			return err
		}
		v, err := domain.ValidateDetails(options, o, in)
		if err != nil {
			return err
		}
		return s.replaceContents(ctx, q, o.ID, v)
	})
}
