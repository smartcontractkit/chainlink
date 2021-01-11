package pipeline_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/eth/contracts"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// ethUSDPairing has the ETH/USD parameters needed when POSTing to the price
// external adapters.
// https://github.com/smartcontractkit/price-adapters

var (
	ethUSDPairing = utils.MustUnmarshalToMap(`{"data":{"coin":"ETH","market":"USD"}}`)
	emptyMeta     = utils.MustUnmarshalToMap("{}")
)

type adapterRequest struct {
	ID   string                   `json:"id"`
	Data pipeline.HttpRequestData `json:"data"`
	Meta pipeline.HttpRequestData `json:"meta"`
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

	content, err := ioutil.ReadFile(file)
	require.NoError(t, err)
	return string(content)
}

func fakePriceResponder(t *testing.T, requestData map[string]interface{}, result decimal.Decimal) http.Handler {
	t.Helper()

	body, err := json.Marshal(requestData)
	require.NoError(t, err)
	var expectedRequest adapterRequest
	err = json.Unmarshal(body, &expectedRequest)
	require.NoError(t, err)
	response := adapterResponse{Data: dataWithResult(t, result)}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody adapterRequest
		payload, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		defer r.Body.Close()
		err = json.Unmarshal(payload, &reqBody)
		require.NoError(t, err)
		require.Equal(t, expectedRequest.Data, reqBody.Data)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(response))
	})
}

func TestBridgeTask_Happy(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	btcUSDPairing := utils.MustUnmarshalToMap(`{"data":{"coin":"BTC","market":"USD"}}`)
	s1 := httptest.NewServer(fakePriceResponder(t, btcUSDPairing, decimal.NewFromInt(9700)))
	defer s1.Close()

	feedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)
	feedWebURL := (*models.WebURL)(feedURL)

	task := pipeline.BridgeTask{
		Name: "foo",
		RequestData: pipeline.HttpRequestData{
			"data": map[string]interface{}{
				"coin":   "BTC",
				"market": "USD",
			},
		},
	}
	task.HelperSetConfigAndTxDB(store.Config, store.DB)

	// Insert bridge
	_, bridge := cltest.NewBridgeType(t, task.Name)
	bridge.URL = *feedWebURL
	require.NoError(t, store.ORM.DB.Create(&bridge).Error)

	result := task.Run(context.Background(), pipeline.TaskRun{
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
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

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

	roundState := contracts.FluxAggregatorRoundState{RoundId: 7, Timeout: 11}
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
	task.HelperSetConfigAndTxDB(store.Config, store.DB)

	_, bridge := cltest.NewBridgeType(t)
	bridge.URL = *feedWebURL
	require.NoError(t, store.ORM.DB.Create(&bridge).Error)

	task.Run(context.Background(), pipeline.TaskRun{
		PipelineRun: pipeline.Run{
			Meta: pipeline.JSONSerializable{request},
		},
	}, nil)
}

func TestBridgeTask_ErrorMessage(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
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
		Name:        "foo",
		RequestData: pipeline.HttpRequestData(ethUSDPairing),
	}
	task.HelperSetConfigAndTxDB(store.Config, store.DB)

	_, bridge := cltest.NewBridgeType(t, task.Name)
	bridge.URL = *feedWebURL
	require.NoError(t, store.ORM.DB.Create(&bridge).Error)

	result := task.Run(context.Background(), pipeline.TaskRun{}, nil)
	require.Error(t, result.Error)
	require.Contains(t, result.Error.Error(), "could not hit data fetcher")
	require.Nil(t, result.Value)
}

func TestBridgeTask_OnlyErrorMessage(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
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
		Name:        "foo",
		RequestData: pipeline.HttpRequestData(ethUSDPairing),
	}
	task.HelperSetConfigAndTxDB(store.Config, store.DB)

	_, bridge := cltest.NewBridgeType(t, task.Name)
	bridge.URL = *feedWebURL
	require.NoError(t, store.ORM.DB.Create(&bridge).Error)

	result := task.Run(context.Background(), pipeline.TaskRun{}, nil)
	require.Error(t, result.Error)
	require.Contains(t, result.Error.Error(), "RequestId")
	require.Nil(t, result.Value)
}

func TestBridgeTask_ErrorIfBridgeMissing(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	task := pipeline.BridgeTask{
		Name: "foo",
		RequestData: pipeline.HttpRequestData{
			"data": map[string]interface{}{
				"coin":   "BTC",
				"market": "USD",
			},
		},
	}
	task.HelperSetConfigAndTxDB(store.Config, store.DB)

	result := task.Run(context.Background(), pipeline.TaskRun{
		PipelineRun: pipeline.Run{
			Meta: pipeline.JSONSerializable{emptyMeta},
		},
	}, nil)
	require.Nil(t, result.Value)
	require.Error(t, result.Error)
	require.Equal(t, "could not find bridge with name 'foo': record not found", result.Error.Error())
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
	store, cleanup := cltest.NewStore(t)
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
		Name:        "test",
		RequestData: pipeline.HttpRequestData(ethUSDPairing),
	}
	// Bridge task should be true by default.
	store.Config.Set("DEFAULT_HTTP_ALLOW_UNRESTRICTED_NETWORK_ACCESS", false)
	task.HelperSetConfigAndTxDB(store.Config, store.DB)

	_, bridge := cltest.NewBridgeType(t)
	bridge.URL = *feedWebURL
	bridge.Name = "test"
	require.NoError(t, store.ORM.DB.Create(&bridge).Error)

	r := task.Run(context.Background(), pipeline.TaskRun{
		PipelineRun: pipeline.Run{
			Meta: pipeline.JSONSerializable{emptyMeta},
		},
	}, nil)
	require.NoError(t, r.Error)
}
