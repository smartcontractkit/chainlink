package services_test

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	clnull "github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

func TestNewRun(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	input := models.JSON{Result: gjson.Parse(`{"address":"0xdfcfc2b9200dbb10952c2b7cce60fc7260e03c6f"}`)}

	_, bt := cltest.NewBridgeType(t, "timecube", "http://http://timecube.2enp.com/")
	bt.MinimumContractPayment = assets.NewLink(10)
	require.NoError(t, store.CreateBridgeType(bt))

	creationHeight := big.NewInt(1000)

	jobSpec := models.NewJob()
	jobSpec.Tasks = []models.TaskSpec{{
		Type: "timecube",
	}}
	jobSpec.Initiators = []models.Initiator{{
		Type: models.InitiatorEthLog,
	}}

	inputResult := models.RunResult{Data: input}
	run, err := services.NewRun(jobSpec, jobSpec.Initiators[0], inputResult, creationHeight, store, nil)
	assert.NoError(t, err)
	assert.Equal(t, string(models.RunStatusInProgress), string(run.Status))
	assert.Len(t, run.TaskRuns, 1)
	assert.Equal(t, input, run.Overrides)
	assert.False(t, run.TaskRuns[0].Confirmations.Valid)
}

func TestNewRun_MeetsMinimumPayment(t *testing.T) {
	tests := []struct {
		name            string
		MinJobPayment   *assets.Link
		RunPayment      *assets.Link
		meetsMinPayment bool
	}{
		{"insufficient payment", assets.NewLink(100), assets.NewLink(10), false},
		{"sufficient payment (strictly greater)", assets.NewLink(1), assets.NewLink(10), true},
		{"sufficient payment (equal)", assets.NewLink(10), assets.NewLink(10), true},
		{"runs that do not accept payments must return true", assets.NewLink(10), nil, true},
		{"return true when minpayment is not specified in jobspec", nil, assets.NewLink(0), true},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			actual := services.MeetsMinimumPayment(test.MinJobPayment, test.RunPayment)
			assert.Equal(t, test.meetsMinPayment, actual)
		})
	}
}

func TestNewRun_jobSpecMinPayment(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	input := models.JSON{Result: gjson.Parse(`{"address":"0xdfcfc2b9200dbb10952c2b7cce60fc7260e03c6f"}`)}

	tests := []struct {
		name           string
		payment        *assets.Link
		minPayment     *assets.Link
		expectedStatus models.RunStatus
	}{
		{"payment < min payment", assets.NewLink(9), assets.NewLink(10), models.RunStatusErrored},
		{"payment = min payment", assets.NewLink(10), assets.NewLink(10), models.RunStatusInProgress},
		{"payment > min payment", assets.NewLink(11), assets.NewLink(10), models.RunStatusInProgress},
		{"payment is nil", nil, assets.NewLink(10), models.RunStatusInProgress},
		{"minPayment is nil", nil, nil, models.RunStatusInProgress},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			jobSpec := models.NewJob()
			jobSpec.Tasks = []models.TaskSpec{{
				Type: "noop",
			}}
			jobSpec.Initiators = []models.Initiator{{
				Type: models.InitiatorEthLog,
			}}
			jobSpec.MinPayment = test.minPayment

			inputResult := models.RunResult{Data: input}

			run, err := services.NewRun(jobSpec, jobSpec.Initiators[0], inputResult, nil, store, test.payment)
			assert.NoError(t, err)
			assert.Equal(t, string(test.expectedStatus), string(run.Status))
		})
	}
}

func TestNewRun_taskSumPayment(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	_, bta := cltest.NewBridgeType(t, "timecube_a", "http://http://timecube.2enp.com/")
	bta.MinimumContractPayment = assets.NewLink(8)
	require.NoError(t, store.CreateBridgeType(bta))

	_, btb := cltest.NewBridgeType(t, "timecube_b", "http://http://timecube.2enp.com/")
	btb.MinimumContractPayment = assets.NewLink(7)
	require.NoError(t, store.CreateBridgeType(btb))

	store.Config.Set("MINIMUM_CONTRACT_PAYMENT", "1")

	input := models.JSON{Result: gjson.Parse(`{"address":"0xdfcfc2b9200dbb10952c2b7cce60fc7260e03c6f"}`)}

	tests := []struct {
		name           string
		payment        *assets.Link
		expectedStatus models.RunStatus
	}{
		{"payment < min payment", assets.NewLink(15), models.RunStatusErrored},
		{"payment = min payment", assets.NewLink(16), models.RunStatusInProgress},
		{"payment > min payment", assets.NewLink(17), models.RunStatusInProgress},
		{"payment is nil", nil, models.RunStatusInProgress},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			jobSpec := models.NewJob()
			jobSpec.Tasks = []models.TaskSpec{
				{Type: "timecube_a"},
				{Type: "timecube_b"},
				{Type: "ethtx"},
				{Type: "noop"},
			}
			jobSpec.Initiators = []models.Initiator{{
				Type: models.InitiatorEthLog,
			}}

			inputResult := models.RunResult{Data: input}

			run, err := services.NewRun(jobSpec, jobSpec.Initiators[0], inputResult, nil, store, test.payment)
			assert.NoError(t, err)
			assert.Equal(t, string(test.expectedStatus), string(run.Status))
		})
	}
}

func TestNewRun_minimumConfirmations(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	input := models.JSON{Result: gjson.Parse(`{"address":"0xdfcfc2b9200dbb10952c2b7cce60fc7260e03c6f"}`)}
	inputResult := models.RunResult{Data: input}

	creationHeight := big.NewInt(1000)

	tests := []struct {
		name                string
		configConfirmations uint32
		taskConfirmations   uint32
		expectedStatus      models.RunStatus
	}{
		{"creates runnable job", 0, 0, models.RunStatusInProgress},
		{"requires minimum task confirmations", 2, 0, models.RunStatusPendingConfirmations},
		{"requires minimum config confirmations", 0, 2, models.RunStatusPendingConfirmations},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			store.Config.Set("MIN_INCOMING_CONFIRMATIONS", test.configConfirmations)

			jobSpec := cltest.NewJobWithLogInitiator()
			jobSpec.Tasks[0].Confirmations = clnull.Uint32From(test.taskConfirmations)

			run, err := services.NewRun(
				jobSpec,
				jobSpec.Initiators[0],
				inputResult,
				creationHeight,
				store,
				nil)
			assert.NoError(t, err)
			assert.Equal(t, string(test.expectedStatus), string(run.Status))
			require.Len(t, run.TaskRuns, 1)
			max := utils.MaxUint32(test.taskConfirmations, test.configConfirmations)
			assert.Equal(t, max, run.TaskRuns[0].MinimumConfirmations.Uint32)
		})
	}
}

func TestNewRun_startAtAndEndAt(t *testing.T) {
	pastTime := cltest.ParseNullableTime(t, "2000-01-01T00:00:00.000Z")
	futureTime := cltest.ParseNullableTime(t, "3000-01-01T00:00:00.000Z")
	nullTime := null.Time{Valid: false}

	tests := []struct {
		name    string
		startAt null.Time
		endAt   null.Time
		errored bool
	}{
		{"job not started", futureTime, nullTime, true},
		{"job started", pastTime, futureTime, false},
		{"job with no time range", nullTime, nullTime, false},
		{"job ended", nullTime, pastTime, true},
	}

	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	clock := cltest.UseSettableClock(store)
	clock.SetTime(time.Now())

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			job := cltest.NewJobWithWebInitiator()
			job.StartAt = test.startAt
			job.EndAt = test.endAt
			assert.Nil(t, store.CreateJob(&job))

			_, err := services.NewRun(job, job.Initiators[0], models.RunResult{}, nil, store, nil)
			if test.errored {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewRun_noTasksErrorsInsteadOfPanic(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	job := cltest.NewJobWithWebInitiator()
	job.Tasks = []models.TaskSpec{}
	require.NoError(t, store.CreateJob(&job))

	jr, err := services.NewRun(job, job.Initiators[0], models.RunResult{}, nil, store, nil)
	assert.NoError(t, err)
	assert.True(t, jr.Status.Errored())
	assert.True(t, jr.Result.HasError())
}

func TestResumePendingTask(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// reject a run with an invalid state
	jobID := models.NewID()
	runID := models.NewID()
	run := &models.JobRun{
		ID:        runID,
		JobSpecID: jobID,
	}
	err := services.ResumePendingTask(run, store, models.RunResult{})
	assert.Error(t, err)

	// reject a run with no tasks
	run = &models.JobRun{
		ID:        runID,
		JobSpecID: jobID,
		Status:    models.RunStatusPendingBridge,
	}
	err = services.ResumePendingTask(run, store, models.RunResult{})
	assert.Error(t, err)

	// input with error errors run
	run.TaskRuns = []models.TaskRun{models.TaskRun{ID: models.NewID(), JobRunID: runID}}
	err = services.ResumePendingTask(run, store, models.RunResult{CachedJobRunID: runID, Status: models.RunStatusErrored})
	assert.Error(t, err)
	assert.True(t, run.FinishedAt.Valid)

	// completed input with remaining tasks should put task into pending
	run = &models.JobRun{
		ID:        runID,
		JobSpecID: jobID,
		Status:    models.RunStatusPendingBridge,
		TaskRuns:  []models.TaskRun{models.TaskRun{ID: models.NewID(), JobRunID: runID}, models.TaskRun{ID: models.NewID(), JobRunID: runID}},
	}
	input := models.JSON{Result: gjson.Parse(`{"address":"0xdfcfc2b9200dbb10952c2b7cce60fc7260e03c6f"}`)}
	err = services.ResumePendingTask(run, store, models.RunResult{CachedJobRunID: runID, Data: input, Status: models.RunStatusCompleted})
	assert.Error(t, err)
	assert.Equal(t, string(models.RunStatusInProgress), string(run.Status))
	assert.Len(t, run.TaskRuns, 2)
	assert.Equal(t, run.ID, run.TaskRuns[0].Result.CachedJobRunID)
	assert.Equal(t, string(models.RunStatusCompleted), string(run.TaskRuns[0].Result.Status))

	// completed input with no remaining tasks should get marked as complete
	run = &models.JobRun{
		ID:        runID,
		JobSpecID: jobID,
		Status:    models.RunStatusPendingBridge,
		TaskRuns:  []models.TaskRun{models.TaskRun{ID: models.NewID(), JobRunID: runID}},
	}
	err = services.ResumePendingTask(run, store, models.RunResult{CachedJobRunID: runID, Data: input, Status: models.RunStatusCompleted})
	assert.Error(t, err)
	assert.Equal(t, string(models.RunStatusCompleted), string(run.Status))
	assert.True(t, run.FinishedAt.Valid)
	assert.Len(t, run.TaskRuns, 1)
	assert.Equal(t, run.ID, run.TaskRuns[0].Result.CachedJobRunID)
	assert.Equal(t, string(models.RunStatusCompleted), string(run.TaskRuns[0].Result.Status))
}

func TestResumeConfirmingTask(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// reject a run with an invalid state
	jobID := models.NewID()
	runID := models.NewID()
	run := &models.JobRun{
		ID:        runID,
		JobSpecID: jobID,
	}
	err := services.ResumeConfirmingTask(run, store, nil)
	assert.Error(t, err)

	// reject a run with no tasks
	run = &models.JobRun{
		ID:        runID,
		JobSpecID: jobID,
		Status:    models.RunStatusPendingConfirmations,
	}
	err = services.ResumeConfirmingTask(run, store, nil)
	assert.Error(t, err)

	jobSpec := models.JobSpec{ID: models.NewID()}
	require.NoError(t, store.ORM.CreateJob(&jobSpec))

	// leave in pending if not enough confirmations have been met yet
	creationHeight := models.NewBig(big.NewInt(0))
	run = &models.JobRun{
		ID:             models.NewID(),
		JobSpecID:      jobSpec.ID,
		CreationHeight: creationHeight,
		Status:         models.RunStatusPendingConfirmations,
		TaskRuns: []models.TaskRun{models.TaskRun{
			ID:                   models.NewID(),
			MinimumConfirmations: clnull.Uint32From(2),
			TaskSpec: models.TaskSpec{
				JobSpecID: jobSpec.ID,
				Type:      adapters.TaskTypeNoOp,
			},
		}},
	}
	require.NoError(t, store.CreateJobRun(run))
	err = services.ResumeConfirmingTask(run, store, creationHeight.ToInt())
	assert.NoError(t, err)
	assert.Equal(t, string(models.RunStatusPendingConfirmations), string(run.Status))
	assert.Equal(t, uint32(1), run.TaskRuns[0].Confirmations.Uint32)

	// input, should go from pending -> in progress and save the input
	run = &models.JobRun{
		ID:             models.NewID(),
		JobSpecID:      jobSpec.ID,
		CreationHeight: creationHeight,
		Status:         models.RunStatusPendingConfirmations,
		TaskRuns: []models.TaskRun{models.TaskRun{
			ID:                   models.NewID(),
			MinimumConfirmations: clnull.Uint32From(1),
			TaskSpec: models.TaskSpec{
				JobSpecID: jobSpec.ID,
				Type:      adapters.TaskTypeNoOp,
			},
		}},
	}
	observedHeight := big.NewInt(1)
	require.NoError(t, store.CreateJobRun(run))
	err = services.ResumeConfirmingTask(run, store, observedHeight)
	assert.NoError(t, err)
	assert.Equal(t, string(models.RunStatusInProgress), string(run.Status))
}

func TestResumeConnectingTask(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// reject a run with an invalid state
	jobID := models.NewID()
	runID := models.NewID()
	run := &models.JobRun{
		ID:        runID,
		JobSpecID: jobID,
	}
	err := services.ResumeConnectingTask(run, store)
	assert.Error(t, err)

	// reject a run with no tasks
	run = &models.JobRun{
		ID:        runID,
		JobSpecID: jobID,
		Status:    models.RunStatusPendingConnection,
	}
	err = services.ResumeConnectingTask(run, store)
	assert.Error(t, err)

	jobSpec := models.JobSpec{ID: models.NewID()}
	require.NoError(t, store.ORM.CreateJob(&jobSpec))

	taskSpec := models.TaskSpec{Type: adapters.TaskTypeNoOp, JobSpecID: jobSpec.ID}
	// input, should go from pending -> in progress and save the input
	run = &models.JobRun{
		ID:        models.NewID(),
		JobSpecID: jobSpec.ID,
		Status:    models.RunStatusPendingConnection,
		TaskRuns: []models.TaskRun{models.TaskRun{
			ID:       models.NewID(),
			TaskSpec: taskSpec,
		}},
	}
	require.NoError(t, store.CreateJobRun(run))
	err = services.ResumeConnectingTask(run, store)
	assert.NoError(t, err)
	assert.Equal(t, string(models.RunStatusInProgress), string(run.Status))
}

func sleepAdapterParams(t testing.TB, n int) models.JSON {
	d := time.Duration(n)
	json := []byte(fmt.Sprintf(`{"until":%v}`, time.Now().Add(d*time.Second).Unix()))
	return cltest.ParseJSON(t, bytes.NewBuffer(json))
}

func TestQueueSleepingTask(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	store.Clock = cltest.NeverClock{}

	// reject a run with an invalid state
	jobID := models.NewID()
	runID := models.NewID()
	run := &models.JobRun{
		ID:        runID,
		JobSpecID: jobID,
	}
	err := services.QueueSleepingTask(run, store)
	assert.Error(t, err)

	// reject a run with no tasks
	run = &models.JobRun{
		ID:        runID,
		JobSpecID: jobID,
		Status:    models.RunStatusPendingSleep,
	}
	err = services.QueueSleepingTask(run, store)
	assert.Error(t, err)

	jobSpec := models.JobSpec{ID: models.NewID()}
	require.NoError(t, store.ORM.CreateJob(&jobSpec))

	// reject a run that is sleeping but its task is not
	run = &models.JobRun{
		ID:        models.NewID(),
		JobSpecID: jobSpec.ID,
		Status:    models.RunStatusPendingSleep,
		TaskRuns: []models.TaskRun{models.TaskRun{
			ID:       models.NewID(),
			TaskSpec: models.TaskSpec{Type: adapters.TaskTypeSleep, JobSpecID: jobSpec.ID},
		}},
	}
	require.NoError(t, store.CreateJobRun(run))
	err = services.QueueSleepingTask(run, store)
	assert.Error(t, err)

	// error decoding params into adapter
	inputFromTheFuture := cltest.ParseJSON(t, bytes.NewBuffer([]byte(`{"until": -1}`)))
	run = &models.JobRun{
		ID:        models.NewID(),
		JobSpecID: jobSpec.ID,
		Status:    models.RunStatusPendingSleep,
		TaskRuns: []models.TaskRun{
			models.TaskRun{
				ID:     models.NewID(),
				Status: models.RunStatusPendingSleep,
				TaskSpec: models.TaskSpec{
					JobSpecID: jobSpec.ID,
					Type:      adapters.TaskTypeSleep,
					Params:    inputFromTheFuture,
				},
			},
		},
	}
	require.NoError(t, store.CreateJobRun(run))
	err = services.QueueSleepingTask(run, store)
	assert.NoError(t, err)
	assert.Equal(t, string(models.RunStatusErrored), string(run.TaskRuns[0].Status))
	assert.Equal(t, string(models.RunStatusErrored), string(run.Status))

	// mark run as pending, task as completed if duration has already elapsed
	run = &models.JobRun{
		ID:        models.NewID(),
		JobSpecID: jobSpec.ID,
		Status:    models.RunStatusPendingSleep,
		TaskRuns: []models.TaskRun{models.TaskRun{
			ID:       models.NewID(),
			Status:   models.RunStatusPendingSleep,
			TaskSpec: models.TaskSpec{Type: adapters.TaskTypeSleep, JobSpecID: jobSpec.ID},
		}},
	}
	require.NoError(t, store.CreateJobRun(run))
	err = services.QueueSleepingTask(run, store)
	assert.NoError(t, err)
	assert.Equal(t, string(models.RunStatusCompleted), string(run.TaskRuns[0].Status))
	assert.Equal(t, string(models.RunStatusInProgress), string(run.Status))

	runRequest, open := <-store.RunChannel.Receive()
	assert.True(t, open)
	assert.Equal(t, run.ID, runRequest.ID)

}

func TestQueueSleepingTaskA_CompletesSleepingTaskAfterDurationElapsed_Happy(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	store.Clock = cltest.NeverClock{}

	jobSpec := models.JobSpec{ID: models.NewID()}
	require.NoError(t, store.ORM.CreateJob(&jobSpec))

	// queue up next run if duration has not elapsed yet
	clock := cltest.UseSettableClock(store)
	store.Clock = clock
	clock.SetTime(time.Time{})

	inputFromTheFuture := sleepAdapterParams(t, 60)
	run := &models.JobRun{
		ID:        models.NewID(),
		JobSpecID: jobSpec.ID,
		Status:    models.RunStatusPendingSleep,
		TaskRuns: []models.TaskRun{
			models.TaskRun{
				ID:     models.NewID(),
				Status: models.RunStatusPendingSleep,
				TaskSpec: models.TaskSpec{
					JobSpecID: jobSpec.ID,
					Type:      adapters.TaskTypeSleep,
					Params:    inputFromTheFuture,
				},
			},
		},
	}
	require.NoError(t, store.CreateJobRun(run))
	err := services.QueueSleepingTask(run, store)
	assert.NoError(t, err)
	assert.Equal(t, string(models.RunStatusPendingSleep), string(run.TaskRuns[0].Status))
	assert.Equal(t, string(models.RunStatusPendingSleep), string(run.Status))

	// force the duration elapse
	clock.SetTime((time.Time{}).Add(math.MaxInt64))
	runRequest, open := <-store.RunChannel.Receive()
	assert.True(t, open)
	assert.Equal(t, run.ID, runRequest.ID)

	*run, err = store.ORM.FindJobRun(run.ID)
	assert.NoError(t, err)
	assert.Equal(t, string(models.RunStatusCompleted), string(run.TaskRuns[0].Status))
	assert.Equal(t, string(models.RunStatusInProgress), string(run.Status))
}

func TestQueueSleepingTaskA_CompletesSleepingTaskAfterDurationElapsed_Archived(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	store.Clock = cltest.NeverClock{}

	jobSpec := models.JobSpec{ID: models.NewID()}
	require.NoError(t, store.ORM.CreateJob(&jobSpec))

	// queue up next run if duration has not elapsed yet
	clock := cltest.UseSettableClock(store)
	store.Clock = clock
	clock.SetTime(time.Time{})

	inputFromTheFuture := sleepAdapterParams(t, 60)
	run := &models.JobRun{
		ID:        models.NewID(),
		JobSpecID: jobSpec.ID,
		Status:    models.RunStatusPendingSleep,
		TaskRuns: []models.TaskRun{
			models.TaskRun{
				ID:     models.NewID(),
				Status: models.RunStatusPendingSleep,
				TaskSpec: models.TaskSpec{
					JobSpecID: jobSpec.ID,
					Type:      adapters.TaskTypeSleep,
					Params:    inputFromTheFuture,
				},
			},
		},
	}
	require.NoError(t, store.CreateJobRun(run))
	require.NoError(t, store.ArchiveJob(jobSpec.ID))

	unscoped := store.Unscoped()
	err := services.QueueSleepingTask(run, unscoped)
	assert.NoError(t, err)
	assert.Equal(t, string(models.RunStatusPendingSleep), string(run.TaskRuns[0].Status))
	assert.Equal(t, string(models.RunStatusPendingSleep), string(run.Status))

	// force the duration elapse
	clock.SetTime((time.Time{}).Add(math.MaxInt64))
	runRequest, open := <-store.RunChannel.Receive()
	assert.True(t, open)
	assert.Equal(t, run.ID, runRequest.ID)

	require.Error(t, utils.JustError(store.FindJobRun(run.ID)), "archived runs should not be visible to normal store")

	*run, err = unscoped.FindJobRun(run.ID)
	assert.NoError(t, err)
	assert.Equal(t, string(models.RunStatusCompleted), string(run.TaskRuns[0].Status))
	assert.Equal(t, string(models.RunStatusInProgress), string(run.Status))
}

func TestExecuteJob_DoesNotSaveToTaskSpec(t *testing.T) {
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()

	store := app.Store
	eth := cltest.MockEthOnStore(t, store)
	eth.Register("eth_chainId", store.Config.ChainID())

	app.Start()

	job := cltest.NewJobWithWebInitiator()
	job.Tasks = []models.TaskSpec{cltest.NewTask(t, "NoOp")} // empty params
	require.NoError(t, store.CreateJob(&job))

	initr := job.Initiators[0]
	jr, err := services.ExecuteJob(
		job,
		initr,
		cltest.RunResultWithData(`{"random": "input"}`),
		nil,
		store,
	)
	require.NoError(t, err)
	cltest.WaitForJobRunToComplete(t, store, *jr)

	retrievedJob, err := store.FindJob(job.ID)
	require.NoError(t, err)
	require.Len(t, job.Tasks, 1)
	require.Len(t, retrievedJob.Tasks, 1)
	assert.Equal(t, job.Tasks[0].Params, retrievedJob.Tasks[0].Params)
}

func TestExecuteJobWithRunRequest(t *testing.T) {
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()

	store := app.Store
	eth := cltest.MockEthOnStore(t, store)
	eth.Register("eth_chainId", store.Config.ChainID())

	app.Start()

	job := cltest.NewJobWithRunLogInitiator()
	job.Tasks = []models.TaskSpec{cltest.NewTask(t, "NoOp")} // empty params
	require.NoError(t, store.CreateJob(&job))

	requestID := "RequestID"
	initr := job.Initiators[0]
	rr := models.NewRunRequest()
	rr.RequestID = &requestID
	jr, err := services.ExecuteJobWithRunRequest(
		job,
		initr,
		cltest.RunResultWithData(`{"random": "input"}`),
		nil,
		store,
		rr,
	)
	require.NoError(t, err)
	updatedJR := cltest.WaitForJobRunToComplete(t, store, *jr)
	assert.Equal(t, rr.RequestID, updatedJR.RunRequest.RequestID)
}

func TestExecuteJobWithRunRequest_fromRunLog_Happy(t *testing.T) {

	initiatingTxHash := cltest.NewHash()
	triggeringBlockHash := cltest.NewHash()
	otherBlockHash := cltest.NewHash()

	tests := []struct {
		name             string
		logBlockHash     common.Hash
		receiptBlockHash common.Hash
		wantStatus       models.RunStatus
	}{
		{
			name:             "main chain",
			logBlockHash:     triggeringBlockHash,
			receiptBlockHash: triggeringBlockHash,
			wantStatus:       models.RunStatusCompleted,
		},
		{
			name:             "ommered chain",
			logBlockHash:     triggeringBlockHash,
			receiptBlockHash: otherBlockHash,
			wantStatus:       models.RunStatusErrored,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, cfgCleanup := cltest.NewConfig(t)
			defer cfgCleanup()
			minimumConfirmations := uint32(2)
			config.Set("MIN_INCOMING_CONFIRMATIONS", minimumConfirmations)
			app, cleanup := cltest.NewApplicationWithConfig(t, config)
			defer cleanup()

			eth := app.MockEthCallerSubscriber()
			app.Start()

			store := app.GetStore()
			job := cltest.NewJobWithRunLogInitiator()
			job.Tasks = []models.TaskSpec{cltest.NewTask(t, "NoOp")}
			require.NoError(t, store.CreateJob(&job))

			creationHeight := big.NewInt(1)
			requestID := "RequestID"
			initr := job.Initiators[0]
			rr := models.NewRunRequest()
			rr.RequestID = &requestID
			rr.TxHash = &initiatingTxHash
			rr.BlockHash = &test.logBlockHash
			jr, err := services.ExecuteJobWithRunRequest(
				job,
				initr,
				cltest.RunResultWithData(`{"random": "input"}`),
				creationHeight,
				store,
				rr,
			)
			require.NoError(t, err)
			cltest.WaitForJobRunToPendConfirmations(t, app.Store, *jr)

			confirmedReceipt := models.TxReceipt{
				Hash:        initiatingTxHash,
				BlockHash:   &test.receiptBlockHash,
				BlockNumber: cltest.Int(3),
			}
			eth.Context("validateOnMainChain", func(ethMock *cltest.EthMock) {
				eth.Register("eth_getTransactionReceipt", confirmedReceipt)
			})

			err = services.ResumeConfirmingTask(jr, store, big.NewInt(2))
			require.NoError(t, err)
			updatedJR := cltest.WaitForJobRunStatus(t, store, *jr, test.wantStatus)
			assert.Equal(t, rr.RequestID, updatedJR.RunRequest.RequestID)
			assert.Equal(t, minimumConfirmations, updatedJR.TaskRuns[0].MinimumConfirmations.Uint32)
			assert.True(t, updatedJR.TaskRuns[0].MinimumConfirmations.Valid)
			assert.Equal(t, minimumConfirmations, updatedJR.TaskRuns[0].Confirmations.Uint32, "task run should track its current confirmations")
			assert.True(t, updatedJR.TaskRuns[0].Confirmations.Valid)
			assert.True(t, eth.AllCalled(), eth.Remaining())
		})
	}
}

func TestExecuteJobWithRunRequest_fromRunLog_ConnectToLaggingEthNode(t *testing.T) {

	initiatingTxHash := cltest.NewHash()
	triggeringBlockHash := cltest.NewHash()

	config, cfgCleanup := cltest.NewConfig(t)
	defer cfgCleanup()
	minimumConfirmations := uint32(2)
	config.Set("MIN_INCOMING_CONFIRMATIONS", minimumConfirmations)
	app, cleanup := cltest.NewApplicationWithConfig(t, config)
	defer cleanup()

	eth := app.MockEthCallerSubscriber()
	app.MockStartAndConnect()

	store := app.GetStore()
	job := cltest.NewJobWithRunLogInitiator()
	job.Tasks = []models.TaskSpec{cltest.NewTask(t, "NoOp")}
	require.NoError(t, store.CreateJob(&job))

	requestID := "RequestID"
	initr := job.Initiators[0]
	rr := models.NewRunRequest()
	rr.RequestID = &requestID
	rr.TxHash = &initiatingTxHash
	rr.BlockHash = &triggeringBlockHash

	futureCreationHeight := big.NewInt(9)
	pastCurrentHeight := big.NewInt(1)

	jr, err := services.ExecuteJobWithRunRequest(
		job,
		initr,
		cltest.RunResultWithData(`{"random": "input"}`),
		futureCreationHeight,
		store,
		rr,
	)
	require.NoError(t, err)
	cltest.WaitForJobRunToPendConfirmations(t, app.Store, *jr)

	err = services.ResumeConfirmingTask(jr, store, pastCurrentHeight)
	require.NoError(t, err)
	updatedJR := cltest.WaitForJobRunToPendConfirmations(t, app.Store, *jr)
	assert.True(t, updatedJR.TaskRuns[0].Confirmations.Valid)
	assert.Equal(t, uint32(0), updatedJR.TaskRuns[0].Confirmations.Uint32)
	assert.True(t, eth.AllCalled(), eth.Remaining())
}
