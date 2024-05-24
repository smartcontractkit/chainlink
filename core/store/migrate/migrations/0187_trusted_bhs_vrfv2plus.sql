-- +goose Up
ALTER TABLE blockhash_store_specs
    ADD COLUMN IF NOT EXISTS "trusted_blockhash_store_address" bytea
    CHECK (octet_length(trusted_blockhash_store_address) = 20);

ALTER TABLE blockhash_store_specs
    ADD COLUMN IF NOT EXISTS "trusted_blockhash_store_batch_size" integer DEFAULT 0;
-- +goose Down
ALTER TABLE blockhash_store_specs DROP COLUMN "trusted_blockhash_store_address";
ALTER TABLE blockhash_store_specs DROP COLUMN "trusted_blockhash_store_batch_size";