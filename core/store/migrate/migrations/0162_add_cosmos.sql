-- +goose Up

-- Table Definition ----------------------------------------------

CREATE TABLE cosmos_msgs (
    id BIGSERIAL PRIMARY KEY,
    cosmos_chain_id text NOT NULL,
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

-- +goose Down

DROP TABLE cosmos_msgs;
