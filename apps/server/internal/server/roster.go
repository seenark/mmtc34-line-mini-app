package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

func (s *Server) mountRoster(r routes) {
	r.member.GET("/campaigns/:campaignId/roster", s.roster)
}

// roster is the Member-facing "who ordered" list. It is built from api.RosterEntry only, so no
// id, real name, address, phone, money or status can leak; eligibility is the effective status.
func (s *Server) roster(c *gin.Context) {
	id, ok := idParam(c, "campaignId")
	if !ok {
		return
	}
	ctx := c.Request.Context()
	if _, err := s.loadCampaign(c, s.q, id); err != nil {
		fail(c, err)
		return
	}
	orders, err := s.loadOrders(ctx, s.q, orderScope{CampaignID: &id})
	if err != nil {
		fail(c, err)
		return
	}

	type line struct{ option, size string }
	type group struct {
		entry api.RosterEntry
		at    map[line]int
	}
	byMember := map[int64]*group{}
	out := []api.RosterEntry{}
	groups := []*group{}
	// loadOrders is newest first; walk backwards so entries follow each Member's earliest Order.
	for i := len(orders) - 1; i >= 0; i-- {
		o := orders[i]
		if st := o.Status(); st != api.OrderPendingReview && st != api.OrderPaid {
			continue
		}
		g := byMember[o.MemberID]
		if g == nil {
			g = &group{entry: api.RosterEntry{Nickname: o.Nickname, Items: []api.RosterItem{}}, at: map[line]int{}}
			byMember[o.MemberID] = g
			groups = append(groups, g)
		}
		for _, it := range o.Items {
			k := line{it.OptionName, it.Size}
			if j, seen := g.at[k]; seen {
				g.entry.Items[j].Quantity += it.Quantity
				continue
			}
			g.at[k] = len(g.entry.Items)
			g.entry.Items = append(g.entry.Items, api.RosterItem{OptionName: it.OptionName, Size: it.Size, Quantity: it.Quantity})
		}
	}
	for _, g := range groups {
		out = append(out, g.entry)
	}
	c.JSON(http.StatusOK, out)
}
