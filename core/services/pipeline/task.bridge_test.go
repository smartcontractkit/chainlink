package pipeline_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
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
		payload, err := ioutil.ReadAll(r.Body)
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

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	btcUSDPairing := utils.MustUnmarshalToMap(`{"data":{"coin":"BTC","market":"USD"}}`)
	s1 := httptest.NewServer(fakePriceResponder(t, btcUSDPairing, decimal.NewFromInt(9700), "", nil))
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

	result := task.Run(context.Background(), pipeline.JSONSerializable{emptyMeta, false}, nil)
	require.NoError(t, result.Error)
	require.NotNil(t, result.Value)
	var x struct {
		Data struct {
			Result decimal.Decimal `json:"result"`
		} `json:"data"`
	}
	json.Unmarshal([]byte(result.Value.(string)), &x)
	require.Equal(t, decimal.NewFromInt(9700), x.Data.Result)
}

func TestBridgeTask_Meta(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	var empty adapterResponse

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req adapterRequest
		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &req)
		require.NoError(t, err)
		require.Equal(t, 10, req.Meta["latestAnswer"])
		require.Equal(t, 1616447984, req.Meta["updatedAt"])
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(empty))
	})

	metaDataForBridge, err := models.MarshalBridgeMetaData(big.NewInt(10), big.NewInt(1616447984))
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

	task.Run(context.Background(), pipeline.JSONSerializable{metaDataForBridge, false}, nil)
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
		{"too many inputs", []pipeline.Result{{Value: decimal.NewFromFloat(123.45)}, {Value: decimal.NewFromFloat(321.45)}}, "result", nil, pipeline.ErrWrongInputCardinality},
		{"input has error", []pipeline.Result{{Error: theErr}}, "result", nil, theErr},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			btcUSDPairing := utils.MustUnmarshalToMap(`{"data":{"coin":"BTC","market":"USD"}}`)
			s1 := httptest.NewServer(fakePriceResponder(t, btcUSDPairing, decimal.NewFromInt(9700), test.includeInputAtKey, test.expectedInput))
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
				IncludeInputAtKey: test.includeInputAtKey,
			}
			task.HelperSetConfigAndTxDB(store.Config, store.DB)

			// Insert bridge
			_, bridge := cltest.NewBridgeType(t, task.Name)
			bridge.URL = *feedWebURL
			require.NoError(t, store.ORM.DB.Create(&bridge).Error)

			result := task.Run(context.Background(), pipeline.JSONSerializable{emptyMeta, false}, test.inputs)
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
				json.Unmarshal([]byte(result.Value.(string)), &x)
				require.Equal(t, decimal.NewFromInt(9700), x.Data.Result)
			}
		})
	}
}

func TestBridgeTask_ErrorMessage(t *testing.T) {
	t.Parallel()

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

	result := task.Run(context.Background(), pipeline.JSONSerializable{}, nil)
	require.Error(t, result.Error)
	require.Contains(t, result.Error.Error(), "could not hit data fetcher")
	require.Nil(t, result.Value)
}

func TestBridgeTask_OnlyErrorMessage(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

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
	feedWebURL := (*models.WebURL)(feedURL)

	task := pipeline.BridgeTask{
		Name:        "foo",
		RequestData: pipeline.HttpRequestData(ethUSDPairing),
	}
	task.HelperSetConfigAndTxDB(store.Config, store.DB)

	_, bridge := cltest.NewBridgeType(t, task.Name)
	bridge.URL = *feedWebURL
	require.NoError(t, store.ORM.DB.Create(&bridge).Error)

	result := task.Run(context.Background(), pipeline.JSONSerializable{}, nil)
	require.Error(t, result.Error)
	require.Contains(t, result.Error.Error(), "RequestId")
	require.Nil(t, result.Value)
}

func TestBridgeTask_ErrorIfBridgeMissing(t *testing.T) {
	t.Parallel()

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

	result := task.Run(context.Background(), pipeline.JSONSerializable{emptyMeta, false}, nil)
	require.Nil(t, result.Value)
	require.Error(t, result.Error)
	require.Equal(t, "could not find bridge with name 'foo': record not found", result.Error.Error())
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
