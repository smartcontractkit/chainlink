-- +goose Up
CREATE TABLE legacy_gas_station_server_specs (
  id BIGSERIAL PRIMARY KEY,
  forwarder_address BYTEA NOT NULL,
  evm_chain_id numeric(78) NOT NULL REFERENCES evm_chains DEFERRABLE,
  ccip_chain_selector numeric(78) NOT NULL,
  from_addresses BYTEA[] DEFAULT '{}' NOT NULL,
  created_at timestamp WITH TIME ZONE NOT NULL,
  updated_at timestamp WITH TIME ZONE NOT NULL,
  CONSTRAINT forwarder_address_len_chk CHECK (
    octet_length(forwarder_address) = 20
  )
);
CREATE TABLE legacy_gas_station_sidecar_specs (
  id BIGSERIAL PRIMARY KEY,
  forwarder_address BYTEA NOT NULL,
  off_ramp_address BYTEA NOT NULL,
  lookback_blocks bigint NOT NULL,
  poll_period bigint NOT NULL,
  run_timeout bigint NOT NULL,
  evm_chain_id numeric(78) NOT NULL REFERENCES evm_chains DEFERRABLE,
  ccip_chain_selector numeric(78) NOT NULL,
  created_at timestamp WITH TIME ZONE NOT NULL,
  updated_at timestamp WITH TIME ZONE NOT NULL,
  CONSTRAINT forwarder_address_len_chk CHECK (
    octet_length(forwarder_address) = 20
  ),
  CONSTRAINT off_ramp_address_len_chk CHECK (
    octet_length(off_ramp_address) = 20
  )
);
CREATE TABLE legacy_gasless_txs (
  legacy_gasless_tx_id TEXT PRIMARY KEY,
  forwarder_address BYTEA NOT NULL,
  from_address BYTEA NOT NULL,
  target_address BYTEA NOT NULL,
  receiver_address BYTEA NOT NULL,
  nonce numeric(78) NOT NULL,
  amount numeric(78) NOT NULL,
  source_chain_id numeric(78) NOT NULL,
  destination_chain_id numeric(78) NOT NULL,
  valid_until_time numeric(78) NOT NULL,
  tx_signature BYTEA NOT NULL,
  tx_status text NOT NULL,
  token_name text NOT NULL,
  token_version text NOT NULL,
  eth_tx_id bigint REFERENCES eth_txes INITIALLY DEFERRED,
  ccip_message_id BYTEA,
  failure_reason text,
  created_at timestamp WITH TIME ZONE NOT NULL,
  updated_at timestamp WITH TIME ZONE NOT NULL,
  CONSTRAINT forwarder_address_len_chk CHECK (
    octet_length(forwarder_address) = 20
  ),
  CONSTRAINT target_address_len_chk CHECK (
    octet_length(target_address) = 20
  ),
  CONSTRAINT receiver_address_len_chk CHECK (
    octet_length(receiver_address) = 20
  ),
  CONSTRAINT ccip_message_id_len_chk CHECK (
    octet_length(ccip_message_id) = 32
  )
);
CREATE INDEX idx_legacy_gasless_txs_source_chain_id_tx_status ON legacy_gasless_txs(source_chain_id, tx_status);
CREATE INDEX idx_legacy_gasless_txs_source_destination_id_tx_status ON legacy_gasless_txs(destination_chain_id, tx_status);
ALTER TABLE
  jobs
ADD
  COLUMN legacy_gas_station_server_spec_id INT REFERENCES legacy_gas_station_server_specs (id),
ADD
  COLUMN legacy_gas_station_sidecar_spec_id INT REFERENCES legacy_gas_station_sidecar_specs (id),
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
      block_header_feeder_spec_id, bootstrap_spec_id
    ) = 1
  );
ALTER TABLE
  jobs
DROP
  COLUMN legacy_gas_station_server_spec_id,
DROP
  COLUMN legacy_gas_station_sidecar_spec_id;
DROP
  TABLE IF EXISTS legacy_gas_station_server_specs;
DROP
  TABLE IF EXISTS legacy_gas_station_sidecar_specs;
DROP
  TABLE IF EXISTS legacy_gasless_txs;
