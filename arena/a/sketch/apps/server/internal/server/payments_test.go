package server_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

// Ticket .5: upload -> reject -> re-upload -> confirm with under-transfer.
func TestSlipRejectReuploadConfirmUnderTransfer(t *testing.T) {
	h := New(t)
	admin, alice := h.Admin(), h.Member("alice")
	j := h.OpenJersey()

	order := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2))) // 520
	order = alice.UploadSlip(order.ID)
	Equal(t, order.Status, api.OrderPendingReview)

	reject := fmt.Sprintf("/api/admin/orders/%d/reject", order.ID)
	Equal(t, admin.POST(reject, api.RejectInput{}).Status(400).Error(), "กรุณาระบุเหตุผล")
	admin.POST(reject, api.RejectInput{Reason: "สลิปไม่ชัด"}).Status(200)

	order = alice.Order(order.ID)
	Equal(t, order.Status, api.OrderAwaitingPayment)
	Equal(t, order.Payments[0].RejectReason, "สลิปไม่ชัด")
	Equal(t, order.Actions.UploadSlip, true)

	order = alice.UploadSlip(order.ID)
	confirm := fmt.Sprintf("/api/admin/orders/%d/confirm", order.ID)
	order = As[api.Order](admin.POST(confirm, api.AmountInput{Amount: 500}).Status(200))
	Equal(t, order.Status, api.OrderPaid)
	Equal(t, order.Balance, int64(20))

	bob := h.Member("bob")
	bob.GET("/api/slips/" + order.Payments[1].Slip).Status(http.StatusNotFound)
	alice.GET("/api/slips/" + order.Payments[1].Slip).Status(200)
}

// Ticket .5: a Slip cannot be uploaded once the close time passes; the Order reads cancelled.
func TestSlipRefusedAfterCloseTime(t *testing.T) {
	h := New(t)
	alice := h.Member("alice")
	j := h.OpenJersey() // closeAt = start + 24h

	order := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "M", 1)))
	h.Advance(25 * time.Hour)

	res := alice.POST(fmt.Sprintf("/api/orders/%d/slip", order.ID), JPEG()).Status(409)
	Equal(t, res.Error(), "แคมเปญปิดแล้ว กรุณาติดต่อแอดมิน")
	Equal(t, alice.Order(order.ID).Status, api.OrderCancelled)
	Equal(t, alice.Order(order.ID).Balance, int64(0))
}
