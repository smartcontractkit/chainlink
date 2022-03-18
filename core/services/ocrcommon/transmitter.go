package ocrcommon

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type txManager interface {
	CreateEthTransaction(newTx txmgr.NewTx, qopts ...pg.QOpt) (etx txmgr.EthTx, err error)
}

type Transmitter interface {
	CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte) error
	FromAddress() common.Address
}

type transmitter struct {
	txm         txManager
	fromAddress common.Address
	gasLimit    uint64
	strategy    txmgr.TxStrategy
	checker     txmgr.TransmitCheckerSpec
}

// NewTransmitter creates a new eth transmitter
func NewTransmitter(txm txManager, fromAddress common.Address, gasLimit uint64, strategy txmgr.TxStrategy, checker txmgr.TransmitCheckerSpec) Transmitter {
	return &transmitter{
		txm:         txm,
		fromAddress: fromAddress,
		gasLimit:    gasLimit,
		strategy:    strategy,
		checker:     checker,
	}
}

func (t *transmitter) CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte) error {
	_, err := t.txm.CreateEthTransaction(txmgr.NewTx{
		FromAddress:    t.fromAddress,
		ToAddress:      toAddress,
		EncodedPayload: payload,
		GasLimit:       t.gasLimit,
		Strategy:       t.strategy,
		Checker:        t.checker,
	}, pg.WithParentCtx(ctx))
	return errors.Wrap(err, "Skipped OCR transmission")
}

func (t *transmitter) FromAddress() common.Address {
	return t.fromAddress
}
