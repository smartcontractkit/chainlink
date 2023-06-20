package config

import (
	"math/big"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/config"
)

// Deprecated, use EVM below
type ChainScopedOnlyConfig interface {
	evmclient.NodeConfig

	BlockEmissionIdleWarningThreshold() time.Duration
	ChainID() *big.Int
	EvmFinalityDepth() uint32
	FlagsContractAddress() string
	ChainType() config.ChainType
	KeySpecificMaxGasPriceWei(addr gethcommon.Address) *assets.Wei
	LinkContractAddress() string
	OperatorFactoryAddress() string
	MinIncomingConfirmations() uint32
	MinimumContractPayment() *assets.Link
}

type EVM interface {
	HeadTracker() HeadTracker
	BalanceMonitor() BalanceMonitor
	Transactions() Transactions
	GasEstimator() GasEstimator
	OCR() OCR
	OCR2() OCR2

	ChainType() config.ChainType
	AutoCreateKey() bool
	BlockBackfillDepth() uint64
	BlockBackfillSkip() bool
	FinalityDepth() uint32
	LogBackfillBatchSize() uint32
	LogPollInterval() time.Duration
	LogKeepBlocksDepth() uint32
	NonceAutoSync() bool
	RPCDefaultBatchSize() uint32
	KeySpecificMaxGasPriceWei(addr gethcommon.Address) *assets.Wei
}

type OCR interface {
	ContractConfirmations() uint16
	ContractTransmitterTransmitTimeout() time.Duration
	ObservationGracePeriod() time.Duration
	DatabaseTimeout() time.Duration
}

type OCR2 interface {
	Automation() OCR2Automation
}

type OCR2Automation interface {
	GasLimit() uint32
}

type HeadTracker interface {
	HistoryDepth() uint32
	MaxBufferSize() uint32
	SamplingInterval() time.Duration
}

type BalanceMonitor interface {
	Enabled() bool
}

type Transactions interface {
	ForwardersEnabled() bool
	ReaperInterval() time.Duration
	ResendAfterThreshold() time.Duration
	ReaperThreshold() time.Duration
	MaxInFlight() uint32
	MaxQueued() uint64
}

type GasEstimator interface {
	BlockHistory() BlockHistory
	LimitJobType() LimitJobType

	EIP1559DynamicFees() bool
	BumpPercent() uint16
	BumpThreshold() uint64
	BumpTxDepth() uint32
	BumpMin() *assets.Wei
	FeeCapDefault() *assets.Wei
	LimitDefault() uint32
	LimitMax() uint32
	LimitMultiplier() float32
	LimitTransfer() uint32
	PriceDefault() *assets.Wei
	TipCapDefault() *assets.Wei
	TipCapMin() *assets.Wei
	PriceMax() *assets.Wei
	PriceMin() *assets.Wei
	Mode() string
}

type LimitJobType interface {
	OCR() *uint32
	OCR2() *uint32
	DR() *uint32
	FM() *uint32
	Keeper() *uint32
	VRF() *uint32
}

type BlockHistory interface {
	BatchSize() uint32
	BlockHistorySize() uint16
	BlockDelay() uint16
	CheckInclusionBlocks() uint16
	CheckInclusionPercentile() uint16
	EIP1559FeeCapBufferBlocks() uint16
	TransactionPercentile() uint16
}

//go:generate mockery --quiet --name ChainScopedConfig --output ./mocks/ --case=underscore
type ChainScopedConfig interface {
	config.AppConfig
	ChainScopedOnlyConfig // Deprecated, to be replaced by EVM() below
	Validate() error

	EVM() EVM
}
