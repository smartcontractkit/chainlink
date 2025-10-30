-- +goose Up
CREATE SCHEMA IF NOT EXISTS ton;

-- Create filters table
CREATE TABLE IF NOT EXISTS ton.log_poller_filters (
  id BIGSERIAL PRIMARY KEY,
  chain_id TEXT NOT NULL,

  name VARCHAR(255) NOT NULL,
  address TEXT NOT NULL, -- TON address in user-friendly format. TODO: consider use BYTEA for address field
  msg_type VARCHAR(20) NOT NULL,
  event_sig BYTEA NOT NULL CHECK (octet_length(event_sig) = 4), -- CRC32 hash as 4-byte binary

  starting_seq_no INTEGER NOT NULL, -- TODO: BIGINT

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,

  CONSTRAINT check_msg_type CHECK (msg_type IN ('INTERNAL', 'EXTERNAL_IN', 'EXTERNAL_OUT'))
);

CREATE UNIQUE INDEX IF NOT EXISTS ton_log_poller_filter_name ON ton.log_poller_filters (chain_id, name) WHERE NOT is_deleted;
CREATE INDEX IF NOT EXISTS idx_filters_address_msgtype ON ton.log_poller_filters(address, msg_type);

-- Create logs table
CREATE TABLE IF NOT EXISTS ton.log_poller_logs (
  id BIGSERIAL PRIMARY KEY,
  filter_id BIGINT NOT NULL,
  chain_id TEXT NOT NULL,

  address TEXT NOT NULL, -- TON address in user-friendly format. TODO: consider use BYTEA for address field
  event_sig BYTEA NOT NULL CHECK (octet_length(event_sig) = 4), -- CRC32 hash as 4-byte binary
  data BYTEA, -- BOC-encoded cell data

  tx_hash BYTEA NOT NULL,
  tx_lt NUMERIC(20, 0) NOT NULL, -- tx_lt is a uint64 which doesn't fit inside a bigint
  tx_timestamp TIMESTAMPTZ NOT NULL,
  msg_lt NUMERIC(20, 0) NOT NULL, -- msg_lt is a uint64 which doesn't fit inside a bigint
  msg_index INTEGER NOT NULL, -- message index within a transaction

  block_workchain INT NOT NULL,
  block_shard BIGINT NOT NULL,
  block_seqno BIGINT NOT NULL,
  block_root_hash BYTEA NOT NULL,
  block_file_hash BYTEA NOT NULL,

  master_block_seqno BIGINT NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_logs_filter FOREIGN KEY (filter_id) REFERENCES ton.log_poller_filters(id) ON DELETE CASCADE
);

-- Unique constraint to prevent duplicate log entries
CREATE UNIQUE INDEX IF NOT EXISTS idx_logs_unique ON ton.log_poller_logs (tx_hash, tx_lt, msg_index);

-- Index for address-scoped message ordering
CREATE INDEX IF NOT EXISTS idx_logs_address_msglt ON ton.log_poller_logs(address, msg_lt);
CREATE INDEX IF NOT EXISTS idx_logs_address_event ON ton.log_poller_logs(address, event_sig);

-- Index first 64 bytes of BOC data for byte-level filtering, covers most common filtering patterns
CREATE INDEX IF NOT EXISTS idx_logs_data_prefix64 ON ton.log_poller_logs (address, event_sig, SUBSTRING(data, 1, 64));

-- +goose Down
DROP INDEX IF EXISTS idx_logs_data_prefix64;
DROP INDEX IF EXISTS idx_logs_address_event;
DROP INDEX IF EXISTS idx_logs_address_msglt;
DROP INDEX IF EXISTS idx_logs_unique;
DROP TABLE IF EXISTS ton.log_poller_logs;
DROP INDEX IF EXISTS idx_filters_address_msgtype;
DROP INDEX IF EXISTS ton_log_poller_filter_name;
DROP TABLE IF EXISTS ton.log_poller_filters;
DROP SCHEMA IF EXISTS ton;

