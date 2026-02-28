-- Migration: Add evm_trigger_cursors table for persistent log event trigger cursors
-- This table stores the last-processed cursor position for each log event trigger.
-- On node restart, triggers resume from the persisted position instead of
-- recalculating from currentBlock - LookbackBlocks.

CREATE TABLE IF NOT EXISTS evm_trigger_cursors (
    trigger_id    TEXT PRIMARY KEY,
    chain_id      TEXT NOT NULL,
    cursor_value  TEXT NOT NULL DEFAULT '',
    block_number  BIGINT NOT NULL DEFAULT 0,
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_evm_trigger_cursors_chain_id ON evm_trigger_cursors (chain_id);
