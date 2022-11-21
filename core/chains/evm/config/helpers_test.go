package config

import "github.com/smartcontractkit/chainlink/core/chains/evm/types"

// Deprecated: https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func UpdatePersistedCfg(cfg ChainScopedConfig, updateFn func(*types.ChainCfg)) {
	c := cfg.(*chainScopedConfig)
	c.persistMu.Lock()
	defer c.persistMu.Unlock()
	updateFn(&c.persistedCfg)
}

func ChainSpecificConfigDefaultSets() map[int64]chainSpecificConfigDefaultSet {
	return chainSpecificConfigDefaultSets
}
