package v02

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	automationTypes "github.com/smartcontractkit/chainlink-automation/pkg/v3/types"

	"github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/encoding"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/utils"

	"github.com/stretchr/testify/assert"
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

func TestV02_SingleFeedRequest(t *testing.T) {
	t.Parallel()
	upkeepId := big.NewInt(123456789)
	tests := []struct {
		name           string
		index          int
		lookup         *mercury.StreamsLookup
		blob           string
		responseBytes  string
		statusCode     int
		lastStatusCode int
		retryNumber    int
		retryable      bool
		errorMessage   string
		streamsErrCode encoding.ErrCode
		state          encoding.PipelineExecutionState
	}{
		{
			name:  "success - mercury responds in the first try",
			index: 0,
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			blob: "0xab2123dc00000012",
		},
		{
			name:  "success - retry for 404",
			index: 0,
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			blob:           "0xab2123dcbabbad",
			retryNumber:    1,
			statusCode:     http.StatusNotFound,
			lastStatusCode: http.StatusOK,
		},
		{
			name:  "success - retry for 500",
			index: 0,
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			blob:           "0xab2123dcbbabad",
			retryNumber:    2,
			statusCode:     http.StatusInternalServerError,
			lastStatusCode: http.StatusOK,
		},
		{
			name:  "failure - returns retryable",
			index: 0,
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			blob:         "0xab2123dc",
			retryNumber:  totalAttempt,
			statusCode:   http.StatusNotFound,
			state:        encoding.MercuryFlakyFailure,
			retryable:    true,
			errorMessage: "failed to request feed for 0x4554482d5553442d415242495452554d2d544553544e45540000000000000000: All attempts fail:\n#1: 404\n#2: 404\n#3: 404",
		},
		{
			name:  "failure StatusGatewayTimeout - returns retryable",
			index: 0,
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			blob:         "0xab2123dc",
			retryNumber:  totalAttempt,
			statusCode:   http.StatusGatewayTimeout,
			state:        encoding.MercuryFlakyFailure,
			retryable:    true,
			errorMessage: "failed to request feed for 0x4554482d5553442d415242495452554d2d544553544e45540000000000000000: All attempts fail:\n#1: 504\n#2: 504\n#3: 504",
		},
		{
			name:  "failure StatusServiceUnavailable - returns retryable",
			index: 0,
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			blob:           "0xab2123dc",
			retryNumber:    totalAttempt,
			statusCode:     http.StatusServiceUnavailable,
			streamsErrCode: encoding.ErrCodeStreamsServiceUnavailable,
			state:          encoding.MercuryFlakyFailure,
			retryable:      true,
			errorMessage:   "failed to request feed for 0x4554482d5553442d415242495452554d2d544553544e45540000000000000000: All attempts fail:\n#1: 503\n#2: 503\n#3: 503",
		},
		{
			name:  "failure StatusBadGateway - returns retryable",
			index: 0,
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			blob:           "0xab2123dc",
			retryNumber:    totalAttempt,
			statusCode:     http.StatusBadGateway,
			streamsErrCode: encoding.ErrCodeStreamsBadGateway,
			state:          encoding.MercuryFlakyFailure,
			retryable:      true,
			errorMessage:   "failed to request feed for 0x4554482d5553442d415242495452554d2d544553544e45540000000000000000: All attempts fail:\n#1: 502\n#2: 502\n#3: 502",
		},
		{
			name:  "failure - returns retryable and then non-retryable",
			index: 0,
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			blob:           "0xab2123dc",
			retryNumber:    1,
			statusCode:     http.StatusNotFound,
			lastStatusCode: http.StatusTooManyRequests,
			streamsErrCode: encoding.ErrCodeStreamsUnknownError,
		},
		{
			name:  "failure - returns not retryable",
			index: 0,
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			blob:           "0xab2123dc",
			statusCode:     http.StatusConflict,
			streamsErrCode: encoding.ErrCodeStreamsUnknownError,
			errorMessage:   "failed to request feed for 0x4554482d5553442d415242495452554d2d544553544e45540000000000000000: All attempts fail:\n#1: at block 123456 upkeep 123456789 received status code 409 for feed 0x4554482d5553442d415242495452554d2d544553544e45540000000000000000",
		},
		{
			name:  "failure invalid json response - returns err code",
			index: 0,
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			responseBytes:  "invalid response",
			streamsErrCode: encoding.ErrCodeStreamsBadResponse,
		},
		{
			name:  "failure invalid blob - returns err code",
			index: 0,
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(123456),
				},
				UpkeepId: upkeepId,
			},
			blob:           "0xgz", // Invalid hex
			streamsErrCode: encoding.ErrCodeStreamsBadResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := setupClient(t)
			defer c.Close()
			hc := new(MockHttpClient)

			mr := MercuryV02Response{ChainlinkBlob: tt.blob}
			b, err := json.Marshal(mr)
			assert.NoError(t, err)

			if tt.responseBytes != "" {
				b = []byte(tt.responseBytes)
			}

			if tt.retryNumber == 0 {
				if tt.errorMessage != "" {
					resp := &http.Response{
						StatusCode: tt.statusCode,
						Body:       io.NopCloser(bytes.NewReader(b)),
					}
					hc.On("Do", mock.Anything).Return(resp, nil).Once()
				} else {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader(b)),
					}
					hc.On("Do", mock.Anything).Return(resp, nil).Once()
				}
			} else if tt.retryNumber > 0 && tt.retryNumber < totalAttempt {
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
			} else {
				resp := &http.Response{
					StatusCode: tt.statusCode,
					Body:       io.NopCloser(bytes.NewReader(b)),
				}
				hc.On("Do", mock.Anything).Return(resp, nil).Times(tt.retryNumber)
			}
			c.httpClient = hc

			ch := make(chan mercury.MercuryData, 1)
			c.singleFeedRequest(testutils.Context(t), ch, tt.index, tt.lookup)

			m := <-ch
			assert.Equal(t, tt.index, m.Index)
			assert.Equal(t, tt.retryable, m.Retryable)
			assert.Equal(t, tt.state, m.State)
			if tt.streamsErrCode != encoding.ErrCodeNil {
				assert.Equal(t, tt.streamsErrCode, m.ErrCode)
				assert.Equal(t, [][]byte(nil), m.Bytes)
			} else if tt.retryNumber >= totalAttempt || tt.errorMessage != "" {
				assert.Equal(t, tt.errorMessage, m.Error.Error())
				assert.Equal(t, [][]byte(nil), m.Bytes)
			} else {
				blobBytes, err := hexutil.Decode(tt.blob)
				assert.NoError(t, err)
				assert.NoError(t, m.Error)
				assert.Equal(t, [][]byte{blobBytes}, m.Bytes)
			}
		})
	}
}

func TestV02_DoMercuryRequestV02(t *testing.T) {
	t.Parallel()
	upkeepId, _ := new(big.Int).SetString("88786950015966611018675766524283132478093844178961698330929478019253453382042", 10)
	pluginRetryKey := "88786950015966611018675766524283132478093844178961698330929478019253453382042|34"
	tests := []struct {
		name                  string
		lookup                *mercury.StreamsLookup
		mockHttpStatusCode    int
		mockChainlinkBlobs    []string
		pluginRetries         int
		upkeepType            automationTypes.UpkeepType
		expectedState         encoding.PipelineExecutionState
		expectedValues        [][]byte
		expectedErrCode       encoding.ErrCode
		expectedRetryable     bool
		expectedRetryInterval time.Duration
		expectedError         error
	}{
		{
			name: "success",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
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
		{
			name: "failure - retryable but out of retries for conditional",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(25880526),
					ExtraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
				},
				UpkeepId: upkeepId,
			},
			upkeepType:            automationTypes.ConditionTrigger,
			mockHttpStatusCode:    http.StatusInternalServerError,
			mockChainlinkBlobs:    []string{"0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"},
			expectedValues:        [][]byte(nil),
			expectedRetryable:     false,
			pluginRetries:         0,
			expectedRetryInterval: 0 * time.Second,
			expectedErrCode:       encoding.ErrCodeStreamsInternalError,
			expectedState:         encoding.NoPipelineError,
		},
		{
			name: "failure - retryable but out of retries for conditional 404 error code",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(25880526),
					ExtraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
				},
				UpkeepId: upkeepId,
			},
			upkeepType:            automationTypes.ConditionTrigger,
			mockHttpStatusCode:    http.StatusNotFound,
			mockChainlinkBlobs:    []string{"0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"},
			expectedValues:        [][]byte(nil),
			expectedRetryable:     false,
			pluginRetries:         0,
			expectedRetryInterval: 0 * time.Second,
			expectedErrCode:       encoding.ErrCodeStreamsNotFound,
			expectedState:         encoding.NoPipelineError,
		},
		{
			name: "failure - retryable and interval is 1s",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(25880526),
					ExtraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
				},
				UpkeepId: upkeepId,
			},
			upkeepType:            automationTypes.LogTrigger,
			pluginRetries:         1,
			mockHttpStatusCode:    http.StatusInternalServerError,
			mockChainlinkBlobs:    []string{"0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"},
			expectedValues:        [][]byte(nil),
			expectedRetryable:     true,
			expectedRetryInterval: 1 * time.Second,
			expectedErrCode:       encoding.ErrCodeStreamsInternalError,
			expectedState:         encoding.MercuryFlakyFailure,
			expectedError:         errors.New("failed to request feed for 0x4554482d5553442d415242495452554d2d544553544e45540000000000000000: All attempts fail:\n#1: 500\n#2: 500\n#3: 500"),
		},
		{
			name: "failure - retryable and interval is 5s",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(25880526),
					ExtraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
				},
				UpkeepId: upkeepId,
			},
			upkeepType:            automationTypes.LogTrigger,
			pluginRetries:         5,
			mockHttpStatusCode:    http.StatusInternalServerError,
			mockChainlinkBlobs:    []string{"0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"},
			expectedValues:        [][]byte(nil),
			expectedRetryable:     true,
			expectedRetryInterval: 5 * time.Second,
			expectedErrCode:       encoding.ErrCodeStreamsInternalError,
			expectedState:         encoding.MercuryFlakyFailure,
			expectedError:         errors.New("failed to request feed for 0x4554482d5553442d415242495452554d2d544553544e45540000000000000000: All attempts fail:\n#1: 500\n#2: 500\n#3: 500"),
		},
		{
			name: "failure - not retryable because there are many plugin retries already",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(25880526),
					ExtraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
				},
				UpkeepId: upkeepId,
			},
			upkeepType:            automationTypes.LogTrigger,
			pluginRetries:         6,
			mockHttpStatusCode:    http.StatusBadGateway,
			mockChainlinkBlobs:    []string{"0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"},
			expectedValues:        [][]byte(nil),
			expectedRetryInterval: 0 * time.Second,
			expectedErrCode:       encoding.ErrCodeStreamsBadGateway,
			expectedRetryable:     false,
			expectedState:         encoding.NoPipelineError,
		},
		{
			name: "failure - not retryable with an error code",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{"0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(25880526),
					ExtraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
				},
				UpkeepId: upkeepId,
			},
			mockHttpStatusCode: http.StatusTooManyRequests,
			mockChainlinkBlobs: []string{"0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"},
			expectedValues:     [][]byte(nil),
			expectedErrCode:    encoding.ErrCodeStreamsUnknownError,
			expectedRetryable:  false,
			expectedState:      encoding.NoPipelineError,
		},
		{
			name: "failure - no feeds",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIdHex,
					Feeds:        []string{},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(25880526),
					ExtraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
				},
				UpkeepId: upkeepId,
			},
			expectedValues:  [][]byte(nil),
			expectedErrCode: encoding.ErrCodeStreamsBadRequest,
		},
		{
			name: "failure - invalid revert data",
			lookup: &mercury.StreamsLookup{
				StreamsLookupError: &mercury.StreamsLookupError{
					FeedParamKey: mercury.FeedIDs,
					Feeds:        []string{},
					TimeParamKey: mercury.BlockNumber,
					Time:         big.NewInt(25880526),
					ExtraData:    []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 100},
				},
				UpkeepId: upkeepId,
			},
			expectedValues:  [][]byte(nil),
			expectedErrCode: encoding.ErrCodeStreamsBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := setupClient(t)
			defer c.Close()
			if tt.pluginRetries != 0 {
				c.mercuryConfig.SetPluginRetry(pluginRetryKey, tt.pluginRetries, cache.DefaultExpiration)
			}
			hc := new(MockHttpClient)

			for _, blob := range tt.mockChainlinkBlobs {
				mr := MercuryV02Response{ChainlinkBlob: blob}
				b, err := json.Marshal(mr)
				assert.NoError(t, err)

				resp := &http.Response{
					StatusCode: tt.mockHttpStatusCode,
					Body:       io.NopCloser(bytes.NewReader(b)),
				}
				if tt.expectedErrCode != encoding.ErrCodeNil {
					hc.On("Do", mock.Anything).Return(resp, nil).Times(totalAttempt)
				} else {
					hc.On("Do", mock.Anything).Return(resp, nil).Once()
				}
			}
			c.httpClient = hc

			state, values, errCode, retryable, retryInterval, reqErr := c.DoRequest(testutils.Context(t), tt.lookup, tt.upkeepType, pluginRetryKey)
			assert.Equal(t, tt.expectedValues, values)
			assert.Equal(t, tt.expectedRetryable, retryable)
			if retryable {
				newRetries, _ := c.mercuryConfig.GetPluginRetry(pluginRetryKey)
				assert.Equal(t, tt.pluginRetries+1, newRetries.(int))
			}
			assert.Equal(t, tt.expectedRetryInterval, retryInterval)
			assert.Equal(t, tt.expectedErrCode, errCode)
			assert.Equal(t, tt.expectedState, state)

			if tt.expectedError != nil {
				assert.True(t, strings.HasPrefix(reqErr.Error(), "failed to request feed for 0x4554482d5553442d415242495452554d2d544553544e45540000000000000000"))
			}
		})
	}
}

func TestV02_DoMercuryRequestV02_MultipleFeedsSuccess(t *testing.T) {
	t.Parallel()
	upkeepId, _ := new(big.Int).SetString("88786950015966611018675766524283132478093844178961698330929478019253453382042", 10)
	pluginRetryKey := "88786950015966611018675766524283132478093844178961698330929478019253453382042|34"

	c := setupClient(t)
	defer c.Close()

	c.mercuryConfig.SetPluginRetry(pluginRetryKey, 0, cache.DefaultExpiration)
	hc := new(MockHttpClient)

	for i := 0; i <= 3; i++ {
		mr := MercuryV02Response{ChainlinkBlob: "0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"}
		b, err := json.Marshal(mr)
		assert.NoError(t, err)

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
	assert.False(t, retryable)
	assert.Equal(t, 0*time.Second, retryInterval)
	assert.Equal(t, encoding.ErrCodeNil, errCode)
	assert.Equal(t, encoding.NoPipelineError, state)
}

func TestV02_DoMercuryRequestV02_Timeout(t *testing.T) {
	t.Parallel()
	upkeepId, _ := new(big.Int).SetString("88786950015966611018675766524283132478093844178961698330929478019253453382042", 10)
	pluginRetryKey := "88786950015966611018675766524283132478093844178961698330929478019253453382042|34"

	c := setupClient(t)
	defer c.Close()

	c.mercuryConfig.SetPluginRetry(pluginRetryKey, 0, cache.DefaultExpiration)
	hc := new(MockHttpClient)

	mr := MercuryV02Response{ChainlinkBlob: "0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"}
	b, err := json.Marshal(mr)
	assert.NoError(t, err)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(b)),
	}
	hc.On("Do", mock.Anything).Return(resp, nil).Once() // First request sends result normally

	resp2 := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}
	// Second request has timeout
	serverTimeout := 15 * time.Second // Server has delay of 15s, higher than mercury.RequestTimeout = 10s
	hc.On("Do", mock.Anything).Run(func(args mock.Arguments) {
		time.Sleep(serverTimeout)
	}).Return(resp2, nil).Once()

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

	start := time.Now()
	state, values, errCode, retryable, retryInterval, _ := c.DoRequest(testutils.Context(t), lookup, automationTypes.ConditionTrigger, pluginRetryKey)
	elapsed := time.Since(start)
	assert.Less(t, elapsed, serverTimeout)
	assert.False(t, retryable)
	assert.Equal(t, 0*time.Second, retryInterval)
	assert.Equal(t, encoding.ErrCodeStreamsTimeout, errCode)
	assert.Equal(t, encoding.NoPipelineError, state)
	assert.Equal(t, [][]byte(nil), values)
}

func TestV02_DoMercuryRequestV02_OneFeedSuccessOneFeedPipelineError(t *testing.T) {
	t.Parallel()
	upkeepId, _ := new(big.Int).SetString("88786950015966611018675766524283132478093844178961698330929478019253453382042", 10)
	pluginRetryKey := "88786950015966611018675766524283132478093844178961698330929478019253453382042|34"

	c := setupClient(t)
	defer c.Close()

	c.mercuryConfig.SetPluginRetry(pluginRetryKey, 0, cache.DefaultExpiration)
	hc := new(MockHttpClient)

	// First request success
	mr := MercuryV02Response{ChainlinkBlob: "0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"}
	b, err := json.Marshal(mr)
	assert.NoError(t, err)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(b)),
	}
	hc.On("Do", mock.Anything).Return(resp, nil).Once()
	// Second request returns MercuryFlakyError
	resp = &http.Response{
		StatusCode: http.StatusServiceUnavailable,
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
	assert.True(t, retryable)
	assert.Equal(t, 1*time.Second, retryInterval)
	assert.Equal(t, encoding.ErrCodeStreamsServiceUnavailable, errCode)
	assert.Equal(t, encoding.MercuryFlakyFailure, state)
	assert.Equal(t, [][]byte(nil), values)
}

func TestV02_DoMercuryRequestV02_OneFeedSuccessOneFeedErrCode(t *testing.T) {
	t.Parallel()
	upkeepId, _ := new(big.Int).SetString("88786950015966611018675766524283132478093844178961698330929478019253453382042", 10)
	pluginRetryKey := "88786950015966611018675766524283132478093844178961698330929478019253453382042|34"

	c := setupClient(t)
	defer c.Close()

	c.mercuryConfig.SetPluginRetry(pluginRetryKey, 0, cache.DefaultExpiration)
	hc := new(MockHttpClient)

	// First request success
	mr := MercuryV02Response{ChainlinkBlob: "0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"}
	b, err := json.Marshal(mr)
	assert.NoError(t, err)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(b)),
	}
	hc.On("Do", mock.Anything).Return(resp, nil).Once()
	// Second request returns invalid response
	resp = &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}
	hc.On("Do", mock.Anything).Return(resp, nil).Once()
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
	assert.False(t, retryable)
	assert.Equal(t, 0*time.Second, retryInterval)
	assert.Equal(t, encoding.ErrCodeStreamsBadResponse, errCode)
	assert.Equal(t, encoding.NoPipelineError, state)
}

func TestV02_DoMercuryRequestV02_OneFeedSuccessOneFeedPipelineErrorConvertedError(t *testing.T) {
	t.Parallel()
	upkeepId, _ := new(big.Int).SetString("88786950015966611018675766524283132478093844178961698330929478019253453382042", 10)
	pluginRetryKey := "88786950015966611018675766524283132478093844178961698330929478019253453382042|34"

	c := setupClient(t)
	defer c.Close()

	c.mercuryConfig.SetPluginRetry(pluginRetryKey, 0, cache.DefaultExpiration)
	hc := new(MockHttpClient)

	// First request success
	mr := MercuryV02Response{ChainlinkBlob: "0x00066dfcd1ed2d95b18c948dbc5bd64c687afe93e4ca7d663ddec14c20090ad80000000000000000000000000000000000000000000000000000000000081401000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002200000000000000000000000000000000000000000000000000000000000000280000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001204554482d5553442d415242495452554d2d544553544e455400000000000000000000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000289ad8d367000000000000000000000000000000000000000000000000000000289acf0b38000000000000000000000000000000000000000000000000000000289b3da40000000000000000000000000000000000000000000000000000000000018ae7ce74d9fa252a8983976eab600dc7590c778d04813430841bc6e765c34cd81a168d00000000000000000000000000000000000000000000000000000000018ae7cb0000000000000000000000000000000000000000000000000000000064891c98000000000000000000000000000000000000000000000000000000000000000260412b94e525ca6cedc9f544fd86f77606d52fe731a5d069dbe836a8bfc0fb8c911963b0ae7a14971f3b4621bffb802ef0605392b9a6c89c7fab1df8633a5ade00000000000000000000000000000000000000000000000000000000000000024500c2f521f83fba5efc2bf3effaaedde43d0a4adff785c1213b712a3aed0d8157642a84324db0cf9695ebd27708d4608eb0337e0dd87b0e43f0fa70c700d911"}
	b, err := json.Marshal(mr)
	assert.NoError(t, err)

	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(b)),
	}
	hc.On("Do", mock.Anything).Return(resp, nil).Once()
	// Second request returns MercuryFlakyError
	resp = &http.Response{
		StatusCode: http.StatusGatewayTimeout,
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

	state, values, errCode, retryable, retryInterval, _ := c.DoRequest(testutils.Context(t), lookup, automationTypes.ConditionTrigger, pluginRetryKey)
	assert.False(t, retryable)
	assert.Equal(t, 0*time.Second, retryInterval)
	assert.Equal(t, encoding.ErrCodeStreamsStatusGatewayTimeout, errCode)
	assert.Equal(t, encoding.NoPipelineError, state)
	assert.Equal(t, [][]byte(nil), values)
}
