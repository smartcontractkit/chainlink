package offchainreporting

import (
	"context"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"gorm.io/gorm"
)

type transmitter struct {
	db                         *gorm.DB
	fromAddress                gethCommon.Address
	gasLimit                   uint64
	maxUnconfirmedTransactions uint64
}

// NewTransmitter creates a new eth transmitter
func NewTransmitter(db *gorm.DB, fromAddress gethCommon.Address, gasLimit, maxUnconfirmedTransactions uint64) Transmitter {
	return &transmitter{
		db:                         db,
		fromAddress:                fromAddress,
		gasLimit:                   gasLimit,
		maxUnconfirmedTransactions: maxUnconfirmedTransactions,
	}
}

func (t *transmitter) CreateEthTransaction(ctx context.Context, toAddress gethCommon.Address, payload []byte) error {
	db := t.db.WithContext(ctx)
	_, err := bulletprooftxmanager.CreateEthTransaction(db, t.fromAddress, toAddress, payload, t.gasLimit, t.maxUnconfirmedTransactions)
	return errors.Wrap(err, "Skipped OCR transmission")
}

func (t *transmitter) FromAddress() gethCommon.Address {
	return t.fromAddress
}
