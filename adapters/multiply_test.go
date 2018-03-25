package adapters_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestMultiply_Perform(t *testing.T) {
	tests := []struct {
		name      string
		params    string
		json      string
		want      string
		errored   bool
		jsonError bool
	}{
		{"string", `{"times":100}`, `{"value":"1.23"}`, "123", false, false},
		{"integer", `{"times":100}`, `{"value":123}`, "12300", false, false},
		{"float", `{"times":100}`, `{"value":1.23}`, "123", false, false},
		{"object", `{"times":100}`, `{"value":{"foo":"bar"}}`, "", true, false},
		{"zero_integer_string", `{"times":0}`, `{"value":"1.23"}`, "0", false, false},
		{"negative_integer_string", `{"times":-5}`, `{"value":"1.23"}`, "-6.15", false, false},

		{"string_string", `{"times":"100"}`, `{"value":"1.23"}`, "123", false, false},
		{"string_integer", `{"times":"100"}`, `{"value":123}`, "12300", false, false},
		{"string_float", `{"times":"100"}`, `{"value":1.23}`, "123", false, false},
		{"string_object", `{"times":"100"}`, `{"value":{"foo":"bar"}}`, "", true, false},
		{"array_string", `{"times":[1, 2, 3]}`, `{"value":"1.23"}`, "", false, true},
		{"rubbish_string", `{"times":"123aaa123"}`, `{"value":"1.23"}`, "", false, true},
		{"zero_string_string", `{"times":"0"}`, `{"value":"1.23"}`, "0", false, false},
		{"negative_string_string", `{"times":"-5"}`, `{"value":"1.23"}`, "-6.15", false, false},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := models.RunResult{
				Data: cltest.JSONFromString(test.json),
			}
			adapter := adapters.Multiply{}
			jsonErr := json.Unmarshal([]byte(test.params), &adapter)
			result := adapter.Perform(input, nil)

			if test.jsonError {
				assert.NotNil(t, jsonErr)
			} else if test.errored {
				assert.NotNil(t, result.GetError())
				assert.Nil(t, jsonErr)
			} else {
				val, err := result.Value()
				assert.Nil(t, err)
				assert.Equal(t, test.want, val)
				assert.Nil(t, result.GetError())
				assert.Nil(t, jsonErr)
			}
		})
	}
}
