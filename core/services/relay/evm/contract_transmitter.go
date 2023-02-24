package evm

import (
	"context"
	"database/sql"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type ContractTransmitter interface {
	services.ServiceCtx
	ocrtypes.ContractTransmitter
}

var _ ContractTransmitter = &contractTransmitter{}

type Transmitter interface {
	CreateEthTransaction(ctx context.Context, toAddress gethcommon.Address, payload []byte) error
	FromAddress() gethcommon.Address
}

type contractTransmitter struct {
	contractAddress     gethcommon.Address
	contractABI         abi.ABI
	transmitter         Transmitter
	transmittedEventSig common.Hash
	contractReader      contractReader
	lp                  logpoller.LogPoller
	lggr                logger.Logger
}

func NewOCRContractTransmitter(
	address gethcommon.Address,
	caller contractReader,
	contractABI abi.ABI,
	transmitter Transmitter,
	lp logpoller.LogPoller,
	lggr logger.Logger,
) (*contractTransmitter, error) {
	transmitted, ok := contractABI.Events["Transmitted"]
	if !ok {
		return nil, errors.New("invalid ABI, missing transmitted")
	}
	filterName := logpoller.FilterName("OCR ContractTransmitter", address.String())
	err := lp.RegisterFilter(logpoller.Filter{Name: filterName, EventSigs: []common.Hash{transmitted.ID}, Addresses: []common.Address{address}})
	if err != nil {
		return nil, err
	}
	return &contractTransmitter{
		contractAddress:     address,
		contractABI:         contractABI,
		transmitter:         transmitter,
		transmittedEventSig: transmitted.ID,
		lp:                  lp,
		contractReader:      caller,
		lggr:                lggr,
	}, nil
}

// Transmit sends the report to the on-chain smart contract's Transmit method.
func (oc *contractTransmitter) Transmit(ctx context.Context, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signatures []ocrtypes.AttributedOnchainSignature) error {
	var rs [][32]byte
	var ss [][32]byte
	var vs [32]byte
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

	oc.lggr.Debugw("Transmitting report", "report", hex.EncodeToString(report), "rawReportCtx", rawReportCtx, "contractAddress", oc.contractAddress)

	payload, err := oc.contractABI.Pack("transmit", rawReportCtx, []byte(report), rs, ss, vs)
	if err != nil {
		return errors.Wrap(err, "abi.Pack failed")
	}

	return errors.Wrap(oc.transmitter.CreateEthTransaction(ctx, oc.contractAddress, payload), "failed to send Eth transaction")
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
	latestConfigDigestAndEpoch, err := callContract(ctx, oc.contractAddress, oc.contractABI, "latestConfigDigestAndEpoch", nil, oc.contractReader)
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
	latest, err := oc.lp.LatestLogByEventSigWithConfs(oc.transmittedEventSig, oc.contractAddress, 1)
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
func (oc *contractTransmitter) FromAccount() ocrtypes.Account {
	return ocrtypes.Account(oc.transmitter.FromAddress().String())
}

func (oc *contractTransmitter) Start(ctx context.Context) error { return nil }
func (oc *contractTransmitter) Close() error                    { return nil }

// Has no state/lifecycle so it's always healthy and ready
func (oc *contractTransmitter) Healthy() error { return nil }
func (oc *contractTransmitter) Ready() error   { return nil }
func (oc *contractTransmitter) HealthReport() map[string]error {
	return map[string]error{oc.Name(): oc.Healthy()}
}
func (oc *contractTransmitter) Name() string { return "" }
