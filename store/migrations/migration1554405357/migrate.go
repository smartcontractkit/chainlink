package migration1554405357

import (
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Migration is the singleton type for this migration
type Migration struct{}

// Migrate adds the sync_events table
func (m Migration) Migrate(tx *gorm.DB) error {
	return tx.AutoMigrate(&models.ExternalInitiator{}).Error
}
