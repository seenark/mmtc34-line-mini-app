package server

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

func (s *Server) mountUnlinked(r routes) {
	r.admin.POST("/members", s.createUnlinked)
	r.admin.GET("/orders/:orderId/link", s.orderLinkInfo)
	r.admin.POST("/orders/:orderId/link", s.linkOrder)
}

// createUnlinked records a classmate who has no LINE account yet: approved, no LINE userId.
func (s *Server) createUnlinked(c *gin.Context) {
	in, ok := bind[api.ProfileInput](c)
	if !ok {
		return
	}
	in, err := domain.ValidateProfile(in, domain.RequireNicknameOnly)
	if err != nil {
		fail(c, err)
		return
	}
	m, err := s.q.CreateMember(c.Request.Context(), db.CreateMemberParams{
		Status:   string(api.MemberApproved),
		Role:     string(api.RoleMember),
		Cohort:   in.Cohort,
		Nickname: in.Nickname,
		RealName: in.RealName,
		Address:  in.Address,
		Phone:    in.Phone,
		Now:      rfc3339(s.now()),
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, memberDTO(m))
}

// checkUnlinkedOwner allows linking only while the Order's owner has no LINE account and is not
// erased. The Super Admin always has a LINE userId, so it can never pass.
func checkUnlinkedOwner(m db.Member) error {
	if m.LineUserID != nil || m.DeletedAt != nil {
		return &domain.Error{Kind: domain.Conflict, Msg: "เจ้าของออเดอร์นี้มีบัญชี LINE แล้วหรือถูกลบไปแล้ว"}
	}
	return nil
}

// orderLinkInfo tells the Admin whom an Order would leave and how many sibling Orders could follow.
func (s *Server) orderLinkInfo(c *gin.Context) {
	id, ok := idParam(c, "orderId")
	if !ok {
		return
	}
	ctx := c.Request.Context()
	owner, err := s.q.GetOrderOwner(ctx, id)
	if err == nil {
		err = checkUnlinkedOwner(owner)
	}
	if err != nil {
		fail(c, err)
		return
	}
	n, err := s.q.CountMemberOrders(ctx, owner.ID)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, api.OrderLinkInfo{
		UnlinkedMember: api.UnlinkedRef{ID: owner.ID, Nickname: owner.Nickname},
		OtherOrders:    n - 1,
	})
}

// linkOrder moves an Unlinked Member's Order (or all their Orders) to the Member who now has a
// LINE account (ADR 0006). Only orders.member_id changes, so Payments, Slips, Adjustments, the
// placed-by-Admin flag and ship snapshots travel with the Order. The emptied Unlinked Member is
// removed in the same transaction; its Sponsor Positions move first so ON DELETE SET NULL
// does not drop them.
func (s *Server) linkOrder(c *gin.Context) {
	id, ok := idParam(c, "orderId")
	if !ok {
		return
	}
	in, ok := bind[api.OrderLinkInput](c)
	if !ok {
		return
	}
	ctx := c.Request.Context()
	var out api.OrderLinkResult
	err := s.tx(ctx, func(q *db.Queries) error {
		owner, err := q.GetOrderOwner(ctx, id)
		if err != nil {
			return err
		}
		if err := checkUnlinkedOwner(owner); err != nil {
			return err
		}
		target, err := q.GetMember(ctx, in.MemberID)
		if errors.Is(err, sql.ErrNoRows) || err == nil &&
			(target.LineUserID == nil || target.DeletedAt != nil || isSuperAdmin(target)) {
			return &domain.Error{Kind: domain.Invalid, Msg: "ไม่พบสมาชิกที่มีบัญชี LINE"}
		}
		if err != nil {
			return err
		}
		if in.AllOrders {
			out.Moved, err = q.MoveMemberOrders(ctx, db.MoveMemberOrdersParams{ToMemberID: target.ID, FromMemberID: owner.ID})
		} else {
			out.Moved, err = q.MoveOrder(ctx, db.MoveOrderParams{ToMemberID: target.ID, OrderID: id})
		}
		if err != nil {
			return err
		}
		left, err := q.CountMemberOrders(ctx, owner.ID)
		if err != nil || left > 0 {
			return err
		}
		if err := q.MoveSponsorPositions(ctx, db.MoveSponsorPositionsParams{ToMemberID: &target.ID, FromMemberID: &owner.ID}); err != nil {
			return err
		}
		out.UnlinkedRemoved = true
		return q.DeleteUnlinkedMember(ctx, owner.ID)
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}
