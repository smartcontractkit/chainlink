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
	res := new(big.Int)
	res, ok := res.SetString(result.Result().String(), 10)
	assert.True(t, ok)
}
