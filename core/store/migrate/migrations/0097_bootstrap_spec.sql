-- +goose Up
-- +goose StatementBegin

CREATE TABLE bootstrap_specs
(
    id                                    SERIAL PRIMARY KEY,
    contract_id                           text                     NOT NULL,
    relay                                 text,
    relay_config                          JSONB,
    monitoring_endpoint                   text,
    blockchain_timeout                    bigint,
    contract_config_tracker_poll_interval bigint,
    contract_config_confirmations         integer                  NOT NULL,
    created_at                            timestamp with time zone NOT NULL,
    updated_at                            timestamp with time zone NOT NULL
);

ALTER TABLE jobs
    ADD COLUMN bootstrap_spec_id INT REFERENCES bootstrap_specs (id),
    DROP CONSTRAINT chk_only_one_spec,
    ADD CONSTRAINT chk_only_one_spec CHECK (
            num_nonnulls(
                    offchainreporting_oracle_spec_id,
                    offchainreporting2_oracle_spec_id,
                    direct_request_spec_id,
                    flux_monitor_spec_id,
                    keeper_spec_id,
                    cron_spec_id,
                    webhook_spec_id,
                    vrf_spec_id,
                    blockhash_store_spec_id,
                    bootstrap_spec_id) = 1
        );

ALTER TABLE offchainreporting2_oracle_specs
    DROP COLUMN contract_config_tracker_subscribe_interval;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE jobs
    DROP CONSTRAINT chk_only_one_spec,
    ADD CONSTRAINT chk_only_one_spec CHECK (
            num_nonnulls(
                    offchainreporting_oracle_spec_id,
                    offchainreporting2_oracle_spec_id,
                    direct_request_spec_id,
                    flux_monitor_spec_id,
                    keeper_spec_id,
                    cron_spec_id,
                    webhook_spec_id,
                    vrf_spec_id,
                    blockhash_store_spec_id) = 1
        );
ALTER TABLE jobs
    DROP COLUMN bootstrap_spec_id;
DROP TABLE IF EXISTS bootstrap_specs;

ALTER TABLE offchainreporting2_oracle_specs
    ADD COLUMN contract_config_tracker_subscribe_interval bigint;
-- +goose StatementEnd
