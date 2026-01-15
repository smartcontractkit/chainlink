-- +goose Up
ALTER TABLE evm.txes ADD COLUMN IF NOT EXISTS max_gas_price NUMERIC(78,0);

-- +goose Down
ALTER TABLE evm.txes DROP COLUMN IF EXISTS max_gas_price;

