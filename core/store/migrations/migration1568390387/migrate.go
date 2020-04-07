package migration1568390387

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Encumbrance struct {
	ID                     uint         `gorm:"primary_key;auto_increment"`
	Payment                *assets.Link `gorm:"type:varchar(255)"`
	Expiration             uint64
	EndAt                  time.Time
	Oracles                string  `gorm:"type:text"`
	Aggregator             string  `gorm:"not null"`
	AggInitiateJobSelector [4]byte `gorm:"not null"`
	AggFulfillSelector     [4]byte `gorm:"not null"`
}

// Migrate amends the encumbrances table to include the aggregator contact details
func Migrate(tx *gorm.DB) error {
	// This table is behind the development flag and so any records are safe to remove
	if err := tx.Exec(`DROP TABLE "encumbrances";`).Error; err != nil {
		return errors.Wrap(err, "could not drop Encumbrances table")
	}

	if err := tx.AutoMigrate(&Encumbrance{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate Encumbrance")
	}

	return nil
}
