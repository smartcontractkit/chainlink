package adapters_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestSleep_Perform(t *testing.T) {
	tests := []struct {
		name         string
		params       string
		parseErrored bool
		errored      bool
	}{
		{"valid duration", `{"seconds":30}`, false, false},
		{"excessive duration", `{"seconds":259201}`, false, true},
		{"invalid json", `{"seconds":"1000h"}`, true, false},
	}

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

			result := adapter.Perform(input, nil)
			if test.errored {
				assert.Error(t, result.GetError())
			} else {
				assert.NoError(t, result.GetError())
			}
		})
	}
}
