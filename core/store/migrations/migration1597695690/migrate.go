package migration1597695690

import (
	"github.com/jinzhu/gorm"
)

// Migrate creates the offchain_reporting_job_specs table
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		--
		-- Pipeline Specs
		--
        CREATE TABLE pipeline_specs (
            id SERIAL PRIMARY KEY,
            dot_dag_source TEXT NOT NULL,
            created_at timestamptz NOT NULL
        );

        CREATE INDEX idx_pipeline_specs_created_at ON pipeline_specs USING BRIN (created_at);

		--
		-- Pipeline Task Specs
		--

        CREATE TABLE pipeline_task_specs (
            id SERIAL PRIMARY KEY,
            dot_id TEXT NOT NULL,
            pipeline_spec_id INT NOT NULL REFERENCES pipeline_specs (id) ON DELETE CASCADE,
            type TEXT NOT NULL,
            json jsonb NOT NULL,
            index INT NOT NULL DEFAULT 0,
            successor_id INT REFERENCES pipeline_task_specs (id),
            created_at timestamptz NOT NULL
        );

		COMMENT ON COLUMN pipeline_task_specs.dot_id IS 'Dot ID is included to help in debugging';

        CREATE INDEX idx_pipeline_task_specs_pipeline_spec_id ON pipeline_task_specs (pipeline_spec_id);
        CREATE INDEX idx_pipeline_task_specs_successor_id ON pipeline_task_specs (successor_id);
        CREATE INDEX idx_pipeline_task_specs_created_at ON pipeline_task_specs USING BRIN (created_at);

		--
		-- Pipeline Runs
		--

        CREATE TABLE pipeline_runs (
            id BIGSERIAL PRIMARY KEY,
            pipeline_spec_id INT NOT NULL REFERENCES pipeline_specs (id) ON DELETE CASCADE,
			-- FIXME: I propose to leave meta out of it for the purposes of cutting scope on this PR
            meta jsonb NOT NULL DEFAULT '{}',
            created_at timestamptz NOT NULL
        );

        CREATE INDEX idx_pipeline_runs_pipeline_spec_id ON pipeline_runs (pipeline_spec_id);
        CREATE INDEX idx_pipeline_runs_created_at ON pipeline_runs USING BRIN (created_at);

		--
		-- Pipeline Task Runs
		--

        CREATE TABLE pipeline_task_runs (
            id BIGSERIAL PRIMARY KEY,
            pipeline_run_id BIGINT NOT NULL REFERENCES pipeline_runs (id) ON DELETE CASCADE,
            output JSONB,
            error TEXT,
            pipeline_task_spec_id INT NOT NULL REFERENCES pipeline_task_specs (id) ON DELETE CASCADE,
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
        CREATE INDEX idx_pipeline_task_runs ON pipeline_task_runs USING BRIN (created_at);

        -- This query is used in the runner to find unstarted task runs
        CREATE INDEX idx_pipeline_task_runs_unfinished ON pipeline_task_runs (finished_at) WHERE finished_at IS NULL;

		--
		-- Offchainreporting Oracle Specs
		--

        ALTER TABLE offchainreporting_oracle_specs
			ADD COLUMN contract_address bytea NOT NULL,
			ADD COLUMN p2p_peer_id text NOT NULL REFERENCES encrypted_p2p_keys (peer_id),
			ADD COLUMN p2p_bootstrap_peers text[] NOT NULL,
			ADD COLUMN is_bootstrap_peer bool NOT NULL,
			ADD COLUMN encrypted_ocr_key_bundle_id bytea NOT NULL REFERENCES encrypted_ocr_key_bundles (id),
			ADD COLUMN monitoring_endpoint TEXT,
			ADD COLUMN transmitter_address bytea NOT NULL REFERENCES keys (address),
			ADD COLUMN observation_timeout bigint NOT NULL,
			ADD COLUMN blockchain_timeout bigint NOT NULL,
			ADD COLUMN contract_config_tracker_subscribe_interval bigint NOT NULL,
			ADD COLUMN contract_config_tracker_poll_interval bigint NOT NULL,
			ADD COLUMN contract_config_confirmations INT NOT NULL,
			ADD COLUMN created_at timestamptz NOT NULL,
			ADD COLUMN updated_at timestamptz NOT NULL,
			ADD CONSTRAINT chk_contract_address_length CHECK (octet_length(contract_address) = 20);

        CREATE UNIQUE INDEX idx_offchainreporting_oracle_specs_unique_key_bundles ON offchainreporting_oracle_specs (encrypted_ocr_key_bundle_id, contract_address);
        CREATE UNIQUE INDEX idx_offchainreporting_oracle_specs_unique_peer_ids ON offchainreporting_oracle_specs (p2p_peer_id, contract_address);

        CREATE INDEX idx_offchainreporting_oracle_specs_created_at ON offchainreporting_oracle_specs USING BRIN (created_at);
        CREATE INDEX idx_offchainreporting_oracle_specs_updated_at ON offchainreporting_oracle_specs USING BRIN (updated_at);

		--
		-- Jobs
		--

        -- NOTE: This will be extended with new IDs when we bring directrequest and fluxmonitor under the new jobspawner umbrella
        -- Only ONE id should ever be present
        CREATE TABLE jobs (
            id SERIAL PRIMARY KEY,
            pipeline_spec_id INT REFERENCES pipeline_specs (id) ON DELETE CASCADE,
            offchainreporting_oracle_spec_id INT REFERENCES offchainreporting_oracle_specs (id) ON DELETE CASCADE,
            CONSTRAINT chk_valid CHECK (
                offchainreporting_oracle_spec_id IS NOT NULL
            )
        );
		CREATE UNIQUE INDEX idx_jobs_unique_offchain_reporting_oracle_spec_id ON jobs (offchainreporting_oracle_spec_id);
		CREATE UNIQUE INDEX idx_jobs_unique_pipeline_spec_id ON jobs (pipeline_spec_id);

		--
		-- Log Consumptions
		--

        ALTER TABLE log_consumptions
			ADD COLUMN job_id_v2 INT REFERENCES jobs (id) ON DELETE CASCADE,
        	ALTER COLUMN job_id DROP NOT NULL,
        	ADD CONSTRAINT chk_log_consumptions_exactly_one_job_id CHECK (
				job_id IS NOT NULL AND job_id_v2 IS NULL
				OR
				job_id_v2 IS NOT NULL AND job_id IS NULL
			);
        DROP INDEX log_consumptions_unique_idx;
        CREATE UNIQUE INDEX log_consumptions_unique_v1_idx ON log_consumptions (job_id, block_hash, log_index);
        CREATE UNIQUE INDEX log_consumptions_unique_v2_idx ON log_consumptions (job_id_v2, block_hash, log_index);


		---- TRIGGERS ----


        ---
        --- Notify the Chainlink node when a new pipeline run has started
        ---

        CREATE OR REPLACE FUNCTION notifyPipelineRunStarted() RETURNS TRIGGER AS $_$
        BEGIN
            PERFORM pg_notify('pipeline_run_started', NEW.id::text);
            RETURN NEW;
        END
        $_$ LANGUAGE 'plpgsql';

        CREATE TRIGGER notify_pipeline_run_started
        AFTER INSERT ON pipeline_runs
        FOR EACH ROW EXECUTE PROCEDURE notifyPipelineRunStarted();

        ---
        --- Notify the Chainlink node when a new job spec is created
        ---

        CREATE OR REPLACE FUNCTION notifyJobCreated() RETURNS TRIGGER AS $_$
        BEGIN
            PERFORM pg_notify('insert_on_jobs', NEW.id::text);
            RETURN NEW;
        END
        $_$ LANGUAGE 'plpgsql';

        CREATE TRIGGER notify_job_created
        AFTER INSERT ON jobs
        FOR EACH ROW EXECUTE PROCEDURE notifyJobCreated();
    `).Error
}
