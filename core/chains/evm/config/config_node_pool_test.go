package config_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	v2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestNodePoolConfig(t *testing.T) {
	gcfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		id := utils.NewBig(big.NewInt(rand.Int63()))
		c.EVM[0] = &v2.EVMConfig{
			ChainID: id,
			Chain:   v2.Defaults(id, &v2.Chain{}),
		}
	})
	cfg := evmtest.NewChainScopedConfig(t, gcfg)

	require.Equal(t, "HighestHead", cfg.EVM().NodePool().SelectionMode())
	require.Equal(t, uint32(5), cfg.EVM().NodePool().SyncThreshold())
	require.Equal(t, time.Duration(10000000000), cfg.EVM().NodePool().PollInterval())
	require.Equal(t, uint32(5), cfg.EVM().NodePool().PollFailureThreshold())
}
