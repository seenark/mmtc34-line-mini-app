package server_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

func slipPath(orderID int64) string    { return fmt.Sprintf("/api/orders/%d/slip", orderID) }
func confirmPath(orderID int64) string { return fmt.Sprintf("/api/admin/orders/%d/confirm", orderID) }
func rejectPath(orderID int64) string  { return fmt.Sprintf("/api/admin/orders/%d/reject", orderID) }

const pendingPath = "/api/admin/payments/pending"

func TestUploadMovesOrderToPendingReview(t *testing.T) {
	h := New(t)
	alice := h.Member("alice")
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2)))

	got := alice.UploadSlip(o.ID)

	Equal(t, got.Status, api.OrderPendingReview)
	Equal(t, len(got.Payments), 1)
	p := got.Payments[0]
	Equal(t, p.Status, api.PaymentPending)
	Equal(t, len(p.Slip), 32+len(".jpg"))
	Equal(t, p.Amount == nil, true)
	Equal(t, p.ReviewedBy, "")
	Equal(t, got.Actions.UploadSlip, false)
	Equal(t, alice.Order(o.ID).Status, api.OrderPendingReview)
	alice.POST(slipPath(o.ID), JPEG()).Status(409)
}

func TestUploadRejectsBadFileAndLeavesOrderAwaiting(t *testing.T) {
	h := New(t)
	alice := h.Member("alice")
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 1)))

	alice.POST(slipPath(o.ID), nil).Status(400)
	alice.POST(slipPath(o.ID), Upload{Filename: "slip.jpg", Bytes: []byte("not a picture")}).Status(400)
	alice.POST(slipPath(o.ID), Upload{Filename: "slip.jpg", Bytes: bytes.Repeat([]byte{0xff}, 6<<20)}).Status(413)

	got := alice.Order(o.ID)
	Equal(t, got.Status, api.OrderAwaitingPayment)
	Equal(t, len(got.Payments), 0)
}

func TestOnlyTheOwnerUploadsASlip(t *testing.T) {
	h := New(t)
	alice, bob, admin := h.Member("alice"), h.Member("bob"), h.Admin()
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 1)))

	bob.POST(slipPath(o.ID), JPEG()).Status(404)
	admin.POST(slipPath(o.ID), JPEG()).Status(403)
	h.Anonymous().POST(slipPath(o.ID), JPEG()).Status(401)
	Equal(t, alice.Order(o.ID).Status, api.OrderAwaitingPayment)
}

func TestRejectNeedsAReasonAndReturnsOrderToAwaiting(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2)))
	alice.UploadSlip(o.ID)

	admin.POST(rejectPath(o.ID), api.RejectInput{}).Status(400)
	Equal(t, admin.POST(rejectPath(o.ID), api.RejectInput{Reason: "  "}).Status(400).Error(), "กรุณาระบุเหตุผล")
	Equal(t, alice.Order(o.ID).Status, api.OrderPendingReview)

	got := As[api.Order](admin.POST(rejectPath(o.ID), api.RejectInput{Reason: "ยอดไม่ตรง"}).Status(200))
	Equal(t, got.Status, api.OrderAwaitingPayment)

	mine := alice.Order(o.ID)
	Equal(t, mine.Status, api.OrderAwaitingPayment)
	Equal(t, mine.Actions.UploadSlip, true)
	Equal(t, len(mine.Payments), 1)
	p := mine.Payments[0]
	Equal(t, p.Status, api.PaymentRejected)
	Equal(t, p.RejectReason, "ยอดไม่ตรง")
	Equal(t, p.ReviewedBy, "admin")
	Equal(t, p.ReviewedAt != nil, true)
	Equal(t, mine.Paid, int64(0))
	admin.POST(rejectPath(o.ID), api.RejectInput{Reason: "ซ้ำ"}).Status(409)

	again := alice.UploadSlip(o.ID)
	Equal(t, again.Status, api.OrderPendingReview)
	Equal(t, len(again.Payments), 2)
	Equal(t, again.Payments[0].Status, api.PaymentRejected)
	Equal(t, again.Payments[1].Status, api.PaymentPending)
	Equal(t, again.Payments[0].Slip != again.Payments[1].Slip, true)
}

func TestConfirmPaysTheOrderAndRecordsTheReviewer(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2)))
	alice.UploadSlip(o.ID)

	admin.POST(confirmPath(o.ID), api.AmountInput{Amount: 0}).Status(400)
	admin.POST(confirmPath(o.ID), api.AmountInput{Amount: -5}).Status(400)
	Equal(t, alice.Order(o.ID).Status, api.OrderPendingReview)

	got := As[api.Order](admin.POST(confirmPath(o.ID), api.AmountInput{Amount: 500}).Status(200))

	Equal(t, got.Status, api.OrderPaid)
	Equal(t, got.Total, int64(520))
	Equal(t, got.Paid, int64(500))
	Equal(t, got.Balance, int64(20))
	mine := alice.Order(o.ID)
	Equal(t, mine.Status, api.OrderPaid)
	p := mine.Payments[0]
	Equal(t, p.Status, api.PaymentConfirmed)
	Equal(t, *p.Amount, int64(500))
	Equal(t, p.ReviewedBy, "admin")
	Equal(t, p.ReviewedAt != nil, true)
	admin.POST(confirmPath(o.ID), api.AmountInput{Amount: 500}).Status(409)
	admin.POST(rejectPath(o.ID), api.RejectInput{Reason: "x"}).Status(409)
	Equal(t, alice.Order(o.ID).Paid, int64(500))
}

func TestSlipIsReadableOnlyByOwnerAndAdmin(t *testing.T) {
	h := New(t)
	alice, bob, admin := h.Member("alice"), h.Member("bob"), h.Admin()
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 1)))
	name := alice.UploadSlip(o.ID).Payments[0].Slip
	path := "/api/slips/" + name

	r := alice.GET(path).Status(200)
	Equal(t, r.Header.Get("X-Content-Type-Options"), "nosniff")
	Equal(t, r.Header.Get("Cache-Control"), "private, max-age=31536000, immutable")
	Equal(t, bytes.Equal(r.Body, JPEG().Bytes), true)
	Equal(t, bytes.Equal(admin.GET(path).Status(200).Body, JPEG().Bytes), true)
	bob.GET(path).Status(404)
	h.Anonymous().GET(path).Status(401)
	h.Applicant("nok").GET(path).Status(403)
	alice.GET("/api/slips/0123456789abcdef0123456789abcdef.jpg").Status(404)
	alice.GET("/api/slips/..%2Fapp.db").Status(404)
}

func TestUploadAfterCloseTimeIs409AndStoresNothing(t *testing.T) {
	h := New(t)
	alice := h.Member("alice")
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 1)))

	h.Advance(25 * time.Hour)

	alice.POST(slipPath(o.ID), JPEG()).Status(409)
	Equal(t, len(alice.Order(o.ID).Payments), 0)
	Equal(t, len(As[[]api.Order](h.Admin().GET(pendingPath).Status(200))), 0)
}

func TestPendingQueueListsOnlyPendingReviewOldestFirst(t *testing.T) {
	h := New(t)
	alice, bob, cat, admin := h.Member("alice"), h.Member("bob"), h.Member("cat"), h.Admin()
	j := h.OpenJersey()
	alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "S", 1)))
	first := bob.Checkout(j.Campaign.ID, Pickup(Shirt(j, "M", 1)))
	second := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 1)))
	done := cat.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 1)))

	alice.UploadSlip(second.ID)
	h.Advance(-time.Minute)
	bob.UploadSlip(first.ID)
	h.Advance(2 * time.Minute)
	cat.UploadSlip(done.ID)
	admin.POST(confirmPath(done.ID), api.AmountInput{Amount: 260}).Status(200)

	queue := As[[]api.Order](admin.GET(pendingPath).Status(200))

	Equal(t, len(queue), 2)
	Equal(t, queue[0].ID, first.ID)
	Equal(t, queue[1].ID, second.ID)
	Equal(t, queue[0].Actions.ReviewSlip, true)
	Equal(t, queue[0].MemberNickname, "bob")
	Equal(t, queue[0].Status, api.OrderPendingReview)
}

func TestMembersCannotUseAdminSlipEndpoints(t *testing.T) {
	h := New(t)
	alice := h.Member("alice")
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 1)))
	alice.UploadSlip(o.ID)

	alice.GET(pendingPath).Status(403)
	alice.POST(confirmPath(o.ID), api.AmountInput{Amount: 260}).Status(403)
	alice.POST(rejectPath(o.ID), api.RejectInput{Reason: "x"}).Status(403)
	h.Anonymous().GET(pendingPath).Status(401)
	Equal(t, alice.Order(o.ID).Status, api.OrderPendingReview)
}
