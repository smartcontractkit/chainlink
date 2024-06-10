-- +goose Up
-- +goose StatementBegin
ALTER TABLE feeds_manager_chain_configs
ADD COLUMN account_address_public_key VARCHAR;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE feeds_manager_chain_configs DROP COLUMN account_address_public_key;

-- +goose StatementEnd
