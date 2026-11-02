-- A Campaign holds any number of Campaign Photos in a chosen order; the first (lowest sort_order)
-- is the Cover. The single campaigns.cover_image becomes the photo at sort_order 0.
-- +goose Up
CREATE TABLE campaign_photos (
    id          INTEGER PRIMARY KEY,
    campaign_id INTEGER NOT NULL REFERENCES campaigns (id) ON DELETE CASCADE,
    name        TEXT    NOT NULL, -- file name under DATA_DIR/files/images
    sort_order  INTEGER NOT NULL,
    UNIQUE (campaign_id, name)
);
CREATE INDEX campaign_photos_campaign ON campaign_photos (campaign_id, sort_order);
INSERT INTO campaign_photos (campaign_id, name, sort_order)
SELECT id, cover_image, 0 FROM campaigns WHERE cover_image <> '';
ALTER TABLE campaigns DROP COLUMN cover_image;
