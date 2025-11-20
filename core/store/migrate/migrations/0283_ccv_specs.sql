-- +goose Up
-- +goose StatementBegin

CREATE TABLE ccv_committee_verifier_specs(
    id BIGSERIAL PRIMARY KEY,
    committee_verifier_config TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE ccv_executor_specs(
    id BIGSERIAL PRIMARY KEY,
    executor_config TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

ALTER TABLE jobs
    ADD COLUMN ccv_committee_verifier_spec_id INT REFERENCES ccv_committee_verifier_specs (id),
    ADD COLUMN ccv_executor_spec_id INT REFERENCES ccv_executor_specs (id),
DROP CONSTRAINT chk_specs,
    ADD CONSTRAINT chk_specs CHECK (
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
        END -- 'stream' type lacks a spec but should not cause validation to fail
      ) = 1
    );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE jobs
DROP CONSTRAINT chk_specs,
     ADD CONSTRAINT chk_specs CHECK (
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
        CASE "type"
	  WHEN 'stream'
	  THEN 1
	  ELSE NULL
        END -- 'stream' type lacks a spec but should not cause validation to fail
      ) = 1
    );

ALTER TABLE jobs
DROP COLUMN ccv_committee_verifier_spec_id;

ALTER TABLE jobs
DROP COLUMN ccv_executor_spec_id;

DROP TABLE ccv_committee_verifier_specs;
DROP TABLE ccv_executor_specs;
-- +goose StatementEnd
