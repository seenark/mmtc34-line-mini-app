-- Slip writes and the Slip ownership lookup. Payments and their Slip images are read for display
-- only through orders.sql; rejecting a pending Payment is orders.sql RejectPendingPayment, shared
-- with Admin cancel.

-- name: InsertPendingPayment :one
INSERT INTO payments (order_id, status, created_at)
VALUES (sqlc.arg(order_id), 'pending', sqlc.arg(now))
RETURNING id;

-- name: InsertSlipImage :exec
INSERT INTO payment_slips (payment_id, name, created_at)
VALUES (sqlc.arg(payment_id), sqlc.arg(name), sqlc.arg(now));

-- name: CountSlipImages :one
SELECT count(*) FROM payment_slips WHERE payment_id = sqlc.arg(payment_id);

-- name: DeleteSlipImage :execrows
DELETE FROM payment_slips WHERE payment_id = sqlc.arg(payment_id) AND name = sqlc.arg(name);

-- name: ConfirmPendingPayment :exec
UPDATE payments
SET status = 'confirmed', amount = sqlc.arg(amount), reviewed_by = sqlc.arg(reviewed_by),
    reviewed_at = sqlc.arg(now)
WHERE order_id = sqlc.arg(order_id) AND status = 'pending';

-- name: SlipOwner :one
-- The Member whose Order the Slip image belongs to; no row = no such Slip.
SELECT o.member_id
FROM payment_slips s
JOIN payments p ON p.id = s.payment_id
JOIN orders o ON o.id = p.order_id
WHERE s.name = sqlc.arg(name);
