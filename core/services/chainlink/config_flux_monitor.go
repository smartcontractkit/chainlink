package chainlink

import "github.com/smartcontractkit/chainlink/v2/core/config/toml"

type fluxMonitorConfig struct {
	c toml.FluxMonitor
}

func (f *fluxMonitorConfig) DefaultTransactionQueueDepth() uint32 {
	return *f.c.DefaultTransactionQueueDepth
}

func (f *fluxMonitorConfig) SimulateTransactions() bool {
	return *f.c.SimulateTransactions
}
