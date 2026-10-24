package domain

import (
	"fmt"
	"slices"
	"strings"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

// ValidItem has its unit price captured: the Option price at checkout, or the donor's amount.
// ID is the input's ItemInput.ID, passed through for Admin edit (.7).
type ValidItem struct {
	ID                      *int64
	OptionID                int64
	Size, PrintName, Number string
	Quantity, UnitPrice     int64
}

// ValidOrder can only be built by ValidateOrder, so the insert path cannot skip validation.
type ValidOrder struct {
	delivery                         api.Delivery
	shipName, shipAddress, shipPhone string
	items                            []ValidItem
}

func (v ValidOrder) Delivery() api.Delivery { return v.delivery }

// Ship returns the shipping snapshot; all "" for pickup (pickup discards a typed address).
func (v ValidOrder) Ship() (name, address, phone string) {
	return v.shipName, v.shipAddress, v.shipPhone
}

func (v ValidOrder) Items() []ValidItem { return v.items }

// ValidateOrder is the checkout rule set for Member and Admin on-behalf checkout: every Item
// takes the current Option price. options are the Campaign's Options, so an unknown OptionID is
// rejected. The first failure is returned with its 1-based Item index in the Thai message.
func ValidateOrder(options []api.Option, in api.OrderInput) (ValidOrder, error) {
	return ValidateEdit(options, in, nil)
}

// ValidateEdit is ValidateOrder for Admin edit. kept holds the Order's existing Items by id: an
// input Item whose ID matches one with the same Option keeps that captured unit price (and a
// kept no-price Item keeps its amount); an ID that matches none is rejected.
func ValidateEdit(options []api.Option, in api.OrderInput, kept map[int64]Item) (ValidOrder, error) {
	bad := func(format string, args ...any) (ValidOrder, error) {
		return ValidOrder{}, &Error{Kind: Invalid, Msg: fmt.Sprintf(format, args...)}
	}
	if len(in.Items) == 0 {
		return bad("กรุณาเลือกสินค้าอย่างน้อย 1 รายการ")
	}
	v := ValidOrder{delivery: in.Delivery, items: make([]ValidItem, len(in.Items))}
	switch in.Delivery {
	case api.DeliveryPickup:
	case api.DeliveryShip:
		v.shipName = strings.TrimSpace(in.ShipName)
		v.shipAddress = strings.TrimSpace(in.ShipAddress)
		v.shipPhone = strings.TrimSpace(in.ShipPhone)
		if v.shipName == "" || v.shipAddress == "" || v.shipPhone == "" {
			return bad("กรุณากรอกชื่อ ที่อยู่ และเบอร์โทรสำหรับจัดส่ง")
		}
	default:
		return bad("วิธีรับสินค้าไม่ถูกต้อง")
	}
	for i, it := range in.Items {
		n := i + 1
		idx := slices.IndexFunc(options, func(o api.Option) bool { return o.ID == it.OptionID })
		if idx < 0 {
			return bad("รายการที่ %d: ไม่พบตัวเลือกนี้ในแคมเปญ", n)
		}
		opt := options[idx]
		var old Item
		same := false
		if it.ID != nil {
			if old, same = kept[*it.ID]; !same {
				return bad("รายการที่ %d: ไม่พบรายการนี้ในคำสั่งซื้อ", n)
			}
			same = old.OptionID == opt.ID
		}
		vi := ValidItem{ID: it.ID, OptionID: opt.ID, Size: strings.TrimSpace(it.Size), Quantity: it.Quantity}
		switch {
		case len(opt.Sizes) > 0 && !slices.Contains(opt.Sizes, vi.Size):
			return bad("รายการที่ %d: กรุณาเลือกไซซ์", n)
		case len(opt.Sizes) == 0 && vi.Size != "":
			return bad("รายการที่ %d: ตัวเลือกนี้ไม่มีไซซ์", n)
		}
		if opt.NeedsPrintName {
			if vi.PrintName = strings.TrimSpace(it.PrintName); vi.PrintName == "" {
				return bad("รายการที่ %d: กรุณากรอกชื่อที่จะพิมพ์", n)
			}
		}
		if opt.NeedsNumber {
			if vi.Number = strings.TrimSpace(it.Number); vi.Number == "" {
				return bad("รายการที่ %d: กรุณากรอกเบอร์ที่จะพิมพ์", n)
			}
		}
		if it.Quantity < 1 || it.Quantity > 99 {
			return bad("รายการที่ %d: จำนวนต้องอยู่ระหว่าง 1 ถึง 99", n)
		}
		switch {
		case opt.Price != nil && it.Amount != nil:
			return bad("รายการที่ %d: ตัวเลือกนี้มีราคาคงที่ ไม่ต้องระบุจำนวนเงิน", n)
		case opt.Price == nil && !same && (it.Amount == nil || *it.Amount <= 0):
			return bad("รายการที่ %d: กรุณาระบุจำนวนเงินที่มากกว่า 0", n)
		case opt.Price == nil && it.Quantity != 1:
			return bad("รายการที่ %d: ตัวเลือกที่ระบุจำนวนเงินเองสั่งได้ครั้งละ 1 ชิ้น", n)
		case same:
			vi.UnitPrice = old.UnitPrice
		case opt.Price != nil:
			vi.UnitPrice = *opt.Price
		default:
			vi.UnitPrice = *it.Amount
		}
		v.items[i] = vi
	}
	return v, nil
}
