-- +goose Up
CREATE TABLE offchainreporting2_oracle_specs (
    id SERIAL PRIMARY KEY,
    contract_address bytea NOT NULL,
    p2p_peer_id text,
    p2p_bootstrap_peers text[],
    is_bootstrap_peer boolean NOT NULL,
    encrypted_ocr_key_bundle_id bytea,
    monitoring_endpoint text,
    transmitter_address bytea,
    observation_timeout bigint,
    blockchain_timeout bigint,
    contract_config_tracker_subscribe_interval bigint,
    contract_config_tracker_poll_interval bigint,
    contract_config_confirmations integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    CONSTRAINT chk_contract_address_length CHECK ((octet_length(contract_address) = 20))
);

ALTER TABLE ONLY offchainreporting2_oracle_specs
    ADD CONSTRAINT offchainreporting2_oracle_specs_unique_contract_addr UNIQUE (contract_address),

    ADD CONSTRAINT offchainreporting2_oracle_specs_p2p_peer_id_fkey
		FOREIGN KEY (p2p_peer_id)
		REFERENCES encrypted_p2p_keys(peer_id),

    ADD CONSTRAINT offchainreporting2_oracle_specs_transmitter_address_fkey
		FOREIGN KEY (transmitter_address)
		REFERENCES keys(address);

CREATE INDEX idx_offchainreporting2_oracle_specs_created_at ON offchainreporting2_oracle_specs USING brin (created_at);
CREATE INDEX idx_offchainreporting2_oracle_specs_updated_at ON offchainreporting2_oracle_specs USING brin (updated_at);

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
		offchainreporting2_oracle_spec_id BIGSERIAL PRIMARY KEY,
    config_digest bytea NOT NULL,
		config_count bigint NOT NULL,
    signers bytea[],
    transmitters text[],
		f smallint NOT NULL,
		onchain_config bytea,
		offchain_config_version bigint NOT NULL,
		offchain_config bytea,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    CONSTRAINT offchainreporting2_contract_configs_config_digest_check CHECK ((octet_length(config_digest) = 32))
);

ALTER TABLE ONLY offchainreporting2_contract_configs
    ADD CONSTRAINT offchainreporting2_contract_co_offchainreporting2_oracle_spe_fkey
		FOREIGN KEY (offchainreporting2_oracle_spec_id)
		REFERENCES offchainreporting2_oracle_specs(id) ON DELETE CASCADE;

CREATE TABLE offchainreporting2_latest_round_requested (
		offchainreporting2_oracle_spec_id BIGSERIAL PRIMARY KEY,
    requester bytea NOT NULL,
    config_digest bytea NOT NULL,
    epoch bigint NOT NULL,
    round bigint NOT NULL,
    raw jsonb NOT NULL,
    CONSTRAINT offchainreporting2_latest_round_requested_config_digest_check CHECK ((octet_length(config_digest) = 32)),
    CONSTRAINT offchainreporting2_latest_round_requested_requester_check CHECK ((octet_length(requester) = 20))
);

ALTER TABLE offchainreporting2_latest_round_requested
	ADD CONSTRAINT offchainreporting2_latest_roun_offchainreporting2_oracle_spe_fkey
	FOREIGN KEY (offchainreporting2_oracle_spec_id)
	REFERENCES offchainreporting2_oracle_specs(id)
	ON DELETE CASCADE DEFERRABLE;

CREATE TABLE offchainreporting2_persistent_states (
    offchainreporting2_oracle_spec_id integer NOT NULL,
    config_digest bytea NOT NULL,
    epoch bigint NOT NULL,
    highest_sent_epoch bigint NOT NULL,
    highest_received_epoch bigint[] NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    CONSTRAINT offchainreporting2_persistent_states_config_digest_check CHECK ((octet_length(config_digest) = 32))
);

ALTER TABLE ONLY offchainreporting2_persistent_states

	ADD CONSTRAINT offchainreporting2_persistent_states_pkey
	PRIMARY KEY (offchainreporting2_oracle_spec_id, config_digest),

	ADD CONSTRAINT offchainreporting2_persistent__offchainreporting2_oracle_spe_fkey
	FOREIGN KEY (offchainreporting2_oracle_spec_id)
	REFERENCES offchainreporting2_oracle_specs(id)
	ON DELETE CASCADE;

CREATE TABLE offchainreporting2_pending_transmissions (
    offchainreporting2_oracle_spec_id integer NOT NULL,
    config_digest bytea NOT NULL,
    epoch bigint NOT NULL,
    round bigint NOT NULL,
    "time" timestamp with time zone NOT NULL,
		extra_hash bytea,
		report bytea,
		attributed_signatures bytea[],
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    CONSTRAINT offchainreporting2_pending_transmissions_config_digest_check CHECK ((octet_length(config_digest) = 32))
);

ALTER TABLE ONLY offchainreporting2_pending_transmissions
	ADD CONSTRAINT offchainreporting2_pending_transmissions_pkey PRIMARY KEY (offchainreporting2_oracle_spec_id, config_digest, epoch, round),

	ADD CONSTRAINT offchainreporting2_pending_tra_offchainreporting2_oracle_spe_fkey FOREIGN KEY (offchainreporting2_oracle_spec_id) REFERENCES offchainreporting2_oracle_specs(id) ON DELETE CASCADE;

 CREATE INDEX idx_offchainreporting2_pending_transmissions_time ON offchainreporting2_pending_transmissions USING btree ("time");

CREATE TABLE encrypted_ocr2_key_bundles (
    id bytea NOT NULL,

		onchain_public_key bytea NOT NULL,
		onchain_signing_address bytea NOT NULL,

		offchain_signing_public_key bytea NOT NULL,
		offchain_encryption_public_key bytea NOT NULL,

    encrypted_private_keys jsonb NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    deleted_at timestamp with time zone
);

ALTER TABLE ONLY encrypted_ocr2_key_bundles
    ADD CONSTRAINT encrypted_ocr2_key_bundles_pkey PRIMARY KEY (id);

ALTER TABLE ONLY offchainreporting2_oracle_specs
    ADD CONSTRAINT offchainreporting2_oracle_specs_encrypted_ocr_key_bundle_id_fkey
		FOREIGN KEY (encrypted_ocr_key_bundle_id)
		REFERENCES encrypted_ocr2_key_bundles(id);

CREATE TABLE offchainreporting2_discoverer_announcements (
    local_peer_id text NOT NULL,
    remote_peer_id text NOT NULL,
    ann bytea NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
		PRIMARY KEY(local_peer_id, remote_peer_id)
);

ALTER TABLE ONLY offchainreporting2_discoverer_announcements
    ADD CONSTRAINT offchainreporting2_discoverer_announcements_local_peer_id_fkey
		FOREIGN KEY (local_peer_id)
		REFERENCES encrypted_p2p_keys(peer_id) DEFERRABLE;

-- +goose Down
DROP TABLE offchainreporting2_discoverer_announcements;
DROP TABLE offchainreporting2_pending_transmissions;
DROP TABLE offchainreporting2_persistent_states;
DROP TABLE offchainreporting2_latest_round_requested;
DROP TABLE offchainreporting2_contract_configs;
DROP TABLE offchainreporting2_oracle_specs;
