package adapters_test

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
)

func TestRandom_Perform(t *testing.T) {
	tests := []struct {
		name      string
		params    string
		errored   bool
		jsonError bool
	}{
		{"simple 1", `{"start":0, "end": 100}`, false, false},
		{"simple 2", `{"start":-100, "end": 0}`, false, false},
		{"simple 3", `{"start":-100, "end": 100}`, false, false},
		{"big nums", `{"start":9223372036854775805, "end": 9223372036854775807}`, false, false},
		{"start neg", `{"start":-25000, "end": 5}`, false, false},
		{"both neg", `{"start":-25000, "end": -5000}`, false, false},
		{"start is 0 by default if not specified", `{"end": 5}`, false, false},
		{"end is 0 by default if not specified", `{"start": -50}`, false, false},
		{"string params", `{"start":"100", "end": "200"}`, false, false},
		{"string params (one neg)", `{"start":"-100", "end": "200"}`, false, false},
		{"string params (both neg)", `{"start":"-500", "end": "-100"}`, false, false},

		{"end less than start (end neg)", `{"start":3, "end": -1}`, true, false},
		{"end less than start (both neg)", `{"start":-5000, "end": -25000}`, true, false},
		{"error on end overflow", `{"start":9223372036854775805, "end": 9223372036854775809}`, false, true},
		{"error on start underflow", `{"start":-9223372036854775809, "end": 0}`, false, true},
		{"start param defaults to 0", `{"end": -5}`, true, false},
		{"end param defaults to 0", `{"start": 5}`, true, false},
		{"no start and no end parameters", `{}`, true, false},
		{"invalid params", `{"start": "1s", "end": "20b"}`, false, true},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := models.RunResult{}
			adapter := adapters.Random{}
			jsonErr := json.Unmarshal([]byte(test.params), &adapter)
			result := adapter.Perform(input, nil)

			if test.jsonError {
				assert.Error(t, jsonErr)
			} else if test.errored {
				assert.Error(t, result.GetError())
				assert.NoError(t, jsonErr)
			} else {
				val, err := result.ResultString()
				assert.NoError(t, err)
				assert.NoError(t, result.GetError())
				assert.NoError(t, jsonErr)
				valInt, err := strconv.ParseInt(val, 10, 64)
				// sanity check the result is in the interval
				assert.GreaterOrEqual(t, valInt, int64(adapter.Start))
				assert.LessOrEqual(t, valInt, int64(adapter.End))
			}
		})
	}
}
