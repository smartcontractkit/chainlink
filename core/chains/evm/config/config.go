package config

import (
	"math/big"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/config"
)

type EVM interface {
	HeadTracker() HeadTracker
	BalanceMonitor() BalanceMonitor
	Transactions() Transactions
	GasEstimator() GasEstimator
	OCR() OCR
	OCR2() OCR2
	NodePool() NodePool

	AutoCreateKey() bool
	BlockBackfillDepth() uint64
	BlockBackfillSkip() bool
	BlockEmissionIdleWarningThreshold() time.Duration
	ChainID() *big.Int
	ChainType() config.ChainType
	FinalityDepth() uint32
	FinalityTagEnabled() bool
	FlagsContractAddress() string
	LinkContractAddress() string
	LogBackfillBatchSize() uint32
	LogKeepBlocksDepth() uint32
	LogPollInterval() time.Duration
	MinContractPayment() *assets.Link
	MinIncomingConfirmations() uint32
	NonceAutoSync() bool
	OperatorFactoryAddress() string
	RPCDefaultBatchSize() uint32
	NodeNoNewHeadsThreshold() time.Duration

	IsEnabled() bool
	TOMLString() (string, error)
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

//go:generate mockery --quiet --name GasEstimator --output ./mocks/ --case=underscore
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
	PriceMaxKey(gethcommon.Address) *assets.Wei
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

type NodePool interface {
	PollFailureThreshold() uint32
	PollInterval() time.Duration
	SelectionMode() string
	SyncThreshold() uint32
	LeaseDuration() time.Duration
}

// TODO BCF-2509 does the chainscopedconfig really need the entire app config?
//
//go:generate mockery --quiet --name ChainScopedConfig --output ./mocks/ --case=underscore
type ChainScopedConfig interface {
	config.AppConfig
	Validate() error

	EVM() EVM
}
