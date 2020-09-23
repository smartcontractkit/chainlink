package migration1598972982

import (
	"github.com/jinzhu/gorm"
)

const up = `
ALTER TABLE keys
ADD FIELD is_rescue BOOLEAN NOT NULL DEFAULT FALSE;
CREATE UNIQUE INDEX only_one_rescue ON keys (is_rescue) WHERE is_rescue = TRUE;
`

const down = `
DROP INDEX only_one_rescue;
ALTER TABLE keys REMOVE FIELD is_rescue;
`

func Migrate(tx *gorm.DB) error {
	return tx.Exec(up).Error
}

func Rollback(tx *gorm.DB) error {
	return tx.Exec(down).Error
}
