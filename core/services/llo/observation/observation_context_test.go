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
	"github.com/smartcontractkit/chainlink-data-streams/llo"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	clhttptest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/httptest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	clnull "github.com/smartcontractkit/chainlink/v2/core/null"
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

func TestObservationContext_Observe(t *testing.T) {
	ctx := t.Context()
	r := &mockRegistry{}
	telem := &mockTelemeter{}
	lggr := logger.TestLogger(t)
	oc := newObservationContext(lggr, r, telem)
	opts := llo.DSOpts(nil)

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
		streamID1:  &mockPipeline{},
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

	t.Run("returns error in case of missing pipeline", func(t *testing.T) {
		_, err := oc.Observe(ctx, missingStreamID, opts)
		require.EqualError(t, err, "no pipeline for stream: 0")
	})
	t.Run("returns error in case of zero results", func(t *testing.T) {
		_, err := oc.Observe(ctx, streamID1, opts)
		require.EqualError(t, err, "invalid number of results, expected: 1 or 3, got: 0")
	})
	t.Run("returns composite value from legacy job with single top-level streamID", func(t *testing.T) {
		val, err := oc.Observe(ctx, streamID2, opts)
		require.NoError(t, err)

		assert.Equal(t, "12.34", val.(*llo.Decimal).String())
	})
	t.Run("returns error in case of erroring pipeline", func(t *testing.T) {
		_, err := oc.Observe(ctx, streamID3, opts)
		require.EqualError(t, err, "pipeline error")
	})
	t.Run("returns values for multiple stream IDs within the same job based on streamID tag with a single pipeline execution", func(t *testing.T) {
		val, err := oc.Observe(ctx, streamID4, opts)
		require.NoError(t, err)
		assert.Equal(t, "12.34", val.(*llo.Decimal).String())

		val, err = oc.Observe(ctx, streamID5, opts)
		require.NoError(t, err)
		assert.Equal(t, "56.78", val.(*llo.Decimal).String())

		val, err = oc.Observe(ctx, streamID6, opts)
		require.NoError(t, err)
		assert.Equal(t, "90.12", val.(*llo.Decimal).String())

		assert.Equal(t, int32(1), multiPipelineDecimal.runCount.Load())

		// returns cached values on subsequent calls
		val, err = oc.Observe(ctx, streamID6, opts)
		require.NoError(t, err)
		assert.Equal(t, "90.12", val.(*llo.Decimal).String())

		assert.Equal(t, int32(1), multiPipelineDecimal.runCount.Load())
	})
	t.Run("returns value from float64 value", func(t *testing.T) {
		val, err := oc.Observe(ctx, streamID7, opts)
		require.NoError(t, err)

		assert.Equal(t, "1.23", val.(*llo.Decimal).String())
	})
	t.Run("returns value from int64 value", func(t *testing.T) {
		val, err := oc.Observe(ctx, streamID8, opts)
		require.NoError(t, err)

		assert.Equal(t, "5", val.(*llo.Decimal).String())
	})
	t.Run("partial extraction failure in multi-stream pipeline", func(t *testing.T) {
		val, err := oc.Observe(ctx, streamID9, opts)
		require.NoError(t, err)
		assert.Equal(t, "100.5", val.(*llo.Decimal).String())

		_, err = oc.Observe(ctx, streamID10, opts)
		require.Error(t, err, "unparseable value should fail extraction")

		val, err = oc.Observe(ctx, streamID11, opts)
		require.NoError(t, err)
		assert.Equal(t, "200.5", val.(*llo.Decimal).String())

		assert.Equal(t, int32(1), multiPipelinePartialFail.runCount.Load())
	})
}

func TestObservationContext_Observe_concurrencyStressTest(t *testing.T) {
	ctx := t.Context()
	r := &mockRegistry{}
	telem := &mockTelemeter{}
	lggr := logger.TestLogger(t)
	oc := newObservationContext(lggr, r, telem)
	opts := llo.DSOpts(nil)

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

func createBridge(t testing.TB, name string, val string, borm bridges.ORM, maxCalls int64) {
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
		jb := job.Job{
			Type:     job.Stream,
			StreamID: &jobStreamID,
			PipelineSpec: &pipeline.Spec{
				DotDagSource: `
// Benchmark Price
result1          [type=memo value="900.0022"];
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
		opts := llo.DSOpts(nil)

		val, err := oc.Observe(ctx, streams.StreamID(1), opts)
		require.NoError(t, err)
		assert.Equal(t, "900.0022", val.(*llo.Decimal).String())
		val, err = oc.Observe(ctx, streams.StreamID(2), opts)
		require.NoError(t, err)
		assert.Equal(t, "123.456", val.(*llo.Decimal).String())
		val, err = oc.Observe(ctx, streams.StreamID(3), opts)
		require.NoError(t, err)
		assert.Equal(t, "124.456", val.(*llo.Decimal).String())

		val, err = oc.Observe(ctx, jobStreamID, opts)
		require.NoError(t, err)
		assert.Equal(t, &llo.Quote{
			Bid:       decimal.NewFromFloat32(123.456),
			Benchmark: decimal.NewFromFloat32(900.0022),
			Ask:       decimal.NewFromFloat32(124.456),
		}, val.(*llo.Quote))
	})
}

func TestObservationContext_Observe_concurrentAtomicOutput(t *testing.T) {
	ctx := t.Context()
	const n = 20

	reg := &mockRegistry{pipelines: make(map[streams.StreamID]*mockPipeline)}
	pipelines := make([]*mockPipeline, n)

	for i := range n {
		ui := uint32(i) //nolint:gosec // i bounded by n=20
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
	opts := llo.DSOpts(nil)

	type result struct {
		strmID uint32
		val    llo.StreamValue
		err    error
	}

	pipelineGroupResults := make([][3]result, n)
	var wg sync.WaitGroup

	for i := range n {
		ui := uint32(i) //nolint:gosec // i bounded by n=20
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
		assert.Equal(t, strconv.Itoa(i*10+1), group[0].val.(*llo.Decimal).String(), "pipeline %d sid1", i)
		assert.Equal(t, strconv.Itoa(i*10+2), group[1].val.(*llo.Decimal).String(), "pipeline %d sid2", i)
		assert.Equal(t, strconv.Itoa(i*10+3), group[2].val.(*llo.Decimal).String(), "pipeline %d sid3", i)
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

	telem := &mockTelemeter{}
	oc := newObservationContext(lggr, r, telem)
	opts := llo.DSOpts(nil)

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
