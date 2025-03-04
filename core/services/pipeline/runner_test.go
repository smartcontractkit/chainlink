package pipeline_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	bridgesMocks "github.com/smartcontractkit/chainlink/v2/core/bridges/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	clhttptest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/httptest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func newRunner(t testing.TB, db *sqlx.DB, bridgeORM bridges.ORM, cfg chainlink.GeneralConfig) (pipeline.Runner, *mocks.ORM) {
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		ChainConfigs:   cfg.EVMConfigs(),
		DatabaseConfig: cfg.Database(),
		FeatureConfig:  cfg.Feature(),
		ListenerConfig: cfg.Database().Listener(),
		DB:             db,
		KeyStore:       ethKeyStore,
	})
	orm := mocks.NewORM(t)
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	r := pipeline.NewRunner(orm, bridgeORM, cfg.JobPipeline(), cfg.WebServer(), legacyChains, ethKeyStore, nil, logger.TestLogger(t), c, c)
	return r, orm
}

func Test_PipelineRunner_ExecuteTaskRuns(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)

	btcUSDPairing := utils.MustUnmarshalToMap(`{"data":{"coin":"BTC","market":"USD"}}`)

	// 1. Setup bridge
	s1 := httptest.NewServer(fakePriceResponder(t, btcUSDPairing, decimal.NewFromInt(9700), "", nil))
	defer s1.Close()

	bridgeFeedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)

	_, bt := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: bridgeFeedURL.String()})

	btORM := bridgesMocks.NewORM(t)
	btORM.On("FindBridge", mock.Anything, bt.Name).Return(*bt, nil).Once()

	// 2. Setup success HTTP
	s2 := httptest.NewServer(fakePriceResponder(t, btcUSDPairing, decimal.NewFromInt(9600), "", nil))
	defer s2.Close()

	s4 := httptest.NewServer(fakeStringResponder(t, "foo-index-1"))
	defer s4.Close()
	s5 := httptest.NewServer(fakeStringResponder(t, "bar-index-2"))
	defer s5.Close()

	r, _ := newRunner(t, db, btORM, cfg)

	s := fmt.Sprintf(`
ds1 [type=bridge name="%s" timeout=0 requestData=<{"data": {"coin": "BTC", "market": "USD"}}>]
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
`, bt.Name.String(), s2.URL, s4.URL, s5.URL)
	d, err := pipeline.Parse(s)
	require.NoError(t, err)

	spec := pipeline.Spec{DotDagSource: s}
	vars := pipeline.NewVarsFrom(nil)

	_, trrs, err := r.ExecuteRun(testutils.Context(t), spec, vars)
	require.NoError(t, err)
	require.Len(t, trrs, len(d.Tasks))

	finalResults := trrs.FinalResult()
	require.Len(t, finalResults.Values, 3)
	require.Len(t, finalResults.AllErrors, 12)
	require.Len(t, finalResults.FatalErrors, 3)
	assert.Equal(t, "9650000000000000000000", finalResults.Values[0].(decimal.Decimal).String())
	assert.Nil(t, finalResults.FatalErrors[0])
	assert.Equal(t, "foo-index-1", finalResults.Values[1].(string))
	assert.Nil(t, finalResults.FatalErrors[1])
	assert.Equal(t, "bar-index-2", finalResults.Values[2].(string))
	assert.Nil(t, finalResults.FatalErrors[2])

	var errorResults []pipeline.TaskRunResult
	for _, trr := range trrs {
		if trr.Result.Error != nil && !trr.IsTerminal() {
			errorResults = append(errorResults, trr)
		}
	}
	// There are three tasks in the erroring pipeline
	require.Len(t, errorResults, 3)
}

func Test_PipelineRunner_ExecuteEthAbiDecode(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)

	mockResult := `{"data":{"result":"0x000000000000000000000000000000000000000000000000105ba6a589b23a81"}}`
	s1 := httptest.NewServer(NewMockHandler(mockResult))
	defer s1.Close()

	bridgeFeedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)

	_, bt := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: bridgeFeedURL.String()})

	btORM := bridgesMocks.NewORM(t)
	btORM.On("FindBridge", mock.Anything, bt.Name).Return(*bt, nil).Once()

	r, _ := newRunner(t, db, btORM, cfg)

	s := fmt.Sprintf(`
		ds1 [type=bridge name="%s" timeout=0 requestData=<{"data": {"address": "0x1234"}}>]
		ds1_parse [type=jsonparse path="data,result"]  
		ds1_decode [type=ethabidecode abi="int256 data" data="$(ds1_parse)"];
		ds1_value [type="multiply" input="$(ds1_decode.data)" times=1]

		ds1->ds1_parse->ds1_decode->ds1_value

`, bt.Name.String())
	d, err := pipeline.Parse(s)
	require.NoError(t, err)

	spec := pipeline.Spec{DotDagSource: s}
	vars := pipeline.NewVarsFrom(nil)

	_, trrs, err := r.ExecuteRun(testutils.Context(t), spec, vars)
	require.NoError(t, err)
	require.Len(t, trrs, len(d.Tasks))

	finalResults := trrs.FinalResult()

	val := finalResults.Values[0].(decimal.Decimal)

	assert.Equal(t, decimal.NewFromInt(1178718957397490305), val)
}

type taskRunWithVars struct {
	bridgeName        string
	ds2URL, ds4URL    string
	submitBridgeName  string
	includeInputAtKey string
}

func (t taskRunWithVars) String() string {
	return fmt.Sprintf(`
        ds1 [type=bridge name="%s" timeout=0 requestData=<{"data": $(foo)}>]
        ds1_parse [type=jsonparse lax=false  path="data,result" data="$(ds1)"]
        ds1_multiply [type=multiply input="$(ds1_parse.result)" times="$(ds1_parse.times)"]

        ds2 [type=http method="POST" url="%s" requestData=<{"data": [ $(bar), $(baz) ]}>]
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

        submit [type=bridge name="%s"
                includeInputAtKey="%s"
                requestData=<{
                    "median": $(median),
                    "fetchedValues": [ $(ds1_parse.result), $(ds2_parse.result) ],
                    "someString": $(ds4)
                }>]

        median -> submit;
        ds4 -> submit;
    `, t.bridgeName, t.ds2URL, t.ds4URL, t.submitBridgeName, t.includeInputAtKey)
}

func Test_PipelineRunner_ExecuteTaskRunsWithVars(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		vars              map[string]interface{}
		meta              map[string]interface{}
		includeInputAtKey string
	}{
		{
			name: "meta + includeInputAtKey",
			vars: map[string]interface{}{
				"foo": []interface{}{float64(123), "chainlink"},
				"bar": float64(123.45),
				"baz": "such oracle",
			},
			meta:              map[string]interface{}{"roundID": float64(456), "latestAnswer": float64(654)},
			includeInputAtKey: "sergey",
		},
		{
			name: "includeInputAtKey",
			vars: map[string]interface{}{
				"foo": *mustDecimal(t, "42.1337"),
				"bar": map[string]interface{}{"steve": "chainlink"},
				"baz": true,
			},
			includeInputAtKey: "best oracles",
		},
		{
			name: "meta",
			vars: map[string]interface{}{
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

			db := pgtest.NewSqlxDB(t)
			cfg := configtest.NewTestGeneralConfig(t)

			expectedRequestDS1 := map[string]interface{}{"data": test.vars["foo"]}
			expectedRequestDS2 := map[string]interface{}{"data": []interface{}{test.vars["bar"], test.vars["baz"]}}
			expectedRequestSubmit := map[string]interface{}{
				"median":        "9650000000000000000000",
				"fetchedValues": []interface{}{"9700", "9600"},
				"someString":    "some random string",
			}
			if test.meta != nil {
				expectedRequestDS1["meta"] = test.meta
				expectedRequestSubmit["meta"] = test.meta
				test.vars["jobRun"] = map[string]interface{}{"meta": test.meta}
			}
			if test.includeInputAtKey != "" {
				expectedRequestSubmit[test.includeInputAtKey] = "9650000000000000000000"
			}

			btORM := bridgesMocks.NewORM(t)

			// 1. Setup bridge
			ds1, bridge := makeBridge(t, db, expectedRequestDS1, map[string]interface{}{
				"data": map[string]interface{}{
					"result": map[string]interface{}{
						"result": decimal.NewFromInt(9700),
						"times":  "1000000000000000000",
					},
				},
			})
			defer ds1.Close()

			btORM.On("FindBridge", mock.Anything, bridge.Name).Return(bridge, nil).Once()

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
			submit, submitBt := makeBridge(t, db, expectedRequestSubmit, map[string]interface{}{"ok": true})
			defer submit.Close()

			btORM.On("FindBridge", mock.Anything, submitBt.Name).Return(submitBt, nil).Once()

			runner, _ := newRunner(t, db, btORM, cfg)
			specStr := taskRunWithVars{
				bridgeName:        bridge.Name.String(),
				ds2URL:            ds2.URL,
				ds4URL:            ds4.URL,
				submitBridgeName:  submitBt.Name.String(),
				includeInputAtKey: test.includeInputAtKey,
			}.String()
			p, err := pipeline.Parse(specStr)
			require.NoError(t, err)

			spec := pipeline.Spec{
				DotDagSource: specStr,
			}
			_, taskRunResults, err := runner.ExecuteRun(testutils.Context(t), spec, pipeline.NewVarsFrom(test.vars))
			require.NoError(t, err)
			require.Len(t, taskRunResults, len(p.Tasks))

			expectedResults := map[string]pipeline.Result{
				"ds1":          {Value: `{"data":{"result":{"result":"9700","times":"1000000000000000000"}}}` + "\n"},
				"ds1_parse":    {Value: map[string]interface{}{"result": "9700", "times": "1000000000000000000"}},
				"ds1_multiply": {Value: *mustDecimal(t, "9700000000000000000000")},
				"ds2":          {Value: `{"data":{"result":"9600","times":"1000000000000000000"}}` + "\n"},
				"ds2_parse":    {Value: map[string]interface{}{"result": "9600", "times": "1000000000000000000"}},
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
				if expected.Error != nil {
					require.Error(t, r.Result.Error)
					require.Contains(t, r.Result.Error.Error(), expected.Error.Error())
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

const (
	CBORDietEmpty = `
decode_log  [type="ethabidecodelog"
             data="$(jobRun.logData)"
             topics="$(jobRun.logTopics)"
             abi="OracleRequest(address requester, bytes32 requestId, uint256 payment, address callbackAddr, bytes4 callbackFunctionId, uint256 cancelExpiration, uint256 dataVersion, bytes cborPayload)"]

decode_cbor [type="cborparse"
             data="$(decode_log.cborPayload)"
			 mode=diet]

decode_log -> decode_cbor;
`
	CBORStdString = `
decode_log  [type="ethabidecodelog"
             data="$(jobRun.logData)"
             topics="$(jobRun.logTopics)"
             abi="OracleRequest(address requester, bytes32 requestId, uint256 payment, address callbackAddr, bytes4 callbackFunctionId, uint256 cancelExpiration, uint256 dataVersion, bytes cborPayload)"]

decode_cbor [type="cborparse"
             data="$(decode_log.cborPayload)"
			 mode=standard]

decode_log -> decode_cbor;
`
)

func Test_PipelineRunner_CBORParse(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	btORM := bridgesMocks.NewORM(t)
	r, _ := newRunner(t, db, btORM, cfg)

	t.Run("diet mode, empty CBOR", func(t *testing.T) {
		s := CBORDietEmpty
		d, err := pipeline.Parse(s)
		require.NoError(t, err)

		spec := pipeline.Spec{DotDagSource: s}
		global := make(map[string]interface{})
		jobRun := make(map[string]interface{})
		global["jobRun"] = jobRun
		jobRun["logData"] = hexutil.MustDecode("0x0000000000000000000000009c26cc46f57667cba75556014c8e0d5ed7c5b83d17a526ff5d8f916fa2f4a218f6ce0a6e410a0d7823f8238979f8579c2145fd6f0000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000009c26cc46f57667cba75556014c8e0d5ed7c5b83d64ef935700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006148ef28000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000")
		jobRun["logTopics"] = []common.Hash{
			common.HexToHash("0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65"),
			common.HexToHash("0x3963386131316165393962363463373161663333376235643831633737353230"),
		}
		vars := pipeline.NewVarsFrom(global)

		_, trrs, err := r.ExecuteRun(testutils.Context(t), spec, vars)
		require.NoError(t, err)
		require.Len(t, trrs, len(d.Tasks))

		finalResults := trrs.FinalResult()
		require.Len(t, finalResults.Values, 1)
		assert.Equal(t, make(map[string]interface{}), finalResults.Values[0])
		require.Len(t, finalResults.FatalErrors, 1)
		assert.Nil(t, finalResults.FatalErrors[0])
	})

	t.Run("standard mode, string value", func(t *testing.T) {
		s := CBORStdString
		d, err := pipeline.Parse(s)
		require.NoError(t, err)

		spec := pipeline.Spec{DotDagSource: s}
		global := make(map[string]interface{})
		jobRun := make(map[string]interface{})
		global["jobRun"] = jobRun
		jobRun["logData"] = hexutil.MustDecode("0x0000000000000000000000009c26cc46f57667cba75556014c8e0d5ed7c5b83d17a526ff5d8f916fa2f4a218f6ce0a6e410a0d7823f8238979f8579c2145fd6f0000000000000000000000000000000000000000000000000de0b6b3a76400000000000000000000000000009c26cc46f57667cba75556014c8e0d5ed7c5b83d64ef935700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006148ef2800000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000463666F6F00000000000000000000000000000000000000000000000000000000")
		jobRun["logTopics"] = []common.Hash{
			common.HexToHash("0xd8d7ecc4800d25fa53ce0372f13a416d98907a7ef3d8d3bdd79cf4fe75529c65"),
			common.HexToHash("0x3963386131316165393962363463373161663333376235643831633737353230"),
		}
		vars := pipeline.NewVarsFrom(global)

		_, trrs, err := r.ExecuteRun(testutils.Context(t), spec, vars)
		require.NoError(t, err)
		require.Len(t, trrs, len(d.Tasks))

		finalResults := trrs.FinalResult()
		require.Len(t, finalResults.Values, 1)
		assert.Equal(t, "foo", finalResults.Values[0])
		require.Len(t, finalResults.FatalErrors, 1)
		assert.Nil(t, finalResults.FatalErrors[0])
	})
}

func Test_PipelineRunner_HandleFaults(t *testing.T) {
	// We want to test the scenario where one or multiple APIs time out,
	// but a sufficient number of them still complete within the desired time frame
	// and so we can still obtain a median.
	db := pgtest.NewSqlxDB(t)

	m1 := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		time.Sleep(100 * time.Millisecond)
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(`{"result":10}`))
		assert.NoError(t, err)
	}))
	m2 := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		_, err := res.Write([]byte(`{"result":11}`))
		assert.NoError(t, err)
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
	cfg := configtest.NewTestGeneralConfig(t)
	btORM := bridgesMocks.NewORM(t)
	r, _ := newRunner(t, db, btORM, cfg)

	// If we cancel before an API is finished, we should still get a median.
	ctx, cancel := context.WithTimeout(testutils.Context(t), 50*time.Millisecond)
	defer cancel()

	spec := pipeline.Spec{DotDagSource: s}
	vars := pipeline.NewVarsFrom(nil)

	_, trrs, err := r.ExecuteRun(ctx, spec, vars)
	require.NoError(t, err)
	for _, trr := range trrs {
		if trr.IsTerminal() {
			require.Equal(t, decimal.RequireFromString("1100"), trr.Result.Value.(decimal.Decimal))
		}
	}
}

func Test_PipelineRunner_HandleFaultsPersistRun(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	orm := mocks.NewORM(t)
	btORM := bridgesMocks.NewORM(t)

	orm.On("InsertFinishedRun", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			args.Get(1).(*pipeline.Run).ID = 1
		}).
		Return(nil)
	cfg := configtest.NewTestGeneralConfig(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		ChainConfigs:   cfg.EVMConfigs(),
		DatabaseConfig: cfg.Database(),
		FeatureConfig:  cfg.Feature(),
		ListenerConfig: cfg.Database().Listener(),
		DB:             db,
		KeyStore:       ethKeyStore,
	})
	lggr := logger.TestLogger(t)
	r := pipeline.NewRunner(orm, btORM, cfg.JobPipeline(), cfg.WebServer(), legacyChains, ethKeyStore, nil, lggr, nil, nil)

	spec := pipeline.Spec{
		ID: 1,
		DotDagSource: `
fail_but_i_dont_care [type=fail]
succeed1             [type=memo value=10]
succeed2             [type=memo value=11]
final                [type=mean]

fail_but_i_dont_care -> final;
succeed1 -> final;
succeed2 -> final;
`}
	vars := pipeline.NewVarsFrom(nil)

	_, taskResults, err := r.ExecuteAndInsertFinishedRun(testutils.Context(t), spec, vars, false)
	finalResult := taskResults.FinalResult()
	require.NoError(t, err)
	assert.True(t, finalResult.HasErrors())
	assert.False(t, finalResult.HasFatalErrors())
	require.Len(t, finalResult.Values, 1)
	assert.Equal(t, "10.5", finalResult.Values[0].(decimal.Decimal).String())
}

func Test_PipelineRunner_ExecuteAndInsertFinishedRun_SavingTheSpec(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	orm := mocks.NewORM(t)
	btORM := bridgesMocks.NewORM(t)

	orm.On("InsertFinishedRunWithSpec", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			args.Get(1).(*pipeline.Run).ID = 1
		}).
		Return(nil)
	cfg := configtest.NewTestGeneralConfig(t)
	ethKeyStore := cltest.NewKeyStore(t, db).Eth()
	legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
		ChainConfigs:   cfg.EVMConfigs(),
		DatabaseConfig: cfg.Database(),
		FeatureConfig:  cfg.Feature(),
		ListenerConfig: cfg.Database().Listener(),
		DB:             db,
		KeyStore:       ethKeyStore,
	})
	lggr := logger.TestLogger(t)
	r := pipeline.NewRunner(orm, btORM, cfg.JobPipeline(), cfg.WebServer(), legacyChains, ethKeyStore, nil, lggr, nil, nil)

	spec := pipeline.Spec{
		DotDagSource: `
fail_but_i_dont_care [type=fail]
succeed1             [type=memo value=10]
succeed2             [type=memo value=11]
final                [type=mean]

fail_but_i_dont_care -> final;
succeed1 -> final;
succeed2 -> final;
`}
	vars := pipeline.NewVarsFrom(nil)

	_, taskResults, err := r.ExecuteAndInsertFinishedRun(testutils.Context(t), spec, vars, false)
	finalResult := taskResults.FinalResult()
	require.NoError(t, err)
	assert.True(t, finalResult.HasErrors())
	assert.False(t, finalResult.HasFatalErrors())
	require.Len(t, finalResult.Values, 1)
	assert.Equal(t, "10.5", finalResult.Values[0].(decimal.Decimal).String())
}

func Test_PipelineRunner_MultipleOutputs(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	btORM := bridgesMocks.NewORM(t)
	r, _ := newRunner(t, db, btORM, cfg)
	input := map[string]interface{}{"val": 2}
	_, trrs, err := r.ExecuteRun(testutils.Context(t), pipeline.Spec{
		DotDagSource: `
a [type=multiply input="$(val)" times=2]
b1 [type=multiply input="$(a)" times=2]
b2 [type=multiply input="$(a)" times=3]
c [type=median values=<[ $(b1), $(b2) ]> index=0]
a->b1->c;
a->b2->c;`,
	}, pipeline.NewVarsFrom(input))
	require.NoError(t, err)
	require.Equal(t, 4, len(trrs))
	assert.Equal(t, false, trrs.FinalResult().HasFatalErrors())

	// a = 4
	// (b1 = 8) + (b2 = 12)
	// c = 20 / 2

	result, err := trrs.FinalResult().SingularResult()
	require.NoError(t, err)
	assert.Equal(t, mustDecimal(t, "10").String(), result.Value.(decimal.Decimal).String())
}

func Test_PipelineRunner_MultipleTerminatingOutputs(t *testing.T) {
	cfg := configtest.NewTestGeneralConfig(t)
	btORM := bridgesMocks.NewORM(t)
	r, _ := newRunner(t, pgtest.NewSqlxDB(t), btORM, cfg)
	input := map[string]interface{}{"val": 2}
	_, trrs, err := r.ExecuteRun(testutils.Context(t), pipeline.Spec{
		DotDagSource: `
a [type=multiply input="$(val)" times=2]
b1 [type=multiply input="$(a)" times=2 index=0]
b2 [type=multiply input="$(a)" times=3 index=1]
a->b1;
a->b2;`,
	}, pipeline.NewVarsFrom(input))
	require.NoError(t, err)
	require.Equal(t, 3, len(trrs))
	result := trrs.FinalResult()
	assert.Equal(t, false, result.HasFatalErrors())

	assert.Equal(t, mustDecimal(t, "8").String(), result.Values[0].(decimal.Decimal).String())
	assert.Equal(t, mustDecimal(t, "12").String(), result.Values[1].(decimal.Decimal).String())
}

func Test_PipelineRunner_AsyncJob_Basic(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	btcUSDPairing := utils.MustUnmarshalToMap(`{"data":{"coin":"BTC","market":"USD"}}`)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody adapterRequest
		payload, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		defer r.Body.Close()
		err = json.Unmarshal(payload, &reqBody)
		require.NoError(t, err)
		// TODO: assert finding the id
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Chainlink-Pending", "true")
		response := map[string]interface{}{}
		require.NoError(t, json.NewEncoder(w).Encode(response))
	})

	// 1. Setup bridge
	s1 := httptest.NewServer(handler)
	defer s1.Close()

	bridgeFeedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)

	cfg := configtest.NewTestGeneralConfig(t)
	_, bt := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: bridgeFeedURL.String()})

	// 2. Setup success HTTP
	s2 := httptest.NewServer(fakePriceResponder(t, btcUSDPairing, decimal.NewFromInt(9600), "", nil))
	defer s2.Close()

	s4 := httptest.NewServer(fakeStringResponder(t, "foo-index-1"))
	defer s4.Close()
	s5 := httptest.NewServer(fakeStringResponder(t, "bar-index-2"))
	defer s5.Close()

	btORM := bridgesMocks.NewORM(t)
	btORM.On("FindBridge", mock.Anything, bt.Name).Return(*bt, nil)

	r, orm := newRunner(t, db, btORM, cfg)
	transactCall := orm.On("Transact", mock.Anything, mock.Anything)
	transactCall.Run(func(args mock.Arguments) {
		fn := args[1].(func(orm pipeline.ORM) error)
		transactCall.ReturnArguments = mock.Arguments{fn(orm)}
	})

	s := fmt.Sprintf(`
ds1 [type=bridge async=true name="%s" timeout=0 requestData=<{"data": {"coin": "BTC", "market": "USD"}}>]
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
`, bt.Name.String(), s2.URL, s4.URL, s5.URL)
	_, err = pipeline.Parse(s)
	require.NoError(t, err)

	spec := pipeline.Spec{DotDagSource: s}

	// Start a new run
	run := pipeline.NewRun(spec, pipeline.NewVarsFrom(nil))
	// we should receive a call to CreateRun because it's contains an async task
	orm.On("CreateRun", mock.Anything, mock.AnythingOfType("*pipeline.Run")).Return(nil).Run(func(args mock.Arguments) {
		run := args.Get(1).(*pipeline.Run)
		run.ID = 1 // give it a valid "id"
	}).Once()
	orm.On("StoreRun", mock.Anything, mock.AnythingOfType("*pipeline.Run")).Return(false, nil).Once()
	incomplete, err := r.Run(testutils.Context(t), run, false, nil)
	require.NoError(t, err)
	require.Len(t, run.PipelineTaskRuns, 9) // 3 tasks are suspended: ds1_parse, ds1_multiply, median. ds1 is present, but contains ErrPending
	require.Equal(t, true, incomplete)      // still incomplete

	// TODO: test a pending run that's not marked async=true, that is not allowed

	// Trigger run resumption with no new data
	orm.On("StoreRun", mock.Anything, mock.AnythingOfType("*pipeline.Run")).Return(false, nil).Once()
	incomplete, err = r.Run(testutils.Context(t), run, false, nil)
	require.NoError(t, err)
	require.Equal(t, true, incomplete) // still incomplete

	// Now simulate a new result coming in
	task := run.ByDotID("ds1")
	task.Error = null.NewString("", false)
	task.Output = jsonserializable.JSONSerializable{
		Val:   `{"data":{"result":"9700"}}` + "\n",
		Valid: true,
	}
	// Trigger run resumption
	orm.On("StoreRun", mock.Anything, mock.AnythingOfType("*pipeline.Run")).Return(false, nil).Once()
	incomplete, err = r.Run(testutils.Context(t), run, false, nil)
	require.NoError(t, err)
	require.Equal(t, false, incomplete) // done
	require.Len(t, run.PipelineTaskRuns, 12)
	require.Equal(t, false, incomplete) // run is complete

	require.Len(t, run.Outputs.Val, 3)
	require.Len(t, run.FatalErrors, 3)
	outputs := run.Outputs.Val.([]interface{})
	assert.Equal(t, "9650000000000000000000", outputs[0].(decimal.Decimal).String())
	assert.True(t, run.FatalErrors[0].IsZero())
	assert.Equal(t, "foo-index-1", outputs[1].(string))
	assert.True(t, run.FatalErrors[1].IsZero())
	assert.Equal(t, "bar-index-2", outputs[2].(string))
	assert.True(t, run.FatalErrors[2].IsZero())

	var errorResults []pipeline.TaskRun
	for _, trr := range run.PipelineTaskRuns {
		if trr.Result().Error != nil {
			errorResults = append(errorResults, trr)
		}
	}
	// There are three tasks in the erroring pipeline
	require.Len(t, errorResults, 3)
}

func Test_PipelineRunner_AsyncJob_InstantRestart(t *testing.T) {
	db := pgtest.NewSqlxDB(t)

	btcUSDPairing := utils.MustUnmarshalToMap(`{"data":{"coin":"BTC","market":"USD"}}`)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody adapterRequest
		payload, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		defer r.Body.Close()
		err = json.Unmarshal(payload, &reqBody)
		require.NoError(t, err)
		require.Contains(t, reqBody.ResponseURL, "http://localhost:6688/v2/resume/")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Chainlink-Pending", "true")
		response := map[string]interface{}{}
		require.NoError(t, json.NewEncoder(w).Encode(response))
	})

	// 1. Setup bridge
	s1 := httptest.NewServer(handler)
	defer s1.Close()

	bridgeFeedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)

	cfg := configtest.NewTestGeneralConfig(t)
	_, bt := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: bridgeFeedURL.String()})

	// 2. Setup success HTTP
	s2 := httptest.NewServer(fakePriceResponder(t, btcUSDPairing, decimal.NewFromInt(9600), "", nil))
	defer s2.Close()

	s4 := httptest.NewServer(fakeStringResponder(t, "foo-index-1"))
	defer s4.Close()
	s5 := httptest.NewServer(fakeStringResponder(t, "bar-index-2"))
	defer s5.Close()

	btORM := bridgesMocks.NewORM(t)
	btORM.On("FindBridge", mock.Anything, bt.Name).Return(*bt, nil)

	r, orm := newRunner(t, db, btORM, cfg)
	transactCall := orm.On("Transact", mock.Anything, mock.Anything)
	transactCall.Run(func(args mock.Arguments) {
		fn := args[1].(func(orm pipeline.ORM) error)
		transactCall.ReturnArguments = mock.Arguments{fn(orm)}
	})

	s := fmt.Sprintf(`
ds1 [type=bridge async=true name="%s" timeout=0 requestData=<{"data": {"coin": "BTC", "market": "USD"}}>]
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
`, bt.Name.String(), s2.URL, s4.URL, s5.URL)
	_, err = pipeline.Parse(s)
	require.NoError(t, err)

	spec := pipeline.Spec{DotDagSource: s}

	// Start a new run
	run := pipeline.NewRun(spec, pipeline.NewVarsFrom(nil))
	// we should receive a call to CreateRun because it's contains an async task
	orm.On("CreateRun", mock.Anything, mock.AnythingOfType("*pipeline.Run")).Return(nil).Run(func(args mock.Arguments) {
		run := args.Get(1).(*pipeline.Run)
		run.ID = 1 // give it a valid "id"
	}).Once()
	// Simulate updated task run data
	orm.On("StoreRun", mock.Anything, mock.AnythingOfType("*pipeline.Run")).Return(true, nil).Run(func(args mock.Arguments) {
		run := args.Get(1).(*pipeline.Run)
		// Now simulate a new result coming in while we were running
		task := run.ByDotID("ds1")
		task.Error = null.NewString("", false)
		task.Output = jsonserializable.JSONSerializable{
			Val:   `{"data":{"result":"9700"}}` + "\n",
			Valid: true,
		}
	}).Once()
	// StoreRun is called again to store the final result
	orm.On("StoreRun", mock.Anything, mock.AnythingOfType("*pipeline.Run")).Return(false, nil).Once()
	incomplete, err := r.Run(testutils.Context(t), run, false, nil)
	require.NoError(t, err)
	require.Len(t, run.PipelineTaskRuns, 12)
	require.Equal(t, false, incomplete) // run is complete

	require.Len(t, run.Outputs.Val, 3)
	require.Len(t, run.FatalErrors, 3)
	outputs := run.Outputs.Val.([]interface{})
	assert.Equal(t, "9650000000000000000000", outputs[0].(decimal.Decimal).String())
	assert.True(t, run.FatalErrors[0].IsZero())
	assert.Equal(t, "foo-index-1", outputs[1].(string))
	assert.True(t, run.FatalErrors[1].IsZero())
	assert.Equal(t, "bar-index-2", outputs[2].(string))
	assert.True(t, run.FatalErrors[2].IsZero())

	var errorResults []pipeline.TaskRun
	for _, trr := range run.PipelineTaskRuns {
		if trr.Result().Error != nil {
			errorResults = append(errorResults, trr)
		}
	}
	// There are three tasks in the erroring pipeline
	require.Len(t, errorResults, 3)
}

func Test_PipelineRunner_LowercaseOutputs(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	btORM := bridgesMocks.NewORM(t)
	r, _ := newRunner(t, db, btORM, cfg)
	input := map[string]interface{}{
		"first":  "camelCase",
		"second": "UPPERCASE",
	}
	_, trrs, err := r.ExecuteRun(testutils.Context(t), pipeline.Spec{
		DotDagSource: `
a [type=lowercase input="$(first)"]
`,
	}, pipeline.NewVarsFrom(input))
	require.NoError(t, err)
	require.Equal(t, 1, len(trrs))
	assert.Equal(t, false, trrs.FinalResult().HasFatalErrors())

	result, err := trrs.FinalResult().SingularResult()
	require.NoError(t, err)
	assert.Equal(t, "camelcase", result.Value.(string))
}

func Test_PipelineRunner_UppercaseOutputs(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	btORM := bridgesMocks.NewORM(t)
	r, _ := newRunner(t, db, btORM, cfg)
	input := map[string]interface{}{
		"first": "somerAnDomTEST",
	}
	_, trrs, err := r.ExecuteRun(testutils.Context(t), pipeline.Spec{
		DotDagSource: `
a [type=uppercase input="$(first)"]
`,
	}, pipeline.NewVarsFrom(input))
	require.NoError(t, err)
	require.Equal(t, 1, len(trrs))
	assert.Equal(t, false, trrs.FinalResult().HasFatalErrors())

	result, err := trrs.FinalResult().SingularResult()
	require.NoError(t, err)
	assert.Equal(t, "SOMERANDOMTEST", result.Value.(string))
}

func Test_PipelineRunner_HexDecodeOutputs(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	btORM := bridgesMocks.NewORM(t)
	r, _ := newRunner(t, db, btORM, cfg)
	input := map[string]interface{}{
		"astring": "0x12345678",
	}
	_, trrs, err := r.ExecuteRun(testutils.Context(t), pipeline.Spec{
		DotDagSource: `
a [type=hexdecode input="$(astring)"]
`,
	}, pipeline.NewVarsFrom(input))
	require.NoError(t, err)
	require.Equal(t, 1, len(trrs))
	assert.Equal(t, false, trrs.FinalResult().HasFatalErrors())

	result, err := trrs.FinalResult().SingularResult()
	require.NoError(t, err)
	assert.Equal(t, []byte{0x12, 0x34, 0x56, 0x78}, result.Value)
}

func Test_PipelineRunner_HexEncodeAndDecode(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	btORM := bridgesMocks.NewORM(t)
	r, _ := newRunner(t, db, btORM, cfg)
	inputBytes := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit")
	input := map[string]interface{}{
		"input_val": inputBytes,
	}
	_, trrs, err := r.ExecuteRun(testutils.Context(t), pipeline.Spec{
		DotDagSource: `
en [type=hexencode input="$(input_val)"]
de [type=hexdecode]
en->de
`,
	}, pipeline.NewVarsFrom(input))
	require.NoError(t, err)
	require.Equal(t, 2, len(trrs))
	assert.Equal(t, false, trrs.FinalResult().HasFatalErrors())

	result, err := trrs.FinalResult().SingularResult()
	require.NoError(t, err)
	assert.Equal(t, inputBytes, result.Value)
}

func Test_PipelineRunner_Base64DecodeOutputs(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	btORM := bridgesMocks.NewORM(t)
	r, _ := newRunner(t, db, btORM, cfg)
	input := map[string]interface{}{
		"astring": "SGVsbG8sIHBsYXlncm91bmQ=",
	}
	_, trrs, err := r.ExecuteRun(testutils.Context(t), pipeline.Spec{
		DotDagSource: `
a [type=base64decode input="$(astring)"]
`,
	}, pipeline.NewVarsFrom(input))
	require.NoError(t, err)
	require.Equal(t, 1, len(trrs))
	assert.Equal(t, false, trrs.FinalResult().HasFatalErrors())

	result, err := trrs.FinalResult().SingularResult()
	require.NoError(t, err)
	assert.Equal(t, []byte("Hello, playground"), result.Value)
}

func Test_PipelineRunner_Base64EncodeAndDecode(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	btORM := bridgesMocks.NewORM(t)
	r, _ := newRunner(t, db, btORM, cfg)
	inputBytes := []byte("[{\"add\": \"weather\", \"during\": true}, 1478647067]")
	input := map[string]interface{}{
		"input_val": inputBytes,
	}
	_, trrs, err := r.ExecuteRun(testutils.Context(t), pipeline.Spec{
		DotDagSource: `
en [type=base64encode input="$(input_val)"]
de [type=base64decode]
en->de
`,
	}, pipeline.NewVarsFrom(input))
	require.NoError(t, err)
	require.Equal(t, 2, len(trrs))
	assert.Equal(t, false, trrs.FinalResult().HasFatalErrors())

	result, err := trrs.FinalResult().SingularResult()
	require.NoError(t, err)
	assert.Equal(t, inputBytes, result.Value)
}

func Test_PipelineRunner_ExecuteRun(t *testing.T) {
	t.Run("uses cached *Pipeline if available", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := configtest.NewTestGeneralConfig(t)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		legacyChains := evmtest.NewLegacyChains(t, evmtest.TestChainOpts{
			ChainConfigs:   cfg.EVMConfigs(),
			DatabaseConfig: cfg.Database(),
			FeatureConfig:  cfg.Feature(),
			ListenerConfig: cfg.Database().Listener(),
			DB:             db,
			KeyStore:       ethKeyStore,
		})
		lggr := logger.TestLogger(t)
		r := pipeline.NewRunner(nil, nil, cfg.JobPipeline(), cfg.WebServer(), legacyChains, ethKeyStore, nil, lggr, nil, nil)

		template := `
succeed             [type=memo value=%d]
succeed;
`

		spec := pipeline.Spec{DotDagSource: fmt.Sprintf(template, 1)}
		vars := pipeline.NewVarsFrom(nil)

		_, trrs, err := r.ExecuteRun(testutils.Context(t), spec, vars)
		require.NoError(t, err)
		require.Len(t, trrs, 1)
		assert.Equal(t, "1", trrs[0].Result.Value.(pipeline.ObjectParam).DecimalValue.Decimal().String())

		// does not automatically cache
		require.Nil(t, spec.Pipeline)

		// initialize it
		spec.Pipeline, err = spec.ParsePipeline()
		require.NoError(t, err)

		// even though this is set to 2, it should use the cached version
		spec.DotDagSource = fmt.Sprintf(template, 2)

		_, trrs, err = r.ExecuteRun(testutils.Context(t), spec, vars)
		require.NoError(t, err)
		require.Len(t, trrs, 1)
		assert.Equal(t, "1", trrs[0].Result.Value.(pipeline.ObjectParam).DecimalValue.Decimal().String())
	})
}
