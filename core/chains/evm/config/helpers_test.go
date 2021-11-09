package config

import "github.com/smartcontractkit/chainlink/core/chains/evm/types"

func UpdatePersistedCfg(cfg ChainScopedConfig, updateFn func(*types.ChainCfg)) {
	c := cfg.(*chainScopedConfig)
	c.persistMu.Lock()
	defer c.persistMu.Unlock()
	updateFn(&c.persistedCfg)
}

func ChainSpecificConfigDefaultSets() map[int64]chainSpecificConfigDefaultSet {
	return chainSpecificConfigDefaultSets
}
