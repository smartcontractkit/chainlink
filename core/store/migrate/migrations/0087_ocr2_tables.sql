-- +goose Up
-- +goose StatementBegin
CREATE TABLE offchainreporting2_oracle_specs (
     id SERIAL PRIMARY KEY,
     contract_address bytea NOT NULL,
     p2p_peer_id text,
     p2p_bootstrap_peers text[] NOT NULL DEFAULT '{}',
     is_bootstrap_peer boolean NOT NULL,
     encrypted_ocr_key_bundle_id bytea,
     monitoring_endpoint text,
     transmitter_address bytea,
     blockchain_timeout bigint,
     evm_chain_id numeric(78,0) REFERENCES evm_chains (id),
     contract_config_tracker_subscribe_interval bigint,
     contract_config_tracker_poll_interval bigint,
     contract_config_confirmations integer NOT NULL,
     juels_per_fee_coin_pipeline text NOT NULL,
     created_at timestamp with time zone NOT NULL,
     updated_at timestamp with time zone NOT NULL,
     CONSTRAINT chk_contract_address_length CHECK ((octet_length(contract_address) = 20))
);

ALTER TABLE ONLY offchainreporting2_oracle_specs
    ADD CONSTRAINT offchainreporting2_oracle_specs_unique_contract_addr
        UNIQUE (contract_address);

CREATE INDEX idx_offchainreporting2_oracle_specs_created_at
    ON offchainreporting2_oracle_specs USING brin (created_at);
CREATE INDEX idx_offchainreporting2_oracle_specs_updated_at
    ON offchainreporting2_oracle_specs USING brin (updated_at);

ALTER TABLE jobs
    ADD COLUMN offchainreporting2_oracle_spec_id integer,
    ADD CONSTRAINT jobs_offchainreporting2_oracle_spec_id_fkey
        FOREIGN KEY (offchainreporting2_oracle_spec_id)
            REFERENCES offchainreporting2_oracle_specs(id)
            ON DELETE CASCADE,
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
                    webhook_spec_id
                ) = 1
        );

CREATE UNIQUE INDEX idx_jobs_unique_offchain2_reporting_oracle_spec_id
    ON jobs
        USING btree (offchainreporting2_oracle_spec_id);

CREATE TABLE offchainreporting2_contract_configs (
     offchainreporting2_oracle_spec_id INTEGER PRIMARY KEY,
     config_digest bytea NOT NULL,
     config_count bigint NOT NULL,
     signers bytea[],
     transmitters text[],
     f smallint NOT NULL,
     onchain_config bytea, -- this field exists in ocr2 but not in ocr1
     offchain_config_version bigint NOT NULL,
     offchain_config bytea,
     created_at timestamp with time zone NOT NULL,
     updated_at timestamp with time zone NOT NULL,
     CONSTRAINT offchainreporting2_contract_configs_config_digest_check CHECK ((octet_length(config_digest) = 32))
);

ALTER TABLE ONLY offchainreporting2_contract_configs
    ADD CONSTRAINT offchainreporting2_contract_configs_oracle_spec_fkey
        FOREIGN KEY (offchainreporting2_oracle_spec_id)
            REFERENCES offchainreporting2_oracle_specs(id)
            ON DELETE CASCADE;

CREATE TABLE offchainreporting2_latest_round_requested (
   offchainreporting2_oracle_spec_id INTEGER PRIMARY KEY,
   requester bytea NOT NULL,
   config_digest bytea NOT NULL,
   epoch bigint NOT NULL,
   round bigint NOT NULL,
   raw jsonb NOT NULL,
   CONSTRAINT offchainreporting2_latest_round_requested_config_digest_check
       CHECK ((octet_length(config_digest) = 32)),
   CONSTRAINT offchainreporting2_latest_round_requested_requester_check
       CHECK ((octet_length(requester) = 20))
);

ALTER TABLE offchainreporting2_latest_round_requested
    ADD CONSTRAINT offchainreporting2_latest_round_oracle_spec_fkey
        FOREIGN KEY (offchainreporting2_oracle_spec_id)
            REFERENCES offchainreporting2_oracle_specs(id)
            ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE;

CREATE TABLE offchainreporting2_persistent_states (
  offchainreporting2_oracle_spec_id integer NOT NULL,
  config_digest bytea NOT NULL,
  epoch bigint NOT NULL,
  highest_sent_epoch bigint NOT NULL,
  highest_received_epoch bigint[] NOT NULL,
  created_at timestamp with time zone NOT NULL,
  updated_at timestamp with time zone NOT NULL,
  CONSTRAINT offchainreporting2_persistent_states_config_digest_check
      CHECK ((octet_length(config_digest) = 32))
);

ALTER TABLE ONLY offchainreporting2_persistent_states
    ADD CONSTRAINT offchainreporting2_persistent_states_pkey
        PRIMARY KEY (offchainreporting2_oracle_spec_id, config_digest),
    ADD CONSTRAINT offchainreporting2_persistent_oracle_spec_fkey
        FOREIGN KEY (offchainreporting2_oracle_spec_id)
            REFERENCES offchainreporting2_oracle_specs(id)
            ON DELETE CASCADE;

CREATE TABLE offchainreporting2_pending_transmissions (
      offchainreporting2_oracle_spec_id integer NOT NULL,
      config_digest bytea NOT NULL,
      epoch bigint NOT NULL,
      round bigint NOT NULL,
      "time" timestamp with time zone NOT NULL,
      extra_hash bytea NOT NULL,
      report bytea NOT NULL,
      attributed_signatures bytea[] NOT NULL,
      created_at timestamp with time zone NOT NULL,
      updated_at timestamp with time zone NOT NULL,
      CONSTRAINT offchainreporting2_pending_transmissions_config_digest_check CHECK ((octet_length(config_digest) = 32))
);

ALTER TABLE ONLY offchainreporting2_pending_transmissions
    ADD CONSTRAINT offchainreporting2_pending_transmissions_pkey
        PRIMARY KEY (offchainreporting2_oracle_spec_id, config_digest, epoch, round),
    ADD CONSTRAINT offchainreporting2_pending_transmission_oracle_spec_fkey
        FOREIGN KEY (offchainreporting2_oracle_spec_id) REFERENCES offchainreporting2_oracle_specs(id)
            ON DELETE CASCADE;

CREATE INDEX idx_offchainreporting2_pending_transmissions_time ON offchainreporting2_pending_transmissions USING btree ("time");

-- After moving to the unified keystore the encrypted_p2p_keys table is no longer used
-- So we have to drop this FK to be able to uses the discoverer (v2) networking stack
ALTER TABLE offchainreporting_discoverer_announcements DROP CONSTRAINT offchainreporting_discoverer_announcements_local_peer_id_fkey;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE offchainreporting2_pending_transmissions;
DROP TABLE offchainreporting2_persistent_states;
DROP TABLE offchainreporting2_latest_round_requested;
DROP TABLE offchainreporting2_contract_configs;
ALTER TABLE jobs DROP CONSTRAINT chk_only_one_spec,
                 ADD CONSTRAINT chk_only_one_spec CHECK (
                         num_nonnulls(
                                 offchainreporting_oracle_spec_id,
                                 direct_request_spec_id,
                                 flux_monitor_spec_id,
                                 keeper_spec_id,
                                 cron_spec_id,
                                 webhook_spec_id,
                                 vrf_spec_id) = 1
                     );
ALTER TABLE jobs DROP COLUMN offchainreporting2_oracle_spec_id;
ALTER TABLE offchainreporting_discoverer_announcements ADD CONSTRAINT offchainreporting_discoverer_announcements_local_peer_id_fkey FOREIGN KEY (local_peer_id) REFERENCES encrypted_p2p_keys(peer_id) DEFERRABLE INITIALLY IMMEDIATE;
DROP TABLE offchainreporting2_oracle_specs;
-- +goose StatementEnd
