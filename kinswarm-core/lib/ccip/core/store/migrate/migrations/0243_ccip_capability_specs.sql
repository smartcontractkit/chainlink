-- +goose Up
-- +goose StatementBegin

-- The ccip_specs table will hold the CCIP capability job specs.
-- For each new CCIP capability version, we will create a new CCIP capability job.
-- A single CCIP capability job manages all CCIP OCR instances across all chains.
CREATE TABLE ccip_specs(
    id BIGSERIAL PRIMARY KEY,

    -- The CCIP capability version, specified in the capability registry.
    capability_version TEXT NOT NULL,

    -- The CCIP capability labelled name, specified in the capability registry.
    capability_labelled_name TEXT NOT NULL,

    -- A mapping of chain family to OCR key bundle ID.
    -- Every chain family will have its own OCR key bundle.
    ocr_key_bundle_ids JSONB NOT NULL,

    -- The P2P ID for the node.
    -- The same P2P ID can be used across many chain families and OCR DONs.
    p2p_key_id TEXT NOT NULL,

    -- The P2P V2 bootstrappers, used to bootstrap the DON network.
    -- These are of the form nodeP2PID@nodeIP:nodePort.
    p2pv2_bootstrappers TEXT[] NOT NULL,

    -- A mapping of chain family to relay configuration for that family.
    -- Relay configuration typically consists of contract reader and contract writer
    -- configurations.
    relay_configs JSONB NOT NULL,

    -- A mapping of ccip plugin type to plugin configuration for that plugin.
    -- For example, the token price pipeline can live in the plugin config of
    -- the commit plugin.
    plugin_config JSONB NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- The ccip_bootstrap_specs table will hold the CCIP capability bootstrap job specs.
-- Similar to the CCIP capability job specs, these specs are scoped to a single CCIP
-- capability version.
-- A single CCIP bootstrap job will be able to bootstrap all CCIP OCR instances across all chains.
CREATE TABLE ccip_bootstrap_specs(
  id BIGSERIAL PRIMARY KEY,

  -- The CCIP capability version, specified in the capability registry.
    capability_version TEXT NOT NULL,

  -- The CCIP capability labelled name, specified in the capability registry.
  capability_labelled_name TEXT NOT NULL,

  -- Relay config of the home chain.
  relay_config JSONB NOT NULL,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

ALTER TABLE jobs
    ADD COLUMN ccip_spec_id INT REFERENCES ccip_specs (id),
    ADD COLUMN ccip_bootstrap_spec_id INT REFERENCES ccip_bootstrap_specs (id),
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
        CASE "type"
	  WHEN 'stream'
	  THEN 1
	  ELSE NULL
        END -- 'stream' type lacks a spec but should not cause validation to fail
      ) = 1
    );

ALTER TABLE jobs
DROP COLUMN ccip_spec_id;

ALTER TABLE jobs
DROP COLUMN ccip_bootstrap_spec_id;

DROP TABLE ccip_specs;
DROP TABLE ccip_bootstrap_specs;
-- +goose StatementEnd
