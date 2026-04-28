-- +goose Up
-- +goose StatementBegin

DELETE FROM jobs WHERE type = 'keeper';

DROP INDEX IF EXISTS idx_jobs_unique_keeper_spec_id;

ALTER TABLE jobs DROP CONSTRAINT IF EXISTS jobs_keeper_spec_id_fkey;

ALTER TABLE jobs DROP COLUMN IF EXISTS keeper_spec_id;

DROP TABLE IF EXISTS upkeep_registrations;
DROP TABLE IF EXISTS keeper_registries;
DROP TABLE IF EXISTS keeper_specs;

ALTER TABLE jobs DROP CONSTRAINT IF EXISTS chk_specs;
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
-- Legacy keeper job support was removed; Down does not restore dropped tables or jobs.
SELECT 1;
-- +goose StatementEnd
