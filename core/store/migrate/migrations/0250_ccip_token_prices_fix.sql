-- +goose Up

-- We need to re-create tables from scratch because of the unique constraint on tokens and chains selectors
DROP TABLE ccip.observed_token_prices;
DROP TABLE ccip.observed_gas_prices;

CREATE TABLE ccip.observed_token_prices
(
    chain_selector NUMERIC(20, 0) NOT NULL,
    token_addr     BYTEA          NOT NULL,
    token_price    NUMERIC(78, 0) NOT NULL,
    updated_at     TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chain_selector, token_addr)
);

CREATE TABLE ccip.observed_gas_prices
(
    chain_selector        NUMERIC(20, 0) NOT NULL,
    source_chain_selector NUMERIC(20, 0) NOT NULL,
    gas_price             NUMERIC(78, 0) NOT NULL,
    updated_at            TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    PRIMARY KEY (chain_selector, source_chain_selector)
);

-- +goose Down
DROP TABLE ccip.observed_token_prices;
DROP TABLE ccip.observed_gas_prices;

-- Restore state from migration 0236_ccip_prices_cache.sql
CREATE TABLE ccip.observed_gas_prices
(
    chain_selector        NUMERIC(20, 0) NOT NULL,
    job_id                INTEGER        NOT NULL,
    source_chain_selector NUMERIC(20, 0) NOT NULL,
    gas_price             NUMERIC(78, 0) NOT NULL,
    created_at            TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE TABLE ccip.observed_token_prices
(
    chain_selector NUMERIC(20, 0) NOT NULL,
    job_id         INTEGER        NOT NULL,
    token_addr     BYTEA          NOT NULL,
    token_price    NUMERIC(78, 0) NOT NULL,
    created_at     TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_ccip_gas_prices_chain_gas_price_timestamp ON ccip.observed_gas_prices (chain_selector, source_chain_selector, created_at DESC);
CREATE INDEX idx_ccip_token_prices_token_price_timestamp ON ccip.observed_token_prices (chain_selector, token_addr, created_at DESC);
