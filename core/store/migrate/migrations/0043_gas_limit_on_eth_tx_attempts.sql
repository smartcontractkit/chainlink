-- +goose Up
ALTER TABLE eth_tx_attempts ADD COLUMN chain_specific_gas_limit bigint;
UPDATE eth_tx_attempts
SET chain_specific_gas_limit = eth_txes.gas_limit
FROM eth_txes
WHERE eth_txes.id = eth_tx_attempts.eth_tx_id;
ALTER TABLE eth_tx_attempts ALTER COLUMN chain_specific_gas_limit SET NOT NULL;

-- +goose Down
ALTER TABLE eth_tx_attempts DROP COLUMN chain_specific_gas_limit;
