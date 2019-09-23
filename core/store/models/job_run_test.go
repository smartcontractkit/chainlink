package models_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
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
	jr.Result = cltest.RunResultWithError(fmt.Errorf("bad idea"))
	err := store.CreateJobRun(&jr)
	require.NoError(t, err)

	run, err := store.FindJobRun(jr.ID)
	assert.NoError(t, err)
	assert.True(t, run.Result.HasError())
	assert.Equal(t, "bad idea", run.Result.Error())
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
	jr.Result = cltest.RunResultWithData(data)
	err = store.CreateJobRun(&jr)
	assert.NoError(t, err)

	run, err := store.FindJobRun(jr.ID)
	assert.NoError(t, err)
	assert.False(t, run.Result.HasError())
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
	jr.Result = cltest.RunResultWithData(data)
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

	jr.Result = cltest.RunResultWithData(`{"result":"11850.00"}`)
	jr.Payment = linkReward
	logsBeforeCompletion := jr.ForLogger()
	require.Len(t, logsBeforeCompletion, 6)
	assert.Equal(t, logsBeforeCompletion[0], "job")
	assert.Equal(t, logsBeforeCompletion[1], jr.JobSpecID)
	assert.Equal(t, logsBeforeCompletion[2], "run")
	assert.Equal(t, logsBeforeCompletion[3], jr.ID)
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

	jrErr := job.NewRun(job.Initiators[0])
	jrErr.Result = cltest.RunResultWithError(fmt.Errorf("bad idea"))
	logsWithErr := jrErr.ForLogger()
	assert.Equal(t, logsWithErr[6], "job_error")
	assert.Equal(t, logsWithErr[7], jrErr.Result.Error())

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

func TestRunResult_Value(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		json        string
		want        string
		wantErrored bool
	}{
		{"string", `{"result": "100", "other": "101"}`, "100", false},
		{"integer", `{"result": 100}`, "", true},
		{"float", `{"result": 100.01}`, "", true},
		{"boolean", `{"result": true}`, "", true},
		{"null", `{"result": null}`, "", true},
		{"no key", `{"other": 100}`, "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var data models.JSON
			assert.NoError(t, json.Unmarshal([]byte(test.json), &data))
			rr := models.RunResult{Data: data}

			val, err := rr.ResultString()
			assert.Equal(t, test.want, val)
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func TestRunResult_Add(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		json  string
		key   string
		value interface{}
		want  string
	}{
		{"string", `{"a": "1"}`, "b", "2", `{"a": "1", "b": "2"}`},
		{"int", `{"a": "1"}`, "b", 2, `{"a": "1", "b": 2}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var data models.JSON
			assert.NoError(t, json.Unmarshal([]byte(test.json), &data))
			rr := models.RunResult{Data: data}

			rr.Add(test.key, test.value)

			assert.JSONEq(t, test.want, rr.Data.String())
		})
	}
}

func TestRunResult_WithError(t *testing.T) {
	t.Parallel()

	rr := models.RunResult{}

	assert.Equal(t, models.RunStatusUnstarted, rr.Status)

	rr.SetError(errors.New("this blew up"))

	assert.Equal(t, models.RunStatusErrored, rr.Status)
	assert.Equal(t, cltest.NullString("this blew up"), rr.ErrorMessage)
}

func TestRunResult_Merge(t *testing.T) {
	t.Parallel()

	inProgress := models.RunStatusInProgress
	pending := models.RunStatusPendingBridge
	errored := models.RunStatusErrored
	completed := models.RunStatusCompleted

	nullString := cltest.NullString(nil)
	tests := []struct {
		name             string
		originalData     string
		originalError    null.String
		originalStatus   models.RunStatus
		inData           string
		inError          null.String
		inStatus         models.RunStatus
		wantData         string
		wantErrorMessage null.String
		wantStatus       models.RunStatus
	}{
		{"merging data",
			`{"result":"old&busted","unique":"1"}`, nullString, inProgress,
			`{"result":"newHotness","and":"!"}`, nullString, inProgress,
			`{"result":"newHotness","unique":"1","and":"!"}`, nullString, inProgress},
		{"completed result",
			`{"result":"old"}`, nullString, inProgress,
			`{}`, nullString, completed,
			`{"result":"old"}`, nullString, completed},
		{"error override",
			`{"result":"old"}`, nullString, inProgress,
			`{}`, cltest.NullString("new problem"), errored,
			`{"result":"old"}`, cltest.NullString("new problem"), errored},
		{"pending override",
			`{"result":"old"}`, nullString, inProgress,
			`{}`, nullString, pending,
			`{"result":"old"}`, nullString, pending},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := models.RunResult{
				Data:         models.JSON{Result: gjson.Parse(test.originalData)},
				ErrorMessage: test.originalError,
				Status:       test.originalStatus,
			}
			in := models.RunResult{
				Data:         cltest.JSONFromString(t, test.inData),
				ErrorMessage: test.inError,
				Status:       test.inStatus,
			}
			merged := original
			merged.Merge(in)

			assert.JSONEq(t, test.originalData, original.Data.String())
			assert.Equal(t, test.originalError, original.ErrorMessage)
			assert.Equal(t, test.originalStatus, original.Status)

			assert.JSONEq(t, test.inData, in.Data.String())
			assert.Equal(t, test.inError, in.ErrorMessage)
			assert.Equal(t, test.inStatus, in.Status)

			assert.JSONEq(t, test.wantData, merged.Data.String())
			assert.Equal(t, test.wantErrorMessage, merged.ErrorMessage)
			assert.Equal(t, test.wantStatus, merged.Status)
		})
	}
}
