package types

import "time"

type TxmConfig interface {
	BroadcasterConfig
	ConfirmerConfig
	ResenderConfig
	ReaperConfig

	SequenceAutoSync() bool
	UseForwarders() bool
	MaxQueuedTransactions() uint64
}

type BroadcasterConfig interface {
	MaxInFlightTransactions() uint32

	// from gas.Config
	IsL2() bool
	MaxFeePrice() string     // logging value
	FeePriceDefault() string // logging value
}

type BroadcasterListenerConfig interface {
	FallbackPollInterval() time.Duration
}

type ConfirmerConfig interface {
	RPCDefaultBatchSize() uint32
	UseForwarders() bool
	FeeBumpTxDepth() uint32
	MaxInFlightTransactions() uint32
	FeeLimitDefault() uint32

	// from gas.Config
	FeeBumpThreshold() uint64
	FinalityDepth() uint32
	MaxFeePrice() string // logging value
	FeeBumpPercent() uint16
}

type ConfirmerDatabaseConfig interface {
	// from pg.QConfig
	DefaultQueryTimeout() time.Duration
}

type ResenderConfig interface {
	TxResendAfterThreshold() time.Duration
	MaxInFlightTransactions() uint32
	RPCDefaultBatchSize() uint32
}

//go:generate mockery --quiet --name ReaperConfig --output ./mocks/ --case=underscore

// ReaperConfig is the config subset used by the reaper
type ReaperConfig interface {
	TxReaperInterval() time.Duration
	TxReaperThreshold() time.Duration

	// gas config
	FinalityDepth() uint32
}
