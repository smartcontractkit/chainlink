package migration1565291711

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/dbutil"
)

// Migrate optimizes the JobRuns table to reduce the cost of IDs
func Migrate(tx *gorm.DB) error {
	if dbutil.IsPostgres(tx) {
		if err := tx.Exec(`
ALTER TABLE job_runs ADD COLUMN "id_uuid" uuid;
UPDATE job_runs
SET
	"id_uuid" = CAST("id" as uuid);
ALTER TABLE job_runs DROP CONSTRAINT "job_runs_pkey" CASCADE;
	`).Error; err != nil {
			return errors.Wrap(err, "failed to add id_uuid on job_runs")
		}

		if err := tx.Exec(`
ALTER TABLE task_runs ADD COLUMN "job_run_id_uuid" uuid;
UPDATE task_runs
SET
	"job_run_id_uuid" = CAST("job_run_id" as uuid);
	`).Error; err != nil {
			return errors.Wrap(err, "failed to add job_run_id_uuid on task_runs")
		}

		if err := tx.Exec(`
ALTER TABLE job_runs DROP COLUMN "id";
ALTER TABLE job_runs RENAME COLUMN "id_uuid" TO "id";
ALTER TABLE job_runs ADD CONSTRAINT "job_run_pkey" PRIMARY KEY ("id");
	`).Error; err != nil {
			return errors.Wrap(err, "failed to rename id on job_runs")
		}

		if err := tx.Exec(`
ALTER TABLE task_runs DROP COLUMN "job_run_id";
ALTER TABLE task_runs RENAME COLUMN "job_run_id_uuid" TO "job_run_id";

ALTER TABLE task_runs ADD CONSTRAINT "task_runs_job_run_id_fkey" FOREIGN KEY ("job_run_id") REFERENCES job_runs ("id") MATCH FULL;
	`).Error; err != nil {
			return errors.Wrap(err, "failed to update job_run_id id on task_runs")
		}

		if err := tx.Exec(`
ALTER TABLE task_runs ADD COLUMN "id_uuid" uuid;
UPDATE task_runs SET "id_uuid" = CAST("id" as uuid);
ALTER TABLE task_runs DROP COLUMN "id";
ALTER TABLE task_runs RENAME COLUMN "id_uuid" TO "id";
	`).Error; err != nil {
			return errors.Wrap(err, "failed to update id on task_runs")
		}
	}

	return nil
}
