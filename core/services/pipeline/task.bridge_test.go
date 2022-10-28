package pipeline_test

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	configtest "github.com/smartcontractkit/chainlink/core/internal/testutils/configtest/v2"
	clhttptest "github.com/smartcontractkit/chainlink/core/internal/testutils/httptest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// ethUSDPairing has the ETH/USD parameters needed when POSTing to the price
// external adapters.
// https://github.com/smartcontractkit/price-adapters

var (
	btcUSDPairing = `{"data":{"coin":"BTC","market":"USD"}}`
	ethUSDPairing = `{"data":{"coin":"ETH","market":"USD"}}`
)

type adapterRequest struct {
	ID          string            `json:"id"`
	Data        pipeline.MapParam `json:"data"`
	Meta        pipeline.MapParam `json:"meta"`
	ResponseURL string            `json:"responseURL"`
}

type adapterResponseData struct {
	Result *decimal.Decimal `json:"result"`
}

// adapterResponse is the HTTP response as defined by the external adapter:
// https://github.com/smartcontractkit/bnc-adapter
type adapterResponse struct {
	Data         adapterResponseData `json:"data"`
	ErrorMessage null.String         `json:"errorMessage"`
}

func (pr adapterResponse) Result() *decimal.Decimal {
	return pr.Data.Result
}

func dataWithResult(t *testing.T, result decimal.Decimal) adapterResponseData {
	t.Helper()
	var data adapterResponseData
	body := []byte(fmt.Sprintf(`{"result":%v}`, result))
	require.NoError(t, json.Unmarshal(body, &data))
	return data
}

func mustReadFile(t testing.TB, file string) string {
	t.Helper()

	content, err := os.ReadFile(file)
	require.NoError(t, err)
	return string(content)
}

func fakePriceResponder(t *testing.T, requestData map[string]interface{}, result decimal.Decimal, inputKey string, expectedInput interface{}) http.Handler {
	t.Helper()

	body, err := json.Marshal(requestData)
	require.NoError(t, err)
	var expectedRequest adapterRequest
	err = json.Unmarshal(body, &expectedRequest)
	require.NoError(t, err)
	response := adapterResponse{Data: dataWithResult(t, result)}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody adapterRequest
		payload, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		defer r.Body.Close()
		err = json.Unmarshal(payload, &reqBody)
		require.NoError(t, err)
		require.Equal(t, expectedRequest.Data, reqBody.Data)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(response))

		if inputKey != "" {
			m := utils.MustUnmarshalToMap(string(payload))
			if expectedInput != nil {
				require.Equal(t, expectedInput, m[inputKey])
			} else {
				require.Nil(t, m[inputKey])
			}
		}
	})
}

func fakeStringResponder(t *testing.T, s string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(s))
		require.NoError(t, err)
	})
}

func TestBridgeTask_Happy(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)

	s1 := httptest.NewServer(fakePriceResponder(t, utils.MustUnmarshalToMap(btcUSDPairing), decimal.NewFromInt(9700), "", nil))
	defer s1.Close()

	feedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)

	orm := bridges.NewORM(db, logger.TestLogger(t), cfg)
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()}, cfg)

	task := pipeline.BridgeTask{
		BaseTask:    pipeline.NewBaseTask(0, "bridge", nil, nil, 0),
		Name:        bridge.Name.String(),
		RequestData: btcUSDPairing,
	}
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	task.HelperSetDependencies(cfg, orm, uuid.UUID{}, c)

	result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), nil)
	assert.False(t, runInfo.IsPending)
	assert.False(t, runInfo.IsRetryable)
	require.NoError(t, result.Error)
	require.NotNil(t, result.Value)
	var x struct {
		Data struct {
			Result decimal.Decimal `json:"result"`
		} `json:"data"`
	}
	err = json.Unmarshal([]byte(result.Value.(string)), &x)
	require.NoError(t, err)
	require.Equal(t, decimal.NewFromInt(9700), x.Data.Result)
}

func TestBridgeTask_AsyncJobPendingState(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)

	id := uuid.NewV4()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody adapterRequest
		payload, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		defer r.Body.Close()

		err = json.Unmarshal(payload, &reqBody)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%s/v2/resume/%v", cfg.BridgeResponseURL(), id.String()), reqBody.ResponseURL)
		w.Header().Set("Content-Type", "application/json")

		// w.Header().Set("X-Chainlink-Pending", "true")
		response := map[string]interface{}{"pending": true}
		require.NoError(t, json.NewEncoder(w).Encode(response))

	})

	server := httptest.NewServer(handler)
	defer server.Close()
	feedURL, err := url.ParseRequestURI(server.URL)
	require.NoError(t, err)

	orm := bridges.NewORM(db, logger.TestLogger(t), cfg)
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()}, cfg)

	task := pipeline.BridgeTask{
		Name:        bridge.Name.String(),
		RequestData: ethUSDPairing,
		Async:       "true",
	}
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	task.HelperSetDependencies(cfg, orm, id, c)

	result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), nil)
	assert.True(t, runInfo.IsPending)
	assert.False(t, runInfo.IsRetryable)

	require.NoError(t, result.Error)
	require.Nil(t, result.Value)
}

func TestBridgeTask_Variables(t *testing.T) {
	t.Parallel()

	validMeta := map[string]interface{}{"theMeta": "yes"}

	tests := []struct {
		name                  string
		requestData           string
		includeInputAtKey     string
		inputs                []pipeline.Result
		vars                  pipeline.Vars
		expectedRequestData   map[string]interface{}
		expectedErrorCause    error
		expectedErrorContains string
	}{
		{
			"requestData (empty) + includeInputAtKey + meta",
			``,
			"input",
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"some_data": map[string]interface{}{"foo": 543.21}}),
			map[string]interface{}{
				"input": 123.45,
				"meta":  validMeta,
			},
			nil,
			"",
		},
		{
			"requestData (pure variable) + includeInputAtKey + meta",
			`$(some_data)`,
			"input",
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"some_data": map[string]interface{}{"foo": 543.21}}),
			map[string]interface{}{
				"foo":   543.21,
				"input": 123.45,
				"meta":  validMeta,
			},
			nil,
			"",
		},
		{
			"requestData (pure variable) + includeInputAtKey",
			`$(some_data)`,
			"input",
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"some_data": map[string]interface{}{"foo": 543.21}}),
			map[string]interface{}{
				"foo":   543.21,
				"input": 123.45,
			},
			nil,
			"",
		},
		{
			"requestData (pure variable) + meta",
			`$(some_data)`,
			"",
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"some_data": map[string]interface{}{"foo": 543.21}}),
			map[string]interface{}{
				"foo":  543.21,
				"meta": validMeta,
			},
			nil,
			"",
		},
		{
			"requestData (pure variable, missing)",
			`$(some_data)`,
			"input",
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"not_some_data": map[string]interface{}{"foo": 543.21}}),
			nil,
			pipeline.ErrKeypathNotFound,
			"requestData",
		},
		{
			"requestData (pure variable, not a map)",
			`$(some_data)`,
			"input",
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"some_data": 543.21}),
			nil,
			pipeline.ErrBadInput,
			"requestData",
		},
		{
			"requestData (interpolation) + includeInputAtKey + meta",
			`{"data":{"result":$(medianize)}}`,
			"input",
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"medianize": 543.21}),
			map[string]interface{}{
				"data":  map[string]interface{}{"result": 543.21},
				"input": 123.45,
				"meta":  validMeta,
			},
			nil,
			"",
		},
		{
			"requestData (interpolation) + includeInputAtKey",
			`{"data":{"result":$(medianize)}}`,
			"input",
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"medianize": 543.21}),
			map[string]interface{}{
				"data":  map[string]interface{}{"result": 543.21},
				"input": 123.45,
			},
			nil,
			"",
		},
		{
			"requestData (interpolation) + meta",
			`{"data":{"result":$(medianize)}}`,
			"",
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"medianize": 543.21}),
			map[string]interface{}{
				"data": map[string]interface{}{"result": 543.21},
				"meta": validMeta,
			},
			nil,
			"",
		},
		{
			"requestData (interpolation, missing)",
			`{"data":{"result":$(medianize)}}`,
			"input",
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"nope": "foo bar"}),
			nil,
			pipeline.ErrKeypathNotFound,
			"requestData",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			db := pgtest.NewSqlxDB(t)
			cfg := configtest.NewTestGeneralConfig(t)

			s1 := httptest.NewServer(fakePriceResponder(t, test.expectedRequestData, decimal.NewFromInt(9700), "", nil))
			defer s1.Close()

			feedURL, err := url.ParseRequestURI(s1.URL)
			require.NoError(t, err)

			orm := bridges.NewORM(db, logger.TestLogger(t), cfg)
			_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()}, cfg)

			task := pipeline.BridgeTask{
				BaseTask:          pipeline.NewBaseTask(0, "bridge", nil, nil, 0),
				Name:              bridge.Name.String(),
				RequestData:       test.requestData,
				IncludeInputAtKey: test.includeInputAtKey,
			}
			c := clhttptest.NewTestLocalOnlyHTTPClient()
			task.HelperSetDependencies(cfg, orm, uuid.UUID{}, c)

			result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), test.vars, test.inputs)
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)
			if test.expectedErrorCause != nil {
				require.Equal(t, test.expectedErrorCause, errors.Cause(result.Error))
				if test.expectedErrorContains != "" {
					require.Contains(t, result.Error.Error(), test.expectedErrorContains)
				}

			} else {
				require.NoError(t, result.Error)
				require.NotNil(t, result.Value)
				var x struct {
					Data struct {
						Result decimal.Decimal `json:"result"`
					} `json:"data"`
				}
				err = json.Unmarshal([]byte(result.Value.(string)), &x)
				require.NoError(t, err)
				require.Equal(t, decimal.NewFromInt(9700), x.Data.Result)
			}
		})
	}
}

func TestBridgeTask_Meta(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)

	var empty adapterResponse

	var httpCalled atomic.Bool
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req adapterRequest
		body, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(body, &req)
		require.NoError(t, err)
		require.Equal(t, float64(10), req.Meta["latestAnswer"])
		require.Equal(t, float64(1616447984), req.Meta["updatedAt"])
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(empty))
		httpCalled.Store(true)
	})

	metaDataForBridge, err := bridges.MarshalBridgeMetaData(big.NewInt(10), big.NewInt(1616447984))
	require.NoError(t, err)

	s1 := httptest.NewServer(handler)

	defer s1.Close()
	feedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)

	orm := bridges.NewORM(db, logger.TestLogger(t), cfg)
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()}, cfg)

	task := pipeline.BridgeTask{
		BaseTask:    pipeline.NewBaseTask(0, "bridge", nil, nil, 0),
		RequestData: ethUSDPairing,
		Name:        bridge.Name.String(),
	}
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	task.HelperSetDependencies(cfg, orm, uuid.UUID{}, c)

	mp := map[string]interface{}{"meta": metaDataForBridge}
	res, _ := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(map[string]interface{}{"jobRun": mp}), nil)
	assert.Nil(t, res.Error)

	assert.True(t, httpCalled.Load())
}

func TestBridgeTask_IncludeInputAtKey(t *testing.T) {
	t.Parallel()

	theErr := errors.New("foo")

	tests := []struct {
		name               string
		inputs             []pipeline.Result
		includeInputAtKey  string
		expectedInput      interface{}
		expectedErrorCause error
	}{
		{"no input, no includeInputAtKey", nil, "", nil, nil},
		{"no input, includeInputAtKey", nil, "result", nil, nil},
		{"input, no includeInputAtKey", []pipeline.Result{{Value: decimal.NewFromFloat(123.45)}}, "", nil, nil},
		{"input, includeInputAtKey", []pipeline.Result{{Value: decimal.NewFromFloat(123.45)}}, "result", "123.45", nil},
		{"input has error", []pipeline.Result{{Error: theErr}}, "result", nil, pipeline.ErrTooManyErrors},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			db := pgtest.NewSqlxDB(t)
			cfg := configtest.NewTestGeneralConfig(t)

			s1 := httptest.NewServer(fakePriceResponder(t, utils.MustUnmarshalToMap(btcUSDPairing), decimal.NewFromInt(9700), test.includeInputAtKey, test.expectedInput))
			defer s1.Close()

			feedURL, err := url.ParseRequestURI(s1.URL)
			require.NoError(t, err)

			orm := bridges.NewORM(db, logger.TestLogger(t), cfg)
			_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()}, cfg)

			task := pipeline.BridgeTask{
				BaseTask:          pipeline.NewBaseTask(0, "bridge", nil, nil, 0),
				Name:              bridge.Name.String(),
				RequestData:       btcUSDPairing,
				IncludeInputAtKey: test.includeInputAtKey,
			}
			c := clhttptest.NewTestLocalOnlyHTTPClient()
			task.HelperSetDependencies(cfg, orm, uuid.UUID{}, c)

			result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), test.inputs)
			assert.False(t, runInfo.IsPending)
			assert.False(t, runInfo.IsRetryable)
			if test.expectedErrorCause != nil {
				require.Equal(t, test.expectedErrorCause, errors.Cause(result.Error))
				require.Nil(t, result.Value)
			} else {
				require.NoError(t, result.Error)
				require.NotNil(t, result.Value)
				var x struct {
					Data struct {
						Result decimal.Decimal `json:"result"`
					} `json:"data"`
				}
				err = json.Unmarshal([]byte(result.Value.(string)), &x)
				require.NoError(t, err)
				require.Equal(t, decimal.NewFromInt(9700), x.Data.Result)
			}
		})
	}
}

func TestBridgeTask_ErrorMessage(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		err := json.NewEncoder(w).Encode(adapterResponse{
			ErrorMessage: null.StringFrom("could not hit data fetcher"),
		})
		require.NoError(t, err)
	})

	server := httptest.NewServer(handler)
	defer server.Close()
	feedURL, err := url.ParseRequestURI(server.URL)
	require.NoError(t, err)

	orm := bridges.NewORM(db, logger.TestLogger(t), cfg)
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()}, cfg)

	task := pipeline.BridgeTask{
		Name:        bridge.Name.String(),
		RequestData: ethUSDPairing,
	}
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	task.HelperSetDependencies(cfg, orm, uuid.UUID{}, c)

	result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), nil)
	assert.False(t, runInfo.IsPending)
	assert.False(t, runInfo.IsRetryable)
	require.Error(t, result.Error)
	require.Contains(t, result.Error.Error(), "could not hit data fetcher")
	require.Nil(t, result.Value)
}

func TestBridgeTask_OnlyErrorMessage(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, err := w.Write([]byte(mustReadFile(t, "../../testdata/apiresponses/coinmarketcap.error.json")))
		require.NoError(t, err)
	})

	server := httptest.NewServer(handler)
	defer server.Close()
	feedURL, err := url.ParseRequestURI(server.URL)
	require.NoError(t, err)

	orm := bridges.NewORM(db, logger.TestLogger(t), cfg)
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()}, cfg)

	task := pipeline.BridgeTask{
		Name:        bridge.Name.String(),
		RequestData: ethUSDPairing,
	}
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	task.HelperSetDependencies(cfg, orm, uuid.UUID{}, c)

	result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), nil)
	assert.False(t, runInfo.IsPending)
	assert.True(t, runInfo.IsRetryable)
	require.Error(t, result.Error)
	require.Contains(t, result.Error.Error(), "RequestId")
	require.Nil(t, result.Value)
}

func TestBridgeTask_ErrorIfBridgeMissing(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)

	task := pipeline.BridgeTask{
		Name:        "foo",
		RequestData: btcUSDPairing,
	}
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	orm := bridges.NewORM(db, logger.TestLogger(t), cfg)
	task.HelperSetDependencies(cfg, orm, uuid.UUID{}, c)

	result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), nil)
	assert.False(t, runInfo.IsPending)
	assert.False(t, runInfo.IsRetryable)
	require.Nil(t, result.Value)
	require.Error(t, result.Error)
	assert.Contains(t, result.Error.Error(), "could not find bridge with name 'foo'")
}

// Sample input taken from
// https://github.com/smartcontractkit/price-adapters#chainlink-price-request-adapters
func TestAdapterResponse_UnmarshalJSON_Happy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name, content string
		expect        decimal.Decimal
	}{
		{"basic", `{"data":{"result":123.4567890},"jobRunID":"1","statusCode":200}`, decimal.NewFromFloat(123.456789)},
		{"bravenewcoin", mustReadFile(t, "../../testdata/apiresponses/bravenewcoin.json"), decimal.NewFromFloat(306.52036004)},
		{"coinmarketcap", mustReadFile(t, "../../testdata/apiresponses/coinmarketcap.json"), decimal.NewFromFloat(305.5574615)},
		{"cryptocompare", mustReadFile(t, "../../testdata/apiresponses/cryptocompare.json"), decimal.NewFromFloat(305.76)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var response adapterResponse
			err := json.Unmarshal([]byte(test.content), &response)
			require.NoError(t, err)
			result := response.Result()
			require.Equal(t, test.expect.String(), result.String())
		})
	}
}
