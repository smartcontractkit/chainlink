-- +goose Up
-- +goose StatementBegin
CREATE TABLE solana_chains (
    id text PRIMARY KEY,
    cfg jsonb NOT NULL DEFAULT '{}',
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    enabled BOOL DEFAULT TRUE NOT NULL
);
CREATE TABLE solana_nodes (
    id serial PRIMARY KEY,
    name varchar(255) NOT NULL CHECK (name != ''),
    solana_chain_id text NOT NULL REFERENCES solana_chains (id) ON DELETE CASCADE,
    solana_url text CHECK (solana_url != ''),
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
CREATE INDEX idx_nodes_solana_chain_id ON solana_nodes (solana_chain_id);
CREATE UNIQUE INDEX idx_solana_nodes_unique_name ON solana_nodes (lower(name));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE solana_nodes;
DROP TABLE solana_chains;
-- +goose StatementEnd
