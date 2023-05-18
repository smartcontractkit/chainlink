package config

type FluxMonitor interface {
	FMDefaultTransactionQueueDepth() uint32
	FMSimulateTransactions() bool
}
