package observation

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"maps"
	"math"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	promtest "github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	lloprotocol "github.com/smartcontractkit/chainlink-data-streams/llo/protocol"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	clhttptest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/httptest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/llo/telem"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
)

type mockPipeline struct {
	run  *pipeline.Run
	trrs pipeline.TaskRunResults
	err  error

	streamIDs []streams.StreamID

	runCount atomic.Int32
}

func (m *mockPipeline) Run(ctx context.Context) (*pipeline.Run, pipeline.TaskRunResults, error) {
	m.runCount.Add(1)
	return m.run, m.trrs, m.err
}

func (m *mockPipeline) StreamIDs() []streams.StreamID {
	return m.streamIDs
}

type mockRegistry struct {
	mu        sync.Mutex
	pipelines map[streams.StreamID]*mockPipeline
}

func (m *mockRegistry) Get(streamID streams.StreamID) (p streams.Pipeline, exists bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	p, exists = m.pipelines[streamID]
	return p, exists
}

func makePipelineWithSingleResult[T any](runID int64, res T, err error) *mockPipeline {
	return &mockPipeline{
		run:  &pipeline.Run{ID: runID},
		trrs: []pipeline.TaskRunResult{{Task: &pipeline.MemoTask{}, Result: pipeline.Result{Value: res}}},
		err:  err,
	}
}

// pipelineForStream sets StreamIDs so the data source schedules Observe for this stream.
// Distinct *mockPipeline values are distinct workers; sibling streams in one job must share one pipeline pointer.
func pipelineForStream(streamID streams.StreamID, runID int64, res *big.Int, err error) *mockPipeline {
	p := makePipelineWithSingleResult[*big.Int](runID, res, err)
	p.streamIDs = []streams.StreamID{streamID}
	return p
}

func makeStreamValues(streamIDs ...llotypes.StreamID) lloprotocol.StreamValues {
	if len(streamIDs) == 0 {
		return lloprotocol.StreamValues{
			1: nil,
			2: nil,
			3: nil,
		}
	}
	vals := lloprotocol.StreamValues{}
	for _, streamID := range streamIDs {
		vals[streamID] = nil
	}
	return vals
}

type mockOpts struct {
	verboseLogging       bool
	seqNr                uint64
	configDigest         ocr2types.ConfigDigest
	observationTimestamp time.Time
	lifeCycleStage       llotypes.LifeCycleStage
}

func (m *mockOpts) VerboseLogging() bool { return m.verboseLogging }
func (m *mockOpts) SeqNr() uint64 {
	if m.seqNr == 0 {
		return 1042
	}
	return m.seqNr
}

func (m *mockOpts) ConfigDigest() ocr2types.ConfigDigest {
	if m.configDigest.Hex() == "" {
		return ocr2types.ConfigDigest{6, 5, 4}
	}
	return m.configDigest
}

func (m *mockOpts) ObservationTimestamp() time.Time {
	if m.observationTimestamp.IsZero() {
		return time.Unix(1737936858, 0)
	}
	return m.observationTimestamp
}

func (m *mockOpts) LifeCycleStage() llotypes.LifeCycleStage {
	if m.lifeCycleStage == "" {
		return lloprotocol.LifeCycleStageProduction
	}
	return m.lifeCycleStage
}

type mockTelemeter struct {
	mu                     sync.Mutex
	v3PremiumLegacyPackets []v3PremiumLegacyPacket
	ch                     chan any
}

type v3PremiumLegacyPacket struct {
	run      *pipeline.Run
	trrs     pipeline.TaskRunResults
	streamID uint32
	opts     telem.DSOpts
	val      lloprotocol.StreamValue
	err      error
}

var _ Telemeter = &mockTelemeter{}

func (m *mockTelemeter) EnqueueV3PremiumLegacy(run *pipeline.Run, trrs pipeline.TaskRunResults, streamID uint32, opts telem.DSOpts, val lloprotocol.StreamValue, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.v3PremiumLegacyPackets = append(m.v3PremiumLegacyPackets, v3PremiumLegacyPacket{run, trrs, streamID, opts, val, err})
}

func (m *mockTelemeter) MakeObservationScopedTelemetryCh(opts telem.DSOpts, size int) (ch chan<- any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ch = make(chan any, size)

	return m.ch
}

func (m *mockTelemeter) GetOutcomeTelemetryCh() chan<- *lloprotocol.LLOOutcomeTelemetry {
	return nil
}
func (m *mockTelemeter) GetReportTelemetryCh() chan<- *lloprotocol.LLOReportTelemetry { return nil }
func (m *mockTelemeter) CaptureEATelemetry() bool                                     { return true }

func (m *mockTelemeter) CaptureObservationTelemetry() bool { return true }

var observationTimeout = 500 * time.Millisecond

type addManyCall struct {
	values map[llotypes.StreamID]lloprotocol.StreamValue
	ttl    time.Duration
}

type mockCache struct {
	StreamValueCache
	mu       sync.Mutex
	addCalls []addManyCall
}

func newMockCache(inner StreamValueCache) *mockCache {
	return &mockCache{StreamValueCache: inner}
}

// AddMany is a spy for the StreamValueCache.AddMany method.
// It records the values and ttl passed to it and then calls the underlying StreamValueCache.AddMany method.
func (s *mockCache) AddMany(values map[llotypes.StreamID]lloprotocol.StreamValue, ttl time.Duration) {
	snapshot := make(map[llotypes.StreamID]lloprotocol.StreamValue, len(values))
	maps.Copy(snapshot, values)
	s.mu.Lock()
	s.addCalls = append(s.addCalls, addManyCall{values: snapshot, ttl: ttl})
	s.mu.Unlock()
	s.StreamValueCache.AddMany(values, ttl)
}

func Test_DataSource(t *testing.T) {
	if testing.Short() {
		t.Skip("too slow for testing.Short")
	}

	t.Parallel()
	lggr := logger.NullLogger
	mainCtx := t.Context()
	opts := &mockOpts{}

	t.Run("Observe", func(t *testing.T) {
		t.Parallel()
		t.Run("doesn't set any values if no streams are defined", func(t *testing.T) {
			t.Parallel()
			reg := &mockRegistry{pipelines: make(map[streams.StreamID]*mockPipeline)}
			ds := newDataSource(lggr, reg, telem.NullTelemeter)

			vals := makeStreamValues()
			ctx, cancel := context.WithTimeout(mainCtx, observationTimeout)
			defer cancel()
			err := ds.Observe(ctx, vals, opts)
			require.NoError(t, err)

			assert.Equal(t, makeStreamValues(), vals)
			ds.Close()
		})

		t.Run("observes each stream with success and returns values matching map argument", func(t *testing.T) {
			t.Parallel()
			reg := &mockRegistry{pipelines: make(map[streams.StreamID]*mockPipeline)}
			ds := newDataSource(lggr, reg, telem.NullTelemeter)

			reg.mu.Lock()
			sids := []streams.StreamID{1, 2, 3}
			multi := makePipelineWithMultipleStreamResults(sids, []any{
				decimal.NewFromInt(2181),
				decimal.NewFromInt(40602),
				decimal.NewFromInt(15),
			})
			reg.pipelines[1] = multi
			reg.pipelines[2] = multi
			reg.pipelines[3] = multi
			reg.mu.Unlock()

			vals := makeStreamValues()
			ctx, cancel := context.WithTimeout(mainCtx, observationTimeout)
			defer cancel()
			err := ds.Observe(ctx, vals, opts)
			require.NoError(t, err)

			assert.Equal(t, lloprotocol.StreamValues{
				1: lloprotocol.ToDecimal(decimal.NewFromInt(2181)),
				2: lloprotocol.ToDecimal(decimal.NewFromInt(40602)),
				3: lloprotocol.ToDecimal(decimal.NewFromInt(15)),
			}, vals, "vals: %v", vals)
			ds.Close()
		})

		t.Run("caches nothing for a pipeline whose quote violates the bid/benchmark/ask invariant", func(t *testing.T) {
			t.Parallel()
			reg := &mockRegistry{pipelines: make(map[streams.StreamID]*mockPipeline)}
			ds := newDataSource(lggr, reg, telem.NullTelemeter)

			reg.mu.Lock()
			sids := []streams.StreamID{1, 2, 3}
			multi := makePipelineWithMultipleStreamResults(sids, []any{
				decimal.NewFromInt(2181),
				decimal.NewFromInt(40602),
				map[string]any{
					"streamValueType": int64(lloprotocol.LLOStreamValue_Quote),
					"bid":             "9.9", // above the benchmark
					"benchmark":       "2.2",
					"ask":             "3.3",
				},
			})
			for _, sid := range sids {
				reg.pipelines[sid] = multi
			}
			reg.mu.Unlock()

			vals := makeStreamValues()
			ctx, cancel := context.WithTimeout(mainCtx, observationTimeout)
			defer cancel()
			require.NoError(t, ds.Observe(ctx, vals, opts))

			// The whole group failed, so nothing was cached: every stream of the
			// pipeline stays nil, not just the one with the bad quote.
			assert.Equal(t, makeStreamValues(), vals, "vals: %v", vals)
			ds.Close()
		})

		t.Run("observes each stream and returns success/errors", func(t *testing.T) {
			t.Parallel()
			reg := &mockRegistry{pipelines: make(map[streams.StreamID]*mockPipeline)}
			ds := newDataSource(lggr, reg, telem.NullTelemeter)

			reg.mu.Lock()
			reg.pipelines[11] = pipelineForStream(11, 11, big.NewInt(21810), errors.New("something exploded"))
			reg.pipelines[12] = pipelineForStream(12, 12, big.NewInt(40602), nil)
			reg.pipelines[13] = pipelineForStream(13, 13, nil, errors.New("something exploded 2"))
			reg.mu.Unlock()

			vals := makeStreamValues(11, 12, 13)
			ctx, cancel := context.WithTimeout(mainCtx, observationTimeout)
			defer cancel()

			err := ds.Observe(ctx, vals, opts)
			require.NoError(t, err)

			assert.Equal(t, lloprotocol.StreamValues{
				11: nil,
				12: lloprotocol.ToDecimal(decimal.NewFromInt(40602)),
				13: nil,
			}, vals, "vals: %v", vals)
			ds.Close()
		})

		t.Run("records telemetry", func(t *testing.T) {
			t.Parallel()
			tm := &mockTelemeter{}
			reg := &mockRegistry{pipelines: make(map[streams.StreamID]*mockPipeline)}
			ds := newDataSource(lggr, reg, tm)

			reg.mu.Lock()
			reg.pipelines[21] = pipelineForStream(21, 100, big.NewInt(2181), nil)
			reg.pipelines[22] = pipelineForStream(22, 101, big.NewInt(40602), nil)
			reg.pipelines[23] = pipelineForStream(23, 102, big.NewInt(15), nil)
			reg.mu.Unlock()

			vals := makeStreamValues(21, 22, 23)
			ctx, cancel := context.WithTimeout(mainCtx, observationTimeout)
			defer cancel()

			err := ds.Observe(ctx, vals, opts)
			tm.mu.Lock()
			ch := tm.ch
			tm.mu.Unlock()

			ds.Close()
			require.NoError(t, err)

			assert.Equal(t, lloprotocol.StreamValues{
				21: lloprotocol.ToDecimal(decimal.NewFromInt(2181)),
				22: lloprotocol.ToDecimal(decimal.NewFromInt(40602)),
				23: lloprotocol.ToDecimal(decimal.NewFromInt(15)),
			}, vals, "vals: %v", vals)

			// Get only the last 3 packets, as those would be the result of the first round of observations.
			tm.mu.Lock()
			require.GreaterOrEqual(t, len(tm.v3PremiumLegacyPackets), 3)
			packets := tm.v3PremiumLegacyPackets[:3]
			tm.mu.Unlock()
			m := make(map[int]v3PremiumLegacyPacket)
			for _, pkt := range packets {
				m[int(pkt.run.ID)] = pkt
			}

			pkt := m[100]
			assert.Equal(t, 100, int(pkt.run.ID))
			assert.Len(t, pkt.trrs, 1)
			assert.Equal(t, 21, int(pkt.streamID))
			assert.Equal(t, opts, pkt.opts)
			assert.Equal(t, "2181", pkt.val.(*lloprotocol.Decimal).String())
			require.NoError(t, pkt.err)

			telems := []any{}

			for p := range ch {
				telems = append(telems, p)
				if len(telems) >= 3 {
					break
				}
			}

			require.Len(t, telems, 3)
			sort.Slice(telems, func(i, j int) bool {
				return telems[i].(*telem.LLOObservationTelemetry).StreamId < telems[j].(*telem.LLOObservationTelemetry).StreamId
			})
			require.IsType(t, &telem.LLOObservationTelemetry{}, telems[0])
			obsTelem := telems[0].(*telem.LLOObservationTelemetry)
			assert.Equal(t, uint32(21), obsTelem.StreamId)
			assert.Equal(t, int32(lloprotocol.LLOStreamValue_Decimal), obsTelem.StreamValueType)
			assert.Equal(t, "00000000020885", hex.EncodeToString(obsTelem.StreamValueBinary))
			assert.Equal(t, "2181", obsTelem.StreamValueText)
			assert.Nil(t, obsTelem.ObservationError)
			assert.Equal(t, int64(1737936858000000000), obsTelem.ObservationTimestamp)
			assert.Greater(t, obsTelem.ObservationFinishedAt, int64(1737936858000000000))
			assert.Equal(t, uint32(0), obsTelem.DonId)
			assert.Equal(t, opts.SeqNr(), obsTelem.SeqNr)
			assert.Equal(t, opts.ConfigDigest().Hex(), hex.EncodeToString(obsTelem.ConfigDigest))
		})

		t.Run("records telemetry for errors", func(t *testing.T) {
			t.Parallel()
			tm := &mockTelemeter{}
			reg := &mockRegistry{pipelines: make(map[streams.StreamID]*mockPipeline)}
			ds := newDataSource(lggr, reg, tm)

			reg.mu.Lock()
			reg.pipelines[31] = pipelineForStream(31, 100, big.NewInt(2181), errors.New("something exploded"))
			reg.pipelines[32] = pipelineForStream(32, 101, big.NewInt(40602), nil)
			reg.pipelines[33] = pipelineForStream(33, 102, nil, errors.New("something exploded 2"))
			reg.mu.Unlock()

			vals := makeStreamValues(31, 32, 33)
			ctx, cancel := context.WithTimeout(mainCtx, observationTimeout)
			defer cancel()
			err := ds.Observe(ctx, vals, opts)
			require.NoError(t, err)

			assert.Equal(t, lloprotocol.StreamValues{
				31: nil,
				32: lloprotocol.ToDecimal(decimal.NewFromInt(40602)),
				33: nil,
			}, vals, "vals: %v", vals)

			m := make(map[int]v3PremiumLegacyPacket)
			tm.mu.Lock()
			for _, pkt := range tm.v3PremiumLegacyPackets {
				m[int(pkt.run.ID)] = pkt
			}
			tm.mu.Unlock()
			pkt := m[100]
			assert.Equal(t, 100, int(pkt.run.ID))
			assert.Len(t, pkt.trrs, 1)
			assert.Equal(t, 31, int(pkt.streamID))
			assert.Equal(t, opts, pkt.opts)
			assert.Nil(t, pkt.val)
			require.Error(t, pkt.err)
			ds.Close()
		})

		t.Run("uses cached values when available", func(t *testing.T) {
			t.Parallel()
			reg := &mockRegistry{pipelines: make(map[streams.StreamID]*mockPipeline)}
			ds := newDataSource(lggr, reg, telem.NullTelemeter)

			// First observation to populate cache
			reg.mu.Lock()
			reg.pipelines[10001] = pipelineForStream(10001, 1, big.NewInt(2181), nil)
			reg.pipelines[20001] = pipelineForStream(20001, 2, big.NewInt(40602), nil)
			reg.mu.Unlock()

			vals := lloprotocol.StreamValues{
				10001: nil,
				20001: nil,
				30001: nil,
			}

			ctx, cancel := context.WithTimeout(mainCtx, observationTimeout)
			defer cancel()
			err := ds.Observe(ctx, vals, opts)
			require.NoError(t, err)

			// Verify initial values
			assert.Equal(t, lloprotocol.StreamValues{
				10001: lloprotocol.ToDecimal(decimal.NewFromInt(2181)),
				20001: lloprotocol.ToDecimal(decimal.NewFromInt(40602)),
				30001: nil,
			}, vals)

			// Change pipeline results
			reg.mu.Lock()
			reg.pipelines[10001] = pipelineForStream(10001, 1, big.NewInt(9999), nil)
			reg.pipelines[20001] = pipelineForStream(20001, 2, big.NewInt(8888), nil)
			reg.mu.Unlock()

			// Second observation should use cached values
			vals = lloprotocol.StreamValues{
				10001: nil,
				20001: nil,
				30001: nil,
			}
			ctx2, cancel := context.WithTimeout(mainCtx, observationTimeout)
			defer cancel()
			err = ds.Observe(ctx2, vals, opts)
			require.NoError(t, err)

			// Should still have original values from cache
			assert.Equal(t, lloprotocol.StreamValues{
				10001: lloprotocol.ToDecimal(decimal.NewFromInt(2181)),
				20001: lloprotocol.ToDecimal(decimal.NewFromInt(40602)),
				30001: nil,
			}, vals)
		})

		t.Run("refreshes cache after expiration", func(t *testing.T) {
			t.Parallel()
			reg := &mockRegistry{pipelines: make(map[streams.StreamID]*mockPipeline)}
			ds := newDataSource(lggr, reg, telem.NullTelemeter)

			// First observation
			reg.mu.Lock()
			reg.pipelines[50002] = pipelineForStream(50002, 1, big.NewInt(100), nil)
			reg.mu.Unlock()
			vals := lloprotocol.StreamValues{50002: nil}

			ctx, cancel := context.WithTimeout(mainCtx, observationTimeout)
			defer cancel()
			err := ds.Observe(ctx, vals, opts)
			require.NoError(t, err)

			// Change pipeline result
			reg.mu.Lock()
			reg.pipelines[50002] = pipelineForStream(50002, 1, big.NewInt(200), nil)
			reg.mu.Unlock()

			// Wait for cache to expire
			time.Sleep(observationTimeout * 3)

			// Second observation should use new value
			vals = lloprotocol.StreamValues{50002: nil}
			ctx2, cancel := context.WithTimeout(mainCtx, observationTimeout*5)
			defer cancel()
			err = ds.Observe(ctx2, vals, opts)
			require.NoError(t, err)

			assert.Equal(t, lloprotocol.StreamValues{50002: lloprotocol.ToDecimal(decimal.NewFromInt(200))}, vals)
		})

		t.Run("handles concurrent cache access", func(t *testing.T) {
			t.Parallel()
			// Create a new data source
			reg := &mockRegistry{pipelines: make(map[streams.StreamID]*mockPipeline)}
			ds := newDataSource(lggr, reg, telem.NullTelemeter)

			// Set up pipeline to return different values
			reg.mu.Lock()
			reg.pipelines[1] = pipelineForStream(1, 1, big.NewInt(100), nil)
			reg.mu.Unlock()

			// First observation to cache
			vals := lloprotocol.StreamValues{1: nil}

			ctx, cancel := context.WithTimeout(mainCtx, observationTimeout)
			defer cancel()
			err := ds.Observe(ctx, vals, opts)
			require.NoError(t, err)

			// Run multiple observations concurrently
			var wg sync.WaitGroup
			for range 10 {
				wg.Go(func() {
					vals := lloprotocol.StreamValues{1: nil}
					err := ds.Observe(ctx, vals, opts)
					assert.NoError(t, err)
					assert.Equal(t, lloprotocol.StreamValues{1: lloprotocol.ToDecimal(decimal.NewFromInt(100))}, vals)
				})
			}
			wg.Wait()

			// Verify pipeline was only called once
			assert.Equal(t, int32(1), reg.pipelines[1].runCount.Load())
		})

		t.Run("cache writes are atomic per pipeline group across observation cycles", func(t *testing.T) {
			t.Parallel()
			reg := &mockRegistry{pipelines: make(map[streams.StreamID]*mockPipeline)}
			ds := newDataSource(lggr, reg, telem.NullTelemeter)
			mc := newMockCache(ds.cache)
			ds.cache = mc
			defer ds.Close()

			sids := []streams.StreamID{1, 2, 3}
			partialPipeline := makePipelineWithMultipleStreamResults(sids, []any{decimal.NewFromFloat(100.0), "not-a-number", decimal.NewFromFloat(300.0)})
			reg.mu.Lock()
			reg.pipelines[1] = partialPipeline
			reg.pipelines[2] = partialPipeline
			reg.pipelines[3] = partialPipeline
			reg.mu.Unlock()

			// Cycle 1: partial extraction failure — entire group should be rejected
			vals := makeStreamValues(1, 2, 3)
			ctx, cancel := context.WithTimeout(mainCtx, observationTimeout)
			defer cancel()
			err := ds.Observe(ctx, vals, opts)
			require.NoError(t, err)

			assert.Equal(t, lloprotocol.StreamValues{1: nil, 2: nil, 3: nil}, vals)

			mc.mu.Lock()
			for _, call := range mc.addCalls {
				for _, sid := range sids {
					assert.NotContains(t, call.values, sid)
				}
			}
			mc.mu.Unlock()

			// Explicitly abort Cycle 1's background observation task before mutating pipelines
			cancel()

			// Fix the pipeline with distinct values so we can verify generation
			fixedPipeline := makePipelineWithMultipleStreamResults(sids, []any{decimal.NewFromFloat(111.0), decimal.NewFromFloat(222.0), decimal.NewFromFloat(333.0)})
			reg.mu.Lock()
			reg.pipelines[1] = fixedPipeline
			reg.pipelines[2] = fixedPipeline
			reg.pipelines[3] = fixedPipeline
			reg.mu.Unlock()

			time.Sleep(observationTimeout * 3)

			// Cycle 2: all streams valid — group should be cached atomically
			vals2 := makeStreamValues(1, 2, 3)
			ctx2, cancel2 := context.WithTimeout(mainCtx, observationTimeout)
			defer cancel2()
			err = ds.Observe(ctx2, vals2, opts)
			require.NoError(t, err)

			expectedCycle2 := lloprotocol.StreamValues{
				1: lloprotocol.ToDecimal(decimal.NewFromFloat(111.0)),
				2: lloprotocol.ToDecimal(decimal.NewFromFloat(222.0)),
				3: lloprotocol.ToDecimal(decimal.NewFromFloat(333.0)),
			}
			assert.Equal(t, expectedCycle2, vals2, "cycle 2: expected a value from fixedPipeline")

			// Verify an atomic write of all 3 streams with correct values from the same generation
			mc.mu.Lock()
			defer mc.mu.Unlock()

			foundAtomicWrite := false
			for _, call := range mc.addCalls {
				v1, has1 := call.values[llotypes.StreamID(1)]
				v2, has2 := call.values[llotypes.StreamID(2)]
				v3, has3 := call.values[llotypes.StreamID(3)]
				if has1 && has2 && has3 {
					assert.Equal(t, expectedCycle2[1], v1, "atomic write: stream 1 value mismatch")
					assert.Equal(t, expectedCycle2[2], v2, "atomic write: stream 2 value mismatch")
					assert.Equal(t, expectedCycle2[3], v3, "atomic write: stream 3 value mismatch")
					foundAtomicWrite = true
					break
				}
			}
			assert.True(t, foundAtomicWrite, "expected one AddMany call containing all 3 streams atomically")
		})

		t.Run("handles cache errors gracefully", func(t *testing.T) {
			t.Parallel()
			reg := &mockRegistry{pipelines: make(map[streams.StreamID]*mockPipeline)}
			ds := newDataSource(lggr, reg, telem.NullTelemeter)

			// First observation with error
			reg.mu.Lock()
			reg.pipelines[1] = pipelineForStream(1, 1, nil, errors.New("pipeline error"))
			reg.mu.Unlock()
			vals := makeStreamValues(1)
			ctx, cancel := context.WithTimeout(mainCtx, observationTimeout)
			defer cancel()

			err := ds.Observe(ctx, vals, opts)
			require.NoError(t, err) // Observe returns nil error even if some streams fail

			// Second observation should try again (not use cache for error case)
			reg.mu.Lock()
			reg.pipelines[1] = pipelineForStream(1, 1, big.NewInt(100), nil)
			reg.mu.Unlock()
			time.Sleep(observationTimeout * 3)

			vals = lloprotocol.StreamValues{1: nil}
			ctx2, cancel := context.WithTimeout(mainCtx, observationTimeout*5)
			defer cancel()
			err = ds.Observe(ctx2, vals, opts)
			require.NoError(t, err)

			assert.Equal(t, lloprotocol.StreamValues{1: lloprotocol.ToDecimal(decimal.NewFromInt(100))}, vals)
		})
	})

	promCacheHitCount.Reset()
	promCacheMissCount.Reset()
	promCacheHitEntryAgeMs.Reset()
	promObservationLoopWaitOutcome.Reset()
}

func Test_DataSource_ObservationLoopWakeSkipsPacing(t *testing.T) { //nolint:paralleltest // resets package-level prometheus metrics and relies on loop timing
	promObservationLoopWaitOutcome.Reset()
	lggr := logger.NullLogger
	mainCtx := t.Context()
	opts := &mockOpts{}

	reg := &mockRegistry{pipelines: make(map[streams.StreamID]*mockPipeline)}
	reg.mu.Lock()
	reg.pipelines[1] = pipelineForStream(1, 1, big.NewInt(42), nil)
	reg.mu.Unlock()

	ds := newDataSource(lggr, reg, telem.NullTelemeter)
	defer ds.Close()

	// Long plugin deadline => large inter-iteration pacing; wake from Observe should advance the loop without waiting.
	longCtx, cancel := context.WithTimeout(mainCtx, 30*time.Second)
	defer cancel()
	vals := makeStreamValues(1)
	require.NoError(t, ds.Observe(longCtx, vals, opts))

	require.Eventually(t, func() bool {
		return promtest.ToFloat64(promObservationLoopWaitOutcome.WithLabelValues("wake")) >= 1
	}, 2*time.Second, 5*time.Millisecond, "expected at least one pacing wait satisfied by plugin wake")
}

func Test_DataSource_ObserveWakeManyConcurrent(t *testing.T) {
	t.Parallel()
	lggr := logger.NullLogger
	mainCtx := t.Context()
	opts := &mockOpts{}

	reg := &mockRegistry{pipelines: make(map[streams.StreamID]*mockPipeline)}
	reg.mu.Lock()
	reg.pipelines[1] = pipelineForStream(1, 1, big.NewInt(1), nil)
	reg.mu.Unlock()

	ds := newDataSource(lggr, reg, telem.NullTelemeter)
	ctx, cancel := context.WithTimeout(mainCtx, observationTimeout)
	defer cancel()
	vals := makeStreamValues(1)
	require.NoError(t, ds.Observe(ctx, vals, opts))

	done := make(chan struct{})
	var wg sync.WaitGroup
	for range 200 {
		wg.Go(func() {
			// Each call needs its own StreamValues map: Observe mutates it in place (UpdateStreamValues).
			localVals := makeStreamValues(1)
			_ = ds.Observe(ctx, localVals, opts)
		})
	}
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("concurrent Observe calls did not complete")
	}
	require.NoError(t, ds.Close())
}

func Test_buildStreamsRefreshPlan(t *testing.T) {
	t.Parallel()
	lggr := logger.NullLogger
	timeout := 100 * time.Millisecond

	pipelineABC := &mockPipeline{streamIDs: []streams.StreamID{1, 2, 3}}
	pipelineSingle := &mockPipeline{streamIDs: []streams.StreamID{10}}
	pipelineDE := &mockPipeline{streamIDs: []streams.StreamID{20, 21}}

	reg := &mockRegistry{pipelines: map[streams.StreamID]*mockPipeline{
		1: pipelineABC, 2: pipelineABC, 3: pipelineABC,
		10: pipelineSingle,
		20: pipelineDE, 21: pipelineDE,
	}}

	t.Run("all streams stale returns all", func(t *testing.T) {
		t.Parallel()
		cache := NewCache(0)
		staleTTL := 1 * time.Millisecond
		cache.Add(1, lloprotocol.ToDecimal(decimal.NewFromInt(100)), staleTTL)
		cache.Add(2, lloprotocol.ToDecimal(decimal.NewFromInt(200)), staleTTL)
		cache.Add(3, lloprotocol.ToDecimal(decimal.NewFromInt(300)), staleTTL)
		ds := &dataSource{lggr: lggr, registry: reg, cache: cache}
		sv := lloprotocol.StreamValues{1: nil, 2: nil, 3: nil}

		result := ds.buildStreamsRefreshPlan(sv, timeout, lggr).streamIDsToRefresh

		assert.Len(t, result, 3)
		for _, id := range []streams.StreamID{1, 2, 3} {
			assert.Contains(t, result, id)
		}
	})

	t.Run("all streams fresh in cache, returns none", func(t *testing.T) {
		t.Parallel()
		cache := NewCache(0)
		cache.Add(1, lloprotocol.ToDecimal(decimal.NewFromInt(100)), time.Hour)
		cache.Add(2, lloprotocol.ToDecimal(decimal.NewFromInt(200)), time.Hour)
		cache.Add(3, lloprotocol.ToDecimal(decimal.NewFromInt(300)), time.Hour)
		ds := &dataSource{lggr: lggr, registry: reg, cache: cache}
		sv := lloprotocol.StreamValues{1: nil, 2: nil, 3: nil}

		result := ds.buildStreamsRefreshPlan(sv, timeout, lggr).streamIDsToRefresh

		assert.Empty(t, result)
	})

	t.Run("one stale driver lists only stale IDs; worker observes all requested streams on that pipeline", func(t *testing.T) {
		t.Parallel()
		cache := NewCache(0)
		cache.Add(1, lloprotocol.ToDecimal(decimal.NewFromInt(100)), time.Hour)
		cache.Add(2, lloprotocol.ToDecimal(decimal.NewFromInt(200)), 1*time.Millisecond)
		cache.Add(3, lloprotocol.ToDecimal(decimal.NewFromInt(300)), time.Hour)
		ds := &dataSource{lggr: lggr, registry: reg, cache: cache}
		sv := lloprotocol.StreamValues{1: nil, 2: nil, 3: nil}

		plan := ds.buildStreamsRefreshPlan(sv, timeout, lggr)

		assert.Equal(t, []streams.StreamID{2}, plan.streamIDsToRefresh, "streamIDsToRefresh is stale plugin-scope keys only")
		require.Len(t, plan.groups, 1)
		for _, sids := range plan.groups {
			assert.ElementsMatch(t, []streams.StreamID{1, 2, 3}, sids, "all three are in plugin scope, so observe list matches pipeline StreamIDs()")
		}
	})

	t.Run("staleStreamIDs lists only stale keys; groups intersect pipeline with plugin scope", func(t *testing.T) {
		t.Parallel()
		cache := NewCache(0)
		cache.Add(1, lloprotocol.ToDecimal(decimal.NewFromInt(100)), 1*time.Millisecond)
		cache.Add(2, lloprotocol.ToDecimal(decimal.NewFromInt(200)), time.Hour)
		// pipeline has {1,2,3}, but only {1,2} in plugin scope
		ds := &dataSource{lggr: lggr, registry: reg, cache: cache}
		sv := lloprotocol.StreamValues{1: nil, 2: nil} // stream 3 not requested

		plan := ds.buildStreamsRefreshPlan(sv, timeout, lggr)

		assert.Equal(t, []streams.StreamID{1}, plan.streamIDsToRefresh)
		assert.NotContains(t, plan.streamIDsToRefresh, streams.StreamID(3), "out-of-scope stream is not a refresh driver")
		require.Len(t, plan.groups, 1)
		for _, sids := range plan.groups {
			assert.Equal(t, []streams.StreamID{1, 2}, sids, "stream 3 is not requested; observe list is intersection with streamValues")
		}
	})

	t.Run("stream not in registry is stale driver only; no pipeline worker", func(t *testing.T) {
		t.Parallel()
		ds := &dataSource{lggr: lggr, registry: reg, cache: NewCache(0)}
		sv := lloprotocol.StreamValues{999: nil} // plugin requested streamId not yet in registry

		plan := ds.buildStreamsRefreshPlan(sv, timeout, lggr)

		assert.Equal(t, []streams.StreamID{999}, plan.streamIDsToRefresh, "stale in-scope keys are listed even without registry")
		assert.Empty(t, plan.groups, "no Observe workers without a pipeline")
		assert.ElementsMatch(t, []streams.StreamID{999}, plan.missingStreamIDs)
	})

	t.Run("empty streamValues returns empty set", func(t *testing.T) {
		t.Parallel()
		ds := &dataSource{lggr: lggr, registry: reg, cache: NewCache(0)}
		sv := lloprotocol.StreamValues{}

		result := ds.buildStreamsRefreshPlan(sv, timeout, lggr).streamIDsToRefresh

		assert.Empty(t, result)
	})

	t.Run("multiple pipelines: only stale keys appear in streamIDsToRefresh", func(t *testing.T) {
		t.Parallel()
		cache := NewCache(0)
		// Pipeline {10}: all fresh
		cache.Add(10, lloprotocol.ToDecimal(decimal.NewFromInt(100)), time.Hour)
		// Pipeline {20,21}: stream 20 stale, stream 21 fresh
		cache.Add(20, lloprotocol.ToDecimal(decimal.NewFromInt(2000)), 1*time.Millisecond)
		cache.Add(21, lloprotocol.ToDecimal(decimal.NewFromInt(2100)), time.Hour)

		ds := &dataSource{lggr: lggr, registry: reg, cache: cache}
		sv := lloprotocol.StreamValues{10: nil, 20: nil, 21: nil}

		plan := ds.buildStreamsRefreshPlan(sv, timeout, lggr)

		assert.NotContains(t, plan.streamIDsToRefresh, streams.StreamID(10), "fresh pipeline should not be refreshed")
		assert.Equal(t, []streams.StreamID{20}, plan.streamIDsToRefresh)
		assert.NotContains(t, plan.streamIDsToRefresh, streams.StreamID(21), "fresh sibling is not a refresh driver")
		require.Len(t, plan.groups, 1)
		for _, sids := range plan.groups {
			assert.ElementsMatch(t, []streams.StreamID{20, 21}, sids)
		}
	})

	promCacheHitCount.Reset()
	promCacheMissCount.Reset()
	promCacheHitEntryAgeMs.Reset()
	promObservationLoopWaitOutcome.Reset()
}

func Test_observationTuningHelpers(t *testing.T) {
	t.Parallel()

	const tuningTestT = 100 * time.Millisecond
	wantStaleSkip := time.Duration(staleRefreshRemainingNumerator) * tuningTestT / time.Duration(staleRefreshRemainingDenominator)

	assert.Equal(t, time.Duration(cacheTTLMultiplier)*tuningTestT, cacheEntryTTL(tuningTestT))
	assert.Equal(t, wantStaleSkip, staleRefreshSkipThreshold(tuningTestT))
	assert.Less(t, staleRefreshSkipThreshold(tuningTestT), cacheEntryTTL(tuningTestT))
	assert.Less(t, staleRefreshSkipThreshold(tuningTestT)+observationLoopPacing(tuningTestT), cacheEntryTTL(tuningTestT))

	// With num/den = 13/5 the invariant cap (cacheTTL−stale−1ns = 7/5·T−1ns) exceeds raw T/2, so pacing is bounded by
	// T/observationLoopPacingDivisor (= T/2), not the invariant.
	assert.Equal(t, 100*time.Millisecond/observationLoopPacingDivisor, observationLoopPacing(100*time.Millisecond))
	assert.Equal(t, 500*time.Millisecond/observationLoopPacingDivisor, observationLoopPacing(500*time.Millisecond))
	assert.Equal(t, observationLoopPacingFloor, observationLoopPacing(0))
	// T/2 = 15ms here, still above the floor and below the invariant cap.
	assert.Equal(t, 30*time.Millisecond/observationLoopPacingDivisor, observationLoopPacing(30*time.Millisecond))
	// Confirm the invariant cap is no longer the binding constraint (raw T/2 < cacheTTL−stale−1ns).
	invMax := cacheEntryTTL(100*time.Millisecond) - staleRefreshSkipThreshold(100*time.Millisecond) - time.Nanosecond
	assert.Greater(t, invMax, 100*time.Millisecond/observationLoopPacingDivisor)
}

func BenchmarkObserve(b *testing.B) {
	lggr := logger.TestLogger(b)
	ctx := b.Context()
	// can enable/disable verbose logging to test performance here
	opts := &mockOpts{verboseLogging: true}

	db := pgtest.NewSqlxDB(b)
	bridgesORM := bridges.NewORM(db)

	if b.N > math.MaxInt32 {
		b.Fatalf("N is too large: %d", b.N)
	}

	n := uint32(b.N) //nolint:gosec // G115 // overflow impossible

	createBridge(b, "foo-bridge", `123.456`, bridgesORM, 0)
	createBridge(b, "bar-bridge", `"124.456"`, bridgesORM, 0)

	c := clhttptest.NewTestLocalOnlyHTTPClient()
	runner := pipeline.NewRunner(
		nil,
		bridgesORM,
		&mockPipelineConfig{},
		&mockBridgeConfig{},
		nil,
		nil,
		nil,
		lggr,
		c,
		c,
	)

	r := streams.NewRegistry(lggr, runner)
	for i := range n {
		jb := job.Job{
			ID:       int32(i),
			Name:     null.StringFrom(fmt.Sprintf("job-%d", i)),
			Type:     job.Stream,
			StreamID: &i,
			PipelineSpec: &pipeline.Spec{
				ID: int32(i * 100),
				DotDagSource: fmt.Sprintf(`
// Benchmark Price
result1          [type=memo value="900.0022"];
multiply2 	  	 [type=multiply times=1 streamID=%d index=0]; // force conversion to decimal

result2          [type=bridge name="foo-bridge" requestData="{\"data\":{\"data\":\"foo\"}}"];
result2_parse    [type=jsonparse path="result" streamID=%d index=1];

result3          [type=bridge name="bar-bridge" requestData="{\"data\":{\"data\":\"bar\"}}"];
result3_parse    [type=jsonparse path="result"];
multiply3 	  	 [type=multiply times=1 streamID=%d index=2]; // force conversion to decimal

result1 -> multiply2;
result2 -> result2_parse;
result3 -> result3_parse -> multiply3;
`, i+n, i+2*n, i+3*n),
			},
		}
		err := r.Register(jb, nil)
		require.NoError(b, err)
	}

	ds := newDataSource(lggr, r, telem.NullTelemeter)
	vals := make(map[llotypes.StreamID]lloprotocol.StreamValue)
	for i := range 4 * n {
		vals[i] = nil
	}

	b.ResetTimer()
	err := ds.Observe(ctx, vals, opts)
	require.NoError(b, err)
	ds.Close()
}

// Test_DataSource_inProduction covers the lifecycle gate: only a Production
// instance observes; staging/retired/unknown and nil opts do not.
func Test_DataSource_inProduction(t *testing.T) {
	t.Parallel()
	reg := &mockRegistry{pipelines: make(map[streams.StreamID]*mockPipeline)}
	ds := newDataSource(logger.NullLogger, reg, telem.NullTelemeter)
	defer ds.Close()

	require.True(t, ds.inProduction(&mockOpts{lifeCycleStage: lloprotocol.LifeCycleStageProduction}))
	require.False(t, ds.inProduction(&mockOpts{lifeCycleStage: lloprotocol.LifeCycleStageStaging}))
	require.False(t, ds.inProduction(&mockOpts{lifeCycleStage: lloprotocol.LifeCycleStageRetired}))
	require.False(t, ds.inProduction(nil))
}

// Test_DataSource_StagingDoesNotObserve asserts a non-Production instance runs
// no pipelines and returns unset stream values.
func Test_DataSource_StagingDoesNotObserve(t *testing.T) {
	t.Parallel()
	reg := &mockRegistry{pipelines: make(map[streams.StreamID]*mockPipeline)}
	reg.pipelines[1] = pipelineForStream(1, 1, big.NewInt(42), nil)
	ds := newDataSource(logger.NullLogger, reg, telem.NullTelemeter)
	defer ds.Close()

	ctx, cancel := context.WithTimeout(t.Context(), observationTimeout)
	defer cancel()
	vals := makeStreamValues(1)
	require.NoError(t, ds.Observe(ctx, vals, &mockOpts{lifeCycleStage: lloprotocol.LifeCycleStageStaging}))

	require.Nil(t, vals[1], "staging instance must not populate stream values")
	require.Zero(t, reg.pipelines[1].runCount.Load(), "staging instance must not run pipelines")
}
