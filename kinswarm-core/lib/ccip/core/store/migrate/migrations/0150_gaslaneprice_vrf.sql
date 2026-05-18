-- +goose Up
ALTER TABLE vrf_specs
    ADD COLUMN "gas_lane_price" NUMERIC(78, 0)
    CHECK (gas_lane_price IS NULL OR gas_lane_price > 0);

-- +goose Down
ALTER TABLE vrf_specs DROP COLUMN "gas_lane_price";
