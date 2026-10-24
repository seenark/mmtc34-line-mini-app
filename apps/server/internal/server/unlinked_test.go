package server_test

import (
	"fmt"
	"testing"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

func linkPath(memberID int64) string { return fmt.Sprintf("/api/admin/members/%d/link", memberID) }

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

func TestLinkingAnApplicantGivesThemTheUnlinkedMembersOrders(t *testing.T) {
	h := New(t)
	admin := h.Admin()
	j := h.OpenJersey()
	u := As[api.Member](admin.POST("/api/admin/members", api.ProfileInput{Nickname: "ติ๋ว", Phone: "0899999999"}).Status(201))
	o := As[api.Order](admin.POST(adminOrdersPath(j.Campaign.ID), api.AdminOrderInput{
		MemberID: u.ID, Order: Pickup(Shirt(j, "M", 2)),
	}).Status(201))
	tiw := h.Applicant("tiw")
	applicantID := tiw.Me().Member.ID
	tiw.GET("/api/orders").Status(403)

	linked := As[api.Member](admin.POST(linkPath(u.ID), api.LinkInput{ApplicantID: applicantID}).Status(200))
	Equal(t, linked.ID, u.ID)
	Equal(t, linked.Linked, true)
	Equal(t, linked.Status, api.MemberApproved)
	Equal(t, linked.Nickname, "ติ๋ว")
	Equal(t, linked.RealName, "Real tiw")
	Equal(t, linked.Cohort, "34")
	Equal(t, linked.Phone, "0899999999")

	me := tiw.Me().Member
	Equal(t, me.ID, u.ID)
	Equal(t, me.Linked, true)
	Equal(t, me.Status, api.MemberApproved)
	orders := As[[]api.Order](tiw.GET("/api/orders").Status(200))
	Equal(t, len(orders), 1)
	Equal(t, orders[0].ID, o.ID)
	admin.GET(fmt.Sprintf("/api/admin/members/%d", applicantID)).Status(404)
}

func TestLinkRejectsLinkedTargetsAndNonApplicants(t *testing.T) {
	h := New(t)
	admin, alice := h.Admin(), h.Member("alice")
	u := As[api.Member](admin.POST("/api/admin/members", api.ProfileInput{Nickname: "ติ๋ว"}).Status(201))
	nok := h.Applicant("nok").Me().Member
	aliceID := alice.Me().Member.ID

	alice.POST(linkPath(u.ID), api.LinkInput{ApplicantID: nok.ID}).Status(403)
	admin.POST(linkPath(aliceID), api.LinkInput{ApplicantID: nok.ID}).Status(409)
	admin.POST(linkPath(9999), api.LinkInput{ApplicantID: nok.ID}).Status(404)
	admin.POST(linkPath(u.ID), api.LinkInput{ApplicantID: aliceID}).Status(400)
	admin.POST(linkPath(u.ID), api.LinkInput{ApplicantID: 9999}).Status(400)

	admin.POST(linkPath(u.ID), api.LinkInput{ApplicantID: nok.ID}).Status(200)
	other := h.Applicant("other").Me().Member
	admin.POST(linkPath(u.ID), api.LinkInput{ApplicantID: other.ID}).Status(409)
	Equal(t, As[api.Member](admin.GET(fmt.Sprintf("/api/admin/members/%d", other.ID)).Status(200)).Status, api.MemberApplicant)
}
