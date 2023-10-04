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

	_, _ = lp.SelectIndexedLogs(common.Address{}, common.Hash{}, 1, []common.Hash{}, 1, pg.WithParentCtx(ctx))
	_, _ = lp.SelectIndexedLogsByBlockRange(0, 1, common.Address{}, common.Hash{}, 1, []common.Hash{}, pg.WithParentCtx(ctx))
	_, _ = lp.SelectIndexedLogsTopicGreaterThan(common.Address{}, common.Hash{}, 1, common.Hash{}, 1, pg.WithParentCtx(ctx))
	_, _ = lp.SelectIndexedLogsTopicRange(common.Address{}, common.Hash{}, 1, common.Hash{}, common.Hash{}, 1, pg.WithParentCtx(ctx))
	_, _ = lp.SelectIndexedLogsWithSigsExcluding(common.Hash{}, common.Hash{}, 1, common.Address{}, 0, 1, 1, pg.WithParentCtx(ctx))
	_, _ = lp.SelectLogsDataWordRange(common.Address{}, common.Hash{}, 0, common.Hash{}, common.Hash{}, 1, pg.WithParentCtx(ctx))
	_, _ = lp.SelectLogsDataWordGreaterThan(common.Address{}, common.Hash{}, 0, common.Hash{}, 1, pg.WithParentCtx(ctx))
	_, _ = lp.SelectLogsCreatedAfter(common.Address{}, common.Hash{}, time.Now(), 0, pg.WithParentCtx(ctx))
	_, _ = lp.SelectLatestLogByEventSigWithConfs(common.Hash{}, common.Address{}, 0, pg.WithParentCtx(ctx))
	_, _ = lp.SelectLatestLogEventSigsAddrsWithConfs(0, []common.Address{{}}, []common.Hash{{}}, 1, pg.WithParentCtx(ctx))
	_, _ = lp.SelectIndexedLogsCreatedAfter(common.Address{}, common.Hash{}, 1, []common.Hash{}, time.Now(), 0, pg.WithParentCtx(ctx))
	_, _ = lp.SelectLogsUntilBlockHashDataWordGreaterThan(common.Address{}, common.Hash{}, 0, common.Hash{}, common.Hash{}, pg.WithParentCtx(ctx))
	_ = lp.InsertLogs([]Log{}, pg.WithParentCtx(ctx))
	_ = lp.InsertBlock(common.Hash{}, 0, time.Now(), pg.WithParentCtx(ctx))

	require.Equal(t, 14, testutil.CollectAndCount(lp.queryDuration))
	require.Equal(t, 10, testutil.CollectAndCount(lp.datasetSize))
	resetMetrics(*lp)
}

func TestShouldPublishDurationInCaseOfError(t *testing.T) {
	ctx := testutils.Context(t)
	lp := createObservedPollLogger(t, 200)
	require.Equal(t, 0, testutil.CollectAndCount(lp.queryDuration))

	_, err := lp.SelectLatestLogByEventSigWithConfs(common.Hash{}, common.Address{}, 0, pg.WithParentCtx(ctx))
	require.Error(t, err)

	require.Equal(t, 1, testutil.CollectAndCount(lp.queryDuration))
	require.Equal(t, 1, counterFromHistogramByLabels(t, lp.queryDuration, "200", "SelectLatestLogByEventSigWithConfs"))

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

func TestMetricsAreProperlyPopulatedForWrites(t *testing.T) {
	lp := createObservedPollLogger(t, 420)
	require.NoError(t, withObservedExec(lp, "execQuery", func() error { return nil }))
	require.Error(t, withObservedExec(lp, "execQuery", func() error { return fmt.Errorf("error") }))

	require.Equal(t, 2, counterFromHistogramByLabels(t, lp.queryDuration, "420", "execQuery"))
}

func createObservedPollLogger(t *testing.T, chainId int64) *ObservedORM {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.ErrorLevel)
	db := pgtest.NewSqlxDB(t)
	return NewObservedORM(
		big.NewInt(chainId), db, lggr, pgtest.NewQConfig(true),
	)
}

func resetMetrics(lp ObservedORM) {
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
