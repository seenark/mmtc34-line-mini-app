package server_test

import (
	"fmt"
	"testing"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

func linkOrderPath(orderID int64) string { return fmt.Sprintf("/api/admin/orders/%d/link", orderID) }

// unlinkedOwnerOfTwoOrders is an Unlinked Member with two Admin-placed Orders.
type unlinkedOwnerOfTwoOrders struct {
	J             Jersey
	Owner         api.Member
	First, Second api.Order
}

func (h *H) UnlinkedWithTwoOrders() unlinkedOwnerOfTwoOrders {
	admin := h.Admin()
	j := h.OpenJersey()
	u := As[api.Member](admin.POST("/api/admin/members", api.ProfileInput{Nickname: "ติ๋ว", Phone: "0899999999"}).Status(201))
	place := func(in api.OrderInput) api.Order {
		return As[api.Order](admin.POST(adminOrdersPath(j.Campaign.ID), api.AdminOrderInput{MemberID: u.ID, Order: in}).Status(201))
	}
	return unlinkedOwnerOfTwoOrders{
		J: j, Owner: u,
		First:  place(Ship(Shirt(j, "M", 2))),
		Second: place(Pickup(Shirt(j, "L", 1))),
	}
}

func TestAdminCreatesUnlinkedMemberFromANicknameAlone(t *testing.T) {
	h := New(t)
	admin, alice := h.Admin(), h.Member("alice")

	admin.POST("/api/admin/members", api.ProfileInput{Nickname: "  "}).Status(400)
	alice.POST("/api/admin/members", api.ProfileInput{Nickname: "ติ๋ว"}).Status(403)

	m := As[api.Member](admin.POST("/api/admin/members", api.ProfileInput{Nickname: " ติ๋ว "}).Status(201))
	Equal(t, m.Nickname, "ติ๋ว")
	Equal(t, m.Linked, false)
	Equal(t, m.Status, api.MemberApproved)
	Equal(t, m.Role, api.RoleMember)
	Equal(t, As[api.Member](admin.GET(fmt.Sprintf("/api/admin/members/%d", m.ID)).Status(200)).Linked, false)

	j := h.OpenJersey()
	o := As[api.Order](admin.POST(adminOrdersPath(j.Campaign.ID), api.AdminOrderInput{
		MemberID: m.ID, Order: Ship(Shirt(j, "M", 1)),
	}).Status(201))
	Equal(t, o.MemberID, m.ID)
	Equal(t, o.MemberNickname, "ติ๋ว")
}

func TestLinkInfoNamesTheUnlinkedOwnerAndCountsTheirOtherOrders(t *testing.T) {
	h := New(t)
	admin, alice := h.Admin(), h.Member("alice")
	u := h.UnlinkedWithTwoOrders()

	info := As[api.OrderLinkInfo](admin.GET(linkOrderPath(u.First.ID)).Status(200))

	Equal(t, info.UnlinkedMember.ID, u.Owner.ID)
	Equal(t, info.UnlinkedMember.Nickname, "ติ๋ว")
	Equal(t, info.OtherOrders, int64(1))
	alice.GET(linkOrderPath(u.First.ID)).Status(403)
	admin.GET(linkOrderPath(9999)).Status(404)
}

func TestLinkingOneOrderLeavesTheUnlinkedMemberUntilTheyAreEmptied(t *testing.T) {
	h := New(t)
	admin, tiw := h.Admin(), h.Member("tiw")
	tiwID := tiw.Me().Member.ID
	u := h.UnlinkedWithTwoOrders()

	got := As[api.OrderLinkResult](admin.POST(linkOrderPath(u.First.ID), api.OrderLinkInput{MemberID: tiwID}).Status(200))

	Equal(t, got.Moved, int64(1))
	Equal(t, got.UnlinkedRemoved, false)
	mine := As[[]api.Order](tiw.GET("/api/orders").Status(200))
	Equal(t, len(mine), 1)
	Equal(t, mine[0].ID, u.First.ID)
	Equal(t, admin.Order(u.Second.ID).MemberID, u.Owner.ID)
	Equal(t, As[api.Member](admin.GET(memberPath(u.Owner.ID)).Status(200)).Linked, false)
	Equal(t, As[api.OrderLinkInfo](admin.GET(linkOrderPath(u.Second.ID)).Status(200)).OtherOrders, int64(0))

	got = As[api.OrderLinkResult](admin.POST(linkOrderPath(u.Second.ID), api.OrderLinkInput{MemberID: tiwID}).Status(200))

	Equal(t, got.Moved, int64(1))
	Equal(t, got.UnlinkedRemoved, true)
	admin.GET(memberPath(u.Owner.ID)).Status(404)
	Equal(t, len(As[[]api.Order](tiw.GET("/api/orders").Status(200))), 2)
}

func TestLinkingAllOrdersMovesEveryOrderAndRemovesTheUnlinkedMember(t *testing.T) {
	h := New(t)
	admin, tiw := h.Admin(), h.Member("tiw")
	u := h.UnlinkedWithTwoOrders()
	sponsor := As[api.SponsorPosition](admin.POST(fmt.Sprintf("/api/admin/campaigns/%d/sponsors", u.J.Campaign.ID), api.SponsorInput{
		Name: "กลางอก", MemberID: &u.Owner.ID,
	}).Status(201))

	got := As[api.OrderLinkResult](admin.POST(linkOrderPath(u.First.ID), api.OrderLinkInput{MemberID: tiw.Me().Member.ID, AllOrders: true}).Status(200))

	Equal(t, got.Moved, int64(2))
	Equal(t, got.UnlinkedRemoved, true)
	admin.GET(memberPath(u.Owner.ID)).Status(404)
	Equal(t, len(As[[]api.Order](tiw.GET("/api/orders").Status(200))), 2)
	positions := As[[]api.SponsorPosition](tiw.GET(fmt.Sprintf("/api/campaigns/%d/sponsors", u.J.Campaign.ID)).Status(200))
	Equal(t, positions[0].ID, sponsor.ID)
	Equal(t, positions[0].MemberNickname, "tiw")
}

func TestLinkedOrderKeepsItsLedgerAndAdminMarkers(t *testing.T) {
	h := New(t)
	admin, tiw := h.Admin(), h.Member("tiw")
	u := h.UnlinkedWithTwoOrders()
	admin.Cash(u.First.ID, 300)
	admin.Cash(u.First.ID, 100)
	before := admin.Adjust(u.First.ID, -20, "ส่วนลด")
	Equal(t, before.Balance, int64(100))

	admin.POST(linkOrderPath(u.First.ID), api.OrderLinkInput{MemberID: tiw.Me().Member.ID}).Status(200)

	after := tiw.Order(u.First.ID)
	Equal(t, after.Balance, before.Balance)
	Equal(t, after.Paid, before.Paid)
	Equal(t, after.Status, before.Status)
	Equal(t, len(after.Payments), 2)
	Equal(t, len(after.Adjustments), 1)
	Equal(t, after.PlacedByAdmin, true)
	Equal(t, after.Delivery, api.DeliveryShip)
	Equal(t, after.ShipAddress, u.First.ShipAddress)
	Equal(t, after.MemberNickname, "tiw")
}

func TestLinkingIsRefusedForBadOwnersTargetsAndCallers(t *testing.T) {
	h := New(t)
	admin, alice, tiw := h.Admin(), h.Member("alice"), h.Member("tiw")
	gone := h.Member("gone")
	goneID := gone.Me().Member.ID
	gone.DELETE("/api/me").Status(204)
	u := h.UnlinkedWithTwoOrders()
	other := As[api.Member](admin.POST("/api/admin/members", api.ProfileInput{Nickname: "อีกคน"}).Status(201))
	tiwID := tiw.Me().Member.ID
	link := func(target int64) api.OrderLinkInput { return api.OrderLinkInput{MemberID: target} }

	aliceOrder := alice.Checkout(u.J.Campaign.ID, Pickup(Shirt(u.J, "M", 1)))
	admin.POST(linkOrderPath(aliceOrder.ID), link(tiwID)).Status(409)
	admin.GET(linkOrderPath(aliceOrder.ID)).Status(409)

	erased := h.Member("erased")
	erasedOrder := erased.Checkout(u.J.Campaign.ID, Pickup(Shirt(u.J, "M", 1)))
	erased.DELETE("/api/me").Status(204)
	admin.POST(linkOrderPath(erasedOrder.ID), link(tiwID)).Status(409)
	admin.GET(linkOrderPath(erasedOrder.ID)).Status(409)

	admin.POST(linkOrderPath(u.First.ID), link(other.ID)).Status(400)
	admin.POST(linkOrderPath(u.First.ID), link(goneID)).Status(400)
	admin.POST(linkOrderPath(u.First.ID), link(h.SuperAdmin().Me().Member.ID)).Status(400)
	admin.POST(linkOrderPath(u.First.ID), link(9999)).Status(400)
	admin.POST(linkOrderPath(9999), link(tiwID)).Status(404)
	alice.POST(linkOrderPath(u.First.ID), link(tiwID)).Status(403)
	h.Anonymous().POST(linkOrderPath(u.First.ID), link(tiwID)).Status(401)

	Equal(t, admin.Order(u.First.ID).MemberID, u.Owner.ID)
	Equal(t, admin.Order(u.Second.ID).MemberID, u.Owner.ID)
	Equal(t, len(As[[]api.Order](tiw.GET("/api/orders").Status(200))), 0)
}
