package v2

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	v2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func ChainEthMainnet(t *testing.T) config.ChainScopedConfig      { return scopedConfig(t, 1) }
func ChainOptimismMainnet(t *testing.T) config.ChainScopedConfig { return scopedConfig(t, 10) }
func ChainOptimismKovan(t *testing.T) config.ChainScopedConfig   { return scopedConfig(t, 69) }
func ChainArbitrumMainnet(t *testing.T) config.ChainScopedConfig { return scopedConfig(t, 42161) }
func ChainArbitrumRinkeby(t *testing.T) config.ChainScopedConfig { return scopedConfig(t, 421611) }

func scopedConfig(t *testing.T, chainID int64) config.ChainScopedConfig {
	id := utils.NewBigI(chainID)
	evmCfg := v2.EVMConfig{ChainID: id, Chain: v2.Defaults(id)}
	return v2.NewTOMLChainScopedConfig(configtest.NewTestGeneralConfig(t), &evmCfg, logger.TestLogger(t))
}
