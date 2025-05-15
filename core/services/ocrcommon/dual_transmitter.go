package ocrcommon

import (
	"context"
	"net/url"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-evm/pkg/forwarders"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
)

type ocr2FeedsDualTransmission struct {
	txm                                txManager
	primaryFromAddresses               []common.Address
	gasLimit                           uint64
	primaryEffectiveTransmitterAddress common.Address
	strategy                           types.TxStrategy
	checker                            txmgr.TransmitCheckerSpec
	keystore                           keys.RoundRobin

	ocr2Aggregator common.Address
	txManagerOCR2

	secondaryContractAddress common.Address
	secondaryFromAddress     common.Address
	secondaryMeta            map[string][]string
}

func (t *ocr2FeedsDualTransmission) forwarderAddress(ctx context.Context, eoa, ocr2Aggregator common.Address) (common.Address, error) {
	//	If effectiveTransmitterAddress is in fromAddresses, then forwarders aren't set.
	if slices.Contains(t.primaryFromAddresses, t.primaryEffectiveTransmitterAddress) {
		return common.Address{}, nil
	}

	forwarderAddress, err := t.GetForwarderForEOAOCR2Feeds(ctx, eoa, ocr2Aggregator)
	if err != nil {
		return common.Address{}, err
	}

	// if forwarder address is in fromAddresses, then none of the forwarders are valid
	if slices.Contains(t.primaryFromAddresses, forwarderAddress) {
		forwarderAddress = common.Address{}
	}

	return forwarderAddress, nil
}

func (t *ocr2FeedsDualTransmission) CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte, txMeta *txmgr.TxMeta) error {
	roundRobinFromAddress, err := t.keystore.GetNextAddress(ctx, t.primaryFromAddresses...)
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

	return errors.Wrap(err, "skipped OCR transmission: skipped primary transmission")
}

func (t *ocr2FeedsDualTransmission) CreateSecondaryEthTransaction(ctx context.Context, payload []byte, txMeta *txmgr.TxMeta) error {
	forwarderAddress, err := t.forwarderAddress(ctx, t.secondaryFromAddress, t.secondaryContractAddress)
	if err != nil {
		return err
	}

	if txMeta == nil {
		txMeta = &txmgr.TxMeta{}
	}

	dualBroadcast := true
	dualBroadcastParams := t.urlParams()

	txMeta.DualBroadcast = &dualBroadcast
	txMeta.DualBroadcastParams = &dualBroadcastParams

	_, err = t.txm.CreateTransaction(ctx, txmgr.TxRequest{
		FromAddress:      t.secondaryFromAddress,
		ToAddress:        t.secondaryContractAddress,
		EncodedPayload:   payload,
		ForwarderAddress: forwarderAddress,
		FeeLimit:         t.gasLimit,
		Strategy:         t.strategy,
		Checker:          t.checker,
		Meta:             txMeta,
	})

	return errors.Wrap(err, "skipped secondary transmission")
}

func (t *ocr2FeedsDualTransmission) FromAddress(ctx context.Context) common.Address {
	roundRobinFromAddress, err := t.keystore.GetNextAddress(ctx, t.primaryFromAddresses...)
	if err != nil {
		return t.primaryEffectiveTransmitterAddress
	}

	forwarderAddress, err := t.GetForwarderForEOAOCR2Feeds(ctx, roundRobinFromAddress, t.ocr2Aggregator)
	if errors.Is(err, forwarders.ErrForwarderForEOANotFound) {
		// if there are no valid forwarders try to fallback to eoa
		return roundRobinFromAddress
	} else if err != nil {
		return t.primaryEffectiveTransmitterAddress
	}

	return forwarderAddress
}

func (t *ocr2FeedsDualTransmission) SecondaryFromAddress(ctx context.Context) (common.Address, error) {
	return t.secondaryFromAddress, nil
}

func (t *ocr2FeedsDualTransmission) urlParams() string {
	values := url.Values{}
	for k, v := range t.secondaryMeta {
		for _, p := range v {
			values.Add(k, p)
		}
	}
	return values.Encode()
}
