package migration1610970281

import "github.com/jinzhu/gorm"

// Migrate adds the pipeline_queue table for better pipeline runner performance
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TABLE pipeline_queue (
			id BIGSERIAL PRIMARY KEY,
			pipeline_task_run_ids bigint[] NOT NULL CHECK (cardinality(pipeline_task_run_ids) > 0)
		);
	`).Error
}
