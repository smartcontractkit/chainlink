package migration1602180905

import "github.com/jinzhu/gorm"

const up = `
ALTER TABLE job_specs
ADD COLUMN name VARCHAR(255) UNIQUE;
`

const down = `
ALTER TABLE job_specs REMOVE FIELD name;
`

func Migrate(tx *gorm.DB) error {
	return tx.Exec(up).Error
}

func Rollback(tx *gorm.DB) error {
	return tx.Exec(down).Error
}
