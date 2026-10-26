package domain

import (
	"strings"
	"time"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

// Campaign is the part of a Campaign that rules depend on.
type Campaign struct {
	Stored  api.CampaignStatus
	CloseAt *time.Time
}

// Status is the ONLY place effective Campaign status is computed: stored open with
// CloseAt <= now reads closed. Draft and closed read as stored. No job flips the stored value.
func (c Campaign) Status(now time.Time) api.CampaignStatus {
	if c.Stored == api.CampaignOpen && c.CloseAt != nil && !c.CloseAt.After(now) {
		return api.CampaignClosed
	}
	return c.Stored
}

// VisibleTo: only Admins see drafts. Callers answer 404, not 403, so drafts are not discoverable.
func (c Campaign) VisibleTo(a Actor) bool {
	return a.Admin || c.Stored != api.CampaignDraft
}

// CanCheckout: a Member needs the Campaign effectively open; an Admin ordering on a Member's
// behalf may order in any non-draft Campaign, closed included.
func (c Campaign) CanCheckout(actorIsAdmin bool, now time.Time) error {
	switch {
	case c.Stored == api.CampaignDraft:
		return &Error{Kind: Conflict, Msg: "อีเวนต์ยังไม่เปิดรับ"}
	case !actorIsAdmin && c.Status(now) != api.CampaignOpen:
		return &Error{Kind: Conflict, Msg: "อีเวนต์ปิดรับแล้ว"}
	}
	return nil
}

// ValidateCampaign trims text fields, checks title, type, status and bank code, and normalises
// CloseAt to UTC RFC3339.
func ValidateCampaign(in api.CampaignInput) (api.CampaignInput, error) {
	in.Title = strings.TrimSpace(in.Title)
	in.Description = strings.TrimSpace(in.Description)
	in.BankCode = strings.TrimSpace(in.BankCode)
	in.BankAccountNo = strings.TrimSpace(in.BankAccountNo)
	in.BankAccountName = strings.TrimSpace(in.BankAccountName)
	switch {
	case in.Title == "":
		return in, &Error{Kind: Invalid, Msg: "กรุณากรอกชื่ออีเวนต์"}
	case in.Type != api.CampaignMerch && in.Type != api.CampaignDonation && in.Type != api.CampaignEvent:
		return in, &Error{Kind: Invalid, Msg: "ประเภทอีเวนต์ไม่ถูกต้อง"}
	case in.Status != api.CampaignDraft && in.Status != api.CampaignOpen && in.Status != api.CampaignClosed:
		return in, &Error{Kind: Invalid, Msg: "สถานะอีเวนต์ไม่ถูกต้อง"}
	case in.BankCode != "" && BankName(in.BankCode) == "":
		return in, &Error{Kind: Invalid, Msg: "ธนาคารไม่ถูกต้อง"}
	}
	if in.CloseAt != nil {
		t, err := time.Parse(time.RFC3339, *in.CloseAt)
		if err != nil {
			return in, &Error{Kind: Invalid, Msg: "เวลาปิดรับไม่ถูกต้อง"}
		}
		s := t.UTC().Format(time.RFC3339)
		in.CloseAt = &s
	}
	return in, nil
}

// ValidateOption trims name and description; a fixed price must be positive (nil = Member
// chooses); sizes are trimmed, non-empty and unique; sortOrder is not negative.
func ValidateOption(in api.OptionInput) (api.OptionInput, error) {
	in.Name = strings.TrimSpace(in.Name)
	in.Description = strings.TrimSpace(in.Description)
	switch {
	case in.Name == "":
		return in, &Error{Kind: Invalid, Msg: "กรุณากรอกชื่อตัวเลือก"}
	case in.Price != nil && *in.Price <= 0:
		return in, &Error{Kind: Invalid, Msg: "ราคาต้องมากกว่า 0"}
	case in.SortOrder < 0:
		return in, &Error{Kind: Invalid, Msg: "ลำดับต้องไม่ติดลบ"}
	}
	sizes := make([]string, 0, len(in.Sizes))
	seen := map[string]bool{}
	for _, sz := range in.Sizes {
		sz = strings.TrimSpace(sz)
		if sz == "" || seen[sz] {
			return in, &Error{Kind: Invalid, Msg: "ไซซ์ต้องไม่ว่างและไม่ซ้ำกัน"}
		}
		seen[sz] = true
		sizes = append(sizes, sz)
	}
	in.Sizes = sizes
	return in, nil
}
