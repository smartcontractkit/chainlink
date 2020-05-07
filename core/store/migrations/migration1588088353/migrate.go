package migration1588088353

import (
	"github.com/jinzhu/gorm"
)

// Migrate converts keys.address into citext to avoid potential duplicates with mismatching capitalisation
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	CREATE EXTENSION IF NOT EXISTS citext;
	ALTER TABLE keys ALTER COLUMN address TYPE citext;
	`).Error
}
