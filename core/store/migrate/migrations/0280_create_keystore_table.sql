-- +goose Up
CREATE TABLE IF NOT EXISTS encrypted_keystore (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    encrypted_data BYTEA NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS encrypted_keystore;