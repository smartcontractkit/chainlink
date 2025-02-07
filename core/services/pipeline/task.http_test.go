package pipeline_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/jsonserializable"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	clhttptest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/httptest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	clhttp "github.com/smartcontractkit/chainlink/v2/core/utils/http"
)

// ethUSDPairing has the ETH/USD parameters needed when POSTing to the price
// external adapters.
// https://github.com/smartcontractkit/price-adapters

func TestHTTPTask_Happy(t *testing.T) {
	t.Parallel()

	config := configtest.NewTestGeneralConfig(t)
	s1 := httptest.NewServer(fakePriceResponder(t, utils.MustUnmarshalToMap(btcUSDPairing), decimal.NewFromInt(9700), "", nil))
	defer s1.Close()

	task := pipeline.HTTPTask{
		BaseTask:    pipeline.NewBaseTask(0, "http", nil, nil, 0),
		Method:      "POST",
		URL:         s1.URL,
		RequestData: btcUSDPairing,
	}
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	task.HelperSetDependencies(config.JobPipeline(), c, c)

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
	err := json.Unmarshal([]byte(result.Value.(string)), &x)
	require.NoError(t, err)
	require.Equal(t, decimal.NewFromInt(9700), x.Data.Result)
}

func TestHTTPTask_Variables(t *testing.T) {
	t.Parallel()

	validMeta := map[string]interface{}{"theMeta": "yes"}

	tests := []struct {
		name                  string
		requestData           string
		meta                  jsonserializable.JSONSerializable
		inputs                []pipeline.Result
		vars                  pipeline.Vars
		expectedRequestData   map[string]interface{}
		expectedErrorCause    error
		expectedErrorContains string
	}{
		{
			"requestData (empty) + meta",
			``,
			jsonserializable.JSONSerializable{Val: validMeta, Valid: true},
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"some_data": map[string]interface{}{"foo": 543.21}}),
			map[string]interface{}{},
			nil,
			"",
		},
		{
			"requestData (pure variable) + meta",
			`$(some_data)`,
			jsonserializable.JSONSerializable{Val: validMeta, Valid: true},
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"some_data": map[string]interface{}{"foo": 543.21}}),
			map[string]interface{}{"foo": 543.21},
			nil,
			"",
		},
		{
			"requestData (pure variable)",
			`$(some_data)`,
			jsonserializable.JSONSerializable{Val: nil, Valid: false},
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"some_data": map[string]interface{}{"foo": 543.21}}),
			map[string]interface{}{"foo": 543.21},
			nil,
			"",
		},
		{
			"requestData (pure variable, missing)",
			`$(some_data)`,
			jsonserializable.JSONSerializable{Val: validMeta, Valid: true},
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"not_some_data": map[string]interface{}{"foo": 543.21}}),
			nil,
			pipeline.ErrKeypathNotFound,
			"requestData",
		},
		{
			"requestData (pure variable, not a map)",
			`$(some_data)`,
			jsonserializable.JSONSerializable{Val: validMeta, Valid: true},
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"some_data": 543.21}),
			nil,
			pipeline.ErrBadInput,
			"requestData",
		},
		{
			"requestData (interpolation) + meta",
			`{"data":{"result":$(medianize)}}`,
			jsonserializable.JSONSerializable{Val: validMeta, Valid: true},
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"medianize": 543.21}),
			map[string]interface{}{"data": map[string]interface{}{"result": 543.21}},
			nil,
			"",
		},
		{
			"requestData (interpolation, missing)",
			`{"data":{"result":$(medianize)}}`,
			jsonserializable.JSONSerializable{Val: validMeta, Valid: true},
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

			orm := bridges.NewORM(db)
			_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()})

			task := pipeline.BridgeTask{
				BaseTask:    pipeline.NewBaseTask(0, "bridge", nil, nil, 0),
				Name:        bridge.Name.String(),
				RequestData: test.requestData,
			}
			c := clhttptest.NewTestLocalOnlyHTTPClient()
			trORM := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
			specID, err := trORM.CreateSpec(testutils.Context(t), pipeline.Pipeline{}, *models.NewInterval(5 * time.Minute))
			require.NoError(t, err)
			task.HelperSetDependencies(cfg.JobPipeline(), cfg.WebServer(), orm, specID, uuid.UUID{}, c)

			err = test.vars.Set("meta", test.meta)
			require.NoError(t, err)

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
				err := json.Unmarshal([]byte(result.Value.(string)), &x)
				require.NoError(t, err)
				require.Equal(t, decimal.NewFromInt(9700), x.Data.Result)
			}
		})
	}
}

func TestHTTPTask_OverrideURLSafe(t *testing.T) {
	t.Parallel()

	config := configtest.NewTestGeneralConfig(t)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("{}"))
		require.NoError(t, err)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	task := pipeline.HTTPTask{
		Method:      "POST",
		URL:         server.URL,
		RequestData: ethUSDPairing,
	}
	// Use real clients here to actually test the local connection blocking
	r := clhttp.NewRestrictedHTTPClient(config.Database(), logger.TestLogger(t))
	u := clhttp.NewUnrestrictedHTTPClient()
	task.HelperSetDependencies(config.JobPipeline(), r, u)

	result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), nil)
	assert.False(t, runInfo.IsPending)
	assert.False(t, runInfo.IsRetryable)
	require.NoError(t, result.Error)

	task.URL = "$(url)"

	vars := pipeline.NewVarsFrom(map[string]interface{}{"url": server.URL})
	result, runInfo = task.Run(testutils.Context(t), logger.TestLogger(t), vars, nil)
	assert.False(t, runInfo.IsPending)
	assert.True(t, runInfo.IsRetryable)
	require.Error(t, result.Error)
	require.Contains(t, result.Error.Error(), "Connections to local/private and multicast networks are disabled")
	require.Nil(t, result.Value)

	task.AllowUnrestrictedNetworkAccess = "true"

	result, runInfo = task.Run(testutils.Context(t), logger.TestLogger(t), vars, nil)
	assert.False(t, runInfo.IsPending)
	assert.False(t, runInfo.IsRetryable)
	require.NoError(t, result.Error)
}

func TestHTTPTask_ErrorMessage(t *testing.T) {
	t.Parallel()

	config := configtest.NewTestGeneralConfig(t)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		resp := &adapterResponse{}
		resp.SetErrorMessage("could not hit data fetcher")
		err := json.NewEncoder(w).Encode(resp)
		require.NoError(t, err)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	c := clhttptest.NewTestLocalOnlyHTTPClient()
	task := pipeline.HTTPTask{
		Method:      "POST",
		URL:         server.URL,
		RequestData: ethUSDPairing,
	}
	task.HelperSetDependencies(config.JobPipeline(), c, c)

	result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), nil)
	assert.False(t, runInfo.IsPending)
	assert.False(t, runInfo.IsRetryable)

	require.Error(t, result.Error)
	require.Contains(t, result.Error.Error(), "could not hit data fetcher")
	require.Nil(t, result.Value)
}

func TestHTTPTask_OnlyErrorMessage(t *testing.T) {
	t.Parallel()

	config := configtest.NewTestGeneralConfig(t)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, err := w.Write([]byte(mustReadFile(t, "../../testdata/apiresponses/coinmarketcap.error.json")))
		require.NoError(t, err)
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	task := pipeline.HTTPTask{
		Method:      "POST",
		URL:         server.URL,
		RequestData: ethUSDPairing,
	}
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	task.HelperSetDependencies(config.JobPipeline(), c, c)

	result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), nil)
	assert.False(t, runInfo.IsPending)
	assert.True(t, runInfo.IsRetryable)
	require.Error(t, result.Error)
	require.Contains(t, result.Error.Error(), "RequestId")
	require.Nil(t, result.Value)
}

func TestHTTPTask_Headers(t *testing.T) {
	allHeaders := func(headers http.Header) (s []string) {
		var keys []string
		for k := range headers {
			keys = append(keys, k)
		}
		// get it in a consistent order
		sort.Strings(keys)
		for _, k := range keys {
			v := headers.Get(k)
			s = append(s, k, v)
		}
		return s
	}

	standardHeaders := []string{"Content-Length", "38", "Content-Type", "application/json", "User-Agent", "Go-http-client/1.1"}

	t.Run("sends headers", func(t *testing.T) {
		config := configtest.NewTestGeneralConfig(t)
		var headers http.Header
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers = r.Header
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"fooresponse": 1}`))
			require.NoError(t, err)
		})

		server := httptest.NewServer(handler)
		defer server.Close()

		task := pipeline.HTTPTask{
			Method:      "POST",
			URL:         server.URL,
			RequestData: ethUSDPairing,
			Headers:     `["X-Header-1", "foo", "X-Header-2", "bar"]`,
		}
		c := clhttptest.NewTestLocalOnlyHTTPClient()
		task.HelperSetDependencies(config.JobPipeline(), c, c)

		result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), nil)
		assert.False(t, runInfo.IsPending)
		assert.Equal(t, `{"fooresponse": 1}`, result.Value)
		assert.NoError(t, result.Error)

		assert.Equal(t, append(standardHeaders, "X-Header-1", "foo", "X-Header-2", "bar"), allHeaders(headers))
	})

	t.Run("errors with odd number of headers", func(t *testing.T) {
		task := pipeline.HTTPTask{
			Method:      "POST",
			URL:         "http://example.com",
			RequestData: ethUSDPairing,
			Headers:     `["X-Header-1", "foo", "X-Header-2", "bar", "odd one out"]`,
		}

		result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), nil)
		assert.False(t, runInfo.IsPending)
		assert.Error(t, result.Error)
		assert.Equal(t, `headers must have an even number of elements`, result.Error.Error())
		assert.Nil(t, result.Value)
	})

	t.Run("allows to override content-type", func(t *testing.T) {
		config := configtest.NewTestGeneralConfig(t)
		var headers http.Header
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headers = r.Header
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"fooresponse": 3}`))
			require.NoError(t, err)
		})

		server := httptest.NewServer(handler)
		defer server.Close()

		task := pipeline.HTTPTask{
			Method:      "POST",
			URL:         server.URL,
			RequestData: ethUSDPairing,
			Headers:     `["X-Header-1", "foo", "Content-Type", "footype", "X-Header-2", "bar"]`,
		}
		c := clhttptest.NewTestLocalOnlyHTTPClient()
		task.HelperSetDependencies(config.JobPipeline(), c, c)

		result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), nil)
		assert.False(t, runInfo.IsPending)
		assert.Equal(t, `{"fooresponse": 3}`, result.Value)
		assert.NoError(t, result.Error)

		assert.Equal(t, []string{"Content-Length", "38", "Content-Type", "footype", "User-Agent", "Go-http-client/1.1", "X-Header-1", "foo", "X-Header-2", "bar"}, allHeaders(headers))
	})
}
