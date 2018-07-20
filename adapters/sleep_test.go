package adapters_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func sleepFor(n int) string {
	d := time.Duration(n)
	return fmt.Sprintf(`{"until":%v}`, time.Now().Add(d*time.Second).Unix())
}

func TestSleep_Perform(t *testing.T) {
	t.Parallel()
	standardInput := models.RunResult{}
	sleptInput := models.RunResult{
		Status: models.RunStatusPendingSleep,
	}

	tests := []struct {
		name         string
		params       string
		input        models.RunResult
		wantStatus   models.RunStatus
		parseErrored bool
		errored      bool
	}{
		{"valid duration", sleepFor(60), standardInput, models.RunStatusPendingSleep, false, false},
		{"past time", sleepFor(-1), standardInput, models.RunStatusCompleted, false, false},
		{"json with iso8601", `{"until":"2222-07-20T12:54:49.000Z"}`, standardInput, models.RunStatusPendingSleep, false, false},
		{"long duration", sleepFor(12592000), standardInput, models.RunStatusPendingSleep, false, false},
		{"invalid json", `{"until":"1000h"}`, standardInput, models.RunStatusPendingSleep, true, false},
		{"already slept", sleepFor(1000), sleptInput, models.RunStatusPendingSleep, false, false},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()
	store.Clock = cltest.InstantClock{}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			adapter := adapters.Sleep{}
			err := json.Unmarshal([]byte(test.params), &adapter)
			if test.parseErrored {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			result := adapter.Perform(test.input, store)
			assert.Equal(t, test.wantStatus, result.Status)
			if test.errored {
				assert.Error(t, result.GetError())
			} else {
				assert.NoError(t, result.GetError())
			}
		})
	}
}
