package migration1586369235

import (
	"github.com/jinzhu/gorm"
)

// Migrate changes text to binary where appropriate and uses numeric type for very large integers
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
	ALTER TABLE run_requests ALTER COLUMN request_id TYPE bytea USING decode(substring(request_id from 3), 'hex');
	ALTER TABLE tx_attempts ALTER COLUMN signed_raw_tx TYPE bytea USING decode(substring(signed_raw_tx from 3), 'hex');
	ALTER TABLE txes ALTER COLUMN signed_raw_tx TYPE bytea USING decode(substring(signed_raw_tx from 3), 'hex');

	ALTER TABLE tx_attempts ALTER COLUMN gas_price TYPE numeric(78, 0) USING gas_price::numeric;
	ALTER TABLE txes ALTER COLUMN gas_price TYPE numeric(78, 0) USING gas_price::numeric;
	ALTER TABLE txes ALTER COLUMN value TYPE numeric(78, 0) USING value::numeric;
	ALTER TABLE encumbrances ALTER COLUMN payment TYPE numeric(78, 0) USING payment::numeric;
	`).Error
}
