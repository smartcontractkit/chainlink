package pipeline_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

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

func TestHTTPTask_Happy(t *testing.T) {
	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	btcUSDPairing := utils.MustUnmarshalToMap(`{"data":{"coin":"BTC","market":"USD"}}`)
	s1 := httptest.NewServer(fakePriceResponder(t, btcUSDPairing, decimal.NewFromInt(9700)))
	defer s1.Close()
	feedURL, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)
	feedWebURL := (*models.WebURL)(feedURL)

	task := pipeline.HTTPTask{
		Method: "POST",
		URL:    *feedWebURL,
		RequestData: pipeline.HttpRequestData{
			"data": map[string]interface{}{
				"coin":   "BTC",
				"market": "USD",
			},
		},
	}
	task.HelperSetConfig(config)

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

func TestHTTPTask_ErrorMessage(t *testing.T) {
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

	task := pipeline.HTTPTask{
		Method:      "POST",
		URL:         models.WebURL(*feedURL),
		RequestData: pipeline.HttpRequestData(ethUSDPairing),
	}
	task.HelperSetConfig(config)

	result := task.Run(context.Background(), pipeline.TaskRun{}, nil)
	require.Error(t, result.Error)
	require.Contains(t, result.Error.Error(), "could not hit data fetcher")
	require.Nil(t, result.Value)
}

func TestHTTPTask_OnlyErrorMessage(t *testing.T) {
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

	task := pipeline.HTTPTask{
		Method:      "POST",
		URL:         models.WebURL(*feedURL),
		RequestData: pipeline.HttpRequestData(ethUSDPairing),
	}
	task.HelperSetConfig(config)

	result := task.Run(context.Background(), pipeline.TaskRun{}, nil)
	require.Error(t, result.Error)
	require.Contains(t, result.Error.Error(), "RequestId")
	require.Nil(t, result.Value)
}
