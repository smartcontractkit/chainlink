package config_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/config"

	"github.com/stretchr/testify/require"
)

func TestConfig_EthGasPriceDefault(t *testing.T) {
	store, cleanup := cltest.NewStore(t)
	defer cleanup()
	cfg := store.Config

	// Get default value
	def := cfg.EthGasPriceDefault()

	// No orm installed
	err := cfg.SetEthGasPriceDefault(big.NewInt(0))
	require.Error(t, err)

	// Install ORM
	orm := config.NewORM(store.DB)
	cfg.SetRuntimeStore(orm)

	// Value still stays as the default
	require.Equal(t, def, cfg.EthGasPriceDefault())

	// Override
	newValue := new(big.Int).Add(def, big.NewInt(1))
	err = cfg.SetEthGasPriceDefault(newValue)
	require.NoError(t, err)

	// Value changes
	require.Equal(t, newValue, cfg.EthGasPriceDefault())

	// Set again
	newerValue := new(big.Int).Add(def, big.NewInt(2))
	err = cfg.SetEthGasPriceDefault(newerValue)
	require.NoError(t, err)

	// Value changes
	require.Equal(t, newerValue, cfg.EthGasPriceDefault())
}
