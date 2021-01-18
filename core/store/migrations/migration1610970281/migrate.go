package migration1610970281

import "github.com/jinzhu/gorm"

// Migrate adds the pipeline_queue table for better pipeline runner performance
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TABLE pipeline_queue (
			pipeline_task_run_id bigint REFERENCES pipeline_task_runs (id) NOT NULL,
			predecessor_task_run_ids bigint[] NOT NULL,
		);

		CREATE UNIQUE INDEX idx_pipeline_queue_task_run_ids ON pipeline_queue (pipeline_task_run_id);
		CREATE INDEX idx_pipeline_queue_predecessor_task_run_ids ON pipeline_queue USING GIN(predecessor_task_run_ids);
	`).Error
}
