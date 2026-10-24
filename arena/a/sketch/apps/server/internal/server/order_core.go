package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

// Ticket .4. The Order core every later ticket calls and none edits.

// orderScope selects Orders; nil fields mean "any". Status filtering is never done in SQL:
// it must use effective status, which only domain.Order.Status knows.
type orderScope struct {
	OrderID, CampaignID, MemberID *int64
}

// loadOrders is the ONLY Order read path: ListOrders + Items/Payments/Adjustments for the
// returned ids (four queries, no N+1), assembled into domain.Order with Now = s.now().
// shortcut: callers filter in Go over whole-Campaign loads; fine for one class (hundreds of
// Orders); push filters into SQL if a Campaign ever exceeds ~10k Items.
func (s *Server) loadOrders(ctx context.Context, q *db.Queries, scope orderScope) ([]domain.Order, error) {
	panic("not implemented")
}

// toAPIOrder is the ONLY domain.Order -> api.Order mapping: Status(), Money(), Actions(viewer).
// RealName is never copied out.
func toAPIOrder(o domain.Order, a domain.Actor) api.Order {
	panic("not implemented")
}

// effect writes a transition's side rows (Payment, Adjustment, Item replacement) inside the
// same transaction as the status change. It receives the Order as loaded before the change.
type effect func(q *db.Queries, o domain.Order) error

// apply is every lifecycle endpoint in one call. In one IMMEDIATE transaction it:
//  1. loads the Order from path param :orderId (404 if missing);
//  2. o.Decide(actor, ev) -> NotFound/Forbidden/Conflict as domain.Error;
//  3. if the status changes: SetOrderStatus(id, from=o.Stored, to) and 0 rows -> 409
//     "คำสั่งซื้อถูกเปลี่ยนสถานะไปแล้ว" (a raced or replayed request applies at most once);
//  4. runs fx (may be nil);
//  5. reloads the Order and writes 200 api.Order.
// On any error it writes the response via fail() and returns ok=false, so handlers can clean
// up (e.g. delete a just-saved Slip file).
func (s *Server) apply(c *gin.Context, ev domain.Event, fx effect) (ok bool) {
	panic("not implemented")
}

// checkout validates and inserts a new Order (status awaiting_payment, unit prices captured)
// and writes 201 api.Order. placedBy is nil for Member checkout, the Admin's id on behalf.
// Shared by POST /campaigns/:id/orders (.4) and POST /admin/campaigns/:id/orders (.7).
func (s *Server) checkout(c *gin.Context, campaignID, memberID int64, placedBy *int64, in api.OrderInput) {
	panic("not implemented")
}
