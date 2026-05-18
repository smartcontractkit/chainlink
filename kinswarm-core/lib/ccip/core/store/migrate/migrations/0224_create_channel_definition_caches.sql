-- +goose Up
CREATE TABLE channel_definitions (
    evm_chain_id NUMERIC(78) NOT NULL,
    addr bytea CHECK (octet_length(addr) = 20),
    definitions JSONB NOT NULL,
    block_num BIGINT NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    PRIMARY KEY (evm_chain_id, addr)
);

-- +goose Down
DROP TABLE channel_definitions;
