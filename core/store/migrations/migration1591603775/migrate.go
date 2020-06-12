package migration1591603775

import (
	"github.com/jinzhu/gorm"
)

// Migrate changes the index on txes.hash to be non-unique again to allow existing buggy code to continue to function
// It also drops the unique nonce index
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		DROP INDEX IF EXISTS idx_txes_hash;
		CREATE INDEX idx_txes_hash ON txes (hash);
		DROP INDEX IF EXISTS idx_txes_unique_nonces_per_account;
		CREATE INDEX IF NOT EXISTS idx_txes_nonce ON txes(nonce);
    `).Error
}
