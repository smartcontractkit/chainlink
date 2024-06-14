-- +goose Up
ALTER TABLE vrf_specs
    ADD COLUMN "max_gas_price_gwei" BIGINT
    CHECK (max_gas_price_gwei >= 0);

-- +goose Down
ALTER TABLE vrf_specs DROP COLUMN "max_gas_price_gwei";
