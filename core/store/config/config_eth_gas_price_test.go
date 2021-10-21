package config_test

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/store/config"

	"github.com/stretchr/testify/require"
)

func TestEVMConfig_EvmGasPriceDefault(t *testing.T) {
	cfg := config.NewEVMConfig(config.NewGeneralConfig())

	// Get default value
	def := cfg.EvmGasPriceDefault()

	// No orm installed
	err := cfg.SetEvmGasPriceDefault(big.NewInt(0))
	require.Error(t, err)

	// Install ORM
	db := pgtest.NewGormDB(t)
	cfg.SetDB(db)

	// Value still stays as the default
	require.Equal(t, def, cfg.EvmGasPriceDefault())

	// Override
	newValue := new(big.Int).Add(def, big.NewInt(1))
	err = cfg.SetEvmGasPriceDefault(newValue)
	require.NoError(t, err)

	// Value changes
	require.Equal(t, newValue, cfg.EvmGasPriceDefault())

	// Set again
	newerValue := new(big.Int).Add(def, big.NewInt(2))
	err = cfg.SetEvmGasPriceDefault(newerValue)
	require.NoError(t, err)

	// Value changes
	require.Equal(t, newerValue, cfg.EvmGasPriceDefault())
}
