package models_test

import (
	"encoding/json"
	"errors"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	null "gopkg.in/guregu/null.v3"
)

func TestJobRuns_RetrievingFromDBWithError(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))
	jr := job.NewRun(job.Initiators[0])
	jr.JobSpecID = job.ID
	jr.Result.ErrorMessage = null.StringFrom("bad idea")
	err := store.CreateJobRun(&jr)
	require.NoError(t, err)

	run, err := store.FindJobRun(jr.ID)
	require.NoError(t, err)
	assert.True(t, run.Result.ErrorMessage.Valid)
	assert.Equal(t, "bad idea", run.Result.ErrorString())
}

func TestJobRuns_RetrievingFromDBWithData(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithWebInitiator()
	err := store.CreateJob(&job)
	initr := job.Initiators[0]
	assert.NoError(t, err)

	jr := job.NewRun(initr)
	data := `{"result":"11850.00"}`
	jr.Result = models.RunResult{Data: cltest.JSONFromString(t, data)}
	err = store.CreateJobRun(&jr)
	assert.NoError(t, err)

	run, err := store.FindJobRun(jr.ID)
	assert.NoError(t, err)
	assert.False(t, run.Result.ErrorMessage.Valid)
	assert.JSONEq(t, data, run.Result.Data.String())
}

func TestJobRuns_SavesASyncEvent(t *testing.T) {
	t.Parallel()
	config, _ := cltest.NewConfig(t)
	config.Set("EXPLORER_URL", "http://localhost:4201")
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	job := cltest.NewJobWithWebInitiator()
	err := store.CreateJob(&job)
	initr := job.Initiators[0]
	assert.NoError(t, err)

	jr := job.NewRun(initr)
	err = store.CreateJobRun(&jr)
	assert.NoError(t, err)

	var events []*models.SyncEvent
	err = store.AllSyncEvents(func(event *models.SyncEvent) error {
		events = append(events, event)
		return nil
	})
	require.NoError(t, err)
	require.Len(t, events, 1)

	var data map[string]interface{}
	err = json.Unmarshal([]byte(events[0].Body), &data)
	require.NoError(t, err)

	var recoveredJobRun models.JobRun
	err = json.Unmarshal([]byte(events[0].Body), &recoveredJobRun)
	require.NoError(t, err)
	assert.Equal(t, jr.Result.Data, recoveredJobRun.Result.Data)

	assert.Contains(t, data, "id")
	assert.Contains(t, data, "runId")
	assert.Contains(t, data, "jobId")
	assert.Contains(t, data, "status")
}

func TestJobRuns_SkipsEventSaveIfURLBlank(t *testing.T) {
	t.Parallel()
	config, _ := cltest.NewConfig(t)
	config.Set("EXPLORER_URL", "")
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	job := cltest.NewJobWithWebInitiator()
	err := store.CreateJob(&job)
	initr := job.Initiators[0]
	assert.NoError(t, err)

	jr := job.NewRun(initr)
	data := `{"result":"921.02"}`
	jr.Result = models.RunResult{Data: cltest.JSONFromString(t, data)}
	err = store.CreateJobRun(&jr)
	assert.NoError(t, err)

	var events []*models.SyncEvent
	err = store.AllSyncEvents(func(event *models.SyncEvent) error {
		events = append(events, event)
		return nil
	})
	require.NoError(t, err)
	require.Len(t, events, 0)
}

func TestForLogger(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithWebInitiator()
	require.NoError(t, store.CreateJob(&job))
	jr := job.NewRun(job.Initiators[0])
	jr.JobSpecID = job.ID
	linkReward := assets.NewLink(5)

	jr.Result = models.RunResult{Data: cltest.JSONFromString(t, `{"result":"11850.00"}`)}
	jr.Payment = linkReward
	logsBeforeCompletion := jr.ForLogger()
	require.Len(t, logsBeforeCompletion, 6)
	assert.Equal(t, logsBeforeCompletion[0], "job")
	assert.Equal(t, logsBeforeCompletion[1], jr.JobSpecID.String())
	assert.Equal(t, logsBeforeCompletion[2], "run")
	assert.Equal(t, logsBeforeCompletion[3], jr.ID.String())
	assert.Equal(t, logsBeforeCompletion[4], "status")
	assert.Equal(t, logsBeforeCompletion[5], jr.Status)

	jr.Status = "completed"
	logsAfterCompletion := jr.ForLogger()
	require.Len(t, logsAfterCompletion, 8)
	assert.Equal(t, logsAfterCompletion[4], "status")
	assert.Equal(t, logsAfterCompletion[5], jr.Status)
	assert.Equal(t, logsAfterCompletion[6], "link_earned")
	assert.Equal(t, logsAfterCompletion[7], linkReward)

	jr.CreationHeight = models.NewBig(big.NewInt(5))
	jr.ObservedHeight = models.NewBig(big.NewInt(10))
	logsWithBlockHeights := jr.ForLogger()
	require.Len(t, logsWithBlockHeights, 12)
	assert.Equal(t, logsWithBlockHeights[6], "creation_height")
	assert.Equal(t, logsWithBlockHeights[7], big.NewInt(5))
	assert.Equal(t, logsWithBlockHeights[8], "observed_height")
	assert.Equal(t, logsWithBlockHeights[9], big.NewInt(10))

	run := job.NewRun(job.Initiators[0])
	run.SetError(errors.New("bad idea"))
	logsWithErr := run.ForLogger()
	assert.Equal(t, logsWithErr[6], "job_error")
	assert.Equal(t, logsWithErr[7], run.Result.ErrorString())
}

func TestJobRun_NextTaskRun(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	jobRunner, cleanup := cltest.NewJobRunner(store)
	defer cleanup()
	jobRunner.Start()

	job := cltest.NewJobWithWebInitiator()
	job.Tasks = []models.TaskSpec{
		{Type: adapters.TaskTypeNoOp},
		{Type: adapters.TaskTypeNoOpPend},
		{Type: adapters.TaskTypeNoOp},
	}
	assert.NoError(t, store.CreateJob(&job))
	run := job.NewRun(job.Initiators[0])
	assert.NoError(t, store.CreateJobRun(&run))
	assert.Equal(t, &run.TaskRuns[0], run.NextTaskRun())

	store.RunChannel.Send(run.ID)
	cltest.WaitForJobRunStatus(t, store, run, models.RunStatusPendingConfirmations)

	run, err := store.FindJobRun(run.ID)
	assert.NoError(t, err)
	assert.Equal(t, &run.TaskRuns[1], run.NextTaskRun())
}

func TestJobRun_ApplyOutput_CompletedWithNoTasksRemaining(t *testing.T) {
	t.Parallel()

	job := cltest.NewJobWithWebInitiator()
	jobRun := job.NewRun(job.Initiators[0])

	result := models.NewRunOutputComplete(models.JSON{})
	jobRun.TaskRuns[0].ApplyOutput(result)
	err := jobRun.ApplyOutput(result)
	assert.NoError(t, err)
	assert.True(t, jobRun.FinishedAt.Valid)
}

func TestJobRun_ApplyOutput_CompletedWithTasksRemaining(t *testing.T) {
	t.Parallel()

	job := cltest.NewJobWithWebInitiator()
	jobRun := job.NewRun(job.Initiators[0])

	result := models.NewRunOutputComplete(models.JSON{})
	err := jobRun.ApplyOutput(result)
	assert.NoError(t, err)
	assert.False(t, jobRun.FinishedAt.Valid)
	assert.Equal(t, jobRun.Status, models.RunStatusInProgress)
}

func TestJobRun_ApplyOutput_ErrorSetsFinishedAt(t *testing.T) {
	t.Parallel()

	job := cltest.NewJobWithWebInitiator()
	jobRun := job.NewRun(job.Initiators[0])
	jobRun.Status = models.RunStatusErrored

	result := models.NewRunOutputError(errors.New("oh futz"))
	err := jobRun.ApplyOutput(result)
	assert.NoError(t, err)
	assert.True(t, jobRun.FinishedAt.Valid)
}
