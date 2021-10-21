-- +goose Up
ALTER TABLE keys DROP COLUMN last_used;
-- +goose Down
ALTER TABLE keys ADD COLUMN last_used timestamptz;
