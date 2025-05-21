package evm

import (
	"context"
	"sync"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"
	gethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"

	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	offchain_aggregator_wrapper "github.com/smartcontractkit/chainlink/v2/core/internal/gethwrappers2/generated/offchainaggregator"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
)

// RequestRoundTracker subscribes to new request round logs.
type RequestRoundTracker struct {
	services.StateMachine

	ethClient        evmclient.Client
	contract         *offchain_aggregator_wrapper.OffchainAggregator
	contractFilterer *ocr2aggregator.OCR2AggregatorFilterer
	logBroadcaster   log.Broadcaster
	jobID            int32
	lggr             logger.SugaredLogger
	odb              RequestRoundDB
	ds               sqlutil.DataSource
	blockTranslator  ocrcommon.BlockTranslator

	// Start/Stop lifecycle
	stopCh          services.StopChan
	unsubscribeLogs func()

	// LatestRoundRequested
	latestRoundRequested ocr2aggregator.OCR2AggregatorRoundRequested
	lrrMu                sync.RWMutex
}

// NewRequestRoundTracker makes a new RequestRoundTracker
func NewRequestRoundTracker(
	contract *offchain_aggregator_wrapper.OffchainAggregator,
	contractFilterer *ocr2aggregator.OCR2AggregatorFilterer,
	ethClient evmclient.Client,
	logBroadcaster log.Broadcaster,
	jobID int32,
	lggr logger.Logger,
	ds sqlutil.DataSource,
	odb RequestRoundDB,
	chain ocrcommon.Config,
) (o *RequestRoundTracker) {
	return &RequestRoundTracker{
		ethClient:        ethClient,
		contract:         contract,
		contractFilterer: contractFilterer,
		logBroadcaster:   logBroadcaster,
		jobID:            jobID,
		lggr:             logger.Sugared(lggr),
		odb:              odb,
		ds:               ds,
		blockTranslator:  ocrcommon.NewBlockTranslator(chain, ethClient, lggr),
		stopCh:           make(chan struct{}),
	}
}

// Start must be called before logs can be delivered
// It ought to be called before starting OCR
func (t *RequestRoundTracker) Start(ctx context.Context) error {
	return t.StartOnce("RequestRoundTracker", func() (err error) {
		t.latestRoundRequested, err = t.odb.LoadLatestRoundRequested(ctx)
		if err != nil {
			return errors.Wrap(err, "RequestRoundTracker#Start: failed to load latest round requested")
		}

		t.unsubscribeLogs = t.logBroadcaster.Register(t, log.ListenerOpts{
			Contract: t.contract.Address(),
			ParseLog: t.contract.ParseLog,
			LogsWithTopics: map[gethCommon.Hash][][]log.Topic{
				offchain_aggregator_wrapper.OffchainAggregatorRoundRequested{}.Topic(): nil,
			},
			MinIncomingConfirmations: 1,
		})
		return nil
	})
}

// Close should be called after teardown of the OCR job relying on this tracker
func (t *RequestRoundTracker) Close() error {
	return t.StopOnce("RequestRoundTracker", func() error {
		close(t.stopCh)
		t.unsubscribeLogs()
		return nil
	})
}

// HandleLog complies with LogListener interface
// It is not thread safe
func (t *RequestRoundTracker) HandleLog(ctx context.Context, lb log.Broadcast) {
	was, err := t.logBroadcaster.WasAlreadyConsumed(ctx, lb)
	if err != nil {
		t.lggr.Errorw("OCRContract: could not determine if log was already consumed", "err", err)
		return
	} else if was {
		return
	}

	raw := lb.RawLog()
	if raw.Address != t.contract.Address() {
		t.lggr.Errorf("log address of 0x%x does not match configured contract address of 0x%x", raw.Address, t.contract.Address())
		t.lggr.ErrorIf(t.logBroadcaster.MarkConsumed(ctx, nil, lb), "unable to mark consumed")
		return
	}
	topics := raw.Topics
	if len(topics) == 0 {
		t.lggr.ErrorIf(t.logBroadcaster.MarkConsumed(ctx, nil, lb), "unable to mark consumed")
		return
	}

	var consumed bool
	switch topics[0] {
	case offchain_aggregator_wrapper.OffchainAggregatorRoundRequested{}.Topic():
		var rr *ocr2aggregator.OCR2AggregatorRoundRequested
		rr, err = t.contractFilterer.ParseRoundRequested(raw)
		if err != nil {
			t.lggr.Errorw("could not parse round requested", "err", err)
			t.lggr.ErrorIf(t.logBroadcaster.MarkConsumed(ctx, nil, lb), "unable to mark consumed")
			return
		}
		if IsLaterThan(raw, t.latestRoundRequested.Raw) {
			err = sqlutil.TransactDataSource(ctx, t.ds, nil, func(tx sqlutil.DataSource) error {
				if err = t.odb.WithDataSource(tx).SaveLatestRoundRequested(ctx, *rr); err != nil {
					return err
				}
				return t.logBroadcaster.MarkConsumed(ctx, tx, lb)
			})
			if err != nil {
				t.lggr.Error(err)
				return
			}
			consumed = true
			t.lrrMu.Lock()
			t.latestRoundRequested = *rr
			t.lrrMu.Unlock()
			t.lggr.Infow("RequestRoundTracker: received new latest RoundRequested event", "latestRoundRequested", *rr)
		} else {
			t.lggr.Warnw("RequestRoundTracker: ignoring out of date RoundRequested event", "latestRoundRequested", t.latestRoundRequested, "roundRequested", rr)
		}
	default:
		t.lggr.Debugw("RequestRoundTracker: got unrecognised log topic", "topic", topics[0])
	}
	if !consumed {
		t.lggr.ErrorIf(t.logBroadcaster.MarkConsumed(ctx, nil, lb), "unable to mark consumed")
	}
}

// IsLaterThan returns true if the first log was emitted "after" the second log
// from the blockchain's point of view
func IsLaterThan(incoming gethTypes.Log, existing gethTypes.Log) bool {
	return incoming.BlockNumber > existing.BlockNumber ||
		(incoming.BlockNumber == existing.BlockNumber && incoming.TxIndex > existing.TxIndex) ||
		(incoming.BlockNumber == existing.BlockNumber && incoming.TxIndex == existing.TxIndex && incoming.Index > existing.Index)
}

// IsV2Job complies with LogListener interface
func (t *RequestRoundTracker) IsV2Job() bool {
	return true
}

// JobID complies with LogListener interface
func (t *RequestRoundTracker) JobID() int32 {
	return t.jobID
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
func (t *RequestRoundTracker) LatestRoundRequested(_ context.Context, lookback time.Duration) (configDigest ocrtypes.ConfigDigest, epoch uint32, round uint8, err error) {
	t.lrrMu.RLock()
	defer t.lrrMu.RUnlock()
	return t.latestRoundRequested.ConfigDigest, t.latestRoundRequested.Epoch, t.latestRoundRequested.Round, nil
}
