-- +goose Up
ALTER TABLE vrf_specs ADD COLUMN from_addresses bytea[] DEFAULT '{}' NOT NULL ;

UPDATE vrf_specs SET from_addresses = from_addresses || from_address
WHERE from_address IS NOT NULL;

ALTER TABLE vrf_specs DROP COLUMN from_address;

-- +goose Down
ALTER TABLE vrf_specs ADD COLUMN from_address bytea;

UPDATE vrf_specs SET from_address = from_addresses[1]
WHERE array_length(from_addresses, 1) > 0;

ALTER TABLE vrf_specs DROP COLUMN from_addresses;
