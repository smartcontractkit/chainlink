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

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	llodatasource "github.com/smartcontractkit/chainlink-data-streams/llo/datasource"
	lloprotocol "github.com/smartcontractkit/chainlink-data-streams/llo/protocol"

	"github.com/smartcontractkit/chainlink/v2/core/services/llo/telem"
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
//     num/den = 13/5 (= 13/5·T stale) and divisor 2 (raw T/2 pacing), the invariant cap is (cacheEntryTTL−stale−1ns) =
//     7/5·T − 1ns, which exceeds raw T/2, so pacing is bounded by T/2 (not the invariant) and the strict inequality
//     holds with slack (13/5·T + T/2 = 31/10·T < 4·T).
//
// Example timings for observationTimeout T = 250ms (cacheTTLMultiplier=4, pacing divisor=2, staleRefresh num/den = 13/5):
//   - cacheEntryTTL = 4·T = 1000ms — TTL applied on successful per-pipeline-group AddMany writes.
//   - staleRefreshSkipThreshold = (13/5)·T = 650ms — a stream in the plugin scope is not a refresh driver while time.Until(expiresAt) > 650ms.
//     Steady-state refresh interval per stream = cacheEntryTTL − staleRefreshSkipThreshold = (4 − 13/5)·T = 7/5·T = 350ms.
//   - observationLoopPacing targets T/2 = 125ms and is bounded by min(T/2, (4−13/5)·T − 1ns) = min(125ms, 350ms − 1ns) = 125ms
//     (≥ observationLoopPacingFloor) — minimum delay between loop iterations after the first (plugin Observe may wake the loop earlier; see loopWakeCh).
//   - per-iteration context uses WithTimeout(..., T) = 250ms — ceiling on wall time for one observation loop iteration (pipeline workers run in parallel under that deadline).
const (
	cacheTTLMultiplier                     = 4
	staleRefreshRemainingNumerator   int64 = 13
	staleRefreshRemainingDenominator int64 = 5

	observationLoopPacingFloor   = 10 * time.Millisecond
	observationLoopPacingDivisor = 2 // pacing targets T/2, capped below by cache invariant
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
// staying responsive relative to T. Scales with T/divisor, clamped to [observationLoopPacingFloor, min(T/2,
// cacheEntryTTL(T)−staleRefreshSkipThreshold(T)−1ns)] so staleRefreshSkipThreshold+observationLoopPacing < cacheEntryTTL.
func observationLoopPacing(observationTimeout time.Duration) time.Duration {
	if observationTimeout <= 0 {
		return observationLoopPacingFloor
	}
	p := observationTimeout / observationLoopPacingDivisor
	invMax := cacheEntryTTL(observationTimeout) - staleRefreshSkipThreshold(observationTimeout) - time.Nanosecond
	maxP := min(invMax, observationTimeout/2)
	if p < observationLoopPacingFloor {
		p = observationLoopPacingFloor
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

type ObservationFailedError struct { //nolint:revive // name matches existing observation failure terminology
	inner    error
	reason   string
	streamID streams.StreamID
	run      *pipeline.Run
}

func (e *ObservationFailedError) Error() string {
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

func (e *ObservationFailedError) String() string {
	return e.Error()
}

func (e *ObservationFailedError) Unwrap() error {
	return e.inner
}

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

var _ llodatasource.DataSource = &dataSource{}

// NewDataSource returns the shared LLO data source. llo/v30 and llo/v31 both
// consume llodatasource.DataSource, so a single implementation serves both OCR
// protocol versions; lifecycle gating is driven by opts.LifeCycleStage().
func NewDataSource(lggr logger.Logger, registry Registry, t Telemeter) llodatasource.DataSource {
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

// Observe gates the background observation loop on the round's lifecycle stage
// (only a Production instance runs pipeline observations), then fills
// streamValues from the in-memory cache. The stage is carried by opts and is
// derived by each plugin version from its own state (v30 from the previous
// outcome, v31 from the KeyValueState), so a single implementation serves both.
func (d *dataSource) Observe(ctx context.Context, streamValues lloprotocol.StreamValues, opts llodatasource.DSOpts) error {
	return d.observe(ctx, streamValues, opts, d.inProduction(opts))
}

// inProduction reports whether this OCR instance is the Production instance (the
// only one that should run pipeline observations).
func (d *dataSource) inProduction(opts llodatasource.DSOpts) bool {
	if opts == nil {
		// setObservableStreams logs the nil-opts case; stay silent here to avoid
		// a duplicate warning per round.
		return false
	}
	if opts.LifeCycleStage() != lloprotocol.LifeCycleStageProduction {
		d.lggr.Debugw("Observe: LLO OCR instance is not in production lifecycle stage",
			"configDigest", opts.ConfigDigest().String(), "stage", opts.LifeCycleStage())
		return false
	}
	return true
}

// observe starts or refreshes the background observation loop for the plugin's stream set, then fills streamValues
// from the in-memory cache (backed by pipeline observations registered for each stream ID). inProduction gates the
// loop: when false the observable stream set is cleared and no pipelines run this round.
func (d *dataSource) observe(ctx context.Context, streamValues lloprotocol.StreamValues, opts telem.DSOpts, inProduction bool) error {
	// Observation loop logic
	{
		// setObservableStreams copies stream IDs and deadline into internal state (the plugin's map is not retained).
		d.setObservableStreams(ctx, streamValues, opts, inProduction)

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
			lggr := logger.With(d.lggr, "observationTimestamp", osv.opts.ObservationTimestamp(), "configDigest", osv.opts.ConfigDigest(), "seqNr", osv.opts.SeqNr())

			var mu sync.Mutex
			var wg sync.WaitGroup
			var errs []ObservationFailedError
			successfulStreamIDs := make([]streams.StreamID, 0, len(osv.streamValues))
			plan := d.buildStreamsRefreshPlan(osv.streamValues, osv.observationTimeout, lggr)
			ttl := cacheEntryTTL(osv.observationTimeout)

			if osv.opts.VerboseLogging() {
				lggr = logger.With(lggr, "staleStreamIDs", plan.streamIDsToRefresh)
				lggr.Debugw("Observing streams")
			}

			// Telemetry
			var telemCh chan<- any
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
					local := make(lloprotocol.StreamValues, len(streamIDs))
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
							errs = append(errs, ObservationFailedError{inner: err, streamID: sid, reason: "failed to observe stream"})
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
func (d *dataSource) buildStreamsRefreshPlan(streamValues lloprotocol.StreamValues, observationTimeout time.Duration, lggr logger.Logger) streamsRefreshPlan {
	candidatesValues := make(lloprotocol.StreamValues, len(streamValues))
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
	opts               telem.DSOpts
	streamValues       lloprotocol.StreamValues
	observationTimeout time.Duration
}

// setObservableStreams updates the stream set and observation deadline (T) used by the background loop. When
// inProduction is false (v30 non-production instance) the observable set is left unchanged/empty so no pipelines run.
func (d *dataSource) setObservableStreams(ctx context.Context, streamValues lloprotocol.StreamValues, opts telem.DSOpts, inProduction bool) {
	if opts == nil || len(streamValues) == 0 {
		d.lggr.Warnw("setObservableStreams: no observable streams to set",
			"opts", opts, "observable_streams", len(streamValues))
		return
	}

	if !inProduction {
		d.lggr.Debugw(
			"setObservableStreams: LLO OCR instance is not in production lifecycle stage",
			"configDigest", opts.ConfigDigest().String())
		return
	}

	osv := &observableStreamValues{
		opts:               opts,
		streamValues:       make(lloprotocol.StreamValues, len(streamValues)),
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
