-- +goose Up
ALTER TABLE cron_specs ADD COLUMN evm_chain_id numeric(78,0);

-- +goose Down
ALTER TABLE cron_specs DROP COLUMN evm_chain_id;
