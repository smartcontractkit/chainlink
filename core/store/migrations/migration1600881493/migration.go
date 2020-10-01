package migration1600881493

import "github.com/jinzhu/gorm"

const up = `
ALTER TABLE keys
ADD COLUMN is_funding BOOLEAN NOT NULL DEFAULT FALSE;
CREATE UNIQUE INDEX idx_keys_only_one_funding ON keys (is_funding) WHERE is_funding = TRUE;
`

const down = `
DROP INDEX idx_keys_only_one_funding;
ALTER TABLE keys REMOVE FIELD is_funding;
`

func Migrate(tx *gorm.DB) error {
	return tx.Exec(up).Error
}

func Rollback(tx *gorm.DB) error {
	return tx.Exec(down).Error
}
