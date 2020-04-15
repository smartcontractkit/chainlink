package migration1585908150

import (
	"github.com/jinzhu/gorm"
)

// Migrate changes all json columns to be jsonb
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE initiators ALTER COLUMN params TYPE jsonb USING params::jsonb;
		ALTER TABLE keys ALTER COLUMN json TYPE jsonb USING json::jsonb;
		ALTER TABLE run_requests ALTER COLUMN request_params DROP DEFAULT;
		ALTER TABLE run_requests ALTER COLUMN request_params TYPE jsonb USING request_params::jsonb;
		ALTER TABLE run_requests ALTER COLUMN request_params SET DEFAULT '{}'::jsonb;
		ALTER TABLE run_results ALTER COLUMN data TYPE jsonb USING data::jsonb;
		ALTER TABLE task_specs ALTER COLUMN params TYPE jsonb USING params::jsonb;
	`).Error
}
