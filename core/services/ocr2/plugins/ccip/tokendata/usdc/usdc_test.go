package usdc

import (
	"context"
	"encoding/json"
	"math/big"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
)

var (
	mockMsgTransmitter = utils.RandomAddress()
)

func TestUSDCReader_callAttestationApi(t *testing.T) {
	ctx := tests.Context(t) //nolint:staticcheck // SA4006 - false positive "unused"
	t.Skipf("Skipping test because it uses the real USDC attestation API")
	usdcMessageHash := "912f22a13e9ccb979b621500f6952b2afd6e75be7eadaed93fc2625fe11c52a2"
	attestationURI, err := url.ParseRequestURI("https://iris-api-sandbox.circle.com")
	require.NoError(t, err)
	lggr := logger.TestLogger(t)
	usdcReader, err := ccipdata.NewUSDCReader(ctx, lggr, "job_123", mockMsgTransmitter, nil, false)
	require.NoError(t, err)
	usdcService := NewUSDCTokenDataReader(lggr, usdcReader, attestationURI, 0, common.Address{}, APIIntervalRateLimitDisabled)

	attestation, err := usdcService.callAttestationApi(ctx, [32]byte(common.FromHex(usdcMessageHash)))
	require.NoError(t, err)

	require.Equal(t, attestationStatusPending, attestation.Status)
	require.Equal(t, "PENDING", attestation.Attestation)
}

func TestUSDCReader_callAttestationApiMock(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	response := attestationResponse{
		Status:      attestationStatusSuccess,
		Attestation: "720502893578a89a8a87982982ef781c18b193",
	}

	ts := getMockUSDCEndpoint(t, response)
	defer ts.Close()
	attestationURI, err := url.ParseRequestURI(ts.URL)
	require.NoError(t, err)

	lggr := logger.TestLogger(t)
	lp := mocks.NewLogPoller(t)
	usdcReader, _ := ccipdata.NewUSDCReader(ctx, lggr, "job_123", mockMsgTransmitter, lp, false)
	usdcService := NewUSDCTokenDataReader(lggr, usdcReader, attestationURI, 0, common.Address{}, APIIntervalRateLimitDisabled)
	attestation, err := usdcService.callAttestationApi(ctx, utils.RandomBytes32())
	require.NoError(t, err)

	require.Equal(t, response.Status, attestation.Status)
	require.Equal(t, response.Attestation, attestation.Attestation)
}

func TestUSDCReader_callAttestationApiMockError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		getTs                func() *httptest.Server
		parentTimeoutSeconds int
		customTimeoutSeconds int
		expectedError        error
	}{
		{
			name: "server error",
			getTs: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
			},
			parentTimeoutSeconds: 60,
			expectedError:        nil,
		},
		{
			name: "default timeout",
			getTs: func() *httptest.Server {
				response := attestationResponse{
					Status:      attestationStatusSuccess,
					Attestation: "720502893578a89a8a87982982ef781c18b193",
				}
				responseBytes, _ := json.Marshal(response)

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(defaultAttestationTimeout + time.Second)
					_, err := w.Write(responseBytes)
					require.NoError(t, err)
				}))
			},
			parentTimeoutSeconds: 60,
			expectedError:        tokendata.ErrTimeout,
		},
		{
			name: "custom timeout",
			getTs: func() *httptest.Server {
				response := attestationResponse{
					Status:      attestationStatusSuccess,
					Attestation: "720502893578a89a8a87982982ef781c18b193",
				}
				responseBytes, _ := json.Marshal(response)

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(2*time.Second + time.Second)
					_, err := w.Write(responseBytes)
					require.NoError(t, err)
				}))
			},
			parentTimeoutSeconds: 60,
			customTimeoutSeconds: 2,
			expectedError:        tokendata.ErrTimeout,
		},
		{
			name: "error response",
			getTs: func() *httptest.Server {
				response := attestationResponse{
					Error: "some error",
				}
				responseBytes, _ := json.Marshal(response)

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, err := w.Write(responseBytes)
					require.NoError(t, err)
				}))
			},
			parentTimeoutSeconds: 60,
			expectedError:        nil,
		},
		{
			name: "invalid status",
			getTs: func() *httptest.Server {
				response := attestationResponse{
					Status:      "",
					Attestation: "720502893578a89a8a87982982ef781c18b193",
				}
				responseBytes, _ := json.Marshal(response)

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, err := w.Write(responseBytes)
					require.NoError(t, err)
				}))
			},
			parentTimeoutSeconds: 60,
			expectedError:        nil,
		},
		{
			name: "rate limit",
			getTs: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusTooManyRequests)
				}))
			},
			parentTimeoutSeconds: 60,
			expectedError:        tokendata.ErrRateLimit,
		},
		{
			name: "parent context timeout",
			getTs: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(defaultAttestationTimeout + time.Second)
				}))
			},
			parentTimeoutSeconds: 1,
			expectedError:        nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ts := test.getTs()
			defer ts.Close()

			attestationURI, err := url.ParseRequestURI(ts.URL)
			require.NoError(t, err)

			lggr := logger.TestLogger(t)
			lp := mocks.NewLogPoller(t)
			ctx := testutils.Context(t)
			usdcReader, _ := ccipdata.NewUSDCReader(ctx, lggr, "job_123", mockMsgTransmitter, lp, false)
			usdcService := NewUSDCTokenDataReader(lggr, usdcReader, attestationURI, test.customTimeoutSeconds, common.Address{}, APIIntervalRateLimitDisabled)
			lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil)
			require.NoError(t, usdcReader.RegisterFilters(ctx))

			parentCtx, cancel := context.WithTimeout(ctx, time.Duration(test.parentTimeoutSeconds)*time.Second)
			defer cancel()

			_, err = usdcService.callAttestationApi(parentCtx, utils.RandomBytes32())
			require.Error(t, err)

			if test.expectedError != nil {
				require.True(t, errors.Is(err, test.expectedError))
			}
			lp.On("UnregisterFilter", mock.Anything, mock.Anything).Return(nil)
			require.NoError(t, usdcReader.Close())
		})
	}
}

func getMockUSDCEndpoint(t *testing.T, response attestationResponse) *httptest.Server {
	responseBytes, err := json.Marshal(response)
	require.NoError(t, err)

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(responseBytes)
		require.NoError(t, err)
	}))
}

func TestGetUSDCMessageBody(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	expectedBody := []byte("0x0000000000000001000000020000000000048d71000000000000000000000000eb08f243e5d3fcff26a9e38ae5520a669f4019d000000000000000000000000023a04d5935ed8bc8e3eb78db3541f0abfb001c6e0000000000000000000000006cb3ed9b441eb674b58495c8b3324b59faff5243000000000000000000000000000000005425890298aed601595a70ab815c96711a31bc65000000000000000000000000ab4f961939bfe6a93567cc57c59eed7084ce2131000000000000000000000000000000000000000000000000000000000000271000000000000000000000000035e08285cfed1ef159236728f843286c55fc0861")
	usdcReader := ccipdatamocks.USDCReader{}
	usdcReader.On("GetUSDCMessagePriorToLogIndexInTx", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(expectedBody, nil)

	usdcTokenAddr := utils.RandomAddress()
	lggr := logger.TestLogger(t)
	usdcService := NewUSDCTokenDataReader(lggr, &usdcReader, nil, 0, usdcTokenAddr, APIIntervalRateLimitDisabled)

	// Make the first call and assert the underlying function is called
	body, err := usdcService.getUSDCMessageBody(ctx, cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
		EVM2EVMMessage: cciptypes.EVM2EVMMessage{
			TokenAmounts: []cciptypes.TokenAmount{
				{
					Token:  ccipcalc.EvmAddrToGeneric(usdcTokenAddr),
					Amount: big.NewInt(rand.Int63()),
				},
			},
		},
	}, 0)
	require.NoError(t, err)
	require.Equal(t, expectedBody, body)

	usdcReader.AssertNumberOfCalls(t, "GetUSDCMessagePriorToLogIndexInTx", 1)
}

func TestTokenDataReader_getUsdcTokenEndOffset(t *testing.T) {
	t.Parallel()
	usdcToken := utils.RandomAddress()
	nonUsdcToken := utils.RandomAddress()

	multipleTokens := []common.Address{
		usdcToken, // 2
		nonUsdcToken,
		nonUsdcToken,
		usdcToken, // 1
		usdcToken, // 0
		nonUsdcToken,
	}

	testCases := []struct {
		name       string
		tokens     []common.Address
		tokenIndex int
		expOffset  int
		expErr     bool
	}{
		{name: "one non usdc token", tokens: []common.Address{nonUsdcToken}, tokenIndex: 0, expOffset: 0, expErr: true},
		{name: "one usdc token", tokens: []common.Address{usdcToken}, tokenIndex: 0, expOffset: 0, expErr: false},
		{name: "one usdc token wrong index", tokens: []common.Address{usdcToken}, tokenIndex: 1, expOffset: 0, expErr: true},
		{name: "multiple tokens 1", tokens: multipleTokens, tokenIndex: 0, expOffset: 2},
		{name: "multiple tokens - non usdc selected", tokens: multipleTokens, tokenIndex: 2, expErr: true},
		{name: "multiple tokens 2", tokens: multipleTokens, tokenIndex: 3, expOffset: 1},
		{name: "multiple tokens 3", tokens: multipleTokens, tokenIndex: 4, expOffset: 0},
		{name: "multiple tokens not found", tokens: multipleTokens, tokenIndex: 5, expErr: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := &TokenDataReader{usdcTokenAddress: usdcToken}
			tokenAmounts := make([]cciptypes.TokenAmount, len(tc.tokens))
			for i := range tokenAmounts {
				tokenAmounts[i] = cciptypes.TokenAmount{
					Token:  ccipcalc.EvmAddrToGeneric(tc.tokens[i]),
					Amount: big.NewInt(rand.Int63()),
				}
			}
			msg := cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{EVM2EVMMessage: cciptypes.EVM2EVMMessage{TokenAmounts: tokenAmounts}}
			offset, err := r.getUsdcTokenEndOffset(msg, tc.tokenIndex)
			if tc.expErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expOffset, offset)
		})
	}
}

func TestUSDCReader_rateLimiting(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		requests     uint64
		rateConfig   time.Duration
		testDuration time.Duration
		timeout      time.Duration
		err          string
	}{
		{
			name:         "no rate limit when disabled",
			requests:     10,
			rateConfig:   APIIntervalRateLimitDisabled,
			testDuration: 1 * time.Millisecond,
		},
		{
			name:         "yes rate limited with default config",
			requests:     5,
			rateConfig:   APIIntervalRateLimitDefault,
			testDuration: 4 * defaultRequestInterval,
		},
		{
			name:         "yes rate limited with config",
			requests:     10,
			rateConfig:   50 * time.Millisecond,
			testDuration: 9 * 50 * time.Millisecond,
		},
		{
			name:         "timeout after first request",
			requests:     5,
			rateConfig:   100 * time.Millisecond,
			testDuration: 1 * time.Millisecond,
			timeout:      1 * time.Millisecond,
			err:          "usdc rate limiting error:",
		},
		{
			name:         "timeout after second request",
			requests:     5,
			rateConfig:   100 * time.Millisecond,
			testDuration: 100 * time.Millisecond,
			timeout:      150 * time.Millisecond,
			err:          "usdc rate limiting error:",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := tests.Context(t)

			response := attestationResponse{
				Status:      attestationStatusSuccess,
				Attestation: "720502893578a89a8a87982982ef781c18b193",
			}

			ts := getMockUSDCEndpoint(t, response)
			defer ts.Close()
			attestationURI, err := url.ParseRequestURI(ts.URL)
			require.NoError(t, err)

			lggr := logger.TestLogger(t)
			lp := mocks.NewLogPoller(t)
			usdcReader, _ := ccipdata.NewUSDCReader(ctx, lggr, "job_123", mockMsgTransmitter, lp, false)
			usdcService := NewUSDCTokenDataReader(lggr, usdcReader, attestationURI, 0, utils.RandomAddress(), tc.rateConfig)

			if tc.timeout > 0 {
				var cf context.CancelFunc
				ctx, cf = context.WithTimeout(ctx, tc.timeout)
				defer cf()
			}

			trigger := make(chan struct{})
			errorChan := make(chan error, tc.requests)
			wg := sync.WaitGroup{}
			for i := 0; i < int(tc.requests); i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					<-trigger
					_, err := usdcService.ReadTokenData(ctx, cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
						EVM2EVMMessage: cciptypes.EVM2EVMMessage{
							TokenAmounts: []cciptypes.TokenAmount{{Token: ccipcalc.EvmAddrToGeneric(utils.ZeroAddress), Amount: nil}}, // trigger failure due to wrong address
						},
					}, 0)

					errorChan <- err
				}()
			}

			// Start the test
			start := time.Now()
			close(trigger)

			// Wait for requests to complete
			wg.Wait()
			finish := time.Now()
			close(errorChan)

			// Collect errors
			errorFound := false
			for err := range errorChan {
				if tc.err != "" && strings.Contains(err.Error(), tc.err) {
					errorFound = true
				} else if err != nil && !strings.Contains(err.Error(), "get usdc token 0 end offset") {
					// Ignore that one error, it's expected because of how mocking is used.
					// Anything else is unexpected.
					require.Fail(t, "unexpected error", err)
				}
			}

			if tc.err != "" {
				assert.True(t, errorFound)
			}
			assert.WithinDuration(t, start.Add(tc.testDuration), finish, 50*time.Millisecond)
		})
	}
}
