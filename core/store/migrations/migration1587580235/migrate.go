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
	BlockHash    common.Hash `gorm:"index;not null"`
	ConsumerType string      `gorm:"index;not null"`
	ConsumerID   *models.ID  `gorm:"index;not null"`
	LogIndex     uint        `gorm:"index;not null"`
	CreatedAt    time.Time
}

// TODO - RYAN add uniqueness constraint on all 4 columns

// Migrate adds the LogConsumption table
func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&logConsumption{}).Error; err != nil {
		return errors.Wrap(err, "could not add log_consumption table")
	}
	return nil
}
