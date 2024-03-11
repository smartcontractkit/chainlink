-- +goose Up
DROP INDEX evm.evm_logs_idx_created_at;

-- +goose Down
CREATE INDEX evm_logs_idx_created_at ON evm.logs (created_at);
