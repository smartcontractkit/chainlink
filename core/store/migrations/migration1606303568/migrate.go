package migration1606303568

import "github.com/jinzhu/gorm"

// Migrate adds name to v2 job specs
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		ALTER TABLE jobs ADD COLUMN name VARCHAR(255), ADD COLUMN schema_version INT, ADD COLUMN type VARCHAR(255);

		UPDATE jobs SET schema_version = 1, type = 'offchainreporting';
		
		ALTER TABLE jobs ALTER COLUMN schema_version SET NOT NULL, ALTER COLUMN type SET NOT NULL,
		ADD CONSTRAINT chk_schema_version CHECK (schema_version > 0),
		ADD CONSTRAINT chk_type CHECK (type != '');
    `).Error
}
