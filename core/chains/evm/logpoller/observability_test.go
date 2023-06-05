package logpoller

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
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
	lp := createObservedPollLogger(t)
	require.Equal(t, 0, testutil.CollectAndCount(lp.histogram))

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

	require.Equal(t, 11, testutil.CollectAndCount(lp.histogram))
	resetMetrics(*lp)
}

func TestShouldPublishMetricInCaseOfError(t *testing.T) {
	ctx := testutils.Context(t)
	lp := createObservedPollLogger(t)
	require.Equal(t, 0, testutil.CollectAndCount(lp.histogram))

	_, err := lp.LatestLogByEventSigWithConfs(common.Hash{}, common.Address{}, 0, pg.WithParentCtx(ctx))
	require.Error(t, err)

	require.Equal(t, 1, testutil.CollectAndCount(lp.histogram))
	resetMetrics(*lp)
}

func TestNotObservedFunctions(t *testing.T) {
	ctx := testutils.Context(t)
	lp := createObservedPollLogger(t)
	require.Equal(t, 0, testutil.CollectAndCount(lp.histogram))

	_, err := lp.Logs(0, 1, common.Hash{}, common.Address{}, pg.WithParentCtx(ctx))
	require.NoError(t, err)

	_, err = lp.LogsWithSigs(0, 1, []common.Hash{{}}, common.Address{}, pg.WithParentCtx(ctx))
	require.NoError(t, err)

	require.Equal(t, 0, testutil.CollectAndCount(lp.histogram))
	resetMetrics(*lp)
}

func createObservedPollLogger(t *testing.T) *ObservedLogPoller {
	lggr, _ := logger.TestLoggerObserved(t, zapcore.ErrorLevel)
	db := pgtest.NewSqlxDB(t)
	orm := NewORM(testutils.NewRandomEVMChainID(), db, lggr, pgtest.NewQConfig(true))
	return NewObservedLogPoller(
		orm, nil, lggr, 1, 1, 1, 1, 1000,
	).(*ObservedLogPoller)
}

func resetMetrics(lp ObservedLogPoller) {
	lp.histogram.Reset()
}
