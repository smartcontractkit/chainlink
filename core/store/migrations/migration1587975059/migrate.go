package migration1587975059

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
	BlockHeight  int64       `gorm:"not null;unique_index:idx_unique_log_consumption"`
	ConsumerType string      `gorm:"not null;unique_index:idx_unique_log_consumption"`
	ConsumerID   *models.ID  `gorm:"not null;unique_index:idx_unique_log_consumption"`
	LogIndex     uint        `gorm:"not null;unique_index:idx_unique_log_consumption"`
	CreatedAt    time.Time
}

// Migrate drops the LogCursor table
// This is already out in the wild as of 0.8.2 so we cannot simply delete the old migration
func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&logConsumption{}).Error; err != nil {
		return errors.Wrap(err, "could not add block_number field to log_consumption table")
	}
	return tx.Exec(`DROP TABLE IF EXISTS log_cursors`).Error
}
