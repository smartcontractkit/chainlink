-- +goose Up
ALTER TABLE evm.tx_attempts ADD COLUMN is_purge_attempt boolean NOT NULL DEFAULT false;
-- +goose Down
ALTER TABLE evm.tx_attempts DROP COLUMN is_purge_attempt;
