-- Owned by ticket .4. Every Order read in the app goes through these four reads
-- (server.loadOrders); no other query file selects from orders/items/payments/adjustments
-- for display. Writes for slips (.5), cash/adjustments (.6), admin edit (.7) and
-- PDPA erase (.11) live in their own query files.

-- name: ListOrders :many
-- All filters optional; NULL means "any". Campaign stored status + close_at come along so
-- domain.Order can compute effective status without a second lookup.
SELECT sqlc.embed(o),
       m.nickname   AS member_nickname,
       m.real_name  AS member_real_name,
       c.title      AS campaign_title,
       c.status     AS campaign_status,
       c.close_at   AS campaign_close_at
FROM orders o
JOIN members   m ON m.id = o.member_id
JOIN campaigns c ON c.id = o.campaign_id
WHERE (sqlc.narg(order_id)    IS NULL OR o.id          = sqlc.narg(order_id))
  AND (sqlc.narg(campaign_id) IS NULL OR o.campaign_id = sqlc.narg(campaign_id))
  AND (sqlc.narg(member_id)   IS NULL OR o.member_id   = sqlc.narg(member_id))
ORDER BY o.id;

-- name: ListItemsForOrders :many
SELECT i.*, op.name AS option_name
FROM items i
JOIN options op ON op.id = i.option_id
WHERE i.order_id IN (sqlc.slice(order_ids))
ORDER BY i.order_id, i.id;

-- name: ListPaymentsForOrders :many
SELECT p.*, COALESCE(r.nickname, '') AS reviewer_nickname
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

-- name: ListOptionsForCheckout :many
SELECT id, campaign_id, name, price, sizes, needs_print_name, needs_number
FROM options
WHERE campaign_id = ?;

-- name: InsertOrder :one
INSERT INTO orders (campaign_id, member_id, placed_by, status, delivery,
                    ship_name, ship_address, ship_phone, created_at, updated_at)
VALUES (@campaign_id, @member_id, @placed_by, 'awaiting_payment', @delivery,
        @ship_name, @ship_address, @ship_phone, @now, @now)
RETURNING id;

-- name: InsertItem :exec
INSERT INTO items (order_id, option_id, size, print_name, number, quantity, unit_price)
VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: SetOrderStatus :execrows
-- Compare-and-set: the only statement that changes orders.status. 0 rows = someone else
-- moved the Order first; server.apply turns that into 409, so a replayed or raced
-- transition can never apply twice.
UPDATE orders
SET status = @to, updated_at = @now
WHERE id = @id AND status = @from;
