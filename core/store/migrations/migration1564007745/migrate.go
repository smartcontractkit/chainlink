package migration1564007745

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&models.Configuration{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Configuration")
	}
	return nil
}
