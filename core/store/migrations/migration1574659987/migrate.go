package migration1574659987

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"chainlink/core/store/models"
)

// Migrate adds VRF proving-key table, and related subtables.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.EncryptedSecretVRFKey{}).Error; err != nil {
		return errors.Wrap(err, "failed to create VRF proving-key table")
	}
	return nil
}
