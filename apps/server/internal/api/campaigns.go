package api

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

// CampaignDetail: Sponsor Positions are served separately (GET /campaigns/:id/sponsors).
type CampaignDetail struct {
	Campaign Campaign `json:"campaign"`
	Options  []Option `json:"options"`
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

type RosterItem struct {
	OptionName string `json:"optionName"`
	Size       string `json:"size"`
	Quantity   int64  `json:"quantity"`
}

// RosterEntry deliberately has no status, money, address, phone or real name.
type RosterEntry struct {
	Nickname string       `json:"nickname"`
	Items    []RosterItem `json:"items"`
}
