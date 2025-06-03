package headreporter_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys/keystest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr/txmgrtest"

	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/headreporter"
)

func Test_PrometheusReporter(t *testing.T) {
	t.Run("with nothing in the database", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)

		backend := headreporter.NewMockPrometheusBackend(t)
		backend.On("SetUnconfirmedTransactions", big.NewInt(0), int64(0)).Return()
		backend.On("SetMaxUnconfirmedAge", big.NewInt(0), float64(0)).Return()
		backend.On("SetMaxUnconfirmedBlocks", big.NewInt(0), int64(0)).Return()

		reporter := headreporter.NewLegacyEVMPrometheusReporter(db, newLegacyChainContainer(t, db))
		reporter.SetBackend(backend)

		head := headreporter.NewHead()
		err := reporter.ReportNewHead(testutils.Context(t), &head)
		require.NoError(t, err)

		backend.On("SetPipelineTaskRunsQueued", 0).Return()
		backend.On("SetPipelineRunsQueued", 0).Return()
		err = reporter.ReportPeriodic(testutils.Context(t))
		require.NoError(t, err)
	})

	t.Run("with null txm", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		backend := headreporter.NewMockPrometheusBackend(t)

		reporter := headreporter.NewLegacyEVMPrometheusReporter(db, newLegacyChainContainerWithNullTxm(t))
		reporter.SetBackend(backend)

		head := headreporter.NewHead()
		err := reporter.ReportNewHead(testutils.Context(t), &head)
		require.NoError(t, err)
	})

	t.Run("with unconfirmed evm.txes", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		txStore := txmgrtest.NewTestTxStore(t, db)
		ethKeyStore := cltest.NewKeyStore(t, db).Eth()
		_, fromAddress := cltest.MustInsertRandomKey(t, ethKeyStore)

		etx := txmgrtest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 0, fromAddress)
		txmgrtest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 1, fromAddress)
		txmgrtest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, txStore, 2, fromAddress)
		require.NoError(t, txStore.UpdateTxAttemptBroadcastBeforeBlockNum(testutils.Context(t), etx.ID, 7))

		backend := headreporter.NewMockPrometheusBackend(t)
		backend.On("SetUnconfirmedTransactions", big.NewInt(0), int64(3)).Return()
		backend.On("SetMaxUnconfirmedAge", big.NewInt(0), mock.MatchedBy(func(s float64) bool {
			return s > 0
		})).Return()
		backend.On("SetMaxUnconfirmedBlocks", big.NewInt(0), int64(35)).Return()

		reporter := headreporter.NewLegacyEVMPrometheusReporter(db, newLegacyChainContainer(t, db))
		reporter.SetBackend(backend)

		head := headreporter.NewHead()
		err := reporter.ReportNewHead(testutils.Context(t), &head)
		require.NoError(t, err)

		backend.On("SetPipelineTaskRunsQueued", 0).Return()
		backend.On("SetPipelineRunsQueued", 0).Return()

		err = reporter.ReportPeriodic(testutils.Context(t))
		require.NoError(t, err)
	})

	t.Run("with unfinished pipeline task runs", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)
		pgtest.MustExec(t, db, `SET CONSTRAINTS pipeline_task_runs_pipeline_run_id_fkey DEFERRED`)

		evmtestutils.MustInsertUnfinishedPipelineTaskRun(t, db, 1)
		evmtestutils.MustInsertUnfinishedPipelineTaskRun(t, db, 1)
		evmtestutils.MustInsertUnfinishedPipelineTaskRun(t, db, 2)

		backend := headreporter.NewMockPrometheusBackend(t)
		backend.On("SetUnconfirmedTransactions", big.NewInt(0), int64(0)).Return()
		backend.On("SetMaxUnconfirmedAge", big.NewInt(0), float64(0)).Return()
		backend.On("SetMaxUnconfirmedBlocks", big.NewInt(0), int64(0)).Return()

		reporter := headreporter.NewLegacyEVMPrometheusReporter(db, newLegacyChainContainer(t, db))
		reporter.SetBackend(backend)

		head := headreporter.NewHead()
		err := reporter.ReportNewHead(testutils.Context(t), &head)
		require.NoError(t, err)

		backend.On("SetPipelineTaskRunsQueued", 3).Return()
		backend.On("SetPipelineRunsQueued", 2).Return()

		err = reporter.ReportPeriodic(testutils.Context(t))
		require.NoError(t, err)
	})
}

func newLegacyChainContainer(t *testing.T, db *sqlx.DB) legacyevm.LegacyChainContainer {
	config, dbConfig, evmConfig := txmgr.MakeTestConfigs(t)
	keyStore := &keystest.FakeChainStore{}
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	estimator, err := gas.NewEstimator(logger.TestLogger(t), ethClient, config.ChainType(), ethClient.ConfiguredChainID(), evmConfig.GasEstimator(), nil)
	require.NoError(t, err)
	lggr := logger.TestLogger(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            2,
		BackfillBatchSize:        3,
		RPCBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	ht := headstest.NewSimulatedHeadTracker(ethClient, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	lp := logpoller.NewLogPoller(logpoller.NewORM(testutils.FixtureChainID, db, lggr), ethClient, lggr, ht, lpOpts)

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
		estimator,
		ht,
		nil)
	require.NoError(t, err)

	cfg := configtest.NewGeneralConfig(t, nil)
	return cltest.NewLegacyChainsWithMockChainAndTxManager(t, ethClient, cfg, txm)
}

func newLegacyChainContainerWithNullTxm(t *testing.T) legacyevm.LegacyChainContainer {
	ethClient := clienttest.NewClientWithDefaultChainID(t)
	txm := &txmgr.NullTxManager{ErrMsg: fmt.Sprintf("TXM disabled for chain %d", ethClient.ConfiguredChainID())}
	cfg := configtest.NewGeneralConfig(t, nil)
	return cltest.NewLegacyChainsWithMockChainAndTxManager(t, ethClient, cfg, txm)
}
