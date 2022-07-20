package ocr

import (
	"context"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/ocrcommon"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"
)

var (
	_ ocrtypes.ContractTransmitter = &OCRContractTransmitter{}
)

type (
	OCRContractTransmitter struct {
		contractAddress gethCommon.Address
		contractABI     abi.ABI
		transmitter     ocrcommon.Transmitter
		contractCaller  *offchainaggregator.OffchainAggregatorCaller
		tracker         *OCRContractTracker
		chainID         *big.Int
	}
)

func NewOCRContractTransmitter(
	address gethCommon.Address,
	contractCaller *offchainaggregator.OffchainAggregatorCaller,
	contractABI abi.ABI,
	transmitter ocrcommon.Transmitter,
	logBroadcaster log.Broadcaster,
	tracker *OCRContractTracker,
	chainID *big.Int,
) *OCRContractTransmitter {
	return &OCRContractTransmitter{
		contractAddress: address,
		contractABI:     contractABI,
		transmitter:     transmitter,
		contractCaller:  contractCaller,
		tracker:         tracker,
		chainID:         chainID,
	}
}

func (oc *OCRContractTransmitter) Transmit(ctx context.Context, report []byte, rs, ss [][32]byte, vs [32]byte) error {
	payload, err := oc.contractABI.Pack("transmit", report, rs, ss, vs)
	if err != nil {
		return errors.Wrap(err, "abi.Pack failed")
	}

	return errors.Wrap(oc.transmitter.CreateEthTransaction(ctx, oc.contractAddress, payload), "failed to send Eth transaction")
}

func (oc *OCRContractTransmitter) LatestTransmissionDetails(ctx context.Context) (configDigest ocrtypes.ConfigDigest, epoch uint32, round uint8, latestAnswer ocrtypes.Observation, latestTimestamp time.Time, err error) {
	opts := bind.CallOpts{Context: ctx, Pending: false}
	result, err := oc.contractCaller.LatestTransmissionDetails(&opts)
	if err != nil {
		return configDigest, 0, 0, ocrtypes.Observation(nil), time.Time{}, errors.Wrap(err, "error getting LatestTransmissionDetails")
	}
	return result.ConfigDigest, result.Epoch, result.Round, ocrtypes.Observation(result.LatestAnswer), time.Unix(int64(result.LatestTimestamp), 0), nil
}

func (oc *OCRContractTransmitter) FromAddress() gethCommon.Address {
	return oc.transmitter.FromAddress()
}

func (oc *OCRContractTransmitter) ChainID() *big.Int {
	return oc.chainID
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
func (oc *OCRContractTransmitter) LatestRoundRequested(ctx context.Context, lookback time.Duration) (configDigest ocrtypes.ConfigDigest, epoch uint32, round uint8, err error) {
	return oc.tracker.LatestRoundRequested(ctx, lookback)
}
