package v2

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func ChainEthMainnet(t *testing.T) config.ChainScopedConfig      { return scopedConfig(t, 1) }
func ChainOptimismMainnet(t *testing.T) config.ChainScopedConfig { return scopedConfig(t, 10) }
func ChainArbitrumMainnet(t *testing.T) config.ChainScopedConfig { return scopedConfig(t, 42161) }
func ChainArbitrumRinkeby(t *testing.T) config.ChainScopedConfig { return scopedConfig(t, 421611) }

func scopedConfig(t *testing.T, chainID int64) config.ChainScopedConfig {
	id := big.NewI(chainID)
	evmCfg := toml.EVMConfig{ChainID: id, Chain: toml.Defaults(id)}
	return config.NewTOMLChainScopedConfig(&evmCfg, logger.TestLogger(t))
}
