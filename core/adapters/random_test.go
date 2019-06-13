package adapters_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
)

func TestRandom_Perform(t *testing.T) {
	tests := []struct {
		name    string
		errored bool
	}{
		{"resulting string parses to a big.Int 1", false},
		{"resulting string parses to a big.Int 2", false},
		{"resulting string parses to a big.Int 3", false},
		{"resulting string parses to a big.Int 4", false},
		{"resulting string parses to a big.Int 5", false},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			input := models.RunResult{}
			adapter := adapters.Random{}
			result := adapter.Perform(input, nil)

			if test.errored {
				assert.Error(t, result.GetError())
			} else {
				val, err := result.ResultString()
				assert.NoError(t, err)
				assert.NoError(t, result.GetError())
				res := new(big.Int)
				res, ok := res.SetString(val, 10)
				assert.True(t, ok)
			}
		})
	}
}
