-- +goose Up
-- +goose StatementBegin

-- Grandfather in old attempts which inadvertently saved with a 0 value, and
-- enforce correctly writing data for future attempts
UPDATE eth_tx_attempts SET chain_specific_gas_limit=1 WHERE chain_specific_gas_limit=0;
ALTER TABLE eth_tx_attempts ADD CONSTRAINT chk_chain_specific_gas_limit_not_zero CHECK (chain_specific_gas_limit > 0);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE eth_tx_attempts DROP CONSTRAINT chk_chain_specific_gas_limit_not_zero;

-- +goose StatementEnd
