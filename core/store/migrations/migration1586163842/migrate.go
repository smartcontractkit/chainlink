package migration1586163842

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds foreign keys that were missing
func Migrate(tx *gorm.DB) error {
	// Add a few more useful indexes while we are here. This also speeds up
	// some queries in this migration
	err := tx.Exec(`
	CREATE INDEX idx_task_specs_job_spec_id ON task_specs (job_spec_id);
	CREATE INDEX idx_job_runs_run_request_id ON job_runs (run_request_id);
	CREATE INDEX idx_job_runs_result_id ON job_runs (result_id);
	CREATE INDEX idx_job_runs_initiator_id ON job_runs (initiator_id);
	CREATE INDEX idx_task_runs_result_id ON task_runs (result_id);
	`).Error
	if err != nil {
		return err
	}

	// Need to cast initiators.job_spec_id to UUID to be compatible with referenced key
	err = tx.Exec(`
	ALTER TABLE initiators ALTER COLUMN job_spec_id TYPE uuid USING job_spec_id::uuid;
	`).Error
	if err != nil {
		return err
	}

	// Before we were using 0 as the null value, need to set to explicit null for FK
	err = tx.Exec(`
	UPDATE job_runs SET result_id = NULL WHERE result_id = 0;
	UPDATE job_runs SET run_request_id = NULL WHERE run_request_id = 0;
	UPDATE job_runs SET initiator_id = NULL WHERE initiator_id = 0;
	UPDATE task_runs SET result_id = NULL WHERE result_id = 0;
	UPDATE task_runs SET task_spec_id = NULL WHERE task_spec_id = 0;
	`).Error
	if err != nil {
		return err
	}

	// This was assumed by the code and probably would have crashed/hung if this assumption was violated.
	// Let's make that explicit.
	err = tx.Exec(`
	DELETE FROM job_runs WHERE initiator_id IS NULL;
	ALTER TABLE job_runs ALTER COLUMN initiator_id SET NOT NULL;
	`).Error
	if err != nil {
		return err
	}

	// Add the foreign keys
	err = tx.Exec(`
	ALTER TABLE initiators ADD CONSTRAINT fk_initiators_job_spec_id FOREIGN KEY (job_spec_id) REFERENCES job_specs (id) ON DELETE RESTRICT;
	ALTER TABLE job_runs ADD CONSTRAINT fk_job_runs_result_id FOREIGN KEY (result_id) REFERENCES run_results (id) ON DELETE CASCADE;
	ALTER TABLE job_runs ADD CONSTRAINT fk_job_runs_run_request_id FOREIGN KEY (run_request_id) REFERENCES run_requests (id) ON DELETE CASCADE;
	ALTER TABLE job_runs ADD CONSTRAINT fk_job_runs_initiator_id FOREIGN KEY (initiator_id) REFERENCES initiators (id) ON DELETE CASCADE;
	ALTER TABLE service_agreements ADD CONSTRAINT fk_service_agreements_encumbrance_id FOREIGN KEY (encumbrance_id) REFERENCES encumbrances (id) ON DELETE RESTRICT;
	ALTER TABLE task_runs ADD CONSTRAINT fk_task_runs_result_id FOREIGN KEY (result_id) REFERENCES run_results (id) ON DELETE CASCADE;
	ALTER TABLE task_runs ADD CONSTRAINT fk_task_runs_task_spec_id FOREIGN KEY (task_spec_id) REFERENCES task_specs (id) ON DELETE CASCADE;
	`).Error
	if err != nil {
		return err
	}
	return nil
}
