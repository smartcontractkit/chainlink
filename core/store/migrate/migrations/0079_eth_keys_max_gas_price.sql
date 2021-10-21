-- +goose Up
ALTER TABLE eth_key_states ADD COLUMN max_gas_gwei numeric(78,0);
-- +goose Down
