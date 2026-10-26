package server

import (
	"database/sql"
	"net/http"
	"slices"
	"strconv"
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
	r.admin.POST("/orders/:orderId/payments/:paymentId/slips", s.addSlipImages)
	r.admin.DELETE("/orders/:orderId/payments/:paymentId/slips/:name", s.removeSlipImage)
}

// uploadSlip reads the body before apply opens the write transaction, so a slow upload never
// holds the SQLite write lock. Files whose Decide or transaction then fails are deleted.
func (s *Server) uploadSlip(c *gin.Context) {
	names, err := s.saveUploads(c, kindSlips, maxSlips)
	if err != nil {
		fail(c, err)
		return
	}
	ok := s.apply(c, domain.EvUploadSlip, func(q *db.Queries, o domain.Order) error {
		now := rfc3339(s.now())
		id, err := q.InsertPendingPayment(c.Request.Context(), db.InsertPendingPaymentParams{OrderID: o.ID, Now: now})
		if err != nil {
			return err
		}
		return s.insertSlipImages(c, q, id, names, now)
	})
	if !ok {
		s.removeFiles(kindSlips, names)
	}
}

func (s *Server) insertSlipImages(c *gin.Context, q *db.Queries, paymentID int64, names []string, now string) error {
	for _, n := range names {
		if err := q.InsertSlipImage(c.Request.Context(), db.InsertSlipImageParams{PaymentID: paymentID, Name: n, Now: now}); err != nil {
			return err
		}
	}
	return nil
}

// slipPayment is the Order's Payment named by :paymentId, which must be a Slip Payment (not cash).
func slipPayment(c *gin.Context, o domain.Order) (domain.Payment, error) {
	id, err := strconv.ParseInt(c.Param("paymentId"), 10, 64)
	if err == nil {
		for _, p := range o.Payments {
			if p.ID == id {
				if len(p.Slips) == 0 {
					return p, &domain.Error{Kind: domain.Conflict, Msg: "รายการนี้เป็นเงินสด ไม่มีสลิป"}
				}
				return p, nil
			}
		}
	}
	return domain.Payment{}, &domain.Error{Kind: domain.NotFound, Msg: "ไม่พบข้อมูล"}
}

// addSlipImages lets an Admin attach more images to a Slip Payment in any status.
func (s *Server) addSlipImages(c *gin.Context) {
	names, err := s.saveUploads(c, kindSlips, maxSlips)
	if err != nil {
		fail(c, err)
		return
	}
	ok := s.apply(c, domain.EvEditSlips, func(q *db.Queries, o domain.Order) error {
		p, err := slipPayment(c, o)
		if err != nil {
			return err
		}
		return s.insertSlipImages(c, q, p.ID, names, rfc3339(s.now()))
	})
	if !ok {
		s.removeFiles(kindSlips, names)
	}
}

// removeSlipImage drops one image of a Slip Payment; the last one stays, so a Slip Payment never
// turns into cash. The file goes only after the transaction commits.
func (s *Server) removeSlipImage(c *gin.Context) {
	name := c.Param("name")
	ok := s.apply(c, domain.EvEditSlips, func(q *db.Queries, o domain.Order) error {
		p, err := slipPayment(c, o)
		if err != nil {
			return err
		}
		if !slices.Contains(p.Slips, name) {
			return &domain.Error{Kind: domain.NotFound, Msg: "ไม่พบข้อมูล"}
		}
		if len(p.Slips) == 1 {
			return &domain.Error{Kind: domain.Conflict, Msg: "ต้องมีรูปสลิปอย่างน้อย 1 รูป"}
		}
		_, err = q.DeleteSlipImage(c.Request.Context(), db.DeleteSlipImageParams{PaymentID: p.ID, Name: name})
		return err
	})
	if ok {
		s.removeFiles(kindSlips, []string{name})
	}
}

// readSlip: only the owning Member or an Admin; anyone else gets the same 404 as a missing Slip.
func (s *Server) readSlip(c *gin.Context) {
	name := c.Param("name")
	if !fileName.MatchString(name) {
		fail(c, sql.ErrNoRows)
		return
	}
	owner, err := s.q.SlipOwner(c.Request.Context(), name)
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
