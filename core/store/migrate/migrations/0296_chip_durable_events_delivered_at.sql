-- +goose Up

ALTER TABLE cre.chip_durable_events
    ADD COLUMN IF NOT EXISTS delivered_at TIMESTAMPTZ NULL;

CREATE INDEX IF NOT EXISTS idx_chip_durable_events_pending_delivery
    ON cre.chip_durable_events (created_at ASC)
    WHERE delivered_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_chip_durable_events_delivered_purge
    ON cre.chip_durable_events (delivered_at ASC)
    WHERE delivered_at IS NOT NULL;

-- +goose Down

DROP INDEX IF EXISTS cre.idx_chip_durable_events_delivered_purge;
DROP INDEX IF EXISTS cre.idx_chip_durable_events_pending_delivery;

ALTER TABLE cre.chip_durable_events DROP COLUMN IF EXISTS delivered_at;
