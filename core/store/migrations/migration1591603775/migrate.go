package migration1591603775

import (
	"github.com/jinzhu/gorm"
)

// Migrate changes the index on txes.hash to be non-unique again to allow existing buggy code to continue to function
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		DROP INDEX IF EXISTS idx_txes_hash;
		CREATE INDEX idx_txes_hash ON txes (hash);
    `).Error
}
