package migration1602180905

import "github.com/jinzhu/gorm"

const up = `
ALTER TABLE job_specs ADD COLUMN name VARCHAR(255);
CREATE UNIQUE INDEX job_specs_name_index on job_specs (LOWER(name));
`

const down = `
DROP INDEX job_specs_name_index;
ALTER TABLE job_specs REMOVE FIELD name;
`

func Migrate(tx *gorm.DB) error {
	return tx.Exec(up).Error
}

func Rollback(tx *gorm.DB) error {
	return tx.Exec(down).Error
}
