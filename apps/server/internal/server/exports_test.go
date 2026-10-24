package server_test

import (
	"encoding/csv"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

func exportPath(campaignID int64, kind string) string {
	return fmt.Sprintf("/api/admin/campaigns/%d/exports/%s.csv", campaignID, kind)
}

func cartItem(j Jersey, size, printName, number string, qty int64) api.ItemInput {
	return api.ItemInput{OptionID: j.Option.ID, Size: size, PrintName: printName, Number: number, Quantity: qty}
}

// pay puts the Order through slip upload and Admin confirmation.
func pay(admin, owner *Client, o api.Order) api.Order {
	owner.UploadSlip(o.ID)
	return As[api.Order](admin.POST(confirmPath(o.ID), api.AmountInput{Amount: o.Total}).Status(200))
}

// records reads a CSV response: BOM stripped, ragged rows allowed. encoding/csv drops blank lines.
func records(t *testing.T, res *Res) [][]string {
	t.Helper()
	body, ok := strings.CutPrefix(res.Text(), "\uFEFF")
	if !ok {
		t.Fatalf("no UTF-8 BOM: %q", res.Text())
	}
	r := csv.NewReader(strings.NewReader(body))
	r.FieldsPerRecord = -1
	recs, err := r.ReadAll()
	if err != nil {
		t.Fatalf("bad CSV: %v\n%s", err, body)
	}
	return recs
}

func firstColumn(recs [][]string) []string {
	var ids []string
	for _, r := range recs[1:] {
		ids = append(ids, r[0])
	}
	return ids
}

func TestFactoryCsvListsPaidItemsWithTotals(t *testing.T) {
	h := New(t)
	alice, bob, admin := h.Member("alice"), h.Member("bob"), h.Admin()
	j := h.OpenJersey()
	a := pay(admin, alice, alice.Checkout(j.Campaign.ID, Pickup(cartItem(j, "L", "ALICE", "7", 2))))
	b := pay(admin, bob, bob.Checkout(j.Campaign.ID, Ship(cartItem(j, "L", "NOK", "9", 1), cartItem(j, "M", "NOK", "9", 1))))

	res := admin.GET(exportPath(j.Campaign.ID, "factory")).Status(200)

	Equal(t, res.Header.Get("Content-Type"), "text/csv; charset=utf-8")
	Equal(t, res.Header.Get("Content-Disposition"), fmt.Sprintf(`attachment; filename="factory-%d.csv"`, j.Campaign.ID))
	Equal(t, strings.Contains(res.Text(), "\n\n"), true)
	id := func(o api.Order) string { return strconv.FormatInt(o.ID, 10) }
	want := [][]string{
		{"สินค้า", "ขนาด", "ชื่อสกรีน", "เบอร์", "จำนวน", "คำสั่งซื้อ", "ผู้สั่ง"},
		{"เสื้อรุ่น", "M", "NOK", "9", "1", id(b), "bob"},
		{"เสื้อรุ่น", "L", "ALICE", "7", "2", id(a), "alice"},
		{"เสื้อรุ่น", "L", "NOK", "9", "1", id(b), "bob"},
		{"สินค้า", "ขนาด", "รวม"},
		{"เสื้อรุ่น", "M", "1"},
		{"เสื้อรุ่น", "L", "3"},
		{"รวมทั้งหมด", "", "4"},
	}
	got := records(t, res)
	Equal(t, len(got), len(want))
	for i := range want {
		Equal(t, slices.Equal(got[i], want[i]), true)
	}
}

func TestExportsOnlyIncludePaidUnlessPendingIsRequested(t *testing.T) {
	h := New(t)
	alice, bob, carol, dave := h.Member("alice"), h.Member("bob"), h.Member("carol"), h.Member("dave")
	admin := h.Admin()
	j := h.OpenJersey()
	paid := pay(admin, alice, alice.Checkout(j.Campaign.ID, Pickup(Shirt(j, "M", 1))))
	pending := bob.Checkout(j.Campaign.ID, Pickup(Shirt(j, "M", 1)))
	bob.UploadSlip(pending.ID)
	awaiting := carol.Checkout(j.Campaign.ID, Pickup(Shirt(j, "M", 1)))
	rejected := dave.Checkout(j.Campaign.ID, Pickup(Shirt(j, "M", 1)))
	dave.UploadSlip(rejected.ID)
	admin.POST(rejectPath(rejected.ID), api.RejectInput{Reason: "สลิปไม่ชัด"}).Status(200)
	str := func(o api.Order) string { return strconv.FormatInt(o.ID, 10) }
	pickup := exportPath(j.Campaign.ID, "pickup")

	factory := records(t, admin.GET(exportPath(j.Campaign.ID, "factory")).Status(200))
	Equal(t, len(factory), 5) // header, the one paid item, totals header, M total, grand total
	Equal(t, slices.Equal(firstColumn(records(t, admin.GET(pickup).Status(200))), []string{str(paid)}), true)
	withPending := firstColumn(records(t, admin.GET(pickup+"?includePending=1").Status(200)))
	slices.Sort(withPending)
	Equal(t, slices.Equal(withPending, []string{str(paid), str(pending)}), true)

	// the Campaign closes: Orders still awaiting payment are now cancelled, and stay out
	h.Advance(25 * time.Hour)
	Equal(t, carol.Order(awaiting.ID).Status, api.OrderCancelled)
	withPending = firstColumn(records(t, admin.GET(pickup+"?includePending=1").Status(200)))
	slices.Sort(withPending)
	Equal(t, slices.Equal(withPending, []string{str(paid), str(pending)}), true)
}

func TestShippingAndPickupCsvSplitByDelivery(t *testing.T) {
	h := New(t)
	alice, bob, admin := h.Member("alice"), h.Member("bob"), h.Admin()
	j := h.OpenJersey()
	a := pay(admin, alice, alice.Checkout(j.Campaign.ID, Pickup(cartItem(j, "L", "ALICE", "7", 2))))
	b := pay(admin, bob, bob.Checkout(j.Campaign.ID, Ship(cartItem(j, "L", "NOK", "9", 2), cartItem(j, "S", "NOK", "9", 1))))

	shipping := admin.GET(exportPath(j.Campaign.ID, "shipping")).Status(200)
	Equal(t, shipping.Header.Get("Content-Disposition"), fmt.Sprintf(`attachment; filename="shipping-%d.csv"`, j.Campaign.ID))
	recs := records(t, shipping)
	Equal(t, slices.Equal(recs[0], []string{"คำสั่งซื้อ", "ชื่อผู้รับ", "ที่อยู่", "เบอร์โทร", "รายการ", "สถานะ"}), true)
	Equal(t, len(recs), 2)
	Equal(t, slices.Equal(recs[1], []string{
		strconv.FormatInt(b.ID, 10), "สมชาย ใจดี", "1 ถนนสุขุมวิท กรุงเทพฯ 10110", "0812345678",
		"เสื้อรุ่น L NOK #9 x2; เสื้อรุ่น S NOK #9 x1", "ชำระแล้ว",
	}), true)

	pickup := admin.GET(exportPath(j.Campaign.ID, "pickup")).Status(200)
	Equal(t, pickup.Header.Get("Content-Disposition"), fmt.Sprintf(`attachment; filename="pickup-%d.csv"`, j.Campaign.ID))
	recs = records(t, pickup)
	Equal(t, slices.Equal(recs[0], []string{"คำสั่งซื้อ", "ชื่อเล่น", "ชื่อจริง", "รายการ", "สถานะ"}), true)
	Equal(t, len(recs), 2)
	Equal(t, recs[1][0], strconv.FormatInt(a.ID, 10))
	Equal(t, recs[1][1], "alice")
	Equal(t, recs[1][3], "เสื้อรุ่น L ALICE #7 x2")
	Equal(t, recs[1][4], "ชำระแล้ว")
}

func TestExportPrefixesFormulaCells(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()
	pay(admin, alice, alice.Checkout(j.Campaign.ID, Pickup(cartItem(j, "M", "=1+1", "7", 1))))

	recs := records(t, admin.GET(exportPath(j.Campaign.ID, "factory")).Status(200))

	Equal(t, recs[1][2], "'=1+1")
}

func TestExportsAreAdminOnlyAndNeedACampaign(t *testing.T) {
	h := New(t)
	alice, admin := h.Member("alice"), h.Admin()
	j := h.OpenJersey()

	for _, kind := range []string{"factory", "shipping", "pickup"} {
		alice.GET(exportPath(j.Campaign.ID, kind)).Status(403)
		h.Anonymous().GET(exportPath(j.Campaign.ID, kind)).Status(401)
		admin.GET(exportPath(9999, kind)).Status(404)
	}
}
