package ocrcommon

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type roundRobinKeystore interface {
	GetRoundRobinAddress(chainID *big.Int, addresses ...common.Address) (address common.Address, err error)
}

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
	keystore                    roundRobinKeystore
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
	keystore roundRobinKeystore,
) (Transmitter, error) {

	// Ensure that a keystore is provided.
	if keystore == nil {
		return nil, errors.New("nil keystore provided to transmitter")
	}

	return &transmitter{
		txm:                         txm,
		fromAddresses:               fromAddresses,
		gasLimit:                    gasLimit,
		effectiveTransmitterAddress: effectiveTransmitterAddress,
		strategy:                    strategy,
		checker:                     checker,
		chainID:                     chainID,
		keystore:                    keystore,
	}, nil
}

func (t *transmitter) CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte) error {

	roundRobinFromAddress, err := t.keystore.GetRoundRobinAddress(t.chainID, t.fromAddresses...)
	if err != nil {
		return errors.Wrap(err, "skipped OCR transmission, error getting round-robin address")
	}

	_, err = t.txm.CreateEthTransaction(txmgr.NewTx{
		FromAddress:      roundRobinFromAddress,
		ToAddress:        toAddress,
		EncodedPayload:   payload,
		GasLimit:         t.gasLimit,
		ForwarderAddress: t.forwarderAddress(),
		Strategy:         t.strategy,
		Checker:          t.checker,
	}, pg.WithParentCtx(ctx))
	return errors.Wrap(err, "skipped OCR transmission")
}

func (t *transmitter) FromAddress() common.Address {
	return t.effectiveTransmitterAddress
}

func (t *transmitter) forwarderAddress() common.Address {
	for _, a := range t.fromAddresses {
		if a == t.effectiveTransmitterAddress {
			return common.Address{}
		}
	}
	return t.effectiveTransmitterAddress
}
