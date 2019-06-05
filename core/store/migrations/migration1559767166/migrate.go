package migration1559767166

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&models.TaskRun{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate TaskRun")
	}
	return nil
}
