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
	txm                         txManager
	fromAddresses               []common.Address
	nextFromAddressIndex        int
	gasLimit                    uint32
	effectiveTransmitterAddress common.Address
	strategy                    txmgr.TxStrategy
	checker                     txmgr.TransmitCheckerSpec
}

// NewTransmitter creates a new eth transmitter
func NewTransmitter(txm txManager, fromAddresses []common.Address, gasLimit uint32, effectiveTransmitterAddress common.Address, strategy txmgr.TxStrategy, checker txmgr.TransmitCheckerSpec) Transmitter {
	return &transmitter{
		txm:                         txm,
		fromAddresses:               fromAddresses,
		gasLimit:                    gasLimit,
		effectiveTransmitterAddress: effectiveTransmitterAddress,
		strategy:                    strategy,
		checker:                     checker,
	}
}

func (t *transmitter) CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte) error {
	_, err := t.txm.CreateEthTransaction(txmgr.NewTx{
		FromAddress:      t.FromAddressForTransaction(),
		ToAddress:        toAddress,
		EncodedPayload:   payload,
		GasLimit:         t.gasLimit,
		ForwarderAddress: t.forwarderAddress(),
		Strategy:         t.strategy,
		Checker:          t.checker,
	}, pg.WithParentCtx(ctx))
	return errors.Wrap(err, "Skipped OCR transmission")
}

func (t *transmitter) FromAddress() common.Address {
	return t.effectiveTransmitterAddress
}

func (t *transmitter) FromAddressForTransaction() common.Address {
	// Use Round-Robin to select the next fromAddress.
	nextFromAddress := t.fromAddresses[t.nextFromAddressIndex]

	// Only apply round-robin logic for multiple sending keys.
	if len(t.fromAddresses) > 1 {
		t.nextFromAddressIndex++
		if t.nextFromAddressIndex >= len(t.fromAddresses) {
			t.nextFromAddressIndex = 0
		}
	}

	return nextFromAddress
}

func (t *transmitter) forwarderAddress() common.Address {
	for _, a := range t.fromAddresses {
		if a == t.effectiveTransmitterAddress {
			return common.Address{}
		}
	}
	return t.effectiveTransmitterAddress
}
