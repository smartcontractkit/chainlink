package migration1560924400

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// Migrate adds the BlockHash column to the RunRequest and TxReceipt tables
func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&RunRequest{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate RunRequest")
	}

	return nil
}

// RunRequest stores the fields used to initiate the parent job run.
// This migration introduces the BlockHash column onto the table.
type RunRequest struct {
	ID        uint `gorm:"primary_key"`
	RequestID *string
	TxHash    *common.Hash
	BlockHash *common.Hash
	Requester *common.Address
	CreatedAt time.Time
}
