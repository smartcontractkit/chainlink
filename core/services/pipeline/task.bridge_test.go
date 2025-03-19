package pipeline_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"sort"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"

	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	clhttptest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/httptest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline/eautils"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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
	eautils.AdapterStatus
	Data adapterResponseData `json:"data"`
}

func (pr *adapterResponse) SetStatusCode(code int) {
	pr.StatusCode = &code
}

func (pr *adapterResponse) UnsetStatusCode() {
	pr.StatusCode = nil
}

func (pr *adapterResponse) SetProviderStatusCode(code int) {
	pr.ProviderStatusCode = &code
}

func (pr *adapterResponse) UnsetProviderStatusCode() {
	pr.ProviderStatusCode = nil
}

func (pr *adapterResponse) SetError(msg string) {
	pr.Error = msg
}

func (pr *adapterResponse) UnsetError() {
	pr.Error = nil
}

func (pr *adapterResponse) SetErrorMessage(msg string) {
	pr.ErrorMessage = &msg
}

func (pr *adapterResponse) UnsetErrorMessage() {
	pr.ErrorMessage = nil
}

func (pr *adapterResponse) Result() *decimal.Decimal {
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

// NewMockHandler returns an http.HandlerFunc that responds with the given payload for any request
func NewMockHandler(payload string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(payload))
		if err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	}
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

func fakeIntermittentlyFailingPriceResponder(t *testing.T, requestData map[string]interface{}, result decimal.Decimal, inputKey string, expectedInput interface{}) http.Handler {
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
		// require.Equal(t, float64(0), reqBody.Meta["id"])

		if reqBody.Meta["shouldFail"].(bool) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			require.NoError(t, json.NewEncoder(w).Encode(errors.New("EA failure")))
			return
		}
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
	telemCh := make(chan interface{}, 1)
	ctx := pipeline.WithTelemetryCh(testutils.Context(t), telemCh)

	s1 := httptest.NewServer(fakePriceResponder(t, utils.MustUnmarshalToMap(btcUSDPairing), decimal.NewFromInt(9700), "", nil))
	defer s1.Close()

	feedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)

	orm := bridges.NewORM(db)
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()})

	task := pipeline.BridgeTask{
		BaseTask:    pipeline.NewBaseTask(0, "bridge", nil, nil, 0),
		Name:        bridge.Name.String(),
		RequestData: btcUSDPairing,
	}
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	trORM := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
	specID, err := trORM.CreateSpec(ctx, pipeline.Pipeline{}, *models.NewInterval(5 * time.Minute))
	require.NoError(t, err)
	task.HelperSetDependencies(cfg.JobPipeline(), cfg.WebServer(), orm, specID, uuid.UUID{}, c)

	result, runInfo := task.Run(ctx, logger.TestLogger(t), pipeline.NewVarsFrom(nil), nil)
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

	telem := <-telemCh
	require.IsType(t, &pipeline.BridgeTelemetry{}, telem)
	btelem := telem.(*pipeline.BridgeTelemetry)
	assert.Equal(t, string(bridge.Name), btelem.Name)
	assert.Equal(t, btcUSDPairing, string(btelem.RequestData))
	assert.Equal(t, `{"errorMessage":null,"error":null,"statusCode":null,"providerStatusCode":null,"data":{"result":"9700"}}
`, string(btelem.ResponseData))
	assert.Nil(t, btelem.ResponseError)
	assert.NotZero(t, btelem.RequestStartTimestamp)
	assert.NotZero(t, btelem.RequestFinishTimestamp)
	assert.Equal(t, 200, btelem.ResponseStatusCode)
	assert.False(t, btelem.LocalCacheHit)
	assert.Equal(t, specID, btelem.SpecID)
	assert.NotEqual(t, uuid.Nil, btelem.StreamID)
	assert.NotEqual(t, uuid.Nil, btelem.DotID)
}

func TestBridgeTask_HandlesIntermittentFailure(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {})

	s1 := httptest.NewServer(fakeIntermittentlyFailingPriceResponder(t, utils.MustUnmarshalToMap(btcUSDPairing), decimal.NewFromInt(9700), "", nil))
	defer s1.Close()

	feedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)

	orm := bridges.NewORM(db)
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()})

	task := pipeline.BridgeTask{
		BaseTask:    pipeline.NewBaseTask(0, "bridge", nil, nil, 0),
		Name:        bridge.Name.String(),
		RequestData: btcUSDPairing,
		CacheTTL:    "30s", // standard duration string format
	}
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	trORM := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
	specID, err := trORM.CreateSpec(testutils.Context(t), pipeline.Pipeline{}, *models.NewInterval(5 * time.Minute))
	require.NoError(t, err)
	task.HelperSetDependencies(cfg.JobPipeline(), cfg.WebServer(), orm, specID, uuid.UUID{}, c)
	result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t),
		pipeline.NewVarsFrom(
			map[string]interface{}{
				"jobRun": map[string]interface{}{
					"meta": map[string]interface{}{
						"shouldFail": false,
					},
				},
			},
		),
		nil)

	assert.False(t, runInfo.IsPending)
	assert.False(t, runInfo.IsRetryable)
	require.NoError(t, result.Error)
	require.NotNil(t, result.Value)

	result2, runInfo2 := task.Run(testutils.Context(t), logger.TestLogger(t),
		pipeline.NewVarsFrom(
			map[string]interface{}{
				"jobRun": map[string]interface{}{
					"meta": map[string]interface{}{
						"shouldFail": true,
					},
				},
			},
		),
		nil)

	require.NoError(t, result2.Error)
	require.Equal(t, result.Value, result2.Value)
	require.Equal(t, runInfo.IsPending, runInfo2.IsPending)
	require.Equal(t, runInfo.IsRetryable, runInfo2.IsRetryable)
}

func TestBridgeTask_DoesNotReturnStaleResults(t *testing.T) {
	t.Parallel()

	ctx := testutils.Context(t)

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.WebServer.BridgeCacheTTL = commonconfig.MustNewDuration(30 * time.Second)
	})

	s1 := httptest.NewServer(fakeIntermittentlyFailingPriceResponder(t, utils.MustUnmarshalToMap(btcUSDPairing), decimal.NewFromInt(9700), "", nil))
	defer s1.Close()

	feedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)

	db := pgtest.NewSqlxDB(t)
	orm := bridges.NewORM(db)
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()})

	task := pipeline.BridgeTask{
		BaseTask:    pipeline.NewBaseTask(0, "bridge", nil, nil, 0),
		Name:        bridge.Name.String(),
		RequestData: btcUSDPairing,
	}
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	trORM := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
	specID, err := trORM.CreateSpec(testutils.Context(t), pipeline.Pipeline{}, *models.NewInterval(5 * time.Minute))
	require.NoError(t, err)
	task.HelperSetDependencies(cfg.JobPipeline(), cfg.WebServer(), orm, specID, uuid.UUID{}, c)

	// Insert entry 1m in the past, stale value, should not be used in case of EA failure.
	_, err = db.ExecContext(ctx, `INSERT INTO bridge_last_value(dot_id, spec_id, value, finished_at)
	VALUES($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT bridge_last_value_pkey
	DO UPDATE SET value = $3, finished_at = $4;`, task.DotID(), specID, big.NewInt(9700).Bytes(), time.Now().Add(-1*time.Minute))
	require.NoError(t, err)

	result2, _ := task.Run(testutils.Context(t), logger.TestLogger(t),
		pipeline.NewVarsFrom(
			map[string]interface{}{
				"jobRun": map[string]interface{}{
					"meta": map[string]interface{}{
						"shouldFail": true,
					},
				},
			},
		),
		nil)

	require.Error(t, result2.Error)
	require.Nil(t, result2.Value)

	// Insert entry 10s in the past, under 30 seconds and should be used in case of failure.
	_, err = db.ExecContext(ctx, `INSERT INTO bridge_last_value(dot_id, spec_id, value, finished_at)
		VALUES($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT bridge_last_value_pkey
		DO UPDATE SET value = $3, finished_at = $4;`, task.DotID(), specID, big.NewInt(9700).Bytes(), time.Now().Add(-10*time.Second))
	require.NoError(t, err)

	result2, _ = task.Run(testutils.Context(t), logger.TestLogger(t),
		pipeline.NewVarsFrom(
			map[string]interface{}{
				"jobRun": map[string]interface{}{
					"meta": map[string]interface{}{
						"shouldFail": true,
					},
				},
			},
		),
		nil)

	require.NoError(t, result2.Error)
	require.Equal(t, string(big.NewInt(9700).Bytes()), result2.Value)

	cfg2 := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.WebServer.BridgeCacheTTL = commonconfig.MustNewDuration(0 * time.Second)
	})
	task.HelperSetDependencies(cfg2.JobPipeline(), cfg2.WebServer(), orm, specID, uuid.UUID{}, c)

	// Even though we have a cached value, this should fail since config now set to 0.
	result2, _ = task.Run(testutils.Context(t), logger.TestLogger(t),
		pipeline.NewVarsFrom(
			map[string]interface{}{
				"jobRun": map[string]interface{}{
					"meta": map[string]interface{}{
						"shouldFail": true,
					},
				},
			},
		),
		nil)

	require.Error(t, result2.Error)
	require.Nil(t, result2.Value)

	task2 := pipeline.BridgeTask{
		BaseTask:    pipeline.NewBaseTask(0, "bridge2", nil, nil, 0),
		Name:        bridge.Name.String(),
		RequestData: btcUSDPairing,
		CacheTTL:    "35m", // more than the stalenessCap 30m
	}
	task2.HelperSetDependencies(cfg2.JobPipeline(), cfg2.WebServer(), orm, specID, uuid.UUID{}, c)

	// Insert entry 32m in the past, under cacheTTL of 35m but more than stalenessCap of 30m.
	_, err = db.ExecContext(ctx, `INSERT INTO bridge_last_value(dot_id, spec_id, value, finished_at)
		VALUES($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT bridge_last_value_pkey
		DO UPDATE SET value = $3, finished_at = $4;`, task2.DotID(), specID, big.NewInt(9700).Bytes(), time.Now().Add(-32*time.Minute))
	require.NoError(t, err)

	// Run fails even though cacheTTL > lastvalue.finished_at because cacheTTL exceeds stalenessCap.
	result2, _ = task2.Run(testutils.Context(t), logger.TestLogger(t),
		pipeline.NewVarsFrom(
			map[string]interface{}{
				"jobRun": map[string]interface{}{
					"meta": map[string]interface{}{
						"shouldFail": true,
					},
				},
			},
		),
		nil)

	require.Error(t, result2.Error)
	require.Nil(t, result2.Value)

	// Insert entry 25m in the past, under stalenessCap
	_, err = db.ExecContext(ctx, `INSERT INTO bridge_last_value(dot_id, spec_id, value, finished_at)
		VALUES($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT bridge_last_value_pkey
		DO UPDATE SET value = $3, finished_at = $4;`, task2.DotID(), specID, big.NewInt(9700).Bytes(), time.Now().Add(-25*time.Minute))
	require.NoError(t, err)

	// Run succeeds using the cached value that's under stalenessCap.
	result2, _ = task2.Run(testutils.Context(t), logger.TestLogger(t),
		pipeline.NewVarsFrom(
			map[string]interface{}{
				"jobRun": map[string]interface{}{
					"meta": map[string]interface{}{
						"shouldFail": true,
					},
				},
			},
		),
		nil)

	require.NoError(t, result2.Error)
	require.Equal(t, string(big.NewInt(9700).Bytes()), result2.Value)
}

func TestBridgeTask_AsyncJobPendingState(t *testing.T) {
	t.Parallel()

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)

	id := uuid.New()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody adapterRequest
		payload, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		defer r.Body.Close()

		err = json.Unmarshal(payload, &reqBody)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%s/v2/resume/%v", cfg.WebServer().BridgeResponseURL(), id.String()), reqBody.ResponseURL)
		w.Header().Set("Content-Type", "application/json")

		// w.Header().Set("X-Chainlink-Pending", "true")
		response := map[string]interface{}{"pending": true}
		require.NoError(t, json.NewEncoder(w).Encode(response))
	})

	server := httptest.NewServer(handler)
	defer server.Close()
	feedURL, err := url.ParseRequestURI(server.URL)
	require.NoError(t, err)

	orm := bridges.NewORM(db)
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()})

	task := pipeline.BridgeTask{
		Name:        bridge.Name.String(),
		RequestData: ethUSDPairing,
		Async:       "true",
	}
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	trORM := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
	specID, err := trORM.CreateSpec(testutils.Context(t), pipeline.Pipeline{}, *models.NewInterval(5 * time.Minute))
	require.NoError(t, err)
	task.HelperSetDependencies(cfg.JobPipeline(), cfg.WebServer(), orm, specID, id, c)

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

			orm := bridges.NewORM(db)
			_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()})

			task := pipeline.BridgeTask{
				BaseTask:          pipeline.NewBaseTask(0, "bridge", nil, nil, 0),
				Name:              bridge.Name.String(),
				RequestData:       test.requestData,
				IncludeInputAtKey: test.includeInputAtKey,
			}
			c := clhttptest.NewTestLocalOnlyHTTPClient()
			trORM := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
			specID, err := trORM.CreateSpec(testutils.Context(t), pipeline.Pipeline{}, *models.NewInterval(5 * time.Minute))
			require.NoError(t, err)
			task.HelperSetDependencies(cfg.JobPipeline(), cfg.WebServer(), orm, specID, uuid.UUID{}, c)

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

	orm := bridges.NewORM(db)
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()})

	task := pipeline.BridgeTask{
		BaseTask:    pipeline.NewBaseTask(0, "bridge", nil, nil, 0),
		RequestData: ethUSDPairing,
		Name:        bridge.Name.String(),
	}
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	trORM := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
	specID, err := trORM.CreateSpec(testutils.Context(t), pipeline.Pipeline{}, *models.NewInterval(5 * time.Minute))
	require.NoError(t, err)
	task.HelperSetDependencies(cfg.JobPipeline(), cfg.WebServer(), orm, specID, uuid.UUID{}, c)

	mp := map[string]interface{}{"meta": metaDataForBridge}
	res, _ := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(map[string]interface{}{"jobRun": mp}), nil)
	assert.NoError(t, res.Error)

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

			orm := bridges.NewORM(db)
			_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()})

			task := pipeline.BridgeTask{
				BaseTask:          pipeline.NewBaseTask(0, "bridge", nil, nil, 0),
				Name:              bridge.Name.String(),
				RequestData:       btcUSDPairing,
				IncludeInputAtKey: test.includeInputAtKey,
			}
			c := clhttptest.NewTestLocalOnlyHTTPClient()
			trORM := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
			specID, err := trORM.CreateSpec(testutils.Context(t), pipeline.Pipeline{}, *models.NewInterval(5 * time.Minute))
			require.NoError(t, err)
			task.HelperSetDependencies(cfg.JobPipeline(), cfg.WebServer(), orm, specID, uuid.UUID{}, c)

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

		resp := &adapterResponse{}
		resp.SetErrorMessage("could not hit data fetcher")
		err := json.NewEncoder(w).Encode(resp)
		require.NoError(t, err)
	})

	server := httptest.NewServer(handler)
	defer server.Close()
	feedURL, err := url.ParseRequestURI(server.URL)
	require.NoError(t, err)

	orm := bridges.NewORM(db)
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()})

	task := pipeline.BridgeTask{
		Name:        bridge.Name.String(),
		RequestData: ethUSDPairing,
	}
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	trORM := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
	specID, err := trORM.CreateSpec(testutils.Context(t), pipeline.Pipeline{}, *models.NewInterval(5 * time.Minute))
	require.NoError(t, err)
	task.HelperSetDependencies(cfg.JobPipeline(), cfg.WebServer(), orm, specID, uuid.UUID{}, c)

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

	orm := bridges.NewORM(db)
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()})

	task := pipeline.BridgeTask{
		Name:        bridge.Name.String(),
		RequestData: ethUSDPairing,
	}
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	trORM := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
	specID, err := trORM.CreateSpec(testutils.Context(t), pipeline.Pipeline{}, *models.NewInterval(5 * time.Minute))
	require.NoError(t, err)
	task.HelperSetDependencies(cfg.JobPipeline(), cfg.WebServer(), orm, specID, uuid.UUID{}, c)

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
	orm := bridges.NewORM(db)
	trORM := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
	specID, err := trORM.CreateSpec(testutils.Context(t), pipeline.Pipeline{}, *models.NewInterval(5 * time.Minute))
	require.NoError(t, err)
	task.HelperSetDependencies(cfg.JobPipeline(), cfg.WebServer(), orm, specID, uuid.UUID{}, c)

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

func TestBridgeTask_Headers(t *testing.T) {
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)

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

	bridgeURL, err := url.ParseRequestURI(server.URL)
	require.NoError(t, err)

	orm := bridges.NewORM(db)
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: bridgeURL.String()})

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
			fmt.Println(k, v)
		}

		return s
	}

	standardHeaders := []string{"Content-Length", "38", "Content-Type", "application/json", "User-Agent", "Go-http-client/1.1"}

	t.Run("sends headers", func(t *testing.T) {
		task := pipeline.BridgeTask{
			BaseTask:    pipeline.NewBaseTask(0, "bridge", nil, nil, 0),
			Name:        bridge.Name.String(),
			RequestData: btcUSDPairing,
			Headers:     `["X-Header-1", "foo", "X-Header-2", "bar"]`,
		}

		c := clhttptest.NewTestLocalOnlyHTTPClient()
		trORM := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
		specID, err := trORM.CreateSpec(testutils.Context(t), pipeline.Pipeline{}, *models.NewInterval(5 * time.Minute))
		require.NoError(t, err)
		task.HelperSetDependencies(cfg.JobPipeline(), cfg.WebServer(), orm, specID, uuid.UUID{}, c)

		result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), nil)
		assert.False(t, runInfo.IsPending)
		assert.Equal(t, `{"fooresponse": 1}`, result.Value)
		assert.NoError(t, result.Error)

		assert.Equal(t, append(standardHeaders, "X-Header-1", "foo", "X-Header-2", "bar"), allHeaders(headers))
	})

	t.Run("errors with odd number of headers", func(t *testing.T) {
		task := pipeline.BridgeTask{
			BaseTask:    pipeline.NewBaseTask(0, "bridge", nil, nil, 0),
			Name:        bridge.Name.String(),
			RequestData: btcUSDPairing,
			Headers:     `["X-Header-1", "foo", "X-Header-2", "bar", "odd one out"]`,
		}

		c := clhttptest.NewTestLocalOnlyHTTPClient()
		trORM := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
		specID, err := trORM.CreateSpec(testutils.Context(t), pipeline.Pipeline{}, *models.NewInterval(5 * time.Minute))
		require.NoError(t, err)
		task.HelperSetDependencies(cfg.JobPipeline(), cfg.WebServer(), orm, specID, uuid.UUID{}, c)

		result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), nil)
		assert.False(t, runInfo.IsPending)
		assert.Error(t, result.Error)
		assert.Equal(t, `headers must have an even number of elements`, result.Error.Error())
		assert.Nil(t, result.Value)
	})

	t.Run("allows to override content-type", func(t *testing.T) {
		task := pipeline.BridgeTask{
			BaseTask:    pipeline.NewBaseTask(0, "bridge", nil, nil, 0),
			Name:        bridge.Name.String(),
			RequestData: btcUSDPairing,
			Headers:     `["X-Header-1", "foo", "Content-Type", "footype", "X-Header-2", "bar"]`,
		}

		c := clhttptest.NewTestLocalOnlyHTTPClient()
		trORM := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
		specID, err := trORM.CreateSpec(testutils.Context(t), pipeline.Pipeline{}, *models.NewInterval(5 * time.Minute))
		require.NoError(t, err)
		task.HelperSetDependencies(cfg.JobPipeline(), cfg.WebServer(), orm, specID, uuid.UUID{}, c)

		result, runInfo := task.Run(testutils.Context(t), logger.TestLogger(t), pipeline.NewVarsFrom(nil), nil)
		assert.False(t, runInfo.IsPending)
		assert.Equal(t, `{"fooresponse": 1}`, result.Value)
		assert.NoError(t, result.Error)

		assert.Equal(t, []string{"Content-Length", "38", "Content-Type", "footype", "User-Agent", "Go-http-client/1.1", "X-Header-1", "foo", "X-Header-2", "bar"}, allHeaders(headers))
	})
}

func TestBridgeTask_AdapterResponseStatusFailure(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.WebServer.BridgeCacheTTL = commonconfig.MustNewDuration(1 * time.Minute)
	})

	testAdapterResponse := &adapterResponse{
		Data: adapterResponseData{Result: &decimal.Zero},
	}

	s1 := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			err := json.NewEncoder(w).Encode(testAdapterResponse)
			require.NoError(t, err)
		}))
	defer s1.Close()

	feedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)

	// orm := bridges.NewORM(db)
	orm := bridges.NewCache(bridges.NewORM(db), logger.TestLogger(t), bridges.DefaultUpsertInterval)

	servicetest.Run(t, orm)

	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()})

	task := pipeline.BridgeTask{
		BaseTask:    pipeline.NewBaseTask(0, "bridge", nil, nil, 0),
		Name:        bridge.Name.String(),
		RequestData: btcUSDPairing,
	}
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	trORM := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
	specID, err := trORM.CreateSpec(ctx, pipeline.Pipeline{}, *models.NewInterval(5 * time.Minute))
	require.NoError(t, err)
	task.HelperSetDependencies(cfg.JobPipeline(), cfg.WebServer(), orm, specID, uuid.UUID{}, c)

	vars := pipeline.NewVarsFrom(
		map[string]interface{}{
			"jobRun": map[string]interface{}{
				"meta": map[string]interface{}{
					"shouldFail": true,
				},
			},
		},
	)

	testAdapterResponse.SetStatusCode(http.StatusInternalServerError)
	testAdapterResponse.Error = map[string]interface{}{
		"name":    "AdapterLWBAError",
		"message": "bid ask violation detected",
	}
	result, runInfo := task.Run(ctx, logger.TestLogger(t), vars, nil)

	require.ErrorContains(t, result.Error, "AdapterLWBAError: bid ask violation detected")
	require.Nil(t, result.Value)
	require.True(t, runInfo.IsRetryable)
	require.False(t, runInfo.IsPending)

	// Insert entry 1m in the past, stale value, should not be used in case of EA failure.
	_, err = db.ExecContext(ctx, `INSERT INTO bridge_last_value(dot_id, spec_id, value, finished_at)
	VALUES($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT bridge_last_value_pkey
	DO UPDATE SET value = $3, finished_at = $4;`, task.DotID(), specID, big.NewInt(9700).Bytes(), time.Now())
	require.NoError(t, err)

	// expect all external adapter response status failures to be served from the cache
	testAdapterResponse.SetStatusCode(http.StatusBadRequest)
	result, runInfo = task.Run(ctx, logger.TestLogger(t), vars, nil)

	require.NoError(t, result.Error)
	require.NotNil(t, result.Value)
	require.False(t, runInfo.IsRetryable)
	require.False(t, runInfo.IsPending)

	testAdapterResponse.SetStatusCode(http.StatusOK)
	testAdapterResponse.SetProviderStatusCode(http.StatusBadRequest)
	result, runInfo = task.Run(ctx, logger.TestLogger(t), vars, nil)

	require.NoError(t, result.Error)
	require.NotNil(t, result.Value)
	require.False(t, runInfo.IsRetryable)
	require.False(t, runInfo.IsPending)

	testAdapterResponse.SetStatusCode(http.StatusOK)
	testAdapterResponse.SetProviderStatusCode(http.StatusOK)
	testAdapterResponse.SetError("some error")
	result, runInfo = task.Run(ctx, logger.TestLogger(t), vars, nil)

	require.NoError(t, result.Error)
	require.NotNil(t, result.Value)
	require.False(t, runInfo.IsRetryable)
	require.False(t, runInfo.IsPending)

	testAdapterResponse.SetStatusCode(http.StatusInternalServerError)
	result, runInfo = task.Run(ctx, logger.TestLogger(t), vars, nil)

	require.NoError(t, result.Error)
	require.NotNil(t, result.Value)
	require.False(t, runInfo.IsRetryable)
	require.False(t, runInfo.IsPending)
}

func TestBridgeTask_AdapterTimeout(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)

	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.WebServer.BridgeCacheTTL = commonconfig.MustNewDuration(1 * time.Minute)
	})

	s1 := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Second) // delay enough to time-out
		}))
	defer s1.Close()

	feedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)

	orm := bridges.NewORM(db)
	_, bridge := cltest.MustCreateBridge(t, db, cltest.BridgeOpts{URL: feedURL.String()})

	task := pipeline.BridgeTask{
		BaseTask:    pipeline.NewBaseTask(0, "bridge", nil, nil, 0),
		Name:        bridge.Name.String(),
		RequestData: btcUSDPairing,
	}
	c := clhttptest.NewTestLocalOnlyHTTPClient()
	trORM := pipeline.NewORM(db, logger.TestLogger(t), cfg.JobPipeline().MaxSuccessfulRuns())
	specID, err := trORM.CreateSpec(ctx, pipeline.Pipeline{}, *models.NewInterval(5 * time.Minute))
	require.NoError(t, err)
	task.HelperSetDependencies(cfg.JobPipeline(), cfg.WebServer(), orm, specID, uuid.UUID{}, c)

	// Insert entry 1m in the past, stale value, should not be used in case of EA failure.
	_, err = db.ExecContext(ctx, `INSERT INTO bridge_last_value(dot_id, spec_id, value, finished_at)
	VALUES($1, $2, $3, $4) ON CONFLICT ON CONSTRAINT bridge_last_value_pkey
	DO UPDATE SET value = $3, finished_at = $4;`, task.DotID(), specID, big.NewInt(9700).Bytes(), time.Now())
	require.NoError(t, err)

	vars := pipeline.NewVarsFrom(
		map[string]interface{}{
			"jobRun": map[string]interface{}{
				"meta": map[string]interface{}{
					"shouldFail": true,
				},
			},
		},
	)

	t.Run("pre-cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(testutils.Context(t))
		cancel() // pre-cancelled
		result, runInfo := task.Run(ctx, logger.TestLogger(t), vars, nil)

		require.NoError(t, result.Error)
		require.NotNil(t, result.Value)
		require.False(t, runInfo.IsRetryable)
		require.False(t, runInfo.IsPending)
	})

	t.Run("short", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(testutils.Context(t), time.Millisecond)
		t.Cleanup(cancel)
		result, runInfo := task.Run(ctx, logger.TestLogger(t), vars, nil)

		require.NoError(t, result.Error)
		require.NotNil(t, result.Value)
		require.False(t, runInfo.IsRetryable)
		require.False(t, runInfo.IsPending)
	})
}

func TestBridgeTask_PipelineAdapterLWBAError(t *testing.T) {
	t.Parallel()

	dag := `
ds [type=bridge name="adapter-error-bridge" timeout="50ms" requestData="{\"data\":{\"from\":\"ETH\",\"to\":\"USD\"}}"];
`

	ctx := testutils.Context(t)
	db := pgtest.NewSqlxDB(t)
	cfg := configtest.NewTestGeneralConfig(t)
	orm := bridges.NewORM(db)
	r, _ := newRunner(t, db, orm, cfg)

	bridge := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		b, herr := io.ReadAll(req.Body)
		require.NoError(t, herr)
		require.JSONEq(t, `{"data":{"from":"ETH","to":"USD"}}`, string(b))

		res.WriteHeader(http.StatusInternalServerError)
		resp := `{"error": {"name":"AdapterLWBAError", "message": "bid ask violation detected"}}`
		_, herr = res.Write([]byte(resp))
		require.NoError(t, herr)
	}))
	t.Cleanup(bridge.Close)
	u, _ := url.Parse(bridge.URL)
	require.NoError(t, orm.CreateBridgeType(ctx, &bridges.BridgeType{
		Name: "adapter-error-bridge",
		URL:  models.WebURL(*u),
	}))

	spec := pipeline.Spec{DotDagSource: dag}
	vars := pipeline.NewVarsFrom(nil)

	_, trrs, err := r.ExecuteRun(ctx, spec, vars)

	require.NoError(t, err)
	require.Len(t, trrs, 1)

	finalResult := trrs[0]

	require.ErrorContains(t, finalResult.Result.Error, "AdapterLWBAError: bid ask violation detected")
	require.Nil(t, finalResult.Result.Value)
}
