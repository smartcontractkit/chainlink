package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestEthBytes32_Perform(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected string
	}{
		{"string", `{"value":"Hello World!"}`, "0x48656c6c6f20576f726c64210000000000000000000000000000000000000000"},
		{"special characters", `{"value":"¡Holá Mündo!"}`, "0xc2a1486f6cc3a1204dc3bc6e646f210000000000000000000000000000000000"},
		{"long string", `{"value":"string that is waaAAAaaay toooo long!!!!!"}`, "0x737472696e672074686174206973207761614141416161617920746f6f6f6f20"},
		{"empty string", `{"value":""}`, "0x0000000000000000000000000000000000000000000000000000000000000000"},
		{"string of number", `{"value":"16800.01"}`, "0x31363830302e3031000000000000000000000000000000000000000000000000"},
		{"float", `{"value":16800.01}`, "0x31363830302e3031000000000000000000000000000000000000000000000000"},
		{"roundable float", `{"value":16800.00}`, "0x3136383030000000000000000000000000000000000000000000000000000000"},
		{"integer", `{"value":16800}`, "0x3136383030000000000000000000000000000000000000000000000000000000"},
		{"boolean true", `{"value":true}`, "0x7472756500000000000000000000000000000000000000000000000000000000"},
		{"boolean false", `{"value":false}`, "0x66616c7365000000000000000000000000000000000000000000000000000000"},
		{"null", `{"value":null}`, "0x0000000000000000000000000000000000000000000000000000000000000000"},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			past := models.RunResult{
				Data: cltest.JSONFromString(test.json),
			}
			adapter := adapters.EthBytes32{}
			result := adapter.Perform(past, nil)

			val, err := result.Value()
			assert.Equal(t, test.expected, val)
			assert.Nil(t, err)
			assert.Nil(t, result.GetError())
		})
	}
}

func TestEthUint256_Perform(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    string
		errored bool
	}{
		{"string", `{"value":"123"}`, "0x000000000000000000000000000000000000000000000000000000000000007b", false},
		{"integer", `{"value":123}`, "0x000000000000000000000000000000000000000000000000000000000000007b", false},
		{"integer", `{"value":"18446744073709551615"}`, "0x000000000000000000000000000000000000000000000000ffffffffffffffff", false},
		{"integer", `{"value":"170141183460469231731687303715884105728"}`, "0x0000000000000000000000000000000080000000000000000000000000000000", false},
		{"float", `{"value":123.0}`, "0x000000000000000000000000000000000000000000000000000000000000007b", false},
		{"rounded float", `{"value":123.99}`, "0x000000000000000000000000000000000000000000000000000000000000007b", false},
		{"negative integer", `{"value":-123}`, "0x000000000000000000000000000000000000000000000000000000000000007b", false},
		{"negative string", `{"value":"-123"}`, "0x000000000000000000000000000000000000000000000000000000000000007b", false},
		{"negative float", `{"value":-123.99}`, "0x000000000000000000000000000000000000000000000000000000000000007b", false},
		{"object", `{"value":{"a": "b"}}`, "", true},
		{"odd length result", `{"value":"1234"}`, "0x00000000000000000000000000000000000000000000000000000000000004d2", false},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := models.RunResult{
				Data: cltest.JSONFromString(test.json),
			}
			adapter := adapters.EthUint256{}
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
