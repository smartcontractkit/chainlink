-- +goose Up

CREATE TABLE log_poller_filters(
    id BIGSERIAL PRIMARY KEY,
    filter_name TEXT,
    address BYTEA CHECK (octet_length(address) = 20) NOT NULL,
    event BYTEA CHECK (octet_length(event) = 32) NOT NULL,
    evm_chain_id numeric(78,0) REFERENCES evm_chains (id) DEFERRABLE INITIALLY IMMEDIATE,
    created_at TIMESTAMPTZ NOT NULL,
    UNIQUE (filter_name, evm_chain_id, address, event)
);

CREATE INDEX idx_log_poller_filters_address_event ON log_poller_filters(filter_name);

-- +goose Down

DROP INDEX IF EXISTS idx_log_poller_filters_address_event;
DROP TABLE log_poller_filters;

