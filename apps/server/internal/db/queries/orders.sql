-- Every Order read in the app goes through the four reads below (server.loadOrders);
-- no other query file selects from orders/items/payments/adjustments for display.
-- Writes for slips, cash/adjustments, admin edit and PDPA erase live in their own query files.

-- name: ListOrders :many
-- Filters are optional (NULL = any). The Campaign row comes along so domain.Order can compute
-- effective status without a second lookup. Newest first.
SELECT sqlc.embed(o),
       sqlc.embed(c),
       m.nickname  AS member_nickname,
       m.real_name AS member_real_name
FROM orders o
JOIN members   m ON m.id = o.member_id
JOIN campaigns c ON c.id = o.campaign_id
WHERE (CAST(sqlc.narg(order_id) AS INTEGER) IS NULL OR o.id = CAST(sqlc.narg(order_id) AS INTEGER))
  AND (CAST(sqlc.narg(campaign_id) AS INTEGER) IS NULL OR o.campaign_id = CAST(sqlc.narg(campaign_id) AS INTEGER))
  AND (CAST(sqlc.narg(member_id) AS INTEGER) IS NULL OR o.member_id = CAST(sqlc.narg(member_id) AS INTEGER))
ORDER BY o.id DESC;

-- name: ListItemsForOrders :many
SELECT i.*, op.name AS option_name
FROM items i
JOIN options op ON op.id = i.option_id
WHERE i.order_id IN (sqlc.slice(order_ids))
ORDER BY i.order_id, i.id;

-- name: ListPaymentsForOrders :many
SELECT p.*, r.nickname AS reviewer_nickname
FROM payments p
LEFT JOIN members r ON r.id = p.reviewed_by
WHERE p.order_id IN (sqlc.slice(order_ids))
ORDER BY p.order_id, p.id;

-- name: ListAdjustmentsForOrders :many
SELECT a.*, m.nickname AS creator_nickname
FROM adjustments a
JOIN members m ON m.id = a.created_by
WHERE a.order_id IN (sqlc.slice(order_ids))
ORDER BY a.order_id, a.id;

-- name: InsertOrder :one
INSERT INTO orders (campaign_id, member_id, placed_by_admin, status, delivery,
                    ship_name, ship_address, ship_phone, created_at, updated_at)
VALUES (sqlc.arg(campaign_id), sqlc.arg(member_id), sqlc.arg(placed_by_admin), 'awaiting_payment',
        sqlc.arg(delivery), sqlc.arg(ship_name), sqlc.arg(ship_address), sqlc.arg(ship_phone),
        sqlc.arg(now), sqlc.arg(now))
RETURNING id;

-- name: InsertItem :exec
INSERT INTO items (order_id, option_id, size, print_name, number, quantity, unit_price)
VALUES (sqlc.arg(order_id), sqlc.arg(option_id), sqlc.arg(size), sqlc.arg(print_name),
        sqlc.arg(number), sqlc.arg(quantity), sqlc.arg(unit_price));

-- name: SetOrderStatus :execrows
-- Compare-and-set: the only statement that changes orders.status. 0 rows = someone else
-- moved the Order first; server.apply turns that into 409, so a replayed or raced
-- transition applies at most once.
UPDATE orders
SET status = sqlc.arg(to_status), updated_at = sqlc.arg(now)
WHERE id = sqlc.arg(order_id) AND status = sqlc.arg(from_status);

-- name: RejectPendingPayment :exec
-- No-op when the Order has no Slip under review.
UPDATE payments
SET status = 'rejected', reject_reason = sqlc.arg(reason), reviewed_by = sqlc.arg(reviewed_by),
    reviewed_at = sqlc.arg(now)
WHERE order_id = sqlc.arg(order_id) AND status = 'pending';
