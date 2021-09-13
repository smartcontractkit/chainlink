package config

import "github.com/smartcontractkit/chainlink/core/chains/evm/types"

func PersistedCfgPtr(cfg ChainScopedConfig) *types.ChainCfg {
	return &cfg.(*chainScopedConfig).persistedCfg
}
