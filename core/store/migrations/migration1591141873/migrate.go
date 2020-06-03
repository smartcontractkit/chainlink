package migration1591141873

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type FluxMonitorRoundStats struct {
	ID              uint64         `gorm:"primary key;not null;auto_increment"`
	Aggregator      common.Address `gorm:"not null"`
	RoundID         uint32         `gorm:"not null"`
	NumNewRoundLogs uint64         `gorm:"not null;default 0"`
	NumSubmissions  uint64         `gorm:"not null;default 0"`
}

func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&FluxMonitorRoundStats{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate FluxMonitorRoundStats")
	}
	return nil
}
