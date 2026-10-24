package server

import (
	"bytes"
	"cmp"
	"encoding/csv"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

// exportBuilders maps the file kind in the URL to the records it holds. A nil record is a blank line.
var exportBuilders = map[string]func(orders []domain.Order, options []api.Option) [][]string{
	"factory":  factoryRecords,
	"shipping": shippingRecords,
	"pickup":   pickupRecords,
}

func (s *Server) mountExports(r routes) {
	for kind := range exportBuilders {
		r.admin.GET("/campaigns/:campaignId/exports/"+kind+".csv", s.exportCSV(kind))
	}
}

// exportCSV serves the paid Orders of a Campaign (plus pending_review with ?includePending=1);
// awaiting_payment and cancelled never appear. Eligibility uses effective Status.
func (s *Server) exportCSV(kind string) gin.HandlerFunc {
	build := exportBuilders[kind]
	return func(c *gin.Context) {
		id, ok := idParam(c, "campaignId")
		if !ok {
			return
		}
		ctx := c.Request.Context()
		if _, err := s.loadCampaign(c, s.q, id); err != nil {
			fail(c, err)
			return
		}
		all, err := s.loadOrders(ctx, s.q, orderScope{CampaignID: &id})
		if err != nil {
			fail(c, err)
			return
		}
		options, err := s.campaignOptions(ctx, s.q, id)
		if err != nil {
			fail(c, err)
			return
		}
		includePending := c.Query("includePending") == "1"
		orders := slices.DeleteFunc(all, func(o domain.Order) bool {
			st := o.Status()
			return st != api.OrderPaid && !(includePending && st == api.OrderPendingReview)
		})

		var buf bytes.Buffer
		buf.WriteString("\uFEFF") // Excel reads Thai as UTF-8 only with a BOM
		w := csv.NewWriter(&buf)
		for _, rec := range build(orders, options) {
			for i, cell := range rec {
				rec[i] = defuseFormula(cell)
			}
			if err := w.Write(rec); err != nil {
				fail(c, err)
				return
			}
		}
		w.Flush()
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s-%d.csv"`, kind, id))
		c.Data(http.StatusOK, "text/csv; charset=utf-8", buf.Bytes())
	}
}

// defuseFormula stops spreadsheets running user-typed text (print names, nicknames, addresses) as a formula.
func defuseFormula(cell string) string {
	if cell != "" && strings.ContainsRune("=+-@", rune(cell[0])) {
		return "'" + cell
	}
	return cell
}

type factoryRow struct {
	optRank, sizeRank int
	item              domain.Item
	order             domain.Order
}

// factoryRecords lists one row per Item ordered by Option sort order, then the Option's size
// list, then Order id; then quantity totals per Option and size.
func factoryRecords(orders []domain.Order, options []api.Option) [][]string {
	optRank := map[int64]int{}
	sizeRank := map[int64]map[string]int{}
	for i, o := range options {
		optRank[o.ID] = i
		sizeRank[o.ID] = map[string]int{}
		for j, sz := range o.Sizes {
			sizeRank[o.ID][sz] = j
		}
	}
	// a size missing from the Option's list (Option edited after ordering) sorts last
	rank := func(it domain.Item) (int, int) {
		sr, ok := sizeRank[it.OptionID][it.Size]
		if !ok {
			sr = len(sizeRank[it.OptionID])
		}
		return optRank[it.OptionID], sr
	}

	var rows []factoryRow
	for _, o := range orders {
		for _, it := range o.Items {
			or, sr := rank(it)
			rows = append(rows, factoryRow{or, sr, it, o})
		}
	}
	slices.SortFunc(rows, func(a, b factoryRow) int {
		return cmp.Or(
			cmp.Compare(a.optRank, b.optRank),
			cmp.Compare(a.sizeRank, b.sizeRank),
			cmp.Compare(a.item.Size, b.item.Size),
			cmp.Compare(a.order.ID, b.order.ID),
		)
	})

	out := [][]string{{"สินค้า", "ขนาด", "ชื่อสกรีน", "เบอร์", "จำนวน", "คำสั่งซื้อ", "ผู้สั่ง"}}
	totals := [][]string{{"สินค้า", "ขนาด", "รวม"}}
	var sum, grand int64
	for i, r := range rows {
		it := r.item
		out = append(out, []string{
			it.OptionName, it.Size, it.PrintName, it.Number,
			strconv.FormatInt(it.Quantity, 10), strconv.FormatInt(r.order.ID, 10), r.order.Nickname,
		})
		sum += it.Quantity
		grand += it.Quantity
		if i+1 == len(rows) || rows[i+1].item.OptionID != it.OptionID || rows[i+1].item.Size != it.Size {
			totals = append(totals, []string{it.OptionName, it.Size, strconv.FormatInt(sum, 10)})
			sum = 0
		}
	}
	totals = append(totals, []string{"รวมทั้งหมด", "", strconv.FormatInt(grand, 10)})
	out = append(out, nil)
	return append(out, totals...)
}

var exportStatusLabel = map[api.OrderStatus]string{
	api.OrderPaid:          "ชำระแล้ว",
	api.OrderPendingReview: "รอตรวจ",
}

func itemsSummary(items []domain.Item) string {
	parts := make([]string, len(items))
	for i, it := range items {
		p := []string{it.OptionName, it.Size}
		if it.PrintName != "" {
			p = append(p, it.PrintName)
		}
		if it.Number != "" {
			p = append(p, "#"+it.Number)
		}
		parts[i] = strings.Join(p, " ") + " x" + strconv.FormatInt(it.Quantity, 10)
	}
	return strings.Join(parts, "; ")
}

func shippingRecords(orders []domain.Order, _ []api.Option) [][]string {
	out := [][]string{{"คำสั่งซื้อ", "ชื่อผู้รับ", "ที่อยู่", "เบอร์โทร", "รายการ", "สถานะ"}}
	for _, o := range orders {
		if o.Delivery == api.DeliveryShip {
			out = append(out, []string{
				strconv.FormatInt(o.ID, 10), o.ShipName, o.ShipAddress, o.ShipPhone,
				itemsSummary(o.Items), exportStatusLabel[o.Status()],
			})
		}
	}
	return out
}

func pickupRecords(orders []domain.Order, _ []api.Option) [][]string {
	out := [][]string{{"คำสั่งซื้อ", "ชื่อเล่น", "ชื่อจริง", "รายการ", "สถานะ"}}
	for _, o := range orders {
		if o.Delivery == api.DeliveryPickup {
			out = append(out, []string{
				strconv.FormatInt(o.ID, 10), o.Nickname, o.RealName,
				itemsSummary(o.Items), exportStatusLabel[o.Status()],
			})
		}
	}
	return out
}
