package offchainreporting

import (
	"context"

	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/utils"
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
	value := *assets.NewEth(0)
	return utils.JustError(
		bulletprooftxmanager.CreateTxIfFunded(ctx, t.db, t.fromAddress, toAddress, value, payload, t.gasLimit, t.maxUnconfirmedTransactions),
	)
}

func (t *transmitter) FromAddress() gethCommon.Address {
	return t.fromAddress
}
