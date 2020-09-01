package migration1598972982

import (
	"github.com/jinzhu/gorm"
)

const up = `
ALTER TABLE log_consumptions
ADD COLUMN block_number BIGINT
`

const down = `
ALTER TABLE log_consumptions
DROP COLUMN block_number
`

func Migrate(tx *gorm.DB) error {
	return tx.Exec(up).Error
}

func Rollback(tx *gorm.DB) error {
	return tx.Exec(down).Error
}
