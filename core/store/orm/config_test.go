package orm_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestConfig_SetEthGasPriceDefault(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	t.Cleanup(cleanup)
	config := store.Config

	config.Set("ETH_MAX_GAS_PRICE_WEI", 1500000000000)

	t.Run("sets the gas price", func(t *testing.T) {
		assert.Equal(t, big.NewInt(20000000000), config.EthGasPriceDefault())

		err := config.SetEthGasPriceDefault(big.NewInt(42000000000))
		assert.NoError(t, err)

		assert.Equal(t, big.NewInt(42000000000), config.EthGasPriceDefault())
	})
	t.Run("is not allowed to set gas price to below EthMinGasPriceWei", func(t *testing.T) {
		assert.Equal(t, big.NewInt(1000000000), config.EthMinGasPriceWei())

		err := config.SetEthGasPriceDefault(big.NewInt(1))
		assert.EqualError(t, err, "cannot set default gas price to 1, it is below the minimum allowed value of 1000000000")

		assert.Equal(t, big.NewInt(42000000000), config.EthGasPriceDefault())
	})
	t.Run("is not allowed to set gas price to above EthMaxGasPriceWei", func(t *testing.T) {
		assert.Equal(t, big.NewInt(1500000000000), config.EthMaxGasPriceWei())

		err := config.SetEthGasPriceDefault(big.NewInt(999999999999999))
		assert.EqualError(t, err, "cannot set default gas price to 999999999999999, it is above the maximum allowed value of 1500000000000")

		assert.Equal(t, big.NewInt(42000000000), config.EthGasPriceDefault())
	})
}
