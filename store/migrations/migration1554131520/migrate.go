package migration1554131520

import (
	"github.com/jinzhu/gorm"
	"github.com/smartcontractkit/chainlink/store/models"
)

// Migration is the singleton type for this migration
type Migration struct{}

// Migrate adds the initiator_runs table
func (m Migration) Migrate(tx *gorm.DB) error {
	return tx.AutoMigrate(&models.InitiatorRun{}).Error
}
