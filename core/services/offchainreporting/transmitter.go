package offchainreporting

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"gorm.io/gorm"
)

type txManager interface {
	CreateEthTransaction(db *gorm.DB, fromAddress, toAddress common.Address, payload []byte, gasLimit uint64, meta interface{}, strategy bulletprooftxmanager.TxStrategy) (etx bulletprooftxmanager.EthTx, err error)
}

type transmitter struct {
	txm         txManager
	db          *gorm.DB
	fromAddress common.Address
	gasLimit    uint64
	strategy    bulletprooftxmanager.TxStrategy
}

// NewTransmitter creates a new eth transmitter
func NewTransmitter(txm txManager, db *gorm.DB, fromAddress common.Address, gasLimit uint64, strategy bulletprooftxmanager.TxStrategy) Transmitter {
	return &transmitter{
		txm:         txm,
		db:          db,
		fromAddress: fromAddress,
		gasLimit:    gasLimit,
		strategy:    strategy,
	}
}

func (t *transmitter) CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte) error {
	db := t.db.WithContext(ctx)
	_, err := t.txm.CreateEthTransaction(db, t.fromAddress, toAddress, payload, t.gasLimit, nil, t.strategy)
	return errors.Wrap(err, "Skipped OCR transmission")
}

func (t *transmitter) FromAddress() common.Address {
	return t.fromAddress
}
