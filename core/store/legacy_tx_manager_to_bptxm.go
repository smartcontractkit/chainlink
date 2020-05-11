package store

import (
	"github.com/jinzhu/gorm"
)

// This file handles migration from the old TxManager
// Note that migration is IRREVERSIBLE
// This file can be deleted after legacy tx manager is deleted

func migrateFromLegacyTxManager(s *Store) error {
	// Set nonce and config (one-way switch)
	if err := setNonceFromLegacyTxManager(s.GetRawDB()); err != nil {
		return err
	}
	return s.Config.PermanentlySetBulletproofTxManagerEnabled()
}

func setNonceFromLegacyTxManager(db *gorm.DB) error {
	return db.Exec(`
	UPDATE keys
	SET next_nonce = COALESCE((
		SELECT max(nonce) FROM txes WHERE txes.from = keys.address
	), 0)
	WHERE next_nonce = 0;
	`).Error
}
