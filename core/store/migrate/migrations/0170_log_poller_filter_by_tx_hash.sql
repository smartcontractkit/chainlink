-- +goose Up
CREATE INDEX evm_logs_idx_by_tx_hash ON evm_logs(evm_chain_id, tx_hash, log_index ASC);

-- +goose Down
DROP INDEX evm_logs_idx_by_tx_hash;
