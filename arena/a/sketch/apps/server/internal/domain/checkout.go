package domain

import "github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"

// OptionRule is what checkout needs from an Option row.
type OptionRule struct {
	ID, CampaignID int64
	Name           string
	Price          *int64 // nil = Member chooses the amount
	Sizes          []string
	NeedsPrintName bool
	NeedsNumber    bool
}

// ValidItem has its unit price captured: the Option price at checkout, the donor's amount,
// or (Admin edit) the price already stored on the Item with the same ID.
type ValidItem struct {
	OptionID                int64
	Size, PrintName, Number string
	Quantity, UnitPrice     int64
}

// ValidOrder can only be built by ValidateOrder, so the insert path cannot skip validation.
type ValidOrder struct {
	delivery                         api.Delivery
	shipName, shipAddress, shipPhone string
	items                            []ValidItem
}

func (v ValidOrder) Delivery() api.Delivery { return v.delivery }

// Ship returns the shipping snapshot; all "" for pickup (pickup discards typed address).
func (v ValidOrder) Ship() (name, address, phone string) {
	return v.shipName, v.shipAddress, v.shipPhone
}

func (v ValidOrder) Items() []ValidItem { return v.items }

// ValidateOrder is the checkout rule set (spec "Checkout validation"), shared by Member
// checkout (.4), Admin on-behalf checkout and Admin edit (.7):
//   - at least one Item; every OptionID is a key of options (i.e. belongs to the Campaign);
//   - size ∈ Option sizes when the list is non-empty, else must be "";
//   - print name / number present when the Option needs them;
//   - quantity 1..99;
//   - Amount > 0 given exactly for no-price Options, absent otherwise;
//   - ship requires name, address and phone (trimmed non-empty).
// keep holds the Order's current Items by ID (Admin edit); nil for new Orders. An input Item
// whose ID is in keep keeps keep[ID].UnitPrice; an unknown ID is Invalid.
// The first failure is returned with its 1-based Item index in the Thai message.
func ValidateOrder(in api.OrderInput, options map[int64]OptionRule, keep map[int64]Item) (ValidOrder, error) {
	panic("not implemented")
}
