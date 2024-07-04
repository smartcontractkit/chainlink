-- +goose Up
CREATE TABLE registry_syncer_state (
    id SERIAL PRIMARY KEY,
    data JSONB NOT NULL,
    data_hash TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose Down
-- +goose StatementBegin
DROP TABLE registry_syncer_state;
-- +goose StatementEnd
