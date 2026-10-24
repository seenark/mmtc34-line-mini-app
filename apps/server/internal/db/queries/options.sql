-- name: ListOptions :many
SELECT * FROM options WHERE campaign_id = sqlc.arg(campaign_id) ORDER BY sort_order, id;

-- name: CreateOption :one
INSERT INTO options (campaign_id, name, description, image, price, sizes, needs_print_name, needs_number, sort_order)
VALUES (sqlc.arg(campaign_id), sqlc.arg(name), sqlc.arg(description), sqlc.arg(image), sqlc.arg(price),
        sqlc.arg(sizes), sqlc.arg(needs_print_name), sqlc.arg(needs_number), sqlc.arg(sort_order))
RETURNING *;

-- name: UpdateOption :one
UPDATE options
SET name = sqlc.arg(name), description = sqlc.arg(description), image = sqlc.arg(image), price = sqlc.arg(price),
    sizes = sqlc.arg(sizes), needs_print_name = sqlc.arg(needs_print_name), needs_number = sqlc.arg(needs_number),
    sort_order = sqlc.arg(sort_order)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteOption :execrows
DELETE FROM options WHERE id = sqlc.arg(id);
