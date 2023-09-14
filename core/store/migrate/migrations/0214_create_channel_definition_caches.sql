-- +goose Up
CREATE TABLE streams_channel_definitions (
    addr bytea PRIMARY KEY CHECK (octet_length(addr) = 20),
    evm_chain_id NUMERIC(78) NOT NULL,
    definitions JSONB NOT NULL,
    block_num BIGINT NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

CREATE INDEX idx_streams_channel_definitions_evm_chain_id_addr ON streams_channel_definitions (evm_chain_id, addr);

-- +goose Down
DROP TABLE streams_channel_definitions;
