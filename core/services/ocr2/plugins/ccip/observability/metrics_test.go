package observability

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func TestProperLabelsArePassed(t *testing.T) {
	histogram := evm2evmOffRampHistogram
	successCounter := 10
	failedCounter := 5

	details := metricDetails{
		histogram:  histogram,
		pluginName: "plugin",
		chainId:    big.NewInt(123),
	}

	for i := 0; i < successCounter; i++ {
		_, err := withObservedContract[string](details, "successFun", successfulContract)
		require.NoError(t, err)
	}

	for i := 0; i < failedCounter; i++ {
		_, err := withObservedContract[string](details, "failedFun", failedContract)
		require.Error(t, err)
	}

	assert.Equal(t, successCounter, counterFromHistogramByLabels(t, histogram, "123", "plugin", "successFun", "true"))
	assert.Equal(t, failedCounter, counterFromHistogramByLabels(t, histogram, "123", "plugin", "failedFun", "false"))

	assert.Equal(t, 0, counterFromHistogramByLabels(t, histogram, "123", "plugin", "failedFun", "true"))
	assert.Equal(t, 0, counterFromHistogramByLabels(t, histogram, "123", "plugin", "successFun", "false"))
}

func TestMetricsSendFromContractDirectly(t *testing.T) {
	expectedCounter := 4
	evmClient := mocks.NewClient(t)
	evmClient.On("ConfiguredChainID").Return(big.NewInt(420), nil)
	evmClient.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return([]byte{}, fmt.Errorf("error"))

	ramp, err := NewObservedEvm2EvmOffRamp(common.HexToAddress("0xa"), "plugin", evmClient)
	require.NoError(t, err)

	for i := 0; i < expectedCounter; i++ {
		_, _ = ramp.GetSupportedTokens(&bind.CallOpts{Context: testutils.Context(t)})
		_, _ = ramp.CurrentRateLimiterState(&bind.CallOpts{Context: testutils.Context(t)})
	}

	assert.Equal(t, expectedCounter, counterFromHistogramByLabels(t, ramp.metric.histogram, "420", "plugin", "GetSupportedTokens", "false"))
	assert.Equal(t, expectedCounter, counterFromHistogramByLabels(t, ramp.metric.histogram, "420", "plugin", "CurrentRateLimiterState", "false"))
	assert.Equal(t, 0, counterFromHistogramByLabels(t, ramp.metric.histogram, "420", "plugin", "GetDestinationTokens", "false"))
	assert.Equal(t, 0, counterFromHistogramByLabels(t, ramp.metric.histogram, "420", "plugin", "GetDestinationTokens", "true"))
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

func successfulContract() (string, error) {
	return "success", nil
}

func failedContract() (string, error) {
	return "", fmt.Errorf("just error")
}
