-- +goose Up

CREATE TABLE IF NOT EXISTS cre.chip_durable_events (
    id         BIGSERIAL   PRIMARY KEY,
    payload    BYTEA       NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_chip_durable_events_created_at
    ON cre.chip_durable_events (created_at ASC);

-- +goose Down

DROP INDEX IF EXISTS cre.idx_chip_durable_events_created_at;
DROP TABLE IF EXISTS cre.chip_durable_events;
