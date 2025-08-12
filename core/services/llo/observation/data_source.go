package observation

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

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
	lggr     logger.Logger
	registry Registry
	t        Telemeter
	cache    *Cache
	timeout  time.Duration

	observationLoopStarted bool
	observationLoopCloseCh chan struct{}
	configDigestToStreamMu sync.Mutex
	configDigestToStream   map[types.ConfigDigest]observableStreamValues
}

var _ io.Closer = dataSource{}

type observableStreamValues struct {
	opts         llo.DSOpts
	streamValues llo.StreamValues
}

func NewDataSource(lggr logger.Logger, registry Registry, t Telemeter) llo.DataSource {
	return newDataSource(lggr, registry, t, true)
}

func newDataSource(lggr logger.Logger, registry Registry, t Telemeter, shouldCache bool) *dataSource {
	return &dataSource{
		lggr:                 logger.Named(lggr, "DataSource"),
		registry:             registry,
		t:                    t,
		cache:                NewCache(500*time.Millisecond, time.Minute),
		configDigestToStream: make(map[types.ConfigDigest]observableStreamValues),
	}
}

// Observe looks up all streams in the registry and populates a map of stream ID => value
func (d *dataSource) Observe(ctx context.Context, streamValues llo.StreamValues, opts llo.DSOpts) error {
	// Observation loop logic
	{
		// Update the list of streams to observe for this config digest and set the timeout
		d.configDigestToStreamMu.Lock()
		d.configDigestToStream[opts.ConfigDigest()] = observableStreamValues{
			opts,
			streamValues,
		}

		deadline, ok := ctx.Deadline()
		if !ok {
			d.timeout = 100 * time.Millisecond
		} else {
			d.timeout = time.Until(deadline)
		}
		d.configDigestToStreamMu.Unlock()

		if !d.observationLoopStarted {
			loopStartedCh := make(chan struct{})
			go d.startObservationLoop(loopStartedCh)
			<-loopStartedCh
		}
	}

	// Fetch the cached observations for all streams.
	for streamID := range streamValues {
		val := d.fromCache(streamID)
		if val != nil {
			streamValues[streamID] = val
		}
	}

	return nil
}

// startObservationLoop continuously makes observations for the streams in d.configDigestToStream and stores those in
// the cache. It does not check for cached versions, it always calculates fresh values.
//
// NOTE: This method needs to be run in a goroutine.
func (d *dataSource) startObservationLoop(loopStartedCh chan struct{}) {
	for {
		now := time.Now()
		opts, streamValues := d.getObservableStreams()

		ctx, cancel := context.WithTimeout(context.Background(), d.timeout)
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
				d.toCache(streamID, val)
			}(streamID)
		}

		wg.Wait()

		// Notify the caller that we've completed our first round of observations.
		if !d.observationLoopStarted {
			d.observationLoopStarted = true
			close(loopStartedCh)
		}

		// After all Observations have returned, nothing else will be sent to the
		// telemetry channel, so it can safely be closed
		if telemCh != nil {
			close(telemCh)
		}

		// Only log on errors or if VerboseLogging is turned on
		if len(errs) > 0 || opts.VerboseLogging() {
			elapsed := time.Since(now)

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

		// Cancel the context, so the linter doesn't complain.
		cancel()

		// If we want to sleep between rounds, here is the place to do it.
		time.Sleep(d.timeout / 20) // turn this know if the EAs are under too much pressure
	}
}

func (d dataSource) Close() error {
	if !d.observationLoopStarted {
		return nil
	}
	if d.observationLoopCloseCh == nil {
		return nil
	}
	select {
	case <-d.observationLoopCloseCh:
		return nil
	default:
	}
	close(d.observationLoopCloseCh)
	return nil
}

func (d *dataSource) fromCache(streamID llotypes.StreamID) llo.StreamValue {
	if streamValue, found := d.cache.Get(streamID); found && streamValue != nil {
		return streamValue
	}
	return nil
}

func (d *dataSource) toCache(streamID llotypes.StreamID, val llo.StreamValue) {
	if val != nil {
		d.cache.Add(streamID, val)
	}
}

func (d *dataSource) getObservableStreams() (llo.DSOpts, llo.StreamValues) {
	d.configDigestToStreamMu.Lock()
	streamsToObserve := make([]observableStreamValues, 0, len(d.configDigestToStream))
	for _, vals := range d.configDigestToStream {
		streamsToObserve = append(streamsToObserve, vals)
	}
	d.configDigestToStreamMu.Unlock()

	var activeOpts llo.DSOpts
	activeStreamValues := make(llo.StreamValues)

	// deduplicate streams and get the active ocr instance options
	for _, vals := range streamsToObserve {
		for streamID := range vals.streamValues {
			activeStreamValues[streamID] = vals.streamValues[streamID]
		}

		if activeOpts == nil {
			outCtx := vals.opts.OutCtx()
			outcome, err := vals.opts.OutcomeCodec().Decode(outCtx.PreviousOutcome)
			if err != nil {
				d.lggr.Errorw("Failed to decode outcome", "error", err)
				continue
			}

			if outcome.LifeCycleStage == llo.LifeCycleStageProduction {
				activeOpts = vals.opts
			}
		}
	}

	return activeOpts, activeStreamValues
}
