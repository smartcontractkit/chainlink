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
		{"with endAt as ISO-8601",
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2006-01-02T15:04:05.000Z"}}`,
			`{"initiators":[{"type":"web"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2006-01-02T15:04:05.000Z"}`},
		{"with endAt as unix timestamp",
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"1522099336"}}`,
			`{"initiators":[{"type":"web"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2018-03-26T21:22:16.000Z"}`},
		{"with runAt as ISO-8601",
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2222-01-02T15:04:05.000Z","runAt":["2016-01-02T15:04:05.000Z","2026-01-02T15:04:05.000Z"]}}`,
			`{"initiators":[{"type":"web"},{"type":"runAt","time":"2016-01-02T15:04:05.000Z"},{"type":"runAt","time":"2026-01-02T15:04:05.000Z"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2222-01-02T15:04:05.000Z"}`},
		{"with runAt as unix timestamp",
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2222-01-02T15:04:05.000Z","runAt":["1522099336","1522109336"]}}`,
			`{"initiators":[{"type":"web"},{"type":"runAt","time":"2018-03-26T21:22:16.000Z"},{"type":"runAt","time":"2018-03-27T00:08:56.000Z"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2222-01-02T15:04:05.000Z"}`},
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
				wantInitiator.ID = actual.ID
				assert.Equal(t, wantInitiator, actual)
			}
		})
	}
}

func TestAssignmentSpec_ConvertToAssignment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{

		{"with endAt as ISO-8601",
			`{"initiators":[{"type":"web"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2006-01-02T15:04:05.000Z"}`,
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2006-01-02T15:04:05.000Z"}}`},
		{"with endAt as unix timestamp",
			`{"initiators":[{"type":"web"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2018-03-26T21:22:16.000Z"}`,
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"1522099336"}}`},
		{"with runAt as ISO-8601",
			`{"initiators":[{"type":"web"},{"type":"runAt","time":"2016-01-02T15:04:05.000Z"},{"type":"runAt","time":"2026-01-02T15:04:05.000Z"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2222-01-02T15:04:05.000Z"}`,
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2222-01-02T15:04:05.000Z","runAt":["2016-01-02T15:04:05.000Z","2026-01-02T15:04:05.000Z"]}}`},
		{"with runAt as unix timestamp",
			`{"initiators":[{"type":"web"},{"type":"runAt","time":"2018-03-26T21:22:16.000Z"},{"type":"runAt","time":"2018-03-27T00:08:56.000Z"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2222-01-02T15:04:05.000Z"}`,
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2222-01-02T15:04:05.000Z","runAt":["1522099336","1522109336"]}}`},
		{"with cron minute",
			`{"initiators":[{"type":"web"},{"type":"cron","schedule":"0 1 * * * *"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2006-01-02T15:04:05.000Z"}`,
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2006-01-02T15:04:05.000Z","minute":"1"}}`},
		{"with cron hour",
			`{"initiators":[{"type":"web"},{"type":"cron","schedule":"0 * 2 * * *"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2006-01-02T15:04:05.000Z"}`,
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2006-01-02T15:04:05.000Z","hour":"2"}}`},
		{"with cron day of month",
			`{"initiators":[{"type":"web"},{"type":"cron","schedule":"0 * * 3 * *"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2006-01-02T15:04:05.000Z"}`,
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2006-01-02T15:04:05.000Z","dayOfMonth":"3"}}`},
		{"with cron month of year",
			`{"initiators":[{"type":"web"},{"type":"cron","schedule":"0 * * * 4 *"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2006-01-02T15:04:05.000Z"}`,
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2006-01-02T15:04:05.000Z","monthOfYear":"4"}}`},
		{"with cron day of week",
			`{"initiators":[{"type":"web"},{"type":"cron","schedule":"0 * * * * 5"}],"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2006-01-02T15:04:05.000Z"}`,
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2006-01-02T15:04:05.000Z","dayOfWeek":"5"}}`,
		},
	}

	_, cleanup := cltest.NewStore()
	defer cleanup()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var js1 models.JobSpec
			assert.Nil(t, json.Unmarshal([]byte(test.input), &js1))

			a1, err := models.ConvertToAssignment(js1)
			assert.Nil(t, err)

			a2 := models.AssignmentSpec{}
			assert.Nil(t, json.Unmarshal([]byte(test.want), &a2))

			for i, wantTask := range a2.Assignment.Subtasks {
				actualTask := a1.Assignment.Subtasks[i]
				assert.Equal(t, strings.ToLower(wantTask.Type), actualTask.Type)
				assert.JSONEq(t, strings.ToLower(wantTask.Params.String()), actualTask.Params.String())
			}

			for i, v := range a1.Schedule.RunAt {
				assert.Equal(t, a2.Schedule.RunAt[i], v)
			}

			assert.Equal(t, a2.Schedule.Minute, a1.Schedule.Minute)
			assert.Equal(t, a2.Schedule.Hour, a1.Schedule.Hour)
			assert.Equal(t, a2.Schedule.DayOfMonth, a1.Schedule.DayOfMonth)
			assert.Equal(t, a2.Schedule.MonthOfYear, a1.Schedule.MonthOfYear)
			assert.Equal(t, a2.Schedule.DayOfWeek, a1.Schedule.DayOfWeek)
			assert.Equal(t, a2.Schedule.EndAt, a1.Schedule.EndAt)
		})
	}
}
