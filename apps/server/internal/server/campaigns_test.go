package server_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

func titles(cs []api.Campaign) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.Title
	}
	return out
}

func TestDraftIsInvisibleToMembersAndVisibleToAdmin(t *testing.T) {
	h := New(t)
	admin, bob := h.Admin(), h.Member("bob")
	d := As[api.Campaign](admin.POST("/api/admin/campaigns", draft("ร่าง")).Status(201))
	pub := As[api.Campaign](admin.POST("/api/admin/campaigns", api.CampaignInput{Type: api.CampaignEvent, Title: "งานเลี้ยง", Status: api.CampaignOpen}).Status(201))

	Equal(t, fmt.Sprint(titles(As[[]api.Campaign](bob.GET("/api/campaigns").Status(200)))), "[งานเลี้ยง]")
	bob.GET(fmt.Sprintf("/api/campaigns/%d", d.ID)).Status(404)
	bob.GET(fmt.Sprintf("/api/campaigns/%d", pub.ID)).Status(200)

	Equal(t, fmt.Sprint(titles(As[[]api.Campaign](admin.GET("/api/admin/campaigns").Status(200)))), "[งานเลี้ยง ร่าง]")
	Equal(t, As[api.CampaignDetail](admin.GET(fmt.Sprintf("/api/admin/campaigns/%d", d.ID)).Status(200)).Campaign.Status, api.CampaignDraft)
}

func TestPublishingADraftShowsItToMembers(t *testing.T) {
	h := New(t)
	admin, bob := h.Admin(), h.Member("bob")
	d := As[api.Campaign](admin.POST("/api/admin/campaigns", draft("ร่าง")).Status(201))

	in := draft("ร่าง")
	in.Status = api.CampaignOpen
	admin.PUT(fmt.Sprintf("/api/admin/campaigns/%d", d.ID), in).Status(200)

	Equal(t, len(As[[]api.Campaign](bob.GET("/api/campaigns").Status(200))), 1)
}

func TestMemberAndApplicantAreRefusedWhereTheyShouldBe(t *testing.T) {
	h := New(t)
	j := h.OpenJersey()
	bob, nok := h.Member("bob"), h.Applicant("nok")

	Equal(t, bob.GET("/api/admin/campaigns").Status(403).Error(), "สำหรับแอดมินเท่านั้น")
	bob.POST("/api/admin/campaigns", draft("x")).Status(403)
	bob.PUT(fmt.Sprintf("/api/admin/campaigns/%d", j.Campaign.ID), draft("x")).Status(403)
	bob.POST(fmt.Sprintf("/api/admin/campaigns/%d/options", j.Campaign.ID), api.OptionInput{Name: "x"}).Status(403)
	bob.PUT(fmt.Sprintf("/api/admin/options/%d", j.Option.ID), api.OptionInput{Name: "x"}).Status(403)
	bob.DELETE(fmt.Sprintf("/api/admin/options/%d", j.Option.ID)).Status(403)
	bob.POST("/api/admin/images", JPEG()).Status(403)

	Equal(t, nok.GET("/api/campaigns").Status(403).Error(), "รอแอดมินอนุมัติ")
	nok.GET(fmt.Sprintf("/api/campaigns/%d", j.Campaign.ID)).Status(403)
	h.Anonymous().GET("/api/campaigns").Status(401)
}

func TestCampaignReadsClosedOncePastCloseAtWithNoWrite(t *testing.T) {
	h := New(t)
	j := h.OpenJersey()
	bob, admin := h.Member("bob"), h.Admin()
	detail := fmt.Sprintf("/api/campaigns/%d", j.Campaign.ID)
	Equal(t, As[api.CampaignDetail](bob.GET(detail).Status(200)).Campaign.Status, api.CampaignOpen)

	h.Advance(25 * time.Hour)

	Equal(t, As[api.CampaignDetail](bob.GET(detail).Status(200)).Campaign.Status, api.CampaignClosed)
	Equal(t, As[[]api.Campaign](bob.GET("/api/campaigns").Status(200))[0].Status, api.CampaignClosed)
	Equal(t, As[[]api.Campaign](admin.GET("/api/admin/campaigns").Status(200))[0].Status, api.CampaignClosed)
}

func TestAdminCloseSetsClosedBeforeCloseAt(t *testing.T) {
	h := New(t)
	j := h.OpenJersey()

	h.Close(j.Campaign.ID)

	Equal(t, As[api.CampaignDetail](h.Member("bob").GET(fmt.Sprintf("/api/campaigns/%d", j.Campaign.ID)).Status(200)).Campaign.Status, api.CampaignClosed)
}

func TestReopenWithLaterCloseAtReadsOpen(t *testing.T) {
	h := New(t)
	j := h.OpenJersey()
	h.Advance(25 * time.Hour)
	later := h.clock().Add(time.Hour).Format(time.RFC3339)
	in := draft(j.Campaign.Title)
	in.Type, in.Status, in.CloseAt = api.CampaignMerch, api.CampaignOpen, &later

	got := As[api.Campaign](h.Admin().PUT(fmt.Sprintf("/api/admin/campaigns/%d", j.Campaign.ID), in).Status(200))

	Equal(t, got.Status, api.CampaignOpen)
}

func TestCampaignsListNewestFirst(t *testing.T) {
	h := New(t)
	admin := h.Admin()
	for _, title := range []string{"แรก", "สอง", "สาม"} {
		in := draft(title)
		in.Status = api.CampaignOpen
		admin.POST("/api/admin/campaigns", in).Status(201)
	}

	got := titles(As[[]api.Campaign](h.Member("bob").GET("/api/campaigns").Status(200)))

	Equal(t, fmt.Sprint(got), "[สาม สอง แรก]")
}

func TestCampaignValidation(t *testing.T) {
	h := New(t)
	admin := h.Admin()
	badClose := "พรุ่งนี้"
	cases := []struct {
		name string
		edit func(*api.CampaignInput)
		want string
	}{
		{"blank title", func(in *api.CampaignInput) { in.Title = "  " }, "กรุณากรอกชื่อแคมเปญ"},
		{"unknown type", func(in *api.CampaignInput) { in.Type = "raffle" }, "ประเภทแคมเปญไม่ถูกต้อง"},
		{"unknown status", func(in *api.CampaignInput) { in.Status = "archived" }, "สถานะแคมเปญไม่ถูกต้อง"},
		{"bad closeAt", func(in *api.CampaignInput) { in.CloseAt = &badClose }, "เวลาปิดรับไม่ถูกต้อง"},
	}
	for _, tc := range cases {
		in := draft("ชื่อ")
		tc.edit(&in)
		Equal(t, admin.POST("/api/admin/campaigns", in).Status(400).Error(), tc.want)
	}
	admin.POST("/api/admin/campaigns", "not an object").Status(400)
	Equal(t, len(As[[]api.Campaign](admin.GET("/api/admin/campaigns").Status(200))), 0)
}

func TestCloseAtIsStoredAsUTC(t *testing.T) {
	h := New(t)
	in := draft("ชื่อ")
	local := "2026-12-01T09:00:00+07:00"
	in.CloseAt = &local

	got := As[api.Campaign](h.Admin().POST("/api/admin/campaigns", in).Status(201))

	Equal(t, *got.CloseAt, "2026-12-01T02:00:00Z")
}

func TestUnknownCampaignIs404(t *testing.T) {
	h := New(t)
	admin, bob := h.Admin(), h.Member("bob")

	bob.GET("/api/campaigns/999").Status(404)
	bob.GET("/api/campaigns/abc").Status(404)
	admin.PUT("/api/admin/campaigns/999", draft("x")).Status(404)
	admin.POST("/api/admin/campaigns/999/options", api.OptionInput{Name: "x"}).Status(404)
}

func TestOptionValidation(t *testing.T) {
	h := New(t)
	j := h.OpenJersey()
	admin := h.Admin()
	zero, neg := int64(0), int64(-1)
	cases := []struct {
		name string
		edit func(*api.OptionInput)
		want string
	}{
		{"blank name", func(in *api.OptionInput) { in.Name = " " }, "กรุณากรอกชื่อตัวเลือก"},
		{"zero price", func(in *api.OptionInput) { in.Price = &zero }, "ราคาต้องมากกว่า 0"},
		{"negative sortOrder", func(in *api.OptionInput) { in.SortOrder = neg }, "ลำดับต้องไม่ติดลบ"},
		{"duplicate size", func(in *api.OptionInput) { in.Sizes = []string{"S", " S "} }, "ไซซ์ต้องไม่ว่างและไม่ซ้ำกัน"},
		{"blank size", func(in *api.OptionInput) { in.Sizes = []string{"S", ""} }, "ไซซ์ต้องไม่ว่างและไม่ซ้ำกัน"},
	}
	for _, tc := range cases {
		in := api.OptionInput{Name: "ตัวเลือก"}
		tc.edit(&in)
		Equal(t, admin.POST(fmt.Sprintf("/api/admin/campaigns/%d/options", j.Campaign.ID), in).Status(400).Error(), tc.want)
		admin.PUT(fmt.Sprintf("/api/admin/options/%d", j.Option.ID), in).Status(400)
	}
}

func TestOptionUpdateReorderAndDelete(t *testing.T) {
	h := New(t)
	j := h.OpenJersey()
	admin, bob := h.Admin(), h.Member("bob")
	create := fmt.Sprintf("/api/admin/campaigns/%d/options", j.Campaign.ID)
	donation := As[api.Option](admin.POST(create, api.OptionInput{Name: "บริจาค", SortOrder: 5}).Status(201))
	Equal(t, donation.Price == nil, true)
	Equal(t, len(donation.Sizes), 0)

	// Moving the jersey behind the donation changes the order Members see.
	price := int64(300)
	updated := As[api.Option](admin.PUT(fmt.Sprintf("/api/admin/options/%d", j.Option.ID), api.OptionInput{
		Name: " เสื้อรุ่น ", Price: &price, Sizes: []string{" M ", "L"}, NeedsPrintName: true, SortOrder: 9,
	}).Status(200))
	Equal(t, updated.Name, "เสื้อรุ่น")
	Equal(t, *updated.Price, int64(300))
	Equal(t, fmt.Sprint(updated.Sizes), "[M L]")
	Equal(t, updated.NeedsNumber, false)
	Equal(t, updated.CampaignID, j.Campaign.ID)

	opts := As[api.CampaignDetail](bob.GET(fmt.Sprintf("/api/campaigns/%d", j.Campaign.ID)).Status(200)).Options
	Equal(t, len(opts), 2)
	Equal(t, opts[0].Name, "บริจาค")
	Equal(t, opts[1].Name, "เสื้อรุ่น")

	del := fmt.Sprintf("/api/admin/options/%d", donation.ID)
	admin.DELETE(del).Status(204)
	admin.DELETE(del).Status(404)
	admin.PUT(del, api.OptionInput{Name: "x"}).Status(404)
	Equal(t, len(As[api.CampaignDetail](bob.GET(fmt.Sprintf("/api/campaigns/%d", j.Campaign.ID)).Status(200)).Options), 1)
}

func TestDraftOptionsAreNotLeakedToMembers(t *testing.T) {
	h := New(t)
	admin := h.Admin()
	d := As[api.Campaign](admin.POST("/api/admin/campaigns", draft("ร่าง")).Status(201))
	admin.POST(fmt.Sprintf("/api/admin/campaigns/%d/options", d.ID), api.OptionInput{Name: "ลับ"}).Status(201)

	h.Member("bob").GET(fmt.Sprintf("/api/campaigns/%d", d.ID)).Status(404)
	Equal(t, len(As[api.CampaignDetail](admin.GET(fmt.Sprintf("/api/admin/campaigns/%d", d.ID)).Status(200)).Options), 1)
}
