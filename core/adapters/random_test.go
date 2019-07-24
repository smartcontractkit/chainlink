package adapters_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
)

func TestRandom_Perform(t *testing.T) {
	input := models.RunResult{}
	adapter := adapters.Random{}
	result := adapter.Perform(input, nil)
	val, err := result.ResultString()
	assert.NoError(t, err)
	assert.NoError(t, result.GetError())
	res := new(big.Int)
	res, ok := res.SetString(val, 10)
	assert.True(t, ok)
}
