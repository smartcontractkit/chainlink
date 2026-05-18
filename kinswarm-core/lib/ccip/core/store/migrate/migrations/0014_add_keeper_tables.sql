-- +goose Up
	CREATE TABLE keeper_specs (
		id BIGSERIAL PRIMARY KEY,
		contract_address bytea NOT NULL,
		from_address bytea NOT NULL,
		created_at timestamp with time zone NOT NULL,
		updated_at timestamp with time zone NOT NULL,
		CONSTRAINT keeper_specs_contract_address_check CHECK ((octet_length(contract_address) = 20)),
		CONSTRAINT keeper_specs_from_address_check CHECK ((octet_length(from_address) = 20))
	);

	ALTER TABLE jobs ADD COLUMN keeper_spec_id INT REFERENCES keeper_specs(id),
	DROP CONSTRAINT chk_only_one_spec,
	ADD CONSTRAINT chk_only_one_spec CHECK (
		num_nonnulls(offchainreporting_oracle_spec_id, direct_request_spec_id, flux_monitor_spec_id, keeper_spec_id) = 1
	);

	CREATE TABLE keeper_registries (
		id BIGSERIAL PRIMARY KEY,
		job_id int UNIQUE NOT NULL REFERENCES jobs(id) ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE,
		keeper_index int NOT NULL,
		contract_address bytea UNIQUE NOT NULL,
		from_address bytea NOT NULL,
		check_gas int NOT NULL,
		block_count_per_turn int NOT NULL,
		num_keepers int NOT NULL
		CONSTRAINT keeper_registries_contract_address_check CHECK ((octet_length(contract_address) = 20))
		CONSTRAINT keeper_registries_from_address_check CHECK ((octet_length(from_address) = 20))
	);

	CREATE INDEX idx_keeper_registries_keeper_index ON keeper_registries(keeper_index);

	CREATE TABLE upkeep_registrations (
		id BIGSERIAL PRIMARY KEY,
		registry_id bigint NOT NULL REFERENCES keeper_registries(id) ON DELETE CASCADE DEFERRABLE INITIALLY IMMEDIATE,
		execute_gas int NOT NULL,
		check_data bytea NOT NULL,
		upkeep_id bigint NOT NULL,
		positioning_constant int NOT NULL
	);

	CREATE UNIQUE INDEX idx_upkeep_registrations_unique_upkeep_ids_per_keeper ON upkeep_registrations(registry_id, upkeep_id);
	CREATE INDEX idx_upkeep_registrations_upkeep_id ON upkeep_registrations(upkeep_id);

-- +goose Down
	DROP TABLE IF EXISTS keeper_specs, keeper_registries, upkeep_registrations;

	ALTER TABLE jobs DROP CONSTRAINT chk_only_one_spec,
		ADD CONSTRAINT chk_only_one_spec CHECK (
			num_nonnulls(offchainreporting_oracle_spec_id, direct_request_spec_id, flux_monitor_spec_id) = 1
		);

	ALTER TABLE jobs DROP COLUMN keeper_spec_id;
