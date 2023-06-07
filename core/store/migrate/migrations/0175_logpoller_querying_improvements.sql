-- +goose Up
-- +goose StatementBegin

--  This index should make the following queries work in almost const time, instead of doing sequential scan.
--  This subquery is heavily used by most of the logpoller's functions
--
--  SELECT * FROM evm_log_poller_blocks
--      WHERE evm_chain_id = 420
--      ORDER BY block_number DESC
--      LIMIT 1;
CREATE INDEX idx_evm_log_poller_blocks_order_by_block
    ON evm_log_poller_blocks (evm_chain_id, block_number DESC);

-- This index optimizes queries used in the following funcitons:
-- * logpoller.LogsCreatedAfter
-- * logpoller.LatestLogByEventSigWithConfs
--
-- Example query:
-- SELECT * FROM evm_logs
--     WHERE evm_chain_id = 420
--     AND address = '\xABC'
--     AND event_sig = '\xABC'
--     AND block_number <= (SELECT COALESCE(block_number, 0) FROM evm_log_poller_blocks WHERE evm_chain_id = 420)
--     AND created_at > '2023-05-31T07:29:11.29Z'
--     ORDER BY created_at ASC;
CREATE INDEX idx_evm_logs_ordered_by_block_and_created_at
    ON evm_logs (evm_chain_id, address, event_sig, block_number, created_at);

-- +goose StatementEnd


-- +goose Down

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_evm_logs_ordered_by_block_and_created_at;
DROP INDEX IF EXISTS idx_evm_log_poller_blocks_order_by_block;
-- +goose StatementEnd
