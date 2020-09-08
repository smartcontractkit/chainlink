package migration1597695690

import (
	"github.com/jinzhu/gorm"
)

// Migrate creates the offchain_reporting_job_specs table
func Migrate(tx *gorm.DB) error {
	// JID                *models.ID         `toml:"jobID"              gorm:"not null;column:job_id"`
	// ContractAddress    common.Address     `toml:"contractAddress"    gorm:"not null"`
	// P2PNodeID          string             `toml:"p2pNodeID"          gorm:"not null"`
	// P2PBootstrapNodes  []P2PBootstrapNode `toml:"p2pBootstrapNodes"  gorm:"not null;type:jsonb"`
	// KeyBundle          string             `toml:"keyBundle"          gorm:"not null"`
	// MonitoringEndpoint string             `toml:"monitoringEndpoint" gorm:"not null"`
	// NodeAddress        common.Address     `toml:"nodeAddress"        gorm:"not null"`
	// ObservationTimeout time.Duration      `toml:"observationTimeout" gorm:"not null"`
	// ObservationSource  pipeline.TaskDAG   `toml:"observationSource"  gorm:"-"`
	// LogLevel           ocrtypes.LogLevel  `toml:"logLevel,omitempty"`
	return tx.Exec(`
	 	CREATE TABLE offchainreporting_key_bundles (
	 		-- NOTE: Key bundle ID is intended to be set by software as keccak256 hash of {onchain sig pubkey, offchain sig pubkey, config decryption pubkey}
	 		id bytea NOT NULL PRIMARY KEY,
	 		CONSTRAINT chk_id_length CHECK (octet_length(id) = 32),
	 		encrypted_priv_key_bundle jsonb NOT NULL,
	 		created_at timestamptz NOT NULL,
	 	);

	 	CREATE TABLE offchainreporting_oracles (
	 		id BIGSERIAL PRIMARY KEY,
	 		contract_address bytea NOT NULL,
	 		CONSTRAINT chk_contract_address_length CHECK (octet_length(contract_address) = 20),
	 		p2p_peer_id text NOT NULL REFERENCES encrypted_p2p_keys (peer_id),
	 		p2p_bootstrap_peers jsonb NOT NULL, -- NOTE: Needs revisiting
	 		key_bundle_id bytea NOT NULL REFERENCES offchainreporting_key_bundles (id),
	 		monitoring_endpoint TEXT,
	 		transmitter_address bytea NOT NULL REFERENCES keys (address),
	 		observation_timeout interval NOT NULL,
			-- TODO: Revisit
	 		-- observation_source_id bigint NOT NULL REFERENCES observation_aggregators(id),
	 		created_at timestamptz NOT NULL
	 	);

		CREATE UNIQUE INDEX idx_offchainreporting_oracles_unique_key_bundles ON offchainreporting_oracles (key_bundle_id, contract_address);
		CREATE UNIQUE INDEX idx_offchainreporting_oracles_unique_peer_ids ON offchainreporting_oracles (p2p_peer_id, contract_address);

		CREATE TABLE offchainreporting_oracles_pipeline_specs (
			offchainreporting_oracle_id BIGINT REFERENCES offchainreporting_oracles (id) NOT NULL,
			pipeline_spec_id BIGINT REFERENCES pipeline_specs (id) NOT NULL,
			PRIMARY KEY(offchainreporting_oracle_id, pipeline_spec_id)
		);

		CREATE INDEX idx_offchainreporting_oracles_pipeline_specs_pipeline_spec_id offchainreporting_oracles_pipeline_specs (pipeline_spec_id);


		CREATE TABLE pipeline_runs (
			id BIGSERIAL PRIMARY KEY,
			pipeline_spec_id BIGINT NOT NULL REFERENCES pipeline_specs (id)
		);

		CREATE TABLE pipeline_task_runs (
			id BIGSERIAL PRIMARY KEY,
			pipeline_run_id BIGINT NOT NULL REFERENCES pipeline_runs (id),
			output JSONB,
			error TEXT, -- TODO: Actually should probabl reference errors instead?
			pipeline_task_id BIGINT NOT NULL REFERENCES pipeline_tasks (id),
			-- TODO: These columns
			-- started_at timestamptz,
			-- finished_at timestamptz,
			CONSTRAINT chk_pipeline_task_run_fsm CHECK (
				error IS NULL AND output IS NULL
				OR
				error IS NULL AND output IS NOT NULL
				OR
				output IS NULL AND error IS NOT NULL
			)
		);

		-- TODO: indexes for pipeline_task_runs

		CREATE TABLE pipeline_task_run_edges (
			parent_id BIGINT NOT NULL REFERENCES pipeline_task_runs (id),
			child_id BIGINT NOT NULL REFERENCES pipeline_task_runs (id),
			PRIMARY KEY(child_id, parent_id)
		);

		CREATE INDEX idx_pipeline_task_run_edges ON pipeline_task_run_edges (parent_id);

		CREATE TABLE pipeline_tasks (
			id BIGSERIAL PRIMARY KEY,

			task jsonb NOT NULL,

			created_at timestamptz NOT NULL,
			updated_at timestamptz NOT NULL
		);

		-- XXX: Do we use task edges or represent it as a serialized DAG?
		CREATE TABLE pipeline_task_edges (
			parent_id BIGINT NOT NULL REFERENCES pipeline_tasks (id),
			child_id BIGINT NOT NULL REFERENCES pipeline_tasks (id),
			PRIMARY KEY(child_id, parent_id)
		);
	`).Error

	// 	CREATE TYPE observation_aggregator_type AS ENUM ('median');

	// 	CREATE TABLE observation_aggregators (
	// 		id BIGSERIAL PRIMARY KEY,
	// 		type observation_aggregator_type NOT NULL,
	// 		created_at timestamptz NOT NULL,
	// 	);

	// 	CREATE TYPE observation_source_type AS ENUM ('http', 'bridge');

	// 	CREATE TABLE observation_sources (
	// 		id BIGSERIAL PRIMARY KEY,
	// 		name TEXT,
	// 		type observation_source_type NOT NULL,
	// 		params jsonb NOT NULL,
	// 		created_at timestamptz NOT NULL,
	// 		updated_at timestamptz NOT NULL,
	// 	);

	// 	CREATE UNIQUE INDEX idx_unique_observation_sources ON observation_sources (params);
	// 	CREATE UNIQUE INDEX idx_unique_observation_source_names ON observation_sources (name);

	// 	CREATE TABLE observation_aggregator_sources (
	// 		observation_aggregator_id BIGINT NOT NULL REFERENCES observation_aggregators (id),
	// 		observation_source_id BIGINT NOT NULL REFERENCES observation_sources (id),
	// 		PRIMARY KEY(observation_aggregator_id, observation_source_id),
	// 		created_at timestamptz NOT NULL
	// 	);

	// 	CREATE INDEX idx_observation_aggregator_sources_observation_source_id ON observation_aggregator_sources (observation_source_id);

	// `).Error

	// err = tx.Exec(`
	//        CREATE TABLE offchain_reporting_persistent_states (
	//            id SERIAL PRIMARY KEY,
	//            job_spec_id uuid NOT NULL,
	//            group_id bytea NOT NULL,
	//            epoch integer NOT NULL,
	//            highest_sent_epoch integer NOT NULL,
	//            highest_received_epoch integer[31] NOT NULL
	//        );
	//        ALTER TABLE offchain_reporting_persistent_states ADD CONSTRAINT "ocr_job_spec_id_fkey" FOREIGN KEY ("job_spec_id") REFERENCES offchain_reporting_job_specs ("id") ON DELETE CASCADE;
	//        ALTER TABLE offchain_reporting_persistent_states ADD CONSTRAINT chk_group_id_length CHECK (
	//            octet_length(group_id) = 32
	//        );
	//        CREATE UNIQUE INDEX ocr_persistent_states_unique_idx ON offchain_reporting_persistent_states ("job_id", "group_id");
	//    `).Error
	// if err != nil {
	// 	return err
	// }

	// err = tx.Exec(`
	//        CREATE TABLE offchain_reporting_configs (
	//            id SERIAL PRIMARY KEY,
	//            job_spec_id uuid NOT NULL,
	//            group_id bytea NOT NULL,
	//            oracles jsonb NOT NULL,
	//            secret bytea NOT NULL,
	//            f integer NOT NULL,
	//            delta_progress integer NOT NULL,
	//            delta_resend integer NOT NULL,
	//            delta_round integer NOT NULL,
	//            delta_observe integer NOT NULL,
	//            delta_c integer NOT NULL,
	//            alpha float NOT NULL,
	//            r_max integer NOT NULL,
	//            delta_stage integer NOT NULL,
	//            schedule integer[] NOT NULL
	//        );
	//        ALTER TABLE offchain_reporting_configs ADD CONSTRAINT "ocr_job_spec_id_fkey" FOREIGN KEY ("job_spec_id") REFERENCES offchain_reporting_job_specs ("id") ON DELETE CASCADE;
	//        ALTER TABLE offchain_reporting_configs ADD CONSTRAINT chk_group_id_length CHECK (
	//            octet_length(group_id) = 27
	//        );
	//        ALTER TABLE offchain_reporting_configs ADD CONSTRAINT chk_secret_length CHECK (
	//            octet_length(secret) = 16
	//        );
	//        CREATE UNIQUE INDEX ocr_configs_unique_idx ON offchain_reporting_configs ("job_id", "group_id");
	//   `).Error
	// if err != nil {
	// 	return err
	// }
}
