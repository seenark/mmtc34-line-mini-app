package api

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

type LinkInput struct {
	ApplicantID int64 `json:"applicantId"`
}
