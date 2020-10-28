package migration1603816329

import "github.com/jinzhu/gorm"

const up = `
CREATE INDEX idx_pipeline_task_runs_optimise_find_predecessor_unfinished_runs ON pipeline_task_runs (pipeline_task_spec_id, pipeline_run_id) INCLUDE (id, finished_at);
CREATE INDEX idx_pipeline_task_runs_optimise_find_results ON pipeline_task_runs (pipeline_run_id);
`

// Migrate adds a couple of specific indexes that are needed to improve job pipeline performance
func Migrate(tx *gorm.DB) error {
	return tx.Exec(up).Error
}
