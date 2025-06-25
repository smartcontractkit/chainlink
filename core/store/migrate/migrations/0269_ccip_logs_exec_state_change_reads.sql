-- +goose Up

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

