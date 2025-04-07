package transmitter

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-framework/chains/txmgr/types"
	"github.com/smartcontractkit/chainlink-integrations/evm/keys"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/statuschecker"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
)

type tronTransmitter struct {
	txm                         txManager
	contractABI                 abi.ABI
	fromAddresses               []common.Address
	gasLimit                    uint64
	effectiveTransmitterAddress common.Address
	strategy                    types.TxStrategy
	checker                     txmgr.TransmitCheckerSpec
	chainID                     *big.Int
	keystore                    keys.RoundRobin
	statuschecker               statuschecker.CCIPTransactionStatusChecker // Used for CCIP's idempotency key generation
}

// NewTronTransmitter creates a new tron transmitter
func NewTronTransmitter(
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

	OCR2AggregatorTransmissionContractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorMetaData.ABI))
	if err != nil {
		return nil, errors.New("unable to load OCR2AggregatorTransmissionContractABI needed for TronTransmitter")
	}

	return &tronTransmitter{
		txm:                         txm,
		contractABI:                 OCR2AggregatorTransmissionContractABI,
		fromAddresses:               fromAddresses,
		gasLimit:                    gasLimit,
		effectiveTransmitterAddress: effectiveTransmitterAddress,
		strategy:                    strategy,
		checker:                     checker,
		chainID:                     chainID,
		keystore:                    keystore,
	}, nil
}

func NewTronTransmitterWithStatusChecker(
	txm txManager,
	fromAddresses []common.Address,
	gasLimit uint64,
	effectiveTransmitterAddress common.Address,
	strategy types.TxStrategy,
	checker txmgr.TransmitCheckerSpec,
	chainID *big.Int,
	keystore keys.RoundRobin,
) (Transmitter, error) {
	t, err := NewTronTransmitter(txm, fromAddresses, gasLimit, effectiveTransmitterAddress, strategy, checker, chainID, keystore)

	if err != nil {
		return nil, err
	}

	transmitter, ok := t.(*tronTransmitter)
	if !ok {
		return nil, errors.New("failed to type assert Transmitter to *transmitter")
	}
	transmitter.statuschecker = statuschecker.NewTxmStatusChecker(txm.GetTransactionStatus)

	return transmitter, nil
}

// Trons Write API is different from the EVM Write API, so this function name is slightly misleading but it's to ensure we conform to the Transmitter interface
func (t *tronTransmitter) CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte, txMeta *txmgr.TxMeta) error {
	// TODO: RoundRobinFromAddress

	if t.IsExecTransmitter() {
		// TODO: Add idempotency key generation
	}

	_, report, rs, ss, vs, err := t.unpackPayload(payload)
	if err != nil {
		return fmt.Errorf("failed to unpack payload: %w", err)
	}

	fmt.Printf("Unpacked payload: report: %s, rs: %v, ss: %v, vs: %v",
		hex.EncodeToString(report),
		rs,
		ss,
		vs)

	fmt.Printf("About to transmit Tron Transaction... noop as Tron TXM is not yet implemented")

	return nil
}

func (t *tronTransmitter) FromAddress(ctx context.Context) common.Address {
	return t.effectiveTransmitterAddress
}

func (t *tronTransmitter) forwarderAddress() common.Address {
	for _, a := range t.fromAddresses {
		if a == t.effectiveTransmitterAddress {
			return common.Address{}
		}
	}
	return t.effectiveTransmitterAddress
}

func (t *tronTransmitter) CreateSecondaryEthTransaction(ctx context.Context, bytes []byte, meta *txmgr.TxMeta) error {
	return errors.New("trying to send a secondary transmission on a non dual transmitter")
}

func (t *tronTransmitter) SecondaryFromAddress(ctx context.Context) (common.Address, error) {
	return common.Address{}, errors.New("trying to get secondary address on a non dual transmitter")
}

func (t *tronTransmitter) IsCommitTransmitter() bool {
	return t.statuschecker == nil
}

func (t *tronTransmitter) IsExecTransmitter() bool {
	return !t.IsCommitTransmitter()
}

func (t *tronTransmitter) unpackPayload(payload []byte) (rawReportCtx [3][32]byte, report []byte, rs [][32]byte, ss [][32]byte, vs [32]byte, err error) {
	// Due to the quirks of the Tron API to send a transaction, we need to specify each parameter of the transmit function
	// To avoid modifying the OCRContractTransmitter, we can just unpack the payload here

	// Skip the first 4 bytes which is the function selector
	result, err := t.contractABI.Methods["transmit"].Inputs.Unpack(payload[4:])
	if err != nil {
		return [3][32]byte{}, nil, nil, nil, [32]byte{}, fmt.Errorf("failed to unpack payload: %w", err)
	}

	rawReportCtx = result[0].([3][32]byte)
	report = result[1].([]byte)
	rs = result[2].([][32]byte)
	ss = result[3].([][32]byte)
	vs = result[4].([32]byte)

	return rawReportCtx, report, rs, ss, vs, nil
}
