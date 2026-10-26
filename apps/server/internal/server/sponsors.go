package server

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

func (s *Server) mountSponsors(r routes) {
	r.member.GET("/campaigns/:campaignId/sponsors", s.listSponsors)
	r.admin.POST("/campaigns/:campaignId/sponsors", s.createSponsor)
	r.admin.PUT("/sponsors/:sponsorId", s.updateSponsor) // includes sortOrder (reorder)
	r.admin.DELETE("/sponsors/:sponsorId", s.deleteSponsor)
}

func toAPISponsor(p db.SponsorPosition, memberNickname string) api.SponsorPosition {
	return api.SponsorPosition{
		ID:             p.ID,
		CampaignID:     p.CampaignID,
		Name:           p.Name,
		SponsorName:    p.SponsorName,
		Logo:           p.Logo,
		MemberID:       p.MemberID,
		MemberNickname: memberNickname,
		SortOrder:      p.SortOrder,
	}
}

func (s *Server) listSponsors(c *gin.Context) {
	id, ok := idParam(c, "campaignId")
	if !ok {
		return
	}
	if _, err := s.loadCampaign(c, s.q, id); err != nil {
		fail(c, err)
		return
	}
	rows, err := s.q.ListSponsors(c.Request.Context(), id)
	if err != nil {
		fail(c, err)
		return
	}
	out := make([]api.SponsorPosition, len(rows))
	for i, row := range rows {
		out[i] = toAPISponsor(row.SponsorPosition, row.MemberNickname)
	}
	c.JSON(http.StatusOK, out)
}

// bindSponsor decodes and validates the input, checks the logo name and resolves the linked
// Member's nickname. ok=false means the response is written.
func (s *Server) bindSponsor(c *gin.Context) (in api.SponsorInput, nickname string, ok bool) {
	in, ok = bind[api.SponsorInput](c)
	if !ok {
		return in, "", false
	}
	in.Name = strings.TrimSpace(in.Name)
	in.SponsorName = strings.TrimSpace(in.SponsorName)
	var err error
	switch {
	case in.Name == "":
		err = &domain.Error{Kind: domain.Invalid, Msg: "กรุณากรอกชื่อตำแหน่ง"}
	case in.SortOrder < 0:
		err = &domain.Error{Kind: domain.Invalid, Msg: "ลำดับต้องไม่ติดลบ"}
	default:
		err = s.requireImages(in.Logo)
	}
	if err == nil && in.MemberID != nil {
		nickname, err = s.sponsorMemberNickname(c, *in.MemberID)
	}
	if err != nil {
		fail(c, err)
		return in, "", false
	}
	return in, nickname, true
}

func (s *Server) sponsorMemberNickname(c *gin.Context, id int64) (string, error) {
	m, err := s.q.GetMember(c.Request.Context(), id)
	switch {
	case errors.Is(err, sql.ErrNoRows), err == nil && m.DeletedAt != nil:
		return "", &domain.Error{Kind: domain.Invalid, Msg: "ไม่พบสมาชิกที่เลือก"}
	case err != nil:
		return "", err
	}
	return m.Nickname, nil
}

func (s *Server) createSponsor(c *gin.Context) {
	campaignID, ok := idParam(c, "campaignId")
	if !ok {
		return
	}
	in, nickname, ok := s.bindSponsor(c)
	if !ok {
		return
	}
	// Explicit lookup: a missing Campaign must read 404, not the FK-violation 409 of fail().
	if _, err := s.q.GetCampaign(c.Request.Context(), campaignID); err != nil {
		fail(c, err)
		return
	}
	row, err := s.q.CreateSponsor(c.Request.Context(), db.CreateSponsorParams{
		CampaignID:  campaignID,
		Name:        in.Name,
		SponsorName: in.SponsorName,
		Logo:        in.Logo,
		MemberID:    in.MemberID,
		SortOrder:   in.SortOrder,
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, toAPISponsor(row, nickname))
}

func (s *Server) updateSponsor(c *gin.Context) {
	id, ok := idParam(c, "sponsorId")
	if !ok {
		return
	}
	in, nickname, ok := s.bindSponsor(c)
	if !ok {
		return
	}
	row, err := s.q.UpdateSponsor(c.Request.Context(), db.UpdateSponsorParams{
		ID:          id,
		Name:        in.Name,
		SponsorName: in.SponsorName,
		Logo:        in.Logo,
		MemberID:    in.MemberID,
		SortOrder:   in.SortOrder,
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, toAPISponsor(row, nickname))
}

func (s *Server) deleteSponsor(c *gin.Context) {
	id, ok := idParam(c, "sponsorId")
	if !ok {
		return
	}
	n, err := s.q.DeleteSponsor(c.Request.Context(), id)
	if err != nil {
		fail(c, err)
		return
	}
	if n == 0 {
		fail(c, &domain.Error{Kind: domain.NotFound, Msg: "ไม่พบข้อมูล"})
		return
	}
	c.Status(http.StatusNoContent)
}
