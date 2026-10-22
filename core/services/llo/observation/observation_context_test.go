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
	"slices"
	"strconv"
	"strings"
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

	t.Run("an invalid job-level quote is not validated and does not fail the tagged streams", func(t *testing.T) {
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

		// Only tagged values are validated, so the tagged streams succeed.
		for _, sid := range []streams.StreamID{11, 12, 13} {
			val, err := oc.Observe(ctx, sid, opts)
			require.NoError(t, err, "stream %d", sid)
			assert.IsType(t, &lloprotocol.Decimal{}, val, "stream %d", sid)
		}

		// The job-level quote is assembled from the terminal outputs and returned
		// as is.
		val, err := oc.Observe(ctx, badJobStreamID, opts)
		require.NoError(t, err)
		q, ok := val.(*lloprotocol.Quote)
		require.True(t, ok, "expected *lloprotocol.Quote, got %T", val)
		assert.False(t, q.IsValid())
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

// makeNamedQuotePipeline returns a pipeline for jobStreamID whose terminal
// outputs, at index 0, 1 and 2, are multiply tasks with the given DOT IDs and
// values. If tags is given, the terminal outputs are tagged with those stream
// IDs in the same order.
func makeNamedQuotePipeline(t *testing.T, jobStreamID streams.StreamID, tags []streams.StreamID, names [3]string, v0, v1, v2 string) *mockPipeline {
	t.Helper()
	return makeNamedTerminalsPipeline(t, jobStreamID, tags, names[:], []string{v0, v1, v2})
}

// makeNamedTerminalsPipeline returns a pipeline for jobStreamID whose terminal
// outputs, in index order, are multiply tasks with the given DOT IDs and values,
// each over a jsonparse. If tags is given, the terminal outputs are tagged with
// those stream IDs in the same order.
func makeNamedTerminalsPipeline(t *testing.T, jobStreamID streams.StreamID, tags []streams.StreamID, names, values []string) *mockPipeline {
	t.Helper()
	var dag strings.Builder
	dag.WriteString("ds [type=memo value=\"{}\"];\n")
	for i, name := range names {
		tag := ""
		if tags != nil {
			tag = fmt.Sprintf(" streamID=%d", tags[i])
		}
		fmt.Fprintf(&dag, "parse_%d [type=jsonparse path=\"data,%d\"];\n", i, i)
		fmt.Fprintf(&dag, "%s [type=multiply times=1 index=%d%s];\n", name, i, tag)
		fmt.Fprintf(&dag, "ds -> parse_%d -> %s;\n", i, name)
	}
	p, err := pipeline.Parse(dag.String())
	require.NoError(t, err)
	var trrs pipeline.TaskRunResults
	for i, name := range names {
		task := p.ByDotID(name)
		require.NotNil(t, task, name)
		trrs = append(trrs, pipeline.TaskRunResult{Task: task, Result: pipeline.Result{Value: decimal.RequireFromString(values[i])}})
	}
	streamIDs := append(slices.Clone(tags), jobStreamID)
	return &mockPipeline{run: &pipeline.Run{}, trrs: trrs, streamIDs: streamIDs, terminalDotIDs: streams.TerminalDotIDs(p)}
}

// makeMedianQuotePipeline returns a pipeline for jobStreamID shaped like the
// feeds-manager quote template: one bridge per entry of sourcePaths, each
// parsed into Benchmark, Bid and Ask with that entry's paths, feeding three
// terminal median tasks. If tags is given, the medians are tagged with those
// stream IDs in Benchmark, Bid, Ask order.
func makeMedianQuotePipeline(t *testing.T, jobStreamID streams.StreamID, tags []streams.StreamID, sourcePaths [][3]string, benchmark, bid, ask string) *mockPipeline {
	t.Helper()
	fields := [3]string{"benchmark", "bid", "ask"}
	var dag strings.Builder
	for i, paths := range sourcePaths {
		fmt.Fprintf(&dag, "ds%d_payload [type=memo value=\"{}\"];\n", i)
		for j, field := range fields {
			fmt.Fprintf(&dag, "ds%d_%s [type=jsonparse path=\"%s\"];\n", i, field, paths[j])
			fmt.Fprintf(&dag, "ds%d_payload -> ds%d_%s -> %s_price;\n", i, i, field, field)
		}
	}
	for j, field := range fields {
		tag := ""
		if tags != nil {
			tag = fmt.Sprintf(" streamID=%d", tags[j])
		}
		fmt.Fprintf(&dag, "%s_price [type=median allowedFaults=1 index=%d%s];\n", field, j, tag)
	}
	p, err := pipeline.Parse(dag.String())
	require.NoError(t, err)
	values := [3]string{benchmark, bid, ask}
	var trrs pipeline.TaskRunResults
	for j, field := range fields {
		task := p.ByDotID(field + "_price")
		require.NotNil(t, task, field)
		trrs = append(trrs, pipeline.TaskRunResult{Task: task, Result: pipeline.Result{Value: decimal.RequireFromString(values[j])}})
	}
	streamIDs := append(slices.Clone(tags), jobStreamID)
	return &mockPipeline{run: &pipeline.Run{}, trrs: trrs, streamIDs: streamIDs, terminalDotIDs: streams.TerminalDotIDs(p)}
}

// makePerpsPriceQuotePipeline returns a pipeline for jobStreamID shaped like
// the feeds-manager perps price template: two venues sampled per field,
// optionally normalized by USDC/USDT conversion factors, decayed by staleness,
// then aggregated into wm_mid, wm_bid and wm_ask, tagged with tags in that
// order. bid and ask are anchored against the winsorized mid.
func makePerpsPriceQuotePipeline(t *testing.T, jobStreamID streams.StreamID, tags [3]streams.StreamID, normalize bool, benchmark, bid, ask string) *mockPipeline {
	t.Helper()
	fields := [3]string{"mid", "bid", "ask"}
	venues := [2]string{"hyperliquid", "binance"}
	var dag strings.Builder
	if normalize {
		for _, f := range [2]string{"usdc", "usdt"} {
			fmt.Fprintf(&dag, "fetch_%s [type=bridge name=\"bridge-data-engine\" requestData=<{\"data\":{\"endpoint\":\"crypto-v3\"}}> cacheTTL=\"300s\"]\n", f)
			fmt.Fprintf(&dag, "%s_factor [type=sample input=\"$(fetch_%s)\" source=\"%s\" weight=\"1\" valuePath=\"data,price\" tsPath=\"timestamps,providerDataReceivedUnixMs\" decimals=\"18\"]\n", f, f, strings.ToUpper(f))
		}
	}
	for _, v := range venues {
		fmt.Fprintf(&dag, "fetch_%s [type=bridge name=\"bridge-coinpaprika\" requestData=<{\"data\":{\"endpoint\":\"markprice\"}}> cacheTTL=\"300s\"]\n", v)
		for _, f := range fields {
			fmt.Fprintf(&dag, "%s_%s [type=sample input=\"$(fetch_%s)\" source=\"%s\" weight=\"1\" unit=\"USDT\" valuePath=\"data,%s\" tsPath=\"timestamps,providerDataReceivedUnixMs\"]\n", v, f, v, v, f)
		}
	}
	for _, f := range fields {
		samples := fmt.Sprintf("[$(hyperliquid_%s), $(binance_%s)]", f, f)
		if normalize {
			fmt.Fprintf(&dag, "norm_%s [type=normalize samples=\"%s\" unitMap=<{\"USDC\":$(usdc_factor),\"USDT\":$(usdt_factor)}> targetUnit=\"USD\" enabled=\"true\" onMissingRate=\"error\"]\n", f, samples)
			samples = fmt.Sprintf("$(norm_%s)", f)
		}
		fmt.Fprintf(&dag, "stale_%s [type=staleness samples=\"%s\" method=\"exp_cooldown\" threshold=\"60s\" halfLife=\"5s\" decayThreshold=\"0.03\" cutoff=\"300s\"]\n", f, samples)
	}
	fmt.Fprintf(&dag, "winsor_mid [type=weightedmean samples=\"$(stale_mid)\" method=\"winsor_samples\" band=\"0.03\"]\n")
	fmt.Fprintf(&dag, "wm_mid [type=weightedmean samples=\"$(winsor_mid)\" minWeightMass=\"1\" streamID=%d index=0]\n", tags[0])
	fmt.Fprintf(&dag, "wm_bid [type=anchor reference=\"$(winsor_mid)\" low=\"$(stale_bid)\" high=\"$(stale_ask)\" select=\"low\" minWeightMass=\"1\" streamID=%d index=1]\n", tags[1])
	fmt.Fprintf(&dag, "wm_ask [type=anchor reference=\"$(winsor_mid)\" low=\"$(stale_bid)\" high=\"$(stale_ask)\" select=\"high\" minWeightMass=\"1\" streamID=%d index=2]\n", tags[2])
	if normalize {
		for _, fa := range [2]string{"usdc", "usdt"} {
			for _, f := range fields {
				fmt.Fprintf(&dag, "fetch_%s -> %s_factor -> norm_%s\n", fa, fa, f)
			}
		}
	}
	for _, v := range venues {
		for _, f := range fields {
			if normalize {
				fmt.Fprintf(&dag, "fetch_%s -> %s_%s -> norm_%s -> stale_%s -> wm_%s\n", v, v, f, f, f, f)
			} else {
				fmt.Fprintf(&dag, "fetch_%s -> %s_%s -> stale_%s -> wm_%s\n", v, v, f, f, f)
			}
		}
	}
	p, err := pipeline.Parse(dag.String())
	require.NoError(t, err)
	values := [3]string{benchmark, bid, ask}
	var trrs pipeline.TaskRunResults
	for i, f := range fields {
		task := p.ByDotID("wm_" + f)
		require.NotNil(t, task, f)
		trrs = append(trrs, pipeline.TaskRunResult{Task: task, Result: pipeline.Result{Value: decimal.RequireFromString(values[i])}})
	}
	require.Len(t, trrs.Terminals(), 3)
	return &mockPipeline{run: &pipeline.Run{}, trrs: trrs, streamIDs: append(tags[:], jobStreamID), terminalDotIDs: streams.TerminalDotIDs(p)}
}

// makeDexCrossQuotePipeline returns a pipeline for jobStreamID shaped like the
// feeds-manager dex cross exchange template: two market rate sources parsed
// into Benchmark, Bid and Ask, aggregated by median, then converted by an
// exchange rate referenced through multiply times.
func makeDexCrossQuotePipeline(t *testing.T, jobStreamID streams.StreamID, benchmark, bid, ask string) *mockPipeline {
	t.Helper()
	p, err := pipeline.Parse(`
exchange_rate [type=bridge name="bridge-x" requestData="{}"];
exchange_rate_parse [type=jsonparse path="result"];
exchange_rate_value [type=multiply input="$(exchange_rate_parse)" times=1];
exchange_rate -> exchange_rate_parse -> exchange_rate_value;

ds1_market_rate_payload [type=bridge name="bridge-y" requestData="{}"];
ds1_market_rate_benchmark [type=jsonparse path="data,mid"];
ds1_market_rate_bid [type=jsonparse path="data,bid"];
ds1_market_rate_ask [type=jsonparse path="data,ask"];
ds1_market_rate_payload -> ds1_market_rate_benchmark -> market_rate_benchmark_median;
ds1_market_rate_payload -> ds1_market_rate_bid -> market_rate_bid_median;
ds1_market_rate_payload -> ds1_market_rate_ask -> market_rate_ask_median;

ds2_market_rate_payload [type=bridge name="bridge-z" requestData="{}"];
ds2_market_rate_benchmark [type=jsonparse path="data,mid"];
ds2_market_rate_bid [type=jsonparse path="data,bid"];
ds2_market_rate_ask [type=jsonparse path="data,ask"];
ds2_market_rate_payload -> ds2_market_rate_benchmark -> market_rate_benchmark_median;
ds2_market_rate_payload -> ds2_market_rate_bid -> market_rate_bid_median;
ds2_market_rate_payload -> ds2_market_rate_ask -> market_rate_ask_median;

market_rate_benchmark_median [type=median allowedFaults=1];
market_rate_bid_median [type=median allowedFaults=1];
market_rate_ask_median [type=median allowedFaults=1];

exchange_rate_multiply_benchmark [type=multiply times="$(exchange_rate_value)"];
exchange_rate_multiply_bid [type=multiply times="$(exchange_rate_value)"];
exchange_rate_multiply_ask [type=multiply times="$(exchange_rate_value)"];

benchmark_price [type=any index=0];
market_rate_benchmark_median -> exchange_rate_multiply_benchmark -> benchmark_price;
bid_price [type=any index=1];
market_rate_bid_median -> exchange_rate_multiply_bid -> bid_price;
ask_price [type=any index=2];
market_rate_ask_median -> exchange_rate_multiply_ask -> ask_price;
`)
	require.NoError(t, err)
	values := [3]string{benchmark, bid, ask}
	var trrs pipeline.TaskRunResults
	for i, id := range [3]string{"benchmark_price", "bid_price", "ask_price"} {
		task := p.ByDotID(id)
		require.NotNil(t, task, id)
		trrs = append(trrs, pipeline.TaskRunResult{Task: task, Result: pipeline.Result{Value: decimal.RequireFromString(values[i])}})
	}
	require.Len(t, trrs.Terminals(), 3)
	return &mockPipeline{run: &pipeline.Run{}, trrs: trrs, streamIDs: []streams.StreamID{jobStreamID}, terminalDotIDs: streams.TerminalDotIDs(p)}
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

	t.Run("does not validate a legacy three-terminal quote", func(t *testing.T) {
		t.Parallel()
		sid := streams.StreamID(1)
		// Ordering is Benchmark, Bid, Ask.
		p := makeLegacyQuotePipeline("2.2", "9.9", "3.3")
		val, err := newOC(t, map[streams.StreamID]*mockPipeline{sid: p}).Observe(ctx, sid, opts)
		require.NoError(t, err)
		q := val.(*lloprotocol.Quote)
		assert.Equal(t, "9.9", q.Bid.String())
		assert.False(t, q.IsValid())
	})

	t.Run("fails a terminal quote with quote named tasks violating the invariant", func(t *testing.T) {
		t.Parallel()
		for name, names := range map[string][3]string{
			"price suffix":  {"benchmark_price", "bid_price", "ask_price"},
			"mid alias":     {"wm_mid", "wm_bid", "wm_ask"},
			"bare names":    {"benchmark", "bid", "ask"},
			"mixed case":    {"Benchmark_Price", "BID", "Ask_Price"},
			"source prefix": {"ds1_benchmark", "ds1_bid", "ds1_ask"},
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				sid := streams.StreamID(1)
				p := makeNamedQuotePipeline(t, sid, nil, names, "2.2", "9.9", "3.3")
				val, err := newOC(t, map[streams.StreamID]*mockPipeline{sid: p}).Observe(ctx, sid, opts)
				assert.Nil(t, val)
				var qerr QuoteInvariantError
				require.ErrorAs(t, err, &qerr)
				assert.Equal(t, sid, qerr.StreamID)
			})
		}
	})

	t.Run("returns a valid terminal quote with quote named tasks", func(t *testing.T) {
		t.Parallel()
		sid := streams.StreamID(1)
		p := makeNamedQuotePipeline(t, sid, nil, [3]string{"benchmark_price", "bid_price", "ask_price"}, "2.2", "1.1", "3.3")
		val, err := newOC(t, map[streams.StreamID]*mockPipeline{sid: p}).Observe(ctx, sid, opts)
		require.NoError(t, err)
		assert.True(t, val.(*lloprotocol.Quote).IsValid())
	})

	t.Run("does not validate a terminal quote without quote named tasks", func(t *testing.T) {
		t.Parallel()
		for name, names := range map[string][3]string{
			"unrelated names":  {"price", "low", "high"},
			"substrings only":  {"benchmarks", "forbid", "task"},
			"camel case":       {"benchmarkPrice", "bidPrice", "askPrice"},
			"ambiguous name":   {"benchmark_price", "bid_ask_spread", "ask_price"},
			"one name missing": {"benchmark_price", "bid_price", "result"},
			"timestamp field":  {"mid_timestamp", "bid_price", "ask_price"},
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				sid := streams.StreamID(1)
				p := makeNamedQuotePipeline(t, sid, nil, names, "2.2", "9.9", "3.3")
				val, err := newOC(t, map[streams.StreamID]*mockPipeline{sid: p}).Observe(ctx, sid, opts)
				require.NoError(t, err)
				assert.False(t, val.(*lloprotocol.Quote).IsValid())
			})
		}
	})

	t.Run("validates terminals named as quote fields regardless of their index", func(t *testing.T) {
		t.Parallel()
		sid := streams.StreamID(1)
		names := [3]string{"ask_price", "benchmark_price", "bid_price"}

		p := makeNamedQuotePipeline(t, sid, nil, names, "3.3", "2.2", "1.1")
		_, err := newOC(t, map[streams.StreamID]*mockPipeline{sid: p}).Observe(ctx, sid, opts)
		require.NoError(t, err)

		p = makeNamedQuotePipeline(t, sid, nil, names, "3.3", "2.2", "9.9")
		_, err = newOC(t, map[streams.StreamID]*mockPipeline{sid: p}).Observe(ctx, sid, opts)
		var qerr QuoteInvariantError
		require.ErrorAs(t, err, &qerr)
		assert.Equal(t, sid, qerr.StreamID)
		assert.Equal(t, "9.9", qerr.Quote.Bid.String())
		assert.Equal(t, "2.2", qerr.Quote.Benchmark.String())
		assert.Equal(t, "3.3", qerr.Quote.Ask.String())
	})

	t.Run("fails every stream of a flex median pipeline with bid, mid and ask", func(t *testing.T) {
		t.Parallel()
		// Flex orders fields alphabetically and adds a timestamp median, as for
		// USD/UAH-Datalink-NCFX-ForexPrice.
		jobSID := streams.StreamID(1000004349)
		for name, tc := range map[string]struct {
			names  []string
			values []string
		}{
			"bid, mid and ask": {
				names:  []string{"price_median_ask", "price_median_bid", "price_median_mid", "price_mid_timestamp"},
				values: []string{"3.3", "9.9", "2.2", "1761148800000"},
			},
			"with volumes and last traded price": {
				names:  []string{"price_median_ask", "price_median_askvolume", "price_median_bid", "price_median_bidvolume", "price_median_lasttradedprice", "price_median_mid", "price_mid_timestamp"},
				values: []string{"3.3", "100", "9.9", "100", "2.5", "2.2", "1761148800000"},
			},
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				tags := make([]streams.StreamID, len(tc.names))
				for i := range tags {
					tags[i] = streams.StreamID(1000004350 + i)
				}
				p := makeNamedTerminalsPipeline(t, jobSID, tags, tc.names, tc.values)
				pipelines := map[streams.StreamID]*mockPipeline{}
				for _, sid := range tags {
					pipelines[sid] = p
				}
				oc := newOC(t, pipelines)

				for _, sid := range tags {
					val, err := oc.Observe(ctx, sid, opts)
					assert.Nil(t, val, "stream %d should have no value", sid)
					var qerr QuoteInvariantError
					require.ErrorAs(t, err, &qerr, "stream %d", sid)
					assert.Equal(t, jobSID, qerr.StreamID)
				}
				assert.Equal(t, int32(1), p.runCount.Load(), "pipeline should only run once")
			})
		}
	})

	t.Run("returns every stream of a valid flex median pipeline with bid, mid and ask", func(t *testing.T) {
		t.Parallel()
		jobSID := streams.StreamID(1000004349)
		names := []string{"price_median_ask", "price_median_bid", "price_median_mid", "price_mid_timestamp"}
		tags := []streams.StreamID{1000004353, 1000004352, 1000004350, 1000004351}
		p := makeNamedTerminalsPipeline(t, jobSID, tags, names, []string{"3.3", "1.1", "2.2", "1761148800000"})
		pipelines := map[streams.StreamID]*mockPipeline{}
		for _, sid := range tags {
			pipelines[sid] = p
		}
		oc := newOC(t, pipelines)

		for _, sid := range tags {
			val, err := oc.Observe(ctx, sid, opts)
			require.NoError(t, err, "stream %d", sid)
			assert.IsType(t, &lloprotocol.Decimal{}, val, "stream %d", sid)
		}
	})

	t.Run("does not validate terminals that do not name each quote field once", func(t *testing.T) {
		t.Parallel()
		for name, names := range map[string][]string{
			"duplicate role": {"bid_a", "bid_b", "benchmark", "ask"},
			"best bid ask":   {"best_bid", "best_ask", "bid_volume", "ask_volume"},
			"no benchmark":   {"price_median_ask", "price_median_bid", "price_median_lasttradedprice"},
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				tags := make([]streams.StreamID, len(names))
				values := make([]string, len(names))
				for i := range names {
					tags[i] = streams.StreamID(100 + i)
					values[i] = strconv.Itoa(len(names) - i)
				}
				p := makeNamedTerminalsPipeline(t, 99, tags, names, values)
				pipelines := map[streams.StreamID]*mockPipeline{}
				for _, sid := range tags {
					pipelines[sid] = p
				}
				oc := newOC(t, pipelines)

				for _, sid := range tags {
					_, err := oc.Observe(ctx, sid, opts)
					require.NoError(t, err, "stream %d", sid)
				}
			})
		}
	})

	t.Run("fails an invalid multi-source median quote", func(t *testing.T) {
		t.Parallel()
		known := [3]string{"data,mid", "data,bid", "data,ask"}
		for name, tags := range map[string][]streams.StreamID{
			"untagged medians": nil,
			"tagged medians":   {11, 12, 13},
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				jobSID := streams.StreamID(15)
				p := makeMedianQuotePipeline(t, jobSID, tags, [][3]string{known, known, known, known}, "2.2", "9.9", "3.3")
				pipelines := map[streams.StreamID]*mockPipeline{jobSID: p}
				for _, sid := range tags {
					pipelines[sid] = p
				}
				oc := newOC(t, pipelines)

				for _, sid := range append(slices.Clone(tags), jobSID) {
					val, err := oc.Observe(ctx, sid, opts)
					assert.Nil(t, val, "stream %d should have no value", sid)
					var qerr QuoteInvariantError
					require.ErrorAs(t, err, &qerr, "stream %d", sid)
					assert.Equal(t, jobSID, qerr.StreamID)
				}
				assert.Equal(t, int32(1), p.runCount.Load(), "pipeline should only run once")
			})
		}
	})

	t.Run("returns a valid multi-source median quote", func(t *testing.T) {
		t.Parallel()
		jobSID := streams.StreamID(15)
		known := [3]string{"data,mid_price", "data,bid_price", "data,ask_price"}
		p := makeMedianQuotePipeline(t, jobSID, nil, [][3]string{known, known, known}, "2.2", "1.1", "3.3")
		val, err := newOC(t, map[streams.StreamID]*mockPipeline{jobSID: p}).Observe(ctx, jobSID, opts)
		require.NoError(t, err)
		assert.True(t, val.(*lloprotocol.Quote).IsValid())
	})

	t.Run("fails an invalid dex cross exchange quote", func(t *testing.T) {
		t.Parallel()
		jobSID := streams.StreamID(1000001458)
		p := makeDexCrossQuotePipeline(t, jobSID, "2.2", "9.9", "3.3")
		val, err := newOC(t, map[streams.StreamID]*mockPipeline{jobSID: p}).Observe(ctx, jobSID, opts)
		assert.Nil(t, val)
		var qerr QuoteInvariantError
		require.ErrorAs(t, err, &qerr)
		assert.Equal(t, jobSID, qerr.StreamID)
	})

	t.Run("returns a valid dex cross exchange quote", func(t *testing.T) {
		t.Parallel()
		jobSID := streams.StreamID(1000001458)
		p := makeDexCrossQuotePipeline(t, jobSID, "2.2", "1.1", "3.3")
		val, err := newOC(t, map[streams.StreamID]*mockPipeline{jobSID: p}).Observe(ctx, jobSID, opts)
		require.NoError(t, err)
		assert.True(t, val.(*lloprotocol.Quote).IsValid())
	})

	t.Run("validates a quote reading one path for all fields", func(t *testing.T) {
		t.Parallel()
		for _, path := range []string{"data,result", "data,index_price"} {
			t.Run(path, func(t *testing.T) {
				t.Parallel()
				jobSID := streams.StreamID(1000007905)
				same := [3]string{path, path, path}

				p := makeMedianQuotePipeline(t, jobSID, nil, [][3]string{same, same}, "2.2", "2.2", "2.2")
				val, err := newOC(t, map[streams.StreamID]*mockPipeline{jobSID: p}).Observe(ctx, jobSID, opts)
				require.NoError(t, err)
				assert.True(t, val.(*lloprotocol.Quote).IsValid())

				p = makeMedianQuotePipeline(t, jobSID, nil, [][3]string{same, same}, "2.2", "9.9", "3.3")
				_, err = newOC(t, map[streams.StreamID]*mockPipeline{jobSID: p}).Observe(ctx, jobSID, opts)
				var qerr QuoteInvariantError
				require.ErrorAs(t, err, &qerr)
				assert.Equal(t, jobSID, qerr.StreamID)
			})
		}
	})

	t.Run("does not fail the tagged streams of a funding rate pipeline", func(t *testing.T) {
		t.Parallel()
		// The implicit job-level quote is Benchmark=fundingRate, Bid=fundingTimestamp,
		// Ask=epochDuration, which always violates the invariant.
		jobSID := streams.StreamID(1000001306)
		tags := []streams.StreamID{1000001007, 1000001015, 1000001016}
		names := [3]string{"ds_funding_rate", "ds_funding_timestamp", "ds_epoch_duration"}
		p := makeNamedQuotePipeline(t, jobSID, tags, names, "0.0001", "1761148800000", "28800")
		pipelines := map[streams.StreamID]*mockPipeline{jobSID: p}
		for _, sid := range tags {
			pipelines[sid] = p
		}
		oc := newOC(t, pipelines)

		for _, sid := range tags {
			val, err := oc.Observe(ctx, sid, opts)
			require.NoError(t, err, "stream %d", sid)
			assert.IsType(t, &lloprotocol.Decimal{}, val, "stream %d", sid)
		}
	})

	t.Run("fails every stream of an invalid perps price quote", func(t *testing.T) {
		t.Parallel()
		for name, normalize := range map[string]bool{
			"without normalize": false,
			"with normalize":    true,
		} {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				jobSID := streams.StreamID(1000008984)
				tags := [3]streams.StreamID{1000008985, 1000008986, 1000008987}
				p := makePerpsPriceQuotePipeline(t, jobSID, tags, normalize, "2.2", "9.9", "3.3")
				pipelines := map[streams.StreamID]*mockPipeline{jobSID: p}
				for _, sid := range tags {
					pipelines[sid] = p
				}
				oc := newOC(t, pipelines)

				for _, sid := range append(tags[:], jobSID) {
					val, err := oc.Observe(ctx, sid, opts)
					assert.Nil(t, val, "stream %d should have no value", sid)
					var qerr QuoteInvariantError
					require.ErrorAs(t, err, &qerr, "stream %d", sid)
					assert.Equal(t, jobSID, qerr.StreamID)
				}
				assert.Equal(t, int32(1), p.runCount.Load(), "pipeline should only run once")
			})
		}
	})

	t.Run("returns every stream of a valid perps price quote", func(t *testing.T) {
		t.Parallel()
		jobSID := streams.StreamID(1000008984)
		tags := [3]streams.StreamID{1000008985, 1000008986, 1000008987}
		p := makePerpsPriceQuotePipeline(t, jobSID, tags, true, "2.2", "1.1", "3.3")
		pipelines := map[streams.StreamID]*mockPipeline{jobSID: p}
		for _, sid := range tags {
			pipelines[sid] = p
		}
		oc := newOC(t, pipelines)

		for _, sid := range tags {
			val, err := oc.Observe(ctx, sid, opts)
			require.NoError(t, err, "stream %d", sid)
			assert.IsType(t, &lloprotocol.Decimal{}, val, "stream %d", sid)
		}
		val, err := oc.Observe(ctx, jobSID, opts)
		require.NoError(t, err)
		assert.True(t, val.(*lloprotocol.Quote).IsValid())
	})

	t.Run("fails the tagged siblings of an invalid terminal quote with quote named tasks", func(t *testing.T) {
		t.Parallel()
		jobSID := streams.StreamID(15)
		tags := []streams.StreamID{11, 12, 13}
		p := makeNamedQuotePipeline(t, jobSID, tags, [3]string{"benchmark_price", "bid_price", "ask_price"}, "2.2", "9.9", "3.3")
		pipelines := map[streams.StreamID]*mockPipeline{jobSID: p}
		for _, sid := range tags {
			pipelines[sid] = p
		}
		oc := newOC(t, pipelines)

		for _, sid := range append(slices.Clone(tags), jobSID) {
			val, err := oc.Observe(ctx, sid, opts)
			assert.Nil(t, val, "stream %d should have no value", sid)
			var qerr QuoteInvariantError
			require.ErrorAs(t, err, &qerr, "stream %d", sid)
			assert.Equal(t, jobSID, qerr.StreamID)
		}
		assert.Equal(t, int32(1), p.runCount.Load(), "pipeline should only run once")
	})

	t.Run("returns the tagged siblings of a valid terminal quote with quote named tasks", func(t *testing.T) {
		t.Parallel()
		jobSID := streams.StreamID(15)
		tags := []streams.StreamID{11, 12, 13}
		p := makeNamedQuotePipeline(t, jobSID, tags, [3]string{"benchmark_price", "bid_price", "ask_price"}, "2.2", "1.1", "3.3")
		pipelines := map[streams.StreamID]*mockPipeline{jobSID: p}
		for _, sid := range tags {
			pipelines[sid] = p
		}
		oc := newOC(t, pipelines)

		for _, sid := range tags {
			val, err := oc.Observe(ctx, sid, opts)
			require.NoError(t, err, "stream %d", sid)
			assert.IsType(t, &lloprotocol.Decimal{}, val, "stream %d", sid)
		}
		val, err := oc.Observe(ctx, jobSID, opts)
		require.NoError(t, err)
		assert.True(t, val.(*lloprotocol.Quote).IsValid())
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
