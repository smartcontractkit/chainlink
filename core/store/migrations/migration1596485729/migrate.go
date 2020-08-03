package migration1596485729

import (
	"github.com/jinzhu/gorm"
)

// Migrate adds last_used to keys
func Migrate(tx *gorm.DB) error {
	// TODO - RYAN
	return tx.Exec(``).Error
}
