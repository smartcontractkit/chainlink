-- +goose Up
-- +goose StatementBegin
ALTER TABLE jobs -- `type_spec` should be made NOT NULL after refactoring
ADD COLUMN type_spec JSONB;
ALTER TABLE jobs DROP CONSTRAINT chk_only_one_spec;
ALTER TABLE jobs DROP COLUMN bootstrap_spec_id;
DROP TABLE bootstrap_contract_configs;
DROP TABLE bootstrap_specs;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE jobs DROP COLUMN type_spec;
CREATE TABLE bootstrap_specs (
    id SERIAL PRIMARY KEY,
    contract_id text NOT NULL,
    relay text,
    relay_config JSONB,
    monitoring_endpoint text,
    blockchain_timeout bigint,
    contract_config_tracker_poll_interval bigint,
    contract_config_confirmations integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    feed_id bytea CHECK (
        feed_id IS NULL
        OR octet_length(feed_id) = 32
    )
);
CREATE TABLE bootstrap_contract_configs (
    bootstrap_spec_id INTEGER PRIMARY KEY,
    config_digest bytea NOT NULL,
    config_count bigint NOT NULL,
    signers bytea [] NOT NULL,
    transmitters text [] NOT NULL,
    f smallint NOT NULL,
    onchain_config bytea,
    offchain_config_version bigint NOT NULL,
    offchain_config bytea,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    CONSTRAINT bootstrap_contract_configs_config_digest_check CHECK ((octet_length(config_digest) = 32))
);
ALTER TABLE ONLY bootstrap_contract_configs
ADD CONSTRAINT bootstrap_contract_configs_oracle_spec_fkey FOREIGN KEY (bootstrap_spec_id) REFERENCES bootstrap_specs (id) ON DELETE CASCADE;
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
