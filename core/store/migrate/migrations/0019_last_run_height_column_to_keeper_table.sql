-- +goose Up
ALTER TABLE upkeep_registrations ADD COLUMN last_run_block_height BIGINT NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE upkeep_registrations DROP COLUMN last_run_block_height;
