package transmitter

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/statuschecker"
)

type roundRobinKeystore interface {
	GetRoundRobinAddress(ctx context.Context, chainID *big.Int, addresses ...common.Address) (address common.Address, err error)
}

type txManager interface {
	CreateTransaction(ctx context.Context, txRequest txmgr.TxRequest) (tx txmgr.Tx, err error)
	GetTransactionStatus(ctx context.Context, transactionID string) (state commontypes.TransactionStatus, err error)
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
	chainID                     *big.Int
	keystore                    keys.RoundRobin
	statuschecker               statuschecker.CCIPTransactionStatusChecker // Used for CCIP's idempotency key generation
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
		chainID:                     chainID,
		keystore:                    keystore,
	}, nil
}

func NewTransmitterWithStatusChecker(
	txm txManager,
	fromAddresses []common.Address,
	gasLimit uint64,
	effectiveTransmitterAddress common.Address,
	strategy types.TxStrategy,
	checker txmgr.TransmitCheckerSpec,
	chainID *big.Int,
	keystore keys.RoundRobin,
) (Transmitter, error) {
	t, err := NewTransmitter(txm, fromAddresses, gasLimit, effectiveTransmitterAddress, strategy, checker, chainID, keystore)

	if err != nil {
		return nil, err
	}

	transmitter, ok := t.(*transmitter)
	if !ok {
		return nil, errors.New("failed to type assert Transmitter to *transmitter")
	}
	transmitter.statuschecker = statuschecker.NewTxmStatusChecker(txm.GetTransactionStatus)

	return transmitter, nil
}

func (t *transmitter) CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte, txMeta *txmgr.TxMeta) error {
	roundRobinFromAddress, err := t.keystore.GetNextAddress(ctx, t.fromAddresses...)
	if err != nil {
		return fmt.Errorf("skipped OCR transmission, error getting round-robin address: %w", err)
	}

	var idempotencyKey *string

	// Define idempotency key for CCIP Execution Plugin
	if len(txMeta.MessageIDs) == 1 && t.statuschecker != nil {
		messageId := txMeta.MessageIDs[0]
		_, count, err1 := t.statuschecker.CheckMessageStatus(ctx, messageId)

		if err1 != nil {
			return errors.Wrap(err, "skipped OCR transmission, error getting message status")
		}
		idempotencyKey = func() *string {
			s := fmt.Sprintf("%s-%d", messageId, count+1)
			return &s
		}()
	}

	_, err = t.txm.CreateTransaction(ctx, txmgr.TxRequest{
		IdempotencyKey:   idempotencyKey,
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

func (t *transmitter) FromAddress(ctx context.Context) common.Address {
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

func (t *transmitter) CreateSecondaryEthTransaction(ctx context.Context, bytes []byte, meta *txmgr.TxMeta) error {
	return errors.New("trying to send a secondary transmission on a non dual transmitter")
}

func (t *transmitter) SecondaryFromAddress(ctx context.Context) (common.Address, error) {
	return common.Address{}, errors.New("trying to get secondary address on a non dual transmitter")
}
