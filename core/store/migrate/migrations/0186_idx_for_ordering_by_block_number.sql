-- +goose Up
-- +goose StatementBegin

-- This index optimizes queries used in the following functions:
-- * logpoller.LatestBlockByEventSigsAddrsWithConfs
--
-- Example query:
-- SELECT COALESCE(MAX(block_number), 0) FROM evm_logs
--         WHERE evm_chain_id = 420 AND
--         event_sig = '\0xA' AND
--         address = '\0xB' AND
--         (block_number) <= (SELECT COALESCE(block_number, 0) FROM evm_log_poller_blocks WHERE evm_chain_id = 420 ORDER BY block_number DESC LIMIT 1);
CREATE INDEX idx_evm_logs_ordered_by_block_number
    on public.evm_logs (evm_chain_id, address, event_sig, block_number DESC);
-- +goose StatementEnd


-- +goose Down

-- +goose StatementBegin
DROP INDEX IF EXISTS idx_evm_logs_ordered_by_block_number;
-- +goose StatementEnd
