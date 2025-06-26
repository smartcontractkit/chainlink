-- +goose Up

-- This index is optimized for CCIP use cases related to fetching ExecutionStateChanged events from
-- the evm.logs table (part of the Execution plugin fetching data from the destination chain).
--   event ExecutionStateChanged(
--     uint64 indexed sourceChainSelector,
--     uint64 indexed sequenceNumber,
--     bytes32 indexed messageId,
--     bytes32 messageHash,
--     Internal.MessageExecutionState state,
--     bytes returnData,
--     uint256 gasUsed
--   );
--
-- By making it a partial-index custom to ExecutionStateChanged event, we can reduce the bloat of the index
-- and improve performance for CCIP use cases. Therefore, it doesn't introduce any perf/space overhead
-- for products that don't need it.
--
-- idx_evm_logs_ccip_exec_state_change_read
-- This is the main index used for reading ExecutionStateChanged events for sequence number range (topic[2])
-- and source chain selector (topic[3]).
CREATE INDEX idx_evm_logs_ccip_exec_state_change_read
    ON evm.logs (
                 address,
                 evm_chain_id,
                 (topics[2]),
                 (topics[3]),
                 block_number,
                 log_index,
                 tx_hash
        )
    WHERE event_sig = '\x05665fe9ad095383d018353f4cbcba77e84db27dd215081bbf7cdf9ae6fbe48b';


-- +goose Down
DROP INDEX IF EXISTS evm.idx_evm_logs_ccip_exec_state_change_read;

