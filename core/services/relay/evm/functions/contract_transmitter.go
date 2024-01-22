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

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/functions/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	evmRelayTypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

type FunctionsContractTransmitter interface {
	services.ServiceCtx
	ocrtypes.ContractTransmitter
}

type Transmitter interface {
	CreateEthTransaction(ctx context.Context, toAddress common.Address, payload []byte, txMeta *txmgr.TxMeta) error
	FromAddress() common.Address
}

type ReportToEthMetadata func([]byte) (*txmgr.TxMeta, error)

func reportToEvmTxMetaNoop([]byte) (*txmgr.TxMeta, error) {
	return nil, nil
}

type contractTransmitter struct {
	contractAddress     atomic.Pointer[common.Address]
	contractABI         abi.ABI
	transmitter         Transmitter
	transmittedEventSig common.Hash
	contractReader      contractReader
	lp                  logpoller.LogPoller
	lggr                logger.Logger
	reportToEvmTxMeta   ReportToEthMetadata
	contractVersion     uint32
	reportCodec         encoding.ReportCodec
}

var _ FunctionsContractTransmitter = &contractTransmitter{}
var _ evmRelayTypes.RouteUpdateSubscriber = &contractTransmitter{}

func transmitterFilterName(addr common.Address) string {
	return logpoller.FilterName("FunctionsOCR2ContractTransmitter", addr.String())
}

func NewFunctionsContractTransmitter(
	caller contractReader,
	contractABI abi.ABI,
	transmitter Transmitter,
	lp logpoller.LogPoller,
	lggr logger.Logger,
	reportToEvmTxMeta ReportToEthMetadata,
	contractVersion uint32,
) (*contractTransmitter, error) {
	transmitted, ok := contractABI.Events["Transmitted"]
	if !ok {
		return nil, errors.New("invalid ABI, missing transmitted")
	}

	if contractVersion != 1 {
		return nil, fmt.Errorf("unsupported contract version: %d", contractVersion)
	}

	if reportToEvmTxMeta == nil {
		reportToEvmTxMeta = reportToEvmTxMetaNoop
	}
	codec, err := encoding.NewReportCodec(contractVersion)
	if err != nil {
		return nil, err
	}
	return &contractTransmitter{
		contractABI:         contractABI,
		transmitter:         transmitter,
		transmittedEventSig: transmitted.ID,
		lp:                  lp,
		contractReader:      caller,
		lggr:                lggr.Named("OCRContractTransmitter"),
		reportToEvmTxMeta:   reportToEvmTxMeta,
		contractVersion:     contractVersion,
		reportCodec:         codec,
	}, nil
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

	txMeta, err := oc.reportToEvmTxMeta(report)
	if err != nil {
		oc.lggr.Warnw("failed to generate tx metadata for report", "err", err)
	}

	var destinationContract common.Address
	if oc.contractVersion == 1 {
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
	} else {
		return fmt.Errorf("unsupported contract version: %d", oc.contractVersion)
	}
	payload, err := oc.contractABI.Pack("transmit", rawReportCtx, []byte(report), rs, ss, vs)
	if err != nil {
		return errors.Wrap(err, "abi.Pack failed")
	}

	oc.lggr.Debugw("FunctionsContractTransmitter: transmitting report", "contractAddress", destinationContract, "txMeta", txMeta, "payloadSize", len(payload))
	return errors.Wrap(oc.transmitter.CreateEthTransaction(ctx, destinationContract, payload, txMeta), "failed to send Eth transaction")
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
	latest, err := oc.lp.LatestLogByEventSigWithConfs(
		oc.transmittedEventSig, *contractAddr, 1, pg.WithParentCtx(ctx))
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
	return ocrtypes.Account(oc.transmitter.FromAddress().String()), nil
}

func (oc *contractTransmitter) Start(ctx context.Context) error { return nil }
func (oc *contractTransmitter) Close() error                    { return nil }

// Has no state/lifecycle so it's always healthy and ready
func (oc *contractTransmitter) Ready() error { return nil }
func (oc *contractTransmitter) HealthReport() map[string]error {
	return map[string]error{oc.Name(): nil}
}
func (oc *contractTransmitter) Name() string { return oc.lggr.Name() }

func (oc *contractTransmitter) UpdateRoutes(activeCoordinator common.Address, proposedCoordinator common.Address) error {
	// transmitter only cares about the active coordinator
	previousContract := oc.contractAddress.Swap(&activeCoordinator)
	if previousContract != nil && *previousContract == activeCoordinator {
		return nil
	}
	oc.lggr.Debugw("FunctionsContractTransmitter: updating routes", "previousContract", previousContract, "activeCoordinator", activeCoordinator)
	err := oc.lp.RegisterFilter(logpoller.Filter{Name: transmitterFilterName(activeCoordinator), EventSigs: []common.Hash{oc.transmittedEventSig}, Addresses: []common.Address{activeCoordinator}})
	if err != nil {
		return err
	}
	// TODO: unregister old filter (needs refactor to get pg.Queryer)
	return nil
}
