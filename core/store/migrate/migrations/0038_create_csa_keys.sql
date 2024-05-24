-- +goose Up
CREATE TABLE csa_keys(
    id BIGSERIAL PRIMARY KEY,
    public_key bytea NOT NULL CHECK (octet_length(public_key) = 32) UNIQUE,
    encrypted_private_key jsonb NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
-- +goose Down
DROP TABLE csa_keys;
