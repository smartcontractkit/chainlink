package pipeline_test

import (
	"context"
	"encoding/json"
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

func TestHTTPTask_Happy(t *testing.T) {
	t.Parallel()

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	s1 := httptest.NewServer(fakePriceResponder(t, utils.MustUnmarshalToMap(btcUSDPairing), decimal.NewFromInt(9700), "", nil))
	defer s1.Close()

	task := pipeline.HTTPTask{
		BaseTask:    pipeline.NewBaseTask(0, "http", nil, nil, 0),
		Method:      "POST",
		URL:         s1.URL,
		RequestData: btcUSDPairing,
	}
	task.HelperSetConfig(config)

	result := task.Run(context.Background(), pipeline.NewVarsFrom(nil), pipeline.JSONSerializable{emptyMeta, false}, nil)
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

func TestHTTPTask_Variables(t *testing.T) {
	t.Parallel()

	validMeta := map[string]interface{}{"theMeta": "yes"}

	tests := []struct {
		name                  string
		requestData           string
		meta                  pipeline.JSONSerializable
		inputs                []pipeline.Result
		vars                  pipeline.Vars
		expectedRequestData   map[string]interface{}
		expectedErrorCause    error
		expectedErrorContains string
	}{
		{
			"requestData (empty) + meta",
			``,
			pipeline.JSONSerializable{validMeta, false},
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"some_data": map[string]interface{}{"foo": 543.21}}),
			map[string]interface{}{},
			nil,
			"",
		},
		{
			"requestData (pure variable) + meta",
			`$(some_data)`,
			pipeline.JSONSerializable{validMeta, false},
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"some_data": map[string]interface{}{"foo": 543.21}}),
			map[string]interface{}{"foo": 543.21},
			nil,
			"",
		},
		{
			"requestData (pure variable)",
			`$(some_data)`,
			pipeline.JSONSerializable{nil, true},
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"some_data": map[string]interface{}{"foo": 543.21}}),
			map[string]interface{}{"foo": 543.21},
			nil,
			"",
		},
		{
			"requestData (pure variable, missing)",
			`$(some_data)`,
			pipeline.JSONSerializable{validMeta, false},
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"not_some_data": map[string]interface{}{"foo": 543.21}}),
			nil,
			pipeline.ErrKeypathNotFound,
			"requestData",
		},
		{
			"requestData (pure variable, not a map)",
			`$(some_data)`,
			pipeline.JSONSerializable{validMeta, false},
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"some_data": 543.21}),
			nil,
			pipeline.ErrBadInput,
			"requestData",
		},
		{
			"requestData (interpolation) + meta",
			`{"data":{"result":$(medianize)}}`,
			pipeline.JSONSerializable{validMeta, false},
			[]pipeline.Result{{Value: 123.45}},
			pipeline.NewVarsFrom(map[string]interface{}{"medianize": 543.21}),
			map[string]interface{}{"data": map[string]interface{}{"result": 543.21}},
			nil,
			"",
		},
		{
			"requestData (interpolation, missing)",
			`{"data":{"result":$(medianize)}}`,
			pipeline.JSONSerializable{validMeta, false},
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

			store, cleanup := cltest.NewStore(t)
			defer cleanup()

			s1 := httptest.NewServer(fakePriceResponder(t, test.expectedRequestData, decimal.NewFromInt(9700), "", nil))
			defer s1.Close()

			feedURL, err := url.ParseRequestURI(s1.URL)
			require.NoError(t, err)
			feedWebURL := (*models.WebURL)(feedURL)

			task := pipeline.BridgeTask{
				BaseTask:    pipeline.NewBaseTask(0, "bridge", nil, nil, 0),
				Name:        "foo",
				RequestData: test.requestData,
			}
			task.HelperSetConfigAndTxDB(store.Config, store.DB)

			// Insert bridge
			_, bridge := cltest.NewBridgeType(t, task.Name)
			bridge.URL = *feedWebURL
			require.NoError(t, store.ORM.DB.Create(&bridge).Error)

			result := task.Run(context.Background(), test.vars, test.meta, test.inputs)
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
				json.Unmarshal([]byte(result.Value.(string)), &x)
				require.Equal(t, decimal.NewFromInt(9700), x.Data.Result)
			}
		})
	}
}

func TestHTTPTask_ErrorMessage(t *testing.T) {
	t.Parallel()

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

	task := pipeline.HTTPTask{
		Method:      "POST",
		URL:         server.URL,
		RequestData: ethUSDPairing,
	}
	task.HelperSetConfig(config)

	result := task.Run(context.Background(), pipeline.NewVarsFrom(nil), pipeline.JSONSerializable{}, nil)
	require.Error(t, result.Error)
	require.Contains(t, result.Error.Error(), "could not hit data fetcher")
	require.Nil(t, result.Value)
}

func TestHTTPTask_OnlyErrorMessage(t *testing.T) {
	t.Parallel()

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

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
	task.HelperSetConfig(config)

	result := task.Run(context.Background(), pipeline.NewVarsFrom(nil), pipeline.JSONSerializable{}, nil)
	require.Error(t, result.Error)
	require.Contains(t, result.Error.Error(), "RequestId")
	require.Nil(t, result.Value)
}
