-- +goose NO TRANSACTION
-- +goose Up

-- NOTE: see 0222_add_finalized_tx_state_part_2 for more details.

ALTER TYPE eth_txes_state ADD VALUE IF NOT EXISTS 'finalized' AFTER 'confirmed';

ALTER TYPE eth_tx_attempts_state ADD VALUE IF NOT EXISTS 'finalized' AFTER 'broadcast';




-- +goose Down

-- removal is handled in 0222_add_finalized_tx_state_part_2