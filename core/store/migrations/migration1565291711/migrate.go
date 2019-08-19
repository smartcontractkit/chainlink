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
ALTER TABLE job_specs ADD COLUMN "id_uuid" uuid;
UPDATE job_specs
SET
	"id_uuid" = CAST("id" as uuid);
ALTER TABLE job_specs DROP CONSTRAINT "job_specs_pkey" CASCADE;
	`).Error; err != nil {
			return errors.Wrap(err, "failed to add id_uuid on job_runs")
		}

		if err := tx.Exec(`
ALTER TABLE job_runs ADD COLUMN "job_spec_id_uuid" uuid;
UPDATE job_runs
SET
	"job_spec_id_uuid" = CAST("job_spec_id" as uuid);
	`).Error; err != nil {
			return errors.Wrap(err, "failed to add job_spec_id_uuid on job_runs")
		}

		if err := tx.Exec(`
ALTER TABLE task_specs ADD COLUMN "job_spec_id_uuid" uuid;
UPDATE task_specs
SET
	"job_spec_id_uuid" = CAST("job_spec_id" as uuid);
	`).Error; err != nil {
			return errors.Wrap(err, "failed to add job_spec_id_uuid on task_specs")
		}

		if err := tx.Exec(`
ALTER TABLE service_agreements ADD COLUMN "job_spec_id_uuid" uuid;
UPDATE service_agreements
SET
	"job_spec_id_uuid" = CAST("job_spec_id" as uuid);
	`).Error; err != nil {
			return errors.Wrap(err, "failed to add job_spec_id_uuid on service_agreements")
		}

		if err := tx.Exec(`
ALTER TABLE link_earned DROP COLUMN "job_spec_id";
	`).Error; err != nil {
			return errors.Wrap(err, "failed to drop job_spec_id from link_earned")
		}

		if err := tx.Exec(`
ALTER TABLE job_specs DROP COLUMN "id";
ALTER TABLE job_specs RENAME COLUMN "id_uuid" TO "id";
ALTER TABLE job_specs ADD CONSTRAINT "job_spec_pkey" PRIMARY KEY ("id");
	`).Error; err != nil {
			return errors.Wrap(err, "failed to rename id on job_specs")
		}

		if err := tx.Exec(`
ALTER TABLE job_runs DROP COLUMN "job_spec_id";
ALTER TABLE job_runs RENAME COLUMN "job_spec_id_uuid" TO "job_spec_id";

ALTER TABLE job_runs ADD CONSTRAINT "job_runs_job_spec_id_fkey" FOREIGN KEY ("job_spec_id") REFERENCES job_specs ("id");
	`).Error; err != nil {
			return errors.Wrap(err, "failed to update job_spec_id id on job_runs")
		}

		if err := tx.Exec(`
ALTER TABLE task_specs DROP COLUMN "job_spec_id";
ALTER TABLE task_specs RENAME COLUMN "job_spec_id_uuid" TO "job_spec_id";

ALTER TABLE task_specs ADD CONSTRAINT "task_specs_job_spec_id_fkey" FOREIGN KEY ("job_spec_id") REFERENCES job_specs ("id");
	`).Error; err != nil {
			return errors.Wrap(err, "failed to update job_spec_id id on task_specs")
		}

		if err := tx.Exec(`
ALTER TABLE service_agreements DROP COLUMN "job_spec_id";
ALTER TABLE service_agreements RENAME COLUMN "job_spec_id_uuid" TO "job_spec_id";

ALTER TABLE service_agreements ADD CONSTRAINT "service_agreements_job_spec_id_fkey" FOREIGN KEY ("job_spec_id") REFERENCES job_specs ("id");
	`).Error; err != nil {
			return errors.Wrap(err, "failed to update job_spec_id id on service_agreements")
		}

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
ALTER TABLE link_earned ADD COLUMN "job_run_id_uuid" uuid;
UPDATE link_earned
SET
	"job_run_id_uuid" = CAST("job_run_id" as uuid);
	`).Error; err != nil {
			return errors.Wrap(err, "failed to add job_run_id_uuid on link_earned")
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

ALTER TABLE task_runs ADD CONSTRAINT "task_runs_job_run_id_fkey" FOREIGN KEY ("job_run_id") REFERENCES job_runs ("id");
	`).Error; err != nil {
			return errors.Wrap(err, "failed to update job_run_id id on task_runs")
		}

		if err := tx.Exec(`
ALTER TABLE link_earned DROP COLUMN "job_run_id";
ALTER TABLE link_earned RENAME COLUMN "job_run_id_uuid" TO "job_run_id";

ALTER TABLE link_earned ADD CONSTRAINT "link_earned_job_run_id_fkey" FOREIGN KEY ("job_run_id") REFERENCES job_runs ("id");
	`).Error; err != nil {
			return errors.Wrap(err, "failed to update job_run_id id on link_earned")
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
