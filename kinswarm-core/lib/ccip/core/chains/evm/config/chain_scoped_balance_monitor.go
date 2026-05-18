package config

import "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"

type balanceMonitorConfig struct {
	c toml.BalanceMonitor
}

func (b *balanceMonitorConfig) Enabled() bool {
	return *b.c.Enabled
}
