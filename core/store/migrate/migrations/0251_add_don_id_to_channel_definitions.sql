-- +goose Up
DELETE FROM channel_definitions;
ALTER TABLE channel_definitions DROP CONSTRAINT channel_definitions_pkey;
ALTER TABLE channel_definitions ADD COLUMN don_id bigint, ADD COLUMN version bigint;
ALTER TABLE channel_definitions RENAME COLUMN evm_chain_id TO chain_selector;
ALTER TABLE channel_definitions ALTER COLUMN chain_selector TYPE NUMERIC(20, 0);
ALTER TABLE channel_definitions ADD PRIMARY KEY (chain_selector, addr, don_id);

-- +goose Down
ALTER TABLE channel_definitions DROP COLUMN don_id, DROP COLUMN version;
ALTER TABLE channel_definitions RENAME COLUMN chain_selector TO evm_chain_id;
ALTER TABLE channel_definitions ALTER COLUMN evm_chain_id TYPE bigint;
ALTER TABLE channel_definitions ADD PRIMARY KEY (evm_chain_id, addr);
