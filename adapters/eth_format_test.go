package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/smartcontractkit/chainlink/utils"
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
	t.Parallel()
	tests := []struct {
		name    string
		json    string
		want    string
		errored bool
	}{
		{"string", `{"value":"123"}`, utils.EVMHexNumber(123), false},
		{"integer", `{"value":123}`, utils.EVMHexNumber(123), false},
		{"integer", `{"value":"18446744073709551615"}`,
			"0x000000000000000000000000000000000000000000000000ffffffffffffffff", false},
		{"integer", `{"value":"170141183460469231731687303715884105728"}`,
			"0x0000000000000000000000000000000080000000000000000000000000000000", false},
		{"integer", `{"value":"170141183460469231731687303715884105729"}`,
			"0x0000000000000000000000000000000080000000000000000000000000000001", false},
		{"2^128", `{"value":"340282366920938463463374607431768211456"}`,
			"0x0000000000000000000000000000000100000000000000000000000000000000", false},
		{"large float precision", `{"value":"115792089237316195423570985008687907853269984665640564039457584007913129639934"}`,
			"0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe", false},
		{"2^256 - 1", `{"value":"115792089237316195423570985008687907853269984665640564039457584007913129639935"}`,
			"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", false},
		{"2^256 - 0.1", `{"value":"115792089237316195423570985008687907853269984665640564039457584007913129639935.9"}`,
			"0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", false},
		{"2^256", `{"value":"115792089237316195423570985008687907853269984665640564039457584007913129639936"}`,
			"", true},
		{"float", `{"value":123.0}`, utils.EVMHexNumber(123), false},
		{"rounded float", `{"value":123.99}`, utils.EVMHexNumber(123), false},
		{"negative integer", `{"value":-123}`, "", true},
		{"negative string", `{"value":"-123"}`, "", true},
		{"negative float", `{"value":-123.99}`, "", true},
		{"object", `{"value":{"a": "b"}}`, "", true},
		{"odd length result", `{"value":"1234"}`, utils.EVMHexNumber(1234), false},
	}

	adapter := adapters.EthUint256{}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := models.RunResult{
				Data: cltest.JSONFromString(test.json),
			}
			result := adapter.Perform(input, nil)

			if test.errored {
				assert.NotNil(t, result.GetError())
			} else {
				val, err := result.Value()
				assert.Nil(t, result.GetError())
				assert.Nil(t, err)
				assert.Equal(t, test.want, val)
			}
		})
	}
}
