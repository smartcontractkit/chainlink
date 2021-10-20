package migrations

import (
	"gorm.io/gorm"
)

const (
	up6 = `
CREATE UNIQUE INDEX idx_pipeline_task_runs_unique_task_spec_id_per_run ON pipeline_task_runs (pipeline_task_spec_id, pipeline_run_id);
`
	down6 = `
DROP INDEX IF EXISTS idx_pipeline_task_runs_unique_task_spec_id_per_run;
`
)

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0006_unique_task_specs_per_pipeline_run",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up6).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down6).Error
		},
	})
}
