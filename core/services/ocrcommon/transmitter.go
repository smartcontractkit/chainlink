package ocrcommon

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
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
	gasLimit                    uint32
	effectiveTransmitterAddress common.Address
	strategy                    txmgr.TxStrategy
	checker                     txmgr.TransmitCheckerSpec
	chainID                     *big.Int
	ethKeyStore                 keystore.Eth
}

// NewTransmitter creates a new eth transmitter
func NewTransmitter(
	txm txManager,
	fromAddresses []common.Address,
	gasLimit uint32,
	effectiveTransmitterAddress common.Address,
	strategy txmgr.TxStrategy,
	checker txmgr.TransmitCheckerSpec,
	chainID *big.Int,
	ethKeyStore keystore.Eth,
) Transmitter {
	return &transmitter{
		txm:                         txm,
		fromAddresses:               fromAddresses,
		gasLimit:                    gasLimit,
		effectiveTransmitterAddress: effectiveTransmitterAddress,
		strategy:                    strategy,
		checker:                     checker,
		chainID:                     chainID,
		ethKeyStore:                 ethKeyStore,
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
	// Default to first sending key.
	nextFromAddress := t.fromAddresses[0]

	// Apply round-robin logic for a valid keystore.
	if t.ethKeyStore != nil {
		fromAddress, err := t.ethKeyStore.GetRoundRobinAddress(t.chainID, t.fromAddresses...)
		if err == nil {
			nextFromAddress = fromAddress
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
