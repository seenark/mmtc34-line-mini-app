-- Owned by ticket .6: cash Payments and Adjustments. Reads go through orders.sql.

-- name: InsertCashPayment :exec
-- Cash is confirmed on entry and carries no Slip (slip NULL).
INSERT INTO payments (order_id, status, amount, created_at, reviewed_by, reviewed_at)
VALUES (sqlc.arg(order_id), 'confirmed', sqlc.arg(amount), sqlc.arg(now), sqlc.arg(reviewed_by), sqlc.arg(now));

-- name: InsertAdjustment :exec
INSERT INTO adjustments (order_id, amount, note, created_by, created_at)
VALUES (sqlc.arg(order_id), sqlc.arg(amount), sqlc.arg(note), sqlc.arg(created_by), sqlc.arg(now));
