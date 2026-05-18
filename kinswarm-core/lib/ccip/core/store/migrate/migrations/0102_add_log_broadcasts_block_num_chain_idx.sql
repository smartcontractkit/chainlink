-- +goose Up
CREATE INDEX idx_log_broadcasts_block_number_evm_chain_id ON log_broadcasts (evm_chain_id, block_number);

-- +goose Down
DROP INDEX IF EXISTS idx_log_broadcasts_block_number_evm_chain_id;
