package migration1553029703

import (
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/store/orm"
)

// Migration is the singleton type for this migration
type Migration struct{}

// Timestamp returns the epoch for this migration
func (m Migration) Timestamp() string {
	return "1553029703"
}

// Migrate adds the sync_events table
func (m Migration) Migrate(orm *orm.ORM) error {
	return orm.DB.AutoMigrate(&models.SyncEvent{}).Error
}
