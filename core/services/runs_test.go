package services_test

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
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
	run, err := services.NewRun(jobSpec, jobSpec.Initiators[0], inputResult, creationHeight, store)
	assert.NoError(t, err)
	assert.Equal(t, string(models.RunStatusInProgress), string(run.Status))
	assert.Len(t, run.TaskRuns, 1)
	assert.Equal(t, input, run.Overrides.Data)
}

func TestNewRun_requiredPayment(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	input := models.JSON{Result: gjson.Parse(`{"address":"0xdfcfc2b9200dbb10952c2b7cce60fc7260e03c6f"}`)}

	_, bt := cltest.NewBridgeType(t, "timecube", "http://http://timecube.2enp.com/")
	bt.MinimumContractPayment = assets.NewLink(10)
	require.NoError(t, store.CreateBridgeType(bt))

	tests := []struct {
		name           string
		payment        *assets.Link
		minimumPayment *assets.Link
		expectedStatus models.RunStatus
	}{
		{"creates runnable job", nil, assets.NewLink(0), models.RunStatusInProgress},
		{"insufficient payment as specified by config", assets.NewLink(9), assets.NewLink(10), models.RunStatusErrored},
		{"sufficient payment as specified by config", assets.NewLink(10), assets.NewLink(10), models.RunStatusInProgress},
		{"insufficient payment as specified by adapter", assets.NewLink(9), assets.NewLink(0), models.RunStatusErrored},
		{"sufficient payment as specified by adapter", assets.NewLink(10), assets.NewLink(0), models.RunStatusInProgress},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			store.Config.Set("MINIMUM_CONTRACT_PAYMENT", test.minimumPayment)

			jobSpec := models.NewJob()
			jobSpec.Tasks = []models.TaskSpec{{
				Type: "timecube",
			}}
			jobSpec.Initiators = []models.Initiator{{
				Type: models.InitiatorEthLog,
			}}

			inputResult := models.RunResult{Data: input, Amount: test.payment}

			run, err := services.NewRun(jobSpec, jobSpec.Initiators[0], inputResult, nil, store)
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
		configConfirmations uint64
		taskConfirmations   uint64
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
			jobSpec.Tasks[0].Confirmations = test.taskConfirmations

			run, err := services.NewRun(
				jobSpec,
				jobSpec.Initiators[0],
				inputResult,
				creationHeight,
				store)
			assert.NoError(t, err)
			assert.Equal(t, string(test.expectedStatus), string(run.Status))
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

			_, err := services.NewRun(job, job.Initiators[0], models.RunResult{}, nil, store)
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

	jr, err := services.NewRun(job, job.Initiators[0], models.RunResult{}, nil, store)
	assert.NoError(t, err)
	assert.True(t, jr.Status.Errored())
	assert.True(t, jr.Result.HasError())
}

func TestResumePendingTask(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()

	// reject a run with an invalid state
	run := &models.JobRun{}
	err := services.ResumePendingTask(run, store, models.RunResult{})
	assert.Error(t, err)

	// reject a run with no tasks
	run = &models.JobRun{Status: models.RunStatusPendingBridge}
	err = services.ResumePendingTask(run, store, models.RunResult{})
	assert.Error(t, err)

	// input with error errors run
	run = &models.JobRun{
		Status:   models.RunStatusPendingBridge,
		TaskRuns: []models.TaskRun{models.TaskRun{}},
	}
	err = services.ResumePendingTask(run, store, models.RunResult{Status: models.RunStatusErrored})
	assert.Error(t, err)
	assert.True(t, run.FinishedAt.Valid)

	// completed input with remaining tasks should put task into pending
	run = &models.JobRun{
		Status:   models.RunStatusPendingBridge,
		TaskRuns: []models.TaskRun{models.TaskRun{}, models.TaskRun{}},
	}
	input := models.JSON{Result: gjson.Parse(`{"address":"0xdfcfc2b9200dbb10952c2b7cce60fc7260e03c6f"}`)}
	err = services.ResumePendingTask(run, store, models.RunResult{Data: input, Status: models.RunStatusCompleted})
	assert.Error(t, err)
	assert.Equal(t, string(models.RunStatusInProgress), string(run.Status))
	assert.Len(t, run.TaskRuns, 2)
	assert.Equal(t, run.ID, run.TaskRuns[0].Result.CachedJobRunID)
	assert.Equal(t, string(models.RunStatusCompleted), string(run.TaskRuns[0].Result.Status))

	// completed input with no remaining tasks should get marked as complete
	run = &models.JobRun{
		Status:   models.RunStatusPendingBridge,
		TaskRuns: []models.TaskRun{models.TaskRun{}},
	}
	err = services.ResumePendingTask(run, store, models.RunResult{Data: input, Status: models.RunStatusCompleted})
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
	run := &models.JobRun{}
	err := services.ResumeConfirmingTask(run, store, nil)
	assert.Error(t, err)

	// reject a run with no tasks
	run = &models.JobRun{Status: models.RunStatusPendingConfirmations}
	err = services.ResumeConfirmingTask(run, store, nil)
	assert.Error(t, err)

	jobSpec := models.JobSpec{ID: utils.NewBytes32ID()}
	require.NoError(t, store.ORM.CreateJob(&jobSpec))

	// leave in pending if not enough confirmations have been met yet
	creationHeight := models.NewBig(big.NewInt(0))
	run = &models.JobRun{
		ID:             utils.NewBytes32ID(),
		JobSpecID:      jobSpec.ID,
		CreationHeight: creationHeight,
		Status:         models.RunStatusPendingConfirmations,
		TaskRuns: []models.TaskRun{models.TaskRun{
			ID:                   utils.NewBytes32ID(),
			MinimumConfirmations: 2,
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

	// input, should go from pending -> in progress and save the input
	run = &models.JobRun{
		ID:             utils.NewBytes32ID(),
		JobSpecID:      jobSpec.ID,
		CreationHeight: creationHeight,
		Status:         models.RunStatusPendingConfirmations,
		TaskRuns: []models.TaskRun{models.TaskRun{
			ID:                   utils.NewBytes32ID(),
			MinimumConfirmations: 1,
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
	run := &models.JobRun{}
	err := services.ResumeConnectingTask(run, store)
	assert.Error(t, err)

	// reject a run with no tasks
	run = &models.JobRun{Status: models.RunStatusPendingConnection}
	err = services.ResumeConnectingTask(run, store)
	assert.Error(t, err)

	jobSpec := models.JobSpec{ID: utils.NewBytes32ID()}
	require.NoError(t, store.ORM.CreateJob(&jobSpec))

	taskSpec := models.TaskSpec{Type: adapters.TaskTypeNoOp, JobSpecID: jobSpec.ID}
	// input, should go from pending -> in progress and save the input
	run = &models.JobRun{
		ID:        utils.NewBytes32ID(),
		JobSpecID: jobSpec.ID,
		Status:    models.RunStatusPendingConnection,
		TaskRuns: []models.TaskRun{models.TaskRun{
			ID:       utils.NewBytes32ID(),
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
	run := &models.JobRun{}
	err := services.QueueSleepingTask(run, store)
	assert.Error(t, err)

	// reject a run with no tasks
	run = &models.JobRun{Status: models.RunStatusPendingSleep}
	err = services.QueueSleepingTask(run, store)
	assert.Error(t, err)

	jobSpec := models.JobSpec{ID: utils.NewBytes32ID()}
	require.NoError(t, store.ORM.CreateJob(&jobSpec))

	// reject a run that is sleeping but its task is not
	run = &models.JobRun{
		ID:        utils.NewBytes32ID(),
		JobSpecID: jobSpec.ID,
		Status:    models.RunStatusPendingSleep,
		TaskRuns: []models.TaskRun{models.TaskRun{
			ID:       utils.NewBytes32ID(),
			TaskSpec: models.TaskSpec{Type: adapters.TaskTypeSleep, JobSpecID: jobSpec.ID},
		}},
	}
	require.NoError(t, store.CreateJobRun(run))
	err = services.QueueSleepingTask(run, store)
	assert.Error(t, err)

	// error decoding params into adapter
	inputFromTheFuture := cltest.ParseJSON(t, bytes.NewBuffer([]byte(`{"until": -1}`)))
	run = &models.JobRun{
		ID:        utils.NewBytes32ID(),
		JobSpecID: jobSpec.ID,
		Status:    models.RunStatusPendingSleep,
		TaskRuns: []models.TaskRun{
			models.TaskRun{
				ID:     utils.NewBytes32ID(),
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
		ID:        utils.NewBytes32ID(),
		JobSpecID: jobSpec.ID,
		Status:    models.RunStatusPendingSleep,
		TaskRuns: []models.TaskRun{models.TaskRun{
			ID:       utils.NewBytes32ID(),
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

	jobSpec := models.JobSpec{ID: utils.NewBytes32ID()}
	require.NoError(t, store.ORM.CreateJob(&jobSpec))

	// queue up next run if duration has not elapsed yet
	clock := cltest.UseSettableClock(store)
	store.Clock = clock
	clock.SetTime(time.Time{})

	inputFromTheFuture := sleepAdapterParams(t, 60)
	run := &models.JobRun{
		ID:        utils.NewBytes32ID(),
		JobSpecID: jobSpec.ID,
		Status:    models.RunStatusPendingSleep,
		TaskRuns: []models.TaskRun{
			models.TaskRun{
				ID:     utils.NewBytes32ID(),
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

	jobSpec := models.JobSpec{ID: utils.NewBytes32ID()}
	require.NoError(t, store.ORM.CreateJob(&jobSpec))

	// queue up next run if duration has not elapsed yet
	clock := cltest.UseSettableClock(store)
	store.Clock = clock
	clock.SetTime(time.Time{})

	inputFromTheFuture := sleepAdapterParams(t, 60)
	run := &models.JobRun{
		ID:        utils.NewBytes32ID(),
		JobSpecID: jobSpec.ID,
		Status:    models.RunStatusPendingSleep,
		TaskRuns: []models.TaskRun{
			models.TaskRun{
				ID:     utils.NewBytes32ID(),
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
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	app.Start()
	store := app.Store

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
	t.Parallel()
	app, cleanup := cltest.NewApplication(t)
	defer cleanup()
	app.Start()
	store := app.Store

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

func TestExecuteJobWithRunRequest_fromRunLog_mainChain(t *testing.T) {
	t.Parallel()

	config, cfgCleanup := cltest.NewConfig(t)
	defer cfgCleanup()
	config.Set("MIN_INCOMING_CONFIRMATIONS", 2)
	app, cleanup := cltest.NewApplicationWithConfig(t, config)
	defer cleanup()

	initiatingTxHash := cltest.NewHash()
	eth := app.MockEthClient()
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
		Hash:        *rr.TxHash,
		BlockNumber: cltest.Int(3),
	}
	eth.Context("validateOnMainChain", func(ethMock *cltest.EthMock) {
		eth.Register("eth_getTransactionReceipt", confirmedReceipt)
	})

	err = services.ResumeConfirmingTask(jr, store, big.NewInt(2))
	require.NoError(t, err)
	updatedJR := cltest.WaitForJobRunToComplete(t, store, *jr)
	assert.Equal(t, rr.RequestID, updatedJR.RunRequest.RequestID)
}

func TestExecuteJobWithRunRequest_fromRunLog_uncled(t *testing.T) {
	t.Parallel()

	config, cfgCleanup := cltest.NewConfig(t)
	defer cfgCleanup()
	config.Set("MIN_INCOMING_CONFIRMATIONS", 2)
	app, cleanup := cltest.NewApplicationWithConfig(t, config)
	defer cleanup()

	initiatingTxHash := cltest.NewHash()
	eth := app.MockEthClient()
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

	unconfirmedReceipt := models.TxReceipt{}
	eth.Context("validateOnMainChain", func(ethMock *cltest.EthMock) {
		eth.Register("eth_getTransactionReceipt", unconfirmedReceipt)
	})

	err = services.ResumeConfirmingTask(jr, store, big.NewInt(2))
	require.NoError(t, err)
	updatedJR := cltest.WaitForJobRunStatus(t, store, *jr, models.RunStatusErrored)
	assert.Equal(t, rr.RequestID, updatedJR.RunRequest.RequestID)
}
