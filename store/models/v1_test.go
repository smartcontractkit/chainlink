package models_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestAssignmentSpec_ConvertToJobSpec(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"basic",
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]}}`,
			`{"initiators":[{"type":"web"}],"tasks":[{"type":"noOp","foo":"bar"}]}`},
		{"with endAt",
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2006-01-02T15:04:05.000Z"}}`,
			`{"initiators":[{"type":"web"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2006-01-02T15:04:05.000Z"}`},
		{"with runAt",
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2222-01-02T15:04:05.000Z","runAt":["2016-01-02T15:04:05.000Z","2026-01-02T15:04:05.000Z"]}}`,
			`{"initiators":[{"type":"web"},{"type":"runAt","time":"2016-01-02T15:04:05.000Z"},{"type":"runAt","time":"2026-01-02T15:04:05.000Z"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2222-01-02T15:04:05.000Z"}`},
		{"with cron minute",
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2006-01-02T15:04:05.000Z","minute":"1"}}`,
			`{"initiators":[{"type":"web"},{"type":"cron","schedule":"0 1 * * * *"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2006-01-02T15:04:05.000Z"}`},
		{"with cron hour",
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2006-01-02T15:04:05.000Z","hour":"2"}}`,
			`{"initiators":[{"type":"web"},{"type":"cron","schedule":"0 * 2 * * *"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2006-01-02T15:04:05.000Z"}`},
		{"with cron day of month",
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2006-01-02T15:04:05.000Z","dayOfMonth":"3"}}`,
			`{"initiators":[{"type":"web"},{"type":"cron","schedule":"0 * * 3 * *"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2006-01-02T15:04:05.000Z"}`},
		{"with cron month of year",
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2006-01-02T15:04:05.000Z","monthOfYear":"4"}}`,
			`{"initiators":[{"type":"web"},{"type":"cron","schedule":"0 * * * 4 *"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2006-01-02T15:04:05.000Z"}`},
		{"with cron day of week",
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2006-01-02T15:04:05.000Z","dayOfWeek":"5"}}`,
			`{"initiators":[{"type":"web"},{"type":"cron","schedule":"0 * * * * 5"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2006-01-02T15:04:05.000Z"}`},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var a models.AssignmentSpec
			assert.Nil(t, json.Unmarshal([]byte(test.input), &a))

			j1, err := a.ConvertToJobSpec()
			assert.Nil(t, err)
			assert.Nil(t, store.SaveJob(&j1))
			j2 := cltest.FindJob(store, j1.ID)

			assert.NotEqual(t, "", j2.ID)
			var want models.JobSpec
			assert.Nil(t, json.Unmarshal([]byte(test.want), &want))
			assert.Equal(t, want.EndAt, j2.EndAt)

			for i, wantTask := range want.Tasks {
				actual := j2.Tasks[i]
				assert.Equal(t, strings.ToLower(wantTask.Type), actual.Type)
				assert.JSONEq(t, wantTask.Params.String(), actual.Params.String())
			}

			for i, wantInitiator := range want.Initiators {
				actual := j2.Initiators[i]
				wantInitiator.JobID = actual.JobID
				assert.Equal(t, wantInitiator, actual)
			}
		})
	}
}
