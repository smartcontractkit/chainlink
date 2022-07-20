-- +goose Up
ALTER TABLE vrf_specs
    ADD COLUMN "backoff_initial_delay" BIGINT
    CHECK (backoff_initial_delay >= 0)
    DEFAULT 0
    NOT NULL;

ALTER TABLE vrf_specs
    ADD COLUMN "backoff_max_delay" BIGINT
    CHECK (backoff_max_delay >= 0)
    DEFAULT 0
    NOT NULL;

-- +goose Down
ALTER TABLE vrf_specs DROP COLUMN "backoff_initial_delay";

ALTER TABLE vrf_specs DROP COLUMN "backoff_max_delay";
