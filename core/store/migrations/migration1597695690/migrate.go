package migration1597695690

import (
	"github.com/jinzhu/gorm"
)

// Migrate creates the offchain_reporting_job_specs table
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	 	CREATE TABLE offchainreporting_key_bundles (
	 		-- NOTE: Key bundle ID is intended to be set by software as sha256 hash of {onchain sig pubkey, offchain sig pubkey, config decryption pubkey}
	 		id bytea NOT NULL PRIMARY KEY,
	 		CONSTRAINT chk_id_length CHECK (octet_length(id) = 32),
	 		encrypted_priv_key_bundle jsonb NOT NULL,
	 		created_at timestamptz NOT NULL,
	 	);

		CREATE INDEX idx_offchainreporting_oracles_unique_key_bundles_created_at ON offchainreporting_key_bundles USING BRIN (created_at);

	 	CREATE TABLE offchainreporting_oracles (
	 		id BIGSERIAL PRIMARY KEY,
	 		contract_address bytea NOT NULL,
	 		CONSTRAINT chk_contract_address_length CHECK (octet_length(contract_address) = 20),
	 		p2p_peer_id text NOT NULL REFERENCES encrypted_p2p_keys (peer_id),
	 		p2p_bootstrap_peers jsonb NOT NULL,
	 		key_bundle_id bytea NOT NULL REFERENCES offchainreporting_key_bundles (id),
	 		monitoring_endpoint TEXT,
	 		transmitter_address bytea NOT NULL REFERENCES keys (address),
	 		observation_timeout interval NOT NULL,
			data_fetch_pipeline_spec_id BIGINT NOT NULL REFERENCES pipeline_specs (id),
	 		created_at timestamptz NOT NULL,
			updated_at timestamptz NOT NULL
	 	);

		CREATE UNIQUE INDEX idx_offchainreporting_oracles_unique_key_bundles ON offchainreporting_oracles (key_bundle_id, contract_address);
		CREATE UNIQUE INDEX idx_offchainreporting_oracles_unique_peer_ids ON offchainreporting_oracles (p2p_peer_id, contract_address);
		CREATE INDEX idx_offchainreporting_oracles_data_fetch_pipeline_spec_id ON offchainreporting_oracles (data_fetch_pipeline_spec_id);
		
		CREATE INDEX idx_offchainreporting_oracles_created_at ON offchainreporting_oracles USING BRIN (created_at);
		CREATE INDEX idx_offchainreporting_oracles_updated_at ON offchainreporting_oracles USING BRIN (updated_at);

		CREATE TABLE pipeline_specs (
			-- id is intended to be a sha256 hash of the json representation of the DOT dag
			id BYTEA PRIMARY KEY,
			CONSTRAINT chk_id_length CHECK (octet_length(id) = 32),
			source_dot_dag TEXT NOT NULL,
			created_at timestamptz NOT NULL
		);

		CREATE INDEX idx_pipeline_specs_created_at ON pipeline_specs USING BRIN (created_at);

		CREATE TABLE pipeline_task_specs (
			id BIGSERIAL PRIMARY KEY,
			pipeline_spec_id BIGINT NOT NULL REFERENCES pipeline_specs (id),

			task_spec jsonb NOT NULL,

			created_at timestamptz NOT NULL,
		);

		CREATE INDEX idx_pipeline_task_specs_created_at ON pipeline_task_specs USING BRIN (created_at);

		CREATE TABLE pipeline_task_spec_edges (
			predecessor_id BIGINT NOT NULL REFERENCES pipeline_task_specs (id),
			successor_id BIGINT NOT NULL REFERENCES pipeline_task_specs (id),
			PRIMARY KEY(successor_id, predecessor_id)
		);

		-- This index is a little confusing, but the result is to only allow one successor (child) for any single predecessor (parent)
		CREATE UNIQUE INDEX idx_pipeline_task_run_edges_only_one_successor_per_predecessor pipeline_task_spec_edges (predecessor_id);

		CREATE TABLE pipeline_runs (
			id BIGSERIAL PRIMARY KEY,
			pipeline_spec_id BIGINT NOT NULL REFERENCES pipeline_specs (id),
			created_at timestamptz NOT NULL
			-- NOTE: Could denormalize here with finished_at/output/error of last task_run if that proves necessary for performance
		);

		CREATE INDEX idx_pipeline_runs_pipeline_spec_id ON pipeline_runs (pipeline_spec_id);
		CREATE INDEX idx_pipeline_runs_created_at ON pipeline_runs USING BRIN (created_at);

		CREATE TABLE pipeline_task_runs (
			id BIGSERIAL PRIMARY KEY,
			pipeline_run_id BIGINT NOT NULL REFERENCES pipeline_runs (id),
			output JSONB,
			error TEXT, 
			pipeline_task_spec_id BIGINT NOT NULL REFERENCES pipeline_task_specs (id),
			created_at timestamptz NOT NULL,
			finished_at timestamptz,
			CONSTRAINT chk_pipeline_task_run_fsm CHECK (
				error IS NULL AND output IS NULL AND finished_at IS NULL
				OR
				error IS NULL AND output IS NOT NULL AND finished_at IS NOT NULL
				OR
				output IS NULL AND error IS NOT NULL AND finished_at IS NOT NULL
			)
		);

		CREATE INDEX idx_pipeline_task_runs ON pipeline_task_tuns USING BRIN (created_at);
		-- TODO: more indexes for pipeline_task_runs (dependent on queries in the pipeline runner)
	`).Error

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
