package models_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestAssignmentSpec_ConvertToJobSpec(t *testing.T) {
	t.Parallel()
	store, cleanup := cltest.NewStore()
	defer cleanup()
	json := cltest.JSONFromString(`{"foo": "bar"}`)

	v1 := models.AssignmentSpec{
		Assignment: models.Assignment{
			Subtasks: []models.Subtask{
				models.Subtask{
					Type:   "noOp",
					Params: json,
				},
			},
		},
	}

	j, err := v1.ConvertToJobSpec()
	assert.Nil(t, err)
	assert.Equal(t, "models.JobSpec", reflect.TypeOf(j).String())
	assert.Equal(t, 1, len(j.Tasks))
	task := j.Tasks[0]
	assert.Equal(t, "noOp", task.Type)
	assert.JSONEq(t, `{"foo": "bar", "type": "noOp"}`, task.Params.String())
	assert.NotEqual(t, "", j.ID)

	assert.Nil(t, store.Save(&j))
	j2 := cltest.FindJob(store, j.ID)
	assert.Equal(t, strings.ToLower(task.Type), j2.Tasks[0].Type)
}
