package models_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
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
	j.Tasks = []models.Task{
		{Type: "NoOp"},
		{Type: "NoOpPend"},
		{Type: "NoOp"},
	}
	assert.Nil(t, store.SaveJob(&j))
	jr := j.NewRun()
	assert.Equal(t, jr.TaskRuns, jr.UnfinishedTaskRuns())

	jr, err := services.ExecuteRun(jr, store, models.JSON{})
	assert.Nil(t, err)
	assert.Equal(t, jr.TaskRuns[1:], jr.UnfinishedTaskRuns())
}

func TestTaskRun_MergeTaskParams(t *testing.T) {
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
				Task: models.Task{
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
		{"no JSON", ``, "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var data models.JSON
			json.Unmarshal([]byte(test.json), &data)
			rr := models.RunResult{Data: data}

			val, err := rr.Value()
			assert.Equal(t, test.want, val)
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func TestRunResult_MergeData(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       string
		want        string
		wantErrored bool
	}{
		{"replace field", `{"value":"new hotness"}`,
			`{"value":"new hotness"}`, false},
		{"add field", `{"extra":1}`,
			`{"value":"old and busted","extra":1}`, false},
		{"replace and add field", `{"value":"new hotness","extra":1}`,
			`{"value":"new hotness","extra":1}`, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			orig := `{"value":"old and busted"}`
			rr := models.RunResult{
				Data: models.JSON{gjson.Parse(orig)},
			}
			input := cltest.JSONFromString(test.input)

			merged, err := rr.MergeData(input)
			assert.Equal(t, test.wantErrored, (err != nil))
			assert.JSONEq(t, test.want, merged.Data.String())
			assert.JSONEq(t, orig, rr.Data.String())
		})
	}
}
