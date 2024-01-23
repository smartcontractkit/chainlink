-- +goose NO TRANSACTION
-- +goose Up

-- NOTE: see 0222_mark_old_txs_finalized for more details.

ALTER TYPE eth_txes_state ADD VALUE IF NOT EXISTS 'finalized' AFTER 'confirmed';

ALTER TYPE eth_tx_attempts_state ADD VALUE IF NOT EXISTS 'finalized' AFTER 'broadcast';




-- +goose Down

