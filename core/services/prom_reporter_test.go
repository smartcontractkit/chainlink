package services_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"
)

func newHead() eth.Head {
	return eth.Head{Number: 42, EVMChainID: utils.NewBigI(0)}
}

func Test_PromReporter_OnNewLongestChain(t *testing.T) {
	t.Run("with nothing in the database", func(t *testing.T) {
		d := pgtest.NewSqlDB(t)

		backend := new(mocks.PrometheusBackend)
		backend.Test(t)
		reporter := services.NewPromReporter(d, logger.TestLogger(t), backend, 10*time.Millisecond)

		var subscribeCalls atomic.Int32

		backend.On("SetUnconfirmedTransactions", big.NewInt(0), int64(0)).Return()
		backend.On("SetMaxUnconfirmedAge", big.NewInt(0), float64(0)).Return()
		backend.On("SetMaxUnconfirmedBlocks", big.NewInt(0), int64(0)).Return()
		backend.On("SetPipelineTaskRunsQueued", 0).Return()
		backend.On("SetPipelineRunsQueued", 0).
			Run(func(args mock.Arguments) {
				subscribeCalls.Inc()
			}).
			Return()

		reporter.Start()
		defer reporter.Close()

		head := newHead()
		reporter.OnNewLongestChain(context.Background(), head)

		require.Eventually(t, func() bool { return subscribeCalls.Load() >= 1 }, 12*time.Second, 100*time.Millisecond)

		backend.AssertExpectations(t)
	})

	t.Run("with unconfirmed eth_txes", func(t *testing.T) {
		db := pgtest.NewGormDB(t)
		sqlxdb := postgres.UnwrapGormDB(db)
		ethKeyStore := cltest.NewKeyStore(t, sqlxdb).Eth()
		_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

		var subscribeCalls atomic.Int32

		backend := new(mocks.PrometheusBackend)
		backend.Test(t)
		backend.On("SetUnconfirmedTransactions", big.NewInt(0), int64(3)).Return()
		backend.On("SetMaxUnconfirmedAge", big.NewInt(0), mock.MatchedBy(func(s float64) bool {
			return s > 0
		})).Return()
		backend.On("SetMaxUnconfirmedBlocks", big.NewInt(0), int64(35)).Return()
		backend.On("SetPipelineTaskRunsQueued", 0).Return()
		backend.On("SetPipelineRunsQueued", 0).
			Run(func(args mock.Arguments) {
				subscribeCalls.Inc()
			}).
			Return()
		d, _ := db.DB()
		reporter := services.NewPromReporter(d, logger.TestLogger(t), backend, 10*time.Millisecond)
		reporter.Start()
		defer reporter.Close()

		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, db, 0, fromAddress)
		cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, db, 1, fromAddress)
		cltest.MustInsertUnconfirmedEthTxWithBroadcastLegacyAttempt(t, db, 2, fromAddress)
		require.NoError(t, db.Exec(`UPDATE eth_tx_attempts SET broadcast_before_block_num = 7 WHERE eth_tx_id = ?`, etx.ID).Error)

		head := newHead()
		reporter.OnNewLongestChain(context.Background(), head)

		require.Eventually(t, func() bool { return subscribeCalls.Load() >= 1 }, 12*time.Second, 100*time.Millisecond)

		backend.AssertExpectations(t)
	})

	t.Run("with unfinished pipeline task runs", func(t *testing.T) {
		db := pgtest.NewGormDB(t)
		d, _ := db.DB()

		_, err := d.Exec(`SET CONSTRAINTS pipeline_task_runs_pipeline_run_id_fkey DEFERRED`)
		require.NoError(t, err)

		backend := new(mocks.PrometheusBackend)
		backend.Test(t)
		reporter := services.NewPromReporter(d, logger.TestLogger(t), backend, 10*time.Millisecond)

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
				subscribeCalls.Inc()
			}).
			Return()
		reporter.Start()
		defer reporter.Close()

		head := newHead()
		reporter.OnNewLongestChain(context.Background(), head)

		require.Eventually(t, func() bool { return subscribeCalls.Load() >= 1 }, 12*time.Second, 100*time.Millisecond)

		backend.AssertExpectations(t)
	})
}
