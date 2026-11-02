package server

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

func (s *Server) mountCampaigns(r routes) {
	r.member.GET("/campaigns", s.listCampaigns)              // non-draft, effective status
	r.member.GET("/campaigns/:campaignId", s.campaignDetail) // draft -> 404 for Members
	r.admin.GET("/campaigns", s.adminListCampaigns)          // includes drafts
	r.admin.GET("/campaigns/:campaignId", s.campaignDetail)
	r.admin.POST("/campaigns", s.createCampaign)
	r.admin.PUT("/campaigns/:campaignId", s.updateCampaign) // also open / close / reopen
	r.admin.POST("/campaigns/:campaignId/options", s.createOption)
	r.admin.PUT("/options/:optionId", s.updateOption) // includes sortOrder (reorder)
	r.admin.DELETE("/options/:optionId", s.deleteOption)
}

// domainCampaign is the rule-facing view of a row. close_at is only ever written by
// ValidateCampaign (normalised RFC3339) or the seed, so a parse failure cannot happen.
func domainCampaign(row db.Campaign) domain.Campaign {
	c := domain.Campaign{Stored: api.CampaignStatus(row.Status)}
	if row.CloseAt != nil {
		if t, err := time.Parse(time.RFC3339, *row.CloseAt); err == nil {
			c.CloseAt = &t
		}
	}
	return c
}

// toAPICampaign is the only db.Campaign -> api.Campaign mapping; it writes the effective status.
// photos are the row's Campaign Photo names in order; nil serializes as [].
func (s *Server) toAPICampaign(row db.Campaign, photos []string) api.Campaign {
	if photos == nil {
		photos = []string{}
	}
	return api.Campaign{
		ID:              row.ID,
		Type:            api.CampaignType(row.Type),
		Title:           row.Title,
		Description:     row.Description,
		Photos:          photos,
		BankCode:        row.BankCode,
		BankName:        domain.BankName(row.BankCode),
		BankAccountNo:   row.BankAccountNo,
		BankAccountName: row.BankAccountName,
		QRImage:         row.QrImage,
		Status:          domainCampaign(row).Status(s.now()),
		CloseAt:         row.CloseAt,
	}
}

// replacePhotos makes names the Campaign's whole photo set; array index is the sort order.
func replacePhotos(ctx context.Context, q *db.Queries, campaignID int64, names []string) error {
	if err := q.DeleteCampaignPhotos(ctx, campaignID); err != nil {
		return err
	}
	for i, name := range names {
		err := q.CreateCampaignPhoto(ctx, db.CreateCampaignPhotoParams{CampaignID: campaignID, Name: name, SortOrder: int64(i)})
		if err != nil {
			return err
		}
	}
	return nil
}

func toAPIOption(o db.Option) api.Option {
	sizes := []string{}
	_ = json.Unmarshal([]byte(o.Sizes), &sizes) // column is written only from validated []string
	return api.Option{
		ID:             o.ID,
		CampaignID:     o.CampaignID,
		Name:           o.Name,
		Description:    o.Description,
		Image:          o.Image,
		SizeChart:      o.SizeChart,
		Price:          o.Price,
		Sizes:          sizes,
		NeedsPrintName: o.NeedsPrintName == 1,
		NeedsNumber:    o.NeedsNumber == 1,
		SortOrder:      o.SortOrder,
	}
}

// loadCampaign returns the Campaign if the actor may see it (domain.Campaign.VisibleTo), else a
// NotFound error. Only valid behind requireApproved or requireAdmin.
func (s *Server) loadCampaign(c *gin.Context, q *db.Queries, id int64) (db.Campaign, error) {
	row, err := q.GetCampaign(c.Request.Context(), id)
	if err != nil {
		return row, err
	}
	if !domainCampaign(row).VisibleTo(actorOf(c)) {
		return row, &domain.Error{Kind: domain.NotFound, Msg: "ไม่พบข้อมูล"}
	}
	return row, nil
}

func (s *Server) listCampaigns(c *gin.Context) {
	rows, err := s.q.ListPublishedCampaigns(c.Request.Context())
	s.writeCampaigns(c, rows, err)
}

func (s *Server) adminListCampaigns(c *gin.Context) {
	rows, err := s.q.ListCampaigns(c.Request.Context())
	s.writeCampaigns(c, rows, err)
}

func (s *Server) writeCampaigns(c *gin.Context, rows []db.Campaign, err error) {
	if err != nil {
		fail(c, err)
		return
	}
	ids := make([]int64, len(rows))
	for i, row := range rows {
		ids[i] = row.ID
	}
	photoRows, err := s.q.ListPhotosForCampaigns(c.Request.Context(), ids)
	if err != nil {
		fail(c, err)
		return
	}
	photos := make(map[int64][]string, len(rows))
	for _, p := range photoRows {
		photos[p.CampaignID] = append(photos[p.CampaignID], p.Name)
	}
	out := make([]api.Campaign, len(rows))
	for i, row := range rows {
		out[i] = s.toAPICampaign(row, photos[row.ID])
	}
	c.JSON(http.StatusOK, out)
}

func (s *Server) campaignDetail(c *gin.Context) {
	id, ok := idParam(c, "campaignId")
	if !ok {
		return
	}
	row, err := s.loadCampaign(c, s.q, id)
	if err != nil {
		fail(c, err)
		return
	}
	opts, err := s.campaignOptions(c.Request.Context(), s.q, id)
	if err != nil {
		fail(c, err)
		return
	}
	photos, err := s.q.ListCampaignPhotos(c.Request.Context(), id)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, api.CampaignDetail{Campaign: s.toAPICampaign(row, photos), Options: opts})
}

// bindCampaign decodes, validates and checks image names. ok=false means the response is written.
func (s *Server) bindCampaign(c *gin.Context) (in api.CampaignInput, ok bool) {
	in, ok = bind[api.CampaignInput](c)
	if !ok {
		return in, false
	}
	in, err := domain.ValidateCampaign(in)
	if err == nil {
		err = s.requirePhotos(in.Photos)
	}
	if err == nil {
		err = s.requireImages(in.QRImage)
	}
	if err != nil {
		fail(c, err)
		return in, false
	}
	return in, true
}

func (s *Server) createCampaign(c *gin.Context) {
	in, ok := s.bindCampaign(c)
	if !ok {
		return
	}
	var row db.Campaign
	err := s.tx(c.Request.Context(), func(q *db.Queries) (err error) {
		row, err = q.CreateCampaign(c.Request.Context(), db.CreateCampaignParams{
			Type:            string(in.Type),
			Title:           in.Title,
			Description:     in.Description,
			BankCode:        in.BankCode,
			BankAccountNo:   in.BankAccountNo,
			BankAccountName: in.BankAccountName,
			QrImage:         in.QRImage,
			Status:          string(in.Status),
			CloseAt:         in.CloseAt,
			CreatedAt:       rfc3339(s.now()),
		})
		if err != nil {
			return err
		}
		return replacePhotos(c.Request.Context(), q, row.ID, in.Photos)
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, s.toAPICampaign(row, in.Photos))
}

func (s *Server) updateCampaign(c *gin.Context) {
	id, ok := idParam(c, "campaignId")
	if !ok {
		return
	}
	in, ok := s.bindCampaign(c)
	if !ok {
		return
	}
	var row db.Campaign
	err := s.tx(c.Request.Context(), func(q *db.Queries) (err error) {
		row, err = q.UpdateCampaign(c.Request.Context(), db.UpdateCampaignParams{
			ID:              id,
			Type:            string(in.Type),
			Title:           in.Title,
			Description:     in.Description,
			BankCode:        in.BankCode,
			BankAccountNo:   in.BankAccountNo,
			BankAccountName: in.BankAccountName,
			QrImage:         in.QRImage,
			Status:          string(in.Status),
			CloseAt:         in.CloseAt,
		})
		if err != nil {
			return err
		}
		return replacePhotos(c.Request.Context(), q, id, in.Photos)
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, s.toAPICampaign(row, in.Photos))
}

// bindOption decodes, validates and checks the image name; sizes is the JSON column value.
func (s *Server) bindOption(c *gin.Context) (in api.OptionInput, sizes string, ok bool) {
	in, ok = bind[api.OptionInput](c)
	if !ok {
		return in, "", false
	}
	in, err := domain.ValidateOption(in)
	if err == nil {
		err = s.requireImages(in.Image, in.SizeChart)
	}
	if err != nil {
		fail(c, err)
		return in, "", false
	}
	raw, _ := json.Marshal(in.Sizes) // []string cannot fail to marshal
	return in, string(raw), true
}

func flag(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func (s *Server) createOption(c *gin.Context) {
	campaignID, ok := idParam(c, "campaignId")
	if !ok {
		return
	}
	in, sizes, ok := s.bindOption(c)
	if !ok {
		return
	}
	// Explicit lookup: a missing Campaign must read 404, not the FK-violation 409 of fail().
	if _, err := s.q.GetCampaign(c.Request.Context(), campaignID); err != nil {
		fail(c, err)
		return
	}
	row, err := s.q.CreateOption(c.Request.Context(), db.CreateOptionParams{
		CampaignID:     campaignID,
		Name:           in.Name,
		Description:    in.Description,
		Image:          in.Image,
		SizeChart:      in.SizeChart,
		Price:          in.Price,
		Sizes:          sizes,
		NeedsPrintName: flag(in.NeedsPrintName),
		NeedsNumber:    flag(in.NeedsNumber),
		SortOrder:      in.SortOrder,
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, toAPIOption(row))
}

func (s *Server) updateOption(c *gin.Context) {
	id, ok := idParam(c, "optionId")
	if !ok {
		return
	}
	in, sizes, ok := s.bindOption(c)
	if !ok {
		return
	}
	row, err := s.q.UpdateOption(c.Request.Context(), db.UpdateOptionParams{
		ID:             id,
		Name:           in.Name,
		Description:    in.Description,
		Image:          in.Image,
		SizeChart:      in.SizeChart,
		Price:          in.Price,
		Sizes:          sizes,
		NeedsPrintName: flag(in.NeedsPrintName),
		NeedsNumber:    flag(in.NeedsNumber),
		SortOrder:      in.SortOrder,
	})
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusOK, toAPIOption(row))
}

// deleteOption: the items.option_id RESTRICT foreign key makes an Option in use fail with the 409 in fail().
func (s *Server) deleteOption(c *gin.Context) {
	id, ok := idParam(c, "optionId")
	if !ok {
		return
	}
	n, err := s.q.DeleteOption(c.Request.Context(), id)
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
