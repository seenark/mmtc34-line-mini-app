package server_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

func ordersPath(campaignID int64) string { return fmt.Sprintf("/api/campaigns/%d/orders", campaignID) }
func cancelPath(orderID int64) string    { return fmt.Sprintf("/api/orders/%d/cancel", orderID) }

func TestCheckoutValidation(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()

	other := As[api.Campaign](admin.POST("/api/admin/campaigns", api.CampaignInput{Type: api.CampaignEvent, Title: "งานเลี้ยง", Status: api.CampaignOpen}).Status(201))
	foreign := As[api.Option](admin.POST(fmt.Sprintf("/api/admin/campaigns/%d/options", other.ID), api.OptionInput{Name: "บัตร", Price: new(int64(100))}).Status(201))

	cases := []struct {
		name string
		in   api.OrderInput
		err  string
	}{
		{"no items", Pickup(), "กรุณาเลือกสินค้าอย่างน้อย 1 รายการ"},
		{"wrong size", Pickup(Shirt(j, "XXL", 1)), "รายการที่ 1: กรุณาเลือกไซซ์"},
		{"missing print name", Pickup(func() api.ItemInput { i := Shirt(j, "M", 1); i.PrintName = " "; return i }()), "รายการที่ 1: กรุณากรอกชื่อที่จะพิมพ์"},
		{"missing number", Pickup(func() api.ItemInput { i := Shirt(j, "M", 1); i.Number = ""; return i }()), "รายการที่ 1: กรุณากรอกเบอร์ที่จะพิมพ์"},
		{"quantity 0", Pickup(Shirt(j, "M", 1), Shirt(j, "M", 0)), "รายการที่ 2: จำนวนต้องอยู่ระหว่าง 1 ถึง 99"},
		{"quantity 100", Pickup(Shirt(j, "M", 100)), "รายการที่ 1: จำนวนต้องอยู่ระหว่าง 1 ถึง 99"},
		{"ship without address", func() api.OrderInput { in := Ship(Shirt(j, "M", 1)); in.ShipAddress = "  "; return in }(), "กรุณากรอกชื่อ ที่อยู่ และเบอร์โทรสำหรับจัดส่ง"},
		{"amount on fixed price", Pickup(func() api.ItemInput { i := Shirt(j, "M", 1); i.Amount = new(int64(5)); return i }()), "รายการที่ 1: ตัวเลือกนี้มีราคาคงที่ ไม่ต้องระบุจำนวนเงิน"},
		{"option of another campaign", Pickup(api.ItemInput{OptionID: foreign.ID, Quantity: 1}), "รายการที่ 1: ไม่พบตัวเลือกนี้ในอีเวนต์"},
		{"unknown delivery", api.OrderInput{Items: []api.ItemInput{Shirt(j, "M", 1)}}, "วิธีรับสินค้าไม่ถูกต้อง"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			Equal(t, alice.POST(ordersPath(j.Campaign.ID), tc.in).Status(400).Error(), tc.err)
		})
	}
	Equal(t, len(As[[]api.Order](alice.GET("/api/orders").Status(200))), 0)
}

func TestDonationAmountIsChosenByTheMember(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	camp := As[api.Campaign](admin.POST("/api/admin/campaigns", api.CampaignInput{Type: api.CampaignDonation, Title: "ระดมทุน", Status: api.CampaignOpen}).Status(201))
	opt := As[api.Option](admin.POST(fmt.Sprintf("/api/admin/campaigns/%d/options", camp.ID), api.OptionInput{Name: "บริจาค"}).Status(201))

	gift := func(amount *int64, qty int64) api.OrderInput {
		return NoDelivery(api.ItemInput{OptionID: opt.ID, Quantity: qty, Amount: amount})
	}
	alice.POST(ordersPath(camp.ID), gift(nil, 1)).Status(400)
	alice.POST(ordersPath(camp.ID), gift(new(int64(0)), 1)).Status(400)
	alice.POST(ordersPath(camp.ID), gift(new(int64(100)), 2)).Status(400)

	o := alice.Checkout(camp.ID, gift(new(int64(150)), 1))
	Equal(t, o.Total, int64(150))
	Equal(t, o.Items[0].UnitPrice, int64(150))
}

func TestCheckoutLimits(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()
	camp := As[api.Campaign](admin.POST("/api/admin/campaigns", api.CampaignInput{Type: api.CampaignDonation, Title: "ระดมทุน", Status: api.CampaignOpen}).Status(201))
	opt := As[api.Option](admin.POST(fmt.Sprintf("/api/admin/campaigns/%d/options", camp.ID), api.OptionInput{Name: "บริจาค"}).Status(201))

	many := make([]api.ItemInput, 51)
	for i := range many {
		many[i] = Shirt(j, "M", 1)
	}
	Equal(t, alice.POST(ordersPath(j.Campaign.ID), Pickup(many...)).Status(400).Error(), "สั่งซื้อได้ไม่เกิน 50 รายการต่อคำสั่งซื้อ")
	Equal(t, len(alice.Checkout(j.Campaign.ID, Pickup(many[:50]...)).Items), 50)

	gift := func(amount int64) api.OrderInput {
		return NoDelivery(api.ItemInput{OptionID: opt.ID, Quantity: 1, Amount: &amount})
	}
	Equal(t, alice.POST(ordersPath(camp.ID), gift(1_000_001)).Status(400).Error(), "รายการที่ 1: จำนวนเงินต้องไม่เกิน 1,000,000 บาท")
	Equal(t, alice.Checkout(camp.ID, gift(1_000_000)).Total, int64(1_000_000))
}

func TestCheckoutTotalsAndCapturedPrice(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()

	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2)))
	Equal(t, o.Status, api.OrderAwaitingPayment)
	Equal(t, o.Total, int64(520))
	Equal(t, o.Paid, int64(0))
	Equal(t, o.Balance, int64(520))
	Equal(t, o.Items[0].UnitPrice, int64(260))
	Equal(t, o.Items[0].OptionName, "เสื้อรุ่น")
	Equal(t, o.MemberNickname, "alice")
	Equal(t, o.PlacedByAdmin, false)
	Equal(t, o.Actions, api.OrderActions{Cancel: true, UploadSlip: true})

	admin.PUT(fmt.Sprintf("/api/admin/options/%d", j.Option.ID), api.OptionInput{
		Name: j.Option.Name, Price: new(int64(300)), Sizes: j.Option.Sizes, NeedsPrintName: true, NeedsNumber: true,
	}).Status(200)

	Equal(t, alice.Order(o.ID).Total, int64(520))
	Equal(t, alice.Order(o.ID).Items[0].UnitPrice, int64(260))
	Equal(t, alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2))).Total, int64(600))
}

func TestShipOrderKeepsAddressAndPickupDiscardsIt(t *testing.T) {
	h := New(t)
	alice := h.Member("alice")
	j := h.OpenJersey()

	shipped := alice.Checkout(j.Campaign.ID, Ship(Shirt(j, "M", 1)))
	Equal(t, shipped.Delivery, api.DeliveryShip)
	Equal(t, shipped.ShipAddress, "1 ถนนสุขุมวิท กรุงเทพฯ 10110")

	in := Pickup(Shirt(j, "M", 1))
	in.ShipName, in.ShipAddress, in.ShipPhone = "x", "y", "z"
	picked := alice.Checkout(j.Campaign.ID, in)
	Equal(t, picked.Delivery, api.DeliveryPickup)
	Equal(t, picked.ShipAddress+picked.ShipName+picked.ShipPhone, "")
}

func TestDonationOnlyOrderHasNoDelivery(t *testing.T) {
	h := New(t)
	alice := h.Member("alice")
	j := h.OpenJersey()
	gift := h.AddDonation(j.Campaign.ID)

	in := NoDelivery(Gift(gift, 150))
	in.ShipName, in.ShipAddress, in.ShipPhone = "x", "y", "z"
	o := alice.Checkout(j.Campaign.ID, in)
	Equal(t, o.Delivery, api.DeliveryNone)
	Equal(t, o.ShipName+o.ShipAddress+o.ShipPhone, "")
	Equal(t, o.Total, int64(150))
	Equal(t, alice.Order(o.ID).Delivery, api.DeliveryNone)

	const wantNone = "คำสั่งซื้อที่มีแต่สมทบทุนไม่ต้องเลือกวิธีรับสินค้า"
	Equal(t, alice.POST(ordersPath(j.Campaign.ID), Ship(Gift(gift, 150))).Status(400).Error(), wantNone)
	Equal(t, alice.POST(ordersPath(j.Campaign.ID), Pickup(Gift(gift, 150))).Status(400).Error(), wantNone)
	Equal(t, len(As[[]api.Order](alice.GET("/api/orders").Status(200))), 1)
}

func TestNoDeliveryIsRejectedForOrdersWithMerch(t *testing.T) {
	h := New(t)
	alice := h.Member("alice")
	j := h.OpenJersey()
	gift := h.AddDonation(j.Campaign.ID)

	const wantMerch = "วิธีรับสินค้า 'ไม่จัดส่ง' ใช้ได้เฉพาะคำสั่งซื้อที่มีแต่สมทบทุน"
	Equal(t, alice.POST(ordersPath(j.Campaign.ID), NoDelivery(Shirt(j, "M", 1))).Status(400).Error(), wantMerch)
	Equal(t, alice.POST(ordersPath(j.Campaign.ID), NoDelivery(Shirt(j, "M", 1), Gift(gift, 100))).Status(400).Error(), wantMerch)

	mixed := alice.Checkout(j.Campaign.ID, Ship(Shirt(j, "M", 1), Gift(gift, 100)))
	Equal(t, mixed.Delivery, api.DeliveryShip)
	Equal(t, mixed.Total, int64(360))
	Equal(t, len(As[[]api.Order](alice.GET("/api/orders").Status(200))), 1)
}

func TestMemberCancelsOwnAwaitingOrder(t *testing.T) {
	h := New(t)
	alice := h.Member("alice")
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2)))

	got := As[api.Order](alice.POST(cancelPath(o.ID), nil).Status(200))
	Equal(t, got.Status, api.OrderCancelled)
	Equal(t, got.Balance, int64(0))
	Equal(t, got.Actions, api.OrderActions{})
	Equal(t, alice.Order(o.ID).Status, api.OrderCancelled)
	alice.POST(cancelPath(o.ID), nil).Status(409)
}

func TestOtherMemberCannotSeeOrCancel(t *testing.T) {
	h := New(t)
	alice, bob := h.Member("alice"), h.Member("bob")
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 1)))

	bob.GET(fmt.Sprintf("/api/orders/%d", o.ID)).Status(404)
	bob.POST(cancelPath(o.ID), nil).Status(404)
	bob.GET("/api/orders/999").Status(404)
	Equal(t, alice.Order(o.ID).Status, api.OrderAwaitingPayment)
}

func TestClosedCampaignReadsUnpaidOrderAsCancelled(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2)))

	h.Advance(25 * time.Hour)

	Equal(t, alice.POST(ordersPath(j.Campaign.ID), Pickup(Shirt(j, "M", 1))).Status(409).Error(), "อีเวนต์ปิดรับแล้ว")
	got := alice.Order(o.ID)
	Equal(t, got.Status, api.OrderCancelled)
	Equal(t, got.Total, int64(520))
	Equal(t, got.Balance, int64(0))
	Equal(t, got.Actions.Cancel, false)
	Equal(t, got.Actions.UploadSlip, false)
	Equal(t, alice.POST(cancelPath(o.ID), nil).Status(409).Error(), "อีเวนต์ปิดแล้ว กรุณาติดต่อแอดมิน")

	Equal(t, admin.Order(o.ID).Status, api.OrderCancelled)
}

func TestAdminCancelsAnyOrderAndIsIdempotent(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 1)))

	seen := admin.Order(o.ID)
	Equal(t, seen.Actions, api.OrderActions{Cancel: true, UploadSlip: true, RecordCash: true, Edit: true, AddAdjust: true})

	Equal(t, As[api.Order](admin.POST(cancelPath(o.ID), nil).Status(200)).Status, api.OrderCancelled)
	Equal(t, As[api.Order](admin.POST(cancelPath(o.ID), nil).Status(200)).Status, api.OrderCancelled)

	late := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 1)))
	h.Advance(25 * time.Hour)
	Equal(t, As[api.Order](admin.POST(cancelPath(late.ID), nil).Status(200)).Status, api.OrderCancelled)
}

func TestCheckoutNeedsAnOpenVisibleCampaign(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()
	d := As[api.Campaign](admin.POST("/api/admin/campaigns", draft("ร่าง")).Status(201))

	alice.POST(ordersPath(d.ID), Pickup(Shirt(j, "M", 1))).Status(404)
	alice.POST(ordersPath(9999), Pickup(Shirt(j, "M", 1))).Status(404)
	h.Close(j.Campaign.ID)
	alice.POST(ordersPath(j.Campaign.ID), Pickup(Shirt(j, "M", 1))).Status(409)
	h.Anonymous().POST(ordersPath(j.Campaign.ID), Pickup(Shirt(j, "M", 1))).Status(401)
	h.Visitor("nok").POST(ordersPath(j.Campaign.ID), Pickup(Shirt(j, "M", 1))).Status(403)
}

func TestMyOrdersListsOnlyMineNewestFirst(t *testing.T) {
	h := New(t)
	alice, bob := h.Member("alice"), h.Member("bob")
	j := h.OpenJersey()
	a1 := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "S", 1)))
	b1 := bob.Checkout(j.Campaign.ID, Pickup(Shirt(j, "M", 1)))
	a2 := alice.Checkout(j.Campaign.ID, Ship(Shirt(j, "L", 3)))

	mine := As[[]api.Order](alice.GET("/api/orders").Status(200))
	Equal(t, len(mine), 2)
	Equal(t, mine[0].ID, a2.ID)
	Equal(t, mine[1].ID, a1.ID)
	Equal(t, mine[0].Total, int64(780))

	theirs := As[[]api.Order](bob.GET("/api/orders").Status(200))
	Equal(t, len(theirs), 1)
	Equal(t, theirs[0].ID, b1.ID)
}
