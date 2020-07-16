package migration1587591248

import (
	"github.com/jinzhu/gorm"
)

// Migrate changes all json columns to be jsonb
func Migrate(tx *gorm.DB) error {
	return tx.Exec(
		`ALTER TABLE initiators ADD COLUMN "absolute_threshold" float;`).Error
}
