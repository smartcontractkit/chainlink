-- +goose Up
-- +goose StatementBegin
ALTER TABLE jobs -- `type_spec` should be made NOT NULL after refactoring
ADD COLUMN type_spec JSONB;
ALTER TABLE jobs DROP CONSTRAINT chk_only_one_spec;
ALTER TABLE jobs DROP COLUMN bootstrap_spec_id;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE jobs DROP COLUMN type_spec;
ALTER TABLE jobs
ADD COLUMN bootstrap_spec_id INT REFERENCES bootstrap_specs (id);
ALTER TABLE jobs
ADD CONSTRAINT chk_only_one_spec CHECK (
        num_nonnulls(
            ocr_oracle_spec_id,
            ocr2_oracle_spec_id,
            direct_request_spec_id,
            flux_monitor_spec_id,
            keeper_spec_id,
            cron_spec_id,
            webhook_spec_id,
            vrf_spec_id,
            blockhash_store_spec_id,
            block_header_feeder_spec_id,
            bootstrap_spec_id,
            gateway_spec_id,
            legacy_gas_station_server_spec_id,
            legacy_gas_station_sidecar_spec_id
        ) = 1
    );
;
-- +goose StatementEnd
