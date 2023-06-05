-- +goose Up
-- +goose StatementBegin
ALTER TABLE jobs
    ADD COLUMN type_spec JSONB NOT NULL DEFAULT '{}',
    DROP COLUMN gateway_spec_id,
    DROP COLUMN block_header_feeder_spec_id,
    DROP COLUMN blockhash_store_spec_id,
    DROP COLUMN bootstrap_spec_id,
    DROP COLUMN cron_spec_id,
    DROP COLUMN direct_request_spec_id,
    DROP COLUMN flux_monitor_spec_id,
    DROP COLUMN keeper_spec_id,
    DROP COLUMN legacy_gas_station_server_spec_id,
    DROP COLUMN legacy_gas_station_sidecar_spec_id,
    DROP COLUMN ocr_oracle_spec_id,
    DROP COLUMN ocr2_oracle_spec_id,
    DROP COLUMN vrf_spec_id,
    DROP COLUMN webhook_spec_id;

-- TODO: Remove this if unnecessary. ALTER TABLE jobs DROP CONSTRAINT chk_only_one_spec;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE jobs
    DROP COLUMN type_spec,
    ADD COLUMN ocr_oracle_spec_id INT REFERENCES ocr_oracle_specs (id),
    ADD COLUMN ocr2_oracle_spec_id INT REFERENCES ocr2_oracle_specs (id),
    ADD COLUMN direct_request_spec_id INT REFERENCES direct_request_specs (id),
    ADD COLUMN flux_monitor_spec_id INT REFERENCES flux_monitor_specs (id),
    ADD COLUMN keeper_spec_id INT REFERENCES keeper_specs (id),
    ADD COLUMN cron_spec_id INT REFERENCES cron_specs (id),
    ADD COLUMN vrf_spec_id INT REFERENCES vrf_specs (id),
    ADD COLUMN webhook_spec_id INT REFERENCES webhook_specs (id),
    ADD COLUMN blockhash_store_spec_id INT REFERENCES blockhash_store_specs (id),
    ADD COLUMN bootstrap_spec_id INT REFERENCES bootstrap_specs (id),
    ADD COLUMN block_header_feeder_spec_id INT REFERENCES block_header_feeder_specs (id),
    ADD COLUMN  gateway_spec_id INT REFERENCES  gateway_specs (id),
    ADD COLUMN legacy_gas_station_server_spec_id INT REFERENCES legacy_gas_station_server_specs (id),
    ADD COLUMN legacy_gas_station_sidecar_spec_id INT REFERENCES legacy_gas_station_sidecar_specs (id),
    ADD CONSTRAINT chk_only_one_spec CHECK (
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
-- +goose StatementEnd
