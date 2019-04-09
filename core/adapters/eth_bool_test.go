package adapters_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/tools/cltest"
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
		{"value is bool false", `{"result":false}`, evmFalse},
		{"value is bool true", `{"result":true}`, evmTrue},
		{"value is null string", `{"result":"null"}`, evmTrue},
		{"value is null", `{"result":null}`, evmFalse},
		{"empty object", `{}`, evmFalse},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			past := models.RunResult{
				Data: cltest.JSONFromString(t, test.json),
			}
			adapter := adapters.EthBool{}
			result := adapter.Perform(past, nil)

			val, err := result.ResultString()
			assert.Equal(t, test.expected, val)
			assert.NoError(t, err)
			assert.NoError(t, result.GetError())
		})
	}
}
