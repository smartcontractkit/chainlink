package types

import (
	"github.com/smartcontractkit/chainlink-terra/pkg/terra/config"
)

var defaultCfg = config.ChainCfg{
	FallbackGasPriceULuna: "0.01",
	GasLimitMultiplier:    1.5,
}

// NewChainScopedConfig returns a ChainCfg with defaults overridden by dbcfg.
func NewChainScopedConfig(dbcfg ChainCfg) (cfg config.ChainCfg) {
	cfg = defaultCfg
	// override defaults, if set
	if v := dbcfg.FallbackGasPriceULuna; v.Valid {
		cfg.FallbackGasPriceULuna = v.String
	}
	if v := dbcfg.GasLimitMultiplier; v.Valid {
		cfg.GasLimitMultiplier = v.Float64
	}
	return
}
