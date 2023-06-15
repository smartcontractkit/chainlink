package logpoller

import (
	"testing"

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

	_, err := lp.IndexedLogs(common.Hash{}, common.Address{}, 1, []common.Hash{}, 1, pg.WithParentCtx(ctx))
	require.NoError(t, err)
	_, err = lp.IndexedLogsByBlockRange(0, 1, common.Hash{}, common.Address{}, 1, []common.Hash{}, pg.WithParentCtx(ctx))
	require.NoError(t, err)
	_, err = lp.IndexedLogsTopicGreaterThan(common.Hash{}, common.Address{}, 1, common.Hash{}, 1, pg.WithParentCtx(ctx))
	require.NoError(t, err)
	_, err = lp.IndexedLogsTopicRange(common.Hash{}, common.Address{}, 1, common.Hash{}, common.Hash{}, 1, pg.WithParentCtx(ctx))
	require.NoError(t, err)
	_, err = lp.IndexedLogsWithSigsExcluding(common.Address{}, common.Hash{}, common.Hash{}, 1, 0, 1, 1, pg.WithParentCtx(ctx))
	require.NoError(t, err)
	_, err = lp.LogsDataWordRange(common.Hash{}, common.Address{}, 0, common.Hash{}, common.Hash{}, 1, pg.WithParentCtx(ctx))
	require.NoError(t, err)
	_, err = lp.LatestLogEventSigsAddrsWithConfs(0, []common.Hash{{}}, []common.Address{{}}, 1, pg.WithParentCtx(ctx))
	require.NoError(t, err)

	require.Equal(t, 7, testutil.CollectAndCount(lp.histogram))
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
