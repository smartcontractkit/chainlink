-- +goose Up
ALTER TABLE vrf_specs
    ADD COLUMN "request_timeout" BIGINT
    CHECK (request_timeout > 0)
    DEFAULT 24 * 60 * 60 * 1e9 -- default of one day in nanoseconds
    NOT NULL;

-- +goose Down
ALTER TABLE vrf_specs
    DROP COLUMN "request_timeout";
