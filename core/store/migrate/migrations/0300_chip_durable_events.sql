-- +goose Up

CREATE TABLE IF NOT EXISTS cre.chip_durable_events (
    id           BIGSERIAL   PRIMARY KEY,
    payload      BYTEA       NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    delivered_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_chip_durable_events_created_at
    ON cre.chip_durable_events (created_at ASC);

CREATE INDEX IF NOT EXISTS idx_chip_durable_events_pending_delivery
    ON cre.chip_durable_events (created_at ASC)
    WHERE delivered_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_chip_durable_events_delivered_purge
    ON cre.chip_durable_events (delivered_at ASC)
    WHERE delivered_at IS NOT NULL;

-- +goose Down

DROP INDEX IF EXISTS cre.idx_chip_durable_events_delivered_purge;
DROP INDEX IF EXISTS cre.idx_chip_durable_events_pending_delivery;
DROP INDEX IF EXISTS cre.idx_chip_durable_events_created_at;
DROP TABLE IF EXISTS cre.chip_durable_events;
