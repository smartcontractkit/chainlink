package types

import "time"

// UNIT - fee unit
type TxmConfig[UNIT any] interface {
	BroadcasterConfig[UNIT]
	ConfirmerConfig[UNIT]
	ResenderConfig
	ReaperConfig

	SequenceAutoSync() bool
	UseForwarders() bool
	MaxQueuedTransactions() uint64
}

// UNIT - fee unit
type BroadcasterConfig[UNIT any] interface {
	TriggerFallbackDBPollInterval() time.Duration
	MaxInFlightTransactions() uint32

	// from gas.Config
	IsL2() bool
	MaxFeePrice() UNIT
	FeePriceDefault() UNIT
}

// UNIT - fee unit
type ConfirmerConfig[UNIT any] interface {
	RPCDefaultBatchSize() uint32
	UseForwarders() bool
	FeeBumpTxDepth() uint16
	MaxInFlightTransactions() uint32
	FeeLimitDefault() uint32

	// from gas.Config
	FeeBumpThreshold() uint64
	FinalityDepth() uint32
	MaxFeePrice() UNIT
	FeeBumpPercent() uint16

	// from pg.QConfig
	DatabaseDefaultQueryTimeout() time.Duration
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
