-- +goose Up
ALTER TABLE heads ADD COLUMN l1_block_number bigint;
-- +goose Down
ALTER TABLE heads DROP COLUMN l1_block_number;
