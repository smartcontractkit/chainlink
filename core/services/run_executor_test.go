package services_test

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/mocks"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunExecutor_Execute(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	keyStore := cltest.NewKeyStore(t, store.DB)

	pusher := new(mocks.StatsPusher)
	pusher.On("PushNow").Return(nil)

	runExecutor := services.NewRunExecutor(store, keyStore, pusher)

	j := models.NewJob()
	i := models.Initiator{Type: models.InitiatorWeb}
	j.Initiators = []models.Initiator{i}
	j.Tasks = []models.TaskSpec{
		cltest.NewTask(t, "noop"),
	}
	assert.NoError(t, store.CreateJob(&j))

	run := cltest.NewJobRun(j)
	run.Payment = assets.NewLink(9117)
	require.NoError(t, store.CreateJobRun(&run))

	err := runExecutor.Execute(run.ID)
	require.NoError(t, err)

	run, err = store.FindJobRun(run.ID)
	require.NoError(t, err)
	assert.Equal(t, models.RunStatusCompleted, run.GetStatus())
	require.Len(t, run.TaskRuns, 1)
	assert.Equal(t, models.RunStatusCompleted, run.TaskRuns[0].Status)

	actual, err := store.LinkEarnedFor(&j)
	require.NoError(t, err)
	assert.Equal(t, assets.NewLink(9117), actual)
}

func TestRunExecutor_Execute_PendingOutgoing(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	keyStore := cltest.NewKeyStore(t, store.DB)

	pusher := new(mocks.StatsPusher)
	pusher.On("PushNow").Return(nil)

	runExecutor := services.NewRunExecutor(store, keyStore, pusher)

	j := models.NewJob()
	i := models.Initiator{Type: models.InitiatorWeb}
	j.Initiators = []models.Initiator{i}
	j.Tasks = []models.TaskSpec{
		cltest.NewTask(t, "noop"),
		cltest.NewTask(t, "nooppendoutgoing"),
	}
	assert.NoError(t, store.CreateJob(&j))

	run := cltest.NewJobRun(j)
	require.NoError(t, store.CreateJobRun(&run))

	err := runExecutor.Execute(run.ID)
	require.NoError(t, err)

	run, err = store.FindJobRun(run.ID)
	require.NoError(t, err)
	assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, run.GetStatus())
	require.Len(t, run.TaskRuns, 2)
	assert.Equal(t, models.RunStatusCompleted, run.TaskRuns[0].Status)
	assert.Equal(t, models.RunStatusPendingOutgoingConfirmations, run.TaskRuns[1].Status)

	actual, err := store.LinkEarnedFor(&j)
	require.NoError(t, err)
	assert.Nil(t, actual)
}

func TestRunExecutor_Execute_RunNotFoundError(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	keyStore := cltest.NewKeyStore(t, store.DB)

	pusher := new(mocks.StatsPusher)
	pusher.On("PushNow").Return(nil)

	runExecutor := services.NewRunExecutor(store, keyStore, pusher)

	err := runExecutor.Execute(uuid.NewV4())
	require.Error(t, err)
}

func TestRunExecutor_Execute_CancelActivelyRunningTask(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	keyStore := cltest.NewKeyStore(t, store.DB)
	// It will never be triggered, so sleep tasks will run forever (or until
	// cancelled)
	clock := cltest.NewTriggerClock(t)
	store.Clock = clock

	pusher := new(mocks.StatsPusher)
	pusher.On("PushNow").Return(nil)

	runExecutor := services.NewRunExecutor(store, keyStore, pusher)

	j := models.NewJob()
	i := models.Initiator{Type: models.InitiatorWeb}
	j.Initiators = []models.Initiator{i}
	j.Tasks = []models.TaskSpec{
		cltest.NewTask(t, "sleep", `{"until": 2147483647}`),
		cltest.NewTask(t, "noop"),
	}
	assert.NoError(t, store.CreateJob(&j))

	run := cltest.NewJobRun(j)
	require.NoError(t, store.CreateJobRun(&run))

	go func() {
		err := runExecutor.Execute(run.ID)
		require.NoError(t, err)
	}()

	// FIXME: Can't think of a better way to do this
	// Make sure Execute has some time to start the sleep task
	time.Sleep(1 * time.Second)

	runQueue := new(mocks.RunQueue)
	runManager := services.NewRunManager(runQueue, store.Config, store.ORM, pusher, store.Clock)
	_, err := runManager.Cancel(run.ID)
	require.NoError(t, err)

	run, err = store.FindJobRun(run.ID)
	require.NoError(t, err)
	assert.Equal(t, models.RunStatusCancelled, run.GetStatus())

	require.Len(t, run.TaskRuns, 2)
	assert.Equal(t, models.RunStatusCancelled, run.TaskRuns[0].Status)
	assert.Equal(t, models.RunStatusUnstarted, run.TaskRuns[1].Status)

	actual, err := store.LinkEarnedFor(&j)
	require.NoError(t, err)
	assert.Nil(t, actual)
}

func TestRunExecutor_InitialTaskLacksConfirmations(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	keyStore := cltest.NewKeyStore(t, store.DB)

	pusher := new(mocks.StatsPusher)
	pusher.On("PushNow").Return(nil)

	runExecutor := services.NewRunExecutor(store, keyStore, pusher)

	j := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{cltest.NewTask(t, "noop")}
	assert.NoError(t, store.CreateJob(&j))

	run := cltest.NewJobRun(j)
	txHash := cltest.NewHash()
	run.RunRequest.TxHash = &txHash
	run.TaskRuns[0].MinRequiredIncomingConfirmations = null.Uint32From(10)
	run.CreationHeight = utils.NewBig(big.NewInt(0))
	run.ObservedHeight = run.CreationHeight
	require.NoError(t, store.CreateJobRun(&run))
	require.NoError(t, runExecutor.Execute(run.ID))

	run, err := store.FindJobRun(run.ID)
	require.NoError(t, err)
	assert.Equal(t, models.RunStatusPendingIncomingConfirmations, run.GetStatus())
	require.Len(t, run.TaskRuns, 1)
	assert.Equal(t, models.RunStatusPendingIncomingConfirmations, run.TaskRuns[0].Status)
}

func TestJobRunner_prioritizeSpecParamsOverRequestParams(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	keyStore := cltest.NewKeyStore(t, store.DB)

	pusher := new(mocks.StatsPusher)
	pusher.On("PushNow").Return(nil)

	runExecutor := services.NewRunExecutor(store, keyStore, pusher)
	requestBase := 2
	requestParameter := 10
	specParameter := 100
	j := cltest.NewJobWithWebInitiator()
	taskParams := cltest.JSONFromString(t, fmt.Sprintf(`{"times":%v}`, specParameter))
	j.Tasks = []models.TaskSpec{{Type: adapters.TaskTypeMultiply, Params: taskParams}}
	assert.NoError(t, store.CreateJob(&j))
	run := cltest.NewJobRun(j)
	run.RunRequest.RequestParams = cltest.JSONFromString(t, fmt.Sprintf(`{"times":%v, "result": %v}`, requestParameter, requestBase))
	assert.NoError(t, store.CreateJobRun(&run))

	require.NoError(t, runExecutor.Execute(run.ID))
	run = cltest.WaitForJobRunToComplete(t, store, run)

	actual := run.Result.Data.Get("result").String()
	expected := strconv.FormatUint(uint64(requestBase*specParameter), 10)
	assert.Equal(t, expected, actual)
}
