package domain

import (
	"time"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

// Order is everything the rules need about one Order, assembled by server.loadOrders.
// Campaign.Stored/CloseAt come from the same row read, so Status() needs no extra lookup.
type Order struct {
	ID          int64
	Campaign    Campaign
	CampaignTitle string
	MemberID    int64
	Nickname    string
	RealName    string // exports (pickup CSV) only; never copied into api.Order
	PlacedBy    *int64 // Admin who placed it on behalf; nil = Member checkout
	Stored      api.OrderStatus
	Delivery    api.Delivery
	ShipName, ShipAddress, ShipPhone string
	Items       []Item
	Payments    []Payment
	Adjustments []Adjustment
	CreatedAt   time.Time
	Now         time.Time // the request's clock reading, captured once by the loader
}

type Item struct {
	ID, OptionID          int64
	OptionName            string
	Size, PrintName, Number string
	Quantity, UnitPrice   int64
}

type Payment struct {
	ID           int64
	Slip         string // "" = cash
	Status       api.PaymentStatus
	Amount       *int64
	RejectReason string
	CreatedAt    time.Time
	ReviewedBy   string
	ReviewedAt   *time.Time
}

type Adjustment struct {
	ID        int64
	Amount    int64
	Note      string
	CreatedBy string
	CreatedAt time.Time
}

// Status is the ONLY place effective Order status is computed:
// stored awaiting_payment + Campaign effectively closed + not placed by an Admin = cancelled.
// Reopening the Campaign makes such Orders read awaiting_payment again; nothing was stored.
func (o Order) Status() api.OrderStatus {
	panic("not implemented")
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
// "Campaign effectively open" for the owner (Admins are never bound by close, story 69/70).
// `to` == "" leaves status unchanged.
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

// rules IS the Order state machine (ADR 0001). Order matters only for readability.
// Side effects of each transition are written by the handler inside server.apply's
// transaction; the table owns only "may this happen, and to which status".
var rules = []rule{
	{EvView, owner | admin, anyStatus, "", false},
	{EvUploadSlip, owner, []api.OrderStatus{awaiting}, reviewing, true},
	{EvUploadSlip, admin, []api.OrderStatus{awaiting}, reviewing, false},
	{EvConfirmSlip, admin, []api.OrderStatus{reviewing}, paid, false},
	{EvRejectSlip, admin, []api.OrderStatus{reviewing}, awaiting, false},
	// Cash settles an unpaid Order even after close (spec "Further notes"), or tops up a paid one.
	// Not from pending_review: review the Slip first, so a pending Payment never outlives its Order's review.
	{EvRecordCash, admin, []api.OrderStatus{awaiting, paid}, paid, false},
	{EvCancel, owner, []api.OrderStatus{awaiting}, cancelled, true},
	// Admin cancel is idempotent: cancelling a cancelled Order succeeds and changes nothing.
	{EvCancel, admin, anyStatus, cancelled, false},
	{EvEdit, admin, anyStatus, "", false},
	{EvAdjust, admin, anyStatus, "", false},
}

// Decide answers "may actor fire ev on this Order now, and what is the new STORED status".
// Errors: NotFound when the actor may not even view it (another Member's Order),
// Forbidden when no rule names the actor's party for ev, Conflict when a rule names the
// party but the status (or the closed Campaign) does not allow it, with a Thai reason
// such as "แคมเปญปิดแล้ว กรุณาติดต่อแอดมิน".
func (o Order) Decide(a Actor, ev Event) (to api.OrderStatus, err error) {
	panic("not implemented")
}

// Actions evaluates Decide for every event, for the response's api.OrderActions.
// Derived from rules, so buttons in the UI can never disagree with the server.
func (o Order) Actions(a Actor) api.OrderActions {
	panic("not implemented")
}
