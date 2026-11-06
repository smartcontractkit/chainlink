-- +goose Up
ALTER TABLE vrf_specs
    ADD COLUMN "gas_overhead" INTEGER
    CHECK (gas_overhead IS NULL OR gas_overhead > 0);

-- +goose Down
ALTER TABLE vrf_specs DROP COLUMN "gas_overhead";
