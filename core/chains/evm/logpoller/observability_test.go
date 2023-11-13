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
	orm := createObservedORM(t, 100)
	t.Cleanup(func() { resetMetrics(*orm) })
	require.Equal(t, 0, testutil.CollectAndCount(orm.queryDuration))

	_, _ = orm.SelectIndexedLogs(common.Address{}, common.Hash{}, 1, []common.Hash{}, 1, pg.WithParentCtx(ctx))
	_, _ = orm.SelectIndexedLogsByBlockRange(0, 1, common.Address{}, common.Hash{}, 1, []common.Hash{}, pg.WithParentCtx(ctx))
	_, _ = orm.SelectIndexedLogsTopicGreaterThan(common.Address{}, common.Hash{}, 1, common.Hash{}, 1, pg.WithParentCtx(ctx))
	_, _ = orm.SelectIndexedLogsTopicRange(common.Address{}, common.Hash{}, 1, common.Hash{}, common.Hash{}, 1, pg.WithParentCtx(ctx))
	_, _ = orm.SelectIndexedLogsWithSigsExcluding(common.Hash{}, common.Hash{}, 1, common.Address{}, 0, 1, 1, pg.WithParentCtx(ctx))
	_, _ = orm.SelectLogsDataWordRange(common.Address{}, common.Hash{}, 0, common.Hash{}, common.Hash{}, 1, pg.WithParentCtx(ctx))
	_, _ = orm.SelectLogsDataWordGreaterThan(common.Address{}, common.Hash{}, 0, common.Hash{}, 1, pg.WithParentCtx(ctx))
	_, _ = orm.SelectLogsCreatedAfter(common.Address{}, common.Hash{}, time.Now(), 0, pg.WithParentCtx(ctx))
	_, _ = orm.SelectLatestLogByEventSigWithConfs(common.Hash{}, common.Address{}, 0, pg.WithParentCtx(ctx))
	_, _ = orm.SelectLatestLogEventSigsAddrsWithConfs(0, []common.Address{{}}, []common.Hash{{}}, 1, pg.WithParentCtx(ctx))
	_, _ = orm.SelectIndexedLogsCreatedAfter(common.Address{}, common.Hash{}, 1, []common.Hash{}, time.Now(), 0, pg.WithParentCtx(ctx))
	_ = orm.InsertLogs([]Log{}, pg.WithParentCtx(ctx))
	_ = orm.InsertBlock(common.Hash{}, 1, time.Now(), 0, pg.WithParentCtx(ctx))

	require.Equal(t, 13, testutil.CollectAndCount(orm.queryDuration))
	require.Equal(t, 10, testutil.CollectAndCount(orm.datasetSize))
}

func TestShouldPublishDurationInCaseOfError(t *testing.T) {
	ctx := testutils.Context(t)
	orm := createObservedORM(t, 200)
	t.Cleanup(func() { resetMetrics(*orm) })
	require.Equal(t, 0, testutil.CollectAndCount(orm.queryDuration))

	_, err := orm.SelectLatestLogByEventSigWithConfs(common.Hash{}, common.Address{}, 0, pg.WithParentCtx(ctx))
	require.Error(t, err)

	require.Equal(t, 1, testutil.CollectAndCount(orm.queryDuration))
	require.Equal(t, 1, counterFromHistogramByLabels(t, orm.queryDuration, "200", "SelectLatestLogByEventSigWithConfs"))
}

func TestMetricsAreProperlyPopulatedWithLabels(t *testing.T) {
	orm := createObservedORM(t, 420)
	t.Cleanup(func() { resetMetrics(*orm) })
	expectedCount := 9
	expectedSize := 2

	for i := 0; i < expectedCount; i++ {
		_, err := withObservedQueryAndResults(orm, "query", func() ([]string, error) { return []string{"value1", "value2"}, nil })
		require.NoError(t, err)
	}

	require.Equal(t, expectedCount, counterFromHistogramByLabels(t, orm.queryDuration, "420", "query"))
	require.Equal(t, expectedSize, counterFromGaugeByLabels(orm.datasetSize, "420", "query"))

	require.Equal(t, 0, counterFromHistogramByLabels(t, orm.queryDuration, "420", "other_query"))
	require.Equal(t, 0, counterFromHistogramByLabels(t, orm.queryDuration, "5", "query"))

	require.Equal(t, 0, counterFromGaugeByLabels(orm.datasetSize, "420", "other_query"))
	require.Equal(t, 0, counterFromGaugeByLabels(orm.datasetSize, "5", "query"))
}

func TestNotPublishingDatasetSizeInCaseOfError(t *testing.T) {
	orm := createObservedORM(t, 420)

	_, err := withObservedQueryAndResults(orm, "errorQuery", func() ([]string, error) { return nil, fmt.Errorf("error") })
	require.Error(t, err)

	require.Equal(t, 1, counterFromHistogramByLabels(t, orm.queryDuration, "420", "errorQuery"))
	require.Equal(t, 0, counterFromGaugeByLabels(orm.datasetSize, "420", "errorQuery"))
}

func TestMetricsAreProperlyPopulatedForWrites(t *testing.T) {
	orm := createObservedORM(t, 420)
	require.NoError(t, withObservedExec(orm, "execQuery", func() error { return nil }))
	require.Error(t, withObservedExec(orm, "execQuery", func() error { return fmt.Errorf("error") }))

	require.Equal(t, 2, counterFromHistogramByLabels(t, orm.queryDuration, "420", "execQuery"))
}

func createObservedORM(t *testing.T, chainId int64) *ObservedORM {
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
