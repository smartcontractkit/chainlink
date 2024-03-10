-- +goose Up
ALTER TABLE ocr2_oracle_specs ADD CONSTRAINT offchainreporting2_oracle_specs_unique_contract_addr UNIQUE (contract_id);
-- +goose Down
ALTER TABLE ocr2_oracle_specs DROP CONSTRAINT offchainreporting2_oracle_specs_unique_contract_addr;
