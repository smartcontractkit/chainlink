package types

import "time"

type TransactionManagerConfig interface {
	BroadcasterConfig
	ConfirmerConfig
	ReaperConfig
}

type TransactionManagerTransactionsConfig interface {
	BroadcasterTransactionsConfig
	ConfirmerTransactionsConfig
	ResenderTransactionsConfig
	ReaperTransactionsConfig

	ForwardersEnabled() bool
	MaxQueued() uint64
}

type BroadcasterConfig interface {
	// from gas.Config
	IsL2() bool
	MaxFeePrice() string     // logging value
	FeePriceDefault() string // logging value
}

type BroadcasterTransactionsConfig interface {
	MaxInFlight() uint32
}

type BroadcasterListenerConfig interface {
	FallbackPollInterval() time.Duration
}

type ConfirmerConfig interface {
	FeeBumpTxDepth() uint32
	FeeLimitDefault() uint32

	// from gas.Config
	FeeBumpThreshold() uint64
	FinalityDepth() uint32
	MaxFeePrice() string // logging value
	FeeBumpPercent() uint16
}

type ConfirmerChainConfig interface {
	RPCDefaultBatchSize() uint32
}

type ConfirmerDatabaseConfig interface {
	// from pg.QConfig
	DefaultQueryTimeout() time.Duration
}

type ConfirmerTransactionsConfig interface {
	MaxInFlight() uint32
	ForwardersEnabled() bool
}

type ResenderChainConfig interface {
	RPCDefaultBatchSize() uint32
}

type ResenderTransactionsConfig interface {
	ResendAfterThreshold() time.Duration
	MaxInFlight() uint32
}

//go:generate mockery --quiet --name ReaperConfig --output ./mocks/ --case=underscore

// ReaperConfig is the config subset used by the reaper
type ReaperConfig interface {
	// gas config
	FinalityDepth() uint32
}

type ReaperTransactionsConfig interface {
	ReaperInterval() time.Duration
	ReaperThreshold() time.Duration
}
