package pipeline_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/stretchr/testify/assert"

	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/require"
)

func Test_PipelineRunner_ExecuteTaskRuns(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	btcUSDPairing := utils.MustUnmarshalToMap(`{"data":{"coin":"BTC","market":"USD"}}`)

	// 1. Setup bridge
	s1 := httptest.NewServer(fakePriceResponder(t, btcUSDPairing, decimal.NewFromInt(9700), "", nil))
	defer s1.Close()

	bridgeFeedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)
	bridgeFeedWebURL := (*models.WebURL)(bridgeFeedURL)

	_, bridge := cltest.NewBridgeType(t, "example-bridge")
	bridge.URL = *bridgeFeedWebURL
	require.NoError(t, store.ORM.DB.Create(&bridge).Error)

	// 2. Setup success HTTP
	s2 := httptest.NewServer(fakePriceResponder(t, btcUSDPairing, decimal.NewFromInt(9600), "", nil))
	defer s2.Close()

	s4 := httptest.NewServer(fakeStringResponder(t, "foo-index-1"))
	defer s4.Close()
	s5 := httptest.NewServer(fakeStringResponder(t, "bar-index-2"))
	defer s5.Close()

	orm := new(mocks.ORM)
	orm.On("DB").Return(store.DB)

	r := pipeline.NewRunner(orm, store.Config)

	d := pipeline.TaskDAG{}
	s := fmt.Sprintf(`
ds1 [type=bridge name="example-bridge" timeout=0 requestData="{\"data\": {\"coin\": \"BTC\", \"market\": \"USD\"}}"]
ds1_parse [type=jsonparse lax=false  path="data,result"]
ds1_multiply [type=multiply times=1000000000000000000]

ds2 [type=http method="GET" url="%s" requestData="{\"data\": {\"coin\": \"BTC\", \"market\": \"USD\"}}"]
ds2_parse [type=jsonparse lax=false  path="data,result"]
ds2_multiply [type=multiply times=1000000000000000000]

ds3 [type=http method="GET" url="blah://test.invalid" requestData="{\"data\": {\"coin\": \"BTC\", \"market\": \"USD\"}}"]
ds3_parse [type=jsonparse lax=false  path="data,result"]
ds3_multiply [type=multiply times=1000000000000000000]

ds1->ds1_parse->ds1_multiply->median;
ds2->ds2_parse->ds2_multiply->median;
ds3->ds3_parse->ds3_multiply->median;

median [type=median index=0]
ds4 [type=http method="GET" url="%s" index=1]
ds5 [type=http method="GET" url="%s" index=2]
`, s2.URL, s4.URL, s5.URL)
	err = d.UnmarshalText([]byte(s))
	require.NoError(t, err)
	ts, err := d.TasksInDependencyOrder()
	require.NoError(t, err)

	spec := pipeline.Spec{
		DotDagSource: s,
	}
	_, trrs, err := r.ExecuteRun(context.Background(), spec, nil, pipeline.JSONSerializable{}, *logger.Default)
	require.NoError(t, err)
	require.Len(t, trrs, len(ts))

	finalResults := trrs.FinalResult()
	require.Len(t, finalResults.Values, 3)
	require.Len(t, finalResults.Errors, 3)
	assert.Equal(t, "9650000000000000000000", finalResults.Values[0].(decimal.Decimal).String())
	assert.Nil(t, finalResults.Errors[0])
	assert.Equal(t, "foo-index-1", finalResults.Values[1].(string))
	assert.Nil(t, finalResults.Errors[1])
	assert.Equal(t, "bar-index-2", finalResults.Values[2].(string))
	assert.Nil(t, finalResults.Errors[2])

	var errorResults []pipeline.TaskRunResult
	for _, trr := range trrs {
		if trr.Result.Error != nil && !trr.IsTerminal {
			errorResults = append(errorResults, trr)
		}
	}
	// There are three tasks in the erroring pipeline
	require.Len(t, errorResults, 3)
}

func Test_PipelineRunner_HandleFaults(t *testing.T) {
	// We want to test the scenario where one or multiple APIs time out,
	// but a sufficient number of them still complete within the desired time frame
	// and so we can still obtain a median.
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	orm := new(mocks.ORM)
	orm.On("DB").Return(store.DB)
	m1 := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		time.Sleep(100 * time.Millisecond)
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(`{"result":10}`))
	}))
	m2 := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(`{"result":11}`))
	}))
	s := fmt.Sprintf(`
ds1          [type=http url="%s"];
ds1_parse    [type=jsonparse path="result"];
ds1_multiply [type=multiply times=100];

ds2          [type=http url="%s"];
ds2_parse    [type=jsonparse path="result"];
ds2_multiply [type=multiply times=100];

ds1 -> ds1_parse -> ds1_multiply -> answer1;
ds2 -> ds2_parse -> ds2_multiply -> answer1;

answer1 [type=median                      index=0];
`, m1.URL, m2.URL)

	r := pipeline.NewRunner(orm, store.Config)

	// If we cancel before an API is finished, we should still get a median.
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	spec := pipeline.Spec{
		DotDagSource: s,
	}
	_, trrs, err := r.ExecuteRun(ctx, spec, nil, pipeline.JSONSerializable{}, *logger.Default)
	require.NoError(t, err)
	for _, trr := range trrs {
		if trr.IsTerminal {
			require.Equal(t, decimal.RequireFromString("1100"), trr.Result.Value.(decimal.Decimal))
		}
	}
}

func TestPanicTask_Run(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	orm := new(mocks.ORM)
	orm.On("DB").Return(store.DB)
	s := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(`{"result":10}`))
	}))
	r := pipeline.NewRunner(orm, store.Config)
	_, trrs, err := r.ExecuteRun(context.Background(), pipeline.Spec{
		DotDagSource: fmt.Sprintf(`
ds1 [type=http url="%s"]
ds_parse [type=jsonparse path="result"]
ds_multiply [type=multiply times=10]
ds_panic [type=panic msg="oh no"]
ds1->ds_parse->ds_multiply->ds_panic;`, s.URL),
	}, nil, pipeline.JSONSerializable{}, *logger.Default)
	require.NoError(t, err)
	require.Equal(t, 4, len(trrs))
	assert.Equal(t, []interface{}{nil}, trrs.FinalResult().Values)
	assert.Equal(t, pipeline.ErrRunPanicked.Error(), trrs.FinalResult().Errors[0].Error())
	for _, trr := range trrs {
		assert.Equal(t, null.NewString("pipeline run panicked", true), trr.Result.ErrorDB())
		assert.Equal(t, true, trr.Result.OutputDB().Null)
	}
}
