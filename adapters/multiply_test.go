package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestMultiply_Perform(t *testing.T) {
	tests := []struct {
		name    string
		times   interface{}
		json    string
		want    string
		errored bool
	}{
		{"string", 100, `{"value":"1.23"}`, "123", false},
		{"integer", 100, `{"value":123}`, "12300", false},
		{"float", 100, `{"value":1.23}`, "123", false},
		{"object", 100, `{"value":{"foo":"bar"}}`, "", true},
		{"string_string", "100", `{"value":"1.23"}`, "123", false},
		{"string_integer", "100", `{"value":123}`, "12300", false},
		{"string_float", "100", `{"value":1.23}`, "123", false},
		{"string_object", "100", `{"value":{"foo":"bar"}}`, "", true},
		{"rubbish_string", "123aaa123", `{"value":"1.23"}`, "", true},
		{"slice_string", []int{1, 2, 3}, `{"value":"1.23"}`, "", true},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := models.RunResult{
				Data: cltest.JSONFromString(test.json),
			}
			adapter := adapters.Multiply{Times: test.times}
			result := adapter.Perform(input, nil)

			if test.errored {
				assert.NotNil(t, result.GetError())
			} else {
				val, err := result.Value()
				assert.Nil(t, err)
				assert.Equal(t, test.want, val)
				assert.Nil(t, result.GetError())
			}
		})
	}
}
