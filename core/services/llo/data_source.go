package llo

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/maps"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	promMissingStreamCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "llo_stream_missing_count",
		Help: "Number of times we tried to observe a stream, but it was missing",
	},
		[]string{"streamID"},
	)
	promObservationErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "llo_stream_observation_error_count",
		Help: "Number of times we tried to observe a stream, but it failed with an error",
	},
		[]string{"streamID"},
	)
)

type Registry interface {
	Get(streamID streams.StreamID) (strm streams.Stream, exists bool)
}

type ErrObservationFailed struct {
	inner    error
	reason   string
	streamID streams.StreamID
	run      *pipeline.Run
}

func (e *ErrObservationFailed) Error() string {
	s := fmt.Sprintf("StreamID: %d; Reason: %s", e.streamID, e.reason)
	if e.inner != nil {
		s += fmt.Sprintf("; Err: %v", e.inner)
	}
	if e.run != nil {
		// NOTE: Could log more info about the run here if necessary
		s += fmt.Sprintf("; RunID: %d; RunErrors: %v", e.run.ID, e.run.AllErrors)
	}
	return s
}

func (e *ErrObservationFailed) String() string {
	return e.Error()
}

func (e *ErrObservationFailed) Unwrap() error {
	return e.inner
}

var _ llo.DataSource = &dataSource{}

type dataSource struct {
	lggr     logger.Logger
	registry Registry

	t Telemeter
}

func NewDataSource(lggr logger.Logger, registry Registry, t Telemeter) llo.DataSource {
	return newDataSource(lggr, registry, t)
}

func newDataSource(lggr logger.Logger, registry Registry, t Telemeter) *dataSource {
	return &dataSource{lggr.Named("DataSource"), registry, t}
}

// Observe looks up all streams in the registry and populates a map of stream ID => value
func (d *dataSource) Observe(ctx context.Context, streamValues llo.StreamValues, opts llo.DSOpts) error {
	var wg sync.WaitGroup
	wg.Add(len(streamValues))
	var svmu sync.Mutex
	var errs []ErrObservationFailed
	var errmu sync.Mutex

	if opts.VerboseLogging() {
		streamIDs := make([]streams.StreamID, 0, len(streamValues))
		for streamID := range streamValues {
			streamIDs = append(streamIDs, streamID)
		}
		sort.Slice(streamIDs, func(i, j int) bool { return streamIDs[i] < streamIDs[j] })
		d.lggr.Debugw("Observing streams", "streamIDs", streamIDs, "configDigest", opts.ConfigDigest(), "seqNr", opts.OutCtx().SeqNr)
	}

	for _, streamID := range maps.Keys(streamValues) {
		go func(streamID llotypes.StreamID) {
			defer wg.Done()

			var val llo.StreamValue

			stream, exists := d.registry.Get(streamID)
			if !exists {
				errmu.Lock()
				errs = append(errs, ErrObservationFailed{streamID: streamID, reason: fmt.Sprintf("missing stream: %d", streamID)})
				errmu.Unlock()
				promMissingStreamCount.WithLabelValues(fmt.Sprintf("%d", streamID)).Inc()
				return
			}
			run, trrs, err := stream.Run(ctx)
			if err != nil {
				errmu.Lock()
				errs = append(errs, ErrObservationFailed{inner: err, run: run, streamID: streamID, reason: "pipeline run failed"})
				errmu.Unlock()
				promObservationErrorCount.WithLabelValues(fmt.Sprintf("%d", streamID)).Inc()
				// TODO: Consolidate/reduce telemetry. We should send all observation results in a single packet
				// https://smartcontract-it.atlassian.net/browse/MERC-6290
				d.t.EnqueueV3PremiumLegacy(run, trrs, streamID, opts, nil, err)
				return
			}
			// TODO: Consolidate/reduce telemetry. We should send all observation results in a single packet
			// https://smartcontract-it.atlassian.net/browse/MERC-6290
			val, err = ExtractStreamValue(trrs)
			if err != nil {
				errmu.Lock()
				errs = append(errs, ErrObservationFailed{inner: err, run: run, streamID: streamID, reason: "failed to extract big.Int"})
				errmu.Unlock()
				return
			}

			d.t.EnqueueV3PremiumLegacy(run, trrs, streamID, opts, val, nil)

			if val != nil {
				svmu.Lock()
				defer svmu.Unlock()
				streamValues[streamID] = val
			}
		}(streamID)
	}

	wg.Wait()

	// Failed observations are always logged at warn level
	var failedStreamIDs []streams.StreamID
	if len(errs) > 0 {
		sort.Slice(errs, func(i, j int) bool { return errs[i].streamID < errs[j].streamID })
		failedStreamIDs = make([]streams.StreamID, len(errs))
		errStrs := make([]string, len(errs))
		for i, e := range errs {
			errStrs[i] = e.String()
			failedStreamIDs[i] = e.streamID
		}
		d.lggr.Warnw("Observation failed for streams", "failedStreamIDs", failedStreamIDs, "errs", errStrs, "configDigest", opts.ConfigDigest(), "seqNr", opts.OutCtx().SeqNr)
	}

	if opts.VerboseLogging() {
		successes := make([]streams.StreamID, 0, len(streamValues))
		for strmID := range streamValues {
			successes = append(successes, strmID)
		}
		sort.Slice(successes, func(i, j int) bool { return successes[i] < successes[j] })
		d.lggr.Debugw("Observation complete", "successfulStreamIDs", successes, "failedStreamIDs", failedStreamIDs, "configDigest", opts.ConfigDigest(), "values", streamValues, "seqNr", opts.OutCtx().SeqNr)
	}

	return nil
}

// ExtractStreamValue extracts a StreamValue from a TaskRunResults
func ExtractStreamValue(trrs pipeline.TaskRunResults) (llo.StreamValue, error) {
	// pipeline.TaskRunResults comes ordered asc by index, this is guaranteed
	// by the pipeline executor
	finaltrrs := trrs.Terminals()

	// TODO: Special handling for missing native/link streams?
	// https://smartcontract-it.atlassian.net/browse/MERC-5949

	// HACK: Right now we rely on the number of outputs to determine whether
	// its a Decimal or a Quote.
	// This isn't very robust or future-proof but is sufficient to support v0.3
	// compat.
	// There are a number of different possible ways to solve this in future.
	// See: https://smartcontract-it.atlassian.net/browse/MERC-5934
	switch len(finaltrrs) {
	case 1:
		res := finaltrrs[0].Result
		if res.Error != nil {
			return nil, res.Error
		}
		val, err := toDecimal(res.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse BenchmarkPrice: %w", err)
		}
		return llo.ToDecimal(val), nil
	case 3:
		// Expect ordering of Benchmark, Bid, Ask
		results := make([]decimal.Decimal, 3)
		for i, trr := range finaltrrs {
			res := trr.Result
			if res.Error != nil {
				return nil, fmt.Errorf("failed to parse stream output into Quote (task index: %d): %w", i, res.Error)
			}
			val, err := toDecimal(res.Value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse decimal: %w", err)
			}
			results[i] = val
		}
		return &llo.Quote{
			Benchmark: results[0],
			Bid:       results[1],
			Ask:       results[2],
		}, nil
	default:
		return nil, fmt.Errorf("invalid number of results, expected: 1 or 3, got: %d", len(finaltrrs))
	}
}

func toDecimal(val interface{}) (decimal.Decimal, error) {
	return utils.ToDecimal(val)
}
