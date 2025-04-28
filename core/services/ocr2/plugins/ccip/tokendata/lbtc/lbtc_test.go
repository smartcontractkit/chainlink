package lbtc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata"
)

var (
	lbtcMessageHash        = "0xbc427abf571a5cfcf7c98799d1f0055f4db25f203f657d30026728a19d16f092"
	lbtcMessageAttestation = "0x0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000000e45c70a5050000000000000000000000000000000000000000000000000000000000aa36a7000000000000000000000000845f8e3c214d8d0e4d83fc094f302aa26a12a0bc0000000000000000000000000000000000000000000000000000000000014a34000000000000000000000000845f8e3c214d8d0e4d83fc094f302aa26a12a0bc00000000000000000000000062f10ce5b727edf787ea45776bd050308a61150800000000000000000000000000000000000000000000000000000000000003e60000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000001200000000000000000000000000000000000000000000000000000000000000040277eeafba008d767c2636d9428f2ebb13ab29ac70337f4fc34b0f5606767cae546f9be3f12160de6d142e5b3c1c3ebd0bf4298662b32b597d0cc5970c7742fc10000000000000000000000000000000000000000000000000000000000000040bbcd60ecc9e06f2effe7c94161219498a1eb435b419387adadb86ec9a52dfb066ce027532517df7216404049d193a25b85c35edfa3e7c5aa4757bfe84887a3980000000000000000000000000000000000000000000000000000000000000040da4a6dc619b5ca2349783cabecc4efdbc910090d3e234d7b8d0430165f8fae532f9a965ceb85c18bb92e059adefa7ce5835850a705761ab9e026d2db4a13ef9a"
	payloadAndProof, _     = hexutil.Decode(lbtcMessageAttestation)
)

func getMockLBTCEndpoint(t *testing.T, response attestationResponse) *httptest.Server {
	responseBytes, err := json.Marshal(response)
	require.NoError(t, err)

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(responseBytes)
		//nolint:testifylint // we need to use require here
		require.NoError(t, err)
	}))
}

func TestLBTCReader_callAttestationApi(t *testing.T) {
	t.Skipf("Skipping test because it uses the real LBTC attestation API")
	attestationURI, err := url.ParseRequestURI("https://bridge-manager.staging.lombard.finance")
	require.NoError(t, err)
	lggr := logger.TestLogger(t)
	lbtcService := NewLBTCTokenDataReader(lggr, attestationURI, 0, common.Address{}, APIIntervalRateLimitDisabled)

	attestation, err := lbtcService.callAttestationAPI(context.Background(), [32]byte(common.FromHex(lbtcMessageHash)))
	require.NoError(t, err)

	require.Equal(t, lbtcMessageHash, attestation.Attestations[0].MessageHash)
	require.Equal(t, attestationStatusSessionApproved, attestation.Attestations[0].Status)
	require.Equal(t, lbtcMessageAttestation, attestation.Attestations[0].Attestation)
}

func TestLBTCReader_callAttestationApiMock(t *testing.T) {
	response := attestationResponse{
		Attestations: []messageAttestationResponse{
			{
				MessageHash: lbtcMessageHash,
				Status:      attestationStatusSessionApproved,
				Attestation: lbtcMessageAttestation,
			},
		},
	}

	ts := getMockLBTCEndpoint(t, response)
	defer ts.Close()
	attestationURI, err := url.ParseRequestURI(ts.URL)
	require.NoError(t, err)

	lggr := logger.TestLogger(t)
	lbtcService := NewLBTCTokenDataReader(lggr, attestationURI, 0, common.Address{}, APIIntervalRateLimitDisabled)
	attestation, err := lbtcService.callAttestationAPI(context.Background(), [32]byte(common.FromHex(lbtcMessageHash)))
	require.NoError(t, err)

	require.Equal(t, response.Attestations[0].Status, attestation.Attestations[0].Status)
	require.Equal(t, response.Attestations[0].Attestation, attestation.Attestations[0].Attestation)
}

func TestLBTCReader_callAttestationApiMockError(t *testing.T) {
	t.Parallel()

	sessionApprovedResponse := attestationResponse{
		Attestations: []messageAttestationResponse{
			{
				MessageHash: lbtcMessageHash,
				Status:      attestationStatusSessionApproved,
				Attestation: lbtcMessageAttestation,
			},
		},
	}

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
				responseBytes, _ := json.Marshal(sessionApprovedResponse)

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(defaultAttestationTimeout + time.Second)
					_, err := w.Write(responseBytes)
					//nolint:testifylint // we need to use require here
					require.NoError(t, err)
				}))
			},
			parentTimeoutSeconds: 60,
			expectedError:        tokendata.ErrTimeout,
		},
		{
			name: "custom timeout",
			getTs: func() *httptest.Server {
				responseBytes, _ := json.Marshal(sessionApprovedResponse)

				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(2*time.Second + time.Second)
					_, err := w.Write(responseBytes)
					//nolint:testifylint // we need to use require here
					require.NoError(t, err)
				}))
			},
			parentTimeoutSeconds: 60,
			customTimeoutSeconds: 2,
			expectedError:        tokendata.ErrTimeout,
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
			lbtcService := NewLBTCTokenDataReader(lggr, attestationURI, test.customTimeoutSeconds, common.Address{}, APIIntervalRateLimitDisabled)

			parentCtx, cancel := context.WithTimeout(context.Background(), time.Duration(test.parentTimeoutSeconds)*time.Second)
			defer cancel()

			_, err = lbtcService.callAttestationAPI(parentCtx, [32]byte(common.FromHex(lbtcMessageHash)))
			require.Error(t, err)

			if test.expectedError != nil {
				require.True(t, errors.Is(err, test.expectedError))
			}
		})
	}
}

func TestLBTCReader_rateLimiting(t *testing.T) {
	sessionApprovedResponse := attestationResponse{
		Attestations: []messageAttestationResponse{
			{
				MessageHash: lbtcMessageHash,
				Status:      attestationStatusSessionApproved,
				Attestation: lbtcMessageAttestation,
			},
		},
	}

	testCases := []struct {
		name          string
		requests      uint64
		rateConfig    time.Duration
		testDuration  time.Duration
		timeout       time.Duration
		err           string
		additionalErr string
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
			name:          "request timeout",
			requests:      5,
			rateConfig:    100 * time.Millisecond,
			testDuration:  1 * time.Millisecond,
			timeout:       1 * time.Millisecond,
			err:           "lbtc rate limiting error:",
			additionalErr: "token data API timed out",
		},
	}

	extraData, err := hexutil.Decode(lbtcMessageHash)
	require.NoError(t, err)

	srcTokenData, err := abihelpers.EncodeAbiStruct[sourceTokenData](sourceTokenData{
		SourcePoolAddress: utils.RandomAddress().Bytes(),
		DestTokenAddress:  utils.RandomAddress().Bytes(),
		ExtraData:         extraData,
	})
	require.NoError(t, err)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ts := getMockLBTCEndpoint(t, sessionApprovedResponse)
			defer ts.Close()
			attestationURI, err := url.ParseRequestURI(ts.URL)
			require.NoError(t, err)

			lggr := logger.TestLogger(t)
			lbtcService := NewLBTCTokenDataReader(lggr, attestationURI, 0, utils.RandomAddress(), tc.rateConfig)

			ctx := context.Background()
			if tc.timeout > 0 {
				var cf context.CancelFunc
				ctx, cf = context.WithTimeout(ctx, tc.timeout)
				defer cf()
			}

			trigger := make(chan struct{})
			errorChan := make(chan error, tc.requests)
			wg := sync.WaitGroup{}
			for i := uint64(0); i < tc.requests; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					<-trigger
					_, err := lbtcService.ReadTokenData(ctx, cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
						EVM2EVMMessage: cciptypes.EVM2EVMMessage{
							SourceTokenData: [][]byte{srcTokenData},
							TokenAmounts:    []cciptypes.TokenAmount{{Token: ccipcalc.EvmAddrToGeneric(utils.ZeroAddress), Amount: nil}}, // trigger failure due to wrong address
						},
					}, 0)

					if err != nil {
						errorChan <- err
					}
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
				//nolint:gocritic // easier to read using ifElse instead of switch
				if tc.err != "" && strings.Contains(err.Error(), tc.err) {
					errorFound = true
				} else if tc.additionalErr != "" && strings.Contains(err.Error(), tc.additionalErr) {
					errorFound = true
				} else if err != nil {
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

func TestLBTCReader_skipApiOnFullPayload(t *testing.T) {
	sessionApprovedResponse := attestationResponse{
		Attestations: []messageAttestationResponse{
			{
				MessageHash: lbtcMessageHash,
				Status:      attestationStatusSessionApproved,
				Attestation: lbtcMessageAttestation,
			},
		},
	}

	srcTokenData, err := abihelpers.EncodeAbiStruct[sourceTokenData](sourceTokenData{
		SourcePoolAddress: utils.RandomAddress().Bytes(),
		DestTokenAddress:  utils.RandomAddress().Bytes(),
		ExtraData:         []byte(lbtcMessageHash), // more than 32 bytes
	})
	require.NoError(t, err)

	ts := getMockLBTCEndpoint(t, sessionApprovedResponse)
	defer ts.Close()
	attestationURI, err := url.ParseRequestURI(ts.URL)
	require.NoError(t, err)

	lggr, logs := logger.TestLoggerObserved(t, zapcore.InfoLevel)
	lbtcService := NewLBTCTokenDataReader(lggr, attestationURI, 0, utils.RandomAddress(), APIIntervalRateLimitDefault)

	ctx := context.Background()

	destTokenData, err := lbtcService.ReadTokenData(ctx, cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
		EVM2EVMMessage: cciptypes.EVM2EVMMessage{
			SourceTokenData: [][]byte{srcTokenData},
			TokenAmounts:    []cciptypes.TokenAmount{{Token: ccipcalc.EvmAddrToGeneric(utils.ZeroAddress), Amount: nil}}, // trigger failure due to wrong address
		},
	}, 0)
	require.NoError(t, err)
	require.EqualValues(t, []byte(lbtcMessageHash), destTokenData)

	require.Equal(t, 1, logs.Len())
	require.Contains(t, logs.All()[0].Message, "SourceTokenData.extraData size is not 32. This is deposit payload, not sha256(payload). Attestation is disabled onchain")
}

func TestLBTCReader_expectedOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		response       attestationResponse
		expectedReturn []byte
		expectedError  error
	}{
		{
			name: "expected payloadAndProof when status SESSION_APPROVED",
			response: attestationResponse{
				Attestations: []messageAttestationResponse{
					{
						MessageHash: lbtcMessageHash,
						Status:      attestationStatusSessionApproved,
						Attestation: lbtcMessageAttestation,
					},
				},
			},
			expectedReturn: payloadAndProof,
			expectedError:  nil,
		},
		{
			name: "expected ErrNotReady on status PENDING",
			response: attestationResponse{
				Attestations: []messageAttestationResponse{
					{
						MessageHash: lbtcMessageHash,
						Status:      attestationStatusPending,
						Attestation: lbtcMessageAttestation,
					},
				},
			},
			expectedReturn: nil,
			expectedError:  tokendata.ErrNotReady,
		},
		{
			name: "expected ErrNotReady on status SUBMITTED",
			response: attestationResponse{
				Attestations: []messageAttestationResponse{
					{
						MessageHash: lbtcMessageHash,
						Status:      attestationStatusSubmitted,
						Attestation: lbtcMessageAttestation,
					},
				},
			},
			expectedReturn: nil,
			expectedError:  tokendata.ErrNotReady,
		},
		{
			name: "expected ErrUnknownResponse on status UNSPECIFIED",
			response: attestationResponse{
				Attestations: []messageAttestationResponse{
					{
						MessageHash: lbtcMessageHash,
						Status:      attestationStatusUnspecified,
						Attestation: lbtcMessageAttestation,
					},
				},
			},
			expectedReturn: nil,
			expectedError:  ErrUnknownResponse,
		},
		{
			name: "expected ErrUnknownResponse on status FAILED",
			response: attestationResponse{
				Attestations: []messageAttestationResponse{
					{
						MessageHash: lbtcMessageHash,
						Status:      attestationStatusFailed,
						Attestation: lbtcMessageAttestation,
					},
				},
			},
			expectedReturn: nil,
			expectedError:  ErrUnknownResponse,
		},
	}

	extraData, err := hexutil.Decode(lbtcMessageHash)
	require.NoError(t, err)

	srcTokenData, err := abihelpers.EncodeAbiStruct[sourceTokenData](sourceTokenData{
		SourcePoolAddress: utils.RandomAddress().Bytes(),
		DestTokenAddress:  utils.RandomAddress().Bytes(),
		ExtraData:         extraData,
	})
	require.NoError(t, err)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ts := getMockLBTCEndpoint(t, tc.response)
			defer ts.Close()
			attestationURI, err := url.ParseRequestURI(ts.URL)
			require.NoError(t, err)

			lggr := logger.TestLogger(t)
			lbtcService := NewLBTCTokenDataReader(lggr, attestationURI, 0, utils.RandomAddress(), APIIntervalRateLimitDefault)

			ctx := context.Background()

			payloadAndProof, err := lbtcService.ReadTokenData(ctx, cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
				EVM2EVMMessage: cciptypes.EVM2EVMMessage{
					SourceTokenData: [][]byte{srcTokenData},
					TokenAmounts:    []cciptypes.TokenAmount{{Token: ccipcalc.EvmAddrToGeneric(utils.ZeroAddress), Amount: nil}}, // trigger failure due to wrong address
				},
			}, 0)

			if tc.expectedReturn != nil {
				require.EqualValues(t, tc.expectedReturn, payloadAndProof)
			} else if tc.expectedError != nil {
				require.Contains(t, err.Error(), tc.expectedError.Error())
			}
		})
	}
}

func Test_DecodeSourceTokenData(t *testing.T) {
	input, err := hexutil.Decode("0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000249f00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000267d40f64ecc4d95f3e8b2237df5f37b10812c250000000000000000000000000000000000000000000000000000000000000020000000000000000000000000c47e4b3124597fdf8dd07843d4a7052f2ee80c3000000000000000000000000000000000000000000000000000000000000000e45c70a5050000000000000000000000000000000000000000000000000000000000aa36a7000000000000000000000000845f8e3c214d8d0e4d83fc094f302aa26a12a0bc0000000000000000000000000000000000000000000000000000000000014a34000000000000000000000000845f8e3c214d8d0e4d83fc094f302aa26a12a0bc00000000000000000000000062f10ce5b727edf787ea45776bd050308a61150800000000000000000000000000000000000000000000000000000000000003e6000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)
	decoded, err := abihelpers.DecodeAbiStruct[sourceTokenData](input)
	require.NoError(t, err)
	expected, err := hexutil.Decode("0x5c70a5050000000000000000000000000000000000000000000000000000000000aa36a7000000000000000000000000845f8e3c214d8d0e4d83fc094f302aa26a12a0bc0000000000000000000000000000000000000000000000000000000000014a34000000000000000000000000845f8e3c214d8d0e4d83fc094f302aa26a12a0bc00000000000000000000000062f10ce5b727edf787ea45776bd050308a61150800000000000000000000000000000000000000000000000000000000000003e60000000000000000000000000000000000000000000000000000000000000006")
	require.NoError(t, err)
	require.Equal(t, expected, decoded.ExtraData)
}
