-- name: ListCampaigns :many
SELECT * FROM campaigns ORDER BY created_at DESC, id DESC;

-- name: ListPublishedCampaigns :many
SELECT * FROM campaigns WHERE status <> 'draft' ORDER BY created_at DESC, id DESC;

-- name: GetCampaign :one
SELECT * FROM campaigns WHERE id = sqlc.arg(id);

-- name: CreateCampaign :one
INSERT INTO campaigns (type, title, description, cover_image, bank_name, bank_account_no, bank_account_name,
                       qr_image, status, close_at, created_at)
VALUES (sqlc.arg(type), sqlc.arg(title), sqlc.arg(description), sqlc.arg(cover_image), sqlc.arg(bank_name),
        sqlc.arg(bank_account_no), sqlc.arg(bank_account_name), sqlc.arg(qr_image), sqlc.arg(status),
        sqlc.arg(close_at), sqlc.arg(created_at))
RETURNING *;

-- name: UpdateCampaign :one
UPDATE campaigns
SET type = sqlc.arg(type), title = sqlc.arg(title), description = sqlc.arg(description),
    cover_image = sqlc.arg(cover_image), bank_name = sqlc.arg(bank_name),
    bank_account_no = sqlc.arg(bank_account_no), bank_account_name = sqlc.arg(bank_account_name),
    qr_image = sqlc.arg(qr_image), status = sqlc.arg(status), close_at = sqlc.arg(close_at)
WHERE id = sqlc.arg(id)
RETURNING *;
