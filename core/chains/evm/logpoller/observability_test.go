package logpoller

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

func TestMultipleMetricsArePublished(t *testing.T) {
	ctx := testutils.Context(t)
	lp := createObservedPollLogger(t, 100)
	require.Equal(t, 0, testutil.CollectAndCount(lp.queryDuration))

	_, _ = lp.IndexedLogs(common.Hash{}, common.Address{}, 1, []common.Hash{}, 1, pg.WithParentCtx(ctx))
	_, _ = lp.IndexedLogsByBlockRange(0, 1, common.Hash{}, common.Address{}, 1, []common.Hash{}, pg.WithParentCtx(ctx))
	_, _ = lp.IndexedLogsTopicGreaterThan(common.Hash{}, common.Address{}, 1, common.Hash{}, 1, pg.WithParentCtx(ctx))
	_, _ = lp.IndexedLogsTopicRange(common.Hash{}, common.Address{}, 1, common.Hash{}, common.Hash{}, 1, pg.WithParentCtx(ctx))
	_, _ = lp.IndexedLogsWithSigsExcluding(common.Address{}, common.Hash{}, common.Hash{}, 1, 0, 1, 1, pg.WithParentCtx(ctx))
	_, _ = lp.LogsDataWordRange(common.Hash{}, common.Address{}, 0, common.Hash{}, common.Hash{}, 1, pg.WithParentCtx(ctx))
	_, _ = lp.LogsDataWordGreaterThan(common.Hash{}, common.Address{}, 0, common.Hash{}, 1, pg.WithParentCtx(ctx))
	_, _ = lp.LogsCreatedAfter(common.Hash{}, common.Address{}, time.Now(), 0, pg.WithParentCtx(ctx))
	_, _ = lp.LatestLogByEventSigWithConfs(common.Hash{}, common.Address{}, 0, pg.WithParentCtx(ctx))
	_, _ = lp.LatestLogEventSigsAddrsWithConfs(0, []common.Hash{{}}, []common.Address{{}}, 1, pg.WithParentCtx(ctx))
	_, _ = lp.IndexedLogsCreatedAfter(common.Hash{}, common.Address{}, 0, []common.Hash{}, time.Now(), 0, pg.WithParentCtx(ctx))

	require.Equal(t, 11, testutil.CollectAndCount(lp.queryDuration))
	require.Equal(t, 10, testutil.CollectAndCount(lp.datasetSize))
	resetMetrics(*lp)
}

func TestShouldPublishDurationInCaseOfError(t *testing.T) {
	ctx := testutils.Context(t)
	lp := createObservedPollLogger(t, 200)
	require.Equal(t, 0, testutil.CollectAndCount(lp.queryDuration))

	_, err := lp.LatestLogByEventSigWithConfs(common.Hash{}, common.Address{}, 0, pg.WithParentCtx(ctx))
	require.Error(t, err)

	require.Equal(t, 1, testutil.CollectAndCount(lp.queryDuration))
	require.Equal(t, 1, counterFromHistogramByLabels(t, lp.queryDuration, "200", "LatestLogByEventSigWithConfs"))

	resetMetrics(*lp)
}

func TestNotObservedFunctions(t *testing.T) {
	ctx := testutils.Context(t)
	lp := createObservedPollLogger(t, 300)
	require.Equal(t, 0, testutil.CollectAndCount(lp.queryDuration))

	_, err := lp.Logs(0, 1, common.Hash{}, common.Address{}, pg.WithParentCtx(ctx))
	require.NoError(t, err)

	_, err = lp.LogsWithSigs(0, 1, []common.Hash{{}}, common.Address{}, pg.WithParentCtx(ctx))
	require.NoError(t, err)

	require.Equal(t, 0, testutil.CollectAndCount(lp.queryDuration))
	require.Equal(t, 0, testutil.CollectAndCount(lp.datasetSize))
	resetMetrics(*lp)
}

func TestMetricsAreProperlyPopulatedWithLabels(t *testing.T) {
	lp := createObservedPollLogger(t, 420)
	expectedCount := 9
	expectedSize := 2

	for i := 0; i < expectedCount; i++ {
		_, err := withObservedQueryAndResults(lp, "query", func() ([]string, error) { return []string{"value1", "value2"}, nil })
		require.NoError(t, err)
	}

	require.Equal(t, expectedCount, counterFromHistogramByLabels(t, lp.queryDuration, "420", "query"))
	require.Equal(t, expectedSize, counterFromGaugeByLabels(lp.datasetSize, "420", "query"))

	require.Equal(t, 0, counterFromHistogramByLabels(t, lp.queryDuration, "420", "other_query"))
	require.Equal(t, 0, counterFromHistogramByLabels(t, lp.queryDuration, "5", "query"))

	require.Equal(t, 0, counterFromGaugeByLabels(lp.datasetSize, "420", "other_query"))
	require.Equal(t, 0, counterFromGaugeByLabels(lp.datasetSize, "5", "query"))

	resetMetrics(*lp)
}

func TestNotPublishingDatasetSizeInCaseOfError(t *testing.T) {
	lp := createObservedPollLogger(t, 420)

	_, err := withObservedQueryAndResults(lp, "errorQuery", func() ([]string, error) { return nil, fmt.Errorf("error") })
	require.Error(t, err)

	require.Equal(t, 1, counterFromHistogramByLabels(t, lp.queryDuration, "420", "errorQuery"))
	require.Equal(t, 0, counterFromGaugeByLabels(lp.datasetSize, "420", "errorQuery"))
}

func createObservedPollLogger(t *testing.T, chainId int64) *ObservedLogPoller {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.ErrorLevel)
	db := pgtest.NewSqlxDB(t)
	orm := NewORM(big.NewInt(chainId), db, lggr, pgtest.NewQConfig(true))
	return NewObservedLogPoller(
		orm, nil, lggr, 1, 1, 1, 1, 1000,
	).(*ObservedLogPoller)
}

func resetMetrics(lp ObservedLogPoller) {
	lp.queryDuration.Reset()
	lp.datasetSize.Reset()
}

func counterFromGaugeByLabels(gaugeVec *prometheus.GaugeVec, labels ...string) int {
	value := testutil.ToFloat64(gaugeVec.WithLabelValues(labels...))
	return int(value)
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
