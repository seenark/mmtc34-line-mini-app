package server_test

import (
	"encoding/base64"
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
	c := h.Member("nok")

	c.POST("/api/auth/logout", nil).Status(204)

	c.GET("/api/me").Status(401)
}

func TestOnboardingCreatesAMemberWhoCanBrowseCampaignsAtOnce(t *testing.T) {
	h := New(t)
	j := h.OpenJersey()
	c := h.Anonymous()
	c.POST("/api/auth/line", api.LineLoginInput{IDToken: "ok:U-nok"}).Status(204)
	c.GET("/api/campaigns").Status(403)

	m := As[api.Member](c.POST("/api/me", api.ProfileInput{Cohort: " 34 ", Nickname: " nok ", RealName: "Nok Jai"}).Status(201))

	Equal(t, m.Status, api.MemberApproved)
	Equal(t, m.Role, api.RoleMember)
	Equal(t, m.Linked, true)
	Equal(t, m.Nickname, "nok")
	Equal(t, m.Cohort, "34")
	Equal(t, c.Me().Member.ID, m.ID)
	campaigns := As[[]api.Campaign](c.GET("/api/campaigns").Status(200))
	Equal(t, len(campaigns), 1)
	Equal(t, campaigns[0].ID, j.Campaign.ID)
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
	c := h.Member("nok")

	c.POST("/api/me", api.ProfileInput{Cohort: "34", Nickname: "again", RealName: "Again"}).Status(409)
}

func TestMemberCannotUseAdminRoutes(t *testing.T) {
	h := New(t)
	nok := h.Member("nok")

	nok.GET("/api/admin/members").Status(403)
	h.Anonymous().GET("/api/admin/members").Status(401)
}

func TestProductionLoginSetsSecureCookie(t *testing.T) {
	h := New(t, func(c *server.Config) { c.DevLogin = false })
	c := h.Anonymous()

	r := c.POST("/api/auth/line", api.LineLoginInput{IDToken: "ok:U-production"}).Status(204)
	cookies := (&http.Response{Header: r.Header}).Cookies()
	var session *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "mmtc34_session" {
			session = cookie
			break
		}
	}
	if session == nil {
		t.Fatal("expected mmtc34_session cookie")
	}
	if !session.Secure {
		t.Error("expected session cookie to be Secure")
	}
	if !session.HttpOnly {
		t.Error("expected session cookie to be HttpOnly")
	}
	if session.SameSite != http.SameSiteLaxMode {
		t.Errorf("expected SameSite=Lax, got %v", session.SameSite)
	}

	c.POST("/api/auth/dev", api.DevLoginInput{LineUserID: "U-production"}).Status(404)
}

func TestDevLoginIs404WhenOff(t *testing.T) {
	h := New(t, func(c *server.Config) { c.DevLogin = false })
	c := h.Anonymous()

	c.POST("/api/auth/dev", api.DevLoginInput{LineUserID: "U-admin"}).Status(404)

	c.GET("/api/me").Status(401)
}

func TestTamperedCookieIs401(t *testing.T) {
	h := New(t)
	login := h.Anonymous().POST("/api/auth/dev", api.DevLoginInput{LineUserID: "U-nok"}).Status(204)
	parts := strings.Split((&http.Response{Header: login.Header}).Cookies()[0].Value, ".")
	parts[0] = base64.RawURLEncoding.EncodeToString([]byte("superadmin"))

	forged := h.Anonymous()
	u, _ := url.Parse(h.srv.URL + basePath)
	forged.http.Jar.SetCookies(u, []*http.Cookie{{Name: "mmtc34_session", Value: strings.Join(parts, "."), Path: basePath}})

	forged.GET("/api/me").Status(401)
}

func TestSessionExpiresAfterThirtyDays(t *testing.T) {
	h := New(t)
	c := h.Member("nok")

	h.Advance(29 * 24 * time.Hour)
	c.GET("/api/me").Status(200)
	h.Advance(2 * 24 * time.Hour)

	c.GET("/api/me").Status(401)
}
