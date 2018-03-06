package services_test

import (
	"encoding/json"
	"errors"
	"fmt"
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
		{"with error in initiator", cltest.LoadJSON("../internal/fixtures/web/run_at_wo_time_job.json"),
			errors.New(`job validation: initiator validation: runat must have time`)},
		{"with error in adapter", cltest.LoadJSON("../internal/fixtures/web/nonexistent_task_job.json"),
			errors.New(`job validation: task validation: idonotexist is not a supported adapter type`)},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var j models.Job
			assert.Nil(t, json.Unmarshal(test.input, &j))
			fmt.Println(j)
			result := services.ValidateJob(j, store)
			assert.Equal(t, test.want, result)
		})
	}
}
