package services_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_PromReporter_OnNewLongestChain(t *testing.T) {
	t.Run("with nothing in the database", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()

		backend := new(mocks.PrometheusBackend)
		d, _ := store.DB.DB()
		reporter := services.NewPromReporter(d, backend)
		reporter.Start()
		defer reporter.Close()

		var subscribeCalls int32

		subscribeCallsPrt := &subscribeCalls

		backend.On("SetUnconfirmedTransactions", int64(0)).Return()
		backend.On("SetMaxUnconfirmedAge", float64(0)).Return()
		backend.On("SetMaxUnconfirmedBlocks", int64(0)).Return()
		backend.On("SetPipelineTaskRunsQueued", 0).Return()
		backend.On("SetPipelineRunsQueued", 0).
			Run(func(args mock.Arguments) {
				atomic.AddInt32(&subscribeCalls, 1)
			}).
			Return()

		head := models.Head{Number: 42}
		reporter.OnNewLongestChain(context.Background(), head)

		require.Eventually(t, func() bool { return atomic.LoadInt32(subscribeCallsPrt) == 1 }, 12*time.Second, 100*time.Millisecond)

		backend.AssertExpectations(t)
	})

	t.Run("with unconfirmed eth_txes", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()
		ethKeyStore := cltest.NewKeyStore(t, store.DB).Eth
		_, fromAddress := cltest.MustAddRandomKeyToKeystore(t, ethKeyStore)

		backend := new(mocks.PrometheusBackend)
		d, _ := store.DB.DB()
		reporter := services.NewPromReporter(d, backend)
		reporter.Start()
		defer reporter.Close()

		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 0, fromAddress)
		cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 1, fromAddress)
		cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 2, fromAddress)
		require.NoError(t, store.DB.Exec(`UPDATE eth_tx_attempts SET broadcast_before_block_num = 7 WHERE eth_tx_id = ?`, etx.ID).Error)

		var subscribeCalls int32
		subscribeCallsPrt := &subscribeCalls

		backend.On("SetUnconfirmedTransactions", int64(3)).Return()
		backend.On("SetMaxUnconfirmedAge", mock.MatchedBy(func(s float64) bool {
			return s > 0
		})).Return()
		backend.On("SetMaxUnconfirmedBlocks", int64(35)).Return()
		backend.On("SetPipelineTaskRunsQueued", 0).Return()
		backend.On("SetPipelineRunsQueued", 0).
			Run(func(args mock.Arguments) {
				atomic.AddInt32(&subscribeCalls, 1)
			}).
			Return()

		head := models.Head{Number: 42}
		reporter.OnNewLongestChain(context.Background(), head)

		require.Eventually(t, func() bool { return atomic.LoadInt32(subscribeCallsPrt) == 1 }, 12*time.Second, 100*time.Millisecond)

		backend.AssertExpectations(t)
	})

	t.Run("with unfinished pipeline task runs", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()

		require.NoError(t, store.DB.Exec(`SET CONSTRAINTS pipeline_task_runs_pipeline_run_id_fkey DEFERRED`).Error)

		backend := new(mocks.PrometheusBackend)
		d, _ := store.DB.DB()
		reporter := services.NewPromReporter(d, backend)
		reporter.Start()
		defer reporter.Close()

		cltest.MustInsertUnfinishedPipelineTaskRun(t, store, 1)
		cltest.MustInsertUnfinishedPipelineTaskRun(t, store, 1)
		cltest.MustInsertUnfinishedPipelineTaskRun(t, store, 2)

		var subscribeCalls int32
		subscribeCallsPrt := &subscribeCalls

		backend.On("SetUnconfirmedTransactions", int64(0)).Return()
		backend.On("SetMaxUnconfirmedAge", float64(0)).Return()
		backend.On("SetMaxUnconfirmedBlocks", int64(0)).Return()
		backend.On("SetPipelineTaskRunsQueued", 3).Return()
		backend.On("SetPipelineRunsQueued", 2).
			Run(func(args mock.Arguments) {
				atomic.AddInt32(&subscribeCalls, 1)
			}).
			Return()

		head := models.Head{Number: 42}
		reporter.OnNewLongestChain(context.Background(), head)

		require.Eventually(t, func() bool { return atomic.LoadInt32(subscribeCallsPrt) == 1 }, 12*time.Second, 100*time.Millisecond)

		backend.AssertExpectations(t)
	})
}
