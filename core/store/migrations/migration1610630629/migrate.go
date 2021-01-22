package migration1610630629

import "github.com/jinzhu/gorm"

// Migrate makes the explicit the pre-existing
// implicit assumption that lowercase external initiator names are unique
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE external_initiators DROP CONSTRAINT external_initiators_name_key;
		CREATE UNIQUE INDEX external_initiators_name_key ON external_initiators (lower(name));
	`).Error
}
