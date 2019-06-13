package migration1560433987

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Migrate changes the confirmations column to be nullable and default to null.
// Must be done in a roundabout way because SQLite does not support altering a column's
// default or nullability.
func Migrate(tx *gorm.DB) error {
	err := tx.Exec(`
		ALTER TABLE task_runs ADD COLUMN confirmations_new1560433987 integer;
		UPDATE task_runs SET confirmations_new1560433987 = confirmations;
		ALTER TABLE task_runs RENAME COLUMN confirmations TO confirmations_old1560433987;
		ALTER TABLE task_runs RENAME COLUMN confirmations_new1560433987 TO confirmations;
	`).Error

	return errors.Wrap(err, "failed to migrate old TaskRuns confirmations to be nullable")
}
