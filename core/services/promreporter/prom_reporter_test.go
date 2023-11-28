package promreporter_test

import (
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/promreporter"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func newHead() evmtypes.Head {
	return evmtypes.Head{Number: 42, EVMChainID: utils.NewBigI(0)}
}

func newLegacyChainContainer(t *testing.T, db *sqlx.DB) evm.LegacyChainContainer {
	config, dbConfig, evmConfig := txmgr.MakeTestConfigs(t)
	keyStore := cltest.NewKeyStore(t, db, dbConfig).Eth()
	ethClient := evmtest.NewEthClientMockWithDefaultChain(t)
	estimator := gas.NewEstimator(logger.TestLogger(t), ethClient, config, evmConfig.GasEstimator())
	lggr := logger.TestLogger(t)
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr, pgtest.NewQConfig(true)), ethClient, lggr, 100*time.Millisecond, false, 2, 3, 2, 1000)

	txm, err := txmgr.NewTxm(
		db,
		evmConfig,
		evmConfig.GasEstimator(),
		evmConfig.Transactions(),
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
		reporter := promreporter.NewPromReporter(db.DB, newLegacyChainContainer(t, db), logger.TestLogger(t), backend, 10*time.Millisecond)

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

		require.NoError(t, reporter.Start(testutils.Context(t)))
		defer func() { assert.NoError(t, reporter.Close()) }()

		head := newHead()
		reporter.OnNewLongestChain(testutils.Context(t), &head)

		require.Eventually(t, func() bool { return subscribeCalls.Load() >= 1 }, 12*time.Second, 100*time.Millisecond)
	})

	t.Run("with unconfirmed evm.txes", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		cfg := configtest.NewGeneralConfig(t, nil)
		txStore := cltest.NewTestTxStore(t, db, cfg.Database())
		ethKeyStore := cltest.NewKeyStore(t, db, cfg.Database()).Eth()
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
		reporter := promreporter.NewPromReporter(db.DB, newLegacyChainContainer(t, db), logger.TestLogger(t), backend, 10*time.Millisecond)
		require.NoError(t, reporter.Start(testutils.Context(t)))
		defer func() { assert.NoError(t, reporter.Close()) }()

		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, fromAddress)
		cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, fromAddress)
		cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, fromAddress)
		require.NoError(t, utils.JustError(db.Exec(`UPDATE evm.tx_attempts SET broadcast_before_block_num = 7 WHERE eth_tx_id = $1`, etx.ID)))

		head := newHead()
		reporter.OnNewLongestChain(testutils.Context(t), &head)

		require.Eventually(t, func() bool { return subscribeCalls.Load() >= 1 }, 12*time.Second, 100*time.Millisecond)
	})

	t.Run("with unfinished pipeline task runs", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		pgtest.MustExec(t, db, `SET CONSTRAINTS pipeline_task_runs_pipeline_run_id_fkey DEFERRED`)

		backend := mocks.NewPrometheusBackend(t)
		reporter := promreporter.NewPromReporter(db.DB, newLegacyChainContainer(t, db), logger.TestLogger(t), backend, 10*time.Millisecond)

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
		require.NoError(t, reporter.Start(testutils.Context(t)))
		defer func() { assert.NoError(t, reporter.Close()) }()

		head := newHead()
		reporter.OnNewLongestChain(testutils.Context(t), &head)

		require.Eventually(t, func() bool { return subscribeCalls.Load() >= 1 }, 12*time.Second, 100*time.Millisecond)
	})
}
