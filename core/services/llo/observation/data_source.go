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
	wg                     sync.WaitGroup
	lggr                   logger.Logger
	registry               Registry
	t                      Telemeter
	cache                  StreamValueCache
	observationLoopStarted atomic.Bool
	observationLoopCloseCh services.StopChan

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

	// Update stream values with the cached observations for all streams.
	d.cache.UpdateStreamValues(streamValues)

	return nil
}

// startObservationLoop continuously makes observations for the streams in this data source
// caching them in memory making the Observe call duration and performance independent
// of the underlying resources providing the observations.
// Based on the expected maxObservationDuration determine the pace of the observation loop
// and for how long to cache the observations.
func (d *dataSource) startObservationLoop(loopStartedCh chan struct{}) {
	// atomically set the observation loop started flag to true
	// or return if it's already started
	if !d.observationLoopStarted.CompareAndSwap(false, true) {
		close(loopStartedCh)
		return
	}

	d.wg.Go(func() {
		loopStarting := true
		var elapsed time.Duration
		stopChanCtx, stopChanCancel := d.observationLoopCloseCh.NewCtx()
		defer stopChanCancel()

		defer func() {
			d.observationLoopStarted.Store(false)
			if loopStarting {
				close(loopStartedCh)
			}
		}()

		for {
			if stopChanCtx.Err() != nil {
				return
			}

			osv := d.getObservableStreams()
			if osv == nil || len(osv.streamValues) == 0 {
				// There is nothing to observe, exit and let the next Observe() call reinitialize the loop.
				d.lggr.Warnw("observation loop: no streams to observe")
				return
			}

			time.Sleep(osv.observationTimeout)
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
			observedValues := make(map[streams.StreamID]llo.StreamValue, len(osv.streamValues))
			var errs []ErrObservationFailed

			var wg sync.WaitGroup
			oc := NewObservationContext(lggr, d.registry, d.t)

			streamsToRefresh := d.getStreamsToRefresh(osv.streamValues, osv.observationTimeout)

			for streamID := range streamsToRefresh {
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

					mu.Lock()
					observedValues[streamID] = val
					successfulStreamIDs = append(successfulStreamIDs, streamID)
					mu.Unlock()
				}(streamID)
			}

			wg.Wait()
			elapsed = time.Since(startTS)

			droppedStreamIDs := d.removeIncompleteGroups(lggr, observedValues, osv.streamValues)

			d.cache.AddMany(observedValues, 4*osv.observationTimeout)

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
					len(observedValues), "nFailedStreams", len(failedStreamIDs), "nDroppedStreams", len(droppedStreamIDs), "errs", errStrs)

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
	})
}

// getStreamsToRefresh returns the set of stream IDs that need to be re-observed.
// When any stream in a pipeline is stale, ALL streams from that pipeline should be
// re-observed to ensure atomic observation of pipeline groups (e.g. bid/mid/ask must be observed together).
func (d *dataSource) getStreamsToRefresh(streamValues llo.StreamValues, observationTimeout time.Duration) map[streams.StreamID]struct{} {
	streamIDs := make(map[streams.StreamID]struct{})
	for streamID := range streamValues {
		if _, exists := streamIDs[streamID]; exists {
			continue
		}
		// refresh stream and associated streams from pipeline if this streamID is stale
		if val, expiresAt := d.cache.Get(streamID); val != nil {
			if time.Until(expiresAt) > 2*observationTimeout {
				continue
			}
		}

		streamIDs[streamID] = struct{}{}

		p, exists := d.registry.Get(streamID)
		if !exists {
			// pipeline isn't registered yet so we can't get associated stream IDs
			// this might happen if the plugin requests observations for streamIDs before
			// the node operator has registered its job spec or before the registry is fully initialized
			continue
		}

		for _, sid := range p.StreamIDs() {
			streamIDs[sid] = struct{}{}
		}
	}
	return streamIDs
}

func (d *dataSource) Close() error {
	close(d.observationLoopCloseCh)
	d.wg.Wait()
	d.cache.Close()

	return nil
}

// removeIncompleteGroups enforces all-or-nothing (atomic) writes per pipeline group.
// Some pipelines produce values that must be used together. For example jobs that output a bid/mid/ask
// must be used together to form a quote. So if any stream in the group failed, we drop
// the entire group to avoid writing a mix of fresh and stale values to the cache.
// Mutates observedValues in place.
func (d *dataSource) removeIncompleteGroups(lggr logger.Logger, observedValues map[streams.StreamID]llo.StreamValue, streamValues llo.StreamValues) []streams.StreamID {
	var dropped []streams.StreamID
	checked := make(map[streams.Pipeline]bool)
	for streamID := range observedValues {
		// we only need to check the pipeline once per group. So if we've already checked this pipeline, skip it.
		p, exists := d.registry.Get(streamID)
		if !exists || checked[p] {
			continue
		}
		checked[p] = true

		// Check that every in-scope stream for this pipeline succeeded.
		// This is because some pipelines might emit values for streams that the plugin is not requesting to be observed
		var missing []streams.StreamID
		for _, sid := range p.StreamIDs() {
			if _, inScope := streamValues[sid]; !inScope {
				continue // not requested this cycle so we can skip evaluating result
			}
			if _, ok := observedValues[sid]; !ok {
				missing = append(missing, sid)
			}
		}

		if len(missing) > 0 {
			var droppedFromGroup []streams.StreamID
			for _, sid := range p.StreamIDs() {
				if _, ok := observedValues[sid]; ok {
					droppedFromGroup = append(droppedFromGroup, sid)
				}
				delete(observedValues, sid)
			}
			dropped = append(dropped, droppedFromGroup...)
			lggr.Debugw("Discarding incomplete pipeline group",
				"pipelineStreamIDs", p.StreamIDs(),
				"missingStreamIDs", missing,
				"droppedStreamIDs", droppedFromGroup,
			)
		}
	}
	return dropped
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
