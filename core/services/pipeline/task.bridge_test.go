package pipeline_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/eth/contracts"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/pipeline/mocks"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// ethUSDPairing has the ETH/USD parameters needed when POSTing to the price
// external adapters.
// https://github.com/smartcontractkit/price-adapters

var (
	ethUSDPairing      = utils.MustUnmarshalToMap(`{"data":{"coin":"ETH","market":"USD"}}`)
	defaultHTTPTimeout = models.MustMakeDuration(15 * time.Second)
	emptyMeta          = utils.MustUnmarshalToMap("{}")
)

func TestBridgeTask_Happy(t *testing.T) {
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	btcUSDPairing := utils.MustUnmarshalToMap(`{"data":{"coin":"BTC","market":"USD"}}`)
	s1 := httptest.NewServer(fakePriceResponder(t, btcUSDPairing, decimal.NewFromInt(9700)))
	defer s1.Close()

	feedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)
	feedWebURL := (*models.WebURL)(feedURL)

	task := pipeline.BridgeTask{
		RequestData: pipeline.HttpRequestData{
			"data": map[string]interface{}{
				"coin":   "BTC",
				"market": "USD",
			},
		},
	}
	orm := new(mocks.ORM)
	orm.On("FindBridge", mock.Anything).Return(models.BridgeType{URL: *feedWebURL}, nil)
	task.HelperSetConfigAndORM(config, orm)

	result := task.Run(pipeline.TaskRun{
		PipelineRun: pipeline.Run{
			Meta: pipeline.JSONSerializable{emptyMeta},
		},
	}, nil)
	require.NoError(t, result.Error)
	require.NotNil(t, result.Value)
	var x struct {
		Data struct {
			Result decimal.Decimal `json:"result"`
		} `json:"data"`
	}
	json.Unmarshal(result.Value.([]byte), &x)
	require.Equal(t, decimal.NewFromInt(9700), x.Data.Result)
}

func TestBridgeTask_Meta(t *testing.T) {
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()
	config.Set("DEFAULT_HTTP_TIMEOUT", defaultHTTPTimeout.String())

	var empty adapterResponse

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req adapterRequest
		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &req)
		require.NoError(t, err)
		require.Equal(t, false, req.Meta["eligibleToSubmit"])
		require.Equal(t, float64(0), req.Meta["oracleCount"])
		require.Equal(t, float64(7), req.Meta["reportableRoundID"])
		require.Equal(t, float64(0), req.Meta["startedAt"])
		require.Equal(t, float64(11), req.Meta["timeout"])
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(empty))
	})

	roundState := contracts.FluxAggregatorRoundState{ReportableRoundID: 7, Timeout: 11}
	request, err := models.MarshalToMap(&roundState)
	require.NoError(t, err)

	s1 := httptest.NewServer(handler)

	defer s1.Close()
	feedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)
	feedWebURL := (*models.WebURL)(feedURL)

	task := pipeline.BridgeTask{
		RequestData: pipeline.HttpRequestData(ethUSDPairing),
	}
	orm := new(mocks.ORM)
	orm.On("FindBridge", mock.Anything).Return(models.BridgeType{URL: *feedWebURL}, nil)
	task.HelperSetConfigAndORM(config, orm)

	task.Run(pipeline.TaskRun{
		PipelineRun: pipeline.Run{
			Meta: pipeline.JSONSerializable{request},
		},
	}, nil)
}

func TestBridgeTask_ErrorMessage(t *testing.T) {
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

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
	feedWebURL := (*models.WebURL)(feedURL)

	task := pipeline.BridgeTask{
		RequestData: pipeline.HttpRequestData(ethUSDPairing),
	}
	orm := new(mocks.ORM)
	orm.On("FindBridge", mock.Anything).Return(models.BridgeType{URL: *feedWebURL}, nil)
	task.HelperSetConfigAndORM(config, orm)

	result := task.Run(pipeline.TaskRun{}, nil)
	require.Error(t, result.Error)
	require.Contains(t, result.Error.Error(), "could not hit data fetcher")
	require.Nil(t, result.Value)
}

func TestBridgeTask_OnlyErrorMessage(t *testing.T) {
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, err := w.Write([]byte(mustReadFile(t, "../testdata/coinmarketcap.error.json")))
		require.NoError(t, err)
	})

	server := httptest.NewServer(handler)
	defer server.Close()
	feedURL, err := url.ParseRequestURI(server.URL)
	require.NoError(t, err)
	feedWebURL := (*models.WebURL)(feedURL)

	task := pipeline.BridgeTask{
		RequestData: pipeline.HttpRequestData(ethUSDPairing),
	}
	orm := new(mocks.ORM)
	orm.On("FindBridge", mock.Anything).Return(models.BridgeType{URL: *feedWebURL}, nil)
	task.HelperSetConfigAndORM(config, orm)

	result := task.Run(pipeline.TaskRun{}, nil)
	require.Error(t, result.Error)
	require.Contains(t, result.Error.Error(), "RequestId")
	require.Nil(t, result.Value)
}

// Sample input taken from
// https://github.com/smartcontractkit/price-adapters#chainlink-price-request-adapters
func TestAdapterResponse_UnmarshalJSON_Happy(t *testing.T) {
	tests := []struct {
		name, content string
		expect        decimal.Decimal
	}{
		{"basic", `{"data":{"result":123.4567890},"jobRunID":"1","statusCode":200}`, decimal.NewFromFloat(123.456789)},
		{"bravenewcoin", mustReadFile(t, "../testdata/bravenewcoin.json"), decimal.NewFromFloat(306.52036004)},
		{"coinmarketcap", mustReadFile(t, "../testdata/coinmarketcap.json"), decimal.NewFromFloat(305.5574615)},
		{"cryptocompare", mustReadFile(t, "../testdata/cryptocompare.json"), decimal.NewFromFloat(305.76)},
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

func TestBridgeTask_AddsID(t *testing.T) {
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	var empty adapterResponse

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req adapterRequest
		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &req)
		require.NoError(t, err)
		require.NotEmpty(t, req.ID)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(empty))
	})

	s1 := httptest.NewServer(handler)
	defer s1.Close()
	feedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)
	feedWebURL := (*models.WebURL)(feedURL)

	task := pipeline.BridgeTask{
		RequestData: pipeline.HttpRequestData(ethUSDPairing),
	}
	orm := new(mocks.ORM)
	orm.On("FindBridge", mock.Anything).Return(models.BridgeType{URL: *feedWebURL}, nil)
	task.HelperSetConfigAndORM(config, orm)

	task.Run(pipeline.TaskRun{
		PipelineRun: pipeline.Run{
			Meta: pipeline.JSONSerializable{emptyMeta},
		},
	}, nil)
}
