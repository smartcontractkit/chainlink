package migration1606141477

import "github.com/jinzhu/gorm"

// Migrate makes a bundle of foreign keys deferrable
// This does not change any existing behaviour, but it does make certain types of test fixture much easier to work with
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
        ALTER TABLE pipeline_task_runs ALTER CONSTRAINT pipeline_task_runs_pipeline_run_id_fkey DEFERRABLE INITIALLY IMMEDIATE;
        ALTER TABLE pipeline_task_runs ALTER CONSTRAINT pipeline_task_runs_pipeline_task_spec_id_fkey DEFERRABLE INITIALLY IMMEDIATE;
        ALTER TABLE pipeline_task_specs ALTER CONSTRAINT pipeline_task_specs_pipeline_spec_id_fkey DEFERRABLE INITIALLY IMMEDIATE;
		ALTER TABLE pipeline_task_specs ALTER CONSTRAINT pipeline_task_specs_successor_id_fkey DEFERRABLE INITIALLY IMMEDIATE;
		ALTER TABLE pipeline_runs ALTER CONSTRAINT pipeline_runs_pipeline_spec_id_fkey DEFERRABLE INITIALLY IMMEDIATE;
    `).Error
}
