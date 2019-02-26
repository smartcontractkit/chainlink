package models_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

func TestJobRuns_RetrievingFromDBWithError(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
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
	store, cleanup := cltest.NewStore()
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

func TestJobRun_NextTaskRun(t *testing.T) {
	t.Parallel()

	store, cleanup := cltest.NewStore()
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

	rr.WithError(errors.New("this blew up"))

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
	jrID := utils.NewBytes32ID()
	tests := []struct {
		name             string
		originalData     string
		originalError    null.String
		originalStatus   models.RunStatus
		originalJRID     string
		inData           string
		inError          null.String
		inStatus         models.RunStatus
		inJRID           string
		wantData         string
		wantErrorMessage null.String
		wantStatus       models.RunStatus
		wantJRID         string
		wantErrored      bool
	}{
		{"merging data",
			`{"result":"old&busted","unique":"1"}`, nullString, inProgress, jrID,
			`{"result":"newHotness","and":"!"}`, nullString, inProgress, jrID,
			`{"result":"newHotness","unique":"1","and":"!"}`, nullString, inProgress, jrID, false},
		{"completed result",
			`{"result":"old"}`, nullString, inProgress, jrID,
			`{}`, nullString, completed, jrID,
			`{"result":"old"}`, nullString, completed, jrID, false},
		{"original error throws",
			`{"result":"old"}`, cltest.NullString("old problem"), errored, jrID,
			`{}`, nullString, inProgress, jrID,
			`{"result":"old"}`, cltest.NullString("old problem"), errored, jrID, true},
		{"error override",
			`{"result":"old"}`, nullString, inProgress, jrID,
			`{}`, cltest.NullString("new problem"), errored, jrID,
			`{"result":"old"}`, cltest.NullString("new problem"), errored, jrID, false},
		{"original job run ID",
			`{"result":"old"}`, nullString, inProgress, jrID,
			`{}`, nullString, inProgress, "",
			`{"result":"old"}`, nullString, inProgress, jrID, false},
		{"job run ID override",
			`{"result":"old"}`, nullString, inProgress, utils.NewBytes32ID(),
			`{}`, nullString, inProgress, jrID,
			`{"result":"old"}`, nullString, inProgress, jrID, false},
		{"original pending",
			`{"result":"old"}`, nullString, pending, jrID,
			`{}`, nullString, inProgress, jrID,
			`{"result":"old"}`, nullString, pending, jrID, false},
		{"pending override",
			`{"result":"old"}`, nullString, inProgress, jrID,
			`{}`, nullString, pending, jrID,
			`{"result":"old"}`, nullString, pending, jrID, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := models.RunResult{
				Data:           models.JSON{Result: gjson.Parse(test.originalData)},
				ErrorMessage:   test.originalError,
				CachedJobRunID: test.originalJRID,
				Status:         test.originalStatus,
			}
			in := models.RunResult{
				Data:           cltest.JSONFromString(t, test.inData),
				ErrorMessage:   test.inError,
				CachedJobRunID: test.inJRID,
				Status:         test.inStatus,
			}
			merged := original
			err := merged.Merge(in)
			if test.wantErrored {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.JSONEq(t, test.originalData, original.Data.String())
			assert.Equal(t, test.originalError, original.ErrorMessage)
			assert.Equal(t, test.originalJRID, original.CachedJobRunID)
			assert.Equal(t, test.originalStatus, original.Status)

			assert.JSONEq(t, test.inData, in.Data.String())
			assert.Equal(t, test.inError, in.ErrorMessage)
			assert.Equal(t, test.inJRID, in.CachedJobRunID)
			assert.Equal(t, test.inStatus, in.Status)

			assert.JSONEq(t, test.wantData, merged.Data.String())
			assert.Equal(t, test.wantErrorMessage, merged.ErrorMessage)
			assert.Equal(t, test.wantJRID, merged.CachedJobRunID)
			assert.Equal(t, test.wantStatus, merged.Status)
		})
	}
}
