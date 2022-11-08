-- +goose Up
ALTER TABLE jobs ADD COLUMN gas_limit BIGINT DEFAULT NULL;
-- +goose Down
ALTER TABLE jobs DROP COLUMN gas_limit;
