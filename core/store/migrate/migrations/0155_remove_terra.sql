-- +goose Up

DROP TABLE terra_msgs;
DROP TABLE terra_nodes;
DROP TABLE terra_chains;

-- +goose Down

-- Table Definition ----------------------------------------------

CREATE TABLE terra_chains (
    id text PRIMARY KEY,
    cfg jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    enabled boolean NOT NULL DEFAULT true
);

-- Table Definition ----------------------------------------------

CREATE TABLE terra_msgs (
    id BIGSERIAL PRIMARY KEY,
    terra_chain_id text NOT NULL REFERENCES terra_chains(id) ON DELETE CASCADE,
    contract_id text NOT NULL,
    raw bytea NOT NULL,
    state text NOT NULL,
    tx_hash text,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    type text NOT NULL DEFAULT '/terra.wasm.v1beta1.MsgExecuteContract'::text,
    CONSTRAINT terra_msgs_check CHECK (tx_hash <> NULL::text OR state <> 'broadcasted'::text AND state <> 'confirmed'::text)
);

-- Indices -------------------------------------------------------

CREATE INDEX idx_terra_msgs_terra_chain_id_state_contract_id ON terra_msgs(terra_chain_id text_ops,state text_ops,contract_id text_ops);

-- Table Definition ----------------------------------------------

CREATE TABLE terra_nodes (
    id SERIAL PRIMARY KEY,
    name character varying(255) NOT NULL CHECK (name::text <> ''::text),
    terra_chain_id text NOT NULL REFERENCES terra_chains(id) ON DELETE CASCADE,
    tendermint_url text CHECK (tendermint_url <> ''::text),
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

-- Indices -------------------------------------------------------

CREATE INDEX idx_nodes_terra_chain_id ON terra_nodes(terra_chain_id text_ops);
CREATE UNIQUE INDEX idx_terra_nodes_unique_name ON terra_nodes((lower(name::text)) text_ops);


