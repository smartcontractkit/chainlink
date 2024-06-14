-- +goose Up
ALTER TABLE blockhash_store_specs ADD COLUMN heartbeat_period bigint DEFAULT 0 NOT NULL;

-- +goose Down
ALTER TABLE blockhash_store_specs DROP COLUMN heartbeat_period;
