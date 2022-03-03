-- +goose Up
-- +goose StatementBegin
ALTER TABLE offchainreporting2_oracle_specs
    ADD COLUMN relay text NOT NULL,
    ADD COLUMN relay_config JSONB NOT NULL DEFAULT '{}',
    ALTER COLUMN contract_address TYPE text,
    ALTER COLUMN transmitter_address TYPE text,
    DROP COLUMN evm_chain_id,
    DROP CONSTRAINT chk_contract_address_length;
ALTER TABLE offchainreporting2_oracle_specs
    RENAME COLUMN contract_address TO contract_id;
ALTER TABLE offchainreporting2_oracle_specs
    RENAME COLUMN transmitter_address TO transmitter_id;
ALTER TABLE offchainreporting2_oracle_specs
    RENAME COLUMN encrypted_ocr_key_bundle_id TO ocr_key_bundle_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE  offchainreporting2_oracle_specs
    DROP COLUMN relay,
    DROP COLUMN relay_config,
    ADD COLUMN evm_chain_id numeric(78,0) REFERENCES evm_chains(id);
ALTER TABLE  offchainreporting2_oracle_specs
    RENAME COLUMN contract_id TO contract_address;
ALTER TABLE offchainreporting2_oracle_specs
    RENAME COLUMN transmitter_id TO transmitter_address;
ALTER TABLE offchainreporting2_oracle_specs
    RENAME COLUMN ocr_key_bundle_id TO encrypted_ocr_key_bundle_id;
-- +goose StatementEnd
