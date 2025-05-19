-- +goose Up
CREATE TABLE liquidity_balancer_specs (
    id BIGSERIAL PRIMARY KEY,
    liquidity_balancer_config JSONB NOT NULL
);

ALTER TABLE
    jobs
ADD COLUMN
    liquidity_balancer_spec_id BIGINT REFERENCES liquidity_balancer_specs(id),
DROP CONSTRAINT chk_only_one_spec,
ADD CONSTRAINT chk_only_one_spec CHECK (
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
      liquidity_balancer_spec_id
    ) = 1
);

-- +goose Down
ALTER TABLE
    jobs
DROP CONSTRAINT chk_only_one_spec,
ADD CONSTRAINT chk_only_one_spec CHECK (
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
ALTER TABLE
    jobs
DROP COLUMN
    liquidity_balancer_spec_id;
DROP TABLE
    liquidity_balancer_specs;
