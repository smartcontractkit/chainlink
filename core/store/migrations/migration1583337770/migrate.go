package migration1583337770

import (
	"github.com/jinzhu/gorm"
)

// Migrate converts task_specs.params into JSONB and grandfathers in
// the legacy followRedirect behaviour
// NOTE: Re-typing columns is fairly expensive, but even very large task_spec
// tables are unlikely to exceed a few thousand rows, so migration time is
// acceptable.
func Migrate(tx *gorm.DB) error {
	err := tx.Exec(`ALTER TABLE task_specs ALTER COLUMN params TYPE jsonb USING params::jsonb`).Error
	if err != nil {
		return err
	}

	return tx.Exec(`
		UPDATE task_specs
		SET params = params || jsonb '{"followRedirects": true}'
		WHERE type IN ('httpget', 'httppost') AND NOT params ? 'followRedirects'
	`).Error
}
