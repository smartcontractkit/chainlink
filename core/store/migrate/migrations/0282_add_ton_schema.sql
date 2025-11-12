-- +goose Up
CREATE SCHEMA IF NOT EXISTS ton;

-- Create filters table
CREATE TABLE IF NOT EXISTS ton.log_poller_filters (
  id BIGSERIAL PRIMARY KEY,
  chain_id TEXT NOT NULL,

  name VARCHAR(255) NOT NULL,
  address BYTEA NOT NULL CHECK (octet_length(address) = 36), -- TON address in raw format (4 bytes workchain + 32 bytes data)
  msg_type VARCHAR(20) NOT NULL,
  event_sig BYTEA NOT NULL CHECK (octet_length(event_sig) = 4), -- CRC32 hash as 4-byte binary

  starting_seq_no BIGINT NOT NULL,

  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,

  CONSTRAINT check_msg_type CHECK (msg_type IN ('INTERNAL', 'EXTERNAL_IN', 'EXTERNAL_OUT'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_filters_name ON ton.log_poller_filters (chain_id, name) WHERE NOT is_deleted;
CREATE INDEX IF NOT EXISTS idx_filters_address_msgtype ON ton.log_poller_filters(address, msg_type);

-- Create logs table
CREATE TABLE IF NOT EXISTS ton.log_poller_logs (
  id BIGSERIAL PRIMARY KEY,
  filter_id BIGINT NOT NULL,
  chain_id TEXT NOT NULL,

  address BYTEA NOT NULL CHECK (octet_length(address) = 36), -- TON address in raw format (4 bytes workchain + 32 bytes data)
  event_sig BYTEA NOT NULL CHECK (octet_length(event_sig) = 4), -- CRC32 hash as 4-byte binary
  data_header BYTEA NOT NULL, -- BOC header (variable size: magic + flags + metadata)
  data_payload BYTEA NOT NULL, -- BOC payload starting with 2-byte cell descriptor

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

-- Generic filtering index: base filter for all log queries
CREATE INDEX IF NOT EXISTS idx_logs_filter ON ton.log_poller_logs(chain_id, address, event_sig);

-- Generic chronological index: time-ordered queries (supports both tx_lt and tx_timestamp)
-- Covers: message ordering, sequence queries, time-based sorting (both logical time and timestamp)
CREATE INDEX IF NOT EXISTS idx_logs_chrono ON ton.log_poller_logs(chain_id, address, event_sig, tx_lt, tx_timestamp);

-- Generic pagination index: cursor-based result pagination
CREATE INDEX IF NOT EXISTS idx_logs_page ON ton.log_poller_logs(chain_id, address, msg_lt);

-- +goose Down
DROP INDEX IF EXISTS ton.idx_logs_page;
DROP INDEX IF EXISTS ton.idx_logs_chrono;
DROP INDEX IF EXISTS ton.idx_logs_filter;
DROP INDEX IF EXISTS ton.idx_logs_unique;
DROP TABLE IF EXISTS ton.log_poller_logs;
DROP INDEX IF EXISTS ton.idx_filters_address_msgtype;
DROP INDEX IF EXISTS ton.idx_filters_name;
DROP TABLE IF EXISTS ton.log_poller_filters;
DROP SCHEMA IF EXISTS ton;

