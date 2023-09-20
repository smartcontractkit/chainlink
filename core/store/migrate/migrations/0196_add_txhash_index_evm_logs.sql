-- +goose Up
CREATE INDEX evm_logs_idx_tx_hash ON evm.logs using brin (tx_hash);

-- +goose Down
DROP INDEX IF EXISTS evm_logs_idx_tx_hash;