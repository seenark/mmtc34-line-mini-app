package db_test

import (
	"context"
	"database/sql"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

func TestFootball2026MigrationKeepsOrdersAndWidensTheSchema(t *testing.T) {
	ctx := context.Background()
	q := url.Values{}
	for _, p := range []string{"foreign_keys(1)", "journal_mode(WAL)", "busy_timeout(5000)"} {
		q.Add("_pragma", p)
	}
	q.Set("_txlock", "immediate")
	sqlDB, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "app.db")+"?"+q.Encode())
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()

	provider, err := goose.NewProvider(goose.DialectSQLite3, sqlDB, os.DirFS("migrations"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := provider.UpTo(ctx, 1); err != nil {
		t.Fatal(err)
	}

	exec := func(query string, args ...any) {
		t.Helper()
		if _, err := sqlDB.ExecContext(ctx, query, args...); err != nil {
			t.Fatalf("%s: %v", query, err)
		}
	}
	exec(`INSERT INTO members (id, line_user_id, status, role, nickname, created_at) VALUES
	        (1, 'U-admin', 'approved', 'admin', 'admin', '2026-10-01T00:00:00Z'),
	        (2, 'U-alice', 'applicant', 'member', 'alice', '2026-10-01T00:00:00Z')`)
	exec(`INSERT INTO campaigns (id, type, title, bank_name, bank_account_no, qr_image, status, created_at) VALUES
	        (10, 'merch', 'เสื้อรุ่น', 'ธนาคารกสิกรไทย', '000-0-00000-0', 'qr.png', 'open', '2026-10-01T00:00:00Z'),
	        (11, 'donation', 'ระดมทุน', 'ธนาคารไทยพาณิชย์', '111-1-11111-1', '', 'open', '2026-10-01T00:00:00Z')`)
	exec(`INSERT INTO options (id, campaign_id, name, price) VALUES (20, 10, 'เสื้อ', 260)`)
	exec(`INSERT INTO orders (id, campaign_id, member_id, status, delivery, ship_name, ship_address, ship_phone, created_at, updated_at)
	        VALUES (30, 10, 2, 'paid', 'ship', 'สมชาย', '1 ถนนสุขุมวิท', '0812345678', '2026-10-02T00:00:00Z', '2026-10-02T00:00:00Z')`)
	exec(`INSERT INTO items (id, order_id, option_id, size, quantity, unit_price) VALUES (40, 30, 20, 'M', 2, 260)`)
	exec(`INSERT INTO payments (id, order_id, slip, status, amount, created_at, reviewed_by, reviewed_at)
	        VALUES (50, 30, 'slip.jpg', 'confirmed', 520, '2026-10-02T00:00:00Z', 1, '2026-10-02T01:00:00Z')`)
	exec(`INSERT INTO adjustments (id, order_id, amount, note, created_by, created_at)
	        VALUES (60, 30, -20, 'ส่วนลด', 1, '2026-10-03T00:00:00Z')`)

	if _, err := provider.Up(ctx); err != nil {
		t.Fatal(err)
	}

	count := func(query string, args ...any) (n int) {
		t.Helper()
		if err := sqlDB.QueryRowContext(ctx, query, args...).Scan(&n); err != nil {
			t.Fatalf("%s: %v", query, err)
		}
		return n
	}
	text := func(query string, args ...any) (s string) {
		t.Helper()
		if err := sqlDB.QueryRowContext(ctx, query, args...).Scan(&s); err != nil {
			t.Fatalf("%s: %v", query, err)
		}
		return s
	}

	if got := count(`SELECT count(*) FROM items WHERE id = 40 AND order_id = 30 AND quantity = 2 AND unit_price = 260`); got != 1 {
		t.Errorf("item 40 of order 30 survived %d times, want 1", got)
	}
	if got := count(`SELECT count(*) FROM payments WHERE id = 50 AND order_id = 30 AND amount = 520`); got != 1 {
		t.Errorf("payment 50 of order 30 survived %d times, want 1", got)
	}
	if got := count(`SELECT count(*) FROM adjustments WHERE id = 60 AND order_id = 30 AND amount = -20`); got != 1 {
		t.Errorf("adjustment 60 of order 30 survived %d times, want 1", got)
	}
	if got := count(`SELECT count(*) FROM orders WHERE id = 30 AND member_id = 2 AND campaign_id = 10 AND status = 'paid'
	                    AND delivery = 'ship' AND ship_name = 'สมชาย' AND ship_address = '1 ถนนสุขุมวิท' AND ship_phone = '0812345678'`); got != 1 {
		t.Errorf("order 30 changed: matched %d rows, want 1", got)
	}
	if got := text(`SELECT status FROM members WHERE id = 2`); got != "approved" {
		t.Errorf("applicant member status = %q, want approved", got)
	}
	if got := text(`SELECT bank_code FROM campaigns WHERE id = 10`); got != "004" {
		t.Errorf("kasikorn campaign bank_code = %q, want 004", got)
	}
	if got := text(`SELECT bank_code FROM campaigns WHERE id = 11`); got != "014" {
		t.Errorf("scb campaign bank_code = %q, want 014", got)
	}
	if got := text(`SELECT qr_image FROM campaigns WHERE id = 10`); got != "qr.png" {
		t.Errorf("campaign qr_image = %q, want qr.png", got)
	}
	if got := count(`SELECT count(*) FROM pragma_table_info('campaigns') WHERE name = 'qr_source'`); got != 0 {
		t.Errorf("campaigns.qr_source survived the full migration")
	}
	if got := text(`SELECT size_chart FROM options WHERE id = 20`); got != "" {
		t.Errorf("existing option size_chart = %q, want empty", got)
	}

	exec(`INSERT INTO orders (id, campaign_id, member_id, status, delivery, created_at, updated_at)
	        VALUES (31, 11, 2, 'awaiting_payment', 'none', '2026-10-04T00:00:00Z', '2026-10-04T00:00:00Z')`)
	if _, err := sqlDB.ExecContext(ctx, `INSERT INTO orders (id, campaign_id, member_id, status, delivery, created_at, updated_at)
	        VALUES (32, 11, 2, 'awaiting_payment', 'courier', '2026-10-04T00:00:00Z', '2026-10-04T00:00:00Z')`); err == nil {
		t.Error("an order with delivery 'courier' was accepted")
	}
	if _, err := sqlDB.ExecContext(ctx, `INSERT INTO orders (id, campaign_id, member_id, status, delivery, created_at, updated_at)
	        VALUES (33, 999, 2, 'awaiting_payment', 'none', '2026-10-04T00:00:00Z', '2026-10-04T00:00:00Z')`); err == nil {
		t.Error("an order of a missing campaign was accepted: foreign keys are off after the migration")
	}
}
