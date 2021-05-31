package migrations

import (
	"gorm.io/gorm"
)

const (
	up22 = `CREATE INDEX idx_unfinished_pipeline_task_runs ON pipeline_task_runs (pipeline_run_id) WHERE finished_at IS NULL;`

	down22 = `DROP INDEX idx_unfinished_pipeline_task_runs;`
)

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0022_unfinished_pipeline_task_run_idx",
		Migrate: func(db *gorm.DB) error {
			return db.Exec(up22).Error
		},
		Rollback: func(db *gorm.DB) error {
			return db.Exec(down22).Error
		},
	})
}
