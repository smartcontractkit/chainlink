package functions

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/encoding"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type roundRobinKeystore interface {
	GetRoundRobinAddress(ctx context.Context, chainID *big.Int, addresses ...common.Address) (address common.Address, err error)
}

type txManager interface {
	CreateTransaction(ctx context.Context, txRequest txmgr.TxRequest) (tx txmgr.Tx, err error)
}

type FunctionsContractTransmitter interface {
	services.ServiceCtx
	ocrtypes.ContractTransmitter
}

type ReportToEthMetadata func([]byte) (*txmgr.TxMeta, error)

type contractTransmitter struct {
	contractAddress             atomic.Pointer[common.Address]
	contractABI                 abi.ABI
	transmittedEventSig         common.Hash
	contractReader              contractReader
	lp                          logpoller.LogPoller
	lggr                        logger.Logger
	contractVersion             uint32
	reportCodec                 encoding.ReportCodec
	txm                         txManager
	fromAddresses               []common.Address
	gasLimit                    uint64
	effectiveTransmitterAddress common.Address
	strategy                    types.TxStrategy
	checker                     txmgr.TransmitCheckerSpec
	chainID                     *big.Int
	keystore                    roundRobinKeystore
}

var _ FunctionsContractTransmitter = &contractTransmitter{}
var _ evmRelayTypes.RouteUpdateSubscriber = &contractTransmitter{}

func transmitterFilterName(addr common.Address) string {
	return logpoller.FilterName("FunctionsOCR2ContractTransmitter", addr.String())
}

func NewFunctionsContractTransmitter(
	caller contractReader,
	contractABI abi.ABI,
	lp logpoller.LogPoller,
	lggr logger.Logger,
	contractVersion uint32,
	txm txManager,
	fromAddresses []common.Address,
	gasLimit uint64,
	effectiveTransmitterAddress common.Address,
	strategy types.TxStrategy,
	checker txmgr.TransmitCheckerSpec,
	chainID *big.Int,
	keystore roundRobinKeystore,
) (*contractTransmitter, error) {
	// Ensure that a keystore is provided.
	if keystore == nil {
		return nil, errors.New("nil keystore provided to transmitter")
	}

	transmitted, ok := contractABI.Events["Transmitted"]
	if !ok {
		return nil, errors.New("invalid ABI, missing transmitted")
	}

	if contractVersion != 1 {
		return nil, fmt.Errorf("unsupported contract version: %d", contractVersion)
	}

	codec, err := encoding.NewReportCodec(contractVersion)
	if err != nil {
		return nil, err
	}
	return &contractTransmitter{
		contractABI:                 contractABI,
		transmittedEventSig:         transmitted.ID,
		lp:                          lp,
		contractReader:              caller,
		lggr:                        logger.Named(lggr, "OCRFunctionsContractTransmitter"),
		contractVersion:             contractVersion,
		reportCodec:                 codec,
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

func (oc *contractTransmitter) createEthTransaction(ctx context.Context, toAddress common.Address, payload []byte) error {
	roundRobinFromAddress, err := oc.keystore.GetRoundRobinAddress(ctx, oc.chainID, oc.fromAddresses...)
	if err != nil {
		return errors.Wrap(err, "skipped OCR transmission, error getting round-robin address")
	}

	_, err = oc.txm.CreateTransaction(ctx, txmgr.TxRequest{
		FromAddress:      roundRobinFromAddress,
		ToAddress:        toAddress,
		EncodedPayload:   payload,
		FeeLimit:         oc.gasLimit,
		ForwarderAddress: oc.forwarderAddress(),
		Strategy:         oc.strategy,
		Checker:          oc.checker,
		Meta:             nil,
	})
	return errors.Wrap(err, "skipped OCR transmission")
}

func (oc *contractTransmitter) forwarderAddress() common.Address {
	for _, a := range oc.fromAddresses {
		if a == oc.effectiveTransmitterAddress {
			return common.Address{}
		}
	}
	return oc.effectiveTransmitterAddress
}

// Transmit sends the report to the on-chain smart contract's Transmit method.
func (oc *contractTransmitter) Transmit(ctx context.Context, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signatures []ocrtypes.AttributedOnchainSignature) error {
	var rs [][32]byte
	var ss [][32]byte
	var vs [32]byte
	if len(signatures) > 32 {
		return errors.New("too many signatures, maximum is 32")
	}
	for i, as := range signatures {
		r, s, v, err := evmutil.SplitSignature(as.Signature)
		if err != nil {
			panic("eventTransmit(ev): error in SplitSignature")
		}
		rs = append(rs, r)
		ss = append(ss, s)
		vs[i] = v
	}
	rawReportCtx := evmutil.RawReportContext(reportCtx)

	var destinationContract common.Address
	switch oc.contractVersion {
	case 1:
		oc.lggr.Debugw("FunctionsContractTransmitter: start", "reportLenBytes", len(report))
		requests, err2 := oc.reportCodec.DecodeReport(report)
		if err2 != nil {
			return errors.Wrap(err2, "FunctionsContractTransmitter: DecodeReport failed")
		}
		if len(requests) == 0 {
			return errors.New("FunctionsContractTransmitter: no requests in report")
		}
		if len(requests[0].CoordinatorContract) != common.AddressLength {
			return fmt.Errorf("FunctionsContractTransmitter: incorrect length of CoordinatorContract field: %d", len(requests[0].CoordinatorContract))
		}
		destinationContract.SetBytes(requests[0].CoordinatorContract)
		if destinationContract == (common.Address{}) {
			return errors.New("FunctionsContractTransmitter: destination coordinator contract is zero")
		}
		// Sanity check - every report should contain requests with the same coordinator contract.
		for _, req := range requests[1:] {
			if !bytes.Equal(req.CoordinatorContract, destinationContract.Bytes()) {
				oc.lggr.Errorw("FunctionsContractTransmitter: non-uniform coordinator addresses in a batch - still sending to a single destination",
					"requestID", hex.EncodeToString(req.RequestID),
					"destinationContract", destinationContract,
					"requestCoordinator", hex.EncodeToString(req.CoordinatorContract),
				)
			}
		}
		oc.lggr.Debugw("FunctionsContractTransmitter: ready", "nRequests", len(requests), "coordinatorContract", destinationContract.Hex())
	default:
		return fmt.Errorf("unsupported contract version: %d", oc.contractVersion)
	}
	payload, err := oc.contractABI.Pack("transmit", rawReportCtx, []byte(report), rs, ss, vs)
	if err != nil {
		return errors.Wrap(err, "abi.Pack failed")
	}

	oc.lggr.Debugw("FunctionsContractTransmitter: transmitting report", "contractAddress", destinationContract, "txMeta", nil, "payloadSize", len(payload))
	return errors.Wrap(oc.createEthTransaction(ctx, destinationContract, payload), "failed to send Eth transaction")
}

type contractReader interface {
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
}

func parseTransmitted(log []byte) ([32]byte, uint32, error) {
	var args abi.Arguments = []abi.Argument{
		{
			Name: "configDigest",
			Type: utils.MustAbiType("bytes32", nil),
		},
		{
			Name: "epoch",
			Type: utils.MustAbiType("uint32", nil),
		},
	}
	transmitted, err := args.Unpack(log)
	if err != nil {
		return [32]byte{}, 0, err
	}
	if len(transmitted) < 2 {
		return [32]byte{}, 0, errors.New("transmitted event log has too few arguments")
	}
	configDigest := *abi.ConvertType(transmitted[0], new([32]byte)).(*[32]byte)
	epoch := *abi.ConvertType(transmitted[1], new(uint32)).(*uint32)
	return configDigest, epoch, err
}

func callContract(ctx context.Context, addr common.Address, contractABI abi.ABI, method string, args []interface{}, caller contractReader) ([]interface{}, error) {
	input, err := contractABI.Pack(method, args...)
	if err != nil {
		return nil, err
	}
	output, err := caller.CallContract(ctx, ethereum.CallMsg{To: &addr, Data: input}, nil)
	if err != nil {
		return nil, err
	}
	return contractABI.Unpack(method, output)
}

// LatestConfigDigestAndEpoch retrieves the latest config digest and epoch from the OCR2 contract.
// It is plugin independent, in particular avoids use of the plugin specific generated evm wrappers
// by using the evm client Call directly for functions/events that are part of OCR2Abstract.
func (oc *contractTransmitter) LatestConfigDigestAndEpoch(ctx context.Context) (ocrtypes.ConfigDigest, uint32, error) {
	contractAddr := oc.contractAddress.Load()
	if contractAddr == nil {
		return ocrtypes.ConfigDigest{}, 0, errors.New("destination contract address not set")
	}
	latestConfigDigestAndEpoch, err := callContract(ctx, *contractAddr, oc.contractABI, "latestConfigDigestAndEpoch", nil, oc.contractReader)
	if err != nil {
		return ocrtypes.ConfigDigest{}, 0, err
	}
	// Panic on these conversions erroring, would mean a broken contract.
	scanLogs := *abi.ConvertType(latestConfigDigestAndEpoch[0], new(bool)).(*bool)
	configDigest := *abi.ConvertType(latestConfigDigestAndEpoch[1], new([32]byte)).(*[32]byte)
	epoch := *abi.ConvertType(latestConfigDigestAndEpoch[2], new(uint32)).(*uint32)
	if !scanLogs {
		return configDigest, epoch, nil
	}

	// Otherwise, we have to scan for the logs.
	if err != nil {
		return ocrtypes.ConfigDigest{}, 0, err
	}
	latest, err := oc.lp.LatestLogByEventSigWithConfs(ctx, oc.transmittedEventSig, *contractAddr, 1)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No transmissions yet
			return configDigest, 0, nil
		}
		return ocrtypes.ConfigDigest{}, 0, err
	}
	return parseTransmitted(latest.Data)
}

// FromAccount returns the account from which the transmitter invokes the contract
func (oc *contractTransmitter) FromAccount() (ocrtypes.Account, error) {
	return ocrtypes.Account(oc.effectiveTransmitterAddress.String()), nil
}

func (oc *contractTransmitter) Start(ctx context.Context) error { return nil }
func (oc *contractTransmitter) Close() error                    { return nil }

// Has no state/lifecycle so it's always healthy and ready
func (oc *contractTransmitter) Ready() error { return nil }
func (oc *contractTransmitter) HealthReport() map[string]error {
	return map[string]error{oc.Name(): nil}
}
func (oc *contractTransmitter) Name() string { return oc.lggr.Name() }

func (oc *contractTransmitter) UpdateRoutes(ctx context.Context, activeCoordinator common.Address, proposedCoordinator common.Address) error {
	// transmitter only cares about the active coordinator
	previousContract := oc.contractAddress.Swap(&activeCoordinator)
	if previousContract != nil && *previousContract == activeCoordinator {
		return nil
	}
	oc.lggr.Debugw("FunctionsContractTransmitter: updating routes", "previousContract", previousContract, "activeCoordinator", activeCoordinator)
	err := oc.lp.RegisterFilter(ctx, logpoller.Filter{Name: transmitterFilterName(activeCoordinator), EventSigs: []common.Hash{oc.transmittedEventSig}, Addresses: []common.Address{activeCoordinator}})
	if err != nil {
		return err
	}
	// TODO: unregister old filter (needs refactor to get pg.Queryer)
	return nil
}
