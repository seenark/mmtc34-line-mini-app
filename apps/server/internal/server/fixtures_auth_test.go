package server_test

import (
	"fmt"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

// Fixtures go through the real HTTP API (dev login, onboarding) and assert each status.
// They are idempotent: asking for the same person again signs in again instead of onboarding twice.

const superAdminPassword = "test-super-admin-password"

// SuperAdmin signs in with the env username and password.
func (h *H) SuperAdmin() *Client {
	c := h.Anonymous()
	c.POST("/api/auth/admin", api.AdminLoginInput{Username: "admin", Password: superAdminPassword}).Status(204)
	return c
}

// Admin is the LINE user "U-admin", onboarded and promoted by the Super Admin.
func (h *H) Admin() *Client {
	c := h.signIn("U-admin", "admin")
	if m := c.Me().Member; m.Role != api.RoleAdmin {
		h.SuperAdmin().PUT(fmt.Sprintf("/api/admin/members/%d/role", m.ID), api.RoleInput{Role: api.RoleAdmin}).Status(200)
	}
	return c
}

// Member signs in as "U-<nick>" and onboards; open membership makes them a Member at once.
func (h *H) Member(nick string) *Client { return h.signIn("U-"+nick, nick) }

// Visitor signs in as "U-<nick>" without filling in a profile, so they have no Member yet.
func (h *H) Visitor(nick string) *Client {
	c := h.Anonymous()
	c.POST("/api/auth/dev", api.DevLoginInput{LineUserID: "U-" + nick}).Status(204)
	return c
}

// Me is GET /api/me.
func (c *Client) Me() api.Me { return As[api.Me](c.GET("/api/me").Status(200)) }

func (h *H) signIn(lineUserID, nick string) *Client {
	c := h.Anonymous()
	c.POST("/api/auth/dev", api.DevLoginInput{LineUserID: lineUserID}).Status(204)
	if c.Me().Member == nil {
		c.POST("/api/me", api.ProfileInput{Cohort: "34", Nickname: nick, RealName: "Real " + nick}).Status(201)
	}
	return c
}
