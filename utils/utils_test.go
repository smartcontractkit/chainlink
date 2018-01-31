package utils_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/utils"
	"github.com/stretchr/testify/assert"
	"math/big"
)

func TestNewBytes32ID(t *testing.T) {
	t.Parallel()

	id := utils.NewBytes32ID()
	assert.NotContains(t, id, "-")
}

func TestWeiToEth(t *testing.T) {
	var numWei *big.Int = new(big.Int).SetInt64(1)
	var expectedNumEth float64 = 1e-18
	actualNumEth := utils.WeiToEth(numWei)
	assert.Equal(t, expectedNumEth, actualNumEth)
}

func TestEthToWei(t *testing.T) {
	var numEth float64 = 1.0
	var expectedNumWei *big.Int = new(big.Int).SetInt64(1e18)
	actualNumWei := utils.EthToWei(numEth)
	assert.Equal(t, actualNumWei, expectedNumWei)
}
