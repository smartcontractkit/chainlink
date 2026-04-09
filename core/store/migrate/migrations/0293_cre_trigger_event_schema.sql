-- +goose Up

CREATE SCHEMA cre;

CREATE TABLE IF NOT EXISTS cre.trigger_pending_events (
    trigger_id    TEXT        NOT NULL,
    event_id      TEXT        NOT NULL,
    payload       BYTEA       NOT NULL,
    first_at      TIMESTAMPTZ NOT NULL,
    last_sent_at  TIMESTAMPTZ NULL,
    attempts      INTEGER     NOT NULL DEFAULT 0,
    PRIMARY KEY (trigger_id, event_id)
);

-- +goose Down

DROP TABLE IF EXISTS cre.trigger_pending_events;
DROP SCHEMA cre;
