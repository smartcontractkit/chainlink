-- +goose Up
ALTER TABLE eth_txes DROP CONSTRAINT chk_eth_txes_fsm;
ALTER TABLE eth_txes ADD CONSTRAINT chk_eth_txes_fsm CHECK (
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
) NOT VALID; -- NOT VALID gives large speedup and this is a relaxing of the constraint so its safe

-- +goose Down
UPDATE eth_txes SET broadcast_at=NULL, initial_broadcast_at=NULL WHERE state='fatal_error';
ALTER TABLE eth_txes DROP CONSTRAINT chk_eth_txes_fsm;
ALTER TABLE eth_txes ADD CONSTRAINT chk_eth_txes_fsm CHECK (
    state = 'unstarted'::eth_txes_state AND nonce IS NULL AND error IS NULL AND broadcast_at IS NULL AND initial_broadcast_at IS NULL
    OR
    state = 'in_progress'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NULL AND initial_broadcast_at IS NULL
    OR
    state = 'fatal_error'::eth_txes_state AND nonce IS NULL AND error IS NOT NULL AND broadcast_at IS NULL AND initial_broadcast_at IS NULL
    OR
    state = 'unconfirmed'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
    OR
    state = 'confirmed'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
    OR
    state = 'confirmed_missing_receipt'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL AND initial_broadcast_at IS NOT NULL
) NOT VALID; -- NOT VALID gives large speedup and we know data is valid because of update above
