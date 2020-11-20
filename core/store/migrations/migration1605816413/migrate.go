package migration1605816413

import "github.com/jinzhu/gorm"

// Note this destroys the ability to have an FK on jobs_specs(LOWER(name)),
// however that's ok given we have serial PK on the table which can be used for FKs.
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
        DROP INDEX job_specs_name_index;
        CREATE UNIQUE INDEX job_specs_name_index_active on job_specs (LOWER(name)) where deleted_at is NULL;
    `).Error
}
