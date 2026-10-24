package server_test

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/server"
)

func TestInvalidIDTokenIs401WithoutCookie(t *testing.T) {
	h := New(t)
	c := h.Anonymous()

	r := c.POST("/api/auth/line", api.LineLoginInput{IDToken: "garbage"}).Status(401)

	Equal(t, r.Header.Get("Set-Cookie"), "")
	c.GET("/api/me").Status(401)
}

func TestLineLoginStartsSessionWithNoMemberYet(t *testing.T) {
	h := New(t)
	c := h.Anonymous()

	r := c.POST("/api/auth/line", api.LineLoginInput{IDToken: "ok:U-nok"}).Status(204)

	cookie := r.Header.Get("Set-Cookie")
	for _, want := range []string{"mmtc34_session=", "HttpOnly", "SameSite=Lax", "Path=" + basePath, "Max-Age=2592000"} {
		if !strings.Contains(cookie, want) {
			t.Fatalf("Set-Cookie %q lacks %q", cookie, want)
		}
	}
	me := c.Me()
	Equal(t, me.LineUserID, "U-nok")
	Equal(t, me.Member == nil, true)
}

func TestLogoutEndsSession(t *testing.T) {
	h := New(t)
	c := h.Applicant("nok")

	c.POST("/api/auth/logout", nil).Status(204)

	c.GET("/api/me").Status(401)
}

func TestOnboardingCreatesApplicant(t *testing.T) {
	h := New(t)
	c := h.Anonymous()
	c.POST("/api/auth/dev", api.DevLoginInput{LineUserID: "U-nok"}).Status(204)

	m := As[api.Member](c.POST("/api/me", api.ProfileInput{Cohort: " 34 ", Nickname: " nok ", RealName: "Nok Jai"}).Status(201))

	Equal(t, m.Status, api.MemberApplicant)
	Equal(t, m.Role, api.RoleMember)
	Equal(t, m.Linked, true)
	Equal(t, m.Nickname, "nok")
	Equal(t, m.Cohort, "34")
	Equal(t, c.Me().Member.ID, m.ID)
}

func TestOnboardingRejectsBlankRequiredFields(t *testing.T) {
	h := New(t)
	c := h.Anonymous()
	c.POST("/api/auth/dev", api.DevLoginInput{LineUserID: "U-nok"}).Status(204)

	for _, in := range []api.ProfileInput{
		{Cohort: "34", Nickname: " ", RealName: "x"},
		{Cohort: "34", Nickname: "nok", RealName: ""},
		{Cohort: "  ", Nickname: "nok", RealName: "x"},
	} {
		if msg := c.POST("/api/me", in).Status(400).Error(); msg == "" {
			t.Fatalf("%+v: empty error message", in)
		}
	}
	Equal(t, c.Me().Member == nil, true)
}

func TestSecondOnboardingIs409(t *testing.T) {
	h := New(t)
	c := h.Applicant("nok")

	c.POST("/api/me", api.ProfileInput{Cohort: "34", Nickname: "again", RealName: "Again"}).Status(409)
}

func TestApplicantCannotUseAdminRoutes(t *testing.T) {
	h := New(t)
	applicant := h.Applicant("nok")

	applicant.GET("/api/admin/members").Status(403)
	h.Anonymous().GET("/api/admin/members").Status(401)
}

func TestAdminApprovalLetsApplicantIn(t *testing.T) {
	h := New(t)
	admin := h.Admin()
	nok := h.Applicant("nok")
	id := nok.Me().Member.ID

	pending := As[[]api.Member](admin.GET("/api/admin/members?status=applicant").Status(200))
	Equal(t, len(pending), 1)
	Equal(t, pending[0].ID, id)

	approved := As[api.Member](admin.POST(fmt.Sprintf("/api/admin/members/%d/approve", id), nil).Status(200))

	Equal(t, approved.Status, api.MemberApproved)
	Equal(t, nok.Me().Member.Status, api.MemberApproved)
	Equal(t, len(As[[]api.Member](admin.GET("/api/admin/members?status=applicant").Status(200))), 0)
	Equal(t, As[api.Member](admin.GET(fmt.Sprintf("/api/admin/members/%d", id)).Status(200)).Status, api.MemberApproved)
	nok.GET("/api/admin/members").Status(403)
}

func TestAdminRejectDeletesApplicant(t *testing.T) {
	h := New(t)
	admin := h.Admin()
	id := h.Applicant("nok").Me().Member.ID

	admin.POST(fmt.Sprintf("/api/admin/members/%d/reject", id), nil).Status(204)

	admin.GET(fmt.Sprintf("/api/admin/members/%d", id)).Status(404)
	again := h.Applicant("nok").Me().Member // a rejected person can apply again
	Equal(t, again.Status, api.MemberApplicant)
}

func TestRejectingApprovedMemberIs409(t *testing.T) {
	h := New(t)
	id := h.Member("nok").Me().Member.ID

	h.Admin().POST(fmt.Sprintf("/api/admin/members/%d/reject", id), nil).Status(409)
}

func TestBootstrapIDOnboardsAsApprovedAdmin(t *testing.T) {
	h := New(t)

	admin := h.Admin()

	m := admin.Me().Member
	Equal(t, m.Status, api.MemberApproved)
	Equal(t, m.Role, api.RoleAdmin)
	admin.GET("/api/admin/members").Status(200)
}

func TestBootstrapLoginPromotesExistingMember(t *testing.T) {
	adminIDs := map[string]bool{} // shared with the server's Config, so the id can be listed after onboarding
	h := New(t, func(c *server.Config) { c.AdminLineUserIDs = adminIDs })
	boss := h.Applicant("boss")
	Equal(t, boss.Me().Member.Role, api.RoleMember)

	adminIDs["U-boss"] = true
	boss.POST("/api/auth/dev", api.DevLoginInput{LineUserID: "U-boss"}).Status(204)

	m := boss.Me().Member
	Equal(t, m.Role, api.RoleAdmin)
	Equal(t, m.Status, api.MemberApproved)
}

func TestDevLoginIs404WhenOff(t *testing.T) {
	h := New(t, func(c *server.Config) { c.DevLogin = false })
	c := h.Anonymous()

	c.POST("/api/auth/dev", api.DevLoginInput{LineUserID: adminLine}).Status(404)

	c.GET("/api/me").Status(401)
}

func TestTamperedCookieIs401(t *testing.T) {
	h := New(t)
	login := h.Anonymous().POST("/api/auth/dev", api.DevLoginInput{LineUserID: "U-nok"}).Status(204)
	parts := strings.Split((&http.Response{Header: login.Header}).Cookies()[0].Value, ".")
	parts[0] = base64.RawURLEncoding.EncodeToString([]byte(adminLine))

	forged := h.Anonymous()
	u, _ := url.Parse(h.srv.URL + basePath)
	forged.http.Jar.SetCookies(u, []*http.Cookie{{Name: "mmtc34_session", Value: strings.Join(parts, "."), Path: basePath}})

	forged.GET("/api/me").Status(401)
}

func TestSessionExpiresAfterThirtyDays(t *testing.T) {
	h := New(t)
	c := h.Applicant("nok")

	h.Advance(29 * 24 * time.Hour)
	c.GET("/api/me").Status(200)
	h.Advance(2 * 24 * time.Hour)

	c.GET("/api/me").Status(401)
}
