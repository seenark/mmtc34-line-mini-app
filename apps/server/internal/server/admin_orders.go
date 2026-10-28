package server

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

func (s *Server) mountAdminOrders(r routes) {
	r.admin.POST("/campaigns/:campaignId/orders", s.createOrderOnBehalf)
	r.admin.PUT("/orders/:orderId", s.editOrder)
}

// createOrderOnBehalf orders for any non-deleted Member (Unlinked included) in any
// non-draft Campaign, closed ones too. Such Orders are placed_by_admin, so a closed Campaign
// never reads them as cancelled.
func (s *Server) createOrderOnBehalf(c *gin.Context) {
	campaignID, ok := idParam(c, "campaignId")
	if !ok {
		return
	}
	in, ok := bind[api.AdminOrderInput](c)
	if !ok {
		return
	}
	ctx := c.Request.Context()
	camp, err := s.q.GetCampaign(ctx, campaignID)
	if err == nil && domainCampaign(camp).Stored == api.CampaignDraft {
		err = &domain.Error{Kind: domain.NotFound, Msg: "ไม่พบข้อมูล"}
	}
	if err != nil {
		fail(c, err)
		return
	}
	m, err := s.q.GetMember(ctx, in.MemberID)
	if err != nil || m.DeletedAt != nil {
		fail(c, &domain.Error{Kind: domain.Invalid, Msg: "ไม่พบสมาชิก"})
		return
	}
	s.checkout(c, campaignID, m.ID, true, in.Order)
}

// editOrder replaces Items, delivery and shipping of an Order in any status. Items sent with the
// id of an existing Item keep its captured price; the rest take the current Option price.
func (s *Server) editOrder(c *gin.Context) {
	in, ok := bind[api.OrderInput](c)
	if !ok {
		return
	}
	s.apply(c, domain.EvEdit, func(q *db.Queries, o domain.Order) error {
		ctx := c.Request.Context()
		options, err := s.campaignOptions(ctx, q, o.CampaignID)
		if err != nil {
			return err
		}
		kept := make(map[int64]domain.Item, len(o.Items))
		for _, it := range o.Items {
			kept[it.ID] = it
		}
		v, err := domain.ValidateEdit(options, in, kept)
		if err != nil {
			return err
		}
		return s.replaceContents(ctx, q, o.ID, v)
	})
}

// replaceContents rewrites an Order's delivery, shipping and Items from a validated body.
func (s *Server) replaceContents(ctx context.Context, q *db.Queries, orderID int64, v domain.ValidOrder) error {
	name, address, phone := v.Ship()
	err := q.UpdateOrderContents(ctx, db.UpdateOrderContentsParams{
		OrderID: orderID, Delivery: string(v.Delivery()), ShipName: name, ShipAddress: address,
		ShipPhone: phone, Now: rfc3339(s.now()),
	})
	if err != nil {
		return err
	}
	if err := q.DeleteItemsForOrder(ctx, orderID); err != nil {
		return err
	}
	for _, it := range v.Items() {
		err := q.InsertItem(ctx, db.InsertItemParams{
			OrderID: orderID, OptionID: it.OptionID, Size: it.Size, PrintName: it.PrintName,
			Number: it.Number, Quantity: it.Quantity, UnitPrice: it.UnitPrice,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) campaignOptions(ctx context.Context, q *db.Queries, campaignID int64) ([]api.Option, error) {
	rows, err := q.ListOptions(ctx, campaignID)
	if err != nil {
		return nil, err
	}
	out := make([]api.Option, len(rows))
	for i, r := range rows {
		out[i] = toAPIOption(r)
	}
	return out, nil
}
