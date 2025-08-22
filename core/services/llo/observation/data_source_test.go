package observation

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	llotypes "github.com/smartcontractkit/chainlink-common/pkg/types/llo"
	"github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
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

	runCount int
}

func (m *mockPipeline) Run(ctx context.Context) (*pipeline.Run, pipeline.TaskRunResults, error) {
	m.runCount++
	return m.run, m.trrs, m.err
}

func (m *mockPipeline) StreamIDs() []streams.StreamID {
	return m.streamIDs
}

type mockRegistry struct {
	pipelines map[streams.StreamID]*mockPipeline
}

func (m *mockRegistry) Get(streamID streams.StreamID) (p streams.Pipeline, exists bool) {
	p, exists = m.pipelines[streamID]
	return
}

func makePipelineWithSingleResult[T any](runID int64, res T, err error) *mockPipeline {
	return &mockPipeline{
		run:  &pipeline.Run{ID: runID},
		trrs: []pipeline.TaskRunResult{{Task: &pipeline.MemoTask{}, Result: pipeline.Result{Value: res}}},
		err:  err,
	}
}

func makeStreamValues(streamIDs ...llotypes.StreamID) llo.StreamValues {
	if len(streamIDs) == 0 {
		return llo.StreamValues{
			1: nil,
			2: nil,
			3: nil,
		}
	}
	vals := llo.StreamValues{}
	for _, streamID := range streamIDs {
		vals[streamID] = nil
	}
	return vals
}

type mockOpts struct {
	verboseLogging       bool
	seqNr                uint64
	outCtx               ocr3types.OutcomeContext
	configDigest         ocr2types.ConfigDigest
	observationTimestamp time.Time
}

func (m *mockOpts) VerboseLogging() bool { return m.verboseLogging }
func (m *mockOpts) SeqNr() uint64 {
	if m.seqNr == 0 {
		return 1042
	}
	return m.seqNr
}
func (m *mockOpts) OutCtx() ocr3types.OutcomeContext {
	if m.outCtx.SeqNr == 0 {
		return ocr3types.OutcomeContext{SeqNr: 1042, PreviousOutcome: []byte("foo")}
	}
	return m.outCtx
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
func (m *mockOpts) OutcomeCodec() llo.OutcomeCodec {
	return mockOutputCodec{}
}

type mockOutputCodec struct{}

func (oc mockOutputCodec) Encode(outcome llo.Outcome) (ocr3types.Outcome, error) {
	return ocr3types.Outcome{}, nil
}
func (oc mockOutputCodec) Decode(encoded ocr3types.Outcome) (outcome llo.Outcome, err error) {
	return llo.Outcome{
		LifeCycleStage: llo.LifeCycleStageProduction,
	}, nil
}

type mockTelemeter struct {
	mu                     sync.Mutex
	v3PremiumLegacyPackets []v3PremiumLegacyPacket
	ch                     chan interface{}
}

type v3PremiumLegacyPacket struct {
	run      *pipeline.Run
	trrs     pipeline.TaskRunResults
	streamID uint32
	opts     llo.DSOpts
	val      llo.StreamValue
	err      error
}

var _ Telemeter = &mockTelemeter{}

func (m *mockTelemeter) EnqueueV3PremiumLegacy(run *pipeline.Run, trrs pipeline.TaskRunResults, streamID uint32, opts llo.DSOpts, val llo.StreamValue, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.v3PremiumLegacyPackets = append(m.v3PremiumLegacyPackets, v3PremiumLegacyPacket{run, trrs, streamID, opts, val, err})
}
func (m *mockTelemeter) MakeObservationScopedTelemetryCh(opts llo.DSOpts, size int) (ch chan<- interface{}) {
	m.ch = make(chan interface{}, size)
	return m.ch
}
func (m *mockTelemeter) GetOutcomeTelemetryCh() chan<- *llo.LLOOutcomeTelemetry {
	return nil
}
func (m *mockTelemeter) GetReportTelemetryCh() chan<- *llo.LLOReportTelemetry { return nil }
func (m *mockTelemeter) CaptureEATelemetry() bool                             { return true }
func (m *mockTelemeter) CaptureObservationTelemetry() bool                    { return true }

func Test_DataSource(t *testing.T) {
	lggr := logger.TestLogger(t)
	ctx := testutils.Context(t)
	opts := &mockOpts{}
	observationLoopSleepDuration := 10 * time.Millisecond

	t.Run("Observe", func(t *testing.T) {
		t.Run("doesn't set any values if no streams are defined", func(t *testing.T) {
			reg := &mockRegistry{make(map[streams.StreamID]*mockPipeline)}
			ds := newDataSource(lggr, reg, telem.NullTelemeter, true)

			vals := makeStreamValues()
			err := ds.Observe(ctx, vals, opts)
			assert.NoError(t, err)

			assert.Equal(t, makeStreamValues(), vals)
			ds.Close()
		})

		t.Run("observes each stream with success and returns values matching map argument", func(t *testing.T) {
			reg := &mockRegistry{make(map[streams.StreamID]*mockPipeline)}
			ds := newDataSource(lggr, reg, telem.NullTelemeter, true)

			reg.pipelines[1] = makePipelineWithSingleResult[*big.Int](1, big.NewInt(2181), nil)
			reg.pipelines[2] = makePipelineWithSingleResult[*big.Int](2, big.NewInt(40602), nil)
			reg.pipelines[3] = makePipelineWithSingleResult[*big.Int](3, big.NewInt(15), nil)

			// Sleep, so the observation loop has time to process the new pipelines.
			time.Sleep(2 * observationLoopSleepDuration)

			vals := makeStreamValues()
			err := ds.Observe(ctx, vals, opts)
			assert.NoError(t, err)

			assert.Equal(t, llo.StreamValues{
				1: llo.ToDecimal(decimal.NewFromInt(2181)),
				2: llo.ToDecimal(decimal.NewFromInt(40602)),
				3: llo.ToDecimal(decimal.NewFromInt(15)),
			}, vals, fmt.Sprintf("vals: %v", vals))
			ds.Close()
		})

		t.Run("observes each stream and returns success/errors", func(t *testing.T) {
			reg := &mockRegistry{make(map[streams.StreamID]*mockPipeline)}
			ds := newDataSource(lggr, reg, telem.NullTelemeter, true)

			reg.pipelines[11] = makePipelineWithSingleResult[*big.Int](11, big.NewInt(21810), errors.New("something exploded"))
			reg.pipelines[12] = makePipelineWithSingleResult[*big.Int](12, big.NewInt(40602), nil)
			reg.pipelines[13] = makePipelineWithSingleResult[*big.Int](13, nil, errors.New("something exploded 2"))

			// Sleep, so the observation loop has time to process the new pipelines.
			time.Sleep(observationLoopSleepDuration)

			vals := makeStreamValues(11, 12, 13)
			err := ds.Observe(ctx, vals, opts)
			assert.NoError(t, err)

			assert.Equal(t, llo.StreamValues{
				11: nil,
				12: llo.ToDecimal(decimal.NewFromInt(40602)),
				13: nil,
			}, vals, fmt.Sprintf("vals: %v", vals))
			ds.Close()
		})

		t.Run("records telemetry", func(t *testing.T) {
			tm := &mockTelemeter{}
			reg := &mockRegistry{make(map[streams.StreamID]*mockPipeline)}
			ds := newDataSource(lggr, reg, tm, true)

			reg.pipelines[21] = makePipelineWithSingleResult[*big.Int](100, big.NewInt(2181), nil)
			reg.pipelines[22] = makePipelineWithSingleResult[*big.Int](101, big.NewInt(40602), nil)
			reg.pipelines[23] = makePipelineWithSingleResult[*big.Int](102, big.NewInt(15), nil)

			// Sleep, so the observation loop has time to process the new pipelines.
			time.Sleep(2 * observationLoopSleepDuration)

			vals := makeStreamValues(21, 22, 23)
			err := ds.Observe(ctx, vals, opts)
			require.NoError(t, err)

			assert.Equal(t, llo.StreamValues{
				21: llo.ToDecimal(decimal.NewFromInt(2181)),
				22: llo.ToDecimal(decimal.NewFromInt(40602)),
				23: llo.ToDecimal(decimal.NewFromInt(15)),
			}, vals, fmt.Sprintf("vals: %v", vals))

			// Get only the last 3 packets, as those would be the result of the first round of observations.
			packets := tm.v3PremiumLegacyPackets[:3]
			m := make(map[int]v3PremiumLegacyPacket)
			for _, pkt := range packets {
				m[int(pkt.run.ID)] = pkt
			}
			pkt := m[100]
			assert.Equal(t, 100, int(pkt.run.ID))
			assert.Len(t, pkt.trrs, 1)
			assert.Equal(t, 21, int(pkt.streamID))
			assert.Equal(t, opts, pkt.opts)
			assert.Equal(t, "2181", pkt.val.(*llo.Decimal).String())
			require.NoError(t, pkt.err)

			telems := []interface{}{}
			for p := range tm.ch {
				telems = append(telems, p)
				if len(telems) >= 3 {
					break
				}
			}

			require.Len(t, telems[:3], 3)
			sort.Slice(telems, func(i, j int) bool {
				return telems[i].(*telem.LLOObservationTelemetry).StreamId < telems[j].(*telem.LLOObservationTelemetry).StreamId
			})
			require.IsType(t, &telem.LLOObservationTelemetry{}, telems[0])
			obsTelem := telems[0].(*telem.LLOObservationTelemetry)
			assert.Equal(t, uint32(21), obsTelem.StreamId)
			assert.Equal(t, int32(llo.LLOStreamValue_Decimal), obsTelem.StreamValueType)
			assert.Equal(t, "00000000020885", hex.EncodeToString(obsTelem.StreamValueBinary))
			assert.Equal(t, "2181", obsTelem.StreamValueText)
			assert.Nil(t, obsTelem.ObservationError)
			assert.Equal(t, int64(1737936858000000000), obsTelem.ObservationTimestamp)
			assert.Greater(t, obsTelem.ObservationFinishedAt, int64(1737936858000000000))
			assert.Equal(t, uint32(0), obsTelem.DonId)
			assert.Equal(t, opts.SeqNr(), obsTelem.SeqNr)
			assert.Equal(t, opts.ConfigDigest().Hex(), hex.EncodeToString(obsTelem.ConfigDigest))
			ds.Close()
		})

		t.Run("records telemetry for errors", func(t *testing.T) {
			tm := &mockTelemeter{}
			reg := &mockRegistry{make(map[streams.StreamID]*mockPipeline)}
			ds := newDataSource(lggr, reg, tm, true)

			reg.pipelines[31] = makePipelineWithSingleResult[*big.Int](100, big.NewInt(2181), errors.New("something exploded"))
			reg.pipelines[32] = makePipelineWithSingleResult[*big.Int](101, big.NewInt(40602), nil)
			reg.pipelines[33] = makePipelineWithSingleResult[*big.Int](102, nil, errors.New("something exploded 2"))

			vals := makeStreamValues(31, 32, 33)
			err := ds.Observe(ctx, vals, opts)
			require.NoError(t, err)

			assert.Equal(t, llo.StreamValues{
				31: nil,
				32: llo.ToDecimal(decimal.NewFromInt(40602)),
				33: nil,
			}, vals, fmt.Sprintf("vals: %v", vals))

			require.Len(t, tm.v3PremiumLegacyPackets, 3)
			m := make(map[int]v3PremiumLegacyPacket)
			for _, pkt := range tm.v3PremiumLegacyPackets {
				m[int(pkt.run.ID)] = pkt
			}
			pkt := m[100]
			assert.Equal(t, 100, int(pkt.run.ID))
			assert.Len(t, pkt.trrs, 1)
			assert.Equal(t, 31, int(pkt.streamID))
			assert.Equal(t, opts, pkt.opts)
			assert.Nil(t, pkt.val)
			assert.Error(t, pkt.err)
			ds.Close()
		})

	})
}

func BenchmarkObserve(b *testing.B) {
	lggr := logger.TestLogger(b)
	ctx := testutils.Context(b)
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
	for i := uint32(0); i < n; i++ {
		i := i
		jb := job.Job{
			ID:       int32(i), //nolint:gosec // G115 // overflow impossible
			Name:     null.StringFrom(fmt.Sprintf("job-%d", i)),
			Type:     job.Stream,
			StreamID: &i,
			PipelineSpec: &pipeline.Spec{
				ID: int32(i * 100), //nolint:gosec // G115 // overflow impossible
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

	ds := newDataSource(lggr, r, telem.NullTelemeter, false)
	vals := make(map[llotypes.StreamID]llo.StreamValue)
	for i := uint32(0); i < 4*n; i++ {
		vals[i] = nil
	}

	b.ResetTimer()
	err := ds.Observe(ctx, vals, opts)
	require.NoError(b, err)
	ds.Close()
}
