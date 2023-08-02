package chainlink

import v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"

type fluxMonitorConfig struct {
	c v2.FluxMonitor
}

func (f *fluxMonitorConfig) DefaultTransactionQueueDepth() uint32 {
	return *f.c.DefaultTransactionQueueDepth
}

func (f *fluxMonitorConfig) SimulateTransactions() bool {
	return *f.c.SimulateTransactions
}
