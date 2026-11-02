package server_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

func openCampaign(title string, photos ...string) api.CampaignInput {
	in := draft(title)
	in.Status = api.CampaignOpen
	in.Photos = photos
	return in
}

func uploadJPEGs(t *testing.T, admin *Client, n int) []string {
	t.Helper()
	names := make([]string, n)
	for i := range names {
		names[i] = uploadJPEG(t, admin)
	}
	return names
}

func memberPhotos(t *testing.T, member *Client, id int64) (detail, list string) {
	t.Helper()
	d := As[api.CampaignDetail](member.GET(fmt.Sprintf("/api/campaigns/%d", id)).Status(200))
	for _, c := range As[[]api.Campaign](member.GET("/api/campaigns").Status(200)) {
		if c.ID == id {
			return fmt.Sprint(d.Campaign.Photos), fmt.Sprint(c.Photos)
		}
	}
	t.Fatalf("campaign %d missing from the Member list", id)
	return "", ""
}

func TestCampaignCreatedWithPhotosKeepsTheirOrderForMembers(t *testing.T) {
	h := New(t)
	admin, bob := h.Admin(), h.Member("bob")
	names := uploadJPEGs(t, admin, 3)
	order := []string{names[2], names[0], names[1]}

	camp := As[api.Campaign](admin.POST("/api/admin/campaigns", openCampaign("เสื้อ", order...)).Status(201))

	want := fmt.Sprint(order)
	Equal(t, fmt.Sprint(camp.Photos), want)
	detail, list := memberPhotos(t, bob, camp.ID)
	Equal(t, detail, want)
	Equal(t, list, want)
}

func TestAdminListsCampaignPhotosInOrder(t *testing.T) {
	h := New(t)
	admin := h.Admin()
	names := uploadJPEGs(t, admin, 2)
	a := As[api.Campaign](admin.POST("/api/admin/campaigns", openCampaign("เอ", names[1], names[0])).Status(201))
	b := As[api.Campaign](admin.POST("/api/admin/campaigns", openCampaign("บี", names[0])).Status(201))

	byID := map[int64]string{}
	for _, c := range As[[]api.Campaign](admin.GET("/api/admin/campaigns").Status(200)) {
		byID[c.ID] = fmt.Sprint(c.Photos)
	}

	Equal(t, byID[a.ID], fmt.Sprint([]string{names[1], names[0]}))
	Equal(t, byID[b.ID], fmt.Sprint([]string{names[0]}))
}

func TestCampaignWithoutPhotosServesAnEmptyArray(t *testing.T) {
	h := New(t)
	admin := h.Admin()

	created := admin.POST("/api/admin/campaigns", openCampaign("ไม่มีรูป")).Status(201)
	id := As[api.Campaign](created).ID

	Equal(t, bytes.Contains(created.Body, []byte(`"photos":[]`)), true)
	Equal(t, bytes.Contains(admin.GET("/api/admin/campaigns").Status(200).Body, []byte(`"photos":[]`)), true)
	Equal(t, bytes.Contains(h.Member("bob").GET(fmt.Sprintf("/api/campaigns/%d", id)).Status(200).Body, []byte(`"photos":[]`)), true)
}

func TestCampaignSavedWithNullPhotosHasNone(t *testing.T) {
	h := New(t)
	admin := h.Admin()
	name := uploadJPEG(t, admin)
	camp := As[api.Campaign](admin.POST("/api/admin/campaigns", openCampaign("เสื้อ", name)).Status(201))

	r := admin.PUT(fmt.Sprintf("/api/admin/campaigns/%d", camp.ID), map[string]any{
		"type": "event", "title": "เสื้อ", "status": "open", "photos": nil,
	}).Status(200)

	Equal(t, bytes.Contains(r.Body, []byte(`"photos":[]`)), true)
}

func TestUpdatingPhotosReplacesTheWholeSet(t *testing.T) {
	h := New(t)
	admin, bob := h.Admin(), h.Member("bob")
	names := uploadJPEGs(t, admin, 3)
	camp := As[api.Campaign](admin.POST("/api/admin/campaigns", openCampaign("เสื้อ", names...)).Status(201))
	path := fmt.Sprintf("/api/admin/campaigns/%d", camp.ID)

	shorter := As[api.Campaign](admin.PUT(path, openCampaign("เสื้อ", names[1])).Status(200))
	detail, list := memberPhotos(t, bob, camp.ID)
	Equal(t, fmt.Sprint(shorter.Photos), fmt.Sprint([]string{names[1]}))
	Equal(t, detail, fmt.Sprint([]string{names[1]}))
	Equal(t, list, fmt.Sprint([]string{names[1]}))

	empty := admin.PUT(path, openCampaign("เสื้อ")).Status(200)
	detail, list = memberPhotos(t, bob, camp.ID)
	Equal(t, bytes.Contains(empty.Body, []byte(`"photos":[]`)), true)
	Equal(t, detail, "[]")
	Equal(t, list, "[]")
}

func TestReorderingPhotosChangesTheCover(t *testing.T) {
	h := New(t)
	admin, bob := h.Admin(), h.Member("bob")
	names := uploadJPEGs(t, admin, 3)
	camp := As[api.Campaign](admin.POST("/api/admin/campaigns", openCampaign("เสื้อ", names...)).Status(201))
	reordered := []string{names[2], names[0], names[1]}

	updated := As[api.Campaign](admin.PUT(fmt.Sprintf("/api/admin/campaigns/%d", camp.ID), openCampaign("เสื้อ", reordered...)).Status(200))

	want := fmt.Sprint(reordered)
	Equal(t, fmt.Sprint(updated.Photos), want)
	detail, list := memberPhotos(t, bob, camp.ID)
	Equal(t, detail, want)
	Equal(t, list, want)
}

func TestRemovedPhotoFileStaysReadable(t *testing.T) {
	h := New(t)
	admin, bob := h.Admin(), h.Member("bob")
	name := uploadJPEG(t, admin)
	camp := As[api.Campaign](admin.POST("/api/admin/campaigns", openCampaign("เสื้อ", name)).Status(201))

	admin.PUT(fmt.Sprintf("/api/admin/campaigns/%d", camp.ID), openCampaign("เสื้อ")).Status(200)

	bob.GET("/api/images/" + name).Status(200)
}

func TestUnknownPhotoNameRejectsTheSaveAndKeepsStoredPhotos(t *testing.T) {
	h := New(t)
	admin, bob := h.Admin(), h.Member("bob")
	names := uploadJPEGs(t, admin, 2)
	camp := As[api.Campaign](admin.POST("/api/admin/campaigns", openCampaign("เสื้อ", names...)).Status(201))
	unknown := "0123456789abcdef0123456789abcdef.jpg"
	path := fmt.Sprintf("/api/admin/campaigns/%d", camp.ID)

	r := admin.PUT(path, openCampaign("เปลี่ยนชื่อ", names[0], unknown, names[1])).Status(400)

	Equal(t, r.Error(), "ไม่พบรูปภาพ กรุณาอัปโหลดใหม่")
	admin.PUT(path, openCampaign("เสื้อ", "")).Status(400)
	detail, _ := memberPhotos(t, bob, camp.ID)
	Equal(t, detail, fmt.Sprint(names))
	Equal(t, As[api.CampaignDetail](bob.GET(fmt.Sprintf("/api/campaigns/%d", camp.ID)).Status(200)).Campaign.Title, "เสื้อ")
}

func TestUnknownPhotoNameOnCreateSavesNothing(t *testing.T) {
	h := New(t)
	admin := h.Admin()
	name := uploadJPEG(t, admin)

	admin.POST("/api/admin/campaigns", openCampaign("เสื้อ", name, "0123456789abcdef0123456789abcdef.jpg")).Status(400)

	Equal(t, len(As[[]api.Campaign](admin.GET("/api/admin/campaigns").Status(200))), 0)
}

func TestDuplicatePhotoNameIsRejected(t *testing.T) {
	h := New(t)
	admin, bob := h.Admin(), h.Member("bob")
	names := uploadJPEGs(t, admin, 2)
	camp := As[api.Campaign](admin.POST("/api/admin/campaigns", openCampaign("เสื้อ", names[0])).Status(201))
	dup := openCampaign("เสื้อ", names[0], names[1], names[0])

	Equal(t, admin.POST("/api/admin/campaigns", dup).Status(400).Error(), "มีรูปภาพซ้ำกัน")
	Equal(t, admin.PUT(fmt.Sprintf("/api/admin/campaigns/%d", camp.ID), dup).Status(400).Error(), "มีรูปภาพซ้ำกัน")

	detail, _ := memberPhotos(t, bob, camp.ID)
	Equal(t, detail, fmt.Sprint([]string{names[0]}))
	Equal(t, len(As[[]api.Campaign](admin.GET("/api/admin/campaigns").Status(200))), 1)
}

func TestUpdatingAMissingCampaignWithPhotosIs404(t *testing.T) {
	h := New(t)
	admin := h.Admin()
	name := uploadJPEG(t, admin)

	admin.PUT("/api/admin/campaigns/999", openCampaign("เสื้อ", name)).Status(404)
}
