package txmgr

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
)

type EvmTxmConfig TxmConfig[*assets.Wei]

type TxmConfig[UNIT any] interface {
	BroadcasterConfig[UNIT]
	ConfirmerConfig[UNIT]
	ResenderConfig
	ReaperConfig

	SequenceAutoSync() bool
	UseForwarders() bool
	MaxQueuedTransactions() uint64
}

type BroadcasterConfig[UNIT any] interface {
	TriggerFallbackDBPollInterval() time.Duration
	MaxInFlightTransactions() uint32

	// from gas.Config
	IsL2() bool
	MaxFeePrice() UNIT
	FeePriceDefault() UNIT
}

type ConfirmerConfig[UNIT any] interface {
	RPCDefaultBatchSize() uint32
	UseForwarders() bool
	FeeBumpTxDepth() uint16
	MaxInFlightTransactions() uint32
	FeeLimitDefault() uint32

	// gas config
	FeeBumpThreshold() uint64
	FinalityDepth() uint32
	MaxFeePrice() UNIT
	FeeBumpPercent() uint16

	// postgres config
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

var _ EvmTxmConfig = (*evmTxmConfig)(nil)

type evmTxmConfig struct {
	Config
}

func NewEvmTxmConfig(c Config) *evmTxmConfig {
	return &evmTxmConfig{c}
}

func (c evmTxmConfig) SequenceAutoSync() bool { return c.EvmNonceAutoSync() }

func (c evmTxmConfig) UseForwarders() bool { return c.EvmUseForwarders() }

func (c evmTxmConfig) MaxQueuedTransactions() uint64 { return c.EvmMaxQueuedTransactions() }

func (c evmTxmConfig) MaxInFlightTransactions() uint32 { return c.EvmMaxInFlightTransactions() }

func (c evmTxmConfig) IsL2() bool { return c.ChainType().IsL2() }

func (c evmTxmConfig) MaxFeePrice() *assets.Wei { return c.EvmMaxGasPriceWei() }

func (c evmTxmConfig) FeePriceDefault() *assets.Wei { return c.EvmGasPriceDefault() }

func (c evmTxmConfig) RPCDefaultBatchSize() uint32 { return c.EvmRPCDefaultBatchSize() }

func (c evmTxmConfig) FeeBumpTxDepth() uint16 { return c.EvmGasBumpTxDepth() }

func (c evmTxmConfig) FeeLimitDefault() uint32 { return c.EvmGasLimitDefault() }

func (c evmTxmConfig) FeeBumpThreshold() uint64 { return c.EvmGasBumpThreshold() }

func (c evmTxmConfig) FinalityDepth() uint32 { return c.EvmFinalityDepth() }

func (c evmTxmConfig) FeeBumpPercent() uint16 { return c.EvmGasBumpPercent() }

func (c evmTxmConfig) TxResendAfterThreshold() time.Duration { return c.EthTxResendAfterThreshold() }

func (c evmTxmConfig) TxReaperInterval() time.Duration { return c.EthTxReaperInterval() }

func (c evmTxmConfig) TxReaperThreshold() time.Duration { return c.EthTxReaperThreshold() }
