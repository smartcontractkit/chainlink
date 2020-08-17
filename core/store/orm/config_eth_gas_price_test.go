package orm_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"

	"github.com/stretchr/testify/require"
)

func TestConfig_EthGasPriceDefault(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	config := store.Config

	// Get default value
	def := config.EthGasPriceDefault()

	// No orm installed
	err := config.SetEthGasPriceDefault(big.NewInt(0))
	require.Error(t, err)

	// Install ORM
	config.SetRuntimeStore(store.ORM)

	// Value still stays as the default
	require.Equal(t, def, config.EthGasPriceDefault())

	// Override
	newValue := new(big.Int).Add(def, big.NewInt(1))
	err = config.SetEthGasPriceDefault(newValue)
	require.NoError(t, err)

	// Value changes
	require.Equal(t, newValue, config.EthGasPriceDefault())

	// Set again
	newerValue := new(big.Int).Add(def, big.NewInt(2))
	err = config.SetEthGasPriceDefault(newerValue)
	require.NoError(t, err)

	// Value changes
	require.Equal(t, newerValue, config.EthGasPriceDefault())
}
