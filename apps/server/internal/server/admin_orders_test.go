package server_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

func adminOrdersPath(campaignID int64) string {
	return fmt.Sprintf("/api/admin/campaigns/%d/orders", campaignID)
}
func adminOrderPath(orderID int64) string { return fmt.Sprintf("/api/admin/orders/%d", orderID) }

func TestAdminOrdersForAMemberEvenWhenTheCampaignIsClosed(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()
	h.Advance(25 * time.Hour)
	in := api.AdminOrderInput{MemberID: alice.Me().Member.ID, Order: Pickup(Shirt(j, "M", 2))}

	alice.POST(adminOrdersPath(j.Campaign.ID), in).Status(403)

	o := As[api.Order](admin.POST(adminOrdersPath(j.Campaign.ID), in).Status(201))
	Equal(t, o.Status, api.OrderAwaitingPayment)
	Equal(t, o.PlacedByAdmin, true)
	Equal(t, o.MemberNickname, "alice")
	Equal(t, o.Total, int64(520))
	Equal(t, o.Balance, int64(520))

	mine := alice.Order(o.ID)
	Equal(t, mine.Status, api.OrderAwaitingPayment)
	Equal(t, mine.Balance, int64(520))
}

func TestAdminOrderTargetAndCampaignMustBeValid(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	nok := h.Applicant("nok").Me().Member
	j := h.OpenJersey()
	d := As[api.Campaign](admin.POST("/api/admin/campaigns", draft("ร่าง")).Status(201))
	in := api.AdminOrderInput{MemberID: alice.Me().Member.ID, Order: Pickup(Shirt(j, "M", 1))}

	admin.POST(adminOrdersPath(d.ID), in).Status(404)
	admin.POST(adminOrdersPath(9999), in).Status(404)

	bad := in
	bad.MemberID = nok.ID
	admin.POST(adminOrdersPath(j.Campaign.ID), bad).Status(400)
	bad.MemberID = 9999
	admin.POST(adminOrdersPath(j.Campaign.ID), bad).Status(400)

	bad = in
	bad.Order = Pickup()
	admin.POST(adminOrdersPath(j.Campaign.ID), bad).Status(400)
}

func TestAdminEditKeepsCapturedPriceOnlyForExistingItems(t *testing.T) {
	h := New(t)
	alice, bob, admin := h.Member("alice"), h.Member("bob"), h.Admin()
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2)))
	other := bob.Checkout(j.Campaign.ID, Pickup(Shirt(j, "M", 1)))

	admin.PUT(fmt.Sprintf("/api/admin/options/%d", j.Option.ID), api.OptionInput{
		Name: j.Option.Name, Price: new(int64(300)), Sizes: j.Option.Sizes, NeedsPrintName: true, NeedsNumber: true,
	}).Status(200)

	kept := Shirt(j, "L", 3)
	kept.ID = &o.Items[0].ID
	in := Ship(kept, Shirt(j, "M", 1))

	alice.PUT(adminOrderPath(o.ID), in).Status(403)

	got := As[api.Order](admin.PUT(adminOrderPath(o.ID), in).Status(200))
	Equal(t, len(got.Items), 2)
	Equal(t, got.Items[0].UnitPrice, int64(260))
	Equal(t, got.Items[0].Quantity, int64(3))
	Equal(t, got.Items[1].UnitPrice, int64(300))
	Equal(t, got.Total, int64(1080))
	Equal(t, got.Delivery, api.DeliveryShip)
	Equal(t, got.ShipAddress, "1 ถนนสุขุมวิท กรุงเทพฯ 10110")
	Equal(t, alice.Order(o.ID).Total, int64(1080))

	unknown := Shirt(j, "L", 1)
	unknown.ID = new(int64(999999))
	admin.PUT(adminOrderPath(o.ID), Pickup(unknown)).Status(400)

	foreign := Shirt(j, "L", 1)
	foreign.ID = &other.Items[0].ID
	admin.PUT(adminOrderPath(o.ID), Pickup(foreign)).Status(400)

	Equal(t, admin.Order(o.ID).Total, int64(1080))
	Equal(t, admin.Order(other.ID).Total, int64(260))
}

func TestAdminEditWorksAfterTheCampaignClosesAndOnCancelledOrders(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 1)))
	admin.POST(cancelPath(o.ID), nil).Status(200)
	h.Advance(25 * time.Hour)

	got := As[api.Order](admin.PUT(adminOrderPath(o.ID), Pickup(Shirt(j, "S", 4))).Status(200))
	Equal(t, got.Status, api.OrderCancelled)
	Equal(t, got.Items[0].Size, "S")
	Equal(t, got.Total, int64(1040))

	admin.PUT(adminOrderPath(9999), Pickup(Shirt(j, "S", 1))).Status(404)
}
