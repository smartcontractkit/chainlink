package services_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestValidateJob(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input []byte
		want  error
	}{
		{"base case", cltest.LoadJSON("../internal/fixtures/web/hello_world_job.json"), nil},
		{"error in runat initr", cltest.LoadJSON("../internal/fixtures/web/run_at_wo_time_job.json"),
			errors.New(`job validation: initiator validation: runat must have time`)},
		{"error in task", cltest.LoadJSON("../internal/fixtures/web/nonexistent_task_job.json"),
			errors.New(`job validation: task validation: idonotexist is not a supported adapter type`)},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var j models.Job
			assert.Nil(t, json.Unmarshal(test.input, &j))
			result := services.ValidateJob(j, store)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestValidateInitiator(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"web", `{"type":"web"}`, false},
		{"ethlog", `{"type":"ethlog"}`, false},
		{"runlog", `{"type":"runlog"}`, false},
		{"runat w/o time", `{"type":"runat"}`, true},
		{"cron w/o schedule", `{"type":"cron"}`, true},
		{"non-existent initiator", `{"type":"doesntExist"}`, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var initr models.Initiator
			assert.Nil(t, json.Unmarshal([]byte(test.input), &initr))
			result := services.ValidateInitiator(initr)
			assert.Equal(t, test.wantError, result != nil)
		})
	}
}
