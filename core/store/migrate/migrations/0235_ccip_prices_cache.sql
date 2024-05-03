-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA ccip;

CREATE TABLE ccip.observed_gas_prices(
    chain_selector NUMERIC(20,0) NOT NULL,
    job_id INTEGER NOT NULL,
    source_chain_selector NUMERIC(20,0) NOT NULL,
    gas_price NUMERIC(78,0) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE ccip.observed_token_prices(
     chain_selector NUMERIC(20,0) NOT NULL,
     job_id INTEGER NOT NULL,
     token_addr BYTEA NOT NULL,
     token_price NUMERIC(78,0) NOT NULL,
     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ccip_gas_prices_chain_gas_price_timestamp ON ccip.observed_gas_prices (chain_selector, source_chain_selector, created_at DESC);
CREATE INDEX idx_ccip_token_prices_token_price_timestamp ON ccip.observed_token_prices (chain_selector, token_addr, created_at DESC);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_ccip_token_prices_token_value;
DROP INDEX IF EXISTS idx_ccip_gas_prices_chain_value;

DROP TABLE ccip.observed_token_prices;
DROP TABLE ccip.observed_gas_prices;

DROP SCHEMA ccip;
-- +goose StatementEnd
