package server_test

import (
	"fmt"
	"time"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

// Jersey is an open merch Campaign with one 260-baht Option, sizes S..XL, print name + number.
type Jersey struct {
	api.CampaignDetail
	Option api.Option
}

// OpenJersey has the Admin create the jersey Campaign through the API; it closes 24h after the clock's now.
func (h *H) OpenJersey() Jersey {
	admin := h.Admin()
	closeAt := h.clock().Add(24 * time.Hour).UTC().Format(time.RFC3339)
	camp := As[api.Campaign](admin.POST("/api/admin/campaigns", api.CampaignInput{
		Type:            api.CampaignMerch,
		Title:           "เสื้อรุ่น 2569",
		Description:     "สั่งเสื้อรุ่น",
		BankCode:        "004",
		BankAccountNo:   "000-0-00000-0",
		BankAccountName: "บัญชีตัวอย่าง",
		Status:          api.CampaignOpen,
		CloseAt:         &closeAt,
	}).Status(201))
	price := int64(260)
	opt := As[api.Option](admin.POST(fmt.Sprintf("/api/admin/campaigns/%d/options", camp.ID), api.OptionInput{
		Name:           "เสื้อรุ่น",
		Price:          &price,
		Sizes:          []string{"S", "M", "L", "XL"},
		NeedsPrintName: true,
		NeedsNumber:    true,
	}).Status(201))
	return Jersey{CampaignDetail: api.CampaignDetail{Campaign: camp, Options: []api.Option{opt}}, Option: opt}
}

// Close sets the stored status to closed (Admin action, not the clock).
func (h *H) Close(campaignID int64) {
	admin := h.Admin()
	path := fmt.Sprintf("/api/admin/campaigns/%d", campaignID)
	c := As[api.CampaignDetail](admin.GET(path).Status(200)).Campaign
	admin.PUT(path, api.CampaignInput{
		Type:            c.Type,
		Title:           c.Title,
		Description:     c.Description,
		Photos:          c.Photos,
		BankCode:        c.BankCode,
		BankAccountNo:   c.BankAccountNo,
		BankAccountName: c.BankAccountName,
		QRImage:         c.QRImage,
		Status:          api.CampaignClosed,
		CloseAt:         c.CloseAt,
	}).Status(200)
}

// draft is the smallest valid CampaignInput.
func draft(title string) api.CampaignInput {
	return api.CampaignInput{Type: api.CampaignDonation, Title: title, Status: api.CampaignDraft}
}
