package config

type FluxMonitor interface {
	DefaultTransactionQueueDepth() uint32
	SimulateTransactions() bool
}
