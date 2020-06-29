package migration1574659987

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/store/models/vrfkey"
)

// Migrate adds VRF proving-key table, and related subtables.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&vrfkey.EncryptedSecretKey{}).Error; err != nil {
		return errors.Wrap(err, "failed to create VRF proving-key table")
	}
	return nil
}
