-- A Slip Payment holds one or more images (payment_slips); a Payment with none is cash.
-- payments.slip sits in a CHECK and a unique index, so SQLite cannot drop it in place:
-- payments is rebuilt. payment_slips is created only after the old table is dropped, because
-- dropping a parent with foreign keys on would cascade-delete its rows.
-- +goose Up
CREATE TABLE payments_new (
    id            INTEGER PRIMARY KEY,
    order_id      INTEGER NOT NULL REFERENCES orders (id),
    status        TEXT    NOT NULL CHECK (status IN ('pending', 'confirmed', 'rejected')),
    amount        INTEGER,                  -- set exactly when confirmed
    reject_reason TEXT    NOT NULL DEFAULT '',
    created_at    TEXT    NOT NULL,
    reviewed_by   INTEGER REFERENCES members (id),
    reviewed_at   TEXT,
    CHECK ((status = 'confirmed') = (amount IS NOT NULL)),
    CHECK (amount IS NULL OR amount > 0),
    CHECK (status <> 'rejected' OR reject_reason <> ''),
    CHECK ((status = 'pending') = (reviewed_by IS NULL))
);
INSERT INTO payments_new (id, order_id, status, amount, reject_reason, created_at, reviewed_by, reviewed_at)
SELECT id, order_id, status, amount, reject_reason, created_at, reviewed_by, reviewed_at FROM payments;
CREATE TABLE old_slips AS SELECT id AS payment_id, slip AS name, created_at FROM payments WHERE slip IS NOT NULL;

DROP TABLE payments;
ALTER TABLE payments_new RENAME TO payments;
CREATE INDEX payments_order ON payments (order_id);
-- At most one Slip under review per Order: a double-tapped upload cannot create two.
CREATE UNIQUE INDEX payments_one_pending ON payments (order_id) WHERE status = 'pending';

CREATE TABLE payment_slips (
    id         INTEGER PRIMARY KEY,
    payment_id INTEGER NOT NULL REFERENCES payments (id) ON DELETE CASCADE,
    name       TEXT    NOT NULL UNIQUE, -- file name under DATA_DIR/files/slips
    created_at TEXT    NOT NULL
);
CREATE INDEX payment_slips_payment ON payment_slips (payment_id, id);
INSERT INTO payment_slips (payment_id, name, created_at) SELECT payment_id, name, created_at FROM old_slips;
DROP TABLE old_slips;
