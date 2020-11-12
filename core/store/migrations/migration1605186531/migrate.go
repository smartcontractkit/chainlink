package migration1605186531

import "github.com/jinzhu/gorm"

// Migrate adds a couple of indexes to pipeline_runs
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE INDEX idx_pipeline_runs_finished_at ON pipeline_runs USING BRIN (finished_at);
		CREATE INDEX idx_pipeline_task_runs_finished_at ON pipeline_task_runs USING BRIN (finished_at);
		CREATE INDEX idx_job_spec_errors_v2_created_at ON job_spec_errors_v2 USING BRIN (created_at);
		CREATE INDEX idx_job_spec_errors_v2_finished_at ON job_spec_errors_v2 USING BRIN (updated_at);
    `).Error
}
