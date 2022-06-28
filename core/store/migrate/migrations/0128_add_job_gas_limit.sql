-- +goose Up
ALTER TABLE jobs ADD COLUMN gas_limit INTEGER DEFAULT NULL;
-- +goose Down
ALTER TABLE jobs DROP COLUMN gas_limit;
