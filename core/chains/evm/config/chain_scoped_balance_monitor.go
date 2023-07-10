package config

import v2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"

type balanceMonitorConfig struct {
	c v2.BalanceMonitor
}

func (b *balanceMonitorConfig) Enabled() bool {
	return *b.c.Enabled
}
