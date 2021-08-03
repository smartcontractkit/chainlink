package migrations

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const up54 = `
ALTER TABLE log_broadcasts DROP COLUMN job_id;
DROP TABLE service_agreements;
DROP TABLE eth_task_run_txes;
DROP TABLE task_runs;
DROP TABLE task_specs;
DROP TABLE flux_monitor_round_stats;
DROP TABLE job_runs;
DROP TABLE job_spec_errors;
DROP TABLE initiators;
DROP TABLE job_specs;

DROP TABLE run_results;
DROP TABLE run_requests;
DROP TABLE sync_events;

ALTER TABLE log_broadcasts RENAME COLUMN job_id_v2 TO job_id;
ALTER TABLE job_spec_errors_v2 RENAME TO job_spec_errors;
`

func init() {
	Migrations = append(Migrations, &Migration{
		ID: "0054_remove_legacy_pipeline",
		Migrate: func(db *gorm.DB) error {
			if err := checkNoLegacyJobs(db); err != nil {
				return err
			}
			return db.Exec(up54).Error
		},
		Rollback: func(db *gorm.DB) error {
			return errors.New("irreversible migration")
		},
	})
}

func checkNoLegacyJobs(db *gorm.DB) error {
	var count int
	if err := db.Raw(`SELECT COUNT(*) FROM job_specs WHERE deleted_at IS NULL`).Scan(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.Errorf("cannot migrate; this release removes support for legacy job specs but there are still %d in the database. Please migrate these job specs to the V2 pipeline and make sure all V1 job_specs are deleted or archived, then run the migration again. Migration instructions found here: https://docs.chain.link/docs/jobs/migration-v1-v2/", count)
	}
	return nil

}
