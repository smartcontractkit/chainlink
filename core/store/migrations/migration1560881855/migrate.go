package migration1560881855

import (
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func Migrate(tx *gorm.DB) error {
	err := tx.AutoMigrate(&models.LinkEarned{}).Error
	return err
}
