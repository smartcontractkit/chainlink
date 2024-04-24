package testutils

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

func NewTestChainScopedConfig(t testing.TB, overrideFn func(c *toml.EVMConfig)) config.ChainScopedConfig {
	var chainID = (*big.Big)(FixtureChainID)
	evmCfg := &toml.EVMConfig{
		ChainID: chainID,
		Chain:   toml.Defaults(chainID),
	}

	if overrideFn != nil {
		overrideFn(evmCfg)
	}

	return config.NewTOMLChainScopedConfig(evmCfg, logger.Test(t))
}
