package migration1606910307

import "github.com/jinzhu/gorm"

func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
        DROP INDEX job_specs_name_index_active;
    `).Error
}
