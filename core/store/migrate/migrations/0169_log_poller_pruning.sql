-- +goose Up
-- +goose StatementBegin
ALTER TABLE evm_log_poller_filters ADD COLUMN retention BIGINT DEFAULT 0;
CREATE INDEX evm_logs_idx_created_at ON evm_logs (created_at);
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP INDEX evm_logs_idx_created_at;
ALTER TABLE evm_log_poller_filters DROP COLUMN retention;
-- +goose StatementEnd
