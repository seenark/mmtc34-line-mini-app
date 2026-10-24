package server_test

import (
	"fmt"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

// Fixtures go through the real HTTP API (dev login, onboarding, approval) and assert each status.
// They are idempotent: asking for the same person again signs in again instead of onboarding twice.

// Admin is the bootstrap Admin, signed in and onboarded.
func (h *H) Admin() *Client { return h.signIn(adminLine, "admin") }

// Applicant signs in as "U-<nick>" and onboards; not yet approved.
func (h *H) Applicant(nick string) *Client { return h.signIn("U-"+nick, nick) }

// Member is Applicant plus Admin approval.
func (h *H) Member(nick string) *Client {
	c := h.Applicant(nick)
	if m := c.Me().Member; m.Status != api.MemberApproved {
		h.Admin().POST(fmt.Sprintf("/api/admin/members/%d/approve", m.ID), nil).Status(200)
	}
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
