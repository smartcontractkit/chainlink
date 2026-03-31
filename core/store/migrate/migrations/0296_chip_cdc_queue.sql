-- +goose Up

CREATE TABLE IF NOT EXISTS cre.chip_queue (
    id BIGSERIAL PRIMARY KEY,
    payload BYTEA NOT NULL,
    attributes JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down

DROP TABLE IF EXISTS cre.chip_queue;
