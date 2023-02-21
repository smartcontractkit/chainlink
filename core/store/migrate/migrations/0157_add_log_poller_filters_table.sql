-- +goose Up

CREATE TABLE evm_log_poller_filters(
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL CHECK (length(name) > 0),
    address BYTEA CHECK (octet_length(address) = 20) NOT NULL,
    event BYTEA CHECK (octet_length(event) = 32) NOT NULL,
    evm_chain_id numeric(78,0) REFERENCES evm_chains (id) DEFERRABLE INITIALLY IMMEDIATE,
    created_at TIMESTAMPTZ NOT NULL,
    UNIQUE (name, evm_chain_id, address, event)
);

-- +goose Down

DROP TABLE evm_log_poller_filters CASCADE;

