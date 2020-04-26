package migration1587580235

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type logConsumption struct {
	ID           *models.ID  `gorm:"primary_key"`
	BlockHash    common.Hash `gorm:"not null;unique_index:idx_unique_log_consumption"`
	ConsumerType string      `gorm:"not null;unique_index:idx_unique_log_consumption"`
	ConsumerID   *models.ID  `gorm:"not null;unique_index:idx_unique_log_consumption"`
	LogIndex     uint        `gorm:"not null;unique_index:idx_unique_log_consumption"`
	CreatedAt    time.Time
}

// Migrate adds the LogConsumption table
func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&logConsumption{}).Error; err != nil {
		return errors.Wrap(err, "could not add log_consumption table")
	}
	return nil
}
