-- +goose Up
create index evm_logs_idx_tx_hash on evm.logs (tx_hash);

-- +goose Down
DROP INDEX IF EXISTS evm_logs_idx_tx_hash;