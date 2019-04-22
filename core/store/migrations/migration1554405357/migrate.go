package migration1554405357

import (
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// Migrate adds the sync_events table
func Migrate(tx *gorm.DB) error {
	return tx.AutoMigrate(&models.ExternalInitiator{}).Error
}
