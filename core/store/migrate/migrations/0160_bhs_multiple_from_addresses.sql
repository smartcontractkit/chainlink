-- +goose Up
ALTER TABLE blockhash_store_specs ADD COLUMN from_addresses bytea[] DEFAULT '{}' NOT NULL ;

UPDATE blockhash_store_specs SET from_addresses = from_addresses || from_address
WHERE from_address IS NOT NULL;

ALTER TABLE blockhash_store_specs DROP COLUMN from_address;

-- +goose Down
ALTER TABLE blockhash_store_specs ADD COLUMN from_address bytea;

UPDATE blockhash_store_specs SET from_address = from_addresses[1]
WHERE array_length(from_addresses, 1) > 0;

ALTER TABLE blockhash_store_specs DROP COLUMN from_addresses;
