package ocrcommon

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type transmitter struct {
	txm                         txManager
	fromAddress                 common.Address
	gasLimit                    uint32
	effectiveTransmitterAddress common.Address
	strategy                    txmgr.TxStrategy
	checker                     txmgr.TransmitCheckerSpec
}

// NewDefaultTransmitter creates a new eth transmitter using a default tx manager
func NewDefaultTransmitter(txm txManager, fromAddress common.Address, gasLimit uint32, effectiveTransmitterAddress common.Address, strategy txmgr.TxStrategy, checker txmgr.TransmitCheckerSpec) Transmitter {
	return &transmitter{
		txm:                         txm,
		fromAddress:                 fromAddress,
		gasLimit:                    gasLimit,
		effectiveTransmitterAddress: effectiveTransmitterAddress,
		strategy:                    strategy,
		checker:                     checker,
	}
}

func (t *transmitter) CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte) error {
	_, err := t.txm.CreateEthTransaction(txmgr.NewTx{
		FromAddress:      t.fromAddress,
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

func (t *transmitter) forwarderAddress() common.Address {
	if t.effectiveTransmitterAddress != t.fromAddress {
		return t.effectiveTransmitterAddress
	}
	return common.Address{}
}
