-- +goose Up
CREATE TABLE block_header_feeder_specs
(
    id                      BIGSERIAL PRIMARY KEY,
    coordinator_v1_address  bytea DEFAULT NULL,
    coordinator_v2_address  bytea DEFAULT NULL,
    wait_blocks             bigint                   NOT NULL,
    lookback_blocks         bigint                   NOT NULL,
    blockhash_store_address bytea                    NOT NULL,
    batch_blockhash_store_address bytea                    NOT NULL,
    poll_period             bigint                   NOT NULL,
    run_timeout             bigint                   NOT NULL,
    evm_chain_id            numeric(78)
        REFERENCES evm_chains
                DEFERRABLE,
    from_addresses          bytea[]                  DEFAULT '{}' NOT NULL,
    get_blockhashes_batch_size    integer                   NOT NULL,
    store_blockhashes_batch_size  integer                   NOT NULL,
    created_at              timestamp with time zone NOT NULL,
    updated_at              timestamp with time zone NOT NULL
        CONSTRAINT coordinator_v1_address_len_chk CHECK (octet_length(coordinator_v1_address) = 20)
        CONSTRAINT coordinator_v2_address_len_chk CHECK (octet_length(coordinator_v2_address) = 20)
        CONSTRAINT blockhash_store_address_len_chk CHECK (octet_length(blockhash_store_address) = 20)
        CONSTRAINT batch_blockhash_store_address_len_chk CHECK (octet_length(batch_blockhash_store_address) = 20)
        CONSTRAINT at_least_one_coordinator_chk CHECK (coordinator_v1_address IS NOT NULL OR coordinator_v2_address IS NOT NULL)
);

ALTER TABLE jobs
    ADD COLUMN block_header_feeder_spec_id INT REFERENCES block_header_feeder_specs (id),
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
    DROP COLUMN block_header_feeder_spec_id;
DROP TABLE IF EXISTS block_header_feeder_specs;