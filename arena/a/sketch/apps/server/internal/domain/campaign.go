package domain

import (
	"time"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

// Campaign is the part of a Campaign that rules depend on.
type Campaign struct {
	ID      int64
	Stored  api.CampaignStatus
	CloseAt *time.Time
}

// Status is the ONLY place effective Campaign status is computed:
// stored open with CloseAt <= now reads closed. Draft and closed read as stored.
func (c Campaign) Status(now time.Time) api.CampaignStatus {
	panic("not implemented")
}

// VisibleTo: Members never see drafts (404, not 403, so drafts are not discoverable).
func (c Campaign) VisibleTo(a Actor) bool {
	panic("not implemented")
}

// CanCheckout: a Member needs the Campaign effectively open; an Admin ordering on
// someone's behalf may order in any non-draft Campaign (story 68).
func (c Campaign) CanCheckout(a Actor, onBehalf bool, now time.Time) error {
	panic("not implemented")
}

// ValidateCampaign checks title, type, status and parses CloseAt.
// Saving status open with CloseAt already past is refused: it would read closed at once,
// and the Admin almost certainly meant to reopen (clear or move the close date).
func ValidateCampaign(in api.CampaignInput, now time.Time) (closeAt *time.Time, err error) {
	panic("not implemented")
}

// ValidateOption checks name, price >= 0 when set, sizes trimmed/unique/non-empty strings.
func ValidateOption(in api.OptionInput) (sizesJSON string, err error) {
	panic("not implemented")
}

// ParseSizes decodes options.sizes; the column CHECK guarantees a JSON array.
func ParseSizes(sizesJSON string) []string {
	panic("not implemented")
}
