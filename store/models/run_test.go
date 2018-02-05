package models_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestRetrievingJobRunsWithErrorsFromDB(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	job := models.NewJob()
	jr := job.NewRun()
	jr.Result = models.RunResultWithError(fmt.Errorf("bad idea"))
	err := store.Save(jr)
	assert.Nil(t, err)

	run := &models.JobRun{}
	err = store.One("ID", jr.ID, run)
	assert.Nil(t, err)
	assert.True(t, run.Result.HasError())
	assert.Equal(t, "bad idea", run.Result.Error())
}

func TestTaskRunsToRun(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()

	j := models.NewJob()
	j.Tasks = []models.Task{
		{Type: "NoOp"},
		{Type: "NoOpPend"},
		{Type: "NoOp"},
	}
	assert.Nil(t, store.SaveJob(j))
	jr := j.NewRun()
	assert.Equal(t, jr.TaskRuns, jr.UnfinishedTaskRuns())

	err := services.ExecuteRun(jr, store, models.Output{})
	assert.Nil(t, err)
	assert.Equal(t, jr.TaskRuns[1:], jr.UnfinishedTaskRuns())
}

func TestOutputUnmarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		json        string
		wantErrored bool
	}{
		{"basic", `{"number": 100, "string": "100", "bool": true}`, false},
		{"invalid JSON", `{`, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var o models.Output
			err := json.Unmarshal([]byte(test.json), &o)
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func TestRunResultValue(t *testing.T) {
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
		{"no JSON", ``, "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output models.Output
			json.Unmarshal([]byte(test.json), &output)
			rr := models.RunResult{Output: &output}

			val, err := rr.Value()
			assert.Equal(t, test.want, val)
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func TestOutputMerge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       string
		want        string
		wantErrored bool
	}{
		{"new field", `{"extra":"fields"}`,
			`{"value":"OLD","other":1,"extra":"fields"}`, false},
		{"overwritting fields", `{"value":["new","new"],"extra":2}`,
			`{"value":["new","new"],"other":1,"extra":2}`, false},
		{"nested JSON", `{"extra":{"fields": ["more", 1]}}`,
			`{"value":"OLD","other":1,"extra":{"fields":["more",1]}}`, false},
		{"empty JSON", `{}`,
			`{"value":"OLD","other":1}`, false},
		{"null values", `{"value":null}`,
			`{"value":null,"other":1}`, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			o1 := cltest.OutputFromString(`{"value":"OLD","other":1}`)
			o2 := cltest.OutputFromString(test.input)

			err := o1.Merge(o2)
			assert.Equal(t, test.wantErrored, (err != nil))
			assert.JSONEq(t, test.want, o1.String())
		})
	}
}

func TestTaskRunMerge(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       string
		want        string
		wantErrored bool
	}{
		{"replace field", `{"url":"https://NEW.example.com/api"}`,
			`{"url":"https://NEW.example.com/api"}`, false},
		{"replace and add field", `{"url":"https://NEW.example.com/api","extra":1}`,
			`{"url":"https://NEW.example.com/api","extra":1}`, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tr := models.TaskRun{
				Task: models.Task{
					Params: json.RawMessage(`{"url":"https://OLD.example.com/api"}`),
					Type:   "httpget",
				},
			}
			input := cltest.OutputFromString(test.input)

			err := tr.Merge(*input)
			assert.Equal(t, test.wantErrored, (err != nil))
			assert.JSONEq(t, test.want, string(tr.Task.Params))
		})
	}
}
