-- +goose Up
CREATE TABLE eal_specs (
  id BIGSERIAL PRIMARY KEY,
  forwarder_address BYTEA NOT NULL,
  evm_chain_id NUMERIC(78) NOT NULL,
  from_addresses BYTEA[] DEFAULT '{}' NOT NULL,
  lookback_blocks BIGINT NOT NULL,
  poll_period BIGINT NOT NULL,
  run_timeout BIGINT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
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

CREATE TABLE eal_txs (
    request_id TEXT PRIMARY KEY,
    forwarder_address BYTEA NOT NULL,
    from_address BYTEA NOT NULL,
    target_address BYTEA NOT NULL,
    evm_chain_id NUMERIC(78) NOT NULL,
    payload BYTEA NOT NULL,
    tx_status TEXT NOT NULL,
    gas_limit BIGINT NOT NULL,
    ccip_message_id BYTEA,
    failure_reason TEXT,
    status_update_url TEXT,
    tx_hash BYTEA,
    tx_id BIGINT REFERENCES txes INITIALLY DEFERRED,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,
    CONSTRAINT forwarder_address_len_chk CHECK (
        octet_length(forwarder_address) = 20
    ),
    CONSTRAINT target_address_len_chk CHECK (
        octet_length(target_address) = 20
    ),
    CONSTRAINT from_address_len_chk CHECK (
        octet_length(from_address) = 20
    ),
    CONSTRAINT ccip_message_id_len_chk CHECK (
        octet_length(ccip_message_id) = 32
    ),
    CONSTRAINT tx_hash_len_chk CHECK (
        octet_length(tx_hash) = 32
    )
);
CREATE INDEX idx_eal_txs_chain_id_tx_status ON eal_txs(evm_chain_id, tx_status);

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
DROP
  TABLE IF EXISTS eal_txs;  
