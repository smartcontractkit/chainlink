package headreporter_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/headreporter"
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

func Test_HeadReporterService(t *testing.T) {
	t.Run("report everything", func(t *testing.T) {
		db := pgtest.NewSqlxDB(t)

		headReporter := mocks.NewHeadReporter(t)
		service := headreporter.NewHeadReporterServiceWithReporters(db, newLegacyChainContainer(t, db), logger.TestLogger(t), []headreporter.HeadReporter{headReporter}, time.Second)
		err := service.Start(testutils.Context(t))
		require.NoError(t, err)

		var reportCalls atomic.Int32
		head := newHead()
		headReporter.On("ReportNewHead", mock.Anything, &head).Run(func(args mock.Arguments) {
			reportCalls.Add(1)
		}).Return(nil)
		headReporter.On("ReportPeriodic", mock.Anything).Run(func(args mock.Arguments) {
			reportCalls.Add(1)
		}).Return(nil)
		service.OnNewLongestChain(testutils.Context(t), &head)

		require.Eventually(t, func() bool { return reportCalls.Load() == 2 }, 5*time.Second, 100*time.Millisecond)
	})
}
