-- +goose Up
-- This is a storage table for the keystore https://github.com/smartcontractkit/chainlink-common/tree/main/keystore entries.
-- Support for versioning (and other automatic handling) can be conveniently done using triggers at the database level.
-- Please do not implement such logic at the go library level.
CREATE TABLE IF NOT EXISTS encrypted_keystore (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    encrypted_data BYTEA NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS encrypted_keystore;