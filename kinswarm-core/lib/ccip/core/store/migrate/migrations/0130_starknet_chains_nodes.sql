-- +goose Up
-- +goose StatementBegin
CREATE TABLE starknet_chains (
                               id text PRIMARY KEY,
                               cfg jsonb NOT NULL DEFAULT '{}',
                               created_at timestamptz NOT NULL,
                               updated_at timestamptz NOT NULL,
                               enabled BOOL DEFAULT TRUE NOT NULL
);
CREATE TABLE starknet_nodes (
                              id serial PRIMARY KEY,
                              name varchar(255) NOT NULL CHECK (name != ''),
    chain_id text NOT NULL REFERENCES starknet_chains (id) ON DELETE CASCADE,
    url text CHECK (url != ''),
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL
);
CREATE INDEX idx_starknet_nodes_chain_id ON starknet_nodes (chain_id);
CREATE UNIQUE INDEX idx_starknet_nodes_unique_name ON starknet_nodes (lower(name));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE starknet_nodes;
DROP TABLE starknet_chains;
-- +goose StatementEnd
