package offchainreporting2

import (
	"context"
	"encoding/hex"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/chains/evmutil"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
)

var (
	_ ocrtypes.ContractTransmitter = &ContractTransmitter{}
	_ median.MedianContract        = &ContractTransmitter{}
)

type Transmitter interface {
	CreateEthTransaction(ctx context.Context, toAddress gethCommon.Address, payload []byte) error
	FromAddress() gethCommon.Address
}

type ContractTransmitter struct {
	contractAddress gethCommon.Address
	contractABI     abi.ABI
	transmitter     Transmitter
	contractCaller  *ocr2aggregator.OCR2AggregatorCaller
	tracker         *ContractTracker
	lggr            logger.Logger
}

func NewOCRContractTransmitter(
	address gethCommon.Address,
	contractCaller *ocr2aggregator.OCR2AggregatorCaller,
	contractABI abi.ABI,
	transmitter Transmitter,
	tracker *ContractTracker,
	lggr logger.Logger,
) *ContractTransmitter {
	return &ContractTransmitter{
		contractAddress: address,
		contractABI:     contractABI,
		transmitter:     transmitter,
		contractCaller:  contractCaller,
		tracker:         tracker,
		lggr:            lggr,
	}
}

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

func (oc *ContractTransmitter) LatestConfigDigestAndEpoch(ctx context.Context) (ocrtypes.ConfigDigest, uint32, error) {
	opts := bind.CallOpts{Context: ctx, Pending: false}
	result, err := oc.contractCaller.LatestTransmissionDetails(&opts)
	return result.ConfigDigest, result.Epoch, errors.Wrap(err, "error getting LatestTransmissionDetails")
}

func (oc *ContractTransmitter) FromAccount() ocrtypes.Account {
	return ocrtypes.Account(oc.transmitter.FromAddress().String())
}

func (oc *ContractTransmitter) LatestTransmissionDetails(ctx context.Context) (ocrtypes.ConfigDigest, uint32, uint8, *big.Int, time.Time, error) {
	opts := bind.CallOpts{Context: ctx, Pending: false}
	result, err := oc.contractCaller.LatestTransmissionDetails(&opts)
	return result.ConfigDigest, result.Epoch, result.Round, result.LatestAnswer, time.Unix(int64(result.LatestTimestamp), 0), errors.Wrap(err, "error getting LatestTransmissionDetails")
}

// LatestRoundRequested returns the configDigest, epoch, and round from the latest
// RoundRequested event emitted by the contract. LatestRoundRequested may or may not
// return a result if the latest such event was emitted in a block b such that
// b.timestamp < tip.timestamp - lookback.
//
// If no event is found, LatestRoundRequested should return zero values, not an error.
// An error should only be returned if an actual error occurred during execution,
// e.g. because there was an error querying the blockchain or the database.
//
// As an optimization, this function may also return zero values, if no
// RoundRequested event has been emitted after the latest NewTransmission event.
func (oc *ContractTransmitter) LatestRoundRequested(ctx context.Context, lookback time.Duration) (ocrtypes.ConfigDigest, uint32, uint8, error) {
	return oc.tracker.LatestRoundRequested(ctx, lookback)
}
