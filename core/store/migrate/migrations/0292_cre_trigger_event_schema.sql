-- +goose Up

CREATE TABLE IF NOT EXISTS trigger_pending_events (
    trigger_id    TEXT        NOT NULL,
    event_id      TEXT        NOT NULL,
    payload       BYTEA       NOT NULL,
    first_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_sent_at  TIMESTAMPTZ NULL,
    attempts      INTEGER     NOT NULL DEFAULT 0,
    PRIMARY KEY (trigger_id, event_id)
);

-- Efficient scans for retransmit scheduler and trigger cleanup
CREATE INDEX IF NOT EXISTS idx_trigger_pending_events_trigger_id
    ON trigger_pending_events (trigger_id);

CREATE INDEX IF NOT EXISTS idx_trigger_pending_events_last_sent_at
    ON trigger_pending_events (last_sent_at);

-- +goose Down

DROP TABLE IF EXISTS trigger_pending_events;

