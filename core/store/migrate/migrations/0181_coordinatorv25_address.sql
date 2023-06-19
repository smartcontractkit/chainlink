-- +goose Up

ALTER TABLE vrf_specs ADD COLUMN coordinator_v25_address bytea;

-- add constraint to check that coordinator_v25_address is exactly 20 bytes long
ALTER TABLE vrf_specs ADD CONSTRAINT coordinator_v25_address_length CHECK (octet_length(coordinator_v25_address) = 20);

-- +goose Down

ALTER TABLE vrf_specs DROP COLUMN coordinator_v25_address;
