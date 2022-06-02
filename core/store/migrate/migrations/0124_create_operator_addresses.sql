-- +goose Up
CREATE TABLE operators (
    id BIGSERIAL PRIMARY KEY,
    address bytea NOT NULL UNIQUE,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    chain_id numeric(78,0) NOT NULL REFERENCES evm_chains(id) ON DELETE CASCADE,
    CONSTRAINT chk_address_length CHECK ((octet_length(address) = 20))
);

CREATE INDEX idx_operators_chain_id ON operators(chain_id);
CREATE INDEX idx_operators_address ON operators(address);
CREATE INDEX idx_operators_created_at ON operators USING brin (created_at);
CREATE INDEX idx_operators_updated_at ON operators USING brin (updated_at);

-- +goose Down
DROP TABLE operators;
