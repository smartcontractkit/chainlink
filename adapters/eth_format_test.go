package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-go/adapters"
	"github.com/smartcontractkit/chainlink-go/models"
	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
)

func TestEthereumBytes32Formatting(t *testing.T) {
	tests := []struct {
		value    null.String
		expected string
	}{
		{
			null.StringFrom("16800.00"),
			"31363830302e3030000000000000000000000000000000000000000000000000",
		},
		{
			null.StringFrom(""),
			"0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			null.StringFrom("Hello World!"),
			"48656c6c6f20576f726c64210000000000000000000000000000000000000000",
		},
		{
			null.StringFrom("string that is waaAAAaaay toooo long!!!!!"),
			"737472696e672074686174206973207761614141416161617920746f6f6f6f20",
		},
		{
			null.StringFrom("¡Hola Mundo!"),
			"c2a1486f6c61204d756e646f2100000000000000000000000000000000000000",
		},
		{
			null.StringFromPtr(nil),
			"0000000000000000000000000000000000000000000000000000000000000000",
		},
	}

	for _, test := range tests {
		past := models.RunResult{
			Output: models.Output{"value": test.value},
		}
		adapter := adapters.EthBytes32{}
		result := adapter.Perform(past)

		assert.Equal(t, test.expected, result.Value())
		assert.Nil(t, result.GetError())
	}
}
