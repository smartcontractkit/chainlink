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

	observableStreamsMu sync.Mutex
	observableStreams   *observableStreamValues
}

func NewDataSource(lggr logger.Logger, registry Registry, t Telemeter) llo.DataSource {
	return newDataSource(lggr, registry, t)
}

func newDataSource(lggr logger.Logger, registry Registry, t Telemeter) *dataSource {
	return &dataSource{
		lggr:                   logger.Named(lggr, "DataSource"),
		registry:               registry,
		t:                      t,
		cache:                  NewCache(time.Minute),
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
		streamValues[streamID], _ = d.cache.Get(streamID)
	}

	return nil
}

// startObservationLoop continuously makes observations for the streams in this data source
// caching them in memory making the Observe call duration and performance independent
// of the underlying resources providing the observations.
// Based on the expected maxObservationDuration determine the pace of the observation loop
// and for how long to cache the observations.
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

		osv := d.getObservableStreams()
		if osv == nil || len(osv.streamValues) == 0 {
			// There is nothing to observe, exit and let the next Observe() call reinitialize the loop.
			d.lggr.Warnw("observation loop: no streams to observe")

			// still at the loop initialization, notify the caller and return
			if loopStarting {
				close(loopStartedCh)
			}
			return
		}

		if d.observationLoopStarted.Load() {
			time.Sleep(osv.observationTimeout)
		}

		startTS := time.Now()
		ctx, cancel := context.WithTimeout(stopChanCtx, osv.observationTimeout)
		lggr := logger.With(d.lggr, "observationTimestamp", osv.opts.ObservationTimestamp(), "configDigest", osv.opts.ConfigDigest(), "seqNr", osv.opts.OutCtx().SeqNr)

		if osv.opts.VerboseLogging() {
			streamIDs := make([]streams.StreamID, 0, len(osv.streamValues))
			for streamID := range osv.streamValues {
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
			telemCh = d.t.MakeObservationScopedTelemetryCh(osv.opts, 10*len(osv.streamValues))
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
		successfulStreamIDs := make([]streams.StreamID, 0, len(osv.streamValues))
		var errs []ErrObservationFailed

		var wg sync.WaitGroup
		oc := NewObservationContext(lggr, d.registry, d.t)

		for streamID := range osv.streamValues {
			if val, expiresAt := d.cache.Get(streamID); val != nil {
				if time.Until(expiresAt) > 2*osv.observationTimeout {
					d.lggr.Debugw("cached stream observation still valid, skipping", "streamID",
						streamID, "expiresAt", expiresAt.Format(time.RFC3339))
					continue
				}
			}

			wg.Add(1)
			go func(streamID llotypes.StreamID) {
				defer wg.Done()
				var val llo.StreamValue
				var err error

				// Observe the stream
				if val, err = oc.Observe(ctx, streamID, osv.opts); err != nil {
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
				d.cache.Add(streamID, val, 4*osv.observationTimeout)

				mu.Lock()
				successfulStreamIDs = append(successfulStreamIDs, streamID)
				mu.Unlock()
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
		if len(errs) > 0 || osv.opts.VerboseLogging() {
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

			if osv.opts.VerboseLogging() {
				lggr = logger.With(lggr, "streamValues", osv.streamValues)
			}
		}

		promObservationLoopDuration.WithLabelValues(
			osv.opts.ConfigDigest().String()).Observe(float64(elapsed.Milliseconds()))

		lggr.Debugw("Observation loop", "elapsed_ms", elapsed.Milliseconds())

		// context cancellation
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
	opts               llo.DSOpts
	streamValues       llo.StreamValues
	observationTimeout time.Duration
}

// setObservableStreams sets the observable streams for the given config digest.
func (d *dataSource) setObservableStreams(ctx context.Context, streamValues llo.StreamValues, opts llo.DSOpts) {
	if opts == nil || len(streamValues) == 0 {
		d.lggr.Warnw("setObservableStreams: no observable streams to set",
			"opts", opts, "observable_streams", len(streamValues))
		return
	}

	outCtx := opts.OutCtx()
	outcome, err := opts.OutcomeCodec().Decode(outCtx.PreviousOutcome)
	if err != nil {
		d.lggr.Errorw("setObservableStreams: failed to decode outcome", "error", err)
		return
	}

	if outcome.LifeCycleStage != llo.LifeCycleStageProduction {
		d.lggr.Debugw(
			"setObservableStreams: LLO OCR instance is not in production lifecycle stage",
			"configDigest", opts.ConfigDigest().String(), "stage", outcome.LifeCycleStage)
		return
	}

	osv := &observableStreamValues{
		opts:               opts,
		streamValues:       make(llo.StreamValues, len(streamValues)),
		observationTimeout: 250 * time.Millisecond,
	}

	for streamID := range streamValues {
		osv.streamValues[streamID] = nil
	}

	if deadline, ok := ctx.Deadline(); ok {
		osv.observationTimeout = time.Until(deadline)
	}

	d.lggr.Debugw("setObservableStreams",
		"timeout_millis", osv.observationTimeout.Milliseconds(),
		"observable_streams", len(osv.streamValues))

	d.observableStreamsMu.Lock()
	defer d.observableStreamsMu.Unlock()

	if d.observableStreams == nil ||
		len(d.observableStreams.streamValues) != len(osv.streamValues) ||
		d.observableStreams.observationTimeout != osv.observationTimeout {
		d.lggr.Infow("setObservableStreams: observable streams changed",
			"timeout_millis", osv.observationTimeout.Milliseconds(),
			"observable_streams", len(osv.streamValues),
		)
	}

	d.observableStreams = osv
}

// getObservableStreams returns the active plugin data source options, the streams to observe and the observation interval
// the observation interval is the maximum time we can spend observing streams. We ensure that we don't exceed this time and
// we wait for the remaining time in the observation loop.
func (d *dataSource) getObservableStreams() *observableStreamValues {
	d.observableStreamsMu.Lock()
	defer d.observableStreamsMu.Unlock()
	return d.observableStreams
}
