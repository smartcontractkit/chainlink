package bulletprooftxmanager_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"

	"github.com/stretchr/testify/assert"
)

func TestBulletproofTxManager_BumpGas(t *testing.T) {
	t.Parallel()

	config, cleanup := cltest.NewConfig(t)
	defer cleanup()

	t.Run("returns gas price based on default if that is largest", func(t *testing.T) {
		config.Set("ETH_GAS_BUMP_PERCENT", 10)
		config.Set("ETH_GAS_BUMP_WEI", 5)
		config.Set("ETH_GAS_PRICE_DEFAULT", 30)

		originalGasPrice := big.NewInt(10)

		assert.Equal(t, big.NewInt(35), bulletprooftxmanager.BumpGas(config, originalGasPrice))
	})

	t.Run("returns percentage bump if that is largest", func(t *testing.T) {
		config.Set("ETH_GAS_BUMP_PERCENT", 10)
		config.Set("ETH_GAS_BUMP_WEI", 5)
		config.Set("ETH_GAS_PRICE_DEFAULT", 30)

		originalGasPrice := big.NewInt(100)

		assert.Equal(t, big.NewInt(110), bulletprooftxmanager.BumpGas(config, originalGasPrice))
	})

	t.Run("returns fixed size bump if that is largest", func(t *testing.T) {
		config.Set("ETH_GAS_BUMP_PERCENT", 10)
		config.Set("ETH_GAS_BUMP_WEI", 5)
		config.Set("ETH_GAS_PRICE_DEFAULT", 25)

		originalGasPrice := big.NewInt(29)

		assert.Equal(t, big.NewInt(34), bulletprooftxmanager.BumpGas(config, originalGasPrice))
	})

	t.Run("caps at EthMaxGasPriceWei if bump would exceed this", func(t *testing.T) {
		config.Set("ETH_GAS_BUMP_PERCENT", 10)
		config.Set("ETH_GAS_BUMP_WEI", 5)
		config.Set("ETH_GAS_PRICE_DEFAULT", 30)
		config.Set("ETH_MAX_GAS_PRICE_WEI", 500000000000)

		originalGasPrice := big.NewInt(480000000000)

		assert.Equal(t, big.NewInt(500000000000), bulletprooftxmanager.BumpGas(config, originalGasPrice))
	})
}
