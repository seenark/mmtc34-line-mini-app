package domain

import (
	"strings"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

// ProfileRequire says which ProfileInput fields a route demands.
type ProfileRequire uint8

const (
	RequireNicknameOnly ProfileRequire = iota // Admin creates an Unlinked Member
	RequireIdentity                           // onboarding/self/admin edit: cohort, nickname, real name
)

// ValidateProfile trims every field and enforces req. Address and phone stay optional.
func ValidateProfile(in api.ProfileInput, req ProfileRequire) (api.ProfileInput, error) {
	in = api.ProfileInput{
		Cohort:   strings.TrimSpace(in.Cohort),
		Nickname: strings.TrimSpace(in.Nickname),
		RealName: strings.TrimSpace(in.RealName),
		Address:  strings.TrimSpace(in.Address),
		Phone:    strings.TrimSpace(in.Phone),
	}
	switch {
	case in.Nickname == "":
		return in, &Error{Kind: Invalid, Msg: "กรุณากรอกชื่อเล่น"}
	case req == RequireIdentity && in.RealName == "":
		return in, &Error{Kind: Invalid, Msg: "กรุณากรอกชื่อ-นามสกุลจริง"}
	case req == RequireIdentity && in.Cohort == "":
		return in, &Error{Kind: Invalid, Msg: "กรุณากรอกรุ่น"}
	}
	return in, nil
}

// IsBootstrapAdmin: LINE userIds in ADMIN_LINE_USER_IDS become approved Admins on login/onboarding.
func IsBootstrapAdmin(lineUserID string, adminIDs map[string]bool) bool {
	return adminIDs[lineUserID]
}
