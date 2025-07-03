-- +goose Up

-- These indexes are optimized for CCIP read use cases related to fetching CCIPMessageSent events from
-- evm.logs table
--
-- event CCIPMessageSent(
--    uint64 indexed destChainSelector, // topics[2]
--    uint64 indexed sequenceNumber,    // topics[3]
--    Internal.EVM2AnyRampMessage message
--  );
--
-- By making it a partial-index custom to CCIPMessageSent event, we can reduce the bloat of the index
-- and improve performance for CCIP use cases. Therefore, it doesn't introduce any perf/space overhead
-- for products that don't need it.
--
-- idx_evm_logs_ccip_message_sent_read_seq
-- This it the main index used for reading CCIPMessageSent events for sequence number range.
-- In order to reduce the heap access cost, we include most of the fields that are needed for the query
-- Data column can't be included in the index because it's unbounded and might hit the Postgres limit of 8kB
-- topics[2] is used to filter by destChainSelector and topics[3] is used to filter by sequenceNumber
-- We start index with address column because this filter should immediately narrow down the result set
-- it's caused by the fact that a single chain might have multiple OnRamp contracts deployed during active-candidate
-- (meaning it should have at least the same efficiency as evm_chain_id)
CREATE INDEX idx_evm_logs_ccip_message_sent_read_seq
    ON evm.logs (
                 address,
                 evm_chain_id,
                 (topics[2]),
                 (topics[3]),
                 block_number
        )
    INCLUDE (
        log_index,
        block_hash,
        event_sig,
        topics,
        tx_hash,
        created_at,
        block_timestamp
        )
    WHERE event_sig = '\x192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f32';


-- idx_evm_logs_ccip_message_sent_read_latest
-- This index is used for fetching the latest CCIPMessageSent which is finalized.
-- Ideally we should be able to fetch only topics[2] from the LogPoller, but that feature is not supported by CR
-- yet and all the fields are always passed to the LogPoller.
-- However, considering that eventually we need to fetch only topics[2] and we select only a single row,
-- we don't add any include clauses with remaining columns to the index.
CREATE INDEX idx_evm_logs_ccip_message_sent_read_latest
    ON evm.logs (
                 address,
                 evm_chain_id,
                 (topics[2]),
                 block_number DESC
        )
    include (
        topics,
        event_sig,
        block_number,
        log_index,
        tx_hash
        )
    WHERE event_sig = '\x192442a2b2adb6a7948f097023cb6b57d29d3a7a5dd33e6666d33c39cc456f32';


-- +goose Down
DROP INDEX IF EXISTS evm.idx_evm_logs_ccip_message_sent_read_seq;
DROP INDEX IF EXISTS evm.idx_evm_logs_ccip_message_sent_read_latest;

