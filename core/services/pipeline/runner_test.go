package pipeline_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

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

	s := fmt.Sprintf(`
ds1 [type=bridge name="example-bridge" timeout=0 requestData=<{"data": {"coin": "BTC", "market": "USD"}}>]
ds1_parse [type=jsonparse lax=false  path="data,result"]
ds1_multiply [type=multiply times=1000000000000000000]

ds2 [type=http method="GET" url="%s" requestData=<{"data": {"coin": "BTC", "market": "USD"}}>]
ds2_parse [type=jsonparse lax=false  path="data,result"]
ds2_multiply [type=multiply times=1000000000000000000]

ds3 [type=http method="GET" url="blah://test.invalid" requestData=<{"data": {"coin": "BTC", "market": "USD"}}>]
ds3_parse [type=jsonparse lax=false  path="data,result"]
ds3_multiply [type=multiply times=1000000000000000000]

ds1->ds1_parse->ds1_multiply->median;
ds2->ds2_parse->ds2_multiply->median;
ds3->ds3_parse->ds3_multiply->median;

median [type=median index=0]
ds4 [type=http method="GET" url="%s" index=1]
ds5 [type=http method="GET" url="%s" index=2]
`, s2.URL, s4.URL, s5.URL)
	d, err := pipeline.Parse(s)
	require.NoError(t, err)

	spec := pipeline.Spec{
		DotDagSource: s,
	}
	_, trrs, err := r.ExecuteRun(context.Background(), spec, nil, pipeline.JSONSerializable{}, *logger.Default)
	require.NoError(t, err)
	require.Len(t, trrs, len(d.Tasks))

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
		if trr.Result.Error != nil && !trr.IsTerminal() {
			errorResults = append(errorResults, trr)
		}
	}
	// There are three tasks in the erroring pipeline
	require.Len(t, errorResults, 3)
}

func Test_PipelineRunner_ExecuteTaskRunsWithVars(t *testing.T) {
	t.Parallel()

	specTemplate := `
        ds1 [type=bridge name="example-bridge" timeout=0 requestData=<{"data": $(input.foo)}>]
        ds1_parse [type=jsonparse lax=false  path="data,result" data="$(ds1)"]
        ds1_multiply [type=multiply input="$(ds1_parse.result)" times="$(ds1_parse.times)"]

        ds2 [type=http method="POST" url="%s" requestData=<{"data": [ $(input.bar), $(input.baz) ]}>]
        ds2_parse [type=jsonparse lax=false  path="data" data="$(ds2)"]
        ds2_multiply [type=multiply input="$(ds2_parse.result)" times="$(ds2_parse.times)"]

        ds3 [type=http method="POST" url="blah://test.invalid" requestData=<{"data": {"coin": "BTC", "market": "USD"}}>]
        ds3_parse [type=jsonparse lax=false  path="data,result" data="$(ds3)"]
        ds3_multiply [type=multiply input="$(ds3_parse.value)" times="$(ds3_parse.times)"]

        ds1->ds1_parse->ds1_multiply->median;
        ds2->ds2_parse->ds2_multiply->median;
        ds3->ds3_parse->ds3_multiply->median;

        median [type=median values=<[ $(ds1_multiply), $(ds2_multiply), $(ds3_multiply) ]> index=0]
        ds4 [type=http method="GET" url="%s" index=1]

        submit [type=bridge name="submit"
                includeInputAtKey="%s"
                requestData=<{
                    "median": $(median),
                    "fetchedValues": [ $(ds1_parse.result), $(ds2_parse.result) ],
                    "someString": $(ds4)
                }>]

        median -> submit;
        ds4 -> submit;
    `

	tests := []struct {
		name              string
		pipelineInput     map[string]interface{}
		meta              map[string]interface{}
		includeInputAtKey string
	}{
		{
			name: "meta + includeInputAtKey",
			pipelineInput: map[string]interface{}{
				"foo": []interface{}{float64(123), "chainlink"},
				"bar": float64(123.45),
				"baz": "such oracle",
			},
			meta:              map[string]interface{}{"roundID": float64(456), "latestAnswer": float64(654)},
			includeInputAtKey: "sergey",
		},
		{
			name: "includeInputAtKey",
			pipelineInput: map[string]interface{}{
				"foo": *mustDecimal(t, "42.1337"),
				"bar": map[string]interface{}{"steve": "chainlink"},
				"baz": true,
			},
			includeInputAtKey: "best oracles",
		},
		{
			name: "meta",
			pipelineInput: map[string]interface{}{
				"foo": []interface{}{"asdf", float64(123)},
				"bar": false,
				"baz": *mustDecimal(t, "42.1337"),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			expectedRequestDS1 := map[string]interface{}{"data": test.pipelineInput["foo"]}
			expectedRequestDS2 := map[string]interface{}{"data": []interface{}{test.pipelineInput["bar"], test.pipelineInput["baz"]}}
			expectedRequestSubmit := map[string]interface{}{
				"median":        "9650000000000000000000",
				"fetchedValues": []interface{}{"9700", "9600"},
				"someString":    "some random string",
			}
			if test.meta != nil {
				expectedRequestDS1["meta"] = test.meta
				expectedRequestSubmit["meta"] = test.meta
			} else {
				expectedRequestDS1["meta"] = nil
				expectedRequestSubmit["meta"] = nil
			}
			if test.includeInputAtKey != "" {
				expectedRequestSubmit[test.includeInputAtKey] = "9650000000000000000000"
			}

			// 1. Setup bridge
			ds1 := makeBridge(t, store, "example-bridge", expectedRequestDS1, map[string]interface{}{
				"data": map[string]interface{}{
					"result": map[string]interface{}{
						"result": decimal.NewFromInt(9700),
						"times":  "1000000000000000000",
					},
				},
			})
			defer ds1.Close()

			// 2. Setup success HTTP
			ds2 := httptest.NewServer(fakeExternalAdapter(t, expectedRequestDS2, map[string]interface{}{
				"data": map[string]interface{}{
					"result": decimal.NewFromInt(9600),
					"times":  "1000000000000000000",
				},
			}))
			defer ds2.Close()

			ds4 := httptest.NewServer(fakeStringResponder(t, "some random string"))
			defer ds4.Close()

			// 3. Setup final bridge task
			submit := makeBridge(t, store, "submit", expectedRequestSubmit, map[string]interface{}{"ok": true})
			defer submit.Close()

			orm := new(mocks.ORM)
			orm.On("DB").Return(store.DB)

			runner := pipeline.NewRunner(orm, store.Config)
			specStr := fmt.Sprintf(specTemplate, ds2.URL, ds4.URL, test.includeInputAtKey)
			p, err := pipeline.Parse(specStr)
			require.NoError(t, err)

			spec := pipeline.Spec{
				DotDagSource: specStr,
			}
			var meta pipeline.JSONSerializable
			if test.meta != nil {
				meta.Val = test.meta
			} else {
				meta.Null = true
			}
			_, taskRunResults, err := runner.ExecuteRun(context.Background(), spec, test.pipelineInput, meta, *logger.Default)
			require.NoError(t, err)
			require.Len(t, taskRunResults, len(p.Tasks))

			type M = map[string]interface{}
			expectedResults := map[string]pipeline.Result{
				"ds1":          {Value: `{"data":{"result":{"result":"9700","times":"1000000000000000000"}}}` + "\n"},
				"ds1_parse":    {Value: M{"result": "9700", "times": "1000000000000000000"}},
				"ds1_multiply": {Value: *mustDecimal(t, "9700000000000000000000")},
				"ds2":          {Value: `{"data":{"result":"9600","times":"1000000000000000000"}}` + "\n"},
				"ds2_parse":    {Value: M{"result": "9600", "times": "1000000000000000000"}},
				"ds2_multiply": {Value: *mustDecimal(t, "9600000000000000000000")},
				"ds3":          {Error: errors.New(`error making http request: Post "blah://test.invalid": unsupported protocol scheme "blah"`)},
				"ds3_parse":    {Error: pipeline.ErrTooManyErrors},
				"ds3_multiply": {Error: pipeline.ErrTooManyErrors},
				"ds4":          {Value: "some random string"},
				"median":       {Value: *mustDecimal(t, "9650000000000000000000")},
				"submit":       {Value: `{"ok":true}` + "\n"},
			}

			for _, r := range taskRunResults {
				expected := expectedResults[r.Task.DotID()]
				if r.Result.Error != nil {
					require.Equal(t, expected.Error.Error(), r.Result.Error.Error())
				} else {
					if d, is := expected.Value.(decimal.Decimal); is {
						require.Equal(t, d.String(), r.Result.Value.(decimal.Decimal).String())
					} else {
						require.Equal(t, expected.Value, r.Result.Value)
					}
				}
			}
		})
	}
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
		if trr.IsTerminal() {
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
	assert.Equal(t, true, trrs.FinalResult().HasErrors())
	assert.IsType(t, pipeline.ErrRunPanicked{}, trrs.FinalResult().Errors[0])
}
