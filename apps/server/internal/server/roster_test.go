package server_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

func rosterPath(campaignID int64) string { return fmt.Sprintf("/api/campaigns/%d/roster", campaignID) }

func TestRosterMergesEachMembersPendingReviewAndPaidOrders(t *testing.T) {
	h := New(t)
	alice, bob, carol, dave, admin := h.Member("alice"), h.Member("bob"), h.Member("carol"), h.Member("dave"), h.Admin()
	j := h.OpenJersey()

	paid := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 2)))
	alice.UploadSlip(paid.ID)
	admin.POST(confirmPath(paid.ID), api.AmountInput{Amount: 520}).Status(200)

	pending := dave.Checkout(j.Campaign.ID, Pickup(Shirt(j, "S", 1)))
	dave.UploadSlip(pending.ID)

	again := alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "L", 1), Shirt(j, "M", 3)))
	alice.UploadSlip(again.ID)

	alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "XL", 5)))
	bob.Checkout(j.Campaign.ID, Pickup(Shirt(j, "S", 9)))
	cancelled := carol.Checkout(j.Campaign.ID, Pickup(Shirt(j, "M", 4)))
	carol.POST(cancelPath(cancelled.ID), nil).Status(200)

	got := As[[]api.RosterEntry](bob.GET(rosterPath(j.Campaign.ID)).Status(200))

	want := []api.RosterEntry{
		{Nickname: "alice", Items: []api.RosterItem{
			{OptionName: j.Option.Name, Size: "L", Quantity: 3},
			{OptionName: j.Option.Name, Size: "M", Quantity: 3},
		}},
		{Nickname: "dave", Items: []api.RosterItem{{OptionName: j.Option.Name, Size: "S", Quantity: 1}}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("roster\n got %+v\nwant %+v", got, want)
	}
}

func TestRosterShowsNothingButNicknameAndItems(t *testing.T) {
	h := New(t)
	alice, bob := h.Member("alice"), h.Member("bob")
	j := h.OpenJersey()
	in := Ship(Shirt(j, "M", 1))
	o := alice.Checkout(j.Campaign.ID, in)
	alice.UploadSlip(o.ID)

	body := bob.GET(rosterPath(j.Campaign.ID)).Status(200).Text()

	for _, leaked := range []string{"address", "phone", "realName", "status", "balance", "memberId", in.ShipAddress, in.ShipPhone, in.ShipName} {
		if strings.Contains(strings.ToLower(body), strings.ToLower(leaked)) {
			t.Fatalf("roster leaks %q: %s", leaked, body)
		}
	}
}

func TestRosterIsEmptyArrayWithoutEligibleOrders(t *testing.T) {
	h := New(t)
	alice := h.Member("alice")
	j := h.OpenJersey()
	alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "M", 1)))

	Equal(t, alice.GET(rosterPath(j.Campaign.ID)).Status(200).Text(), "[]")
}

func TestRosterIsForMembersOfVisibleCampaigns(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()
	d := As[api.Campaign](admin.POST("/api/admin/campaigns", draft("ร่าง")).Status(201))

	h.Visitor("nok").GET(rosterPath(j.Campaign.ID)).Status(403)
	h.Anonymous().GET(rosterPath(j.Campaign.ID)).Status(401)
	alice.GET(rosterPath(d.ID)).Status(404)
	alice.GET(rosterPath(9999)).Status(404)
	admin.GET(rosterPath(d.ID)).Status(200)
}
