package migration1559081901

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func Migrate(tx *gorm.DB) error {
	// If any of these tables somehow ended up with a duplicate hash, moving to a
	// unique constraint puts us in a bit of a difficult position. Do we delete
	// old txes? Use some algorithm to choose one to keep? For now, archive, and
	// make new, stricter tables.
	if err := tx.Exec(
		`DROP INDEX idx_txes_from;
		 DROP INDEX idx_txes_nonce;
		 DROP INDEX idx_tx_attempts_tx_id;
		 DROP INDEX idx_tx_attempts_created_at;
		 ALTER TABLE txes RENAME TO txes_archive;
		 ALTER TABLE tx_attempts RENAME TO tx_attempts_archive;
		 `,
	).Error; err != nil {
		return errors.Wrap(err, "failed to drop txes and txattempts")
	}
	if err := tx.AutoMigrate(&models.Tx{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Tx")
	}
	if err := tx.AutoMigrate(&models.TxAttempt{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate TxAttempt")
	}
	return nil
}
