package llo

import (
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-data-streams/llo"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	clhttptest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/httptest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/null"
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

func makePipelineWithMultipleStreamResults(streamIDs []streams.StreamID, results []interface{}) *mockPipeline {
	if len(streamIDs) != len(results) {
		panic("streamIDs and results must have the same length")
	}
	trrs := make([]pipeline.TaskRunResult, len(streamIDs))
	for i, res := range results {
		trrs[i] = pipeline.TaskRunResult{Task: &pipeline.MemoTask{BaseTask: pipeline.BaseTask{StreamID: null.Uint32From(streamIDs[i])}}, Result: pipeline.Result{Value: res}}
	}
	return &mockPipeline{
		run:       &pipeline.Run{},
		trrs:      trrs,
		err:       nil,
		streamIDs: streamIDs,
	}
}

func TestObservationContext_Observe(t *testing.T) {
	ctx := tests.Context(t)
	r := &mockRegistry{}
	telem := &mockTelemeter{}
	oc := newObservationContext(r, telem)
	opts := llo.DSOpts(nil)

	missingStreamID := streams.StreamID(0)
	streamID1 := streams.StreamID(1)
	streamID2 := streams.StreamID(2)
	streamID3 := streams.StreamID(3)
	streamID4 := streams.StreamID(4)
	streamID5 := streams.StreamID(5)
	streamID6 := streams.StreamID(6)

	multiPipelineDecimal := makePipelineWithMultipleStreamResults([]streams.StreamID{streamID4, streamID5, streamID6}, []interface{}{decimal.NewFromFloat(12.34), decimal.NewFromFloat(56.78), decimal.NewFromFloat(90.12)})

	r.pipelines = map[streams.StreamID]*mockPipeline{
		streamID1: &mockPipeline{},
		streamID2: makePipelineWithSingleResult[decimal.Decimal](rand.Int64(), decimal.NewFromFloat(12.34), nil),
		streamID3: makeErroringPipeline(),
		streamID4: multiPipelineDecimal,
		streamID5: multiPipelineDecimal,
		streamID6: multiPipelineDecimal,
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

		assert.Equal(t, 1, multiPipelineDecimal.runCount)

		// returns cached values on subsequent calls
		val, err = oc.Observe(ctx, streamID6, opts)
		require.NoError(t, err)
		assert.Equal(t, "90.12", val.(*llo.Decimal).String())

		assert.Equal(t, 1, multiPipelineDecimal.runCount)
	})
}

type mockPipelineConfig struct{}

func (m *mockPipelineConfig) DefaultHTTPLimit() int64 { return 10000 }
func (m *mockPipelineConfig) DefaultHTTPTimeout() commonconfig.Duration {
	return *commonconfig.MustNewDuration(1 * time.Hour)
}
func (m *mockPipelineConfig) MaxRunDuration() time.Duration  { return 1 * time.Hour }
func (m *mockPipelineConfig) ReaperInterval() time.Duration  { return 0 }
func (m *mockPipelineConfig) ReaperThreshold() time.Duration { return 0 }
func (m *mockPipelineConfig) VerboseLogging() bool           { return true }

type mockBridgeConfig struct{}

func (m *mockBridgeConfig) BridgeResponseURL() *url.URL {
	return nil
}
func (m *mockBridgeConfig) BridgeCacheTTL() time.Duration {
	return 0
}

func createBridge(t *testing.T, name string, val string, borm bridges.ORM) {
	callcount := 0
	bridge := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		callcount++
		if callcount > 1 {
			t.Fatal("expected only one call to the bridge")
		}
		_, herr := io.ReadAll(req.Body)
		if herr != nil {
			t.Fatal(herr)
		}

		res.WriteHeader(http.StatusOK)
		resp := fmt.Sprintf(`{"result": %s}`, val)
		_, herr = res.Write([]byte(resp))
		if herr != nil {
			t.Fatal(herr)
		}
	}))
	t.Cleanup(bridge.Close)
	u, _ := url.Parse(bridge.URL)
	require.NoError(t, borm.CreateBridgeType(tests.Context(t), &bridges.BridgeType{
		Name: bridges.BridgeName(name),
		URL:  models.WebURL(*u),
	}))
}

func TestObservationContext_Observe_integrationRealPipeline(t *testing.T) {
	ctx := tests.Context(t)
	lggr := logger.TestLogger(t)
	db := pgtest.NewSqlxDB(t)
	bridgesORM := bridges.NewORM(db)

	createBridge(t, "foo-bridge", `123.456`, bridgesORM)
	createBridge(t, "bar-bridge", `"124.456"`, bridgesORM)

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
		oc := newObservationContext(r, telem)
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
