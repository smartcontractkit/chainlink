package migrations

import (
	"gorm.io/gorm"
)

const up43 = `
	-- Truncate tables to ease bigint -> UUID migration
	TRUNCATE TABLE pipeline_runs, pipeline_task_runs, flux_monitor_round_stats_v2;

	-- Migrate pipeline_task_runs to UUID
	ALTER TABLE pipeline_task_runs DROP CONSTRAINT pipeline_task_runs_pkey; 
	ALTER TABLE pipeline_task_runs DROP COLUMN id; 
	ALTER TABLE pipeline_task_runs ADD COLUMN id uuid PRIMARY KEY;

	-- Drop dependent constraints and indices
	-- DROP INDEX flux_monitor_round_stats_v2_pipeline_run_id_idx;
	-- DROP INDEX idx_unfinished_pipeline_task_runs;
	-- DROP INDEX pipeline_task_runs_pipeline_run_id_dot_id_idx;
	-- DROP CONSTRAINT flux_monitor_round_stats_v2_pipeline_run_id_fkey;
	-- DROP CONSTRAINT pipeline_task_runs_pipeline_run_id_fkey;

	-- ALTER TABLE pipeline_runs DROP CONSTRAINT pipeline_runs_pkey; 
	-- ALTER TABLE pipeline_runs DROP COLUMN id; 
	-- ALTER TABLE pipeline_runs ADD COLUMN id uuid PRIMARY KEY;

	-- CREATE INDEX flux_monitor_round_stats_v2_pipeline_run_id_idx ON public.flux_monitor_round_stats_v2 USING btree (pipeline_run_id)
	-- CREATE INDEX idx_unfinished_pipeline_task_runs ON public.pipeline_task_runs USING btree (pipeline_run_id) WHERE (finished_at IS NULL);
	-- CREATE UNIQUE INDEX pipeline_task_runs_pipeline_run_id_dot_id_idx ON public.pipeline_task_runs USING btree (pipeline_run_id, dot_id);
	-- ALTER TABLE ONLY public.flux_monitor_round_stats_v2
	--     ADD CONSTRAINT flux_monitor_round_stats_v2_pipeline_run_id_fkey FOREIGN KEY (pipeline_run_id) REFERENCES public.pipeline_runs(id) ON DELETE CASCADE;
	-- ALTER TABLE ONLY public.pipeline_task_runs
	--     ADD CONSTRAINT pipeline_task_runs_pipeline_run_id_fkey FOREIGN KEY (pipeline_run_id) REFERENCES public.pipeline_runs(id) ON DELETE CASCADE DEFERRABLE;
	--     pipeline_run_id on flux_monitor_round_stats_v2, pipeline_task_runs


	-- Add state & inputs to pipeline_runs
	ALTER TABLE pipeline_runs ADD COLUMN inputs jsonb;
	CREATE TYPE pipeline_runs_state AS ENUM (
	    'running',
	    'suspended',
	    'errored',
	    'completed'
	);
	ALTER TABLE pipeline_runs ADD COLUMN state pipeline_runs_state DEFAULT 'completed';

	ALTER TABLE pipeline_runs DROP CONSTRAINT pipeline_runs_check;
	ALTER TABLE pipeline_runs ADD CONSTRAINT pipeline_runs_check CHECK (
		((state IN ('completed', 'errored')) AND (finished_at IS NOT NULL) AND (num_nulls(outputs, errors) = 0))
			OR 
		((state IN ('running', 'suspended')) AND num_nulls(finished_at, outputs, errors) = 3)
	);
`

const down43 = `
	ALTER TABLE pipeline_runs DROP CONSTRAINT pipeline_runs_check;
	ALTER TABLE pipeline_runs ADD CONSTRAINT pipeline_runs_check CHECK (
		(((outputs IS NULL) AND (errors IS NULL) AND (finished_at IS NULL))
		OR ((outputs IS NOT NULL) AND (errors IS NOT NULL) AND (finished_at IS NOT NULL)))
	)
	DROP CONSTRAINT IF EXISTS pipeline_task_runs_run_id_key;
	ALTER TABLE pipeline_task_runs DROP COLUMN run_id;
	ALTER TABLE pipeline_runs DROP COLUMN inputs;
	DROP TYPE pipeline_runs_state;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0043_add_uuid_to_pipeline_task_runs",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up43).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down43).Error
		},
	})
}
