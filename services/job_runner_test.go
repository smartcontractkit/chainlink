package services_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

func TestJobRunner_resumeSleepingRuns(t *testing.T) {
	store, cleanup := cltest.NewStore()
	defer cleanup()
	rm, cleanup := cltest.NewJobRunner(store)
	defer cleanup()

	j := models.NewJob()
	i := models.Initiator{Type: models.InitiatorWeb}
	j.Initiators = []models.Initiator{i}
	json := fmt.Sprintf(`{"until":"%v"}`, utils.ISO8601UTC(time.Now().Add(time.Second)))
	j.Tasks = []models.TaskSpec{cltest.NewTask("sleep", json)}
	assert.NoError(t, store.Save(&j))

	jr := j.NewRun(i)
	jr.Status = models.RunStatusPendingSleep
	jr.Result.Data = cltest.JSONFromString(`{"foo":"bar"}`)
	assert.NoError(t, store.Save(&jr))

	assert.NoError(t, services.ExportedResumeSleepingRuns(rm))
	rr, open := <-store.RunChannel.Receive()
	assert.Equal(t, jr.Result, rr.Input)
	assert.True(t, open)
}

func TestJobRunner_ChannelForRun_equalityBetweenRuns(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
	defer cleanup()
	rm, cleanup := cltest.NewJobRunner(store)
	defer cleanup()

	job, initr := cltest.NewJobWithWebInitiator()
	run1 := job.NewRun(initr)
	run2 := job.NewRun(initr)

	chan1a := services.ExportedChannelForRun(rm, run1.ID)
	chan2 := services.ExportedChannelForRun(rm, run2.ID)
	chan1b := services.ExportedChannelForRun(rm, run1.ID)

	assert.NotEqual(t, chan1a, chan2)
	assert.Equal(t, chan1a, chan1b)
	assert.NotEqual(t, chan2, chan1b)
}

func TestJobRunner_ChannelForRun_sendAfterClosing(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore()
	defer cleanup()
	rm, cleanup := cltest.NewJobRunner(s)
	defer cleanup()
	assert.NoError(t, rm.Start())

	j, initr := cltest.NewJobWithWebInitiator()
	assert.NoError(t, s.SaveJob(&j))
	jr := j.NewRun(initr)
	assert.NoError(t, s.Save(&jr))

	chan1 := services.ExportedChannelForRun(rm, jr.ID)
	chan1 <- store.RunRequest{Input: jr.Result}
	cltest.WaitForJobRunToComplete(t, s, jr)

	gomega.NewGomegaWithT(t).Eventually(func() chan<- store.RunRequest {
		return services.ExportedChannelForRun(rm, jr.ID)
	}).Should(gomega.Not(gomega.Equal(chan1))) // eventually deletes the channel

	chan2 := services.ExportedChannelForRun(rm, jr.ID)
	chan2 <- store.RunRequest{Input: jr.Result} // does not panic
}

func TestJobRunner_ChannelForRun_equalityWithoutClosing(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore()
	defer cleanup()
	rm, cleanup := cltest.NewJobRunner(s)
	defer cleanup()
	assert.NoError(t, rm.Start())

	j, initr := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{cltest.NewTask("nooppend")}
	assert.NoError(t, s.SaveJob(&j))
	jr := j.NewRun(initr)
	assert.NoError(t, s.Save(&jr))

	chan1 := services.ExportedChannelForRun(rm, jr.ID)

	chan1 <- store.RunRequest{}
	cltest.WaitForJobRunToPendConfirmations(t, s, jr)

	chan2 := services.ExportedChannelForRun(rm, jr.ID)
	assert.Equal(t, chan1, chan2)
}

func TestJobRunner_Stop(t *testing.T) {
	t.Parallel()

	s, cleanup := cltest.NewStore()
	defer cleanup()
	rm, cleanup := cltest.NewJobRunner(s)
	defer cleanup()
	j, initr := cltest.NewJobWithWebInitiator()
	jr := j.NewRun(initr)

	services.ExportedChannelForRun(rm, jr.ID)
	assert.Equal(t, 1, services.ExportedWorkerCount(rm))

	rm.Stop()

	gomega.NewGomegaWithT(t).Eventually(func() int {
		return services.ExportedWorkerCount(rm)
	}).Should(gomega.Equal(0))
}

func TestJobRunner_ExecuteRun(t *testing.T) {
	t.Parallel()

	bridgeName := "auctionBidding"
	tests := []struct {
		name       string
		bridgeType string
		input      string
		runResult  string
		wantStatus models.RunStatus
		wantData   string
	}{
		{"success", bridgeName, `{}`, `{"data":{"value":"100"}}`, models.RunStatusCompleted, `{"value":"100"}`},
		{"errored", bridgeName, `{}`, `{"error":"too much"}`, models.RunStatusErrored, `{}`},
		{"errored with a value", bridgeName, `{}`, `{"error":"too much", "data":{"value":"99"}}`, models.RunStatusErrored, `{"value":"99"}`},
		{"overriding bridge type params", bridgeName, `{"url":"hack"}`, `{"data":{"value":"100"}}`, models.RunStatusCompleted, `{"value":"100","url":"hack"}`},
		{"type parameter does not override", bridgeName, `{"type":"0"}`, `{"data":{"value":"100"}}`, models.RunStatusCompleted, `{"value":"100","type":"0"}`},
		{"non-existent bridge type", "non-existent", `{}`, `{}`, models.RunStatusErrored, `{}`},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {

			var run models.JobRun
			mockServer, _ := cltest.NewHTTPMockServer(t, 200, "POST", test.runResult,
				func(body string) {
					want := fmt.Sprintf(`{"id":"%v","data":%v}`, run.ID, test.input)
					assert.JSONEq(t, want, body)
				})
			bt := cltest.NewBridgeType(bridgeName, mockServer.URL)
			assert.Nil(t, store.Save(&bt))

			job, initr := cltest.NewJobWithWebInitiator()
			job.Tasks = []models.TaskSpec{
				cltest.NewTask(test.bridgeType),
				cltest.NewTask("noop"),
			}
			assert.Nil(t, store.Save(&job))

			run = job.NewRun(initr)
			input := models.RunResult{Data: cltest.JSONFromString(test.input)}
			run, err := services.ExecuteRun(run, store, input)
			assert.NoError(t, err)

			store.One("ID", run.ID, &run)
			assert.Equal(t, test.wantStatus, run.Status)
			assert.JSONEq(t, test.wantData, run.Result.Data.String())

			tr1 := run.TaskRuns[0]
			assert.Equal(t, test.wantStatus, tr1.Status)
			assert.Equal(t, test.wantStatus, tr1.Result.Status)
			assert.JSONEq(t, test.wantData, tr1.Result.Data.String())

			if test.wantStatus == models.RunStatusCompleted {
				tr2 := run.TaskRuns[1]
				assert.JSONEq(t, test.wantData, tr2.Result.Data.String())
				assert.True(t, run.CompletedAt.Valid)
			}
		})
	}
}

func TestExecuteRun_TransitionToPendingConfirmations(t *testing.T) {
	t.Parallel()

	config, cfgCleanup := cltest.NewConfig()
	defer cfgCleanup()
	config.MinIncomingConfirmations = 10

	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()
	creationHeight := 1000
	configMin := int(store.Config.MinIncomingConfirmations)

	tests := []struct {
		name           string
		confirmations  int
		triggeringConf int
	}{
		{"not defined in task spec", 0, configMin},
		{"task spec > global min confs", configMin + 1, configMin + 1},
		{"task spec == global min confs", configMin, configMin},
		{"task spec < global min confs", configMin - 1, configMin},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			job, initr := cltest.NewJobWithLogInitiator()
			job.Tasks = []models.TaskSpec{
				{
					Type:          adapters.TaskTypeNoOp,
					Confirmations: uint64(test.confirmations),
				},
			}

			run := job.NewRun(initr)
			run, err := store.SaveCreationHeight(run, cltest.IndexableBlockNumber(creationHeight))
			assert.NoError(t, err)

			early := cltest.IndexableBlockNumber(creationHeight + test.triggeringConf - 2)
			initialData := models.JSON{Result: gjson.Parse(`{"address":"0xdfcfc2b9200dbb10952c2b7cce60fc7260e03c6f"}`)}
			runLogInitialInput := models.RunResult{Data: initialData}
			run, err = services.ExecuteRunAtBlock(run, store, runLogInitialInput, early)
			assert.NoError(t, err)

			store.One("ID", run.ID, &run)
			assert.Equal(t, models.RunStatusPendingConfirmations, run.Status)
			assert.Equal(t, initialData, run.Result.Data)

			trigger := cltest.IndexableBlockNumber(creationHeight + test.triggeringConf - 1)
			run, err = services.ExecuteRunAtBlock(run, store, models.RunResult{}, trigger)
			assert.NoError(t, err)
			assert.Equal(t, models.RunStatusCompleted, run.Status)
			assert.Equal(t, initialData, run.Result.Data)
		})
	}
}

func TestExecuteRun_TransitionToPendingConfirmations_WithBridgeTask(t *testing.T) {
	t.Parallel()

	config, cfgCleanup := cltest.NewConfig()
	defer cfgCleanup()
	config.MinIncomingConfirmations = 10
	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()
	creationHeight := 1000
	configMin := int(store.Config.MinIncomingConfirmations)

	tests := []struct {
		name                    string
		bridgeTypeConfirmations int
		taskSpecConfirmations   int
		triggeringConf          int
	}{
		{"not defined in task spec or bridge type", 0, 0, configMin},
		{"bridge type confirmations > task spec confirmations", configMin + 1, configMin, configMin + 1},
		{"bridge type confirmations = task spec confirmations", configMin, configMin, configMin},
		{"bridge type confirmations < task spec confirmations", configMin - 2, configMin - 1, configMin},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			job, initr := cltest.NewJobWithLogInitiator()
			job.Tasks = []models.TaskSpec{
				{
					Type:          models.NewTaskType("randomNumber"),
					Confirmations: uint64(test.taskSpecConfirmations),
				},
			}

			run := job.NewRun(initr)
			mockServer, _ := cltest.NewHTTPMockServer(t, 200, "POST", "{\"todo\": \"todo\"}",
				func(body string) {
					want := fmt.Sprintf(`{"id":"%v","data":%v}`, run.ID, "{}")
					assert.JSONEq(t, want, body)
				})
			bt := cltest.NewBridgeTypeWithDefaultConfirmations(uint64(test.bridgeTypeConfirmations), "randomNumber", mockServer.URL)
			assert.Nil(t, store.Save(&bt))

			run, err := store.SaveCreationHeight(run, cltest.IndexableBlockNumber(creationHeight))
			assert.NoError(t, err)

			early := cltest.IndexableBlockNumber(creationHeight + test.triggeringConf - 2)
			run, err = services.ExecuteRunAtBlock(run, store, models.RunResult{}, early)
			assert.NoError(t, err)

			store.One("ID", run.ID, &run)
			assert.Equal(t, models.RunStatusPendingConfirmations, run.Status)

			trigger := cltest.IndexableBlockNumber(creationHeight + test.triggeringConf - 1)
			run, err = services.ExecuteRunAtBlock(run, store, models.RunResult{}, trigger)
			assert.NoError(t, err)
			assert.Equal(t, models.RunStatusCompleted, run.Status)
		})
	}
}

func TestJobRunner_ExecuteRun_TransitionToPending(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	job, initr := cltest.NewJobWithWebInitiator()
	job.Tasks = []models.TaskSpec{cltest.NewTask("NoOpPend")}

	run := job.NewRun(initr)
	run, err := services.ExecuteRun(run, store, models.RunResult{})
	assert.NoError(t, err)

	store.One("ID", run.ID, &run)
	assert.Equal(t, models.RunStatusPendingConfirmations, run.Status)
}

func TestJobRunner_ExecuteRun_ErrorsWithNoRuns(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	job, initr := cltest.NewJobWithWebInitiator()
	job.Tasks = []models.TaskSpec{}
	run := job.NewRun(initr)
	run, err := services.ExecuteRun(run, store, models.RunResult{})
	assert.Error(t, err)
}

func TestJobRunner_BeginRunWithAmount(t *testing.T) {
	tests := []struct {
		name   string
		amount *big.Int
		status models.RunStatus
	}{
		{"job with no amount", nil, models.RunStatusCompleted},
		{"job with insufficient amount", big.NewInt(9), models.RunStatusErrored},
		{"job with exact amount", big.NewInt(10), models.RunStatusCompleted},
		{"job with valid amount", big.NewInt(11), models.RunStatusCompleted},
	}

	config, cfgCleanup := cltest.NewConfig()
	defer cfgCleanup()
	config.MinimumContractPayment = *big.NewInt(10)

	store, cleanup := cltest.NewStoreWithConfig(config)
	defer cleanup()

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			job, initr := cltest.NewJobWithWebInitiator()
			assert.Nil(t, store.SaveJob(&job))

			runResult := models.RunResult{
				Amount: test.amount,
			}
			run, _ := services.BeginRun(job, initr, runResult, store)
			assert.Equal(t, test.status, run.Status)
		})
	}
}

func TestJobRunner_BeginRun(t *testing.T) {
	pastTime := cltest.ParseNullableTime("2000-01-01T00:00:00.000Z")
	futureTime := cltest.ParseNullableTime("3000-01-01T00:00:00.000Z")
	nullTime := null.Time{Valid: false}

	tests := []struct {
		name     string
		startAt  null.Time
		endAt    null.Time
		errored  bool
		runCount int
	}{
		{"job not started", futureTime, nullTime, true, 0},
		{"job started", pastTime, futureTime, false, 1},
		{"job with no time range", nullTime, nullTime, false, 1},
		{"job ended", nullTime, pastTime, true, 0},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			job, initr := cltest.NewJobWithWebInitiator()
			job.StartAt = test.startAt
			job.EndAt = test.endAt
			assert.Nil(t, store.SaveJob(&job))

			_, err := services.BeginRun(job, initr, models.RunResult{}, store)

			if test.errored {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			jrs, err := store.JobRunsFor(job.ID)
			assert.NoError(t, err)
			assert.Equal(t, test.runCount, len(jrs))
		})
	}
}

func TestJobRunner_BuildRun(t *testing.T) {
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

			_, err := services.BuildRun(job, initr, store)

			if test.errored {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
