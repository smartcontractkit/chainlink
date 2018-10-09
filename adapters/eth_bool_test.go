package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

var evmFalse = "0x0000000000000000000000000000000000000000000000000000000000000000"
var evmTrue = "0x0000000000000000000000000000000000000000000000000000000000000001"

func TestEthBool_Perform(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected string
	}{
		{"bool false", `{"value":false}`, evmFalse},
		{"bool true", `{"value":true}`, evmTrue},
		{"true string", `{"value":"true"}`, evmTrue},
		{"false string", `{"value":"false"}`, evmTrue},
		{"number 5", `{"value":5}`, evmTrue},
		{"number 0", `{"value":0}`, evmTrue},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			past := models.RunResult{
				Data: cltest.JSONFromString(test.json),
			}
			adapter := adapters.EthBool{}
			result := adapter.Perform(past, nil)

			val, err := result.Value()
			assert.Equal(t, test.expected, val)
			assert.NoError(t, err)
			assert.NoError(t, result.GetError())
		})
	}
}
