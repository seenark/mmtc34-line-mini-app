-- name: ListCampaignPhotos :many
SELECT name FROM campaign_photos WHERE campaign_id = sqlc.arg(campaign_id) ORDER BY sort_order;

-- name: ListPhotosForCampaigns :many
SELECT campaign_id, name FROM campaign_photos
WHERE campaign_id IN (sqlc.slice(campaign_ids))
ORDER BY campaign_id, sort_order;

-- name: DeleteCampaignPhotos :exec
DELETE FROM campaign_photos WHERE campaign_id = sqlc.arg(campaign_id);

-- name: CreateCampaignPhoto :exec
INSERT INTO campaign_photos (campaign_id, name, sort_order)
VALUES (sqlc.arg(campaign_id), sqlc.arg(name), sqlc.arg(sort_order));
