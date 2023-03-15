-- +goose Up
ALTER TABLE evm_logs ADD COLUMN block_timestamp timestamptz NOT NULL DEFAULT now();
ALTER TABLE evm_logs ALTER COLUMN block_timestamp DROP DEFAULT;
ALTER TABLE evm_log_poller_blocks ADD COLUMN block_timestamp timestamptz NOT NULL DEFAULT now();
ALTER TABLE evm_log_poller_blocks ALTER COLUMN block_timestamp DROP DEFAULT;
-- +goose Down
ALTER TABLE evm_log_poller_blocks DROP COLUMN block_timestamp;
ALTER TABLE evm_logs DROP COLUMN block_timestamp;
