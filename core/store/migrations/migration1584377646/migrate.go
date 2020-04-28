package migration1584377646

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type logCursor struct {
	Name        string `gorm:"primary_key"`
	Initialized bool   `gorm:"not null;default true"`
	BlockIndex  uint   `gorm:"not null;default 0"`
	LogIndex    uint   `gorm:"not null;default 0"`
}

// Migrate adds the LogCursor table
func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&logCursor{}).Error; err != nil {
		return errors.Wrap(err, "could not add log_cursor table")
	}
	return nil
}
