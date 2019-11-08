package services_test

import (
	"math/big"
	"testing"
	"time"

	"chainlink/core/internal/cltest"
	"chainlink/core/internal/mocks"
	"chainlink/core/services"
	"chainlink/core/store/assets"
	"chainlink/core/store/models"

	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunExecutor_Execute(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runExecutor := services.NewRunExecutor(store)

	j := models.NewJob()
	i := models.Initiator{Type: models.InitiatorWeb}
	j.Initiators = []models.Initiator{i}
	j.Tasks = []models.TaskSpec{
		cltest.NewTask(t, "noop"),
	}
	assert.NoError(t, store.CreateJob(&j))

	run := j.NewRun(i)
	run.Payment = assets.NewLink(9117)
	require.NoError(t, store.CreateJobRun(&run))

	err := runExecutor.Execute(run.ID)
	require.NoError(t, err)

	run, err = store.FindJobRun(run.ID)
	require.NoError(t, err)
	assert.Equal(t, models.RunStatusCompleted, run.Status)
	require.Len(t, run.TaskRuns, 1)
	assert.Equal(t, models.RunStatusCompleted, run.TaskRuns[0].Status)

	actual, err := store.LinkEarnedFor(&j)
	require.NoError(t, err)
	assert.Equal(t, assets.NewLink(9117), actual)
}

func TestRunExecutor_Execute_Pending(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runExecutor := services.NewRunExecutor(store)

	j := models.NewJob()
	i := models.Initiator{Type: models.InitiatorWeb}
	j.Initiators = []models.Initiator{i}
	j.Tasks = []models.TaskSpec{
		cltest.NewTask(t, "noop"),
		cltest.NewTask(t, "nooppend"),
	}
	assert.NoError(t, store.CreateJob(&j))

	run := j.NewRun(i)
	run.Payment = assets.NewLink(9117)
	require.NoError(t, store.CreateJobRun(&run))

	err := runExecutor.Execute(run.ID)
	require.NoError(t, err)

	run, err = store.FindJobRun(run.ID)
	require.NoError(t, err)
	assert.Equal(t, models.RunStatusPendingConfirmations, run.Status)
	require.Len(t, run.TaskRuns, 2)
	assert.Equal(t, models.RunStatusCompleted, run.TaskRuns[0].Status)
	assert.Equal(t, models.RunStatusPendingConfirmations, run.TaskRuns[1].Status)

	actual, err := store.LinkEarnedFor(&j)
	require.NoError(t, err)
	assert.Nil(t, actual)
}

func TestRunExecutor_Execute_RunNotFoundError(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runExecutor := services.NewRunExecutor(store)

	err := runExecutor.Execute(models.NewID())
	require.Error(t, err)
}

func TestRunExecutor_Execute_RunNotRunnableError(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	runExecutor := services.NewRunExecutor(store)

	j := models.NewJob()
	i := models.Initiator{Type: models.InitiatorWeb}
	j.Initiators = []models.Initiator{i}
	j.Tasks = []models.TaskSpec{
		cltest.NewTask(t, "noop"),
	}
	assert.NoError(t, store.CreateJob(&j))

	run := j.NewRun(i)
	run.Status = models.RunStatusPendingConfirmations
	require.NoError(t, store.CreateJobRun(&run))

	err := runExecutor.Execute(run.ID)
	require.Error(t, err)
}

func TestRunExecutor_Execute_CancelActivelyRunningTask(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	clock := cltest.NewTriggerClock(t)
	store.Clock = clock

	runExecutor := services.NewRunExecutor(store)

	j := models.NewJob()
	i := models.Initiator{Type: models.InitiatorWeb}
	j.Initiators = []models.Initiator{i}
	j.Tasks = []models.TaskSpec{
		cltest.NewTask(t, "sleep", `{"until": 2147483647}`),
		cltest.NewTask(t, "noop"),
	}
	assert.NoError(t, store.CreateJob(&j))

	run := j.NewRun(i)
	run.Payment = assets.NewLink(19238)
	require.NoError(t, store.CreateJobRun(&run))

	go func() {
		err := runExecutor.Execute(run.ID)
		require.NoError(t, err)
	}()

	// FIXME: Can't think of a better way to do this
	// Make sure Execute has some time to start the sleep task
	time.Sleep(300 * time.Millisecond)

	runQueue := new(mocks.RunQueue)
	runManager := services.NewRunManager(runQueue, store.Config, store.ORM, store.TxManager, clock)
	runManager.Cancel(run.ID)

	clock.Trigger()

	run, err := store.FindJobRun(run.ID)
	require.NoError(t, err)
	assert.Equal(t, models.RunStatusCancelled, run.Status)

	require.Len(t, run.TaskRuns, 2)
	assert.Equal(t, models.RunStatusCancelled, run.TaskRuns[0].Status)
	assert.Equal(t, models.RunStatusUnstarted, run.TaskRuns[1].Status)

	actual, err := store.LinkEarnedFor(&j)
	require.NoError(t, err)
	assert.Nil(t, actual)
}

func TestRunExecutor_UncleForkDoesNotCompleteJob(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	txManager := new(mocks.TxManager)
	store.TxManager = txManager

	runExecutor := services.NewRunExecutor(store)

	j := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{cltest.NewTask(t, "noop"), cltest.NewTask(t, "nooppend")}
	assert.NoError(t, store.CreateJob(&j))

	run := j.NewRun(j.Initiators[0])
	txHash := cltest.NewHash()
	run.RunRequest.TxHash = &txHash
	run.TaskRuns[0].MinimumConfirmations = null.Uint32From(10)
	run.ObservedHeight = models.NewBig(big.NewInt(0))
	require.NoError(t, store.CreateJobRun(&run))
	txManager.On("GetTxReceipt", txHash).Return(&models.TxReceipt{}, nil)
	require.NoError(t, runExecutor.Execute(run.ID))

	run, err := store.FindJobRun(run.ID)
	require.NoError(t, err)
	assert.Equal(t, models.RunStatusErrored, run.Status)
	require.Len(t, run.TaskRuns, 2)
	assert.Equal(t, models.RunStatusErrored, run.TaskRuns[0].Status)
	assert.Equal(t, models.RunStatusUnstarted, run.TaskRuns[1].Status)

	txManager.AssertExpectations(t)
}
