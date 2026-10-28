package api

type OrderStatus string

const (
	OrderAwaitingPayment OrderStatus = "awaiting_payment"
	OrderPendingReview   OrderStatus = "pending_review"
	OrderPaid            OrderStatus = "paid"
	OrderCancelled       OrderStatus = "cancelled"
)

type Delivery string

const (
	DeliveryShip   Delivery = "ship"
	DeliveryPickup Delivery = "pickup"
	DeliveryNone   Delivery = "none" // donation-only Orders: required there, invalid elsewhere
)

// AddressChoice is one shipping triple offered at checkout: the profile, or the snapshot of an
// earlier Order. Choosing it only copies it onto the new Order.
type AddressChoice struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type ItemInput struct {
	ID        *int64 `json:"id"` // Admin edit only: existing Item keeps its captured unit price
	OptionID  int64  `json:"optionId"`
	Size      string `json:"size"`
	PrintName string `json:"printName"`
	Number    string `json:"number"`
	Quantity  int64  `json:"quantity"`
	Amount    *int64 `json:"amount"` // only for Options with price null
}

// OrderInput is the checkout body (Member) and the Admin edit body.
type OrderInput struct {
	Delivery    Delivery    `json:"delivery"`
	ShipName    string      `json:"shipName"`
	ShipAddress string      `json:"shipAddress"`
	ShipPhone   string      `json:"shipPhone"`
	Items       []ItemInput `json:"items"`
}

// ItemDetailsInput is the owner-editable part of an existing Item; quantity, Option and price are not accepted.
type ItemDetailsInput struct {
	ID        int64  `json:"id"`
	Size      string `json:"size"`
	PrintName string `json:"printName"`
	Number    string `json:"number"`
}

// OrderDetailsInput is the body of PUT /orders/:orderId/details. Items must list every existing Item once.
type OrderDetailsInput struct {
	ShipName    string             `json:"shipName"`
	ShipAddress string             `json:"shipAddress"`
	ShipPhone   string             `json:"shipPhone"`
	Items       []ItemDetailsInput `json:"items"`
}

// AdminOrderInput is the body of POST /admin/campaigns/:campaignId/orders (on behalf).
type AdminOrderInput struct {
	MemberID int64      `json:"memberId"`
	Order    OrderInput `json:"order"`
}

type Item struct {
	ID         int64  `json:"id"`
	OptionID   int64  `json:"optionId"`
	OptionName string `json:"optionName"`
	Size       string `json:"size"`
	PrintName  string `json:"printName"`
	Number     string `json:"number"`
	Quantity   int64  `json:"quantity"`
	UnitPrice  int64  `json:"unitPrice"`
}

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentConfirmed PaymentStatus = "confirmed"
	PaymentRejected  PaymentStatus = "rejected"
)

type Payment struct {
	ID           int64         `json:"id"`
	OrderID      int64         `json:"orderId"`
	Slips        []string      `json:"slips"` // slip image names, oldest first; empty = cash recorded by an Admin
	Status       PaymentStatus `json:"status"`
	Amount       *int64        `json:"amount"` // set when confirmed
	RejectReason string        `json:"rejectReason"`
	CreatedAt    string        `json:"createdAt"`
	ReviewedBy   string        `json:"reviewedBy"` // admin nickname, "" = not reviewed
	ReviewedAt   *string       `json:"reviewedAt"`
}

type Adjustment struct {
	ID        int64  `json:"id"`
	Amount    int64  `json:"amount"` // + = Member owes more, - = discount/refund
	Note      string `json:"note"`
	CreatedBy string `json:"createdBy"`
	CreatedAt string `json:"createdAt"`
}

// OrderActions is domain.Order.Actions(viewer): which lifecycle events the viewer may fire
// right now. The web app shows buttons from this and never encodes the state machine itself.
type OrderActions struct {
	Cancel      bool `json:"cancel"`
	UploadSlip  bool `json:"uploadSlip"`
	ReviewSlip  bool `json:"reviewSlip"` // confirm or reject
	EditSlips   bool `json:"editSlips"`  // Admin adds or removes images of a Slip Payment
	RecordCash  bool `json:"recordCash"`
	Edit        bool `json:"edit"`
	EditDetails bool `json:"editDetails"` // the owner fixes size, print name, number and shipping; never money
	AddAdjust   bool `json:"addAdjust"`
}

type Order struct {
	ID             int64        `json:"id"`
	CampaignID     int64        `json:"campaignId"`
	CampaignTitle  string       `json:"campaignTitle"`
	MemberID       int64        `json:"memberId"`
	MemberNickname string       `json:"memberNickname"`
	PlacedByAdmin  bool         `json:"placedByAdmin"` // true = an Admin placed it on the Member's behalf
	Status         OrderStatus  `json:"status"`        // effective (domain.Order.Status)
	Delivery       Delivery     `json:"delivery"`
	ShipName       string       `json:"shipName"`
	ShipAddress    string       `json:"shipAddress"`
	ShipPhone      string       `json:"shipPhone"`
	Items          []Item       `json:"items"`
	Payments       []Payment    `json:"payments"`
	Adjustments    []Adjustment `json:"adjustments"`
	Total          int64        `json:"total"`   // Σ quantity × unitPrice
	Paid           int64        `json:"paid"`    // Σ confirmed amounts
	Balance        int64        `json:"balance"` // domain.Order.Money: >0 owes, 0 settled, <0 overpaid
	Actions        OrderActions `json:"actions"`
	CreatedAt      string       `json:"createdAt"`
}

type AmountInput struct {
	Amount int64 `json:"amount"`
}

type RejectInput struct {
	Reason string `json:"reason"`
}

type AdjustmentInput struct {
	Amount int64  `json:"amount"`
	Note   string `json:"note"`
}
