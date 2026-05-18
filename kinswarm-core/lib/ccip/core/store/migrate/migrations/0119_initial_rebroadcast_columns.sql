-- +goose Up
ALTER TABLE eth_txes ADD COLUMN initial_broadcast_at timestamptz;
UPDATE eth_txes SET initial_broadcast_at = broadcast_at; -- Not perfect but this mirrors the old behaviour and will sort itself out in time when the old eth_txes are reaped
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
);
CREATE INDEX idx_eth_txes_initial_broadcast_at ON eth_txes USING BRIN (initial_broadcast_at timestamptz_minmax_ops);


-- +goose Down
ALTER TABLE eth_txes DROP COLUMN initial_broadcast_at;
ALTER TABLE eth_txes ADD CONSTRAINT chk_eth_txes_fsm CHECK (state = 'unstarted'::eth_txes_state AND nonce IS NULL AND error IS NULL AND broadcast_at IS NULL OR state = 'in_progress'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NULL OR state = 'fatal_error'::eth_txes_state AND nonce IS NULL AND error IS NOT NULL AND broadcast_at IS NULL OR state = 'unconfirmed'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL OR state = 'confirmed'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL OR state = 'confirmed_missing_receipt'::eth_txes_state AND nonce IS NOT NULL AND error IS NULL AND broadcast_at IS NOT NULL);
