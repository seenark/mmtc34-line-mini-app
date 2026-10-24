-- Admin edit replaces an Order's Items and delivery fields.

-- name: DeleteItemsForOrder :exec
DELETE FROM items WHERE order_id = sqlc.arg(order_id);

-- name: UpdateOrderContents :exec
UPDATE orders
SET delivery = sqlc.arg(delivery), ship_name = sqlc.arg(ship_name), ship_address = sqlc.arg(ship_address),
    ship_phone = sqlc.arg(ship_phone), updated_at = sqlc.arg(now)
WHERE id = sqlc.arg(order_id);
