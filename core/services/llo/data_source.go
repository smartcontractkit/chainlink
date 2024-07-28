package llo

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/exp/maps"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
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
}

func newDataSource(lggr logger.Logger, registry Registry) llo.DataSource {
	return &dataSource{lggr.Named("DataSource"), registry}
}

// Observe looks up all streams in the registry and populates a map of stream ID => value
func (d *dataSource) Observe(ctx context.Context, streamValues llo.StreamValues, opts llo.DSOpts) error {
	var wg sync.WaitGroup
	wg.Add(len(streamValues))
	var svmu sync.Mutex
	var errors []ErrObservationFailed
	var errmu sync.Mutex

	if opts.VerboseLogging() {
		streamIDs := make([]streams.StreamID, 0, len(streamValues))
		for streamID := range streamValues {
			streamIDs = append(streamIDs, streamID)
		}
		sort.Slice(streamIDs, func(i, j int) bool { return streamIDs[i] < streamIDs[j] })
		d.lggr.Debugw("Observing streams", "streamIDs", streamIDs, "seqNr", opts.SeqNr())
	}

	for _, streamID := range maps.Keys(streamValues) {
		go func(streamID llotypes.StreamID) {
			defer wg.Done()

			var val *big.Int

			stream, exists := d.registry.Get(streamID)
			if !exists {
				errmu.Lock()
				errors = append(errors, ErrObservationFailed{streamID: streamID, reason: fmt.Sprintf("missing stream: %d", streamID)})
				errmu.Unlock()
				promMissingStreamCount.WithLabelValues(fmt.Sprintf("%d", streamID)).Inc()
				return
			}
			run, trrs, err := stream.Run(ctx)
			if err != nil {
				errmu.Lock()
				errors = append(errors, ErrObservationFailed{inner: err, run: run, streamID: streamID, reason: "pipeline run failed"})
				errmu.Unlock()
				promObservationErrorCount.WithLabelValues(fmt.Sprintf("%d", streamID)).Inc()
				return
			}
			// TODO: support types other than *big.Int
			// https://smartcontract-it.atlassian.net/browse/MERC-3525
			val, err = streams.ExtractBigInt(trrs)
			if err != nil {
				errmu.Lock()
				errors = append(errors, ErrObservationFailed{inner: err, run: run, streamID: streamID, reason: "failed to extract big.Int"})
				errmu.Unlock()
				return
			}

			if val != nil {
				svmu.Lock()
				defer svmu.Unlock()
				streamValues[streamID] = nil
			}
		}(streamID)
	}

	wg.Wait()

	// Failed observations are always logged at warn level
	var failedStreamIDs []streams.StreamID
	if len(errors) > 0 {
		sort.Slice(errors, func(i, j int) bool { return errors[i].streamID < errors[j].streamID })
		failedStreamIDs = make([]streams.StreamID, len(errors))
		for i, e := range errors {
			failedStreamIDs[i] = e.streamID
		}
		d.lggr.Warnw("Observation failed for streams", "failedStreamIDs", failedStreamIDs, "errors", errors, "seqNr", opts.SeqNr())
	}

	if opts.VerboseLogging() {
		successes := make([]streams.StreamID, 0, len(streamValues))
		for strmID, res := range streamValues {
			if res != nil {
				successes = append(successes, strmID)
			}
		}
		sort.Slice(successes, func(i, j int) bool { return successes[i] < successes[j] })
		d.lggr.Debugw("Observation complete", "successfulStreamIDs", successes, "failedStreamIDs", failedStreamIDs, "values", streamValues, "seqNr", opts.SeqNr())
	}

	return nil
}
