-- +goose Up

CREATE SCHEMA ccip;

CREATE TABLE ccip.observed_gas_prices(
    chain_selector NUMERIC(20,0) NOT NULL,
    job_id INTEGER NOT NULL,
    source_chain_selector NUMERIC(20,0) NOT NULL,
    value NUMERIC(78,0) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
);

CREATE TABLE ccip.observed_token_prices(
     chain_selector NUMERIC(20,0) NOT NULL,
     job_id INTEGER NOT NULL,
     token_addr BYTEA CHECK (octet_length(token_addr) = 20) NOT NULL,
     value NUMERIC(78,0) NOT NULL,
     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
);

CREATE INDEX idx_ccip_gas_prices_chain_value_timestamp ON ccip.observed_gas_prices (chain_selector, source_chain_selector, created_at DESC);
CREATE INDEX idx_ccip_token_prices_token_value_timestamp ON ccip.observed_token_prices (chain_selector, token_addr, created_at DESC);


-- +goose Down

DROP INDEX IF EXISTS idx_ccip_token_prices_token_value;
DROP INDEX IF EXISTS idx_ccip_gas_prices_chain_value;

DROP TABLE "ccip".observed_token_prices;
DROP TABLE "ccip".observed_gas_prices;

DROP SCHEMA "ccip";