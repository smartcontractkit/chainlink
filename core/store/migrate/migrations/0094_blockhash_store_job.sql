-- +goose Up
CREATE TABLE blockhash_store_specs
(
    id                      BIGSERIAL PRIMARY KEY,
    coordinator_v1_address  bytea DEFAULT NULL,
    coordinator_v2_address  bytea DEFAULT NULL,
    wait_blocks             bigint                   NOT NULL,
    lookback_blocks         bigint                   NOT NULL,
    blockhash_store_address bytea                    NOT NULL,
    poll_period             bigint                   NOT NULL,
    run_timeout             bigint                   NOT NULL,
    evm_chain_id            numeric(78)
        REFERENCES evm_chains
            DEFERRABLE,
    from_address            bytea                    DEFAULT NULL,
    created_at              timestamp with time zone NOT NULL,
    updated_at              timestamp with time zone NOT NULL
        CONSTRAINT coordinator_v1_address_len_chk CHECK (octet_length(coordinator_v1_address) = 20)
        CONSTRAINT coordinator_v2_address_len_chk CHECK (octet_length(coordinator_v2_address) = 20)
        CONSTRAINT blockhash_store_address_len_chk CHECK (octet_length(blockhash_store_address) = 20)
        CONSTRAINT at_least_one_coordinator_chk CHECK (coordinator_v1_address IS NOT NULL OR coordinator_v2_address IS NOT NULL)
);
ALTER TABLE jobs
    ADD COLUMN blockhash_store_spec_id INT REFERENCES blockhash_store_specs (id),
    DROP CONSTRAINT chk_only_one_spec,
    ADD CONSTRAINT chk_only_one_spec CHECK (
            num_nonnulls(
                    offchainreporting_oracle_spec_id,
                    offchainreporting2_oracle_spec_id,
                    direct_request_spec_id,
                    flux_monitor_spec_id,
                    keeper_spec_id,
                    cron_spec_id,
                    vrf_spec_id,
                    webhook_spec_id,
                    blockhash_store_spec_id) = 1);

-- +goose Down
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
                    vrf_spec_id,
                    webhook_spec_id) = 1);

ALTER TABLE jobs
    DROP COLUMN blockhash_store_spec_id;
DROP TABLE IF EXISTS blockhash_store_specs;
