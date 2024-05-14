package promreporter_test

import (
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/promreporter"
)

func newHead() evmtypes.Head {
	return evmtypes.Head{Number: 42, EVMChainID: ubig.NewI(0)}
}

func newLegacyChainContainer(t *testing.T, db *sqlx.DB) legacyevm.LegacyChainContainer {
	config, dbConfig, evmConfig := txmgr.MakeTestConfigs(t)
	keyStore := cltest.NewKeyStore(t, db).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	estimator := gas.NewEstimator(logger.TestLogger(t), ethClient, config, evmConfig.GasEstimator())
	lggr := logger.TestLogger(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            2,
		BackfillBatchSize:        3,
		RpcBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr), ethClient, lggr, lpOpts)

	txm, err := txmgr.NewTxm(
		db,
		evmConfig,
		evmConfig.GasEstimator(),
		evmConfig.Transactions(),
		nil,
		dbConfig,
		dbConfig.Listener(),
		ethClient,
		lggr,
		lp,
		keyStore,
		estimator)
	require.NoError(t, err)

	cfg := configtest.NewGeneralConfig(t, nil)
	return cltest.NewLegacyChainsWithMockChainAndTxManager(t, ethClient, cfg, txm)
}

func Test_PromReporter_OnNewLongestChain(t *testing.T) {
	t.Run("with nothing in the database", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)

		backend := mocks.NewPrometheusBackend(t)
		reporter := promreporter.NewPromReporter(db, newLegacyChainContainer(t, db), logger.TestLogger(t), backend, 10*time.Millisecond)

		var subscribeCalls atomic.Int32

		backend.On("SetUnconfirmedTransactions", big.NewInt(0), int64(0)).Return()
		backend.On("SetMaxUnconfirmedAge", big.NewInt(0), float64(0)).Return()
		backend.On("SetMaxUnconfirmedBlocks", big.NewInt(0), int64(0)).Return()
		backend.On("SetPipelineTaskRunsQueued", 0).Return()
		backend.On("SetPipelineRunsQueued", 0).
			Run(func(args mock.Arguments) {
				subscribeCalls.Add(1)
			}).
			Return()

		servicetest.Run(t, reporter)

		head := newHead()
		reporter.OnNewLongestChain(testutils.Context(t), &head)

		require.Eventually(t, func() bool { return subscribeCalls.Load() >= 1 }, 12*time.Second, 100*time.Millisecond)
	})

	t.Run("with unconfirmed evm.txes", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		txStore := cltest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

		var subscribeCalls atomic.Int32

		backend := mocks.NewPrometheusBackend(t)
		backend.On("SetUnconfirmedTransactions", big.NewInt(0), int64(3)).Return()
		backend.On("SetMaxUnconfirmedAge", big.NewInt(0), mock.MatchedBy(func(s float64) bool {
			return s > 0
		})).Return()
		backend.On("SetMaxUnconfirmedBlocks", big.NewInt(0), int64(35)).Return()
		backend.On("SetPipelineTaskRunsQueued", 0).Return()
		backend.On("SetPipelineRunsQueued", 0).
			Run(func(args mock.Arguments) {
				subscribeCalls.Add(1)
			}).
			Return()
		reporter := promreporter.NewPromReporter(db, newLegacyChainContainer(t, db), logger.TestLogger(t), backend, 10*time.Millisecond)
		servicetest.Run(t, reporter)

		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, fromAddress)
		cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, fromAddress)
		cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, fromAddress)
		require.NoError(t, txStore.UpdateTxAttemptBroadcastBeforeBlockNum(testutils.Context(t), etx.ID, 7))

		head := newHead()
		reporter.OnNewLongestChain(testutils.Context(t), &head)

		require.Eventually(t, func() bool { return subscribeCalls.Load() >= 1 }, 12*time.Second, 100*time.Millisecond)
	})

	t.Run("with unfinished pipeline task runs", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		pgtest.MustExec(t, db, `SET CONSTRAINTS pipeline_task_runs_pipeline_run_id_fkey DEFERRED`)

		backend := mocks.NewPrometheusBackend(t)
		reporter := promreporter.NewPromReporter(db, newLegacyChainContainer(t, db), logger.TestLogger(t), backend, 10*time.Millisecond)

		cltest.MustInsertUnfinishedPipelineTaskRun(t, db, 1)
		cltest.MustInsertUnfinishedPipelineTaskRun(t, db, 1)
		cltest.MustInsertUnfinishedPipelineTaskRun(t, db, 2)

		var subscribeCalls atomic.Int32

		backend.On("SetUnconfirmedTransactions", big.NewInt(0), int64(0)).Return()
		backend.On("SetMaxUnconfirmedAge", big.NewInt(0), float64(0)).Return()
		backend.On("SetMaxUnconfirmedBlocks", big.NewInt(0), int64(0)).Return()
		backend.On("SetPipelineTaskRunsQueued", 3).Return()
		backend.On("SetPipelineRunsQueued", 2).
			Run(func(args mock.Arguments) {
				subscribeCalls.Add(1)
			}).
			Return()
		servicetest.Run(t, reporter)

		head := newHead()
		reporter.OnNewLongestChain(testutils.Context(t), &head)

		require.Eventually(t, func() bool { return subscribeCalls.Load() >= 1 }, 12*time.Second, 100*time.Millisecond)
	})
}
