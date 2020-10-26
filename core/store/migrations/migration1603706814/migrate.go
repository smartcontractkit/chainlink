package migration1603706814

import "github.com/jinzhu/gorm"

const up = `
ALTER TABLE offchainreporting_oracle_specs
ADD COLUMN name text UNIQUE NOT NULL;
`

func Migrate(tx *gorm.DB) error {
	return tx.Exec(up).Error
}
