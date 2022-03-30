package evm

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	offchain_aggregator_wrapper "github.com/smartcontractkit/chainlink/core/internal/gethwrappers2/generated/offchainaggregator"
	"github.com/smartcontractkit/chainlink/core/logger"
)

var _ median.MedianContract = &medianContract{}

type medianContract struct {
	contractCaller *ocr2aggregator.OCR2AggregatorCaller
	tracker        *RequestRoundTracker
}

func newMedianContract(contractAddress common.Address, chain evm.Chain, specID int32, db *sqlx.DB, lggr logger.Logger) (*medianContract, error) {
	contract, err := offchain_aggregator_wrapper.NewOffchainAggregator(contractAddress, chain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate NewOffchainAggregator")
	}

	contractFilterer, err := ocr2aggregator.NewOCR2AggregatorFilterer(contractAddress, chain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate NewOffchainAggregatorFilterer")
	}

	contractCaller, err := ocr2aggregator.NewOCR2AggregatorCaller(contractAddress, chain.Client())
	if err != nil {
		return nil, errors.Wrap(err, "could not instantiate NewOffchainAggregatorCaller")
	}

	return &medianContract{
		contractCaller: contractCaller,
		tracker: NewRequestRoundTracker(
			contract,
			contractFilterer,
			chain.Client(),
			chain.LogBroadcaster(),
			specID,
			lggr,
			db,
			NewRoundRequestedDB(db.DB, specID, lggr),
			chain.Config(),
		),
	}, nil
}

func (oc *medianContract) Start() error {
	return oc.tracker.Start()
}

func (oc *medianContract) Close() error {
	return oc.tracker.Close()
}

func (oc *medianContract) LatestTransmissionDetails(ctx context.Context) (ocrtypes.ConfigDigest, uint32, uint8, *big.Int, time.Time, error) {
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
func (oc *medianContract) LatestRoundRequested(ctx context.Context, lookback time.Duration) (ocrtypes.ConfigDigest, uint32, uint8, error) {
	return oc.tracker.LatestRoundRequested(ctx, lookback)
}
