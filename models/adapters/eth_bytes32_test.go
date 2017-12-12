package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-go/models/adapters"
	"github.com/stretchr/testify/assert"
)

func TestEthereumBytes32Formatting(t *testing.T) {
	tests := []struct {
		value    string
		expected string
	}{
		{"16800.00", "31363830302e3030000000000000000000000000000000000000000000000000"},
		{"", "0000000000000000000000000000000000000000000000000000000000000000"},
		{"Hello World!", "48656c6c6f20576f726c64210000000000000000000000000000000000000000"},
		{"string that is waaAAAaaay toooo long!!!!!", "737472696e672074686174206973207761614141416161617920746f6f6f6f20"},
	}

	for _, test := range tests {
		past := adapters.RunResultWithValue(test.value)
		adapter := adapters.EthBytes32{}
		result := adapter.Perform(past)

		assert.Equal(t, test.expected, result.Value())
		assert.Nil(t, result.Error)
	}
}
