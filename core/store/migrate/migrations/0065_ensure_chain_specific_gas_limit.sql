-- +goose Up
-- +goose StatementBegin

-- NOT VALID is here to grandfather in old attempts which inadvertently saved
-- with a 0 value, but to enforce correctly writing data for future attempts
ALTER TABLE eth_tx_attempts ADD CONSTRAINT chk_chain_specific_gas_limit_not_zero CHECK (chain_specific_gas_limit > 0) NOT VALID;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE eth_tx_attempts DROP CONSTRAINT chk_chain_specific_gas_limit_not_zero;

-- +goose StatementEnd
