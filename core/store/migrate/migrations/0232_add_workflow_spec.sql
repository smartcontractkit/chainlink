-- +goose Up
-- +goose StatementBegin
CREATE TABLE workflow_specs (
    id              SERIAL PRIMARY KEY,
    workflow_id     varchar(64) NOT NULL,
    workflow        text NOT NULL,
    workflow_owner  varchar(40) NOT NULL,
    created_at      timestamp with time zone NOT NULL,
    updated_at      timestamp with time zone NOT NULL
);

ALTER TABLE jobs
    ADD COLUMN workflow_spec_id INT REFERENCES workflow_specs (id),
    DROP CONSTRAINT chk_specs;

ALTER TABLE jobs
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
        CASE "type" WHEN 'stream' THEN 1 ELSE NULL END, -- 'stream' type lacks a spec but should not cause validation to fail
        CASE "type" WHEN 'workflow' THEN 1 ELSE NULL END -- 'workflow' type lacks a spec but should not cause validation to fail
      ) = 1
    );

ALTER TABLE jobs
    DROP COLUMN workflow_spec_id;

DROP TABLE workflow_specs;
-- +goose StatementEnd
