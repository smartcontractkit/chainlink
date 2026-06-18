-- +goose Up
-- +goose StatementBegin

DELETE FROM jobs WHERE type = 'keeper';

DROP INDEX IF EXISTS idx_jobs_unique_keeper_spec_id;

ALTER TABLE jobs DROP CONSTRAINT IF EXISTS jobs_keeper_spec_id_fkey;

-- chk_specs still references keeper_spec_id from prior migrations; drop it before the column.
ALTER TABLE jobs DROP CONSTRAINT IF EXISTS chk_specs;

ALTER TABLE jobs DROP COLUMN IF EXISTS keeper_spec_id;

DROP TABLE IF EXISTS upkeep_registrations;
DROP TABLE IF EXISTS keeper_registries;
DROP TABLE IF EXISTS keeper_specs;

ALTER TABLE jobs ADD CONSTRAINT chk_specs CHECK (
      num_nonnulls(
        ocr_oracle_spec_id, ocr2_oracle_spec_id,
        direct_request_spec_id, flux_monitor_spec_id,
        cron_spec_id, webhook_spec_id,
        vrf_spec_id, blockhash_store_spec_id,
        block_header_feeder_spec_id, bootstrap_spec_id,
        gateway_spec_id,
        legacy_gas_station_server_spec_id,
        legacy_gas_station_sidecar_spec_id,
        eal_spec_id,
        workflow_spec_id,
        standard_capabilities_spec_id,
        ccip_spec_id,
        ccip_bootstrap_spec_id,
        cre_settings_spec_id,
        ccv_committee_verifier_spec_id,
        ccv_executor_spec_id,
        CASE "type"
	  WHEN 'stream'
	  THEN 1
	  ELSE NULL
        END
      ) = 1
    );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

CREATE TABLE keeper_specs (
    id BIGSERIAL PRIMARY KEY,
    contract_address bytea NOT NULL,
    from_address bytea NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    min_incoming_confirmations integer,
    evm_chain_id numeric(78,0) NOT NULL,
    CONSTRAINT keeper_specs_contract_address_check CHECK ((octet_length(contract_address) = 20)),
    CONSTRAINT keeper_specs_from_address_check CHECK ((octet_length(from_address) = 20))
);

CREATE TABLE keeper_registries (
    id BIGSERIAL PRIMARY KEY,
    job_id int UNIQUE NOT NULL REFERENCES jobs (id) ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE,
    keeper_index int NOT NULL,
    contract_address bytea UNIQUE NOT NULL,
    from_address bytea NOT NULL,
    check_gas int NOT NULL,
    block_count_per_turn int NOT NULL,
    num_keepers int NOT NULL,
    keeper_index_map jsonb DEFAULT NULL,
    CONSTRAINT keeper_registries_contract_address_check CHECK ((octet_length(contract_address) = 20)),
    CONSTRAINT keeper_registries_from_address_check CHECK ((octet_length(from_address) = 20))
);

CREATE INDEX idx_keeper_registries_keeper_index ON keeper_registries (keeper_index);

CREATE TABLE upkeep_registrations (
    id BIGSERIAL PRIMARY KEY,
    registry_id bigint NOT NULL REFERENCES keeper_registries (id) ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE,
    execute_gas int NOT NULL,
    check_data bytea NOT NULL,
    upkeep_id numeric(78,0) NOT NULL,
    positioning_constant int NOT NULL,
    last_run_block_height bigint NOT NULL DEFAULT 0,
    last_keeper_index integer DEFAULT NULL
);

CREATE UNIQUE INDEX idx_upkeep_registrations_unique_upkeep_ids_per_keeper ON upkeep_registrations (registry_id, upkeep_id);
CREATE INDEX idx_upkeep_registrations_upkeep_id ON upkeep_registrations (upkeep_id);

ALTER TABLE jobs ADD COLUMN keeper_spec_id INT;

ALTER TABLE jobs DROP CONSTRAINT IF EXISTS chk_specs;
ALTER TABLE jobs ADD CONSTRAINT chk_specs CHECK (
      num_nonnulls(
        ocr_oracle_spec_id, ocr2_oracle_spec_id,
        direct_request_spec_id, flux_monitor_spec_id,
        keeper_spec_id, cron_spec_id, webhook_spec_id,
        vrf_spec_id, blockhash_store_spec_id,
        block_header_feeder_spec_id, bootstrap_spec_id,
        gateway_spec_id,
        legacy_gas_station_server_spec_id,
        legacy_gas_station_sidecar_spec_id,
        eal_spec_id,
        workflow_spec_id,
        standard_capabilities_spec_id,
        ccip_spec_id,
        ccip_bootstrap_spec_id,
        cre_settings_spec_id,
        ccv_committee_verifier_spec_id,
        ccv_executor_spec_id,
        CASE "type"
	  WHEN 'stream'
	  THEN 1
	  ELSE NULL
        END
      ) = 1
    );

ALTER TABLE jobs ADD CONSTRAINT jobs_keeper_spec_id_fkey FOREIGN KEY (keeper_spec_id) REFERENCES keeper_specs (id) ON DELETE CASCADE;

CREATE UNIQUE INDEX idx_jobs_unique_keeper_spec_id ON jobs (keeper_spec_id);

-- +goose StatementEnd
