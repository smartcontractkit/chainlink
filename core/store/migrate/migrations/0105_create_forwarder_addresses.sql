-- +goose Up
CREATE TABLE evm_forwarders (
    id BIGSERIAL PRIMARY KEY,
    address bytea NOT NULL UNIQUE,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    evm_chain_id numeric(78,0) NOT NULL REFERENCES evm_chains(id) ON DELETE CASCADE,
    CONSTRAINT chk_address_length CHECK ((octet_length(address) = 20))
);

CREATE INDEX idx_forwarders_evm_chain_id ON evm_forwarders(evm_chain_id);
CREATE INDEX idx_forwarders_evm_address ON evm_forwarders(address);
CREATE INDEX idx_forwarders_created_at ON evm_forwarders USING brin (created_at);
CREATE INDEX idx_forwarders_updated_at ON evm_forwarders USING brin (updated_at);

-- +goose Down
DROP TABLE evm_forwarders;
