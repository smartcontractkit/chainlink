-- +goose Up
DROP TABLE evm_chains CASCADE;
DROP TABLE solana_chains CASCADE;
DROP TABLE starknet_chains CASCADE;

-- +goose Down
-- evm_chains definition
CREATE TABLE evm_chains (
    id numeric(78) NOT NULL,
    cfg jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    enabled bool NOT NULL DEFAULT true,
    CONSTRAINT evm_chains_pkey PRIMARY KEY (id)
);

-- evm_chains foreign keys
ALTER TABLE evm_log_poller_filters ADD CONSTRAINT evm_log_poller_filters_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id) DEFERRABLE;
ALTER TABLE evm_log_poller_blocks ADD CONSTRAINT evm_log_poller_blocks_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id) ON DELETE CASCADE DEFERRABLE NOT VALID;
ALTER TABLE log_broadcasts ADD CONSTRAINT log_broadcasts_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id) ON DELETE CASCADE DEFERRABLE NOT VALID;
ALTER TABLE block_header_feeder_specs ADD CONSTRAINT block_header_feeder_specs_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id) DEFERRABLE;
ALTER TABLE direct_request_specs ADD CONSTRAINT direct_request_specs_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id) ON DELETE CASCADE DEFERRABLE NOT VALID;
ALTER TABLE evm_logs ADD CONSTRAINT evm_logs_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id) ON DELETE CASCADE DEFERRABLE NOT VALID;
ALTER TABLE vrf_specs ADD CONSTRAINT vrf_specs_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id) ON DELETE CASCADE DEFERRABLE NOT VALID;
ALTER TABLE evm_heads ADD CONSTRAINT heads_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id) ON DELETE CASCADE DEFERRABLE NOT VALID;
ALTER TABLE evm_forwarders ADD CONSTRAINT evm_forwarders_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id) ON DELETE CASCADE DEFERRABLE NOT VALID;
ALTER TABLE blockhash_store_specs ADD CONSTRAINT blockhash_store_specs_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id) ON DELETE CASCADE DEFERRABLE NOT VALID;
ALTER TABLE evm_key_states ADD CONSTRAINT eth_key_states_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id) ON DELETE CASCADE DEFERRABLE NOT VALID;
ALTER TABLE log_broadcasts_pending ADD CONSTRAINT log_broadcasts_pending_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id) ON DELETE CASCADE DEFERRABLE NOT VALID;
ALTER TABLE eth_txes ADD CONSTRAINT eth_txes_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id) ON DELETE CASCADE DEFERRABLE NOT VALID;
ALTER TABLE keeper_specs ADD CONSTRAINT keeper_specs_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id) ON DELETE CASCADE DEFERRABLE NOT VALID;
ALTER TABLE flux_monitor_specs ADD CONSTRAINT flux_monitor_specs_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id) ON DELETE CASCADE DEFERRABLE NOT VALID;
ALTER TABLE ocr_oracle_specs ADD CONSTRAINT offchainreporting_oracle_specs_evm_chain_id_fkey FOREIGN KEY (evm_chain_id) REFERENCES evm_chains(id) ON DELETE CASCADE DEFERRABLE NOT VALID;

-- solana_chains definition
CREATE TABLE solana_chains (
    id text NOT NULL,
    cfg jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    enabled bool NOT NULL DEFAULT true,
    CONSTRAINT solana_chains_pkey PRIMARY KEY (id)
);

-- starknet_chains definition
CREATE TABLE starknet_chains (
    id text NOT NULL,
    cfg jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    enabled bool NOT NULL DEFAULT true,
    CONSTRAINT starknet_chains_pkey PRIMARY KEY (id)
);