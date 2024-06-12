-- +goose Up
-- +goose StatementBegin
ALTER TABLE evm.txes ADD COLUMN finalized boolean NOT NULL DEFAULT false;
ALTER TABLE evm.txes ADD CONSTRAINT chk_eth_txes_state_finalized CHECK (
    state <> 'confirmed'::eth_txes_state AND finalized = false
    OR
    state = 'confirmed'::eth_txes_state
) NOT VALID;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE evm.txes DROP CONSTRAINT chk_eth_txes_state_finalized;
ALTER TABLE evm.txes DROP COLUMN finalized;
-- +goose StatementEnd
