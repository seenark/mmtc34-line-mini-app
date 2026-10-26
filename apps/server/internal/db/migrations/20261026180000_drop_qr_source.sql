-- Generated Thai QR is gone: no banking app accepted the bank-account payload, so a Campaign
-- only shows its uploaded qr_image.
-- +goose Up
ALTER TABLE campaigns DROP COLUMN qr_source;
