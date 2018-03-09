package models_test

import (
	"encoding/json"
	"fmt"
	"testing"

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

	job := models.NewJob()
	jr := job.NewRun()
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

	j := models.NewJob()
	j.Tasks = []models.TaskSpec{
		{Type: "NoOp"},
		{Type: "NoOpPend"},
		{Type: "NoOp"},
	}
	assert.Nil(t, store.SaveJob(&j))
	jr := j.NewRun()
	assert.Equal(t, jr.TaskRuns, jr.UnfinishedTaskRuns())

	jr, err := services.ExecuteRun(jr, store, models.RunResult{})
	assert.Nil(t, err)
	assert.Equal(t, jr.TaskRuns[1:], jr.UnfinishedTaskRuns())
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
					Params: models.JSON{gjson.Parse(orig)},
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
		{"string", `{"value": "100", "other": "101"}`, "100", false},
		{"integer", `{"value": 100}`, "", true},
		{"float", `{"value": 100.01}`, "", true},
		{"boolean", `{"value": true}`, "", true},
		{"null", `{"value": null}`, "", true},
		{"no key", `{"other": 100}`, "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var data models.JSON
			assert.Nil(t, json.Unmarshal([]byte(test.json), &data))
			rr := models.RunResult{Data: data}

			val, err := rr.Value()
			assert.Equal(t, test.want, val)
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func TestRunResult_Merge(t *testing.T) {
	t.Parallel()

	nullString := cltest.NullString(nil)
	jrID := utils.NewBytes32ID()
	tests := []struct {
		name             string
		originalData     string
		originalError    null.String
		originalPending  bool
		originalJRID     string
		inData           string
		inError          null.String
		inPending        bool
		inJRID           string
		wantData         string
		wantErrorMessage null.String
		wantPending      bool
		wantJRID         string
		wantErrored      bool
	}{
		{"merging data",
			`{"value":"old&busted","unique":"1"}`, nullString, false, jrID,
			`{"value":"newHotness","and":"!"}`, nullString, false, jrID,
			`{"value":"newHotness","unique":"1","and":"!"}`, nullString, false, jrID, false},
		{"original error throws",
			`{"value":"old"}`, cltest.NullString("old problem"), false, jrID,
			`{}`, nullString, false, jrID,
			`{"value":"old"}`, cltest.NullString("old problem"), false, jrID, true},
		{"error override",
			`{"value":"old"}`, nullString, false, jrID,
			`{}`, cltest.NullString("new problem"), false, jrID,
			`{"value":"old"}`, cltest.NullString("new problem"), false, jrID, false},
		{"original job run ID",
			`{"value":"old"}`, nullString, false, jrID,
			`{}`, nullString, false, "",
			`{"value":"old"}`, nullString, false, jrID, false},
		{"job run ID override",
			`{"value":"old"}`, nullString, false, utils.NewBytes32ID(),
			`{}`, nullString, false, jrID,
			`{"value":"old"}`, nullString, false, jrID, false},
		{"original pending",
			`{"value":"old"}`, nullString, true, jrID,
			`{}`, nullString, false, jrID,
			`{"value":"old"}`, nullString, true, jrID, false},
		{"pending override",
			`{"value":"old"}`, nullString, false, jrID,
			`{}`, nullString, true, jrID,
			`{"value":"old"}`, nullString, true, jrID, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := models.RunResult{
				Data:         models.JSON{gjson.Parse(test.originalData)},
				ErrorMessage: test.originalError,
				JobRunID:     test.originalJRID,
				Pending:      test.originalPending,
			}
			in := models.RunResult{
				Data:         cltest.JSONFromString(test.inData),
				ErrorMessage: test.inError,
				JobRunID:     test.inJRID,
				Pending:      test.inPending,
			}
			merged, err := original.Merge(in)
			assert.Equal(t, test.wantErrored, err != nil)

			assert.JSONEq(t, test.originalData, original.Data.String())
			assert.Equal(t, test.originalError, original.ErrorMessage)
			assert.Equal(t, test.originalJRID, original.JobRunID)
			assert.Equal(t, test.originalPending, original.Pending)

			assert.JSONEq(t, test.inData, in.Data.String())
			assert.Equal(t, test.inError, in.ErrorMessage)
			assert.Equal(t, test.inJRID, in.JobRunID)
			assert.Equal(t, test.inPending, in.Pending)

			assert.JSONEq(t, test.wantData, merged.Data.String())
			assert.Equal(t, test.wantErrorMessage, merged.ErrorMessage)
			assert.Equal(t, test.wantJRID, merged.JobRunID)
			assert.Equal(t, test.wantPending, merged.Pending)
		})
	}
}
