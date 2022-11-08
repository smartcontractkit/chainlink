-- +goose Up
ALTER TABLE vrf_specs
    ADD COLUMN "gas_lane_price_gwei" BIGINT
    CHECK (gas_lane_price_gwei >= 0);

-- +goose Down
ALTER TABLE vrf_specs DROP COLUMN "gas_lane_price_gwei";
