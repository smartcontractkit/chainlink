-- +goose Up
-- +goose StatementBegin

ALTER TABLE offchainreporting2_oracle_specs
    ADD COLUMN plugin_config JSONB NOT NULL DEFAULT '{}',
    ADD COLUMN plugin_type   text  NOT NULL default '';

-- migrate existing juels_per_fee_coin_pipeline settings to json format and set plugin_type to median as the only
-- plugins that are supported before this version are median plugins.
UPDATE offchainreporting2_oracle_specs
SET plugin_type   = 'median',
    plugin_config = jsonb_build_object('juelsPerFeeCoinSource', juels_per_fee_coin_pipeline);

ALTER TABLE offchainreporting2_oracle_specs
    DROP COLUMN juels_per_fee_coin_pipeline;

-- rename OCR2 tables
ALTER TABLE jobs
    RENAME COLUMN offchainreporting2_oracle_spec_id TO ocr2_oracle_spec_id;
ALTER TABLE offchainreporting2_oracle_specs
    RENAME TO ocr2_oracle_specs;

ALTER TABLE offchainreporting2_contract_configs
    RENAME TO ocr2_contract_configs;
ALTER TABLE ocr2_contract_configs
    RENAME COLUMN offchainreporting2_oracle_spec_id TO ocr2_oracle_spec_id;

ALTER TABLE offchainreporting2_latest_round_requested
    RENAME TO ocr2_latest_round_requested;
ALTER TABLE ocr2_latest_round_requested
    RENAME COLUMN offchainreporting2_oracle_spec_id TO ocr2_oracle_spec_id;

ALTER TABLE offchainreporting2_pending_transmissions
    RENAME TO ocr2_pending_transmissions;
ALTER TABLE ocr2_pending_transmissions
    RENAME COLUMN offchainreporting2_oracle_spec_id TO ocr2_oracle_spec_id;

ALTER TABLE offchainreporting2_persistent_states
    RENAME TO ocr2_persistent_states;
ALTER TABLE ocr2_persistent_states
    RENAME COLUMN offchainreporting2_oracle_spec_id TO ocr2_oracle_spec_id;

-- rename OCR tables
ALTER TABLE jobs
    RENAME COLUMN offchainreporting_oracle_spec_id TO ocr_oracle_spec_id;
ALTER TABLE offchainreporting_oracle_specs
    RENAME TO ocr_oracle_specs;

ALTER TABLE offchainreporting_contract_configs
    RENAME TO ocr_contract_configs;
ALTER TABLE ocr_contract_configs
    RENAME COLUMN offchainreporting_oracle_spec_id TO ocr_oracle_spec_id;

-- this table does not have offchainreporting_oracle_spec_id
ALTER TABLE offchainreporting_discoverer_announcements
    RENAME TO ocr_discoverer_announcements;

ALTER TABLE offchainreporting_latest_round_requested
    RENAME TO ocr_latest_round_requested;
ALTER TABLE ocr_latest_round_requested
    RENAME COLUMN offchainreporting_oracle_spec_id TO ocr_oracle_spec_id;

ALTER TABLE offchainreporting_pending_transmissions
    RENAME TO ocr_pending_transmissions;
ALTER TABLE ocr_pending_transmissions
    RENAME COLUMN offchainreporting_oracle_spec_id TO ocr_oracle_spec_id;

ALTER TABLE offchainreporting_persistent_states
    RENAME TO ocr_persistent_states;
ALTER TABLE ocr_persistent_states
    RENAME COLUMN offchainreporting_oracle_spec_id TO ocr_oracle_spec_id;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE ocr2_oracle_specs
    ADD COLUMN juels_per_fee_coin_pipeline text NOT NULL default '';

UPDATE ocr2_oracle_specs
SET juels_per_fee_coin_pipeline = plugin_config ->> 'juelsPerFeeCoinSource';

ALTER TABLE ocr2_oracle_specs
    DROP COLUMN plugin_config,
    DROP COLUMN plugin_type;

-- rename OCR2 tables
ALTER TABLE jobs
    RENAME COLUMN ocr2_oracle_spec_id TO offchainreporting2_oracle_spec_id;

ALTER TABLE ocr2_oracle_specs
    RENAME TO offchainreporting2_oracle_specs;

ALTER TABLE ocr2_contract_configs
    RENAME TO offchainreporting2_contract_configs;
ALTER TABLE offchainreporting2_contract_configs
    RENAME COLUMN ocr2_oracle_spec_id TO offchainreporting2_oracle_spec_id;

ALTER TABLE ocr2_latest_round_requested
    RENAME TO offchainreporting2_latest_round_requested;
ALTER TABLE offchainreporting2_latest_round_requested
    RENAME COLUMN ocr2_oracle_spec_id TO offchainreporting2_oracle_spec_id;

ALTER TABLE ocr2_pending_transmissions
    RENAME TO offchainreporting2_pending_transmissions;
ALTER TABLE offchainreporting2_pending_transmissions
    RENAME COLUMN ocr2_oracle_spec_id TO offchainreporting2_oracle_spec_id;

ALTER TABLE ocr2_persistent_states
    RENAME TO offchainreporting2_persistent_states;
ALTER TABLE offchainreporting2_persistent_states
    RENAME COLUMN ocr2_oracle_spec_id TO offchainreporting2_oracle_spec_id;

-- rename OCR tables
ALTER TABLE jobs
    RENAME COLUMN ocr_oracle_spec_id TO offchainreporting_oracle_spec_id;
ALTER TABLE ocr_oracle_specs
    RENAME TO offchainreporting_oracle_specs;

ALTER TABLE ocr_contract_configs
    RENAME TO offchainreporting_contract_configs;
ALTER TABLE offchainreporting_contract_configs
    RENAME COLUMN ocr_oracle_spec_id TO offchainreporting_oracle_spec_id;

ALTER TABLE ocr_discoverer_announcements
    RENAME TO offchainreporting_discoverer_announcements;

ALTER TABLE ocr_latest_round_requested
    RENAME TO offchainreporting_latest_round_requested;
ALTER TABLE offchainreporting_latest_round_requested
    RENAME COLUMN ocr_oracle_spec_id TO offchainreporting_oracle_spec_id;

ALTER TABLE ocr_pending_transmissions
    RENAME TO offchainreporting_pending_transmissions;
ALTER TABLE offchainreporting_pending_transmissions
    RENAME COLUMN ocr_oracle_spec_id TO offchainreporting_oracle_spec_id;

ALTER TABLE ocr_persistent_states
    RENAME TO offchainreporting_persistent_states;
ALTER TABLE offchainreporting_persistent_states
    RENAME COLUMN ocr_oracle_spec_id TO offchainreporting_oracle_spec_id;

-- +goose StatementEnd
