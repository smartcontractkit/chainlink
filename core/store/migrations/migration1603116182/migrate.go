package migration1603116182

import "github.com/jinzhu/gorm"

const up = `
ALTER INDEX idx_pipeline_task_runs RENAME TO idx_pipeline_task_runs_created_at;
ALTER TABLE pipeline_task_runs DROP CONSTRAINT chk_pipeline_task_run_fsm;
ALTER TABLE pipeline_task_runs ADD CONSTRAINT chk_pipeline_task_run_fsm CHECK (
	type != 'result' AND (
		finished_at IS NULL AND error IS NULL AND output IS NULL OR (
			finished_at IS NOT NULL AND NOT (error IS NOT NULL AND OUTPUT IS NOT NULL)
		)
	)
	OR
	type = 'result' AND (
		output IS NULL AND error IS NULL AND finished_at IS NULL
		OR
		output IS NOT NULL AND error IS NOT NULL AND finished_at IS NOT NULL
	)
);
`

func Migrate(tx *gorm.DB) error {
	return tx.Exec(up).Error
}
