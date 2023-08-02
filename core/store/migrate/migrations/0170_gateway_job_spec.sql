-- +goose Up
-- +goose StatementBegin
CREATE TABLE gateway_specs (
    id              SERIAL PRIMARY KEY,
    gateway_config  JSONB NOT NULL DEFAULT '{}',
    created_at      timestamp with time zone NOT NULL,
    updated_at      timestamp with time zone NOT NULL
);

ALTER TABLE jobs
    ADD COLUMN gateway_spec_id INT REFERENCES gateway_specs (id),
    DROP CONSTRAINT chk_only_one_spec,
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
                    gateway_spec_id) = 1);
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
ALTER TABLE jobs
    DROP CONSTRAINT chk_only_one_spec,
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
                    bootstrap_spec_id) = 1);

ALTER TABLE jobs
    DROP COLUMN gateway_spec_id;

DROP TABLE gateway_specs;
-- +goose StatementEnd