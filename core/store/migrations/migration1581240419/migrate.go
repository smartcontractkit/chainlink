package migration1581240419

import (
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/store/dbutil"
)

// Migrate moves job_runs.overrides to run_requests.request_params
func Migrate(tx *gorm.DB) error {
	err := tx.Exec(`ALTER TABLE run_requests ADD request_params text NOT NULL DEFAULT '{}'`).Error
	if err != nil {
		return err
	}
	if dbutil.IsPostgres(tx) {
		err = tx.Exec(`
			UPDATE run_requests
			SET request_params = job_runs.overrides
			FROM job_runs
			WHERE job_runs.run_request_id = run_requests.id AND job_runs.overrides IS NOT NULL
		`).Error
	} else {
		// WARNING: This is slow on Sqlite since there is no index on job_runs.run_request_id
		err = tx.Exec(`
			UPDATE run_requests
			SET request_params = COALESCE((
				SELECT overrides
				FROM job_runs
				WHERE run_request_id = run_requests.id
			), '{}')
		`).Error
	}
	if err != nil {
		return err
	}
	if dbutil.IsPostgres(tx) {
		return tx.Exec(`ALTER TABLE job_runs DROP COLUMN IF EXISTS overrides`).Error
	}
	return nil
}
