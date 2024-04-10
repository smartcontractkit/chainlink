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

INSERT INTO workflow_specs (workflow_id, workflow, workflow_owner, created_at, updated_at)
    VALUES (
	'15c631d295ef5e32deb99a10ee6804bc4af1385568f9b3363f6552ac6dbb2cef',
	'triggers:
  - type: "mercury-trigger"
    config:
      feedIds:
        - "0x1111111111111111111100000000000000000000000000000000000000000000"
        - "0x2222222222222222222200000000000000000000000000000000000000000000"
        - "0x3333333333333333333300000000000000000000000000000000000000000000"

consensus:
  - type: "offchain_reporting"
    ref: "evm_median"
    inputs:
      observations:
        - "$(trigger.outputs)"
    config:
      aggregation_method: "data_feeds_2_0"
      aggregation_config:
        "0x1111111111111111111100000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: 3600
        "0x2222222222222222222200000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: 3600
        "0x3333333333333333333300000000000000000000000000000000000000000000":
          deviation: "0.001"
          heartbeat: 3600
      encoder: "EVM"
      encoder_config:
        abi: "mercury_reports bytes[]"

targets:
  - type: "write_polygon-testnet-mumbai"
    inputs:
      report: "$(evm_median.outputs.report)"
    config:
      address: "0x3F3554832c636721F1fD1822Ccca0354576741Ef"
      params: ["$(report)"]
      abi: "receive(report bytes)"
  - type: "write_ethereum-testnet-sepolia"
    inputs:
      report: "$(evm_median.outputs.report)"
    config:
      address: "0x54e220867af6683aE6DcBF535B4f952cB5116510"
      params: ["$(report)"]
      abi: "receive(report bytes)"',
        '00000000000000000000000000000000000000aa',
	NOW(),
	NOW()
    );

ALTER TABLE jobs
    ADD COLUMN workflow_spec_id INT REFERENCES workflow_specs (id),
    DROP CONSTRAINT chk_specs;

UPDATE jobs SET workflow_spec_id = (select id from workflow_specs limit 1) where type = 'workflow';

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
