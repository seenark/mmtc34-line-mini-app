// Package api holds the JSON request/response shapes shared with the web app.
// web/src/api/types.ts is generated from this file by `go tool tygo generate`.
package api

// Money is whole Thai baht everywhere.

type Config struct {
	LiffID   string `json:"liffId"`
	DevLogin bool   `json:"devLogin"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// --- people ---

type MemberStatus string

const (
	MemberApplicant MemberStatus = "applicant"
	MemberApproved  MemberStatus = "approved"
)

type Role string

const (
	RoleMember Role = "member"
	RoleAdmin  Role = "admin"
)

type Member struct {
	ID       int64        `json:"id"`
	Linked   bool         `json:"linked"` // has a LINE account; false = Unlinked Member
	Status   MemberStatus `json:"status"`
	Role     Role         `json:"role"`
	Cohort   string       `json:"cohort"`
	Nickname string       `json:"nickname"`
	RealName string       `json:"realName"`
	Address  string       `json:"address"`
	Phone    string       `json:"phone"`
	Deleted  bool         `json:"deleted"`
}

type Me struct {
	LineUserID string  `json:"lineUserId"`
	Member     *Member `json:"member"` // null until onboarding is submitted
}

type ProfileInput struct {
	Cohort   string `json:"cohort"`
	Nickname string `json:"nickname"`
	RealName string `json:"realName"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
}

type LineLoginInput struct {
	IDToken string `json:"idToken"`
}

type DevLoginInput struct {
	LineUserID string `json:"lineUserId"`
}

type RoleInput struct {
	Role Role `json:"role"`
}

type LinkInput struct {
	ApplicantID int64 `json:"applicantId"`
}

// --- campaigns ---

type CampaignType string

const (
	CampaignMerch    CampaignType = "merch"
	CampaignDonation CampaignType = "donation"
	CampaignEvent    CampaignType = "event"
)

type CampaignStatus string

const (
	CampaignDraft  CampaignStatus = "draft"
	CampaignOpen   CampaignStatus = "open"
	CampaignClosed CampaignStatus = "closed"
)

type Campaign struct {
	ID              int64          `json:"id"`
	Type            CampaignType   `json:"type"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	CoverImage      string         `json:"coverImage"` // image name, "" = none
	BankName        string         `json:"bankName"`
	BankAccountNo   string         `json:"bankAccountNo"`
	BankAccountName string         `json:"bankAccountName"`
	QRImage         string         `json:"qrImage"` // image name, "" = none
	Status          CampaignStatus `json:"status"`  // effective: open past closeAt reads as closed
	CloseAt         *string        `json:"closeAt"` // RFC3339
}

type CampaignInput struct {
	Type            CampaignType   `json:"type"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	CoverImage      string         `json:"coverImage"`
	BankName        string         `json:"bankName"`
	BankAccountNo   string         `json:"bankAccountNo"`
	BankAccountName string         `json:"bankAccountName"`
	QRImage         string         `json:"qrImage"`
	Status          CampaignStatus `json:"status"`
	CloseAt         *string        `json:"closeAt"`
}

type Option struct {
	ID             int64    `json:"id"`
	CampaignID     int64    `json:"campaignId"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Image          string   `json:"image"`
	Price          *int64   `json:"price"` // null = Member chooses the amount (donation)
	Sizes          []string `json:"sizes"` // empty = no size
	NeedsPrintName bool     `json:"needsPrintName"`
	NeedsNumber    bool     `json:"needsNumber"`
	SortOrder      int64    `json:"sortOrder"`
}

type OptionInput struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Image          string   `json:"image"`
	Price          *int64   `json:"price"`
	Sizes          []string `json:"sizes"`
	NeedsPrintName bool     `json:"needsPrintName"`
	NeedsNumber    bool     `json:"needsNumber"`
	SortOrder      int64    `json:"sortOrder"`
}

type SponsorPosition struct {
	ID             int64  `json:"id"`
	CampaignID     int64  `json:"campaignId"`
	Name           string `json:"name"`
	SponsorName    string `json:"sponsorName"`
	Logo           string `json:"logo"`
	MemberID       *int64 `json:"memberId"`
	MemberNickname string `json:"memberNickname"`
	SortOrder      int64  `json:"sortOrder"`
}

type SponsorInput struct {
	Name        string `json:"name"`
	SponsorName string `json:"sponsorName"`
	Logo        string `json:"logo"`
	MemberID    *int64 `json:"memberId"`
	SortOrder   int64  `json:"sortOrder"`
}

type CampaignDetail struct {
	Campaign Campaign          `json:"campaign"`
	Options  []Option          `json:"options"`
	Sponsors []SponsorPosition `json:"sponsors"`
}

type RosterItem struct {
	OptionName string `json:"optionName"`
	Size       string `json:"size"`
	Quantity   int64  `json:"quantity"`
}

type RosterEntry struct {
	Nickname string       `json:"nickname"`
	Items    []RosterItem `json:"items"`
}

type ImageUpload struct {
	Name string `json:"name"`
}

// --- orders & money ---

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
)

type ItemInput struct {
	OptionID  int64  `json:"optionId"`
	Size      string `json:"size"`
	PrintName string `json:"printName"`
	Number    string `json:"number"`
	Quantity  int64  `json:"quantity"`
	Amount    *int64 `json:"amount"` // only for options with price null
}

type OrderInput struct {
	MemberID    *int64      `json:"memberId"` // admin ordering on someone's behalf; ignored for members
	Delivery    Delivery    `json:"delivery"`
	ShipName    string      `json:"shipName"`
	ShipAddress string      `json:"shipAddress"`
	ShipPhone   string      `json:"shipPhone"`
	Items       []ItemInput `json:"items"`
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
	Slip         string        `json:"slip"` // slip name, "" = cash recorded by an Admin
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

type Order struct {
	ID             int64        `json:"id"`
	CampaignID     int64        `json:"campaignId"`
	CampaignTitle  string       `json:"campaignTitle"`
	MemberID       int64        `json:"memberId"`
	MemberNickname string       `json:"memberNickname"`
	Status         OrderStatus  `json:"status"` // effective: awaiting_payment in a closed campaign reads as cancelled
	Delivery       Delivery     `json:"delivery"`
	ShipName       string       `json:"shipName"`
	ShipAddress    string       `json:"shipAddress"`
	ShipPhone      string       `json:"shipPhone"`
	Items          []Item       `json:"items"`
	Payments       []Payment    `json:"payments"`
	Adjustments    []Adjustment `json:"adjustments"`
	Total          int64        `json:"total"`   // sum of quantity * unitPrice
	Paid           int64        `json:"paid"`    // sum of confirmed payment amounts
	Balance        int64        `json:"balance"` // total - paid + adjustments
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
