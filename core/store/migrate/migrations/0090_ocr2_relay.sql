-- +goose Up
-- +goose StatementBegin
ALTER TABLE offchainreporting2_oracle_specs
    ADD COLUMN relay text,
    ADD COLUMN relay_config JSONB;
ALTER TABLE offchainreporting2_oracle_specs
    RENAME COLUMN contract_address TO contract_id;
ALTER TABLE offchainreporting2_oracle_specs
    RENAME COLUMN transmitter_address TO transmitter_id;
ALTER TABLE offchainreporting2_oracle_specs
    RENAME COLUMN key_bundle_id TO ocr_key_bundle_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE  offchainreporting2_oracle_specs
    DROP COLUMN relay,
    DROP COLUMN relay_config;
ALTER TABLE  offchainreporting2_oracle_specs
    RENAME COLUMN contract_id TO contract_address;
ALTER TABLE offchainreporting2_oracle_specs
    RENAME COLUMN transmitter_id TO transmitter_address;
ALTER TABLE offchainreporting2_oracle_specs
    RENAME COLUMN ocr_key_bundle_id TO key_bundle_id;
-- +goose StatementEnd
