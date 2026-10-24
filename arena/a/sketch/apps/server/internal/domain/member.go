package domain

import "github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"

// ProfileRequire says which ProfileInput fields a route demands.
type ProfileRequire uint8

const (
	RequireNicknameOnly ProfileRequire = iota // Admin creates an Unlinked Member
	RequireIdentity                           // onboarding/self/admin edit: cohort, nickname, real name
)

// ValidateProfile trims every field and enforces req. Address and phone stay optional.
func ValidateProfile(in api.ProfileInput, req ProfileRequire) (api.ProfileInput, error) {
	panic("not implemented")
}

// IsBootstrapAdmin: LINE userIds in ADMIN_LINE_USER_IDS become approved Admins on login/onboarding.
func IsBootstrapAdmin(lineUserID string, adminIDs map[string]bool) bool {
	return adminIDs[lineUserID]
}
