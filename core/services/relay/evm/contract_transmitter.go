package evm

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/logger"
)

var _ ocrtypes.ContractTransmitter = &ContractTransmitter{}

type Transmitter interface {
	CreateEthTransaction(ctx context.Context, toAddress gethcommon.Address, payload []byte) error
	FromAddress() gethcommon.Address
}

type ContractTransmitter struct {
	contractAddress gethcommon.Address
	contractABI     abi.ABI
	transmitter     Transmitter
	contractReader  contractReader
	lggr            logger.Logger
}

func NewOCRContractTransmitter(
	address gethcommon.Address,
	caller contractReader,
	contractABI abi.ABI,
	transmitter Transmitter,
	lggr logger.Logger,
) *ContractTransmitter {
	return &ContractTransmitter{
		contractAddress: address,
		contractABI:     contractABI,
		transmitter:     transmitter,
		contractReader:  caller,
		lggr:            lggr,
	}
}

// Transmit sends the report to the on-chain smart contract's Transmit method.
func (oc *ContractTransmitter) Transmit(ctx context.Context, reportCtx ocrtypes.ReportContext, report ocrtypes.Report, signatures []ocrtypes.AttributedOnchainSignature) error {
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
	FilterLogs(ctx context.Context, q ethereum.FilterQuery) ([]types.Log, error)
}

func parseTransmitted(log []byte) ([32]byte, uint32, error) {
	mustType := func(ts string) abi.Type {
		ty, _ := abi.NewType(ts, "", nil)
		return ty
	}
	var args abi.Arguments = []abi.Argument{
		{
			Name: "configDigest",
			Type: mustType("bytes32"),
		},
		{
			Name: "epoch",
			Type: mustType("uint32"),
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

// LatestConfigDigestAndEpoch retrieves the latest config digest and epoch from the OCR2 contract.
// It is plugin independent, in particular avoids use of the plugin specific generated evm wrappers
// by using the evm client Call directly for functions/events that are part of OCR2Abstract.
func (oc *ContractTransmitter) LatestConfigDigestAndEpoch(ctx context.Context) (ocrtypes.ConfigDigest, uint32, error) {
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

	// Otherwise, we have to scan for the logs. First get the latest config block as a log lower bound.
	latestConfigDetails, err := callContract(ctx, oc.contractAddress, oc.contractABI, "latestConfigDetails", nil, oc.contractReader)
	if err != nil {
		return ocrtypes.ConfigDigest{}, 0, err
	}
	configBlock := *abi.ConvertType(latestConfigDetails[1], new(uint32)).(*uint32)
	configDigest = *abi.ConvertType(latestConfigDetails[2], new([32]byte)).(*[32]byte)
	topics, err := abi.MakeTopics([]interface{}{oc.contractABI.Events["Transmitted"].ID})
	if err != nil {
		return ocrtypes.ConfigDigest{}, 0, err
	}
	query := ethereum.FilterQuery{
		Addresses: []gethcommon.Address{oc.contractAddress},
		Topics:    topics,
		FromBlock: new(big.Int).SetUint64(uint64(configBlock)),
	}
	logs, err := oc.contractReader.FilterLogs(ctx, query)
	if err != nil {
		return ocrtypes.ConfigDigest{}, 0, err
	}
	// No transmissions yet
	if len(logs) == 0 {
		return configDigest, 0, nil
	}
	// Logs come back ordered https://github.com/ethereum/go-ethereum/blob/d78590560d0107e727a44d0eea088eeb4d280bab/eth/filters/filter.go#L215
	// If there is a transmission, we take the latest one
	return parseTransmitted(logs[len(logs)-1].Data)
}

// FromAccount returns the account from which the transmitter invokes the contract
func (oc *ContractTransmitter) FromAccount() ocrtypes.Account {
	return ocrtypes.Account(oc.transmitter.FromAddress().String())
}
