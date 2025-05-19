-- +goose Up
-- +goose StatementBegin
-- This removes evm from the search path, to avoid risk of an unqualified query intended to act on a Solana table accidentally
-- acting on the corresponding evm table. This makes schema qualification mandatory for all chain-specific queries.
SET search_path TO public;
CREATE SCHEMA solana;

CREATE TABLE solana.log_poller_filters (
    id BIGSERIAL PRIMARY KEY,
    chain_id TEXT NOT NULL, -- use human-readable name instead of genesis block hash to reduce index size by 50%
    name TEXT NOT NULL,
    address BYTEA NOT NULL,
    event_name TEXT NOT NULL,
    event_sig BYTEA NOT NULL,
    starting_block BIGINT NOT NULL,
    event_idl TEXT,
    subkey_paths json, -- A list of subkeys to be indexed, represented by their json paths in the event struct. Forced to use json's 2d array as text[][] requires all paths to have equal length.
    retention BIGINT NOT NULL DEFAULT 0, -- we don’t have to implement this initially, but good to include it in the schema
    max_logs_kept BIGINT NOT NULL DEFAULT 0, -- same as retention, no need to implement yet
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE UNIQUE INDEX IF NOT EXISTS solana_log_poller_filter_name ON solana.log_poller_filters (chain_id, name) WHERE NOT is_deleted;

CREATE TABLE solana.logs (
    id               BIGSERIAL PRIMARY KEY,
    filter_id        BIGINT NOT NULL REFERENCES solana.log_poller_filters (id) ON DELETE CASCADE,
    chain_id         TEXT                     not null, -- use human-readable name instead of genesis block hash to reduce index size by 50%
    log_index        bigint                    not null,
    block_hash       bytea                     not null,
    block_number     bigint                    not null CHECK (block_number > 0),
    block_timestamp  timestamp with time zone  not null,
    address          bytea                     not null,
    event_sig        bytea                     not null,
    subkey_values    bytea[]                   not null,
    tx_hash          bytea                     not null,
    data             bytea                     not null,
    created_at       timestamp with time zone  not null,
    expires_at       timestamp with time zone  null, -- null to indicate no timebase expiry
    sequence_num     bigint                    not null
);

CREATE INDEX IF NOT EXISTS solana_logs_idx_by_timestamp ON solana.logs (chain_id, address, event_sig, block_timestamp, block_number);

CREATE INDEX IF NOT EXISTS solana_logs_idx ON solana.logs (chain_id, block_number, address, event_sig);

CREATE INDEX IF NOT EXISTS solana_logs_idx_subkey_one ON solana.logs ((subkey_values[1]));
CREATE INDEX IF NOT EXISTS solana_logs_idx_subkey_two ON solana.logs ((subkey_values[2]));
CREATE INDEX IF NOT EXISTS solana_logs_idx_subkey_three ON solana.logs ((subkey_values[3]));
CREATE INDEX IF NOT EXISTS solana_logs_idx_subkey_four ON solana.logs ((subkey_values[4]));

CREATE INDEX IF NOT EXISTS solana_logs_idx_tx_hash ON solana.logs (tx_hash);

-- Useful for the current form of those queries:  WHERE chain_id = $1 AND address = $2 AND event_sig = $3 ... ORDER BY block_number, log_index
--  log_index is not included in this index because it increases the index size by ~ 10% for a likely negligible performance benefit
CREATE INDEX IF NOT EXISTS solana_logs_idx_chain_address_event_block ON solana.logs (chain_id, address, event_sig, block_number);

CREATE UNIQUE INDEX IF NOT EXISTS solana_logs_idx_chain_filter_block_logindex ON solana.logs USING btree (chain_id, filter_id, block_number, log_index);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS solana.logs;
DROP TABLE IF EXISTS solana.log_poller_filters;

DROP SCHEMA solana;

SET search_path TO public,evm;
-- +goose StatementEnd
