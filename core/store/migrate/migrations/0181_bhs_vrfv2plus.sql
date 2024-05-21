-- +goose Up
ALTER TABLE blockhash_store_specs
    ADD COLUMN IF NOT EXISTS "coordinator_v2_plus_address" bytea
    CHECK (octet_length(coordinator_v2_plus_address) = 20);

ALTER TABLE block_header_feeder_specs
    ADD COLUMN IF NOT EXISTS "coordinator_v2_plus_address" bytea
    CHECK (octet_length(coordinator_v2_plus_address) = 20);

-- +goose Down
ALTER TABLE blockhash_store_specs DROP COLUMN "coordinator_v2_plus_address";
ALTER TABLE block_header_feeder_specs DROP COLUMN "coordinator_v2_plus_address";