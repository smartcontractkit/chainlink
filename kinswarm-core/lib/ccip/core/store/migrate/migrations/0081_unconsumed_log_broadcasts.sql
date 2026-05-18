-- +goose Up
ALTER TABLE log_broadcasts ADD COLUMN updated_at timestamp with time zone NOT NULL DEFAULT NOW();
DROP INDEX IF EXISTS log_consumptions_unique_v2_idx;
CREATE UNIQUE INDEX log_broadcasts_unique_idx ON log_broadcasts(job_id, block_hash, log_index, evm_chain_id);
CREATE TABLE log_broadcasts_pending (
    evm_chain_id numeric(78,0) REFERENCES evm_chains (id) DEFERRABLE INITIALLY IMMEDIATE PRIMARY KEY,
    block_number int8,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);
CREATE INDEX idx_log_broadcasts_unconsumed on log_broadcasts(evm_chain_id, block_number) WHERE consumed = false AND block_number IS NOT NULL;
-- +goose Down
DROP INDEX IF EXISTS idx_log_broadcasts_unconsumed;
DROP TABLE IF EXISTS log_broadcasts_pending;
ALTER TABLE log_broadcasts DROP COLUMN updated_at;
DROP INDEX IF EXISTS log_broadcasts_unique_idx;
CREATE UNIQUE INDEX log_consumptions_unique_v2_idx ON log_broadcasts(job_id, block_hash, log_index, consumed, evm_chain_id) WHERE job_id IS NOT NULL;
