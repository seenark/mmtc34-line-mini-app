package server_test

import (
	"testing"
	"time"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

func TestCampaignOrdersListsNewestFirstAndFiltersOnEffectiveStatus(t *testing.T) {
	h := New(t)
	alice, bob, admin := h.Member("alice"), h.Member("bob"), h.Admin()
	j := h.OpenJersey()
	first := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2)))
	second := bob.Checkout(j.Campaign.ID, Pickup(Shirt(j, "M", 1)))
	bob.UploadSlip(second.ID)

	all := admin.CampaignOrders(j.Campaign.ID, "")
	Equal(t, len(all), 2)
	Equal(t, all[0].ID, second.ID)
	Equal(t, all[1].ID, first.ID)
	Equal(t, all[0].MemberNickname, "bob")
	Equal(t, admin.CampaignOrders(j.Campaign.ID, "pending_review")[0].ID, second.ID)
	Equal(t, len(admin.CampaignOrders(j.Campaign.ID, "awaiting_payment")), 1)
	Equal(t, len(admin.CampaignOrders(j.Campaign.ID, "paid")), 0)

	h.Advance(25 * time.Hour)

	Equal(t, len(admin.CampaignOrders(j.Campaign.ID, "awaiting_payment")), 0)
	gone := admin.CampaignOrders(j.Campaign.ID, "cancelled")
	Equal(t, len(gone), 1)
	Equal(t, gone[0].ID, first.ID)
	Equal(t, gone[0].Status, api.OrderCancelled)
	Equal(t, admin.CampaignOrders(j.Campaign.ID, "pending_review")[0].ID, second.ID)
}

func TestCampaignOrdersRejectsBadStatusMissingCampaignAndMembers(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()

	admin.GET(campaignOrdersPath(j.Campaign.ID) + "?status=shipped").Status(400)
	admin.GET(campaignOrdersPath(9999)).Status(404)
	alice.GET(campaignOrdersPath(j.Campaign.ID)).Status(403)
	h.Anonymous().GET(campaignOrdersPath(j.Campaign.ID)).Status(401)
}

func TestCampaignOrdersIncludesDraftCampaignsForAdmin(t *testing.T) {
	h := New(t)
	admin := h.Admin()
	d := As[api.Campaign](admin.POST("/api/admin/campaigns", draft("ร่าง")).Status(201))

	Equal(t, len(admin.CampaignOrders(d.ID, "")), 0)
}

func TestCashSettlesAnAwaitingOrderWithoutASlip(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2)))

	admin.POST(cashPath(o.ID), api.AmountInput{Amount: 0}).Status(400)
	admin.POST(cashPath(o.ID), api.AmountInput{Amount: -5}).Status(400)
	Equal(t, alice.Order(o.ID).Status, api.OrderAwaitingPayment)
	alice.POST(cashPath(o.ID), api.AmountInput{Amount: 520}).Status(403)

	got := admin.Cash(o.ID, 520)

	Equal(t, got.Status, api.OrderPaid)
	Equal(t, got.Total, int64(520))
	Equal(t, got.Paid, int64(520))
	Equal(t, got.Balance, int64(0))
	Equal(t, len(got.Payments), 1)
	p := got.Payments[0]
	Equal(t, p.Slip, "")
	Equal(t, p.Status, api.PaymentConfirmed)
	Equal(t, *p.Amount, int64(520))
	Equal(t, p.ReviewedBy, "admin")
	Equal(t, p.ReviewedAt != nil && *p.ReviewedAt == p.CreatedAt, true)
	Equal(t, alice.Order(o.ID).Status, api.OrderPaid)
}

func TestCashStillSettlesAfterCloseAndTopsUpAPaidOrder(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2)))
	h.Advance(25 * time.Hour)
	Equal(t, admin.Order(o.ID).Status, api.OrderCancelled)

	partial := admin.Cash(o.ID, 200)
	Equal(t, partial.Status, api.OrderPaid)
	Equal(t, partial.Balance, int64(320))

	done := admin.Cash(o.ID, 320)
	Equal(t, done.Status, api.OrderPaid)
	Equal(t, done.Paid, int64(520))
	Equal(t, done.Balance, int64(0))
	Equal(t, len(done.Payments), 2)
}

func TestCashIsRefusedWhileASlipIsPendingReview(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 1)))
	alice.UploadSlip(o.ID)

	admin.POST(cashPath(o.ID), api.AmountInput{Amount: 260}).Status(409)
	Equal(t, len(admin.Order(o.ID).Payments), 1)
}

func TestCashCannotExceedTheBalance(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2)))

	admin.POST(cashPath(o.ID), api.AmountInput{Amount: 600}).Status(409)
	Equal(t, len(admin.Order(o.ID).Payments), 0)

	admin.Cash(o.ID, 520)
	admin.POST(cashPath(o.ID), api.AmountInput{Amount: 520}).Status(409)

	got := admin.Order(o.ID)
	Equal(t, got.Paid, int64(520))
	Equal(t, got.Balance, int64(0))
	Equal(t, len(got.Payments), 1)
}

func TestAdjustmentsMoveTheBalanceOfPaidAndAwaitingOrders(t *testing.T) {
	h := New(t)
	alice, bob, admin := h.Member("alice"), h.Member("bob"), h.Admin()
	j := h.OpenJersey()
	paid := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2)))
	admin.Cash(paid.ID, 520)
	waiting := bob.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2)))

	got := admin.Adjust(paid.ID, -50, "ส่วนลด")
	Equal(t, got.Status, api.OrderPaid)
	Equal(t, got.Balance, int64(-50))
	Equal(t, len(got.Adjustments), 1)
	Equal(t, got.Adjustments[0].Amount, int64(-50))
	Equal(t, got.Adjustments[0].Note, "ส่วนลด")
	Equal(t, got.Adjustments[0].CreatedBy, "admin")

	more := admin.Adjust(waiting.ID, 30, "ค่าจัดส่งเพิ่ม")
	Equal(t, more.Status, api.OrderAwaitingPayment)
	Equal(t, more.Balance, int64(550))
	Equal(t, alice.Order(paid.ID).Balance, int64(-50))
}

func TestAdjustmentNeedsAnAmountAndANote(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2)))

	Equal(t, admin.POST(adjustPath(o.ID), api.AdjustmentInput{Amount: 30}).Status(400).Error(), "กรุณาระบุหมายเหตุ")
	Equal(t, admin.POST(adjustPath(o.ID), api.AdjustmentInput{Amount: 30, Note: "  "}).Status(400).Error(), "กรุณาระบุหมายเหตุ")
	Equal(t, admin.POST(adjustPath(o.ID), api.AdjustmentInput{Amount: 0, Note: "x"}).Status(400).Error(), "กรุณาระบุจำนวนเงินที่ถูกต้อง")
	alice.POST(adjustPath(o.ID), api.AdjustmentInput{Amount: 30, Note: "x"}).Status(403)
	admin.POST(adjustPath(9999), api.AdjustmentInput{Amount: 30, Note: "x"}).Status(404)

	got := alice.Order(o.ID)
	Equal(t, len(got.Adjustments), 0)
	Equal(t, got.Balance, int64(520))
}

func TestCancelledOrderWithAConfirmedPaymentOwesARefund(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2)))
	admin.Cash(o.ID, 520)

	Equal(t, As[api.Order](admin.POST(cancelPath(o.ID), nil).Status(200)).Status, api.OrderCancelled)

	got := alice.Order(o.ID)
	Equal(t, got.Status, api.OrderCancelled)
	Equal(t, got.Paid, int64(520))
	Equal(t, got.Balance, int64(-520))
}
