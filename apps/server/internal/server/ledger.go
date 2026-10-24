package server

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

func (s *Server) mountLedger(r routes) {
	r.admin.GET("/campaigns/:campaignId/orders", s.campaignOrders)
	r.admin.POST("/orders/:orderId/cash", s.recordCash)
	r.admin.POST("/orders/:orderId/adjustments", s.addAdjustment)
}

// campaignOrders lists every Order of a Campaign (drafts too), newest first. The status filter
// is applied in Go on the effective status.
func (s *Server) campaignOrders(c *gin.Context) {
	id, ok := idParam(c, "campaignId")
	if !ok {
		return
	}
	status := api.OrderStatus(c.Query("status"))
	if status != "" && !slices.Contains([]api.OrderStatus{
		api.OrderAwaitingPayment, api.OrderPendingReview, api.OrderPaid, api.OrderCancelled,
	}, status) {
		fail(c, &domain.Error{Kind: domain.Invalid, Msg: "สถานะไม่ถูกต้อง"})
		return
	}
	ctx := c.Request.Context()
	if _, err := s.loadCampaign(c, s.q, id); err != nil {
		fail(c, err)
		return
	}
	orders, err := s.loadOrders(ctx, s.q, orderScope{CampaignID: &id})
	if err != nil {
		fail(c, err)
		return
	}
	a := actorOf(c)
	out := make([]api.Order, 0, len(orders))
	for _, o := range orders {
		if status == "" || o.Status() == status {
			out = append(out, toAPIOrder(o, a))
		}
	}
	c.JSON(http.StatusOK, out)
}

// recordCash settles cash received outside the Slip flow: a confirmed Payment with no Slip.
func (s *Server) recordCash(c *gin.Context) {
	in, ok := bind[api.AmountInput](c)
	if !ok {
		return
	}
	if in.Amount <= 0 {
		fail(c, &domain.Error{Kind: domain.Invalid, Msg: "กรุณาระบุจำนวนเงินที่ถูกต้อง"})
		return
	}
	me := actorOf(c)
	s.apply(c, domain.EvRecordCash, func(q *db.Queries, o domain.Order) error {
		return q.InsertCashPayment(c.Request.Context(), db.InsertCashPaymentParams{
			OrderID: o.ID, Amount: &in.Amount, Now: rfc3339(s.now()), ReviewedBy: &me.MemberID,
		})
	})
}

// addAdjustment: positive amounts raise what the Member owes, negative ones lower it.
func (s *Server) addAdjustment(c *gin.Context) {
	in, ok := bind[api.AdjustmentInput](c)
	if !ok {
		return
	}
	note := strings.TrimSpace(in.Note)
	if in.Amount == 0 {
		fail(c, &domain.Error{Kind: domain.Invalid, Msg: "กรุณาระบุจำนวนเงินที่ถูกต้อง"})
		return
	}
	if note == "" {
		fail(c, &domain.Error{Kind: domain.Invalid, Msg: "กรุณาระบุหมายเหตุ"})
		return
	}
	me := actorOf(c)
	s.apply(c, domain.EvAdjust, func(q *db.Queries, o domain.Order) error {
		return q.InsertAdjustment(c.Request.Context(), db.InsertAdjustmentParams{
			OrderID: o.ID, Amount: in.Amount, Note: note, CreatedBy: me.MemberID, Now: rfc3339(s.now()),
		})
	})
}
