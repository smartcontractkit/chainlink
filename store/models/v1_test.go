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
			`{"tasks":[{"type":"noOp","foo":"bar"}]}`},
		{"withEndAt",
			`{"assignment":{"subtasks":[{"adapterType":"noOp","adapterParams":{"foo":"bar"}}]},"schedule":{"endAt":"2006-01-02T15:04:05.000Z"}}`,
			`{"tasks":[{"type":"noOp","foo":"bar"}],"endAt":"2006-01-02T15:04:05.000Z"}`},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var a models.AssignmentSpec
			assert.Nil(t, json.Unmarshal([]byte(test.input), &a))

			j1, err := a.ConvertToJobSpec()
			assert.Nil(t, err)
			assert.Nil(t, store.Save(&j1))
			j2 := cltest.FindJob(store, j1.ID)

			assert.NotEqual(t, "", j2.ID)
			var want models.JobSpec
			assert.Nil(t, json.Unmarshal([]byte(test.want), &want))
			assert.Equal(t, len(want.Tasks), len(j2.Tasks))
			assert.Equal(t, want.EndAt, j2.EndAt)
			for i, wantTask := range want.Tasks {
				actual := j2.Tasks[i]
				assert.Equal(t, strings.ToLower(wantTask.Type), actual.Type)
				assert.JSONEq(t, wantTask.Params.String(), actual.Params.String())
			}
		})
	}
}
