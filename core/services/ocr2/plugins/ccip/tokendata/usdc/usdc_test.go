package usdc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/cciptypes"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
)

var (
	mockMsgTransmitter = utils.RandomAddress()
)

func TestUSDCReader_callAttestationApi(t *testing.T) {
	t.Skipf("Skipping test because it uses the real USDC attestation API")
	usdcMessageHash := "912f22a13e9ccb979b621500f6952b2afd6e75be7eadaed93fc2625fe11c52a2"
	attestationURI, err := url.ParseRequestURI("https://iris-api-sandbox.circle.com")
	require.NoError(t, err)
	lggr := logger.TestLogger(t)
	usdcReader, _ := ccipdata.NewUSDCReader(lggr, "job_123", mockMsgTransmitter, nil, false)
	usdcService := NewUSDCTokenDataReader(lggr, usdcReader, attestationURI, 0)

	attestation, err := usdcService.callAttestationApi(context.Background(), [32]byte(common.FromHex(usdcMessageHash)))
	require.NoError(t, err)

	require.Equal(t, attestationStatusPending, attestation.Status)
	require.Equal(t, "PENDING", attestation.Attestation)
}

func TestUSDCReader_callAttestationApiMock(t *testing.T) {
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
	usdcReader, _ := ccipdata.NewUSDCReader(lggr, "job_123", mockMsgTransmitter, lp, false)
	usdcService := NewUSDCTokenDataReader(lggr, usdcReader, attestationURI, 0)
	attestation, err := usdcService.callAttestationApi(context.Background(), utils.RandomBytes32())
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
			ts := test.getTs()
			defer ts.Close()

			attestationURI, err := url.ParseRequestURI(ts.URL)
			require.NoError(t, err)

			lggr := logger.TestLogger(t)
			lp := mocks.NewLogPoller(t)
			usdcReader, _ := ccipdata.NewUSDCReader(lggr, "job_123", mockMsgTransmitter, lp, false)
			usdcService := NewUSDCTokenDataReader(lggr, usdcReader, attestationURI, test.customTimeoutSeconds)
			lp.On("RegisterFilter", mock.Anything).Return(nil)
			require.NoError(t, usdcReader.RegisterFilters())

			parentCtx, cancel := context.WithTimeout(context.Background(), time.Duration(test.parentTimeoutSeconds)*time.Second)
			defer cancel()

			_, err = usdcService.callAttestationApi(parentCtx, utils.RandomBytes32())
			require.Error(t, err)

			if test.expectedError != nil {
				require.True(t, errors.Is(err, test.expectedError))
			}
			lp.On("UnregisterFilter", mock.Anything).Return(nil)
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
	expectedBody := []byte("0x0000000000000001000000020000000000048d71000000000000000000000000eb08f243e5d3fcff26a9e38ae5520a669f4019d000000000000000000000000023a04d5935ed8bc8e3eb78db3541f0abfb001c6e0000000000000000000000006cb3ed9b441eb674b58495c8b3324b59faff5243000000000000000000000000000000005425890298aed601595a70ab815c96711a31bc65000000000000000000000000ab4f961939bfe6a93567cc57c59eed7084ce2131000000000000000000000000000000000000000000000000000000000000271000000000000000000000000035e08285cfed1ef159236728f843286c55fc0861")
	usdcReader := ccipdatamocks.USDCReader{}
	usdcReader.On("GetLastUSDCMessagePriorToLogIndexInTx", mock.Anything, mock.Anything, mock.Anything).Return(expectedBody, nil)

	lggr := logger.TestLogger(t)
	usdcService := NewUSDCTokenDataReader(lggr, &usdcReader, nil, 0)

	// Make the first call and assert the underlying function is called
	body, err := usdcService.getUSDCMessageBody(context.Background(), cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{})
	require.NoError(t, err)
	require.Equal(t, body, expectedBody)

	usdcReader.AssertNumberOfCalls(t, "GetLastUSDCMessagePriorToLogIndexInTx", 1)
}
