package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestEthBytes32Formatting(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected string
	}{
		{"string", `{"value":"Hello World!"}`, "48656c6c6f20576f726c64210000000000000000000000000000000000000000"},
		{"special characters", `{"value":"¡Holá Mündo!"}`, "c2a1486f6cc3a1204dc3bc6e646f210000000000000000000000000000000000"},
		{"long string", `{"value":"string that is waaAAAaaay toooo long!!!!!"}`, "737472696e672074686174206973207761614141416161617920746f6f6f6f20"},
		{"empty string", `{"value":""}`, "0000000000000000000000000000000000000000000000000000000000000000"},
		{"string of number", `{"value":"16800.01"}`, "31363830302e3031000000000000000000000000000000000000000000000000"},
		{"float", `{"value":16800.01}`, "31363830302e3031000000000000000000000000000000000000000000000000"},
		{"roundable float", `{"value":16800.00}`, "3136383030000000000000000000000000000000000000000000000000000000"},
		{"integer", `{"value":16800}`, "3136383030000000000000000000000000000000000000000000000000000000"},
		{"boolean true", `{"value":true}`, "7472756500000000000000000000000000000000000000000000000000000000"},
		{"boolean false", `{"value":false}`, "66616c7365000000000000000000000000000000000000000000000000000000"},
		{"null", `{"value":null}`, "0000000000000000000000000000000000000000000000000000000000000000"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			past := models.RunResult{Output: &models.Output{test.json}}
			adapter := adapters.EthBytes32{}
			result := adapter.Perform(past, nil)

			val, err := result.Value()
			assert.Equal(t, test.expected, val)
			assert.Nil(t, err)
			assert.Nil(t, result.GetError())
		})
	}
}
