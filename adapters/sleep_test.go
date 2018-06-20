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

	tests := []struct {
		name         string
		params       string
		parseErrored bool
		errored      bool
	}{
		{"excessive duration", sleepFor(2592010), false, true},
		{"valid duration", sleepFor(1), false, false},
		{"past time", sleepFor(-1), false, false},
		{"json with iso8601", `{"until":"2018-06-19T12:54:49.000Z"}`, false, false},
		{"max duration", sleepFor(2592000), false, false},
		{"invalid json", `{"until":"1000h"}`, true, false},
	}

	store, cleanup := cltest.NewStore()
	defer cleanup()
	store.Clock = cltest.InstantClock{}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			input := models.RunResult{}
			adapter := adapters.Sleep{}
			err := json.Unmarshal([]byte(test.params), &adapter)
			if test.parseErrored {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			result := adapter.Perform(input, store)
			if test.errored {
				assert.Error(t, result.GetError())
			} else {
				assert.NoError(t, result.GetError())
			}
		})
	}
}
