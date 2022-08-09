-- +goose Up
-- +goose StatementBegin
ALTER TABLE starknet_nodes RENAME COLUMN chain_id TO starknet_chain_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE starknet_nodes RENAME COLUMN starknet_chain_id TO chain_id;
-- +goose StatementEnd
