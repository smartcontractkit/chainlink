package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
)

func TestEthBytes32_Perform(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected string
	}{
		{"string", `{"result":"Hello World!"}`, "0x48656c6c6f20576f726c64210000000000000000000000000000000000000000"},
		{"special characters", `{"result":"¡Holá Mündo!"}`, "0xc2a1486f6cc3a1204dc3bc6e646f210000000000000000000000000000000000"},
		{"long string", `{"result":"string that is waaAAAaaay toooo long!!!!!"}`, "0x737472696e672074686174206973207761614141416161617920746f6f6f6f20"},
		{"empty string", `{"result":""}`, "0x0000000000000000000000000000000000000000000000000000000000000000"},
		{"string of number", `{"result":"16800.01"}`, "0x31363830302e3031000000000000000000000000000000000000000000000000"},
		{"float", `{"result":16800.01}`, "0x31363830302e3031000000000000000000000000000000000000000000000000"},
		{"scientific float", `{"result":1.68e+4}`,
			"0x3136383030000000000000000000000000000000000000000000000000000000"},
		{"roundable float", `{"result":16800.00}`, "0x3136383030000000000000000000000000000000000000000000000000000000"},
		{"integer", `{"result":16800}`, "0x3136383030000000000000000000000000000000000000000000000000000000"},
		{"boolean true", `{"result":true}`, "0x7472756500000000000000000000000000000000000000000000000000000000"},
		{"boolean false", `{"result":false}`, "0x66616c7365000000000000000000000000000000000000000000000000000000"},
		{"null", `{"result":null}`, "0x0000000000000000000000000000000000000000000000000000000000000000"},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := cltest.JSONFromString(t, test.json)
			result := models.RunResult{}
			adapter := adapters.EthBytes32{}
			result = adapter.Perform(input, result, nil)

			val, err := result.ResultString()
			assert.NoError(t, err)
			assert.NoError(t, result.GetError())
			assert.Equal(t, test.expected, val)
		})
	}
}

func TestEthInt256_Perform(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		json    string
		want    string
		errored bool
	}{
		{"string", `{"result":"123"}`,
			"0x000000000000000000000000000000000000000000000000000000000000007b", false},
		{"integer", `{"result":123}`,
			"0x000000000000000000000000000000000000000000000000000000000000007b", false},
		{"float", `{"result":123.0}`,
			"0x000000000000000000000000000000000000000000000000000000000000007b", false},
		{"rounded float", `{"result":123.99}`,
			"0x000000000000000000000000000000000000000000000000000000000000007b", false},
		{"negative string", `{"result":"-123"}`,
			"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff85", false},
		{"scientific", `{"result":1.68e+4}`,
			"0x00000000000000000000000000000000000000000000000000000000000041a0", false},
		{"scientific string", `{"result":"1.68e+4"}`,
			"0x00000000000000000000000000000000000000000000000000000000000041a0", false},
		{"negative float", `{"result":-123.99}`,
			"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff85", false},
		{"object", `{"result":{"a": "b"}}`, "", true},
	}

	adapter := adapters.EthInt256{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := cltest.JSONFromString(t, test.json)
			result := models.RunResult{}
			result = adapter.Perform(input, result, nil)

			if test.errored {
				assert.Error(t, result.GetError())
			} else {
				val, err := result.ResultString()
				assert.NoError(t, result.GetError())
				assert.NoError(t, err)
				assert.Equal(t, test.want, val)
			}
		})
	}
}

func TestEthUint256_Perform(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		json    string
		want    string
		errored bool
	}{
		{"string", `{"result":"123"}`,
			"0x000000000000000000000000000000000000000000000000000000000000007b", false},
		{"integer", `{"result":123}`,
			"0x000000000000000000000000000000000000000000000000000000000000007b", false},
		{"float", `{"result":123.0}`,
			"0x000000000000000000000000000000000000000000000000000000000000007b", false},
		{"rounded float", `{"result":123.99}`,
			"0x000000000000000000000000000000000000000000000000000000000000007b", false},
		{"scientific", `{"result":1.68e+4}`,
			"0x00000000000000000000000000000000000000000000000000000000000041a0", false},
		{"scientific string", `{"result":"1.68e+4"}`,
			"0x00000000000000000000000000000000000000000000000000000000000041a0", false},
		{"negative integer", `{"result":-123}`, "", true},
		{"negative string", `{"result":"-123"}`, "", true},
		{"negative float", `{"result":-123.99}`, "", true},
		{"object", `{"result":{"a": "b"}}`, "", true},
	}

	adapter := adapters.EthUint256{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := cltest.JSONFromString(t, test.json)
			result := models.RunResult{}
			result = adapter.Perform(input, result, nil)

			if test.errored {
				assert.Error(t, result.GetError())
			} else {
				val, err := result.ResultString()
				assert.NoError(t, result.GetError())
				assert.NoError(t, err)
				assert.Equal(t, test.want, val)
			}
		})
	}
}
