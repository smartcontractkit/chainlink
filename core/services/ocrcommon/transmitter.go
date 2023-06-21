package ocrcommon

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

type roundRobinKeystore interface {
	GetRoundRobinAddress(chainID *big.Int, addresses ...common.Address) (address common.Address, err error)
}

type txManager interface {
	CreateTransaction(txRequest txmgr.EvmTxRequest, qopts ...pg.QOpt) (tx txmgr.EvmTx, err error)
}

type Transmitter interface {
	CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte, txMeta *txmgr.EvmTxMeta) error
	FromAddress() common.Address
}

type transmitter struct {
	txm                         txManager
	fromAddresses               []common.Address
	gasLimit                    uint32
	effectiveTransmitterAddress common.Address
	strategy                    types.TxStrategy
	checker                     txmgr.EvmTransmitCheckerSpec
	chainID                     *big.Int
	keystore                    roundRobinKeystore
}

// NewTransmitter creates a new eth transmitter
func NewTransmitter(
	txm txManager,
	fromAddresses []common.Address,
	gasLimit uint32,
	effectiveTransmitterAddress common.Address,
	strategy types.TxStrategy,
	checker txmgr.EvmTransmitCheckerSpec,
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

func (t *transmitter) CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte, txMeta *txmgr.EvmTxMeta) error {

	roundRobinFromAddress, err := t.keystore.GetRoundRobinAddress(t.chainID, t.fromAddresses...)
	if err != nil {
		return errors.Wrap(err, "skipped OCR transmission, error getting round-robin address")
	}

	_, err = t.txm.CreateTransaction(txmgr.EvmTxRequest{
		FromAddress:      roundRobinFromAddress,
		ToAddress:        toAddress,
		EncodedPayload:   payload,
		FeeLimit:         t.gasLimit,
		ForwarderAddress: t.forwarderAddress(),
		Strategy:         t.strategy,
		Checker:          t.checker,
		Meta:             txMeta,
	}, pg.WithParentCtx(ctx))
	return errors.Wrap(err, "skipped OCR transmission")
}

func (t *transmitter) FromAddress() common.Address {
	log.Println("$$$$$ Effective transmitter address here is ", t.effectiveTransmitterAddress)
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
