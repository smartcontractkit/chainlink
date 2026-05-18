-- +goose Up
ALTER TABLE feeds_managers
ADD COLUMN is_ocr_bootstrap_peer boolean NOT NULL DEFAULT false;
-- +goose Down
ALTER TABLE feeds_managers
DROP COLUMN is_ocr_bootstrap_peer;
