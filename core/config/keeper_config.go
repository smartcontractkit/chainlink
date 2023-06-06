package config

import "time"

type Registry interface {
	CheckGasOverhead() uint32
	PerformGasOverhead() uint32
	MaxPerformDataSize() uint32
	SyncInterval() time.Duration
	SyncUpkeepQueueSize() uint32
}

type Keeper interface {
	DefaultTransactionQueueDepth() uint32
	GasPriceBufferPercent() uint16
	GasTipCapBufferPercent() uint16
	BaseFeeBufferPercent() uint16
	MaxGracePeriod() int64
	TurnLookBack() int64
	Registry() Registry
}
