package observability

import (
	"errors"
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
		readerName:          "reader",
		chainId:             123,
	}

	for range successCounter {
		_, err := withObservedInteraction[string](details, "successFun", successfulContract)
		require.NoError(t, err)
	}

	for range failedCounter {
		_, err := withObservedInteraction[string](details, "failedFun", failedContract)
		require.Error(t, err)
	}

	assert.Equal(t, successCounter, counterFromHistogramByLabels(t, histogram, "123", "reader", "successFun"))
	assert.Equal(t, failedCounter, counterFromHistogramByLabels(t, histogram, "123", "reader", "failedFun"))
}

func TestMetricsSendFromContractDirectly(t *testing.T) {
	expectedCounter := 4
	ctx := testutils.Context(t)
	chainId := int64(420)

	mockedOfframp := ccipdatamocks.NewOffRampReader(t)
	mockedOfframp.On("GetTokens", ctx).Return(cciptypes.OffRampTokens{}, errors.New("execution error"))

	observedOfframp := NewObservedOffRampReader(mockedOfframp, chainId, "plugin")

	for range expectedCounter {
		_, _ = observedOfframp.GetTokens(ctx)
	}

	assert.Equal(t, expectedCounter, counterFromHistogramByLabels(t, observedOfframp.metric.interactionDuration, "420", "OffRampReader", "GetTokens"))
	assert.Equal(t, 0, counterFromHistogramByLabels(t, observedOfframp.metric.interactionDuration, "420", "OffRampReader", "GetPoolByDestToken"))
	assert.Equal(t, 0, counterFromHistogramByLabels(t, observedOfframp.metric.interactionDuration, "420", "OffRampReader", "GetPoolByDestToken"))
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
	return "", errors.New("just error")
}
