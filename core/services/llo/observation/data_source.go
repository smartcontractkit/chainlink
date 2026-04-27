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
	"github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
)

// Observation cache and loop tuning (all durations scale with observationTimeout T from the plugin Observe deadline).
//
// Invariants:
//   - cacheEntryTTL(T) = cacheTTLMultiplier·T — how long successful values remain valid after cache write.
//   - staleRefreshSkipThreshold(T) = (staleRefreshRemainingNumerator/Denominator)·T — a stream is not a refresh driver
//     while time.Until(expiresAt) is strictly greater than this (see buildStreamsRefreshPlan). A larger threshold
//     makes that inequality fail sooner as TTL decays, so the stream becomes a refresh driver earlier after each write
//     (higher freshness, more pipeline work); a smaller threshold lengthens the no-driver interval (staler reads, less load).
//   - Keep staleRefreshSkipThreshold(T)+observationLoopPacing(T) < cacheEntryTTL(T) (same T throughout). With
//     num/den = 8/5 and default pacing = T/10, (8/5+1/10)·T = 1.7·T < 2·T.
//
// Example timings for observationTimeout T = 250ms (cacheTTLMultiplier=2, pacing divisor=10, staleRefresh num/den = 8/5):
//   - cacheEntryTTL = 2·T = 500ms — TTL applied on successful per-pipeline-group AddMany writes.
//   - staleRefreshSkipThreshold = (8/5)·T = 400ms — a stream in the plugin scope is not a refresh driver while time.Until(expiresAt) > 400ms.
//   - observationLoopPacing = T/10 = 25ms (≥ observationLoopPacingMin and ≤ T/2) — minimum delay between loop iterations after the first (plugin Observe may wake the loop earlier; see loopWakeCh).
//   - per-iteration context uses WithTimeout(..., T) = 250ms — ceiling on wall time for one observation loop iteration (pipeline workers run in parallel under that deadline).
const (
	cacheTTLMultiplier                     = 2
	staleRefreshRemainingNumerator   int64 = 8
	staleRefreshRemainingDenominator int64 = 5

	observationLoopPacingMin     = 10 * time.Millisecond
	observationLoopPacingDivisor = 10 // pacing default = T/10, capped below
)

func cacheEntryTTL(observationTimeout time.Duration) time.Duration {
	return time.Duration(cacheTTLMultiplier) * observationTimeout
}

// staleRefreshSkipThreshold returns (staleRefreshRemainingNumerator/staleRefreshRemainingDenominator)·T.
// buildStreamsRefreshPlan treats a cached stream as still fresh (not a refresh driver) while time.Until(expiresAt)
// is strictly greater than this value. A larger fraction (e.g. higher numerator) raises the threshold, so the stream
// becomes a refresh driver again sooner after each successful write (more pipeline work, fresher cache entries).
func staleRefreshSkipThreshold(observationTimeout time.Duration) time.Duration {
	return (time.Duration(staleRefreshRemainingNumerator) * observationTimeout) / time.Duration(staleRefreshRemainingDenominator)
}

// observationLoopPacing returns the minimum time between observation loop iterations to cap CPU while
// staying responsive relative to T. Scales with T, clamped to [observationLoopPacingMin, T/2].
func observationLoopPacing(observationTimeout time.Duration) time.Duration {
	if observationTimeout <= 0 {
		return observationLoopPacingMin
	}
	p := observationTimeout / observationLoopPacingDivisor
	maxP := observationTimeout / 2
	if p < observationLoopPacingMin {
		p = observationLoopPacingMin
	}
	if p > maxP {
		p = maxP
	}
	return p
}

var (
	promMissingStreamCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "llo",
		Subsystem: "datasource",
		Name:      "stream_missing_count",
		Help:      "Number of times a stream had no pipeline in the registry (including refresh planning before Observe)",
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
		Help:      "Wall time for one observation loop iteration (pacing excluded)",
		Buckets: []float64{
			10, 25, 50, 100, 250, 500, 750, 1000,
		},
	},
		[]string{"configDigest"},
	)
	promObservationLoopWaitOutcome = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "llo",
		Subsystem: "datasource",
		Name:      "observation_loop_wait_outcome_count",
		Help:      "How the observation loop ended its inter-iteration wait: timer (pacing), wake (plugin Observe hint), or shutdown",
	},
		[]string{"outcome"},
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

	// loopWakeCh coalesces plugin Observe hints (buffer 1): the loop may skip pacing when the plugin is active.
	loopWakeCh chan struct{}
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
		loopWakeCh:             make(chan struct{}, 1),
	}
}

// signalObservationLoopWake notifies the background observation loop that the plugin called Observe (non-blocking,
// coalesced). Safe when the loop is not running (no-op).
func (d *dataSource) signalObservationLoopWake() {
	if d.loopWakeCh == nil || !d.observationLoopStarted.Load() {
		return
	}
	select {
	case d.loopWakeCh <- struct{}{}:
	default:
	}
}

// Observe starts or refreshes the background observation loop for the plugin's stream set, then fills streamValues
// from the in-memory cache (backed by pipeline observations registered for each stream ID).
func (d *dataSource) Observe(ctx context.Context, streamValues llo.StreamValues, opts llo.DSOpts) error {
	// Observation loop logic
	{
		// setObservableStreams copies stream IDs and deadline into internal state (the plugin's map is not retained).
		d.setObservableStreams(ctx, streamValues, opts)

		if !d.observationLoopStarted.Load() {
			loopStartedCh := make(chan struct{})
			go d.startObservationLoop(loopStartedCh)
			<-loopStartedCh
		}
		d.signalObservationLoopWake()
	}

	// Update stream values with the cached observations for all streams.
	d.cache.UpdateStreamValues(streamValues)

	return nil
}

// startObservationLoop continuously runs pipeline observations for this data source and caches results in memory,
// so the plugin Observe path stays fast regardless of adapter latency.
// Loop pacing (observationLoopPacing), per-round deadline (observationTimeout), cache TTL, and stale-refresh
// threshold (see package constants) are chosen so refresh tends to run before entries expire without busy-spinning.
// Each iteration schedules one goroutine per pipeline that needs work; that worker calls Observe once per
// stream ID in the intersection of p.StreamIDs() with the plugin's requested keys (see buildStreamsRefreshPlan), then
// AddMany for that batch if every Observe in the worker succeeded. If any Observe fails, the worker skips AddMany for that pipeline.
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
			osv := d.getObservableStreams()
			if osv == nil || len(osv.streamValues) == 0 {
				// There is nothing to observe, exit and let the next Observe() call reinitialize the loop.
				d.lggr.Warnw("observation loop: no streams to observe")
				return
			}

			if !loopStarting {
				// Pace the loop to bound CPU; plugin Observe can wake via loopWakeCh to skip sleep when data is needed sooner.
				pacing := observationLoopPacing(osv.observationTimeout)
				t := time.NewTimer(pacing)
				select {
				case <-stopChanCtx.Done():
					promObservationLoopWaitOutcome.WithLabelValues("shutdown").Inc()
					if !t.Stop() {
						select {
						case <-t.C:
						default:
						}
					}
					return
				case <-t.C:
					promObservationLoopWaitOutcome.WithLabelValues("timer").Inc()
				case <-d.loopWakeCh:
					promObservationLoopWaitOutcome.WithLabelValues("wake").Inc()
					if !t.Stop() {
						select {
						case <-t.C:
						default:
						}
					}
				}
			}

			if stopChanCtx.Err() != nil {
				return
			}

			startTS := time.Now()
			ctx, cancel := context.WithTimeout(stopChanCtx, osv.observationTimeout)
			lggr := logger.With(d.lggr, "observationTimestamp", osv.opts.ObservationTimestamp(), "configDigest", osv.opts.ConfigDigest(), "seqNr", osv.opts.OutCtx().SeqNr)

			var mu sync.Mutex
			var wg sync.WaitGroup
			var errs []ErrObservationFailed
			successfulStreamIDs := make([]streams.StreamID, 0, len(osv.streamValues))
			plan := d.buildStreamsRefreshPlan(osv.streamValues, osv.observationTimeout, lggr)
			ttl := cacheEntryTTL(osv.observationTimeout)

			if osv.opts.VerboseLogging() {
				lggr = logger.With(lggr, "staleStreamIDs", plan.streamIDsToRefresh)
				lggr.Debugw("Observing streams")
			}

			// Telemetry
			var telemCh chan<- interface{}
			{
				// Size needs to accommodate the max number of telemetry events that could be generated
				// Standard case might be about 3 bridge requests per spec and one stream<=>spec
				// Overallocate for safety (to avoid dropping packets). Sized from stale driver count; total Observe
				// calls can be higher when one pipeline has several plugin-requested streams (same worker, same iteration).
				telemCh = d.t.MakeObservationScopedTelemetryCh(osv.opts, 10*len(plan.streamIDsToRefresh))
				if telemCh != nil {
					if d.t.CaptureEATelemetry() {
						ctx = pipeline.WithTelemetryCh(ctx, telemCh)
					}
					if d.t.CaptureObservationTelemetry() {
						ctx = WithObservationTelemetryCh(ctx, telemCh)
					}
				}
			}

			oc := NewObservationContext(lggr, d.registry, d.t)
			for p := range plan.groups {
				wg.Add(1)
				go func(streamIDs []streams.StreamID) {
					defer wg.Done()
					local := make(llo.StreamValues, len(streamIDs))
					var hadErr bool
					for _, sid := range streamIDs {
						local[sid] = nil
						val, err := oc.Observe(ctx, sid, osv.opts)
						if err != nil {
							hadErr = true
							streamIDStr := strconv.FormatUint(uint64(sid), 10)
							if errors.As(err, &MissingStreamError{}) {
								promMissingStreamCount.WithLabelValues(streamIDStr).Inc()
							}
							promObservationErrorCount.WithLabelValues(streamIDStr).Inc()
							mu.Lock()
							errs = append(errs, ErrObservationFailed{inner: err, streamID: sid, reason: "failed to observe stream"})
							mu.Unlock()
							continue
						}
						mu.Lock()
						successfulStreamIDs = append(successfulStreamIDs, sid)
						mu.Unlock()
						local[sid] = val
					}
					if !hadErr {
						d.cache.AddMany(local, ttl)
					}
				}(plan.groups[p])
			}

			wg.Wait()
			elapsed = time.Since(startTS)

			// Unblock the first Observe caller once the initial observation round has finished.
			if loopStarting {
				loopStarting = false
				close(loopStartedCh)
			}

			// After all Observe calls return, nothing else is sent on the telemetry channel for this iteration.
			if telemCh != nil {
				close(telemCh)
			}

			slices.Sort(successfulStreamIDs)
			hasRefreshWork := len(plan.groups) > 0 || len(plan.missingStreamIDs) > 0
			elapsedMs := elapsed.Milliseconds()
			nOK := len(successfulStreamIDs)

			switch {
			case len(errs) > 0:
				sort.Slice(errs, func(i, j int) bool { return errs[i].streamID < errs[j].streamID })
				errStrs := make([]string, len(errs))
				failedStreamIDs := make([]streams.StreamID, len(errs))
				for i, e := range errs {
					errStrs[i] = e.String()
					failedStreamIDs[i] = e.streamID
				}
				wl := logger.With(lggr,
					"elapsed_ms", elapsedMs,
					"nSuccessfulStreams", nOK,
					"nFailedStreams", len(failedStreamIDs),
					"failedStreamIDs", failedStreamIDs,
					"errs", errStrs,
				)
				if osv.opts.VerboseLogging() {
					wl = logger.With(wl, "staleStreamIDs", plan.streamIDsToRefresh)
				}
				wl.Warnw("Observation loop completed with observation errors")
			case osv.opts.VerboseLogging():
				logger.With(lggr,
					"elapsed_ms", elapsedMs,
					"nSuccessfulStreams", nOK,
					"staleStreamIDs", plan.streamIDsToRefresh,
				).Debugw("Observation loop")
			case hasRefreshWork:
				lggr.Debugw("Observation loop",
					"elapsed_ms", elapsedMs,
					"nSuccessfulStreams", nOK,
				)
			}

			promObservationLoopDuration.WithLabelValues(
				osv.opts.ConfigDigest().String()).Observe(float64(elapsedMs))

			cancel()
		}
	})
}

// streamsRefreshPlan is the refresh scope for one observation loop iteration.
// groups maps each pipeline that needs work to the ordered list of stream IDs that worker will Observe: the
// intersection of p.StreamIDs() with the plugin's requested keys (streamValues), in pipeline order.
// streamIDsToRefresh is the sorted list of plugin-scope keys that are refresh drivers this round (stale or uncached).
// missingStreamIDs lists drivers with no registry entry (no Observe worker is started for them).
type streamsRefreshPlan struct {
	groups             map[streams.Pipeline][]streams.StreamID
	streamIDsToRefresh []streams.StreamID
	missingStreamIDs   []streams.StreamID
}

// buildStreamsRefreshPlan derives pipeline work groups and refresh-driver IDs from the cache and streamValues keys.
// A key is a refresh driver when the cache has no live value for it, or when time.Until(expiresAt) is not greater than
// staleRefreshSkipThreshold (same rule as package-level tuning comments). Each registered driver records its pipeline
// in groups; the worker observe list is built from p.StreamIDs() filtered to keys present in streamValues so we never
// run Observe for pipeline siblings the plugin did not request this round. Unregistered drivers go to missingStreamIDs;
// each increments promMissingStreamCount and triggers a single Warn when missingStreamIDs is non-empty.
func (d *dataSource) buildStreamsRefreshPlan(streamValues llo.StreamValues, observationTimeout time.Duration, lggr logger.Logger) streamsRefreshPlan {
	candidatesValues := make(llo.StreamValues, len(streamValues))
	for streamID := range streamValues {
		// Plugin-scope keys that need refresh become drivers; pipelines are collected below and scoped to these keys.
		if val, expiresAt := d.cache.Get(streamID); val != nil {
			if time.Until(expiresAt) > staleRefreshSkipThreshold(observationTimeout) {
				continue
			}
		}
		candidatesValues[streamID] = nil
	}

	// Observe all streams for the pipelines that have at least one candidate stream
	// that are in the plugin scope
	groups := make(map[streams.Pipeline][]streams.StreamID, len(candidatesValues))
	missingSet := []streams.StreamID{}
	for sid := range candidatesValues {
		p, ok := d.registry.Get(sid)
		if !ok {
			missingSet = append(missingSet, sid)
			continue
		}

		if _, ok := groups[p]; !ok {
			for _, sid := range p.StreamIDs() {
				if _, ok := streamValues[sid]; ok {
					groups[p] = append(groups[p], sid)
				}
			}
		}
	}

	var candidates = make([]streams.StreamID, 0, len(candidatesValues))
	for sid := range candidatesValues {
		candidates = append(candidates, sid)
	}
	slices.Sort(candidates)

	for _, sid := range missingSet {
		streamIDStr := strconv.FormatUint(uint64(sid), 10)
		promMissingStreamCount.WithLabelValues(streamIDStr).Inc()
	}
	if len(missingSet) > 0 {
		lggr.Warnw("observation loop: streams have no pipeline in registry; discarding",
			"missingStreamIDs", missingSet,
			"nMissingStreams", len(missingSet),
		)
	}

	return streamsRefreshPlan{
		groups:             groups,
		streamIDsToRefresh: candidates,
		missingStreamIDs:   missingSet,
	}
}

func (d *dataSource) Close() error {
	close(d.observationLoopCloseCh)
	d.wg.Wait()
	d.cache.Close()

	return nil
}

type observableStreamValues struct {
	opts               llo.DSOpts
	streamValues       llo.StreamValues
	observationTimeout time.Duration
}

// setObservableStreams updates the stream set and observation deadline (T) used by the background loop when in production.
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

// getObservableStreams returns the current observableStreamValues (opts, stream key set, and T), or nil if unset.
// T is the plugin Observe deadline; the loop uses a separate pacing timer between iterations and WithTimeout(..., T) per round.
func (d *dataSource) getObservableStreams() *observableStreamValues {
	d.observableStreamsMu.Lock()
	defer d.observableStreamsMu.Unlock()
	return d.observableStreams
}
