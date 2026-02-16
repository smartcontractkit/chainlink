-- +goose Up

-- Add pruning configuration columns to TON log poller filters
-- Enables per-filter retention policy: how long to keep logs (log_retention) and max count (max_logs_kept)
ALTER TABLE ton.log_poller_filters
    ADD COLUMN log_retention BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN max_logs_kept BIGINT NOT NULL DEFAULT 0;

-- Update unique constraint to include filter_id
-- This allows multiple filters to store the same blockchain event (same tx_hash, tx_lt, msg_index)
-- Query-time deduplication handles returning unique events to callers
DROP INDEX IF EXISTS ton.idx_logs_unique;
CREATE UNIQUE INDEX IF NOT EXISTS idx_logs_unique ON ton.log_poller_logs (chain_id, filter_id, tx_hash, tx_lt, msg_index);

-- Update filters index to include chain_id for efficient multi-chain filtering
DROP INDEX IF EXISTS ton.idx_filters_address_msgtype;
CREATE INDEX IF NOT EXISTS idx_filters_address_msgtype ON ton.log_poller_filters(chain_id, address, msg_type);

-- Checkpoint resumption index
CREATE INDEX IF NOT EXISTS idx_logs_master_block ON ton.log_poller_logs (chain_id, master_block_seqno DESC);

-- Add expires_at column for time-based pruning
ALTER TABLE ton.log_poller_logs ADD COLUMN expires_at TIMESTAMPTZ;
CREATE INDEX IF NOT EXISTS idx_logs_expires ON ton.log_poller_logs(expires_at) WHERE expires_at IS NOT NULL;

-- +goose Down

DROP INDEX IF EXISTS ton.idx_logs_expires;
ALTER TABLE ton.log_poller_logs DROP COLUMN IF EXISTS expires_at;
DROP INDEX IF EXISTS ton.idx_logs_master_block;
DROP INDEX IF EXISTS ton.idx_filters_address_msgtype;
CREATE INDEX IF NOT EXISTS idx_filters_address_msgtype ON ton.log_poller_filters(address, msg_type);
DROP INDEX IF EXISTS ton.idx_logs_unique;
CREATE UNIQUE INDEX IF NOT EXISTS idx_logs_unique ON ton.log_poller_logs (tx_hash, tx_lt, msg_index);

-- Remove pruning columns
ALTER TABLE ton.log_poller_filters
    DROP COLUMN IF EXISTS log_retention,
    DROP COLUMN IF EXISTS max_logs_kept;
