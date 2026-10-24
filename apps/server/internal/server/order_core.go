package server

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

// orderScope selects Orders; nil fields mean "any". Status filtering is never done in SQL:
// it must use effective status, which only domain.Order.Status knows.
type orderScope struct {
	OrderID, CampaignID, MemberID *int64
}

// loadOrders is the ONLY Order read path: ListOrders (newest first) plus Items, Payments and
// Adjustments for the returned ids in three IN queries, assembled into domain.Order with
// Now = s.now() read once.
// shortcut: callers filter in Go over whole-Campaign loads; fine for one class (hundreds of
// Orders); push filters into SQL if a Campaign ever exceeds ~10k Items.
func (s *Server) loadOrders(ctx context.Context, q *db.Queries, scope orderScope) ([]domain.Order, error) {
	rows, err := q.ListOrders(ctx, db.ListOrdersParams{
		OrderID: scope.OrderID, CampaignID: scope.CampaignID, MemberID: scope.MemberID,
	})
	if err != nil || len(rows) == 0 {
		return nil, err
	}
	now := s.now()
	orders := make([]domain.Order, len(rows))
	ids := make([]int64, len(rows))
	at := make(map[int64]*domain.Order, len(rows))
	for i, r := range rows {
		o := r.Order
		orders[i] = domain.Order{
			ID:            o.ID,
			CampaignID:    o.CampaignID,
			CampaignTitle: r.Campaign.Title,
			Campaign:      domainCampaign(r.Campaign),
			MemberID:      o.MemberID,
			Nickname:      r.MemberNickname,
			RealName:      r.MemberRealName,
			PlacedByAdmin: o.PlacedByAdmin == 1,
			Stored:        api.OrderStatus(o.Status),
			Delivery:      api.Delivery(o.Delivery),
			ShipName:      o.ShipName,
			ShipAddress:   o.ShipAddress,
			ShipPhone:     o.ShipPhone,
			CreatedAt:     o.CreatedAt,
			Now:           now,
		}
		ids[i] = o.ID
		at[o.ID] = &orders[i]
	}

	items, err := q.ListItemsForOrders(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, it := range items {
		o := at[it.OrderID]
		o.Items = append(o.Items, domain.Item{
			ID: it.ID, OptionID: it.OptionID, OptionName: it.OptionName, Size: it.Size,
			PrintName: it.PrintName, Number: it.Number, Quantity: it.Quantity, UnitPrice: it.UnitPrice,
		})
	}
	payments, err := q.ListPaymentsForOrders(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, p := range payments {
		o := at[p.OrderID]
		o.Payments = append(o.Payments, domain.Payment{
			ID:           p.ID,
			Slip:         emptyIfNil(p.Slip),
			Status:       api.PaymentStatus(p.Status),
			Amount:       p.Amount,
			RejectReason: p.RejectReason,
			CreatedAt:    p.CreatedAt,
			ReviewedBy:   emptyIfNil(p.ReviewerNickname),
			ReviewedAt:   p.ReviewedAt,
		})
	}
	adjustments, err := q.ListAdjustmentsForOrders(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, a := range adjustments {
		o := at[a.OrderID]
		o.Adjustments = append(o.Adjustments, domain.Adjustment{
			ID: a.ID, Amount: a.Amount, Note: a.Note, CreatedBy: a.CreatorNickname, CreatedAt: a.CreatedAt,
		})
	}
	return orders, nil
}

func emptyIfNil(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// loadOrder is loadOrders for one id; a missing Order is sql.ErrNoRows (404 via fail).
func (s *Server) loadOrder(ctx context.Context, q *db.Queries, id int64) (domain.Order, error) {
	os, err := s.loadOrders(ctx, q, orderScope{OrderID: &id})
	if err != nil {
		return domain.Order{}, err
	}
	if len(os) == 0 {
		return domain.Order{}, sql.ErrNoRows
	}
	return os[0], nil
}

// toAPIOrder is the ONLY domain.Order -> api.Order mapping: Status(), Money(), Actions(viewer).
// RealName is never copied out.
func toAPIOrder(o domain.Order, a domain.Actor) api.Order {
	m := o.Money()
	out := api.Order{
		ID:             o.ID,
		CampaignID:     o.CampaignID,
		CampaignTitle:  o.CampaignTitle,
		MemberID:       o.MemberID,
		MemberNickname: o.Nickname,
		PlacedByAdmin:  o.PlacedByAdmin,
		Status:         o.Status(),
		Delivery:       o.Delivery,
		ShipName:       o.ShipName,
		ShipAddress:    o.ShipAddress,
		ShipPhone:      o.ShipPhone,
		Items:          make([]api.Item, len(o.Items)),
		Payments:       make([]api.Payment, len(o.Payments)),
		Adjustments:    make([]api.Adjustment, len(o.Adjustments)),
		Total:          m.Total,
		Paid:           m.Paid,
		Balance:        m.Balance,
		Actions:        o.Actions(a),
		CreatedAt:      o.CreatedAt,
	}
	for i, it := range o.Items {
		out.Items[i] = api.Item{
			ID: it.ID, OptionID: it.OptionID, OptionName: it.OptionName, Size: it.Size,
			PrintName: it.PrintName, Number: it.Number, Quantity: it.Quantity, UnitPrice: it.UnitPrice,
		}
	}
	for i, p := range o.Payments {
		out.Payments[i] = api.Payment{
			ID: p.ID, OrderID: o.ID, Slip: p.Slip, Status: p.Status, Amount: p.Amount,
			RejectReason: p.RejectReason, CreatedAt: p.CreatedAt, ReviewedBy: p.ReviewedBy, ReviewedAt: p.ReviewedAt,
		}
	}
	for i, ad := range o.Adjustments {
		out.Adjustments[i] = api.Adjustment{
			ID: ad.ID, Amount: ad.Amount, Note: ad.Note, CreatedBy: ad.CreatedBy, CreatedAt: ad.CreatedAt,
		}
	}
	return out
}

// effect writes a transition's side rows (Payment, Adjustment, Item replacement) inside the
// same transaction as the status change. It receives the Order as loaded before the change.
type effect func(q *db.Queries, o domain.Order) error

// apply is every lifecycle endpoint in one call. For the Order named by path param :orderId it:
//  1. loads it (404 if missing);
//  2. o.Decide(actor, ev): NotFound/Forbidden/Conflict as domain.Error;
//  3. if the stored status changes, compare-and-sets it; 0 rows -> 409 (a raced or replayed
//     request applies at most once);
//  4. runs fx (may be nil);
//  5. reloads the Order and writes 200 api.Order.
//
// Steps 1-5 run in one IMMEDIATE transaction, except EvView, which never writes and so runs
// without one. On any error it writes the response via fail() and returns false, so handlers
// can clean up (e.g. delete a just-saved Slip file).
func (s *Server) apply(c *gin.Context, ev domain.Event, fx effect) (ok bool) {
	id, ok := idParam(c, "orderId")
	if !ok {
		return false
	}
	ctx := c.Request.Context()
	a := actorOf(c)
	var out api.Order
	run := func(q *db.Queries) error {
		o, err := s.loadOrder(ctx, q, id)
		if err != nil {
			return err
		}
		next, err := o.Decide(a, ev)
		if err != nil {
			return err
		}
		if next == o.Stored && fx == nil {
			out = toAPIOrder(o, a)
			return nil
		}
		if next != o.Stored {
			n, err := q.SetOrderStatus(ctx, db.SetOrderStatusParams{
				OrderID: id, FromStatus: string(o.Stored), ToStatus: string(next), Now: rfc3339(s.now()),
			})
			if err != nil {
				return err
			}
			if n == 0 {
				return &domain.Error{Kind: domain.Conflict, Msg: "คำสั่งซื้อถูกเปลี่ยนสถานะไปแล้ว"}
			}
		}
		if fx != nil {
			if err := fx(q, o); err != nil {
				return err
			}
		}
		if o, err = s.loadOrder(ctx, q, id); err != nil {
			return err
		}
		out = toAPIOrder(o, a)
		return nil
	}
	var err error
	if ev == domain.EvView {
		err = run(s.q)
	} else {
		err = s.tx(ctx, run)
	}
	if err != nil {
		fail(c, err)
		return false
	}
	c.JSON(http.StatusOK, out)
	return true
}

// checkout validates and inserts a new Order (status awaiting_payment, unit prices captured)
// and writes 201 api.Order. placedByAdmin is false for Member checkout (the Campaign must be
// effectively open) and true when an Admin orders on a Member's behalf (any non-draft Campaign).
// Shared by POST /campaigns/:id/orders and the Admin on-behalf route.
func (s *Server) checkout(c *gin.Context, campaignID, memberID int64, placedByAdmin bool, in api.OrderInput) {
	ctx := c.Request.Context()
	var order domain.Order
	err := s.tx(ctx, func(q *db.Queries) error {
		camp, err := s.loadCampaign(c, q, campaignID)
		if err != nil {
			return err
		}
		if err := domainCampaign(camp).CanCheckout(placedByAdmin, s.now()); err != nil {
			return err
		}
		options, err := s.campaignOptions(ctx, q, campaignID)
		if err != nil {
			return err
		}
		v, err := domain.ValidateOrder(options, in)
		if err != nil {
			return err
		}
		name, address, phone := v.Ship()
		id, err := q.InsertOrder(ctx, db.InsertOrderParams{
			CampaignID:    campaignID,
			MemberID:      memberID,
			PlacedByAdmin: flag(placedByAdmin),
			Delivery:      string(v.Delivery()),
			ShipName:      name,
			ShipAddress:   address,
			ShipPhone:     phone,
			Now:           rfc3339(s.now()),
		})
		if err != nil {
			return err
		}
		for _, it := range v.Items() {
			err := q.InsertItem(ctx, db.InsertItemParams{
				OrderID: id, OptionID: it.OptionID, Size: it.Size, PrintName: it.PrintName,
				Number: it.Number, Quantity: it.Quantity, UnitPrice: it.UnitPrice,
			})
			if err != nil {
				return err
			}
		}
		order, err = s.loadOrder(ctx, q, id)
		return err
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, toAPIOrder(order, actorOf(c)))
}
