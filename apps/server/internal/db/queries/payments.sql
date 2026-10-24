-- Slip writes and the Slip ownership lookup. Payments are read for display only through
-- orders.sql (ListPaymentsForOrders); rejecting a pending Payment is orders.sql
-- RejectPendingPayment, shared with Admin cancel.

-- name: InsertPendingPayment :exec
INSERT INTO payments (order_id, slip, status, created_at)
VALUES (sqlc.arg(order_id), sqlc.arg(slip), 'pending', sqlc.arg(now));

-- name: ConfirmPendingPayment :exec
UPDATE payments
SET status = 'confirmed', amount = sqlc.arg(amount), reviewed_by = sqlc.arg(reviewed_by),
    reviewed_at = sqlc.arg(now)
WHERE order_id = sqlc.arg(order_id) AND status = 'pending';

-- name: SlipOwner :one
-- The Member whose Order the Slip belongs to; no row = no such Slip.
SELECT o.member_id
FROM payments p
JOIN orders o ON o.id = p.order_id
WHERE p.slip = sqlc.arg(slip);
