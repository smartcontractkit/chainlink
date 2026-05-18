-- +goose Up
CREATE TABLE IF NOT EXISTS lease_lock (
    client_id uuid NOT NULL,
    expires_at timestamptz NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS only_one_lease_lock ON lease_lock ((client_id IS NOT NULL));

-- +goose Down
DROP TABLE lease_lock;
