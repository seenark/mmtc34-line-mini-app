package server

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

func (s *Server) mountProfile(r routes) {
	r.self.PUT("/me", s.editMe)
	r.self.DELETE("/me", s.deleteMe)
	r.admin.PUT("/members/:memberId", s.editMember)
	r.admin.PUT("/members/:memberId/role", s.setRole)
	r.member.GET("/me/addresses", s.myAddresses)
	r.admin.GET("/members/:memberId/addresses", s.memberAddresses)
}

// ownMember is the viewer's Member, or 404 before onboarding. A deleted Member has no
// line_user_id, so it can never be a viewer's Member.
func ownMember(c *gin.Context) (db.Member, bool) {
	m := viewerOf(c).Member
	if m == nil {
		fail(c, &domain.Error{Kind: domain.NotFound, Msg: "ไม่พบข้อมูล"})
		return db.Member{}, false
	}
	return *m, true
}

// editMe updates the profile only. Orders keep their own shipping snapshot.
func (s *Server) editMe(c *gin.Context) {
	m, ok := ownMember(c)
	if !ok {
		return
	}
	in, ok := bind[api.ProfileInput](c)
	if !ok {
		return
	}
	s.saveProfile(c, m.ID, in, domain.RequireIdentity)
}

func (s *Server) myAddresses(c *gin.Context) {
	m, ok := ownMember(c)
	if !ok {
		return
	}
	s.writeAddressChoices(c, m)
}

func (s *Server) memberAddresses(c *gin.Context) {
	id, ok := idParam(c, "memberId")
	if !ok {
		return
	}
	m, err := s.q.GetMember(c.Request.Context(), id)
	if err != nil {
		fail(c, err)
		return
	}
	s.writeAddressChoices(c, m)
}

// writeAddressChoices lists the Member's complete profile triple, then the distinct shipping
// snapshots of their earlier Orders, newest first. A triple is identified by its trimmed fields;
// the first occurrence wins, so the profile address is listed once even if an Order reused it.
func (s *Server) writeAddressChoices(c *gin.Context, m db.Member) {
	snaps, err := s.q.ListShipSnapshots(c.Request.Context(), m.ID)
	if err != nil {
		fail(c, err)
		return
	}
	out := []api.AddressChoice{}
	add := func(name, address, phone string) {
		ch := api.AddressChoice{
			Name: strings.TrimSpace(name), Address: strings.TrimSpace(address), Phone: strings.TrimSpace(phone),
		}
		if ch.Name != "" && ch.Address != "" && ch.Phone != "" && !slices.Contains(out, ch) {
			out = append(out, ch)
		}
	}
	add(m.RealName, m.Address, m.Phone)
	for _, sn := range snaps {
		add(sn.ShipName, sn.ShipAddress, sn.ShipPhone)
	}
	c.JSON(http.StatusOK, out)
}

// editMember is the Admin edit. An Unlinked Member only ever has a nickname.
func (s *Server) editMember(c *gin.Context) {
	id, ok := idParam(c, "memberId")
	if !ok {
		return
	}
	in, ok := bind[api.ProfileInput](c)
	if !ok {
		return
	}
	m, err := s.q.GetMember(c.Request.Context(), id)
	if err != nil {
		fail(c, err)
		return
	}
	req := domain.RequireIdentity
	if m.LineUserID == nil {
		req = domain.RequireNicknameOnly
	}
	s.saveProfile(c, id, in, req)
}

// saveProfile validates and writes; a deleted Member matches no row and reads as 404.
func (s *Server) saveProfile(c *gin.Context, id int64, in api.ProfileInput, req domain.ProfileRequire) {
	in, err := domain.ValidateProfile(in, req)
	if err != nil {
		fail(c, err)
		return
	}
	m, err := s.q.UpdateMemberProfile(c.Request.Context(), db.UpdateMemberProfileParams{
		ID: id, Cohort: in.Cohort, Nickname: in.Nickname, RealName: in.RealName, Address: in.Address, Phone: in.Phone,
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, memberDTO(m))
}

// deleteMe erases contact data (PDPA) but keeps the row, its nickname and every Order, Item,
// Payment and Adjustment.
func (s *Server) deleteMe(c *gin.Context) {
	m, ok := ownMember(c)
	if !ok {
		return
	}
	if isSuperAdmin(m) {
		fail(c, &domain.Error{Kind: domain.Conflict, Msg: "ลบบัญชีผู้ดูแลระบบหลักไม่ได้"})
		return
	}
	ctx := c.Request.Context()
	err := s.tx(ctx, func(q *db.Queries) error {
		now := rfc3339(s.now())
		if err := q.EraseMemberShipping(ctx, m.ID); err != nil {
			return err
		}
		return q.EraseMember(ctx, db.EraseMemberParams{ID: m.ID, Now: &now})
	})
	if err != nil {
		fail(c, err)
		return
	}
	s.clearSession(c)
	c.Status(http.StatusNoContent)
}

// setRole promotes or demotes a Member. Demoting oneself is allowed.
func (s *Server) setRole(c *gin.Context) {
	id, ok := idParam(c, "memberId")
	if !ok {
		return
	}
	in, ok := bind[api.RoleInput](c)
	if !ok {
		return
	}
	if in.Role != api.RoleMember && in.Role != api.RoleAdmin {
		fail(c, &domain.Error{Kind: domain.Invalid, Msg: "บทบาทไม่ถูกต้อง"})
		return
	}
	ctx := c.Request.Context()
	var out db.Member
	err := s.tx(ctx, func(q *db.Queries) error {
		m, err := q.GetMember(ctx, id)
		if err != nil {
			return err
		}
		if m.DeletedAt != nil {
			return &domain.Error{Kind: domain.NotFound, Msg: "ไม่พบข้อมูล"}
		}
		if isSuperAdmin(m) {
			return &domain.Error{Kind: domain.Conflict, Msg: "เปลี่ยนบทบาทผู้ดูแลระบบหลักไม่ได้"}
		}
		out, err = q.SetMemberRole(ctx, db.SetMemberRoleParams{ID: id, Role: string(in.Role)})
		return err
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, memberDTO(out))
}
