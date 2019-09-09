package migration1565139192

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func Migrate(tx *gorm.DB) error {
	if err := tx.Exec(`ALTER TABLE job_specs ADD min_payment varchar(255)`).Error; err != nil {
		return errors.Wrap(err, "failed to add MinPayment to JobSpec")
	}
	return nil
}
