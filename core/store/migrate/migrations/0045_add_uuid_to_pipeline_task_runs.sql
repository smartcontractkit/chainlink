-- +goose Up
-- Truncate tables to ease bigint -> UUID migration
	TRUNCATE TABLE pipeline_runs, pipeline_task_runs, flux_monitor_round_stats_v2;

	-- Migrate pipeline_task_runs to UUID
	ALTER TABLE pipeline_task_runs DROP CONSTRAINT pipeline_task_runs_pkey; 
	ALTER TABLE pipeline_task_runs DROP COLUMN id; 
	ALTER TABLE pipeline_task_runs ADD COLUMN id uuid PRIMARY KEY;

	-- Add state & inputs to pipeline_runs
	ALTER TABLE pipeline_runs ADD COLUMN inputs jsonb;
	CREATE TYPE pipeline_runs_state AS ENUM (
	    'running',
	    'suspended',
	    'errored',
	    'completed'
	);
	ALTER TABLE pipeline_runs ADD COLUMN state pipeline_runs_state NOT NULL DEFAULT 'completed';

	ALTER TABLE pipeline_runs DROP CONSTRAINT pipeline_runs_check;
	ALTER TABLE pipeline_runs ADD CONSTRAINT pipeline_runs_check CHECK (
		((state IN ('completed', 'errored')) AND (finished_at IS NOT NULL) AND (num_nulls(outputs, errors) = 0))
			OR 
		((state IN ('running', 'suspended')) AND num_nulls(finished_at, outputs, errors) = 3)
	);

-- +goose Down
	ALTER TABLE pipeline_runs DROP CONSTRAINT pipeline_runs_check;
	ALTER TABLE pipeline_runs ADD CONSTRAINT pipeline_runs_check CHECK (
		(((outputs IS NULL) AND (errors IS NULL) AND (finished_at IS NULL))
		OR ((outputs IS NOT NULL) AND (errors IS NOT NULL) AND (finished_at IS NOT NULL)))
	)
	DROP CONSTRAINT IF EXISTS pipeline_task_runs_run_id_key;
	ALTER TABLE pipeline_task_runs DROP COLUMN run_id;
	ALTER TABLE pipeline_runs DROP COLUMN inputs;
	DROP TYPE pipeline_runs_state;
