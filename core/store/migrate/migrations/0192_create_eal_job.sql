-- +goose Up
CREATE TABLE eal_specs (
  id BIGSERIAL PRIMARY KEY,
  forwarder_address BYTEA NOT NULL,
  evm_chain_id numeric(78) NOT NULL,
  ccip_chain_selector numeric(78) NOT NULL,
  from_addresses BYTEA[] DEFAULT '{}' NOT NULL,
  created_at timestamp WITH TIME ZONE NOT NULL,
  updated_at timestamp WITH TIME ZONE NOT NULL,
  CONSTRAINT forwarder_address_len_chk CHECK (
    octet_length(forwarder_address) = 20
  )
);

ALTER TABLE
  jobs
ADD
  COLUMN eal_spec_id INT REFERENCES eal_specs (id),
DROP
  CONSTRAINT chk_only_one_spec,
ADD
  CONSTRAINT chk_only_one_spec CHECK (
    num_nonnulls(
      ocr_oracle_spec_id, ocr2_oracle_spec_id,
      direct_request_spec_id, flux_monitor_spec_id,
      keeper_spec_id, cron_spec_id, webhook_spec_id,
      vrf_spec_id, blockhash_store_spec_id,
      block_header_feeder_spec_id, bootstrap_spec_id,
      gateway_spec_id,
      legacy_gas_station_server_spec_id,
      legacy_gas_station_sidecar_spec_id,
      eal_spec_id
    ) = 1
  );

-- +goose Down
ALTER TABLE
  jobs
DROP
  CONSTRAINT chk_only_one_spec,
ADD
  CONSTRAINT chk_only_one_spec CHECK (
    num_nonnulls(
      ocr_oracle_spec_id, ocr2_oracle_spec_id,
      direct_request_spec_id, flux_monitor_spec_id,
      keeper_spec_id, cron_spec_id, webhook_spec_id,
      vrf_spec_id, blockhash_store_spec_id,
      block_header_feeder_spec_id, bootstrap_spec_id,
      gateway_spec_id,
      legacy_gas_station_server_spec_id,
      legacy_gas_station_sidecar_spec_id
    ) = 1
  );
ALTER TABLE
  jobs
DROP
  COLUMN eal_spec_id;
DROP
  TABLE IF EXISTS eal_specs;