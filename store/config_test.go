package store_test

import (
	"math/big"
	"testing"

	strpkg "github.com/smartcontractkit/chainlink/store"
	"github.com/stretchr/testify/assert"
)

func TestStore_ConfigDefaults(t *testing.T) {
	config := strpkg.NewConfig()
	assert.Equal(t, uint64(0), config.ChainID)
	assert.Equal(t, *big.NewInt(20000000000), config.EthGasPriceDefault)
}
