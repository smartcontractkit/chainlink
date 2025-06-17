-- +goose Up
CREATE INDEX IF NOT EXISTS idx_txes_evm_chain_id_state_nonce
    ON evm.txes (evm_chain_id, state, nonce);

CREATE INDEX IF NOT EXISTS idx_tx_attempts_eth_tx_id_state_hash
    ON evm.tx_attempts (eth_tx_id, state, hash);

CREATE INDEX IF NOT EXISTS idx_receipts_tx_hash_id
    ON evm.receipts (tx_hash, id);

CREATE INDEX IF NOT EXISTS idx_receipts_tx_hash
    ON evm.receipts (tx_hash);

CREATE INDEX IF NOT EXISTS idx_tx_attempts_eth_tx_id_hash
    ON evm.tx_attempts (eth_tx_id, hash);

-- +goose Down
DROP INDEX IF EXISTS idx_txes_evm_chain_id_state_nonce;
DROP INDEX IF EXISTS idx_tx_attempts_eth_tx_id_state_hash;
DROP INDEX IF EXISTS idx_receipts_tx_hash_id;
DROP INDEX IF EXISTS idx_receipts_tx_hash;
DROP INDEX IF EXISTS idx_tx_attempts_eth_tx_id_hash;
