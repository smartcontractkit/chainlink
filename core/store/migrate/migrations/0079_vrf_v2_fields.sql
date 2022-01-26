-- +goose Up
ALTER TABLE vrf_specs
    ADD COLUMN from_address bytea,
    ADD COLUMN poll_period bigint NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE vrf_specs
    DROP COLUMN from_address,
    DROP COLUMN poll_period;
