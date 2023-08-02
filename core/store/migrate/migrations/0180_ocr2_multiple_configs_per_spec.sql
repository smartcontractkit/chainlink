-- +goose Up
ALTER TABLE ocr2_contract_configs
  ADD COLUMN plugin_id INTEGER NOT NULL DEFAULT 0,
  DROP CONSTRAINT offchainreporting2_contract_configs_pkey,
  ADD CONSTRAINT ocr2_contract_configs_unique_id_pair UNIQUE (ocr2_oracle_spec_id, plugin_id);

-- +goose Down
ALTER TABLE ocr2_contract_configs
  DROP CONSTRAINT ocr2_contract_configs_unique_id_pair,
  ADD CONSTRAINT offchainreporting2_contract_configs_pkey PRIMARY KEY (ocr2_oracle_spec_id),
  DROP COLUMN plugin_id;
