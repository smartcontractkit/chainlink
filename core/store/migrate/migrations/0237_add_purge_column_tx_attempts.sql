-- +goose Up
ALTER TABLE evm.tx_attempts ADD COLUMN is_purge_attempt boolean NOT NULL DEFAULT false;
ALTER TABLE evm.txes DROP CONSTRAINT chk_eth_txes_fsm;
ALTER TABLE evm.txes ADD CONSTRAINT chk_eth_txes_fsm CHECK (
    state = 'unstarted'::eth_txes_state AND nonce IS NULL AND error IS NULL AND broadcast_at IS NULL AND initial_broadcast_at IS NULL
    OR
    state = 'in_progress'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NULL AND initial_broadcast_at IS NULL
    OR
    state = 'fatal_error'::eth_txes_state AND error IS NOT NULL
    OR
    state = 'unconfirmed'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
    OR
    state = 'confirmed'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
    OR
    state = 'confirmed_missing_receipt'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
) NOT VALID;
-- +goose Down
ALTER TABLE evm.tx_attempts DROP COLUMN is_purge_attempt;
ALTER TABLE evm.txes DROP CONSTRAINT chk_eth_txes_fsm;
ALTER TABLE evm.txes ADD CONSTRAINT chk_eth_txes_fsm CHECK (
    state = 'unstarted'::eth_txes_state AND nonce IS NULL AND error IS NULL AND broadcast_at IS NULL AND initial_broadcast_at IS NULL
    OR
    state = 'in_progress'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NULL AND initial_broadcast_at IS NULL
    OR
    state = 'fatal_error'::eth_txes_state AND nonce IS NULL AND error IS NOT NULL
    OR
    state = 'unconfirmed'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
    OR
    state = 'confirmed'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
    OR
    state = 'confirmed_missing_receipt'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
) NOT VALID;