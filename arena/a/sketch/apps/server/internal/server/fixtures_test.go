package server_test

// Sketch shows three tickets' fixture files in one file for brevity. In the repo each block is
// its own file owned by its ticket: fixtures_auth_test.go (.2), fixtures_campaign_test.go (.3),
// fixtures_order_test.go (.4), fixtures_payment_test.go (.5). Every fixture goes through the real
// HTTP API and asserts the expected status, so a fixture is also a smoke test of its endpoint.

import "github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"

// ---- fixtures_auth_test.go (.2)

// Admin signs in as the bootstrap Admin and onboards. Same Client on repeated calls.
func (h *H) Admin() *Client { panic("not implemented") }

// Applicant signs in as "U-<nick>" and onboards; not yet approved.
func (h *H) Applicant(nick string) *Client { panic("not implemented") }

// Member is Applicant + Admin approval.
func (h *H) Member(nick string) *Client { panic("not implemented") }

// Me is GET /api/me.
func (c *Client) Me() api.Me { panic("not implemented") }

// ---- fixtures_campaign_test.go (.3)

// Jersey is an open merch Campaign with one 260-baht Option, sizes S..XL, print name + number.
type Jersey struct {
	api.CampaignDetail
	Option api.Option
}

func (h *H) OpenJersey() Jersey { panic("not implemented") }

// Close sets the stored status to closed (Admin action, not the clock).
func (h *H) Close(campaignID int64) { panic("not implemented") }

// ---- fixtures_order_test.go (.4)

// Shirt is one jersey ItemInput.
func Shirt(j Jersey, size string, qty int64) api.ItemInput { panic("not implemented") }

func Pickup(items ...api.ItemInput) api.OrderInput { panic("not implemented") }

func Ship(items ...api.ItemInput) api.OrderInput { panic("not implemented") }

// Checkout POSTs the Order and requires 201.
func (c *Client) Checkout(campaignID int64, in api.OrderInput) api.Order { panic("not implemented") }

// Order GETs one Order and requires 200.
func (c *Client) Order(id int64) api.Order { panic("not implemented") }

// ---- fixtures_payment_test.go (.5)

// UploadSlip POSTs JPEG() and requires 200.
func (c *Client) UploadSlip(orderID int64) api.Order { panic("not implemented") }
