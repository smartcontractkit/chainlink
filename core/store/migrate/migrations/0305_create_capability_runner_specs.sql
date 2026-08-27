-- +goose Up
-- +goose StatementBegin

CREATE TABLE capability_runner_specs(
    id BIGSERIAL PRIMARY KEY,

    command TEXT NOT NULL,
    args TEXT[] NOT NULL DEFAULT '{}',

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

ALTER TABLE jobs
    ADD COLUMN capability_runner_spec_id INT REFERENCES capability_runner_specs (id),
DROP CONSTRAINT chk_specs,
    ADD CONSTRAINT chk_specs CHECK (
      num_nonnulls(
        ocr_oracle_spec_id, ocr2_oracle_spec_id,
        direct_request_spec_id, flux_monitor_spec_id,
        cron_spec_id, webhook_spec_id,
        vrf_spec_id, blockhash_store_spec_id,
        block_header_feeder_spec_id, bootstrap_spec_id,
        gateway_spec_id,
        eal_spec_id,
        workflow_spec_id,
        standard_capabilities_spec_id,
        ccip_spec_id,
        ccip_bootstrap_spec_id,
        cre_settings_spec_id,
        ccv_committee_verifier_spec_id,
        ccv_executor_spec_id,
        capability_runner_spec_id,
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
ALTER TABLE jobs
DROP CONSTRAINT chk_specs,
    ADD CONSTRAINT chk_specs CHECK (
      num_nonnulls(
        ocr_oracle_spec_id, ocr2_oracle_spec_id,
        direct_request_spec_id, flux_monitor_spec_id,
        cron_spec_id, webhook_spec_id,
        vrf_spec_id, blockhash_store_spec_id,
        block_header_feeder_spec_id, bootstrap_spec_id,
        gateway_spec_id,
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

ALTER TABLE jobs
DROP COLUMN capability_runner_spec_id;

DROP TABLE capability_runner_specs;

-- +goose StatementEnd
