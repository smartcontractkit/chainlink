package observation

import (
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"golang.org/x/sync/errgroup"
	"gopkg.in/guregu/null.v4"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	clnull "github.com/smartcontractkit/chainlink-common/pkg/utils/null"
	lloprotocol "github.com/smartcontractkit/chainlink-data-streams/llo/protocol"
	llov30 "github.com/smartcontractkit/chainlink-data-streams/llo/v30"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	clhttptest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/httptest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/streams"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func makeErroringPipeline() *mockPipeline {
	return &mockPipeline{
		err: errors.New("pipeline error"),
	}
}

func makePipelineWithMultipleStreamResults(streamIDs []streams.StreamID, results []any) *mockPipeline {
	if len(streamIDs) != len(results) {
		panic("streamIDs and results must have the same length")
	}
	trrs := make([]pipeline.TaskRunResult, len(streamIDs))
	for i, res := range results {
		trrs[i] = pipeline.TaskRunResult{Task: &pipeline.MemoTask{BaseTask: pipeline.BaseTask{StreamID: clnull.Uint32From(streamIDs[i])}}, Result: pipeline.Result{Value: res}}
	}
	return &mockPipeline{
		run:       &pipeline.Run{},
		trrs:      trrs,
		err:       nil,
		streamIDs: streamIDs,
	}
}

func TestObservationContext_Observe(t *testing.T) { //nolint:paralleltest // subtests share one ObservationContext and pipeline run counters
	ctx := t.Context()
	r := &mockRegistry{}
	telem := &mockTelemeter{}
	lggr := logger.TestLogger(t)
	oc := newObservationContext(lggr, r, telem)
	opts := llov30.DSOpts(nil)

	missingStreamID := streams.StreamID(0)
	streamID1 := streams.StreamID(1)
	streamID2 := streams.StreamID(2)
	streamID3 := streams.StreamID(3)
	streamID4 := streams.StreamID(4)
	streamID5 := streams.StreamID(5)
	streamID6 := streams.StreamID(6)
	streamID7 := streams.StreamID(7)
	streamID8 := streams.StreamID(8)

	multiPipelineDecimal := makePipelineWithMultipleStreamResults([]streams.StreamID{streamID4, streamID5, streamID6}, []any{decimal.NewFromFloat(12.34), decimal.NewFromFloat(56.78), decimal.NewFromFloat(90.12)})

	streamID9 := streams.StreamID(9)
	streamID10 := streams.StreamID(10)
	streamID11 := streams.StreamID(11)
	multiPipelinePartialFail := &mockPipeline{
		run: &pipeline.Run{},
		trrs: []pipeline.TaskRunResult{
			{Task: &pipeline.MemoTask{BaseTask: pipeline.BaseTask{StreamID: clnull.Uint32From(streamID9)}}, Result: pipeline.Result{Value: decimal.NewFromFloat(100.5)}},
			{Task: &pipeline.MemoTask{BaseTask: pipeline.BaseTask{StreamID: clnull.Uint32From(streamID10)}}, Result: pipeline.Result{Value: "not-a-number"}},
			{Task: &pipeline.MemoTask{BaseTask: pipeline.BaseTask{StreamID: clnull.Uint32From(streamID11)}}, Result: pipeline.Result{Value: decimal.NewFromFloat(200.5)}},
		},
		streamIDs: []streams.StreamID{streamID9, streamID10, streamID11},
	}

	r.pipelines = map[streams.StreamID]*mockPipeline{
		streamID1:  {},
		streamID2:  makePipelineWithSingleResult[decimal.Decimal](rand.Int64(), decimal.NewFromFloat(12.34), nil),
		streamID3:  makeErroringPipeline(),
		streamID4:  multiPipelineDecimal,
		streamID5:  multiPipelineDecimal,
		streamID6:  multiPipelineDecimal,
		streamID7:  makePipelineWithSingleResult[float64](rand.Int64(), 1.23, nil),
		streamID8:  makePipelineWithSingleResult[int64](rand.Int64(), 5, nil),
		streamID9:  multiPipelinePartialFail,
		streamID10: multiPipelinePartialFail,
		streamID11: multiPipelinePartialFail,
	}

	t.Run("returns error in case of missing pipeline", func(t *testing.T) { //nolint:paralleltest // shares ObservationContext setup
		_, err := oc.Observe(ctx, missingStreamID, opts)
		require.EqualError(t, err, "no pipeline for stream: 0")
	})
	t.Run("returns error in case of zero results", func(t *testing.T) { //nolint:paralleltest // shares ObservationContext setup
		_, err := oc.Observe(ctx, streamID1, opts)
		require.EqualError(t, err, "invalid number of results, expected: 1 or 3, got: 0")
	})
	t.Run("returns composite value from legacy job with single top-level streamID", func(t *testing.T) { //nolint:paralleltest // shares ObservationContext setup
		val, err := oc.Observe(ctx, streamID2, opts)
		require.NoError(t, err)

		assert.Equal(t, "12.34", val.(*lloprotocol.Decimal).String())
	})
	t.Run("returns error in case of erroring pipeline", func(t *testing.T) { //nolint:paralleltest // shares ObservationContext setup
		_, err := oc.Observe(ctx, streamID3, opts)
		require.EqualError(t, err, "pipeline error")
	})
	t.Run("returns values for multiple stream IDs within the same job based on streamID tag with a single pipeline execution", func(t *testing.T) { //nolint:paralleltest // shares ObservationContext setup
		val, err := oc.Observe(ctx, streamID4, opts)
		require.NoError(t, err)
		assert.Equal(t, "12.34", val.(*lloprotocol.Decimal).String())

		val, err = oc.Observe(ctx, streamID5, opts)
		require.NoError(t, err)
		assert.Equal(t, "56.78", val.(*lloprotocol.Decimal).String())

		val, err = oc.Observe(ctx, streamID6, opts)
		require.NoError(t, err)
		assert.Equal(t, "90.12", val.(*lloprotocol.Decimal).String())

		assert.Equal(t, int32(1), multiPipelineDecimal.runCount.Load())

		// returns cached values on subsequent calls
		val, err = oc.Observe(ctx, streamID6, opts)
		require.NoError(t, err)
		assert.Equal(t, "90.12", val.(*lloprotocol.Decimal).String())

		assert.Equal(t, int32(1), multiPipelineDecimal.runCount.Load())
	})
	t.Run("returns value from float64 value", func(t *testing.T) { //nolint:paralleltest // shares ObservationContext setup
		val, err := oc.Observe(ctx, streamID7, opts)
		require.NoError(t, err)

		assert.Equal(t, "1.23", val.(*lloprotocol.Decimal).String())
	})
	t.Run("returns value from int64 value", func(t *testing.T) { //nolint:paralleltest // shares ObservationContext setup
		val, err := oc.Observe(ctx, streamID8, opts)
		require.NoError(t, err)

		assert.Equal(t, "5", val.(*lloprotocol.Decimal).String())
	})
	t.Run("partial extraction failure in multi-stream pipeline", func(t *testing.T) { //nolint:paralleltest // shares ObservationContext setup
		val, err := oc.Observe(ctx, streamID9, opts)
		require.NoError(t, err)
		assert.Equal(t, "100.5", val.(*lloprotocol.Decimal).String())

		_, err = oc.Observe(ctx, streamID10, opts)
		require.Error(t, err, "unparseable value should fail extraction")

		val, err = oc.Observe(ctx, streamID11, opts)
		require.NoError(t, err)
		assert.Equal(t, "200.5", val.(*lloprotocol.Decimal).String())

		assert.Equal(t, int32(1), multiPipelinePartialFail.runCount.Load())
	})
}

func TestObservationContext_Observe_concurrencyStressTest(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	r := &mockRegistry{}
	telem := &mockTelemeter{}
	lggr := logger.TestLogger(t)
	oc := newObservationContext(lggr, r, telem)
	opts := llov30.DSOpts(nil)

	streamID := streams.StreamID(1)
	val := decimal.NewFromFloat(123.456)

	// observes the same pipeline 1000 times to try and detect races etc
	r.pipelines = make(map[streams.StreamID]*mockPipeline)
	r.pipelines[streamID] = makePipelineWithSingleResult[decimal.Decimal](0, val, nil)
	g, ctx := errgroup.WithContext(ctx)
	for range 1000 {
		g.Go(func() error {
			_, err := oc.Observe(ctx, streamID, opts)
			return err
		})
	}
	if err := g.Wait(); err != nil {
		t.Fatalf("Observation failed: %v", err)
	}
}

type mockPipelineConfig struct{}

func (m *mockPipelineConfig) DefaultHTTPLimit() int64 { return 10000 }
func (m *mockPipelineConfig) DefaultHTTPTimeout() commonconfig.Duration {
	return *commonconfig.MustNewDuration(1 * time.Hour)
}
func (m *mockPipelineConfig) MaxRunDuration() time.Duration  { return 1 * time.Hour }
func (m *mockPipelineConfig) ReaperInterval() time.Duration  { return 0 }
func (m *mockPipelineConfig) ReaperThreshold() time.Duration { return 0 }

// func (m *mockPipelineConfig) VerboseLogging() bool           { return true }
func (m *mockPipelineConfig) VerboseLogging() bool { return false }

type mockBridgeConfig struct{}

func (m *mockBridgeConfig) BridgeResponseURL() *url.URL {
	return nil
}

func (m *mockBridgeConfig) BridgeCacheTTL() time.Duration {
	return 0
}

func createBridge(t testing.TB, name, val string, borm bridges.ORM, maxCalls int64) {
	callcount := atomic.NewInt64(0)
	bridge := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		n := callcount.Inc()
		if maxCalls > 0 && n > maxCalls {
			panic("too many calls to bridge" + name)
		}
		_, herr := io.ReadAll(req.Body)
		if herr != nil {
			panic(herr)
		}

		res.WriteHeader(http.StatusOK)
		resp := fmt.Sprintf(`{"result": %s}`, val)
		_, herr = res.Write([]byte(resp))
		if herr != nil {
			panic(herr)
		}
	}))
	t.Cleanup(bridge.Close)
	u, _ := url.Parse(bridge.URL)
	require.NoError(t, borm.CreateBridgeType(t.Context(), &bridges.BridgeType{
		Name: bridges.BridgeName(name),
		URL:  models.WebURL(*u),
	}))
}

func TestObservationContext_Observe_integrationRealPipeline(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	bridgesORM := bridges.NewORM(db)

	createBridge(t, "foo-bridge", `123.456`, bridgesORM, 1)
	createBridge(t, "bar-bridge", `"124.456"`, bridgesORM, 1)

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

	jobStreamID := streams.StreamID(5)

	t.Run("using only streamID attributes", func(t *testing.T) {
		t.Parallel()
		jb := job.Job{
			Type:     job.Stream,
			StreamID: &jobStreamID,
			PipelineSpec: &pipeline.Spec{
				DotDagSource: `
// Benchmark Price
result1          [type=memo value="123.9"];
multiply2 	  	 [type=multiply times=1 streamID=1 index=0]; // force conversion to decimal

result2          [type=bridge name="foo-bridge" requestData="{\"data\":{\"data\":\"foo\"}}"];
result2_parse    [type=jsonparse path="result" streamID=2 index=1];

result3          [type=bridge name="bar-bridge" requestData="{\"data\":{\"data\":\"bar\"}}"];
result3_parse    [type=jsonparse path="result"];
multiply3 	  	 [type=multiply times=1 streamID=3 index=2]; // force conversion to decimal

result1 -> multiply2;
result2 -> result2_parse;
result3 -> result3_parse -> multiply3; 
`,
			},
		}
		err := r.Register(jb, nil)
		require.NoError(t, err)

		telem := &mockTelemeter{}
		oc := newObservationContext(lggr, r, telem)
		opts := llov30.DSOpts(nil)

		val, err := oc.Observe(ctx, streams.StreamID(1), opts)
		require.NoError(t, err)
		assert.Equal(t, "123.9", val.(*lloprotocol.Decimal).String())
		val, err = oc.Observe(ctx, streams.StreamID(2), opts)
		require.NoError(t, err)
		assert.Equal(t, "123.456", val.(*lloprotocol.Decimal).String())
		val, err = oc.Observe(ctx, streams.StreamID(3), opts)
		require.NoError(t, err)
		assert.Equal(t, "124.456", val.(*lloprotocol.Decimal).String())

		val, err = oc.Observe(ctx, jobStreamID, opts)
		require.NoError(t, err)
		assert.Equal(t, &lloprotocol.Quote{
			Bid:       decimal.NewFromFloat32(123.456),
			Benchmark: decimal.NewFromFloat32(123.9),
			Ask:       decimal.NewFromFloat32(124.456),
		}, val.(*lloprotocol.Quote))
	})

	t.Run("an invalid job-level quote fails the tagged streams of the same pipeline", func(t *testing.T) {
		t.Parallel()
		badJobStreamID := streams.StreamID(15)
		jb := job.Job{
			Type:     job.Stream,
			StreamID: &badJobStreamID,
			PipelineSpec: &pipeline.Spec{
				// Benchmark deliberately above the ask. No bridges here: the shared
				// ones only serve a single request.
				DotDagSource: `
result1          [type=memo value="900.0022"];
multiply1 	  	 [type=multiply times=1 streamID=11 index=0];

result2          [type=memo value="123.456"];
multiply2 	  	 [type=multiply times=1 streamID=12 index=1];

result3          [type=memo value="124.456"];
multiply3 	  	 [type=multiply times=1 streamID=13 index=2];

result1 -> multiply1;
result2 -> multiply2;
result3 -> multiply3;
`,
			},
		}
		badReg := streams.NewRegistry(lggr, runner)
		require.NoError(t, badReg.Register(jb, nil))

		oc := newObservationContext(lggr, badReg, &mockTelemeter{})
		opts := llov30.DSOpts(nil)

		// The quote is only assembled for the job-level stream, but the whole set
		// fails: the tagged streams come from the same pipeline run.
		for _, sid := range []streams.StreamID{11, 12, 13, badJobStreamID} {
			val, err := oc.Observe(ctx, sid, opts)
			assert.Nil(t, val, "stream %d should have no value", sid)
			var qerr QuoteInvariantError
			require.ErrorAs(t, err, &qerr, "stream %d", sid)
			assert.Equal(t, badJobStreamID, qerr.StreamID)
		}
	})
}

func TestObservationContext_Observe_concurrentAtomicOutput(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	const n = 20

	reg := &mockRegistry{pipelines: make(map[streams.StreamID]*mockPipeline)}
	pipelines := make([]*mockPipeline, n)

	for i := range n {
		ui := uint32(i)
		sid1 := ui*3 + 1
		sid2 := ui*3 + 2
		sid3 := ui*3 + 3
		val1 := decimal.NewFromInt(int64(i*10 + 1))
		val2 := decimal.NewFromInt(int64(i*10 + 2))
		val3 := decimal.NewFromInt(int64(i*10 + 3))

		mp := makePipelineWithMultipleStreamResults(
			[]streams.StreamID{sid1, sid2, sid3},
			[]any{val1, val2, val3},
		)
		pipelines[i] = mp
		reg.pipelines[sid1] = mp
		reg.pipelines[sid2] = mp
		reg.pipelines[sid3] = mp
	}

	lggr := logger.TestLogger(t)
	telem := &mockTelemeter{}
	oc := newObservationContext(lggr, reg, telem)
	opts := llov30.DSOpts(nil)

	type result struct {
		strmID uint32
		val    lloprotocol.StreamValue
		err    error
	}

	pipelineGroupResults := make([][3]result, n)
	var wg sync.WaitGroup

	for i := range n {
		ui := uint32(i)
		sid1 := ui*3 + 1
		sid2 := ui*3 + 2
		sid3 := ui*3 + 3
		for j, strmID := range [3]streams.StreamID{sid1, sid2, sid3} {
			wg.Go(func() {
				val, err := oc.Observe(ctx, strmID, opts)
				pipelineGroupResults[i][j] = result{strmID, val, err}
			})
		}
	}
	wg.Wait()

	for i, group := range pipelineGroupResults {
		for _, r := range group {
			require.NoError(t, r.err, "pipeline %d, stream %d", i, r.strmID)
			require.NotNil(t, r.val, "pipeline %d, stream %d: nil value", i, r.strmID)
		}
		assert.Equal(t, strconv.Itoa(i*10+1), group[0].val.(*lloprotocol.Decimal).String(), "pipeline %d sid1", i)
		assert.Equal(t, strconv.Itoa(i*10+2), group[1].val.(*lloprotocol.Decimal).String(), "pipeline %d sid2", i)
		assert.Equal(t, strconv.Itoa(i*10+3), group[2].val.(*lloprotocol.Decimal).String(), "pipeline %d sid3", i)
		assert.Equal(t, int32(1), pipelines[i].runCount.Load(), "pipeline %d should have run exactly once", i)
	}
}

func BenchmarkObservationContext_Observe_integrationRealPipeline_concurrencyStressTest_manyStreams(b *testing.B) {
	ctx := b.Context()
	lggr := logger.TestLogger(b)
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

	telem := &mockTelemeter{}
	oc := newObservationContext(lggr, r, telem)
	opts := llov30.DSOpts(nil)

	// concurrency stress test
	b.ResetTimer()
	g, ctx := errgroup.WithContext(ctx)
	for i := range n {
		for _, strmID := range []uint32{i, i + n, i + 2*n, i + 3*n} {
			g.Go(func() error {
				// ignore errors, only care about races
				oc.Observe(ctx, strmID, opts) //nolint:errcheck // ignore error
				return nil
			})
		}
	}
	if err := g.Wait(); err != nil {
		b.Fatalf("Observation failed: %v", err)
	}
}

// quoteResultMap builds the map form of a tagged Quote stream result.
// benchmarkKey lets a test exercise the "mid" alias.
func quoteResultMap(bid, benchmarkKey, benchmark, ask string) map[string]any {
	m := map[string]any{
		"streamValueType": int64(lloprotocol.LLOStreamValue_Quote),
		"bid":             bid,
		"ask":             ask,
	}
	if benchmarkKey != "" {
		m[benchmarkKey] = benchmark
	}
	return m
}

func makeLegacyQuotePipeline(benchmark, bid, ask string) *mockPipeline {
	// Untagged terminal tasks, ordered Benchmark, Bid, Ask.
	trrs := make(pipeline.TaskRunResults, 0, 3)
	for _, v := range []string{benchmark, bid, ask} {
		trrs = append(trrs, pipeline.TaskRunResult{Task: &pipeline.MemoTask{}, Result: pipeline.Result{Value: v}})
	}
	return &mockPipeline{run: &pipeline.Run{}, trrs: trrs}
}

func TestObservationContext_Observe_quotes(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	opts := llov30.DSOpts(nil)

	newOC := func(t *testing.T, pipelines map[streams.StreamID]*mockPipeline) ObservationContext {
		r := &mockRegistry{pipelines: pipelines}
		return newObservationContext(logger.TestLogger(t), r, &mockTelemeter{})
	}

	t.Run("parses a tagged quote", func(t *testing.T) {
		t.Parallel()
		sid := streams.StreamID(1)
		p := makePipelineWithMultipleStreamResults([]streams.StreamID{sid}, []any{quoteResultMap("1.1", "benchmark", "2.2", "3.3")})
		val, err := newOC(t, map[streams.StreamID]*mockPipeline{sid: p}).Observe(ctx, sid, opts)
		require.NoError(t, err)
		q, ok := val.(*lloprotocol.Quote)
		require.True(t, ok, "expected *lloprotocol.Quote, got %T", val)
		assert.Equal(t, "1.1", q.Bid.String())
		assert.Equal(t, "2.2", q.Benchmark.String())
		assert.Equal(t, "3.3", q.Ask.String())
	})

	t.Run("accepts mid as an alias for benchmark", func(t *testing.T) {
		t.Parallel()
		sid := streams.StreamID(1)
		p := makePipelineWithMultipleStreamResults([]streams.StreamID{sid}, []any{quoteResultMap("1.1", "mid", "2.2", "3.3")})
		val, err := newOC(t, map[streams.StreamID]*mockPipeline{sid: p}).Observe(ctx, sid, opts)
		require.NoError(t, err)
		assert.Equal(t, "2.2", val.(*lloprotocol.Quote).Benchmark.String())
	})

	t.Run("rejects both benchmark and mid", func(t *testing.T) {
		t.Parallel()
		sid := streams.StreamID(1)
		m := quoteResultMap("1.1", "benchmark", "2.2", "3.3")
		m["mid"] = "2.2"
		p := makePipelineWithMultipleStreamResults([]streams.StreamID{sid}, []any{m})
		_, err := newOC(t, map[streams.StreamID]*mockPipeline{sid: p}).Observe(ctx, sid, opts)
		require.ErrorContains(t, err, "expected exactly one of 'benchmark' or 'mid', got both")
	})

	t.Run("errors on a malformed tagged quote", func(t *testing.T) {
		t.Parallel()
		for name, m := range map[string]map[string]any{
			"missing benchmark": quoteResultMap("1.1", "", "", "3.3"),
			"missing ask":       {"streamValueType": int64(lloprotocol.LLOStreamValue_Quote), "bid": "1.1", "benchmark": "2.2"},
			"unparseable bid":   quoteResultMap("not-a-number", "benchmark", "2.2", "3.3"),
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				sid := streams.StreamID(1)
				p := makePipelineWithMultipleStreamResults([]streams.StreamID{sid}, []any{m})
				_, err := newOC(t, map[streams.StreamID]*mockPipeline{sid: p}).Observe(ctx, sid, opts)
				require.ErrorContains(t, err, "failed to parse Quote")
			})
		}
	})

	t.Run("fails a tagged quote violating the bid/benchmark/ask invariant", func(t *testing.T) {
		t.Parallel()
		for name, m := range map[string]map[string]any{
			"bid above benchmark": quoteResultMap("3.3", "benchmark", "2.2", "4.4"),
			"benchmark above ask": quoteResultMap("1.1", "benchmark", "5.5", "4.4"),
			"bid above ask":       quoteResultMap("5.5", "mid", "4.4", "3.3"),
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				sid := streams.StreamID(1)
				p := makePipelineWithMultipleStreamResults([]streams.StreamID{sid}, []any{m})
				_, err := newOC(t, map[streams.StreamID]*mockPipeline{sid: p}).Observe(ctx, sid, opts)
				var qerr QuoteInvariantError
				require.ErrorAs(t, err, &qerr)
				assert.Equal(t, sid, qerr.StreamID)
				assert.Contains(t, err.Error(), "quote invariant violation for stream 1")
			})
		}
	})

	t.Run("accepts a quote with equal bid, benchmark and ask", func(t *testing.T) {
		t.Parallel()
		sid := streams.StreamID(1)
		p := makePipelineWithMultipleStreamResults([]streams.StreamID{sid}, []any{quoteResultMap("2.2", "benchmark", "2.2", "2.2")})
		_, err := newOC(t, map[streams.StreamID]*mockPipeline{sid: p}).Observe(ctx, sid, opts)
		require.NoError(t, err)
	})

	t.Run("fails every stream of a pipeline when one tagged quote is invalid", func(t *testing.T) {
		t.Parallel()
		plain, badQuote, goodQuote := streams.StreamID(4), streams.StreamID(5), streams.StreamID(6)
		p := makePipelineWithMultipleStreamResults(
			[]streams.StreamID{plain, badQuote, goodQuote},
			[]any{
				decimal.NewFromFloat(12.34),
				quoteResultMap("9.9", "benchmark", "2.2", "3.3"), // bid > benchmark
				quoteResultMap("1.1", "benchmark", "2.2", "3.3"),
			},
		)
		oc := newOC(t, map[streams.StreamID]*mockPipeline{plain: p, badQuote: p, goodQuote: p})

		for _, sid := range []streams.StreamID{plain, badQuote, goodQuote} {
			val, err := oc.Observe(ctx, sid, opts)
			assert.Nil(t, val, "stream %d should have no value", sid)
			var qerr QuoteInvariantError
			require.ErrorAs(t, err, &qerr, "stream %d should fail with QuoteInvariantError", sid)
			// The violation is always attributed to the offending stream, not the
			// stream that was asked for.
			assert.Equal(t, badQuote, qerr.StreamID)
		}
		assert.Equal(t, int32(1), p.runCount.Load(), "pipeline should only run once")
	})

	t.Run("fails every stream concurrently observing a pipeline with an invalid quote", func(t *testing.T) {
		t.Parallel()
		sids := []streams.StreamID{1, 2, 3}
		p := makePipelineWithMultipleStreamResults(sids, []any{
			decimal.NewFromFloat(1),
			decimal.NewFromFloat(2),
			quoteResultMap("9.9", "benchmark", "2.2", "3.3"),
		})
		pipelines := map[streams.StreamID]*mockPipeline{}
		for _, sid := range sids {
			pipelines[sid] = p
		}
		oc := newOC(t, pipelines)

		var mu sync.Mutex
		errsBySID := map[streams.StreamID]error{}
		var wg sync.WaitGroup
		for _, sid := range sids {
			for range 10 {
				wg.Go(func() {
					_, err := oc.Observe(ctx, sid, opts)
					mu.Lock()
					defer mu.Unlock()
					errsBySID[sid] = err
				})
			}
		}
		wg.Wait()

		require.Len(t, errsBySID, len(sids))
		for sid, err := range errsBySID {
			var qerr QuoteInvariantError
			require.ErrorAs(t, err, &qerr, "stream %d", sid)
			assert.Equal(t, streams.StreamID(3), qerr.StreamID)
		}
		assert.Equal(t, int32(1), p.runCount.Load())
	})

	t.Run("fails a legacy three-terminal quote violating the invariant", func(t *testing.T) {
		t.Parallel()
		sid := streams.StreamID(1)
		// Ordering is Benchmark, Bid, Ask.
		p := makeLegacyQuotePipeline("2.2", "9.9", "3.3")
		val, err := newOC(t, map[streams.StreamID]*mockPipeline{sid: p}).Observe(ctx, sid, opts)
		assert.Nil(t, val)
		var qerr QuoteInvariantError
		require.ErrorAs(t, err, &qerr)
		assert.Equal(t, sid, qerr.StreamID)
	})

	t.Run("returns a valid legacy three-terminal quote", func(t *testing.T) {
		t.Parallel()
		sid := streams.StreamID(1)
		p := makeLegacyQuotePipeline("2.2", "1.1", "3.3")
		val, err := newOC(t, map[streams.StreamID]*mockPipeline{sid: p}).Observe(ctx, sid, opts)
		require.NoError(t, err)
		q := val.(*lloprotocol.Quote)
		assert.Equal(t, "1.1", q.Bid.String())
		assert.Equal(t, "2.2", q.Benchmark.String())
		assert.Equal(t, "3.3", q.Ask.String())
	})
}
