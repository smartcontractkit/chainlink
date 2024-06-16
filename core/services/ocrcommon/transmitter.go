package ocrcommon

import (
	"context"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
)

type roundRobinKeystore interface {
	GetRoundRobinAddress(ctx context.Context, chainID *big.Int, addresses ...common.Address) (address common.Address, err error)
}

type txManager interface {
	CreateTransaction(ctx context.Context, txRequest txmgr.TxRequest) (tx txmgr.Tx, err error)
}

type Transmitter interface {
	CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte, txMeta *txmgr.TxMeta) error
	FromAddress() common.Address
}

type transmitter struct {
	txm                         txManager
	fromAddresses               []common.Address
	gasLimit                    uint64
	effectiveTransmitterAddress common.Address
	strategy                    types.TxStrategy
	checker                     txmgr.TransmitCheckerSpec
	chainID                     *big.Int
	keystore                    roundRobinKeystore
}

// NewTransmitter creates a new eth transmitter
func NewTransmitter(
	txm txManager,
	fromAddresses []common.Address,
	gasLimit uint64,
	effectiveTransmitterAddress common.Address,
	strategy types.TxStrategy,
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

type txManagerOCR2 interface {
	CreateTransaction(ctx context.Context, txRequest txmgr.TxRequest) (tx txmgr.Tx, err error)
	GetForwarderForEOAOCR2Feeds(ctx context.Context, eoa, ocr2AggregatorID common.Address) (forwarder common.Address, err error)
}

type ocr2FeedsTransmitter struct {
	ocr2Aggregator common.Address
	txManagerOCR2
	transmitter
}

// NewOCR2FeedsTransmitter creates a new eth transmitter that handles OCR2 Feeds specific logic surrounding forwarders.
// ocr2FeedsTransmitter validates forwarders before every transmission, enabling smooth onchain config changes without job restarts.
func NewOCR2FeedsTransmitter(
	txm txManagerOCR2,
	fromAddresses []common.Address,
	ocr2Aggregator common.Address,
	gasLimit uint64,
	effectiveTransmitterAddress common.Address,
	strategy types.TxStrategy,
	checker txmgr.TransmitCheckerSpec,
	chainID *big.Int,
	keystore roundRobinKeystore,
) (Transmitter, error) {
	// Ensure that a keystore is provided.
	if keystore == nil {
		return nil, errors.New("nil keystore provided to transmitter")
	}

	return &ocr2FeedsTransmitter{
		ocr2Aggregator: ocr2Aggregator,
		txManagerOCR2:  txm,
		transmitter: transmitter{
			txm:                         txm,
			fromAddresses:               fromAddresses,
			gasLimit:                    gasLimit,
			effectiveTransmitterAddress: effectiveTransmitterAddress,
			strategy:                    strategy,
			checker:                     checker,
			chainID:                     chainID,
			keystore:                    keystore,
		},
	}, nil
}

func (t *transmitter) CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte, txMeta *txmgr.TxMeta) error {
	roundRobinFromAddress, err := t.keystore.GetRoundRobinAddress(ctx, t.chainID, t.fromAddresses...)
	if err != nil {
		return errors.Wrap(err, "skipped OCR transmission, error getting round-robin address")
	}

	_, err = t.txm.CreateTransaction(ctx, txmgr.TxRequest{
		FromAddress:      roundRobinFromAddress,
		ToAddress:        toAddress,
		EncodedPayload:   payload,
		FeeLimit:         t.gasLimit,
		ForwarderAddress: t.forwarderAddress(),
		Strategy:         t.strategy,
		Checker:          t.checker,
		Meta:             txMeta,
	})
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

func (t *ocr2FeedsTransmitter) CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte, txMeta *txmgr.TxMeta) error {
	roundRobinFromAddress, err := t.keystore.GetRoundRobinAddress(ctx, t.chainID, t.fromAddresses...)
	if err != nil {
		return errors.Wrap(err, "skipped OCR transmission, error getting round-robin address")
	}

	forwarderAddress, err := t.forwarderAddress(ctx, roundRobinFromAddress, toAddress)
	if err != nil {
		return err
	}

	_, err = t.txm.CreateTransaction(ctx, txmgr.TxRequest{
		FromAddress:      roundRobinFromAddress,
		ToAddress:        toAddress,
		EncodedPayload:   payload,
		FeeLimit:         t.gasLimit,
		ForwarderAddress: forwarderAddress,
		Strategy:         t.strategy,
		Checker:          t.checker,
		Meta:             txMeta,
	})

	return errors.Wrap(err, "skipped OCR transmission")
}

// FromAddress for ocr2FeedsTransmitter returns valid forwarder or effectiveTransmitterAddress if forwarders are not set.
func (t *ocr2FeedsTransmitter) FromAddress() common.Address {
	roundRobinFromAddress, err := t.keystore.GetRoundRobinAddress(context.Background(), t.chainID, t.fromAddresses...)
	if err != nil {
		return t.effectiveTransmitterAddress
	}

	forwarderAddress, err := t.GetForwarderForEOAOCR2Feeds(context.Background(), roundRobinFromAddress, t.ocr2Aggregator)
	if errors.Is(err, forwarders.ErrForwarderForEOANotFound) {
		// if there are no valid forwarders try to fallback to eoa
		return roundRobinFromAddress
	} else if err != nil {
		return t.effectiveTransmitterAddress
	}

	return forwarderAddress
}

func (t *ocr2FeedsTransmitter) forwarderAddress(ctx context.Context, eoa, ocr2Aggregator common.Address) (common.Address, error) {
	//	If effectiveTransmitterAddress is in fromAddresses, then forwarders aren't set.
	if slices.Contains(t.fromAddresses, t.effectiveTransmitterAddress) {
		return common.Address{}, nil
	}

	forwarderAddress, err := t.GetForwarderForEOAOCR2Feeds(ctx, eoa, ocr2Aggregator)
	if err != nil {
		return common.Address{}, err
	}

	// if forwarder address is in fromAddresses, then none of the forwarders are valid
	if slices.Contains(t.fromAddresses, forwarderAddress) {
		forwarderAddress = common.Address{}
	}

	return forwarderAddress, nil
}
