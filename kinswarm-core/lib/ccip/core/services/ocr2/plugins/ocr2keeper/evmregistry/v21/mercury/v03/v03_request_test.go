package v03

import (
	"bytes"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"testing"
	"time"

	automationTypes "github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

const (
	defaultPluginRetryExpiration = 30 * time.Minute
	cleanupInterval              = 5 * time.Minute
)

type MockMercuryConfigProvider struct {
	cache *cache.Cache
	mock.Mock
}

func NewMockMercuryConfigProvider() *MockMercuryConfigProvider {
	return &MockMercuryConfigProvider{
		cache: cache.New(defaultPluginRetryExpiration, cleanupInterval),
	}
}

func (m *MockMercuryConfigProvider) Credentials() *types.MercuryCredentials {
	mc := &types.MercuryCredentials{
		LegacyURL: "https://google.old.com",
		URL:       "https://google.com",
		Username:  "FakeClientID",
		Password:  "FakeClientKey",
	}
	return mc
}

func (m *MockMercuryConfigProvider) IsUpkeepAllowed(s string) (interface{}, bool) {
	args := m.Called(s)
	return args.Get(0), args.Bool(1)
}

func (m *MockMercuryConfigProvider) SetUpkeepAllowed(s string, i interface{}, d time.Duration) {
	m.Called(s, i, d)
}

func (m *MockMercuryConfigProvider) GetPluginRetry(s string) (interface{}, bool) {
	if value, found := m.cache.Get(s); found {
		return value, true
	}

	return nil, false
}

func (m *MockMercuryConfigProvider) SetPluginRetry(s string, i interface{}, d time.Duration) {
	m.cache.Set(s, i, d)
}

type MockHttpClient struct {
	mock.Mock
}

func (mock *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	args := mock.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

// setups up a client object for tests.
func setupClient(t *testing.T) *client {
	lggr := logger.TestLogger(t)
	mockHttpClient := new(MockHttpClient)
	mercuryConfig := NewMockMercuryConfigProvider()
	threadCtl := utils.NewThreadControl()

	client := NewClient(
		mercuryConfig,
		mockHttpClient,
		threadCtl,
		lggr,
	)
	return client
}

func TestV03_DoMercuryRequestV03(t *testing.T) {
	t.Parallel()
	upkeepId, _ := new(big.Int).SetString("88786950015966611018675766524283132478093844178961698330929478019253453382042", 10)

	tests := []struct {
		name                  string
		lookup                *mercury.StreamsLookup
		mockHttpStatusCode    int
		mockChainlinkBlobs    []string
		pluginRetryKey        string
		expectedValues        [][]byte
		expectedRetryable     bool
		expectedRetryInterval time.Duration
		expectedErrCode       encoding.ErrCode
		expectedError         error
		state                 encoding.PipelineExecutionState
	}{
		{
			name: "success v0.3",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIDs,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(25880526),
					ExtraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
				},
				UpkeepId: upkeepId,
			},
			mockHttpStatusCode: http.StatusOK,
			mockChainlinkBlobs: []string{"0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"},
			expectedValues:     [][]byte{{0, 6, 109, 252, 209, 237, 45, 149, 177, 140, 148, 141, 188, 91, 214, 76, 104, 122, 254, 147, 228, 202, 125, 102, 61, 222, 193, 76, 32, 9, 10, 216, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8, 20, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 128, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 32, 69, 84, 72, 45, 85, 83, 68, 45, 65, 82, 66, 73, 84, 82, 85, 77, 45, 84, 69, 83, 84, 78, 69, 84, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 137, 28, 152, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 40, 154, 216, 211, 103, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 40, 154, 207, 11, 56, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 40, 155, 61, 164, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 138, 231, 206, 116, 217, 250, 37, 42, 137, 131, 151, 110, 171, 96, 13, 199, 89, 12, 119, 141, 4, 129, 52, 48, 132, 27, 198, 231, 101, 195, 76, 216, 26, 22, 141, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 138, 231, 203, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100, 137, 28, 152, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 96, 65, 43, 148, 229, 37, 202, 108, 237, 201, 245, 68, 253, 134, 247, 118, 6, 213, 47, 231, 49, 165, 208, 105, 219, 232, 54, 168, 191, 192, 251, 140, 145, 25, 99, 176, 174, 122, 20, 151, 31, 59, 70, 33, 191, 251, 128, 46, 240, 96, 83, 146, 185, 166, 200, 156, 127, 171, 29, 248, 99, 58, 90, 222, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 69, 0, 194, 245, 33, 248, 63, 186, 94, 252, 43, 243, 239, 250, 174, 221, 228, 61, 10, 74, 223, 247, 133, 193, 33, 59, 113, 42, 58, 237, 13, 129, 87, 100, 42, 132, 50, 77, 176, 207, 150, 149, 235, 210, 119, 8, 212, 96, 142, 176, 51, 126, 13, 216, 123, 14, 67, 240, 250, 112, 199, 0, 217, 17}},
			expectedRetryable:  false,
			expectedError:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupClient(t)
			defer c.Close()
			hc := mocks.NewHttpClient(t)

			mr := MercuryV03Response{}
			for i, blob := range tt.mockChainlinkBlobs {
				r := MercuryV03Report{
					FeedID:                tt.lookup.Feeds[i],
					ValidFromTimestamp:    0,
					ObservationsTimestamp: 0,
					FullReport:            blob,
				}
				mr.Reports = append(mr.Reports, r)
			}

			b, err := json.Marshal(mr)
			assert.Nil(t, err)
			resp := &http.Response{
				StatusCode: tt.mockHttpStatusCode,
				Body:       io.NopCloser(bytes.NewReader(b)),
			}
			if tt.expectedError != nil && tt.expectedRetryable {
				hc.On("Do", mock.Anything).Return(resp, nil).Times(totalAttempt)
			} else {
				hc.On("Do", mock.Anything).Return(resp, nil).Once()
			}
			c.httpClient = hc

			state, values, errCode, retryable, retryInterval, reqErr := c.DoRequest(testutils.Context(t), tt.lookup, automationTypes.ConditionTrigger, tt.pluginRetryKey)

			assert.Equal(t, tt.expectedValues, values)
			assert.Equal(t, tt.expectedRetryable, retryable)
			assert.Equal(t, tt.expectedRetryInterval, retryInterval)
			assert.Equal(t, tt.state, state)
			assert.Equal(t, tt.expectedErrCode, errCode)
			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError.Error(), reqErr.Error())
			}
		})
	}
}

func TestV03_DoMercuryRequestV03_MultipleFeedsSuccess(t *testing.T) {
	t.Parallel()
	upkeepId, _ := new(big.Int).SetString("88786950015966611018675766524283132478093844178961698330929478019253453382042", 10)
	pluginRetryKey := "88786950015966611018675766524283132478093844178961698330929478019253453382042|34"

	c := setupClient(t)
	defer c.Close()

	c.mercuryConfig.SetPluginRetry(pluginRetryKey, 0, cache.DefaultExpiration)
	hc := new(MockHttpClient)

	for i := 0; i <= 3; i++ {
		mr := MercuryV03Response{
			Reports: []MercuryV03Report{
				{
					FeedID:                "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
					ValidFromTimestamp:    123456,
					ObservationsTimestamp: 123456,
					FullReport:            "0xab2123dc00000012",
				},
			},
		}
		b, err := json.Marshal(mr)
		assert.Nil(t, err)

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(b)),
		}
		hc.On("Do", mock.Anything).Return(resp, nil).Once()
	}

	c.httpClient = hc

	lookup := &mercury.StreamsLookup{
		StreamsLookupError: &mercury.StreamsLookupError{
			FeedParamKey: mercury.FeedIdHex,
			Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
			TimeParamKey: mercury.BlockNumber,
			Time:         big.NewInt(25880526),
			ExtraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
		},
		UpkeepId: upkeepId,
	}

	state, _, errCode, retryable, retryInterval, _ := c.DoRequest(testutils.Context(t), lookup, automationTypes.ConditionTrigger, pluginRetryKey)
	assert.Equal(t, false, retryable)
	assert.Equal(t, 0*time.Second, retryInterval)
	assert.Equal(t, encoding.ErrCodeNil, errCode)
	assert.Equal(t, encoding.NoPipelineError, state)
}

func TestV03_DoMercuryRequestV03_Timeout(t *testing.T) {
	t.Skip("TODO: MERC-5965")
	t.Parallel()
	upkeepId, _ := new(big.Int).SetString("88786950015966611018675766524283132478093844178961698330929478019253453382042", 10)
	pluginRetryKey := "88786950015966611018675766524283132478093844178961698330929478019253453382042|34"

	c := setupClient(t)
	defer c.Close()

	c.mercuryConfig.SetPluginRetry(pluginRetryKey, 0, cache.DefaultExpiration)
	hc := new(MockHttpClient)

	mr := MercuryV03Response{
		Reports: []MercuryV03Report{
			{
				FeedID:                "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
				ValidFromTimestamp:    123456,
				ObservationsTimestamp: 123456,
				FullReport:            "0xab2123dc00000012",
			},
		},
	}
	b, err := json.Marshal(mr)
	assert.Nil(t, err)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(b)),
	}
	serverTimeout := 15 * time.Second // Server has delay of 15s, higher than mercury.RequestTimeout = 10s
	hc.On("Do", mock.Anything).Run(func(args mock.Arguments) {
		time.Sleep(serverTimeout)
	}).Return(resp, nil).Once()

	c.httpClient = hc

	lookup := &mercury.StreamsLookup{
		StreamsLookupError: &mercury.StreamsLookupError{
			FeedParamKey: mercury.FeedIdHex,
			Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
			TimeParamKey: mercury.BlockNumber,
			Time:         big.NewInt(25880526),
			ExtraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
		},
		UpkeepId: upkeepId,
	}

	start := time.Now()
	state, values, errCode, retryable, retryInterval, _ := c.DoRequest(testutils.Context(t), lookup, automationTypes.ConditionTrigger, pluginRetryKey)
	elapsed := time.Since(start)
	assert.True(t, elapsed < serverTimeout)
	assert.Equal(t, false, retryable)
	assert.Equal(t, 0*time.Second, retryInterval)
	assert.Equal(t, encoding.ErrCodeStreamsTimeout, errCode)
	assert.Equal(t, encoding.NoPipelineError, state)
	assert.Equal(t, [][]byte(nil), values)
}

func TestV03_DoMercuryRequestV03_OneFeedSuccessOneFeedPipelineError(t *testing.T) {
	t.Parallel()
	upkeepId, _ := new(big.Int).SetString("88786950015966611018675766524283132478093844178961698330929478019253453382042", 10)
	pluginRetryKey := "88786950015966611018675766524283132478093844178961698330929478019253453382042|34"

	c := setupClient(t)
	defer c.Close()

	c.mercuryConfig.SetPluginRetry(pluginRetryKey, 0, cache.DefaultExpiration)
	hc := new(MockHttpClient)

	// First request success
	mr := MercuryV03Response{
		Reports: []MercuryV03Report{
			{
				FeedID:                "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
				ValidFromTimestamp:    123456,
				ObservationsTimestamp: 123456,
				FullReport:            "0xab2123dc00000012",
			},
		},
	}
	b, err := json.Marshal(mr)
	assert.Nil(t, err)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(b)),
	}
	hc.On("Do", mock.Anything).Return(resp, nil).Once()
	// Second request returns MercuryFlakyError
	resp = &http.Response{
		StatusCode: http.StatusBadGateway,
		Body:       io.NopCloser(bytes.NewReader(b)),
	}
	hc.On("Do", mock.Anything).Return(resp, nil).Times(totalAttempt)
	c.httpClient = hc

	lookup := &mercury.StreamsLookup{
		StreamsLookupError: &mercury.StreamsLookupError{
			FeedParamKey: mercury.FeedIdHex,
			Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
			TimeParamKey: mercury.BlockNumber,
			Time:         big.NewInt(25880526),
			ExtraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
		},
		UpkeepId: upkeepId,
	}

	state, values, errCode, retryable, retryInterval, _ := c.DoRequest(testutils.Context(t), lookup, automationTypes.LogTrigger, pluginRetryKey)
	assert.Equal(t, true, retryable)
	assert.Equal(t, 1*time.Second, retryInterval)
	assert.Equal(t, encoding.ErrCodeStreamsBadGateway, errCode)
	assert.Equal(t, encoding.MercuryFlakyFailure, state)
	assert.Equal(t, [][]byte(nil), values)
}

func TestV03_DoMercuryRequestV03_OneFeedSuccessOneFeedErrCode(t *testing.T) {
	t.Parallel()
	upkeepId, _ := new(big.Int).SetString("88786950015966611018675766524283132478093844178961698330929478019253453382042", 10)
	pluginRetryKey := "88786950015966611018675766524283132478093844178961698330929478019253453382042|34"

	c := setupClient(t)
	defer c.Close()

	c.mercuryConfig.SetPluginRetry(pluginRetryKey, 0, cache.DefaultExpiration)
	hc := new(MockHttpClient)

	// First request success
	mr := MercuryV03Response{
		Reports: []MercuryV03Report{
			{
				FeedID:                "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
				ValidFromTimestamp:    123456,
				ObservationsTimestamp: 123456,
				FullReport:            "0xab2123dc00000012",
			},
		},
	}
	b, err := json.Marshal(mr)
	assert.Nil(t, err)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(b)),
	}
	hc.On("Do", mock.Anything).Return(resp, nil).Once()

	// Second request returns invalid response
	invalidResponse := MercuryV03Response{
		Reports: []MercuryV03Report{
			{
				FeedID:                "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
				ValidFromTimestamp:    123456,
				ObservationsTimestamp: 123456,
				FullReport:            "random", // invalid hex
			},
		},
	}
	b, err = json.Marshal(invalidResponse)
	assert.Nil(t, err)

	resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(b)),
	}
	hc.On("Do", mock.Anything).Return(resp, nil).Times(totalAttempt)
	c.httpClient = hc

	lookup := &mercury.StreamsLookup{
		StreamsLookupError: &mercury.StreamsLookupError{
			FeedParamKey: mercury.FeedIdHex,
			Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
			TimeParamKey: mercury.BlockNumber,
			Time:         big.NewInt(25880526),
			ExtraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
		},
		UpkeepId: upkeepId,
	}

	state, values, errCode, retryable, retryInterval, _ := c.DoRequest(testutils.Context(t), lookup, automationTypes.LogTrigger, pluginRetryKey)
	assert.Equal(t, [][]byte(nil), values)
	assert.Equal(t, false, retryable)
	assert.Equal(t, 0*time.Second, retryInterval)
	assert.Equal(t, encoding.ErrCodeStreamsBadResponse, errCode)
	assert.Equal(t, encoding.NoPipelineError, state)
}

func TestV03_MultiFeedRequest(t *testing.T) {
	t.Parallel()
	upkeepId := big.NewInt(123456789)
	tests := []struct {
		name           string
		lookup         *mercury.StreamsLookup
		statusCode     int
		lastStatusCode int
		pluginRetries  int
		pluginRetryKey string
		retryNumber    int
		retryable      bool
		errorMessage   string
		firstResponse  *MercuryV03Response
		response       *MercuryV03Response
		streamsErrCode encoding.ErrCode
		state          encoding.PipelineExecutionState
	}{
		{
			name: "success - mercury responds in the first try",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIDs,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.Timestamp,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			response: &MercuryV03Response{
				Reports: []MercuryV03Report{
					{
						FeedID:                "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
						ValidFromTimestamp:    123456,
						ObservationsTimestamp: 123456,
						FullReport:            "0xab2123dc00000012",
					},
					{
						FeedID:                "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000",
						ValidFromTimestamp:    123458,
						ObservationsTimestamp: 123458,
						FullReport:            "0xab2123dc00000016",
					},
				},
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success - mercury responds in the first try with blocknumber",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIDs,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			response: &MercuryV03Response{
				Reports: []MercuryV03Report{
					{
						FeedID:                "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
						ValidFromTimestamp:    123456,
						ObservationsTimestamp: 123456,
						FullReport:            "0xab2123dc00000012",
					},
					{
						FeedID:                "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000",
						ValidFromTimestamp:    123458,
						ObservationsTimestamp: 123458,
						FullReport:            "0xab2123dc00000016",
					},
				},
			},
			statusCode: http.StatusOK,
		},
		{
			name: "success - retry 206",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIDs,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.Timestamp,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			firstResponse: &MercuryV03Response{
				Reports: []MercuryV03Report{
					{
						FeedID:                "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
						ValidFromTimestamp:    123456,
						ObservationsTimestamp: 123456,
						FullReport:            "0xab2123dc00000012",
					},
				},
			},
			response: &MercuryV03Response{
				Reports: []MercuryV03Report{
					{
						FeedID:                "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
						ValidFromTimestamp:    123456,
						ObservationsTimestamp: 123456,
						FullReport:            "0xab2123dc00000012",
					},
					{
						FeedID:                "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000",
						ValidFromTimestamp:    123458,
						ObservationsTimestamp: 123458,
						FullReport:            "0xab2123dc00000019",
					},
				},
			},
			retryNumber:    1,
			statusCode:     http.StatusPartialContent,
			lastStatusCode: http.StatusOK,
		},
		{
			name: "success - retry for 500",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIDs,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.Timestamp,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			retryNumber:    2,
			statusCode:     http.StatusInternalServerError,
			lastStatusCode: http.StatusOK,
			response: &MercuryV03Response{
				Reports: []MercuryV03Report{
					{
						FeedID:                "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
						ValidFromTimestamp:    123456,
						ObservationsTimestamp: 123456,
						FullReport:            "0xab2123dc00000012",
					},
					{
						FeedID:                "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000",
						ValidFromTimestamp:    123458,
						ObservationsTimestamp: 123458,
						FullReport:            "0xab2123dc00000019",
					},
				},
			},
		},
		{
			name: "failure - invalid response and fail to decode reportBlob",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIDs,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.Timestamp,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			response: &MercuryV03Response{
				Reports: []MercuryV03Report{
					{
						FeedID:                "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
						ValidFromTimestamp:    123456,
						ObservationsTimestamp: 123456,
						FullReport:            "qerwiu", // invalid hex blob
					},
					{
						FeedID:                "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000",
						ValidFromTimestamp:    123458,
						ObservationsTimestamp: 123458,
						FullReport:            "0xab2123dc00000016",
					},
				},
			},
			statusCode:     http.StatusOK,
			retryable:      false,
			errorMessage:   "All attempts fail:\n#1: hex string without 0x prefix",
			streamsErrCode: encoding.ErrCodeStreamsBadResponse,
			state:          encoding.NoPipelineError,
		},
		{
			name: "failure - returns retryable with 1s plugin retry interval",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIDs,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.Timestamp,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			retryNumber:  totalAttempt,
			statusCode:   http.StatusInternalServerError,
			retryable:    true,
			errorMessage: "All attempts fail:\n#1: 500\n#2: 500\n#3: 500",
		},
		{
			name: "failure - returns retryable with 5s plugin retry interval",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIDs,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.Timestamp,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			pluginRetries: 6,
			retryNumber:   totalAttempt,
			statusCode:    http.StatusInternalServerError,
			retryable:     true,
			errorMessage:  "All attempts fail:\n#1: 500\n#2: 500\n#3: 500",
		},
		{
			name: "failure - returns retryable and then non-retryable",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIDs,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.Timestamp,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			retryNumber:    1,
			statusCode:     http.StatusInternalServerError,
			lastStatusCode: http.StatusUnauthorized,
			streamsErrCode: encoding.ErrCodeStreamsUnauthorized,
			state:          encoding.NoPipelineError,
		},
		{
			name: "failure - returns status code 422 not retryable",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIDs,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.Timestamp,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			statusCode:     http.StatusUnprocessableEntity,
			streamsErrCode: encoding.ErrCodeStreamsUnknownError,
			state:          encoding.NoPipelineError,
		},
		{
			name: "failure - StatusGatewayTimeout - returns retryable",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIDs,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.Timestamp,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			retryNumber:    totalAttempt,
			statusCode:     http.StatusGatewayTimeout,
			state:          encoding.MercuryFlakyFailure,
			retryable:      true,
			streamsErrCode: encoding.ErrCodeStreamsStatusGatewayTimeout,
			errorMessage:   "All attempts fail:\n#1: 504\n#2: 504\n#3: 504",
		},
		{
			name: "failure - StatusServiceUnavailable - returns retryable",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIDs,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.Timestamp,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			retryNumber:    totalAttempt,
			statusCode:     http.StatusServiceUnavailable,
			state:          encoding.MercuryFlakyFailure,
			retryable:      true,
			streamsErrCode: encoding.ErrCodeStreamsServiceUnavailable,
			errorMessage:   "All attempts fail:\n#1: 503\n#2: 503\n#3: 503",
		},
		{
			name: "failure - StatusBadGateway - returns retryable",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIDs,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.Timestamp,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			retryNumber:    totalAttempt,
			statusCode:     http.StatusBadGateway,
			streamsErrCode: encoding.ErrCodeStreamsBadGateway,
			state:          encoding.MercuryFlakyFailure,
			retryable:      true,
			errorMessage:   "All attempts fail:\n#1: 502\n#2: 502\n#3: 502",
		},

		{
			name: "failure - partial content three times with status ok",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIDs,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			response: &MercuryV03Response{
				Reports: []MercuryV03Report{
					{
						FeedID:                "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
						ValidFromTimestamp:    123456,
						ObservationsTimestamp: 123456,
						FullReport:            "0xab2123dc00000012",
					},
				},
			},
			statusCode:     http.StatusOK,
			retryNumber:    totalAttempt,
			retryable:      true,
			streamsErrCode: encoding.ErrCodeStreamsPartialContent,
			errorMessage:   "All attempts fail:\n#1: 404\n#2: 404\n#3: 404",
			state:          encoding.MercuryFlakyFailure,
		},
		{
			name: "failure - partial content three times with status partial content",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIDs,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			response: &MercuryV03Response{
				Reports: []MercuryV03Report{
					{
						FeedID:                "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
						ValidFromTimestamp:    123456,
						ObservationsTimestamp: 123456,
						FullReport:            "0xab2123dc00000012",
					},
				},
			},
			statusCode:   http.StatusPartialContent,
			retryNumber:  totalAttempt,
			retryable:    true,
			errorMessage: "All attempts fail:\n#1: 206\n#2: 206\n#3: 206",
			state:        encoding.MercuryFlakyFailure,
		},
		{
			name: "success - retry when reports length does not match feeds length",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIDs,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000", "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.Timestamp,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			firstResponse: &MercuryV03Response{
				Reports: []MercuryV03Report{
					{
						FeedID:                "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
						ValidFromTimestamp:    123456,
						ObservationsTimestamp: 123456,
						FullReport:            "0xab2123dc00000012",
					},
				},
			},
			response: &MercuryV03Response{
				Reports: []MercuryV03Report{
					{
						FeedID:                "0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
						ValidFromTimestamp:    123456,
						ObservationsTimestamp: 123456,
						FullReport:            "0xab2123dc00000012",
					},
					{
						FeedID:                "0x4254432d5553442d415242495452554d2d544553544e45540000000000000000",
						ValidFromTimestamp:    123458,
						ObservationsTimestamp: 123458,
						FullReport:            "0xab2123dc00000019",
					},
				},
			},
			retryNumber:    1,
			statusCode:     http.StatusOK,
			lastStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := setupClient(t)
			defer c.Close()
			if tt.pluginRetries != 0 {
				c.mercuryConfig.SetPluginRetry(tt.pluginRetryKey, tt.pluginRetries, cache.DefaultExpiration)
			}

			hc := new(MockHttpClient)
			b, err := json.Marshal(tt.response)
			assert.Nil(t, err)

			if tt.retryNumber == 0 {
				resp := &http.Response{
					StatusCode: tt.statusCode,
					Body:       io.NopCloser(bytes.NewReader(b)),
				}
				hc.On("Do", mock.Anything).Return(resp, nil).Once()
			} else if tt.retryNumber < totalAttempt {
				if tt.firstResponse != nil && tt.response != nil {
					b0, err := json.Marshal(tt.firstResponse)
					assert.Nil(t, err)
					resp0 := &http.Response{
						StatusCode: tt.statusCode,
						Body:       io.NopCloser(bytes.NewReader(b0)),
					}
					b1, err := json.Marshal(tt.response)
					assert.Nil(t, err)
					resp1 := &http.Response{
						StatusCode: tt.lastStatusCode,
						Body:       io.NopCloser(bytes.NewReader(b1)),
					}
					hc.On("Do", mock.Anything).Return(resp0, nil).Once().On("Do", mock.Anything).Return(resp1, nil).Once()
				} else {
					retryResp := &http.Response{
						StatusCode: tt.statusCode,
						Body:       io.NopCloser(bytes.NewReader(b)),
					}
					hc.On("Do", mock.Anything).Return(retryResp, nil).Times(tt.retryNumber)

					resp := &http.Response{
						StatusCode: tt.lastStatusCode,
						Body:       io.NopCloser(bytes.NewReader(b)),
					}
					hc.On("Do", mock.Anything).Return(resp, nil).Once()
				}
			} else {
				for i := 1; i <= tt.retryNumber; i++ {
					resp := &http.Response{
						StatusCode: tt.statusCode,
						Body:       io.NopCloser(bytes.NewReader(b)),
					}
					hc.On("Do", mock.Anything).Return(resp, nil).Once()
				}
			}
			c.httpClient = hc

			ch := make(chan mercury.MercuryData, 1)
			c.multiFeedsRequest(testutils.Context(t), ch, tt.lookup)

			m := <-ch
			assert.Equal(t, 0, m.Index)
			assert.Equal(t, tt.retryable, m.Retryable)
			if tt.streamsErrCode != encoding.ErrCodeNil {
				assert.Equal(t, tt.streamsErrCode, m.ErrCode)
				assert.Equal(t, tt.state, m.State)
				assert.Equal(t, [][]byte(nil), m.Bytes)
			} else if tt.retryNumber >= totalAttempt || tt.errorMessage != "" {
				assert.Equal(t, tt.errorMessage, m.Error.Error())
				assert.Equal(t, [][]byte(nil), m.Bytes)
			} else {
				assert.Nil(t, m.Error)
				var reports [][]byte
				for _, rsp := range tt.response.Reports {
					b, _ := hexutil.Decode(rsp.FullReport)
					reports = append(reports, b)
				}
				assert.Equal(t, reports, m.Bytes)
			}
		})
	}
}
