package migration1588088353

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds a case-insensitive unique index on keys.address to avoid potential duplicates with mismatching capitalisation
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	CREATE UNIQUE INDEX idx_unique_case_insensitive_keys_addresses ON keys(lower(address))
	`).Error
}
