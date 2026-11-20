-- +goose Up
ALTER TABLE vrf_specs
    ADD COLUMN "gas_overhead" INTEGER
    CHECK (gas_overhead IS NULL OR gas_overhead > 0);

ALTER TABLE vrf_specs
    ADD COLUMN "premium_percentage" INTEGER
    CHECK (premium_percentage IS NULL OR (premium_percentage > 0 AND premium_percentage <= 155));

-- +goose Down
ALTER TABLE vrf_specs DROP COLUMN "premium_percentage";
ALTER TABLE vrf_specs DROP COLUMN "gas_overhead";
