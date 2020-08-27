package migration1598521075

import (
	"github.com/jinzhu/gorm"
)

// Migrate creates the p2p_peerstore table
func Migrate(tx *gorm.DB) error {
	return tx.Exec(`
		CREATE TABLE IF NOT EXISTS p2p_peerstore (
			key TEXT PRIMARY KEY,
			data BYTEA NOT NULL
		);
	`).Error
}
