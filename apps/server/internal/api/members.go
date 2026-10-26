package api

// MemberStatus is stored for schema stability; open membership (ADR 0006) means every Member is approved.
type MemberStatus string

const MemberApproved MemberStatus = "approved"

type Role string

const (
	RoleMember Role = "member"
	RoleAdmin  Role = "admin"
)

type Member struct {
	ID         int64        `json:"id"`
	Linked     bool         `json:"linked"` // has a LINE account; false = Unlinked Member
	Status     MemberStatus `json:"status"`
	Role       Role         `json:"role"`
	Cohort     string       `json:"cohort"`
	Nickname   string       `json:"nickname"`
	RealName   string       `json:"realName"`
	Address    string       `json:"address"`
	Phone      string       `json:"phone"`
	Deleted    bool         `json:"deleted"`
	SuperAdmin bool         `json:"superAdmin"` // the env-configured account: role is fixed, not a classmate
}

type Me struct {
	LineUserID string  `json:"lineUserId"`
	Member     *Member `json:"member"` // null until onboarding is submitted (or after account deletion)
}

// ProfileInput is used for onboarding, self edit, Admin edit and creating an Unlinked Member
// (nickname only). Which fields are required depends on the route; domain.ValidateProfile decides.
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

type AdminLoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RoleInput struct {
	Role Role `json:"role"`
}

// UnlinkedRef names the Unlinked Member who owns an Order.
type UnlinkedRef struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
}

// OrderLinkInfo is GET /admin/orders/:orderId/link: what linking this Order would affect.
type OrderLinkInfo struct {
	UnlinkedMember UnlinkedRef `json:"unlinkedMember"`
	OtherOrders    int64       `json:"otherOrders"` // the owner's Orders besides this one
}

// OrderLinkInput is POST /admin/orders/:orderId/link. MemberID is the Member who has the LINE account.
type OrderLinkInput struct {
	MemberID  int64 `json:"memberId"`
	AllOrders bool  `json:"allOrders"` // also move the owner's other Orders
}

type OrderLinkResult struct {
	Moved           int64 `json:"moved"`
	UnlinkedRemoved bool  `json:"unlinkedRemoved"` // the Unlinked Member owned nothing left and was deleted
}
