package migration1589470036

import (
	"github.com/jinzhu/gorm"
)

// Migrate converts keys.address from text to bytea
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	ALTER TABLE keys ALTER COLUMN address TYPE bytea USING decode(substring(address from 3), 'hex');
	-- it's no longer necessary since the key is stored as a binary
	DROP INDEX idx_unique_case_insensitive_keys_addresses;
	ALTER TABLE encumbrances ALTER COLUMN aggregator TYPE bytea USING decode(substring(aggregator from 3), 'hex');
	`).Error
}
