package migration1565210496

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/dbutil"
)

// Migrate optimizes the JobRuns table to reduce the cost of IDs
func Migrate(tx *gorm.DB) error {
	if dbutil.IsPostgres(tx) {
		if err := tx.Exec(`
	`).Error; err != nil {
			return errors.Wrap(err, "failed to update ids on job_runs")
		}
	}

	return nil
}
