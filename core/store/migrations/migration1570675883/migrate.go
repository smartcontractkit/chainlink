package migration1570675883

import (
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type JobRun struct {
	models.JobRun
	Overrides models.JSON
}

func (JobRun) TableName() string {
	return "job_runs"
}

func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
ALTER TABLE job_runs ADD COLUMN "overrides" text;
UPDATE job_runs
SET "overrides" = (
	SELECT data
	FROM run_results
	WHERE overrides_id = run_results.id
);
DELETE FROM run_results
WHERE id IN (
	SELECT overrides_id
	FROM job_runs
);`).Error
}
