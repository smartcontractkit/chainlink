package services_test

import (
	"context"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/require"
)

func Test_PromReporter_OnNewLongestChain(t *testing.T) {
	t.Run("with nothing in the database", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()

		backend := new(mocks.PrometheusBackend)
		reporter := services.NewPromReporter(store.DB.DB(), backend)

		backend.On("SetUnconfirmedTransactions", int64(0)).Return()
		backend.On("SetMaxUnconfirmedBlocks", int64(0)).Return()
		backend.On("SetPipelineTaskRunsQueued", 0).Return()
		backend.On("SetPipelineRunsQueued", 0).Return()

		head := models.Head{Number: 42}
		reporter.OnNewLongestChain(context.Background(), head)

		backend.AssertExpectations(t)
	})

	t.Run("with unconfirmed eth_txes", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()

		backend := new(mocks.PrometheusBackend)
		reporter := services.NewPromReporter(store.DB.DB(), backend)

		etx := cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 0)
		cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 1)
		cltest.MustInsertUnconfirmedEthTxWithBroadcastAttempt(t, store, 2)
		require.NoError(t, store.DB.Exec(`UPDATE eth_tx_attempts SET broadcast_before_block_num = 7 WHERE eth_tx_id = ?`, etx.ID).Error)

		backend.On("SetUnconfirmedTransactions", int64(3)).Return()
		backend.On("SetMaxUnconfirmedBlocks", int64(35)).Return()
		backend.On("SetPipelineTaskRunsQueued", 0).Return()
		backend.On("SetPipelineRunsQueued", 0).Return()

		head := models.Head{Number: 42}
		reporter.OnNewLongestChain(context.Background(), head)

		backend.AssertExpectations(t)
	})

	t.Run("with unfinished pipeline task runs", func(t *testing.T) {
		store, cleanup := cltest.NewStore(t)
		defer cleanup()

		require.NoError(t, store.DB.Exec(`SET CONSTRAINTS pipeline_task_runs_pipeline_run_id_fkey DEFERRED`).Error)
		require.NoError(t, store.DB.Exec(`SET CONSTRAINTS pipeline_task_runs_pipeline_task_spec_id_fkey DEFERRED`).Error)

		backend := new(mocks.PrometheusBackend)
		reporter := services.NewPromReporter(store.DB.DB(), backend)

		cltest.MustInsertUnfinishedPipelineTaskRun(t, store, 1)
		cltest.MustInsertUnfinishedPipelineTaskRun(t, store, 1)
		cltest.MustInsertUnfinishedPipelineTaskRun(t, store, 2)

		backend.On("SetUnconfirmedTransactions", int64(0)).Return()
		backend.On("SetMaxUnconfirmedBlocks", int64(0)).Return()
		backend.On("SetPipelineTaskRunsQueued", 3).Return()
		backend.On("SetPipelineRunsQueued", 2).Return()

		head := models.Head{Number: 42}
		reporter.OnNewLongestChain(context.Background(), head)

		backend.AssertExpectations(t)
	})
}
