package observability

import (
	"context"
	"encoding/json"
	"math/big"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	http2 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata/http"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/tokendata/usdc"
)

type expected struct {
	status string
	result string
	count  int
}

func TestUSDCClientMonitoring(t *testing.T) {
	tests := []struct {
		name     string
		server   *httptest.Server
		requests int
		expected []expected
	}{
		{
			name:     "success",
			server:   newSuccessServer(t),
			requests: 5,
			expected: []expected{
				{"200", "true", 5},
				{"429", "false", 0},
			},
		},
		{
			name:     "rate_limited",
			server:   newRateLimitedServer(),
			requests: 26,
			expected: []expected{
				{"200", "true", 0},
				{"429", "false", 1},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testMonitoring(t, test.name, test.server, test.requests, test.expected, logger.TestLogger(t))
		})
	}
}

func testMonitoring(t *testing.T, name string, server *httptest.Server, requests int, expected []expected, log logger.Logger) {
	server.Start()
	defer server.Close()
	attestationURI, err := url.ParseRequestURI(server.URL)
	require.NoError(t, err)

	// Define test histogram (avoid side effects from other tests if using the real usdcHistogram).
	histogram := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "test_client_histogram_" + name,
		Help:    "Latency of calls to the USDC mock client",
		Buckets: []float64{float64(250 * time.Millisecond), float64(1 * time.Second), float64(5 * time.Second)},
	}, []string{"status", "success"})

	// Mock USDC reader.
	usdcReader := mocks.NewUSDCReader(t)
	msgBody := []byte{0xb0, 0xd1}
	usdcReader.On("GetUSDCMessagePriorToLogIndexInTx", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(msgBody, nil)

	// Service with monitored http client.
	usdcTokenAddr := utils.RandomAddress()
	observedHttpClient := http2.NewObservedIHttpClientWithMetric(&http2.HttpClient{}, histogram)
	tokenDataReaderDefault := usdc.NewUSDCTokenDataReader(log, usdcReader, attestationURI, 0, usdcTokenAddr, usdc.APIIntervalRateLimitDisabled)
	tokenDataReader := usdc.NewUSDCTokenDataReaderWithHttpClient(*tokenDataReaderDefault, observedHttpClient, usdcTokenAddr, usdc.APIIntervalRateLimitDisabled)
	require.NotNil(t, tokenDataReader)

	for i := 0; i < requests; i++ {
		_, _ = tokenDataReader.ReadTokenData(context.Background(), cciptypes.EVM2EVMOnRampCCIPSendRequestedWithMeta{
			EVM2EVMMessage: cciptypes.EVM2EVMMessage{
				TokenAmounts: []cciptypes.TokenAmount{
					{
						Token:  ccipcalc.EvmAddrToGeneric(usdcTokenAddr),
						Amount: big.NewInt(rand.Int63()),
					},
				},
			},
		}, 0)
	}

	// Check that the metrics are updated as expected.
	for _, e := range expected {
		assert.Equal(t, e.count, counterFromHistogramByLabels(t, histogram, e.status, e.result))
	}
}

func counterFromHistogramByLabels(t *testing.T, histogramVec *prometheus.HistogramVec, labels ...string) int {
	observer, err := histogramVec.GetMetricWithLabelValues(labels...)
	require.NoError(t, err)

	metricCh := make(chan prometheus.Metric, 1)
	observer.(prometheus.Histogram).Collect(metricCh)
	close(metricCh)

	metric := <-metricCh
	pb := &io_prometheus_client.Metric{}
	err = metric.Write(pb)
	require.NoError(t, err)

	return int(pb.GetHistogram().GetSampleCount())
}

func newSuccessServer(t *testing.T) *httptest.Server {
	return httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		response := struct {
			Status      string `json:"status"`
			Attestation string `json:"attestation"`
		}{
			Status:      "complete",
			Attestation: "720502893578a89a8a87982982ef781c18b193",
		}
		responseBytes, err := json.Marshal(response)
		require.NoError(t, err)
		_, err = w.Write(responseBytes)
		require.NoError(t, err)
	}))
}

func newRateLimitedServer() *httptest.Server {
	return httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
}
