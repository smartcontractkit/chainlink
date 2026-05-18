-- +goose Up
-- +goose StatementBegin
DROP INDEX IF EXISTS log_broadcasts_unique_idx;
ALTER TABLE log_broadcasts DROP COLUMN tx_index;
CREATE UNIQUE INDEX log_broadcasts_unique_idx ON log_broadcasts USING BTREE (job_id, block_hash, log_index, evm_chain_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX log_broadcasts_unique_idx;
ALTER TABLE log_broadcasts ADD COLUMN tx_index BIGINT NOT NULL DEFAULT -1;
CREATE UNIQUE INDEX log_broadcasts_unique_idx ON log_broadcasts USING BTREE (job_id, block_hash, tx_index, log_index, evm_chain_id);
-- +goose StatementEnd
