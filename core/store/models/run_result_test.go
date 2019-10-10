package models_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	null "gopkg.in/guregu/null.v3"
)

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
