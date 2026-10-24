package domain

import "github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"

// Money is the ONLY place Total, Paid and Balance are computed. api.Order, exports and the
// confirm default amount all read it; SQL never sums money.
type Money struct {
	Total   int64 // Σ quantity × unit price
	Paid    int64 // Σ amounts of confirmed Payments
	Balance int64 // owed − Paid + Σ Adjustments; >0 owes, 0 settled, <0 overpaid
}

// Money computes the Order's Money. owed is Total, except an Order whose effective status
// is cancelled owes nothing: an unpaid Order cancelled by close must not read "ค้างจ่าย 520 บาท",
// and a cancelled paid Order reads as overpaid (refund due) until an Adjustment settles it.
func (o Order) Money() Money {
	var m Money
	for _, it := range o.Items {
		m.Total += it.Quantity * it.UnitPrice
	}
	for _, p := range o.Payments {
		if p.Status == api.PaymentConfirmed && p.Amount != nil {
			m.Paid += *p.Amount
		}
	}
	m.Balance = -m.Paid
	if o.Status() != cancelled {
		m.Balance += m.Total
	}
	for _, a := range o.Adjustments {
		m.Balance += a.Amount
	}
	return m
}
