package server_test

import (
	"fmt"
	"testing"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

func memberPath(id int64) string { return fmt.Sprintf("/api/admin/members/%d", id) }

func TestMemberEditsOwnProfile(t *testing.T) {
	h := New(t)
	alice := h.Member("alice")

	edit := api.ProfileInput{Cohort: "35", Nickname: "อลิซ", RealName: "อลิซ ใจดี", Address: "9 ถนนสีลม", Phone: "0899999999"}
	got := As[api.Member](alice.PUT("/api/me", edit).Status(200))
	Equal(t, got.Nickname, "อลิซ")
	Equal(t, got.Cohort, "35")
	Equal(t, got.Address, "9 ถนนสีลม")
	Equal(t, alice.Me().Member.RealName, "อลิซ ใจดี")

	edit.RealName = " "
	alice.PUT("/api/me", edit).Status(400)
	Equal(t, alice.Me().Member.RealName, "อลิซ ใจดี")
}

func TestProfileEditKeepsOrderShipSnapshot(t *testing.T) {
	h := New(t)
	alice, j := h.Member("alice"), h.OpenJersey()
	o := alice.Checkout(j.Campaign.ID, Ship(Shirt(j, "M", 1)))

	alice.PUT("/api/me", api.ProfileInput{Cohort: "34", Nickname: "alice", RealName: "Alice New", Address: "ที่อยู่ใหม่", Phone: "0800000000"}).Status(200)

	got := alice.Order(o.ID)
	Equal(t, got.ShipName, o.ShipName)
	Equal(t, got.ShipAddress, o.ShipAddress)
	Equal(t, got.ShipPhone, o.ShipPhone)
}

func TestDeleteAccountErasesContactDataKeepsOrders(t *testing.T) {
	h := New(t)
	alice, admin, j := h.Member("alice"), h.Admin(), h.OpenJersey()
	aliceID := alice.Me().Member.ID
	alice.PUT("/api/me", api.ProfileInput{Cohort: "34", Nickname: "alice", RealName: "Alice Real", Address: "1 ถนนสุขุมวิท", Phone: "0811111111"}).Status(200)
	o := alice.Checkout(j.Campaign.ID, Ship(Shirt(j, "M", 2)))

	alice.DELETE("/api/me").Status(204)
	alice.GET("/api/me").Status(401)

	m := As[api.Member](admin.GET(memberPath(aliceID)).Status(200))
	Equal(t, m.Deleted, true)
	Equal(t, m.Linked, false)
	Equal(t, m.Nickname, "alice")
	Equal(t, m.RealName+m.Cohort+m.Address+m.Phone, "")

	got := As[api.Order](admin.GET(fmt.Sprintf("/api/orders/%d", o.ID)).Status(200))
	Equal(t, got.ShipName+got.ShipAddress+got.ShipPhone, "")
	Equal(t, got.MemberNickname, "alice")
	Equal(t, got.Total, o.Total)
	Equal(t, got.Balance, o.Balance)
	Equal(t, len(got.Items), len(o.Items))

	admin.PUT(memberPath(aliceID), api.ProfileInput{Cohort: "34", Nickname: "alice", RealName: "x"}).Status(404)
	admin.PUT(memberPath(aliceID)+"/role", api.RoleInput{Role: api.RoleAdmin}).Status(404)
}

func TestSameLineAccountOnboardsFreshAfterDeletion(t *testing.T) {
	h := New(t)
	alice := h.Member("alice")
	oldID := alice.Me().Member.ID
	alice.DELETE("/api/me").Status(204)

	again := h.Anonymous()
	again.POST("/api/auth/dev", api.DevLoginInput{LineUserID: "U-alice"}).Status(204)
	if again.Me().Member != nil {
		t.Fatal("deleted account still resolves to a Member")
	}
	fresh := As[api.Member](again.POST("/api/me", api.ProfileInput{Cohort: "34", Nickname: "alice2", RealName: "Alice Again"}).Status(201))
	Equal(t, fresh.Status, api.MemberApproved)
	if fresh.ID == oldID {
		t.Fatal("onboarding reused the deleted Member")
	}
	Equal(t, len(As[[]api.Order](again.GET("/api/orders").Status(200))), 0)
}

func TestAdminEditsMemberProfile(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	id := alice.Me().Member.ID

	got := As[api.Member](admin.PUT(memberPath(id), api.ProfileInput{Cohort: "36", Nickname: "alice", RealName: "Alice Fixed", Phone: "0822222222"}).Status(200))
	Equal(t, got.Cohort, "36")
	Equal(t, got.RealName, "Alice Fixed")
	Equal(t, alice.Me().Member.Phone, "0822222222")

	admin.PUT(memberPath(id), api.ProfileInput{Cohort: "36", Nickname: "alice"}).Status(400)
	admin.PUT(memberPath(9999), api.ProfileInput{Cohort: "36", Nickname: "x", RealName: "x"}).Status(404)
	alice.PUT(memberPath(id), api.ProfileInput{Cohort: "36", Nickname: "alice", RealName: "x"}).Status(403)
}

func TestRoles(t *testing.T) {
	h := New(t)
	bob, admin := h.Member("bob"), h.Admin()
	bobID := bob.Me().Member.ID
	rolePath := memberPath(bobID) + "/role"

	bob.PUT(rolePath, api.RoleInput{Role: api.RoleAdmin}).Status(403)
	bob.GET("/api/admin/members").Status(403)

	got := As[api.Member](admin.PUT(rolePath, api.RoleInput{Role: api.RoleAdmin}).Status(200))
	Equal(t, got.Role, api.RoleAdmin)
	bob.GET("/api/admin/members").Status(200)

	got = As[api.Member](admin.PUT(rolePath, api.RoleInput{Role: api.RoleMember}).Status(200))
	Equal(t, got.Role, api.RoleMember)
	bob.GET("/api/admin/members").Status(403)

	admin.PUT(rolePath, api.RoleInput{Role: "owner"}).Status(400)
}
