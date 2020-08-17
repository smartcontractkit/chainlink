package migration1565210496

import (
	"github.com/smartcontractkit/chainlink/core/store/dbutil"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Migrate optimizes the JobRuns table for reduced on disk footprint
func Migrate(tx *gorm.DB) error {
	if dbutil.IsPostgres(tx) {
		if err := tx.Exec(`
ALTER TABLE job_runs ADD COLUMN "creation_height_numeric" numeric(78, 0);
ALTER TABLE job_runs ADD COLUMN "observed_height_numeric" numeric(78, 0);
UPDATE job_runs
SET
	"creation_height_numeric" = CAST("creation_height" as numeric),
	"observed_height_numeric" = CAST("observed_height" as numeric);
ALTER TABLE job_runs DROP COLUMN "creation_height";
ALTER TABLE job_runs DROP COLUMN "observed_height";
ALTER TABLE job_runs RENAME COLUMN "creation_height_numeric" TO "creation_height";
ALTER TABLE job_runs RENAME COLUMN "observed_height_numeric" TO "observed_height";
	`).Error; err != nil {
			return errors.Wrap(err, "failed to change height columns on job_runs")
		}
	}

	return nil
}
