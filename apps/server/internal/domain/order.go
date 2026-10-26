package domain

import (
	"slices"
	"time"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

// Order is everything the rules need about one Order, assembled by server.loadOrders.
// Timestamps stay in the RFC3339 form they are stored in; only Now is a time.Time.
type Order struct {
	ID            int64
	CampaignID    int64
	CampaignTitle string
	Campaign      Campaign
	MemberID      int64
	Nickname      string
	RealName      string // exports only; never copied into api.Order
	PlacedByAdmin bool
	Stored        api.OrderStatus
	Delivery      api.Delivery
	ShipName      string
	ShipAddress   string
	ShipPhone     string
	Items         []Item
	Payments      []Payment
	Adjustments   []Adjustment
	CreatedAt     string
	Now           time.Time // the request's clock reading, captured once by the loader
}

type Item struct {
	ID, OptionID            int64
	OptionName              string
	Size, PrintName, Number string
	Quantity, UnitPrice     int64
}

type Payment struct {
	ID           int64
	Slip         string // "" = cash
	Status       api.PaymentStatus
	Amount       *int64
	RejectReason string
	CreatedAt    string
	ReviewedBy   string
	ReviewedAt   *string
}

type Adjustment struct {
	ID        int64
	Amount    int64
	Note      string
	CreatedBy string
	CreatedAt string
}

// Status is the ONLY place effective Order status is computed:
// stored awaiting_payment + Campaign effectively closed + not placed by an Admin = cancelled.
// Reopening the Campaign makes such Orders read awaiting_payment again; nothing was stored.
func (o Order) Status() api.OrderStatus {
	if o.Stored == awaiting && !o.PlacedByAdmin && o.Campaign.Status(o.Now) == api.CampaignClosed {
		return cancelled
	}
	return o.Stored
}

// Event is a lifecycle action. Checkout is not an Event: it creates the Order (Campaign.CanCheckout).
type Event string

const (
	EvView        Event = "view"
	EvUploadSlip  Event = "upload_slip"
	EvConfirmSlip Event = "confirm_slip"
	EvRejectSlip  Event = "reject_slip"
	EvRecordCash  Event = "record_cash"
	EvCancel      Event = "cancel"
	EvEdit        Event = "edit"
	EvAdjust      Event = "adjust"
)

type party uint8

const (
	owner party = 1 << iota // the Order's Member
	admin
)

// rule: `by` may fire `ev` when the STORED status is in `from`. `ownerNeedsOpen` adds
// "Campaign effectively open"; Admins are never bound by close. `to` == "" leaves status unchanged.
type rule struct {
	ev             Event
	by             party
	from           []api.OrderStatus
	to             api.OrderStatus
	ownerNeedsOpen bool
}

var (
	awaiting  = api.OrderAwaitingPayment
	reviewing = api.OrderPendingReview
	paid      = api.OrderPaid
	cancelled = api.OrderCancelled
	anyStatus = []api.OrderStatus{awaiting, reviewing, paid, cancelled}
)

// rules IS the Order state machine (ADR 0001). The table owns only "may this happen, and to
// which status"; each handler writes the transition's side rows inside server.apply's transaction.
var rules = []rule{
	{EvView, owner | admin, anyStatus, "", false},
	{EvUploadSlip, owner, []api.OrderStatus{awaiting}, reviewing, true},
	// Admins upload the Slip a Member sent them another way, e.g. for an on-behalf Order.
	{EvUploadSlip, admin, []api.OrderStatus{awaiting}, reviewing, false},
	{EvConfirmSlip, admin, []api.OrderStatus{reviewing}, paid, false},
	{EvRejectSlip, admin, []api.OrderStatus{reviewing}, awaiting, false},
	// Cash settles an unpaid Order even after close, or tops up a paid one. Not from
	// pending_review: review the Slip first, so a pending Payment never outlives its review.
	{EvRecordCash, admin, []api.OrderStatus{awaiting, paid}, paid, false},
	{EvCancel, owner, []api.OrderStatus{awaiting}, cancelled, true},
	// Admin cancel is idempotent: cancelling a cancelled Order succeeds and changes nothing.
	{EvCancel, admin, anyStatus, cancelled, false},
	{EvEdit, admin, anyStatus, "", false},
	{EvAdjust, admin, anyStatus, "", false},
}

// Decide answers "may actor fire ev on this Order now, and what is the new STORED status".
// NotFound when the actor may not even view it (another Member's Order); Forbidden when no
// rule names the actor's party for ev; Conflict when a rule names the party but the status
// or the closed Campaign does not allow it.
func (o Order) Decide(a Actor, ev Event) (api.OrderStatus, error) {
	var p party
	if a.MemberID == o.MemberID {
		p |= owner
	}
	if a.Admin {
		p |= admin
	}
	if p == 0 {
		return "", &Error{Kind: NotFound, Msg: "ไม่พบข้อมูล"}
	}
	open := o.Campaign.Status(o.Now) == api.CampaignOpen
	named, closedBlocked := false, false
	for _, r := range rules {
		if r.ev != ev || r.by&p == 0 {
			continue
		}
		named = true
		if !slices.Contains(r.from, o.Stored) {
			continue
		}
		if r.ownerNeedsOpen && !open {
			closedBlocked = true
			continue
		}
		if r.to == "" {
			return o.Stored, nil
		}
		return r.to, nil
	}
	switch {
	case !named:
		return "", &Error{Kind: Forbidden, Msg: "คุณไม่มีสิทธิ์ทำรายการนี้"}
	case closedBlocked:
		return "", &Error{Kind: Conflict, Msg: "อีเวนต์ปิดแล้ว กรุณาติดต่อแอดมิน"}
	}
	return "", &Error{Kind: Conflict, Msg: "สถานะคำสั่งซื้อไม่รองรับการทำรายการนี้"}
}

// Actions evaluates Decide for every event, for api.Order.actions. Derived from rules, so the
// buttons in the UI can never disagree with the server.
func (o Order) Actions(a Actor) api.OrderActions {
	can := func(ev Event) bool { _, err := o.Decide(a, ev); return err == nil }
	return api.OrderActions{
		Cancel:     can(EvCancel),
		UploadSlip: can(EvUploadSlip),
		ReviewSlip: can(EvConfirmSlip) || can(EvRejectSlip),
		RecordCash: can(EvRecordCash),
		Edit:       can(EvEdit),
		AddAdjust:  can(EvAdjust),
	}
}
