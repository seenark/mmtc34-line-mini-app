-- MMTC Events schema change:
--   * open membership (ADR 0006): every Applicant becomes approved; the status value stays.
--   * options.size_chart: stored image name of the Option's size chart ('' = none).
--   * campaigns.bank_code (3-digit Bank of Thailand code) replaces free-text bank_name;
--     campaigns.qr_source picks the generated Thai QR or the uploaded qr_image.
--   * orders.delivery gains 'none' for donation-only Orders.
-- Widening a CHECK in SQLite means rebuilding orders. Dropping orders with foreign keys on would
-- cascade-delete its items, so the whole change runs as ONE statement on ONE connection with
-- foreign keys off, inside its own transaction.

-- +goose NO TRANSACTION
-- +goose Up
-- +goose StatementBegin
PRAGMA foreign_keys = OFF;
BEGIN;

UPDATE members SET status = 'approved' WHERE status = 'applicant';

ALTER TABLE options ADD COLUMN size_chart TEXT NOT NULL DEFAULT '';

ALTER TABLE campaigns ADD COLUMN bank_code TEXT NOT NULL DEFAULT '';
UPDATE campaigns SET bank_code = CASE
    WHEN bank_name LIKE '%กสิกร%' OR lower(bank_name) LIKE '%kasikorn%' OR lower(bank_name) LIKE '%kbank%' THEN '004'
    WHEN bank_name LIKE '%กรุงเทพ%' OR lower(bank_name) LIKE '%bangkok bank%' OR lower(bank_name) LIKE '%bbl%' THEN '002'
    WHEN bank_name LIKE '%กรุงไทย%' OR lower(bank_name) LIKE '%krungthai%' OR lower(bank_name) LIKE '%ktb%' THEN '006'
    WHEN bank_name LIKE '%ทหารไทย%' OR bank_name LIKE '%ทีทีบี%' OR lower(bank_name) LIKE '%ttb%' THEN '011'
    WHEN bank_name LIKE '%ไทยพาณิชย์%' OR lower(bank_name) LIKE '%scb%' OR lower(bank_name) LIKE '%siam commercial%' THEN '014'
    WHEN bank_name LIKE '%กรุงศรี%' OR lower(bank_name) LIKE '%krungsri%' OR lower(bank_name) LIKE '%ayudhya%' THEN '025'
    WHEN bank_name LIKE '%ออมสิน%' OR lower(bank_name) LIKE '%gsb%' OR lower(bank_name) LIKE '%government savings%' THEN '030'
    ELSE ''
END;
ALTER TABLE campaigns DROP COLUMN bank_name;

ALTER TABLE campaigns ADD COLUMN qr_source TEXT NOT NULL DEFAULT 'generated'
    CHECK (qr_source IN ('generated', 'uploaded'));
-- Campaigns that already show an uploaded QR keep showing it.
UPDATE campaigns SET qr_source = 'uploaded' WHERE qr_image <> '';

CREATE TABLE orders_new (
    id           INTEGER PRIMARY KEY,
    campaign_id  INTEGER NOT NULL REFERENCES campaigns (id),
    member_id    INTEGER NOT NULL REFERENCES members (id),
    -- 1 = an Admin placed the Order on the Member's behalf. Such Orders are exempt from
    -- "awaiting_payment in a closed Campaign reads cancelled".
    placed_by_admin INTEGER NOT NULL DEFAULT 0 CHECK (placed_by_admin IN (0, 1)),
    status       TEXT    NOT NULL CHECK (status IN ('awaiting_payment', 'pending_review', 'paid', 'cancelled')),
    -- 'none' = the Order holds only donations, so nothing is delivered.
    delivery     TEXT    NOT NULL CHECK (delivery IN ('ship', 'pickup', 'none')),
    ship_name    TEXT    NOT NULL DEFAULT '',
    ship_address TEXT    NOT NULL DEFAULT '',
    ship_phone   TEXT    NOT NULL DEFAULT '',
    created_at   TEXT    NOT NULL,
    updated_at   TEXT    NOT NULL
);
INSERT INTO orders_new (id, campaign_id, member_id, placed_by_admin, status, delivery,
                        ship_name, ship_address, ship_phone, created_at, updated_at)
SELECT id, campaign_id, member_id, placed_by_admin, status, delivery,
       ship_name, ship_address, ship_phone, created_at, updated_at
FROM orders;
DROP TABLE orders;
ALTER TABLE orders_new RENAME TO orders;
CREATE INDEX orders_campaign ON orders (campaign_id);
CREATE INDEX orders_member ON orders (member_id);

COMMIT;
PRAGMA foreign_keys = ON;
-- +goose StatementEnd
