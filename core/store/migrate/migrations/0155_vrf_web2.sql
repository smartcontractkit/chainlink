-- +goose Up

CREATE TABLE vrf_web2_requests(
    client_request_id bytea NOT NULL,
    lottery_type INTEGER NOT NULL,
    vrf_external_request_id bytea NOT NULL, -- UUID
    lottery_contract_address bytea CHECK (octet_length(lottery_contract_address) = 20) NOT NULL,
    evm_chain_id numeric(78,0) NOT NULL REFERENCES evm_chains (id) DEFERRABLE,
    request_tx_hash bytea NOT NULL CHECK (octet_length(request_tx_hash) = 32),
    PRIMARY KEY (
        client_request_id,
        lottery_type,
        vrf_external_request_id,
        lottery_contract_address,
        evm_chain_id
    )
);

CREATE TABLE vrf_web2_fulfillments(
    client_request_id TEXT NOT NULL,
    lottery_type INTEGER NOT NULL,
    vrf_external_request_id TEXT NOT NULL, -- UUID
    lottery_contract_address bytea CHECK (octet_length(lottery_contract_address) = 20) NOT NULL,
    evm_chain_id numeric(78,0) NOT NULL REFERENCES evm_chains (id) DEFERRABLE,
    winning_numbers integer[] NOT NULL,
    fulfillment_tx_hash bytea NOT NULL CHECK (octet_length(fulfillment_tx_hash) = 32),
    PRIMARY KEY (
        client_request_id,
        lottery_type,
        vrf_external_request_id,
        lottery_contract_address,
        evm_chain_id
    )
);

CREATE TABLE vrf_web2_specs
(
    id BIGSERIAL PRIMARY KEY,
    lottery_consumer_address bytea NOT NULL,
    evm_chain_id numeric(78) REFERENCES evm_chains DEFERRABLE,
    from_addresses bytea[] NOT NULL,
    created_at              timestamp with time zone NOT NULL,
    updated_at              timestamp with time zone NOT NULL
    CONSTRAINT lottery_consumer_address CHECK (octet_length(lottery_consumer_address) = 20)
);

ALTER TABLE jobs
    ADD COLUMN vrf_web2_spec_id INT REFERENCES vrf_web2_specs (id),
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
                    bootstrap_spec_id,
                    vrf_web2_spec_id) = 1
        );

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
                    bootstrap_spec_id) = 1
        );
ALTER TABLE jobs
    DROP COLUMN vrf_web2_spec_id;

DROP TABLE vrf_web2_specs;
DROP TABLE vrf_web2_requests;
DROP TABLE vrf_web2_fulfillments;
