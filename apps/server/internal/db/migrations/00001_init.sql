-- Complete milestone-1 schema.
-- A later schema change adds a NEW file named with a goose timestamp version
-- (e.g. 20261101093000_add_x.sql) so two parallel branches never pick the same number.
--
-- Conventions:
--   * Money: whole baht INTEGER.
--   * Time: TEXT, RFC3339 UTC ("2026-10-24T09:00:00Z"), always written from the injected clock.
--     There is deliberately NO `DEFAULT CURRENT_TIMESTAMP` anywhere: tests move the clock.
--   * Booleans: INTEGER 0/1 with CHECK.
--   * Effective statuses are NOT stored. `orders.status` and `campaigns.status` are the stored values only.

-- +goose Up

-- Exactly one Group exists (spec: "schema keeps one group row only").
CREATE TABLE class_groups (
    id   INTEGER PRIMARY KEY CHECK (id = 1),
    name TEXT    NOT NULL
);
INSERT INTO class_groups (id, name) VALUES (1, 'MMTC34');

CREATE TABLE members (
    id           INTEGER PRIMARY KEY,
    -- NULL = Unlinked Member (ADR 0002) or deleted Member. The session cookie holds this value,
    -- so moving it (link) or clearing it (delete) re-points / ends every session at once.
    line_user_id TEXT    UNIQUE,
    status       TEXT    NOT NULL CHECK (status IN ('applicant', 'approved')),
    role         TEXT    NOT NULL DEFAULT 'member' CHECK (role IN ('member', 'admin')),
    cohort       TEXT    NOT NULL DEFAULT '',
    nickname     TEXT    NOT NULL CHECK (nickname <> ''),
    real_name    TEXT    NOT NULL DEFAULT '',
    address      TEXT    NOT NULL DEFAULT '',
    phone        TEXT    NOT NULL DEFAULT '',
    deleted_at   TEXT,
    created_at   TEXT    NOT NULL,
    -- A deleted Member keeps only its nickname (PDPA).
    CHECK (deleted_at IS NULL OR (line_user_id IS NULL AND real_name = '' AND cohort = ''
                                  AND address = '' AND phone = '')),
    -- An Applicant always came through LINE; Unlinked Members are created approved.
    CHECK (status = 'approved' OR line_user_id IS NOT NULL)
);

CREATE TABLE campaigns (
    id                INTEGER PRIMARY KEY,
    group_id          INTEGER NOT NULL DEFAULT 1 REFERENCES class_groups (id),
    type              TEXT    NOT NULL CHECK (type IN ('merch', 'donation', 'event')),
    title             TEXT    NOT NULL CHECK (title <> ''),
    description       TEXT    NOT NULL DEFAULT '',
    cover_image       TEXT    NOT NULL DEFAULT '',
    bank_name         TEXT    NOT NULL DEFAULT '',
    bank_account_no   TEXT    NOT NULL DEFAULT '',
    bank_account_name TEXT    NOT NULL DEFAULT '',
    qr_image          TEXT    NOT NULL DEFAULT '',
    status            TEXT    NOT NULL CHECK (status IN ('draft', 'open', 'closed')),
    close_at          TEXT,
    created_at        TEXT    NOT NULL
);

CREATE TABLE options (
    id               INTEGER PRIMARY KEY,
    campaign_id      INTEGER NOT NULL REFERENCES campaigns (id),
    name             TEXT    NOT NULL CHECK (name <> ''),
    description      TEXT    NOT NULL DEFAULT '',
    image            TEXT    NOT NULL DEFAULT '',
    price            INTEGER CHECK (price IS NULL OR price >= 0), -- NULL = Member chooses the amount
    sizes            TEXT    NOT NULL DEFAULT '[]' CHECK (json_valid(sizes) AND json_type(sizes) = 'array'),
    needs_print_name INTEGER NOT NULL DEFAULT 0 CHECK (needs_print_name IN (0, 1)),
    needs_number     INTEGER NOT NULL DEFAULT 0 CHECK (needs_number IN (0, 1)),
    sort_order       INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX options_campaign ON options (campaign_id, sort_order);

CREATE TABLE sponsor_positions (
    id           INTEGER PRIMARY KEY,
    campaign_id  INTEGER NOT NULL REFERENCES campaigns (id),
    name         TEXT    NOT NULL CHECK (name <> ''),
    sponsor_name TEXT    NOT NULL DEFAULT '',
    logo         TEXT    NOT NULL DEFAULT '',
    member_id    INTEGER REFERENCES members (id) ON DELETE SET NULL,
    sort_order   INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX sponsor_positions_campaign ON sponsor_positions (campaign_id, sort_order);

CREATE TABLE orders (
    id           INTEGER PRIMARY KEY,
    campaign_id  INTEGER NOT NULL REFERENCES campaigns (id),
    member_id    INTEGER NOT NULL REFERENCES members (id),
    -- 1 = an Admin placed the Order on the Member's behalf. Such Orders are exempt from
    -- "awaiting_payment in a closed Campaign reads cancelled".
    placed_by_admin INTEGER NOT NULL DEFAULT 0 CHECK (placed_by_admin IN (0, 1)),
    status       TEXT    NOT NULL CHECK (status IN ('awaiting_payment', 'pending_review', 'paid', 'cancelled')),
    delivery     TEXT    NOT NULL CHECK (delivery IN ('ship', 'pickup')),
    ship_name    TEXT    NOT NULL DEFAULT '',
    ship_address TEXT    NOT NULL DEFAULT '',
    ship_phone   TEXT    NOT NULL DEFAULT '',
    created_at   TEXT    NOT NULL,
    updated_at   TEXT    NOT NULL
);
CREATE INDEX orders_campaign ON orders (campaign_id);
CREATE INDEX orders_member ON orders (member_id);

CREATE TABLE items (
    id          INTEGER PRIMARY KEY,
    order_id    INTEGER NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    -- RESTRICT enforces that an Option used by an Item cannot be deleted.
    option_id   INTEGER NOT NULL REFERENCES options (id) ON DELETE RESTRICT,
    size        TEXT    NOT NULL DEFAULT '',
    print_name  TEXT    NOT NULL DEFAULT '',
    number      TEXT    NOT NULL DEFAULT '',
    quantity    INTEGER NOT NULL CHECK (quantity BETWEEN 1 AND 99),
    unit_price  INTEGER NOT NULL CHECK (unit_price >= 0) -- captured at checkout, never re-read from options
);
CREATE INDEX items_order ON items (order_id);

CREATE TABLE payments (
    id            INTEGER PRIMARY KEY,
    order_id      INTEGER NOT NULL REFERENCES orders (id),
    slip          TEXT,                     -- file name under DATA_DIR/files/slips; NULL = cash
    status        TEXT    NOT NULL CHECK (status IN ('pending', 'confirmed', 'rejected')),
    amount        INTEGER,                  -- set exactly when confirmed
    reject_reason TEXT    NOT NULL DEFAULT '',
    created_at    TEXT    NOT NULL,
    reviewed_by   INTEGER REFERENCES members (id),
    reviewed_at   TEXT,
    CHECK ((status = 'confirmed') = (amount IS NOT NULL)),
    CHECK (amount IS NULL OR amount > 0),
    CHECK (status <> 'rejected' OR reject_reason <> ''),
    CHECK (slip IS NOT NULL OR status = 'confirmed'),          -- cash is born confirmed
    CHECK ((status = 'pending') = (reviewed_by IS NULL))
);
CREATE INDEX payments_order ON payments (order_id);
CREATE UNIQUE INDEX payments_slip ON payments (slip) WHERE slip IS NOT NULL;
-- At most one Slip under review per Order: a double-tapped upload cannot create two.
CREATE UNIQUE INDEX payments_one_pending ON payments (order_id) WHERE status = 'pending';

CREATE TABLE adjustments (
    id         INTEGER PRIMARY KEY,
    order_id   INTEGER NOT NULL REFERENCES orders (id),
    amount     INTEGER NOT NULL CHECK (amount <> 0), -- + Member owes more, - discount/refund
    note       TEXT    NOT NULL CHECK (note <> ''),
    created_by INTEGER NOT NULL REFERENCES members (id),
    created_at TEXT    NOT NULL
);
CREATE INDEX adjustments_order ON adjustments (order_id);

-- +goose Down
DROP TABLE adjustments;
DROP TABLE payments;
DROP TABLE items;
DROP TABLE orders;
DROP TABLE sponsor_positions;
DROP TABLE options;
DROP TABLE campaigns;
DROP TABLE members;
DROP TABLE class_groups;
