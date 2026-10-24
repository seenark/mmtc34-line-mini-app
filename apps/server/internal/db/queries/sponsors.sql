-- name: ListSponsors :many
SELECT sqlc.embed(sp), COALESCE(m.nickname, '') AS member_nickname
FROM sponsor_positions sp
LEFT JOIN members m ON m.id = sp.member_id
WHERE sp.campaign_id = sqlc.arg(campaign_id)
ORDER BY sp.sort_order, sp.id;

-- name: CreateSponsor :one
INSERT INTO sponsor_positions (campaign_id, name, sponsor_name, logo, member_id, sort_order)
VALUES (sqlc.arg(campaign_id), sqlc.arg(name), sqlc.arg(sponsor_name), sqlc.arg(logo), sqlc.arg(member_id), sqlc.arg(sort_order))
RETURNING *;

-- name: UpdateSponsor :one
UPDATE sponsor_positions
SET name = sqlc.arg(name), sponsor_name = sqlc.arg(sponsor_name), logo = sqlc.arg(logo),
    member_id = sqlc.arg(member_id), sort_order = sqlc.arg(sort_order)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteSponsor :execrows
DELETE FROM sponsor_positions WHERE id = sqlc.arg(id);
