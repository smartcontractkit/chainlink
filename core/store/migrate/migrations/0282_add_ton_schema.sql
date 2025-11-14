-- +goose Up
CREATE SCHEMA IF NOT EXISTS ton;

-- TON log poller filters: defines which messages to capture
CREATE TABLE IF NOT EXISTS ton.log_poller_filters (
  id BIGSERIAL PRIMARY KEY,
  chain_id TEXT NOT NULL, -- chain identifier
  name VARCHAR(255) NOT NULL, -- filter name
  address BYTEA NOT NULL CHECK (octet_length(address) = 36), -- TON address: 4 bytes workchain + 32 bytes hash (codec.ToRawAddr)
  msg_type VARCHAR(20) NOT NULL, -- INTERNAL/EXTERNAL_IN/EXTERNAL_OUT (tlb.MsgType)
  event_sig BYTEA NOT NULL CHECK (octet_length(event_sig) = 4), -- CRC32 hash, BE
  starting_seq_no BIGINT NOT NULL, -- start polling from this block seqno
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE,

  CONSTRAINT check_msg_type CHECK (msg_type IN ('INTERNAL', 'EXTERNAL_IN', 'EXTERNAL_OUT'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_filters_name ON ton.log_poller_filters (chain_id, name) WHERE NOT is_deleted;
CREATE INDEX IF NOT EXISTS idx_filters_address_msgtype ON ton.log_poller_filters(address, msg_type);

-- TON log poller logs: captured messages matching filters
CREATE TABLE IF NOT EXISTS ton.log_poller_logs (
  id BIGSERIAL PRIMARY KEY,
  filter_id BIGINT NOT NULL, -- filter that captured this log
  chain_id TEXT NOT NULL, -- chain identifier
  address BYTEA NOT NULL CHECK (octet_length(address) = 36), -- contract address (36 bytes)
  event_sig BYTEA NOT NULL CHECK (octet_length(event_sig) = 4), -- event type (4-byte CRC32)
  data_header BYTEA NOT NULL, -- BOC header: magic + flags + metadata (boc.HeaderLen)
  data_payload BYTEA NOT NULL, -- BOC payload: 2-byte cell descriptor + data (queryparser.go uses this)
  tx_hash BYTEA NOT NULL, -- transaction hash (32 bytes)
  tx_lt NUMERIC(20, 0) NOT NULL, -- transaction logical time (uint64 as NUMERIC, strconv.FormatUint)
  tx_timestamp TIMESTAMPTZ NOT NULL, -- transaction wall-clock time
  msg_lt NUMERIC(20, 0) NOT NULL, -- message logical time (uint64, used for pagination cursor)
  msg_index INTEGER NOT NULL, -- message index within transaction (deterministic ordering)
  block_workchain INT NOT NULL, -- workchain ID (usually 0=basechain)
  block_shard BIGINT NOT NULL, -- shard block id 
  block_seqno BIGINT NOT NULL, -- shard block sequence number
  block_root_hash BYTEA NOT NULL, -- shard block root hash (32 bytes)
  block_file_hash BYTEA NOT NULL, -- shard block file hash (32 bytes)
  master_block_seqno BIGINT NOT NULL, -- masterchain block sequence number
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_logs_filter FOREIGN KEY (filter_id) REFERENCES ton.log_poller_filters(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_logs_unique ON ton.log_poller_logs (tx_hash, tx_lt, msg_index); -- uniqueness: one message per (tx_hash, tx_lt, msg_index)
CREATE INDEX IF NOT EXISTS idx_logs_filter ON ton.log_poller_logs(chain_id, address, event_sig); -- query: WHERE chain_id = ? AND address = ? AND event_sig = ?
CREATE INDEX IF NOT EXISTS idx_logs_chrono ON ton.log_poller_logs(chain_id, address, event_sig, tx_lt, tx_timestamp); -- query: ORDER BY tx_lt, tx_timestamp (time-ordered events)
CREATE INDEX IF NOT EXISTS idx_logs_page ON ton.log_poller_logs(chain_id, address, msg_lt); -- query: WHERE (address, msg_lt) > (?, ?) (cursor pagination)

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

