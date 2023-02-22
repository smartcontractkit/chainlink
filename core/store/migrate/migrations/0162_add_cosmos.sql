-- +goose Up

-- Table Definition ----------------------------------------------

CREATE TABLE cosmos_chains (
    id text PRIMARY KEY,
    cfg jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    enabled boolean NOT NULL DEFAULT true
);

-- Table Definition ----------------------------------------------

CREATE TABLE cosmos_msgs (
    id BIGSERIAL PRIMARY KEY,
    cosmos_chain_id text NOT NULL REFERENCES cosmos_chains(id) ON DELETE CASCADE,
    contract_id text NOT NULL,
    raw bytea NOT NULL,
    state text NOT NULL,
    tx_hash text,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    type text NOT NULL DEFAULT '/cosmwasm.wasm.v1.MsgExecuteContract'::text,
    CONSTRAINT cosmos_msgs_check CHECK (tx_hash <> NULL::text OR state <> 'broadcasted'::text AND state <> 'confirmed'::text)
);

-- Indices -------------------------------------------------------

CREATE INDEX idx_cosmos_msgs_cosmos_chain_id_state_contract_id ON cosmos_msgs(cosmos_chain_id text_ops,state text_ops,contract_id text_ops);

-- Table Definition ----------------------------------------------

CREATE TABLE cosmos_nodes (
    id SERIAL PRIMARY KEY,
    name character varying(255) NOT NULL CHECK (name::text <> ''::text),
    cosmos_chain_id text NOT NULL REFERENCES cosmos_chains(id) ON DELETE CASCADE,
    tendermint_url text CHECK (tendermint_url <> ''::text),
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

-- Indices -------------------------------------------------------

CREATE INDEX idx_nodes_cosmos_chain_id ON cosmos_nodes(cosmos_chain_id text_ops);
CREATE UNIQUE INDEX idx_cosmos_nodes_unique_name ON cosmos_nodes((lower(name::text)) text_ops);

-- +goose Down

DROP TABLE cosmos_msgs;
DROP TABLE cosmos_nodes;
DROP TABLE cosmos_chains;
