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
            created_at timestamptz NOT NULL
        );

        CREATE INDEX idx_offchainreporting_oracle_specs_unique_key_bundles_created_at ON offchainreporting_key_bundles USING BRIN (created_at);

        -- NOTE: This will be extended with new IDs when we bring directrequest and fluxmonitor under the new jobspawner umbrella
        -- Only ONE id should ever be present
        CREATE TABLE jobs (
            id SERIAL PRIMARY KEY
        );

        CREATE TABLE offchainreporting_oracle_specs (
            id SERIAL PRIMARY KEY,
            job_id INT NOT NULL REFERENCES jobs (id),
            contract_address bytea NOT NULL,
            CONSTRAINT chk_contract_address_length CHECK (octet_length(contract_address) = 20),
            p2p_peer_id text NOT NULL REFERENCES encrypted_p2p_keys (peer_id),
            p2p_bootstrap_peers jsonb NOT NULL,
            key_bundle_id bytea NOT NULL REFERENCES offchainreporting_key_bundles (id),
            monitoring_endpoint TEXT,
            transmitter_address bytea NOT NULL REFERENCES keys (address),
            observation_timeout interval NOT NULL,
            data_fetch_pipeline_spec_id INT NOT NULL REFERENCES pipeline_specs (id),
            created_at timestamptz NOT NULL,
            updated_at timestamptz NOT NULL
        );

        ALTER TABLE jobs ADD COLUMN
            offchainreporting_oracle_spec_id INT REFERENCES offchainreporting_oracle_specs (id)
            CONSTRAINT chk_valid CHECK (
                offchainreporting_oracle_spec_id IS NOT NULL
            );
        CREATE UNIQUE INDEX idx_jobs_unique_offchainreporting_oracle_spec_ids ON jobs (offchainreporting_oracle_spec_id);

        CREATE UNIQUE INDEX idx_offchainreporting_oracle_specs_unique_key_bundles ON offchainreporting_oracle_specs (key_bundle_id, contract_address);
        CREATE UNIQUE INDEX idx_offchainreporting_oracle_specs_unique_peer_ids ON offchainreporting_oracle_specs (p2p_peer_id, contract_address);
        CREATE UNIQUE INDEX idx_offchainreporting_oracle_specs_unique_job_ids ON offchainreporting_oracle_specs (job_id);
        CREATE INDEX idx_offchainreporting_oracle_specs_data_fetch_pipeline_spec_id ON offchainreporting_oracle_specs (data_fetch_pipeline_spec_id);

        CREATE INDEX idx_offchainreporting_oracle_specs_created_at ON offchainreporting_oracle_specs USING BRIN (created_at);
        CREATE INDEX idx_offchainreporting_oracle_specs_updated_at ON offchainreporting_oracle_specs USING BRIN (updated_at);

        CREATE TABLE pipeline_specs (
            id SERIAL PRIMARY KEY,
            source_dot_dag TEXT NOT NULL,
            created_at timestamptz NOT NULL
        );

        CREATE INDEX idx_pipeline_specs_created_at ON pipeline_specs USING BRIN (created_at);

        CREATE TABLE pipeline_task_specs (
            id SERIAL PRIMARY KEY,
            pipeline_spec_id INT NOT NULL REFERENCES pipeline_specs (id),

            task_type TEXT NOT NULL,
            task_json jsonb NOT NULL,

            successor_id INT REFERENCES pipeline_task_specs (id),

            created_at timestamptz NOT NULL
        );

        CREATE INDEX idx_pipeline_task_specs_pipeline_spec_id ON pipeline_task_specs (pipeline_spec_id);
        CREATE INDEX idx_pipeline_task_specs_successor_id ON pipeline_task_specs (successor_id);
        CREATE INDEX idx_pipeline_task_specs_created_at ON pipeline_task_specs USING BRIN (created_at);

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

        -- NOTE: This table is large and insert/update heavy so we must be efficient with indexes

        CREATE INDEX idx_pipeline_task_runs ON pipeline_task_tuns USING BRIN (created_at);
        -- CREATE INDEX idx_pipeline_task_runs_finished_at ON pipeline_task_tuns WHERE finished_at IS NOT NULL USING BRIN (finished_at);

        -- This query is used in the runner to find unstarted task runs
        -- CREATE INDEX idx_pipeline_task_runs_unfinished ON pipeline_task_runs (finished_at) WHERE finished_at IS NULL;
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
	//  return err
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
	//  return err
	// }
}
