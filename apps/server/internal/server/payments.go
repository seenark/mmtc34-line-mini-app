package server

import (
	"database/sql"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

func (s *Server) mountPayments(r routes) {
	r.member.POST("/orders/:orderId/slip", s.uploadSlip)
	r.member.GET("/slips/:name", s.readSlip)
	r.admin.GET("/payments/pending", s.pendingSlips)
	r.admin.POST("/orders/:orderId/confirm", s.confirmSlip)
	r.admin.POST("/orders/:orderId/reject", s.rejectSlip)
}

// uploadSlip saves the file inside the apply effect, i.e. after Decide, so a refused upload
// stores nothing. A file whose transaction then fails is deleted.
func (s *Server) uploadSlip(c *gin.Context) {
	var saved string
	ok := s.apply(c, domain.EvUploadSlip, func(q *db.Queries, o domain.Order) error {
		name, err := s.saveUpload(c, kindSlips)
		if err != nil {
			return err
		}
		saved = name
		return q.InsertPendingPayment(c.Request.Context(), db.InsertPendingPaymentParams{
			OrderID: o.ID, Slip: &name, Now: rfc3339(s.now()),
		})
	})
	if !ok && saved != "" {
		os.Remove(s.filePath(kindSlips, saved))
	}
}

// readSlip: only the owning Member or an Admin; anyone else gets the same 404 as a missing Slip.
func (s *Server) readSlip(c *gin.Context) {
	name := c.Param("name")
	if !fileName.MatchString(name) {
		fail(c, sql.ErrNoRows)
		return
	}
	owner, err := s.q.SlipOwner(c.Request.Context(), &name)
	if err != nil {
		fail(c, err)
		return
	}
	if a := actorOf(c); owner != a.MemberID && !a.Admin {
		fail(c, sql.ErrNoRows)
		return
	}
	c.Header("Cache-Control", "private, max-age=31536000, immutable") // names are random, never reused
	c.Header("X-Content-Type-Options", "nosniff")
	c.File(s.filePath(kindSlips, name))
}

// pendingSlips filters in Go on effective status, oldest pending Payment first.
func (s *Server) pendingSlips(c *gin.Context) {
	orders, err := s.loadOrders(c.Request.Context(), s.q, orderScope{})
	if err != nil {
		fail(c, err)
		return
	}
	orders = slices.DeleteFunc(orders, func(o domain.Order) bool { return o.Status() != api.OrderPendingReview })
	waiting := func(o domain.Order) string {
		for _, p := range o.Payments {
			if p.Status == api.PaymentPending {
				return p.CreatedAt
			}
		}
		return ""
	}
	slices.SortStableFunc(orders, func(a, b domain.Order) int { return strings.Compare(waiting(a), waiting(b)) })
	a := actorOf(c)
	out := make([]api.Order, len(orders))
	for i, o := range orders {
		out[i] = toAPIOrder(o, a)
	}
	c.JSON(http.StatusOK, out)
}

func (s *Server) confirmSlip(c *gin.Context) {
	in, ok := bind[api.AmountInput](c)
	if !ok {
		return
	}
	if in.Amount <= 0 {
		fail(c, &domain.Error{Kind: domain.Invalid, Msg: "กรุณาระบุจำนวนเงินที่ถูกต้อง"})
		return
	}
	me := actorOf(c)
	s.apply(c, domain.EvConfirmSlip, func(q *db.Queries, o domain.Order) error {
		now := rfc3339(s.now())
		return q.ConfirmPendingPayment(c.Request.Context(), db.ConfirmPendingPaymentParams{
			Amount: &in.Amount, ReviewedBy: &me.MemberID, Now: &now, OrderID: o.ID,
		})
	})
}

func (s *Server) rejectSlip(c *gin.Context) {
	in, ok := bind[api.RejectInput](c)
	if !ok {
		return
	}
	reason := strings.TrimSpace(in.Reason)
	if reason == "" {
		fail(c, &domain.Error{Kind: domain.Invalid, Msg: "กรุณาระบุเหตุผล"})
		return
	}
	me := actorOf(c)
	s.apply(c, domain.EvRejectSlip, func(q *db.Queries, o domain.Order) error {
		now := rfc3339(s.now())
		return q.RejectPendingPayment(c.Request.Context(), db.RejectPendingPaymentParams{
			Reason: reason, ReviewedBy: &me.MemberID, Now: &now, OrderID: o.ID,
		})
	})
}
