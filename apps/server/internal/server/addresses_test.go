package server_test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

func memberAddressesPath(memberID int64) string {
	return fmt.Sprintf("/api/admin/members/%d/addresses", memberID)
}

func shipTo(c api.AddressChoice, items ...api.ItemInput) api.OrderInput {
	return api.OrderInput{
		Delivery: api.DeliveryShip, ShipName: c.Name, ShipAddress: c.Address, ShipPhone: c.Phone, Items: items,
	}
}

func addressChoices(res *Res) []api.AddressChoice {
	res.Status(200)
	return As[[]api.AddressChoice](res)
}

func TestAddressChoicesAreTheProfileThenDistinctOrderAddressesNewestFirst(t *testing.T) {
	h := New(t)
	alice := h.Member("alice")
	j := h.OpenJersey()
	profile := api.AddressChoice{Name: "Alice Real", Address: "1 ถนนสุขุมวิท", Phone: "0811111111"}
	x := api.AddressChoice{Name: "แม่ของอลิซ", Address: "9 ถนนพหลโยธิน", Phone: "0822222222"}
	y := api.AddressChoice{Name: "ที่ทำงาน", Address: "5 ถนนสีลม", Phone: "0833333333"}
	alice.PUT("/api/me", api.ProfileInput{Cohort: "34", Nickname: "alice", RealName: profile.Name, Address: profile.Address, Phone: profile.Phone}).Status(200)

	alice.Checkout(j.Campaign.ID, shipTo(x, Shirt(j, "M", 1)))
	alice.Checkout(j.Campaign.ID, shipTo(y, Shirt(j, "M", 1)))
	padded := x
	padded.Address = "  " + x.Address + " "
	alice.Checkout(j.Campaign.ID, shipTo(padded, Shirt(j, "M", 1)))
	alice.Checkout(j.Campaign.ID, shipTo(profile, Shirt(j, "M", 1)))
	alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "M", 1)))

	got := addressChoices(alice.GET("/api/me/addresses"))
	Equal(t, slices.Equal(got, []api.AddressChoice{profile, x, y}), true)
}

func TestAddressChoicesOmitAnIncompleteProfile(t *testing.T) {
	h := New(t)
	alice := h.Member("alice")
	j := h.OpenJersey()
	x := api.AddressChoice{Name: "แม่ของอลิซ", Address: "9 ถนนพหลโยธิน", Phone: "0822222222"}

	res := alice.GET("/api/me/addresses").Status(200)
	Equal(t, strings.TrimSpace(res.Text()), "[]")

	alice.PUT("/api/me", api.ProfileInput{Cohort: "34", Nickname: "alice", RealName: "Alice Real", Address: "1 ถนนสุขุมวิท"}).Status(200)
	Equal(t, len(addressChoices(alice.GET("/api/me/addresses"))), 0)

	alice.Checkout(j.Campaign.ID, shipTo(x, Shirt(j, "M", 1)))
	Equal(t, slices.Equal(addressChoices(alice.GET("/api/me/addresses")), []api.AddressChoice{x}), true)
}

func TestAddressChoicesNeverIncludeAnotherMembersOrders(t *testing.T) {
	h := New(t)
	alice, bob := h.Member("alice"), h.Member("bob")
	j := h.OpenJersey()
	z := api.AddressChoice{Name: "บ๊อบ", Address: "7 ถนนรัชดา", Phone: "0844444444"}
	bob.Checkout(j.Campaign.ID, shipTo(z, Shirt(j, "M", 1)))

	Equal(t, len(addressChoices(alice.GET("/api/me/addresses"))), 0)
	Equal(t, slices.Equal(addressChoices(bob.GET("/api/me/addresses")), []api.AddressChoice{z}), true)
}

func TestCheckoutWithANewAddressLeavesTheProfileAlone(t *testing.T) {
	h := New(t)
	alice := h.Member("alice")
	j := h.OpenJersey()
	alice.PUT("/api/me", api.ProfileInput{Cohort: "34", Nickname: "alice", RealName: "Alice Real", Address: "1 ถนนสุขุมวิท", Phone: "0811111111"}).Status(200)
	before := *alice.Me().Member

	alice.Checkout(j.Campaign.ID, Ship(Shirt(j, "M", 1)))

	Equal(t, *alice.Me().Member, before)
}

func TestAdminReadsAddressChoicesOfAnotherMember(t *testing.T) {
	h := New(t)
	alice, bob, admin := h.Member("alice"), h.Member("bob"), h.Admin()
	j := h.OpenJersey()
	x := api.AddressChoice{Name: "แม่ของอลิซ", Address: "9 ถนนพหลโยธิน", Phone: "0822222222"}
	aliceID := alice.Me().Member.ID
	alice.PUT("/api/me", api.ProfileInput{Cohort: "34", Nickname: "alice", RealName: "Alice Real", Address: "1 ถนนสุขุมวิท", Phone: "0811111111"}).Status(200)
	alice.Checkout(j.Campaign.ID, shipTo(x, Shirt(j, "M", 1)))

	got := addressChoices(admin.GET(memberAddressesPath(aliceID)))
	Equal(t, slices.Equal(got, []api.AddressChoice{{Name: "Alice Real", Address: "1 ถนนสุขุมวิท", Phone: "0811111111"}, x}), true)

	bob.GET(memberAddressesPath(aliceID)).Status(403)
	h.Anonymous().GET(memberAddressesPath(aliceID)).Status(401)
	h.Anonymous().GET("/api/me/addresses").Status(401)
	admin.GET(memberAddressesPath(9999)).Status(404)

	unlinked := As[api.Member](admin.POST("/api/admin/members", api.ProfileInput{Nickname: "ติ๋ว"}).Status(201))
	Equal(t, len(addressChoices(admin.GET(memberAddressesPath(unlinked.ID)))), 0)
	y := api.AddressChoice{Name: "ติ๋ว", Address: "3 ถนนเพชรบุรี", Phone: "0855555555"}
	admin.POST(adminOrdersPath(j.Campaign.ID), api.AdminOrderInput{MemberID: unlinked.ID, Order: shipTo(y, Shirt(j, "M", 1))}).Status(201)
	Equal(t, slices.Equal(addressChoices(admin.GET(memberAddressesPath(unlinked.ID))), []api.AddressChoice{y}), true)
}

func TestErasedMemberOffersNoAddresses(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()
	aliceID := alice.Me().Member.ID
	alice.PUT("/api/me", api.ProfileInput{Cohort: "34", Nickname: "alice", RealName: "Alice Real", Address: "1 ถนนสุขุมวิท", Phone: "0811111111"}).Status(200)
	alice.Checkout(j.Campaign.ID, Ship(Shirt(j, "M", 1)))
	Equal(t, len(addressChoices(admin.GET(memberAddressesPath(aliceID)))), 2)

	alice.DELETE("/api/me").Status(204)

	Equal(t, len(addressChoices(admin.GET(memberAddressesPath(aliceID)))), 0)
}
