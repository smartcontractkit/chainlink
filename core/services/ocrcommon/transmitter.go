package ocrcommon

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/pkg/forwarders"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type RoundRobinKeyLocker interface {
	keys.RoundRobin
	keys.Locker
}

type txManager interface {
	CreateTransaction(ctx context.Context, txRequest txmgr.TxRequest) (tx txmgr.Tx, err error)
}

type Transmitter interface {
	CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte, txMeta *txmgr.TxMeta) error
	FromAddress(context.Context) common.Address

	CreateSecondaryEthTransaction(context.Context, []byte, *txmgr.TxMeta) error
	SecondaryFromAddress(context.Context) (common.Address, error)
}

type transmitter struct {
	txm                         txManager
	fromAddresses               []common.Address
	gasLimit                    uint64
	effectiveTransmitterAddress common.Address
	strategy                    types.TxStrategy
	checker                     txmgr.TransmitCheckerSpec
	keystore                    keys.RoundRobin
}

// NewTransmitter creates a new eth transmitter
func NewTransmitter(
	txm txManager,
	fromAddresses []common.Address,
	gasLimit uint64,
	effectiveTransmitterAddress common.Address,
	strategy types.TxStrategy,
	checker txmgr.TransmitCheckerSpec,
	keystore keys.RoundRobin,
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
	ks RoundRobinKeyLocker,
	dualTransmissionConfig *evmtypes.DualTransmissionConfig,
) (Transmitter, error) {
	// Ensure that a keystore is provided.
	if ks == nil {
		return nil, errors.New("nil keystore provided to transmitter")
	}

	if keyHasLock(ks, effectiveTransmitterAddress, keys.TXMv2) {
		return nil, fmt.Errorf("key %s is used as a secondary transmitter in another job. primary and secondary transmitters cannot be mixed", effectiveTransmitterAddress.String())
	}

	if dualTransmissionConfig != nil {
		if keyHasLock(ks, dualTransmissionConfig.TransmitterAddress, keys.TXMv1) {
			return nil, fmt.Errorf("key %s is used as a primary transmitter in another job. primary and secondary transmitters cannot be mixed", effectiveTransmitterAddress.String())
		}
		return &ocr2FeedsDualTransmission{
			ocr2Aggregator:                     ocr2Aggregator,
			txm:                                txm,
			txManagerOCR2:                      txm,
			primaryFromAddresses:               fromAddresses,
			gasLimit:                           gasLimit,
			primaryEffectiveTransmitterAddress: effectiveTransmitterAddress,
			strategy:                           strategy,
			checker:                            checker,
			keystore:                           ks,
			secondaryContractAddress:           dualTransmissionConfig.ContractAddress,
			secondaryFromAddress:               dualTransmissionConfig.TransmitterAddress,
			secondaryMeta:                      dualTransmissionConfig.Meta,
		}, nil
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
			keystore:                    ks,
		},
	}, nil
}

func (t *transmitter) CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte, txMeta *txmgr.TxMeta) error {
	roundRobinFromAddress, err := t.keystore.GetNextAddress(ctx, t.fromAddresses...)
	if err != nil {
		return fmt.Errorf("skipped OCR transmission, error getting round-robin address: %w", err)
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
	if err != nil {
		return fmt.Errorf("skipped OCR transmission: %w", err)
	}
	return nil
}

func (t *transmitter) CreateSecondaryEthTransaction(ctx context.Context, bytes []byte, meta *txmgr.TxMeta) error {
	return errors.New("trying to send a secondary transmission on a non dual transmitter")
}
func (t *transmitter) SecondaryFromAddress(ctx context.Context) (common.Address, error) {
	return common.Address{}, errors.New("trying to get secondary address on a non dual transmitter")
}

func (t *transmitter) FromAddress(context.Context) common.Address {
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
	roundRobinFromAddress, err := t.keystore.GetNextAddress(ctx, t.fromAddresses...)
	if err != nil {
		return fmt.Errorf("skipped OCR transmission, error getting round-robin address: %w", err)
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

	if err != nil {
		return fmt.Errorf("skipped OCR transmission: %w", err)
	}
	return nil
}

// FromAddress for ocr2FeedsTransmitter returns valid forwarder or effectiveTransmitterAddress if forwarders are not set.
func (t *ocr2FeedsTransmitter) FromAddress(ctx context.Context) common.Address {
	roundRobinFromAddress, err := t.keystore.GetNextAddress(ctx, t.fromAddresses...)
	if err != nil {
		return t.effectiveTransmitterAddress
	}

	forwarderAddress, err := t.GetForwarderForEOAOCR2Feeds(ctx, roundRobinFromAddress, t.ocr2Aggregator)
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

func (t *ocr2FeedsTransmitter) CreateSecondaryEthTransaction(ctx context.Context, bytes []byte, meta *txmgr.TxMeta) error {
	return errors.New("trying to send a secondary transmission on a non dual transmitter")
}

func keyHasLock(ks keys.Locker, address common.Address, service keys.ServiceType) bool {
	return ks.GetMutex(address).IsLocked(service)
}
