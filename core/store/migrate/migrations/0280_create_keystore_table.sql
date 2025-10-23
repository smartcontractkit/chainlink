-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS encrypted_keystore (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    encrypted_data BYTEA NOT NULL
);

CREATE OR REPLACE FUNCTION encrypted_keystore_set_updated_at() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER encrypted_keystore_set_updated_at_trigger
    BEFORE UPDATE ON encrypted_keystore
    FOR EACH ROW
EXECUTE FUNCTION encrypted_keystore_set_updated_at();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS encrypted_keystore_set_updated_at_trigger on encrypted_keystore;
DROP FUNCTION IF EXISTS encrypted_keystore_set_updated_at;
DROP TABLE IF EXISTS encrypted_keystore;
-- +goose StatementEnd