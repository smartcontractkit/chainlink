package models_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

func TestJobRuns_RetrievingFromDBWithError(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	job, initr := cltest.NewJobWithWebInitiator()
	jr := job.NewRun(initr)
	jr.Result = cltest.RunResultWithError(fmt.Errorf("bad idea"))
	err := store.Save(&jr)
	assert.Nil(t, err)

	run := &models.JobRun{}
	err = store.One("ID", jr.ID, run)
	assert.Nil(t, err)
	assert.True(t, run.Result.HasError())
	assert.Equal(t, "bad idea", run.Result.Error())
}

func TestJobRun_UnfinishedTaskRuns(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j, i := cltest.NewJobWithWebInitiator()
	j.Tasks = []models.TaskSpec{
		{Type: "NoOp"},
		{Type: "NoOpPend"},
		{Type: "NoOp"},
	}
	assert.Nil(t, store.SaveJob(&j))
	jr := j.NewRun(i)
	assert.Equal(t, jr.TaskRuns, jr.UnfinishedTaskRuns())

	jr, err := services.ExecuteRun(jr, store, models.RunResult{})
	assert.Nil(t, err)
	assert.Equal(t, jr.TaskRuns[1:], jr.UnfinishedTaskRuns())
}

func TestTaskRun_Runnable(t *testing.T) {
	t.Parallel()

	job, initr := cltest.NewJobWithLogInitiator()
	tests := []struct {
		name                 string
		creationHeight       *hexutil.Big
		currentHeight        *models.IndexableBlockNumber
		minimumConfirmations uint64
		want                 bool
	}{
		{"unset nil 0", nil, nil, 0, true},
		{"1 nil 0", cltest.NewBigHexInt(1), nil, 0, true},
		{"1 1 minconf 0", cltest.NewBigHexInt(1), cltest.IndexableBlockNumber(1), 0, true},
		{"1 1 minconf 1", cltest.NewBigHexInt(1), cltest.IndexableBlockNumber(1), 1, false},
		{"1 2 minconf 1", cltest.NewBigHexInt(1), cltest.IndexableBlockNumber(2), 1, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jr := job.NewRun(initr)
			if test.creationHeight != nil {
				jr.CreationHeight = test.creationHeight
			}

			assert.Equal(t, test.want, jr.Runnable(test.currentHeight, test.minimumConfirmations))
		})
	}
}

func TestTaskRun_Merge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       string
		want        string
		wantErrored bool
	}{
		{"replace field", `{"url":"https://NEW.example.com/api"}`,
			`{"url":"https://NEW.example.com/api"}`, false},
		{"add field", `{"extra":1}`,
			`{"url":"https://OLD.example.com/api","extra":1}`, false},
		{"replace and add field", `{"url":"https://NEW.example.com/api","extra":1}`,
			`{"url":"https://NEW.example.com/api","extra":1}`, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			orig := `{"url":"https://OLD.example.com/api"}`
			tr := models.TaskRun{
				Task: models.TaskSpec{
					Params: models.JSON{Result: gjson.Parse(orig)},
					Type:   "httpget",
				},
			}
			input := cltest.JSONFromString(test.input)

			merged, err := tr.MergeTaskParams(input)
			assert.Equal(t, test.wantErrored, (err != nil))
			assert.JSONEq(t, test.want, merged.Task.Params.String())
			assert.JSONEq(t, orig, tr.Task.Params.String())
		})
	}
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
			assert.Nil(t, json.Unmarshal([]byte(test.json), &data))
			rr := models.RunResult{Data: data}

			val, err := rr.Result()
			assert.Equal(t, test.want, val)
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func TestRunResult_WithError(t *testing.T) {
	t.Parallel()

	rr := models.RunResult{}

	assert.Equal(t, models.RunStatusUnstarted, rr.Status)

	rr = rr.WithError(errors.New("this blew up"))

	assert.Equal(t, models.RunStatusErrored, rr.Status)
	assert.Equal(t, cltest.NullString("this blew up"), rr.ErrorMessage)
}

func TestRunResult_Merge(t *testing.T) {
	t.Parallel()

	inProgress := models.RunStatusInProgress
	pending := models.RunStatusPendingBridge
	errored := models.RunStatusErrored

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
				Data:         models.JSON{Result: gjson.Parse(test.originalData)},
				ErrorMessage: test.originalError,
				JobRunID:     test.originalJRID,
				Status:       test.originalStatus,
			}
			in := models.RunResult{
				Data:         cltest.JSONFromString(test.inData),
				ErrorMessage: test.inError,
				JobRunID:     test.inJRID,
				Status:       test.inStatus,
			}
			merged, err := original.Merge(in)
			assert.Equal(t, test.wantErrored, err != nil)

			assert.JSONEq(t, test.originalData, original.Data.String())
			assert.Equal(t, test.originalError, original.ErrorMessage)
			assert.Equal(t, test.originalJRID, original.JobRunID)
			assert.Equal(t, test.originalStatus, original.Status)

			assert.JSONEq(t, test.inData, in.Data.String())
			assert.Equal(t, test.inError, in.ErrorMessage)
			assert.Equal(t, test.inJRID, in.JobRunID)
			assert.Equal(t, test.inStatus, in.Status)

			assert.JSONEq(t, test.wantData, merged.Data.String())
			assert.Equal(t, test.wantErrorMessage, merged.ErrorMessage)
			assert.Equal(t, test.wantJRID, merged.JobRunID)
			assert.Equal(t, test.wantStatus, merged.Status)
		})
	}
}
