package types

import "time"

// FEE_UNIT - fee unit
type TxmConfig[FEE_UNIT Unit] interface {
	BroadcasterConfig[FEE_UNIT]
	ConfirmerConfig[FEE_UNIT]
	ResenderConfig
	ReaperConfig

	SequenceAutoSync() bool
	UseForwarders() bool
	MaxQueuedTransactions() uint64
}

// FEE_UNIT - fee unit
type BroadcasterConfig[FEE_UNIT Unit] interface {
	TriggerFallbackDBPollInterval() time.Duration
	MaxInFlightTransactions() uint32

	// from gas.Config
	IsL2() bool
	MaxFeePrice() FEE_UNIT
	FeePriceDefault() FEE_UNIT
}

// FEE_UNIT - fee unit
type ConfirmerConfig[FEE_UNIT Unit] interface {
	RPCDefaultBatchSize() uint32
	UseForwarders() bool
	FeeBumpTxDepth() uint32
	MaxInFlightTransactions() uint32
	FeeLimitDefault() uint32

	// from gas.Config
	FeeBumpThreshold() uint64
	FinalityDepth() uint32
	MaxFeePrice() FEE_UNIT
	FeeBumpPercent() uint16

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
