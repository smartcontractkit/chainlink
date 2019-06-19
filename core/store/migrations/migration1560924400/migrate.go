package migration1560924400

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// Migrate adds the BlockHash column to the RunRequest and TxReceipt tables
func Migrate(tx *gorm.DB) error {
	if err := tx.AutoMigrate(&RunRequest{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate RunRequest")
	}

	if err := tx.AutoMigrate(&TxReceipt{}).Error; err != nil {
		return errors.Wrap(err, "failed to auto migrate TxReceipt")
	}

	return nil
}

type RunRequest struct {
	ID        uint `gorm:"primary_key"`
	RequestID *string
	TxHash    *common.Hash
	BlockHash *common.Hash
	Requester *common.Address
	CreatedAt time.Time
}

type TxReceipt struct {
	BlockNumber *models.Big  `json:"blockNumber" gorm:"type:numeric"`
	BlockHash   common.Hash  `json:"blockHash"`
	Hash        common.Hash  `json:"transactionHash"`
	Logs        []models.Log `json:"logs"`
}
