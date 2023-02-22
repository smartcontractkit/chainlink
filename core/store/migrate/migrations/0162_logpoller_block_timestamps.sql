-- +goose Up
ALTER TABLE evm_logs ADD COLUMN block_timestamp timestamptz NOT NULL DEFAULT now();

-- +goose Down
ALTER TABLE evm_logs DROP COLUMN block_timestamp;
