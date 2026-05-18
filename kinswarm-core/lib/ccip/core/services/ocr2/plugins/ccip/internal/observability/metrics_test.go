package observability

import (
	"fmt"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
)

func TestProperLabelsArePassed(t *testing.T) {
	histogram := readerHistogram
	successCounter := 10
	failedCounter := 5

	details := metricDetails{
		interactionDuration: histogram,
		pluginName:          "plugin",
		readerName:          "reader",
		chainId:             123,
	}

	for i := 0; i < successCounter; i++ {
		_, err := withObservedInteraction[string](details, "successFun", successfulContract)
		require.NoError(t, err)
	}

	for i := 0; i < failedCounter; i++ {
		_, err := withObservedInteraction[string](details, "failedFun", failedContract)
		require.Error(t, err)
	}

	assert.Equal(t, successCounter, counterFromHistogramByLabels(t, histogram, "123", "plugin", "reader", "successFun", "true"))
	assert.Equal(t, failedCounter, counterFromHistogramByLabels(t, histogram, "123", "plugin", "reader", "failedFun", "false"))

	assert.Equal(t, 0, counterFromHistogramByLabels(t, histogram, "123", "plugin", "reader", "failedFun", "true"))
	assert.Equal(t, 0, counterFromHistogramByLabels(t, histogram, "123", "plugin", "reader", "successFun", "false"))
}

func TestMetricsSendFromContractDirectly(t *testing.T) {
	expectedCounter := 4
	ctx := testutils.Context(t)
	chainId := int64(420)

	mockedOfframp := ccipdatamocks.NewOffRampReader(t)
	mockedOfframp.On("GetTokens", ctx).Return(cciptypes.OffRampTokens{}, fmt.Errorf("execution error"))

	observedOfframp := NewObservedOffRampReader(mockedOfframp, chainId, "plugin")

	for i := 0; i < expectedCounter; i++ {
		_, _ = observedOfframp.GetTokens(ctx)
	}

	assert.Equal(t, expectedCounter, counterFromHistogramByLabels(t, observedOfframp.metric.interactionDuration, "420", "plugin", "OffRampReader", "GetTokens", "false"))
	assert.Equal(t, 0, counterFromHistogramByLabels(t, observedOfframp.metric.interactionDuration, "420", "plugin", "OffRampReader", "GetPoolByDestToken", "false"))
	assert.Equal(t, 0, counterFromHistogramByLabels(t, observedOfframp.metric.interactionDuration, "420", "plugin", "OffRampReader", "GetPoolByDestToken", "true"))
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
