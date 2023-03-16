-- +goose Up
CREATE TABLE transmission_specs
(
    id                      BIGSERIAL PRIMARY KEY,
    rpc_port                numeric(5)               NOT NULL,
    evm_chain_id            numeric(78)
        REFERENCES evm_chains
                DEFERRABLE,
    from_addresses          bytea[]                  DEFAULT '{}' NOT NULL,
    created_at              timestamp with time zone NOT NULL,
    updated_at              timestamp with time zone NOT NULL
);
ALTER TABLE jobs
    ADD COLUMN transmission_spec_id INT REFERENCES transmission_specs (id),
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
                    transmission_spec_id) = 1);


-- +goose Down
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
                    bootstrap_spec_id) = 1);

ALTER TABLE jobs
    DROP COLUMN transmission_spec_id;
DROP TABLE IF EXISTS transmission_specs;
