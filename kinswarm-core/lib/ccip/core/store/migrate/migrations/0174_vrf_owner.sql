-- +goose Up

ALTER TABLE vrf_specs
    ADD COLUMN IF NOT EXISTS "vrf_owner_address" bytea
    CHECK (octet_length(vrf_owner_address) = 20);

-- +goose Down
ALTER TABLE vrf_specs DROP COLUMN "vrf_owner_address";
