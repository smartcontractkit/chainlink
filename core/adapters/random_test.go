package adapters_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandom_Perform(t *testing.T) {
	adapter := adapters.Random{}
	result := adapter.Perform(models.RunInput{}, nil)
	require.NoError(t, result.Error())
	val, err := result.ResultString()
	require.NoError(t, err)
	res := new(big.Int)
	res, ok := res.SetString(val, 10)
	assert.True(t, ok)
}
