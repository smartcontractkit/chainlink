package services_test

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

func TestNewRun(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	input := models.JSON{Result: gjson.Parse(`{"address":"0xdfcfc2b9200dbb10952c2b7cce60fc7260e03c6f"}`)}

	bt := cltest.NewBridgeType("timecube", "http://http://timecube.2enp.com/")
	bt.MinimumContractPayment = *assets.NewLink(10)
	assert.Nil(t, store.Save(&bt))

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
	store, cleanup := cltest.NewStore()
	defer cleanup()

	input := models.JSON{Result: gjson.Parse(`{"address":"0xdfcfc2b9200dbb10952c2b7cce60fc7260e03c6f"}`)}

	bt := cltest.NewBridgeType("timecube", "http://http://timecube.2enp.com/")
	bt.MinimumContractPayment = *assets.NewLink(10)
	assert.Nil(t, store.Save(&bt))

	tests := []struct {
		name           string
		payment        *assets.Link
		minimumPayment assets.Link
		expectedStatus models.RunStatus
	}{
		{"creates runnable job", nil, *assets.NewLink(0), models.RunStatusInProgress},
		{"insufficient payment as specified by config", assets.NewLink(9), *assets.NewLink(10), models.RunStatusErrored},
		{"sufficient payment as specified by config", assets.NewLink(10), *assets.NewLink(10), models.RunStatusInProgress},
		{"insufficient payment as specified by adapter", assets.NewLink(9), *assets.NewLink(0), models.RunStatusErrored},
		{"sufficient payment as specified by adapter", assets.NewLink(10), *assets.NewLink(0), models.RunStatusInProgress},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			store.Config.MinimumContractPayment = test.minimumPayment

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
	store, cleanup := cltest.NewStore()
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
			store.Config.MinIncomingConfirmations = test.configConfirmations

			jobSpec, initiator := cltest.NewJobWithLogInitiator()
			jobSpec.Tasks[0].Confirmations = test.taskConfirmations

			run, err := services.NewRun(jobSpec, initiator, inputResult, creationHeight, store)
			assert.NoError(t, err)
			assert.Equal(t, string(test.expectedStatus), string(run.Status))
		})
	}
}

func TestNewRun_startAtAndEndAt(t *testing.T) {
	pastTime := cltest.ParseNullableTime("2000-01-01T00:00:00.000Z")
	futureTime := cltest.ParseNullableTime("3000-01-01T00:00:00.000Z")
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

	store, cleanup := cltest.NewStore()
	defer cleanup()
	clock := cltest.UseSettableClock(store)
	clock.SetTime(time.Now())

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			job, initr := cltest.NewJobWithWebInitiator()
			job.StartAt = test.startAt
			job.EndAt = test.endAt
			assert.Nil(t, store.SaveJob(&job))

			_, err := services.NewRun(job, initr, models.RunResult{}, nil, store)
			if test.errored {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestResumePendingTask(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	// reject a run with an invalid state
	run := models.JobRun{}
	err := services.ResumePendingTask(&run, store, models.RunResult{})
	assert.Error(t, err)

	// reject a run with no tasks
	run = models.JobRun{Status: models.RunStatusPendingBridge}
	err = services.ResumePendingTask(&run, store, models.RunResult{})
	assert.Error(t, err)

	// input with error errors run
	run = models.JobRun{
		Status:   models.RunStatusPendingBridge,
		TaskRuns: []models.TaskRun{models.TaskRun{}},
	}
	err = services.ResumePendingTask(&run, store, models.RunResult{Status: models.RunStatusErrored})
	assert.Error(t, err)

	// completed input with remaining tasks should put task into pending
	run = models.JobRun{
		Status:   models.RunStatusPendingBridge,
		TaskRuns: []models.TaskRun{models.TaskRun{}, models.TaskRun{}},
	}
	input := models.JSON{Result: gjson.Parse(`{"address":"0xdfcfc2b9200dbb10952c2b7cce60fc7260e03c6f"}`)}
	err = services.ResumePendingTask(&run, store, models.RunResult{Data: input, Status: models.RunStatusCompleted})
	assert.Error(t, err)
	assert.Equal(t, string(models.RunStatusInProgress), string(run.Status))
	assert.Len(t, run.TaskRuns, 2)
	assert.Equal(t, run.ID, run.TaskRuns[0].Result.JobRunID)
	assert.Equal(t, string(models.RunStatusCompleted), string(run.TaskRuns[0].Result.Status))

	// completed input with no remaining tasks should get marked as complete
	run = models.JobRun{
		Status:   models.RunStatusPendingBridge,
		TaskRuns: []models.TaskRun{models.TaskRun{}},
	}
	err = services.ResumePendingTask(&run, store, models.RunResult{Data: input, Status: models.RunStatusCompleted})
	assert.Error(t, err)
	assert.Equal(t, string(models.RunStatusCompleted), string(run.Status))
	assert.Len(t, run.TaskRuns, 1)
	assert.Equal(t, run.ID, run.TaskRuns[0].Result.JobRunID)
	assert.Equal(t, string(models.RunStatusCompleted), string(run.TaskRuns[0].Result.Status))
}

func TestResumeConfirmingTask(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()

	// reject a run with an invalid state
	run := models.JobRun{}
	err := services.ResumeConfirmingTask(&run, store, big.NewInt(0))
	assert.Error(t, err)

	// reject a run with no tasks
	run = models.JobRun{Status: models.RunStatusPendingConfirmations}
	err = services.ResumeConfirmingTask(&run, store, big.NewInt(0))
	assert.Error(t, err)

	// leave in pending if not enough confirmations have been met yet
	creationHeight := cltest.BigHexInt(0)
	run = models.JobRun{
		ID:             utils.NewBytes32ID(),
		CreationHeight: &creationHeight,
		Status:         models.RunStatusPendingConfirmations,
		TaskRuns:       []models.TaskRun{models.TaskRun{MinimumConfirmations: 2, Task: models.TaskSpec{Type: adapters.TaskTypeNoOp}}},
	}
	err = services.ResumeConfirmingTask(&run, store, big.NewInt(0))
	assert.NoError(t, err)
	assert.Equal(t, string(models.RunStatusPendingConfirmations), string(run.Status))

	// input, should go from pending -> in progress and save the input
	creationHeight = cltest.BigHexInt(0)
	run = models.JobRun{
		ID:             utils.NewBytes32ID(),
		CreationHeight: &creationHeight,
		Status:         models.RunStatusPendingConfirmations,
		TaskRuns:       []models.TaskRun{models.TaskRun{MinimumConfirmations: 1, Task: models.TaskSpec{Type: adapters.TaskTypeNoOp}}},
	}
	err = services.ResumeConfirmingTask(&run, store, big.NewInt(1))
	assert.NoError(t, err)
	assert.Equal(t, string(models.RunStatusInProgress), string(run.Status))
}

func sleepAdapterParams(n int) models.JSON {
	d := time.Duration(n)
	json := []byte(fmt.Sprintf(`{"until":%v}`, time.Now().Add(d*time.Second).Unix()))
	return cltest.ParseJSON(bytes.NewBuffer(json))
}

func TestQueueSleepingTask(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()
	store.Clock = cltest.NeverClock{}

	// reject a run with an invalid state
	run := models.JobRun{}
	err := services.QueueSleepingTask(&run, store)
	assert.Error(t, err)

	// reject a run with no tasks
	run = models.JobRun{Status: models.RunStatusPendingSleep}
	err = services.QueueSleepingTask(&run, store)
	assert.Error(t, err)

	// reject a run that is sleeping but its task is not
	run = models.JobRun{
		ID:       utils.NewBytes32ID(),
		Status:   models.RunStatusPendingSleep,
		TaskRuns: []models.TaskRun{models.TaskRun{Task: models.TaskSpec{Type: adapters.TaskTypeSleep}}},
	}
	err = services.QueueSleepingTask(&run, store)
	assert.Error(t, err)

	// error decoding params into adapter
	inputFromTheFuture := cltest.ParseJSON(bytes.NewBuffer([]byte(`{"until": -1}`)))
	run = models.JobRun{
		ID:     utils.NewBytes32ID(),
		Status: models.RunStatusPendingSleep,
		TaskRuns: []models.TaskRun{
			models.TaskRun{
				Status: models.RunStatusPendingSleep,
				Task: models.TaskSpec{
					Type:   adapters.TaskTypeSleep,
					Params: inputFromTheFuture,
				},
			},
		},
	}
	err = services.QueueSleepingTask(&run, store)
	assert.NoError(t, err)
	assert.Equal(t, string(models.RunStatusErrored), string(run.TaskRuns[0].Status))
	assert.Equal(t, string(models.RunStatusErrored), string(run.Status))

	// mark run as pending, task as completed if duration has already elapsed
	run = models.JobRun{
		ID:       utils.NewBytes32ID(),
		Status:   models.RunStatusPendingSleep,
		TaskRuns: []models.TaskRun{models.TaskRun{Status: models.RunStatusPendingSleep, Task: models.TaskSpec{Type: adapters.TaskTypeSleep}}},
	}
	err = services.QueueSleepingTask(&run, store)
	assert.NoError(t, err)
	assert.Equal(t, string(models.RunStatusCompleted), string(run.TaskRuns[0].Status))
	assert.Equal(t, string(models.RunStatusInProgress), string(run.Status))

	runRequest, open := <-store.RunChannel.Receive()
	assert.True(t, open)
	assert.Equal(t, run.ID, runRequest.ID)

	// queue up next run if duration has not elapsed yet
	clock := cltest.UseSettableClock(store)
	store.Clock = clock
	clock.SetTime(time.Time{})

	inputFromTheFuture = sleepAdapterParams(60)
	run = models.JobRun{
		ID:     utils.NewBytes32ID(),
		Status: models.RunStatusPendingSleep,
		TaskRuns: []models.TaskRun{
			models.TaskRun{
				Status: models.RunStatusPendingSleep,
				Task: models.TaskSpec{
					Type:   adapters.TaskTypeSleep,
					Params: inputFromTheFuture,
				},
			},
		},
	}
	err = services.QueueSleepingTask(&run, store)
	assert.NoError(t, err)
	assert.Equal(t, string(models.RunStatusPendingSleep), string(run.TaskRuns[0].Status))
	assert.Equal(t, string(models.RunStatusPendingSleep), string(run.Status))

	// force the duration elapse
	clock.SetTime((time.Time{}).Add(math.MaxInt64))
	runRequest, open = <-store.RunChannel.Receive()
	assert.True(t, open)
	assert.Equal(t, run.ID, runRequest.ID)

	run, err = store.ORM.FindJobRun(run.ID)
	assert.NoError(t, err)
	assert.Equal(t, string(models.RunStatusCompleted), string(run.TaskRuns[0].Status))
	assert.Equal(t, string(models.RunStatusInProgress), string(run.Status))
}
