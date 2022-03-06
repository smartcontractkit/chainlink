-- +goose Up
CREATE TABLE evm_forwarders (
    id BIGSERIAL PRIMARY KEY,
    address bytea NOT NULL UNIQUE,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    evm_chain_id numeric(78,0) NOT NULL REFERENCES evm_chains(id),
    CONSTRAINT chk_address_length CHECK ((octet_length(address) = 20))
);

CREATE INDEX idx_forwarders_evm_chain_id ON evm_forwarders(evm_chain_id);

-- +goose Down
DROP TABLE evm_forwarders;
