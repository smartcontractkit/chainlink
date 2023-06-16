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

	AutoCreateKey() bool
	BlockBackfillDepth() uint64
	BlockBackfillSkip() bool
	BlockEmissionIdleWarningThreshold() time.Duration
	ChainID() *big.Int
	EvmEIP1559DynamicFees() bool
	EvmFinalityDepth() uint32
	EvmGasBumpPercent() uint16
	EvmGasBumpThreshold() uint64
	EvmGasBumpTxDepth() uint32
	EvmGasBumpWei() *assets.Wei
	EvmGasFeeCapDefault() *assets.Wei
	EvmGasLimitDefault() uint32
	EvmGasLimitMax() uint32
	EvmGasLimitMultiplier() float32
	EvmGasLimitTransfer() uint32
	EvmGasLimitOCRJobType() *uint32
	EvmGasLimitOCR2JobType() *uint32
	EvmGasLimitDRJobType() *uint32
	EvmGasLimitVRFJobType() *uint32
	EvmGasLimitFMJobType() *uint32
	EvmGasLimitKeeperJobType() *uint32
	EvmGasPriceDefault() *assets.Wei
	EvmGasTipCapDefault() *assets.Wei
	EvmGasTipCapMinimum() *assets.Wei
	EvmLogBackfillBatchSize() uint32
	EvmLogKeepBlocksDepth() uint32
	EvmLogPollInterval() time.Duration
	EvmMaxGasPriceWei() *assets.Wei
	EvmMinGasPriceWei() *assets.Wei
	EvmNonceAutoSync() bool
	EvmRPCDefaultBatchSize() uint32
	FlagsContractAddress() string
	GasEstimatorMode() string
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
