-- +goose Up
-- +goose StatementBegin
CREATE TABLE bootstrap_contract_configs
(
    bootstrap_spec_id       INTEGER PRIMARY KEY,
    config_digest           bytea                    NOT NULL,
    config_count            bigint                   NOT NULL,
    signers                 bytea[]                  NOT NULL,
    transmitters            text[]                   NOT NULL,
    f                       smallint                 NOT NULL,
    onchain_config          bytea,
    offchain_config_version bigint                   NOT NULL,
    offchain_config         bytea,
    created_at              timestamp with time zone NOT NULL,
    updated_at              timestamp with time zone NOT NULL,
    CONSTRAINT bootstrap_contract_configs_config_digest_check CHECK ((octet_length(config_digest) = 32))
);

ALTER TABLE ONLY bootstrap_contract_configs
    ADD CONSTRAINT bootstrap_contract_configs_oracle_spec_fkey
        FOREIGN KEY (bootstrap_spec_id)
            REFERENCES bootstrap_specs (id)
            ON DELETE CASCADE;

-- add missing unique constraint for bootstrap specs
CREATE UNIQUE INDEX idx_jobs_unique_bootstrap_spec_id ON jobs USING btree (bootstrap_spec_id);

-- migrate existing OCR2 bootstrap jobs to the new bootstrap spec
-- create helper column
ALTER TABLE bootstrap_specs
    ADD COLUMN job_id INTEGER;

-- insert bootstrap specs
INSERT INTO bootstrap_specs (contract_id, relay, relay_config, monitoring_endpoint, blockchain_timeout,
                             contract_config_tracker_poll_interval, contract_config_confirmations, created_at,
                             updated_at, job_id)
SELECT ocr2.contract_id,
       ocr2.relay,
       ocr2.relay_config,
       ocr2.monitoring_endpoint,
       ocr2.blockchain_timeout,
       ocr2.contract_config_tracker_poll_interval,
       ocr2.contract_config_confirmations,
       ocr2.created_at,
       ocr2.updated_at,
       jobs.id
FROM jobs
         INNER JOIN offchainreporting2_oracle_specs AS ocr2 ON jobs.offchainreporting2_oracle_spec_id = ocr2.id
WHERE ocr2.is_bootstrap_peer IS true;

-- point jobs to new bootstrap specs
UPDATE jobs
SET type                              = 'bootstrap',
    offchainreporting2_oracle_spec_id = null,
    bootstrap_spec_id                 = (SELECT id FROM bootstrap_specs WHERE jobs.id = bootstrap_specs.job_id)
WHERE (SELECT COUNT(*) FROM bootstrap_specs WHERE jobs.id = bootstrap_specs.job_id) > 0;

-- cleanup
-- delete old ocr2 bootstrap specs
DELETE
FROM offchainreporting2_oracle_specs
WHERE is_bootstrap_peer IS true;

ALTER TABLE offchainreporting2_oracle_specs
    DROP COLUMN is_bootstrap_peer;
ALTER TABLE bootstrap_specs
    DROP COLUMN job_id;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE bootstrap_contract_configs;

-- create helper column
ALTER TABLE offchainreporting2_oracle_specs
    ADD COLUMN is_bootstrap_peer bool not null default false,
    ADD COLUMN job_id            INTEGER;

-- insert ocr2 specs
INSERT INTO offchainreporting2_oracle_specs (contract_id, is_bootstrap_peer, ocr_key_bundle_id, monitoring_endpoint,
                                             transmitter_id, blockchain_timeout, contract_config_tracker_poll_interval,
                                             contract_config_confirmations, juels_per_fee_coin_pipeline, created_at,
                                             updated_at, relay, relay_config, job_id)
SELECT bootstrap_specs.contract_id,
       true,
       null,
       bootstrap_specs.monitoring_endpoint,
       '',
       bootstrap_specs.blockchain_timeout,
       bootstrap_specs.contract_config_tracker_poll_interval,
       bootstrap_specs.contract_config_confirmations,
       '',
       bootstrap_specs.created_at,
       bootstrap_specs.updated_at,
       bootstrap_specs.relay,
       bootstrap_specs.relay_config,
       jobs.id
FROM jobs
         INNER JOIN bootstrap_specs ON jobs.bootstrap_spec_id = bootstrap_specs.id
WHERE jobs.bootstrap_spec_id is not null;

-- point jobs to new ocr2 specs
UPDATE jobs
SET type                              = 'offchainreporting2',
    bootstrap_spec_id                 = null,
    offchainreporting2_oracle_spec_id = (SELECT id
                                         FROM offchainreporting2_oracle_specs
                                         WHERE jobs.id = offchainreporting2_oracle_specs.job_id)
WHERE (SELECT COUNT(*) FROM offchainreporting2_oracle_specs WHERE jobs.id = offchainreporting2_oracle_specs.job_id) > 0;

-- cleanup
DELETE
FROM bootstrap_specs;

ALTER TABLE offchainreporting2_oracle_specs
    DROP COLUMN job_id;
-- +goose StatementEnd