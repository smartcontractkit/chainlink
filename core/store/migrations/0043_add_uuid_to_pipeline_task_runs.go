package migrations

import (
	"gorm.io/gorm"
)

const up43 = `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	ALTER TABLE pipeline_task_runs ADD COLUMN task_run_id uuid NOT NULL UNIQUE DEFAULT uuid_generate_v4 ();
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

// TODO: update the state machine constraint to include state

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
