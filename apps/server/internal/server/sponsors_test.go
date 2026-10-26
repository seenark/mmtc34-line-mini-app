package server_test

import (
	"fmt"
	"testing"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

func TestAdminManagesSponsorPositionsAndMemberReadsThemInOrder(t *testing.T) {
	h := New(t)
	j := h.OpenJersey()
	admin, bob := h.Admin(), h.Member("bob")
	list := fmt.Sprintf("/api/campaigns/%d/sponsors", j.Campaign.ID)
	create := fmt.Sprintf("/api/admin/campaigns/%d/sponsors", j.Campaign.ID)
	logo := As[api.ImageUpload](admin.POST("/api/admin/images", JPEG()).Status(201)).Name
	bobID := bob.Me().Member.ID

	back := As[api.SponsorPosition](admin.POST(create, api.SponsorInput{Name: "หลัง", SortOrder: 20}).Status(201))
	chest := As[api.SponsorPosition](admin.POST(create, api.SponsorInput{
		Name: " กลางอก ", SponsorName: "TK.LINE", Logo: logo, MemberID: &bobID, SortOrder: 10,
	}).Status(201))
	Equal(t, chest.Name, "กลางอก")
	Equal(t, chest.MemberNickname, "bob")

	got := As[[]api.SponsorPosition](bob.GET(list).Status(200))
	Equal(t, len(got), 2)
	Equal(t, got[0].ID, chest.ID)
	Equal(t, got[0].SponsorName, "TK.LINE")
	Equal(t, got[0].Logo, logo)
	Equal(t, got[0].MemberNickname, "bob")
	Equal(t, got[1].ID, back.ID)
	Equal(t, got[1].MemberNickname, "")

	updated := As[api.SponsorPosition](admin.PUT(fmt.Sprintf("/api/admin/sponsors/%d", back.ID), api.SponsorInput{
		Name: "หลัง", SponsorName: "CodeSook", SortOrder: 5,
	}).Status(200))
	Equal(t, updated.SponsorName, "CodeSook")
	got = As[[]api.SponsorPosition](bob.GET(list).Status(200))
	Equal(t, got[0].ID, back.ID)

	admin.DELETE(fmt.Sprintf("/api/admin/sponsors/%d", back.ID)).Status(204)
	admin.DELETE(fmt.Sprintf("/api/admin/sponsors/%d", back.ID)).Status(404)
	Equal(t, len(As[[]api.SponsorPosition](bob.GET(list).Status(200))), 1)
}

func TestMemberCannotEditSponsorPositions(t *testing.T) {
	h := New(t)
	j := h.OpenJersey()
	bob := h.Member("bob")
	p := As[api.SponsorPosition](h.Admin().POST(fmt.Sprintf("/api/admin/campaigns/%d/sponsors", j.Campaign.ID), api.SponsorInput{Name: "แขนขวา"}).Status(201))

	bob.POST(fmt.Sprintf("/api/admin/campaigns/%d/sponsors", j.Campaign.ID), api.SponsorInput{Name: "x"}).Status(403)
	bob.PUT(fmt.Sprintf("/api/admin/sponsors/%d", p.ID), api.SponsorInput{Name: "x"}).Status(403)
	bob.DELETE(fmt.Sprintf("/api/admin/sponsors/%d", p.ID)).Status(403)
	h.Visitor("nok").GET(fmt.Sprintf("/api/campaigns/%d/sponsors", j.Campaign.ID)).Status(403)
}

func TestSponsorPositionRejectsBadInput(t *testing.T) {
	h := New(t)
	j := h.OpenJersey()
	admin := h.Admin()
	create := fmt.Sprintf("/api/admin/campaigns/%d/sponsors", j.Campaign.ID)
	nok := h.Member("nok")
	nokID := nok.Me().Member.ID
	nok.DELETE("/api/me").Status(204)
	ghost := int64(9999)

	Equal(t, admin.POST(create, api.SponsorInput{Name: "  "}).Status(400).Error(), "กรุณากรอกชื่อตำแหน่ง")
	Equal(t, admin.POST(create, api.SponsorInput{Name: "x", Logo: "nope.jpg"}).Status(400).Error(), "ไม่พบรูปภาพ กรุณาอัปโหลดใหม่")
	Equal(t, admin.POST(create, api.SponsorInput{Name: "x", MemberID: &ghost}).Status(400).Error(), "ไม่พบสมาชิกที่เลือก")
	Equal(t, admin.POST(create, api.SponsorInput{Name: "x", MemberID: &nokID}).Status(400).Error(), "ไม่พบสมาชิกที่เลือก")
	admin.POST("/api/admin/campaigns/9999/sponsors", api.SponsorInput{Name: "x"}).Status(404)
	admin.PUT("/api/admin/sponsors/9999", api.SponsorInput{Name: "x"}).Status(404)
}

func TestSponsorPositionsOfDraftCampaignAreHiddenFromMembers(t *testing.T) {
	h := New(t)
	admin, bob := h.Admin(), h.Member("bob")
	d := As[api.Campaign](admin.POST("/api/admin/campaigns", draft("ร่าง")).Status(201))

	bob.GET(fmt.Sprintf("/api/campaigns/%d/sponsors", d.ID)).Status(404)
}
