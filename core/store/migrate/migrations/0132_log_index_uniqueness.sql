-- +goose Up

-- Add tx_index column to log_broadcasts
ALTER TABLE log_broadcasts ADD COLUMN tx_index BIGINT;

DROP INDEX IF EXISTS log_broadcasts_unique_idx;
CREATE UNIQUE INDEX log_broadcasts_unique_idx ON log_broadcasts USING BTREE (job_id, block_hash, tx_index, log_index, evm_chain_id);

-- +goose Down

DROP INDEX IF EXISTS log_broadcasts_unique_idx;
ALTER TABLE log_broadcasts DROP COLUMN tx_index;
CREATE UNIQUE INDEX log_broadcasts_unique_idx ON log_broadcasts USING BTREE (job_id, block_hash, log_index, evm_chain_id);
