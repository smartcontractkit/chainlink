package observation

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-common/pkg/services"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
)

var (
	promMissingStreamCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "llo",
		Subsystem: "datasource",
		Name:      "stream_missing_count",
		Help:      "Number of times we tried to observe a stream, but it was missing",
	},
		[]string{"streamID"},
	)
	promObservationErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "llo",
		Subsystem: "datasource",
		Name:      "stream_observation_error_count",
		Help:      "Number of times we tried to observe a stream, but it failed with an error",
	},
		[]string{"streamID"},
	)
	promObservationLoopDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "llo",
		Subsystem: "datasource",
		Name:      "observation_loop_duration_ms",
		Help:      "Duration of the observation loop",
		Buckets: []float64{
			10, 25, 50, 100, 250, 500, 750, 1000,
		},
	},
		[]string{"configDigest"},
	)
)

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
	lggr                   logger.Logger
	registry               Registry
	t                      Telemeter
	cache                  *Cache
	observationLoopStarted atomic.Bool
	observationLoopCloseCh services.StopChan
	observationLoopDoneCh  chan struct{} // will be closed when we exit the observation loop

	configDigestToStreamMu sync.Mutex
	configDigestToStream   map[types.ConfigDigest]observableStreamValues
}

func NewDataSource(lggr logger.Logger, registry Registry, t Telemeter) llo.DataSource {
	return newDataSource(lggr, registry, t, true)
}

func newDataSource(lggr logger.Logger, registry Registry, t Telemeter, shouldCache bool) *dataSource {
	return &dataSource{
		lggr:                   logger.Named(lggr, "DataSource"),
		registry:               registry,
		t:                      t,
		cache:                  NewCache(500*time.Millisecond, time.Minute),
		configDigestToStream:   make(map[types.ConfigDigest]observableStreamValues),
		observationLoopCloseCh: make(chan struct{}),
		observationLoopDoneCh:  make(chan struct{}),
	}
}

// Observe looks up all streams in the registry and populates a map of stream ID => value
func (d *dataSource) Observe(ctx context.Context, streamValues llo.StreamValues, opts llo.DSOpts) error {
	// Observation loop logic
	{
		// Update the list of streams to observe for this config digest and set the timeout
		// StreamValues  needs a copy to avoid concurrent access
		d.setObservableStreams(ctx, streamValues, opts)

		if !d.observationLoopStarted.Load() {
			loopStartedCh := make(chan struct{})
			go d.startObservationLoop(loopStartedCh)
			<-loopStartedCh
		}
	}

	// Fetch the cached observations for all streams.
	for streamID := range streamValues {
		streamValues[streamID] = d.cache.Get(streamID)
	}

	return nil
}

// startObservationLoop continuously makes observations for the streams in d.configDigestToStream and stores those in
// the cache. It does not check for cached versions, it always calculates fresh values.
//
// NOTE: This method needs to be run in a goroutine.
func (d *dataSource) startObservationLoop(loopStartedCh chan struct{}) {
	if !d.observationLoopStarted.CompareAndSwap(false, true) {
		close(loopStartedCh)
		return
	}

	loopStarting := true
	var elapsed time.Duration
	stopChanCtx, stopChanCancel := d.observationLoopCloseCh.NewCtx()
	defer stopChanCancel()

	for {
		if stopChanCtx.Err() != nil {
			close(d.observationLoopDoneCh)
			return
		}

		startTS := time.Now()
		opts, streamValues, observationInterval := d.getObservableStreams()
		if len(streamValues) == 0 || opts == nil {
			// There is nothing to observe, exit and let the next Observe() call reinitialize the loop.
			d.lggr.Debugw("invalid observation loop parameters", "opts", opts, "streamValues", streamValues)

			// still at the loop initialization, notify the caller and return
			if loopStarting {
				close(loopStartedCh)
			}
			return
		}

		ctx, cancel := context.WithTimeout(stopChanCtx, observationInterval)
		lggr := logger.With(d.lggr, "observationTimestamp", opts.ObservationTimestamp(), "configDigest", opts.ConfigDigest(), "seqNr", opts.OutCtx().SeqNr)

		if opts.VerboseLogging() {
			streamIDs := make([]streams.StreamID, 0, len(streamValues))
			for streamID := range streamValues {
				streamIDs = append(streamIDs, streamID)
			}
			sort.Slice(streamIDs, func(i, j int) bool { return streamIDs[i] < streamIDs[j] })
			lggr = logger.With(lggr, "streamIDs", streamIDs)
			lggr.Debugw("Observing streams")
		}

		// Telemetry
		var telemCh chan<- interface{}
		{
			// Size needs to accommodate the max number of telemetry events that could be generated
			// Standard case might be about 3 bridge requests per spec and one stream<=>spec
			// Overallocate for safety (to avoid dropping packets)
			telemCh = d.t.MakeObservationScopedTelemetryCh(opts, 10*len(streamValues))
			if telemCh != nil {
				if d.t.CaptureEATelemetry() {
					ctx = pipeline.WithTelemetryCh(ctx, telemCh)
				}
				if d.t.CaptureObservationTelemetry() {
					ctx = WithObservationTelemetryCh(ctx, telemCh)
				}
			}
		}

		var mu sync.Mutex
		successfulStreamIDs := make([]streams.StreamID, 0, len(streamValues))
		var errs []ErrObservationFailed

		var wg sync.WaitGroup
		wg.Add(len(streamValues))

		oc := NewObservationContext(lggr, d.registry, d.t)

		for streamID := range streamValues {
			go func(streamID llotypes.StreamID) {
				defer wg.Done()
				var val llo.StreamValue
				var err error

				// Observe the stream
				if val, err = oc.Observe(ctx, streamID, opts); err != nil {
					streamIDStr := strconv.FormatUint(uint64(streamID), 10)
					if errors.As(err, &MissingStreamError{}) {
						promMissingStreamCount.WithLabelValues(streamIDStr).Inc()
					}
					promObservationErrorCount.WithLabelValues(streamIDStr).Inc()
					mu.Lock()
					errs = append(errs, ErrObservationFailed{inner: err, streamID: streamID, reason: "failed to observe stream"})
					mu.Unlock()
					return
				}

				// cache the observed value
				d.cache.Add(streamID, val)
			}(streamID)
		}

		wg.Wait()
		elapsed = time.Since(startTS)

		// notify the caller that we've completed our first round of observations.
		if loopStarting {
			loopStarting = false
			close(loopStartedCh)
		}

		// After all Observations have returned, nothing else will be sent to the
		// telemetry channel, so it can safely be closed
		if telemCh != nil {
			close(telemCh)
		}

		// Only log on errors or if VerboseLogging is turned on
		if len(errs) > 0 || opts.VerboseLogging() {
			slices.Sort(successfulStreamIDs)
			sort.Slice(errs, func(i, j int) bool { return errs[i].streamID < errs[j].streamID })

			failedStreamIDs := make([]streams.StreamID, len(errs))
			errStrs := make([]string, len(errs))
			for i, e := range errs {
				errStrs[i] = e.String()
				failedStreamIDs[i] = e.streamID
			}

			lggr = logger.With(lggr, "elapsed", elapsed, "nSuccessfulStreams",
				len(successfulStreamIDs), "nFailedStreams", len(failedStreamIDs), "errs", errStrs)

			if opts.VerboseLogging() {
				lggr = logger.With(lggr, "streamValues", streamValues)
			}

			if len(errs) == 0 && opts.VerboseLogging() {
				lggr.Infow("Observation succeeded for all streamsToObserve")
			} else if len(errs) > 0 {
				lggr.Warnw("Observation failed for streamsToObserve")
			}
		}

		promObservationLoopDuration.WithLabelValues(
			opts.ConfigDigest().String()).Observe(float64(elapsed.Milliseconds()))

		if elapsed < observationInterval {
			lggr.Debugw("Observation loop sleep", "elapsed_ms", elapsed.Milliseconds(),
				"interval_ms", observationInterval.Milliseconds(), "sleep_ms", observationInterval-elapsed)
			time.Sleep(observationInterval - elapsed)
		} else {
			lggr.Debugw("Observation loop", "elapsed_ms", elapsed.Milliseconds(), "interval_ms", observationInterval.Milliseconds())
		}

		// Cancel the context, so the linter doesn't complain.
		cancel()
	}
}

func (d *dataSource) Close() error {
	close(d.observationLoopCloseCh)
	d.observationLoopStarted.Store(false)
	<-d.observationLoopDoneCh

	return nil
}

type observableStreamValues struct {
	opts                llo.DSOpts
	streamValues        llo.StreamValues
	observationInterval time.Duration
}

func (o *observableStreamValues) IsActive() (bool, error) {
	outCtx := o.opts.OutCtx()
	outcome, err := o.opts.OutcomeCodec().Decode(outCtx.PreviousOutcome)
	if err != nil {
		return false, fmt.Errorf("observable stream value: failed to decode outcome: %w", err)
	}

	if outcome.LifeCycleStage == llo.LifeCycleStageProduction {
		return true, nil
	}

	return false, nil
}

// setObservableStreams sets the observable streams for the given config digest.
func (d *dataSource) setObservableStreams(ctx context.Context, streamValues llo.StreamValues, opts llo.DSOpts) {
	values := make(llo.StreamValues, len(streamValues))
	for streamID := range streamValues {
		values[streamID] = nil
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(100 * time.Millisecond)
	}

	streamVals := make(llo.StreamValues)
	for streamID := range values {
		streamVals[streamID] = values[streamID]
	}

	d.configDigestToStreamMu.Lock()
	d.configDigestToStream[opts.ConfigDigest()] = observableStreamValues{
		opts:                opts,
		streamValues:        streamVals,
		observationInterval: time.Until(deadline),
	}
	d.configDigestToStreamMu.Unlock()
}

// getObservableStreams returns the active plugin data source options, the streams to observe and the observation interval
// the observation interval is the maximum time we can spend observing streams. We ensure that we don't exceed this time and
// we wait for the remaining time in the observation loop.
func (d *dataSource) getObservableStreams() (llo.DSOpts, llo.StreamValues, time.Duration) {
	d.configDigestToStreamMu.Lock()
	streamsToObserve := make([]observableStreamValues, 0, len(d.configDigestToStream))
	for _, vals := range d.configDigestToStream {
		streamsToObserve = append(streamsToObserve, vals)
	}
	d.configDigestToStreamMu.Unlock()

	for _, vals := range streamsToObserve {
		active, err := vals.IsActive()
		if err != nil {
			d.lggr.Errorw("getObservableStreams: failed to check if OCR instance is active", "error", err)
			continue
		}

		if active {
			return vals.opts, vals.streamValues, vals.observationInterval
		}

	}

	d.lggr.Errorw("getObservableStreams: no active OCR instance found")
	return nil, nil, 0
}
