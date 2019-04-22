package migration1549496047

import (
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func Migrate(tx *gorm.DB) error {
	return tx.AutoMigrate(&models.Key{}).Error
}
