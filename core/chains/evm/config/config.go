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
	BlockHistoryEstimatorBatchSize() (size uint32)
	BlockHistoryEstimatorBlockDelay() uint16
	BlockHistoryEstimatorBlockHistorySize() uint16
	BlockHistoryEstimatorCheckInclusionBlocks() uint16
	BlockHistoryEstimatorCheckInclusionPercentile() uint16
	BlockHistoryEstimatorEIP1559FeeCapBufferBlocks() uint16
	BlockHistoryEstimatorTransactionPercentile() uint16
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
	EvmHeadTrackerHistoryDepth() uint32
	EvmHeadTrackerMaxBufferSize() uint32
	EvmHeadTrackerSamplingInterval() time.Duration
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

	// OCR2 chain specific config
	OCR2AutomationGasLimit() uint32
}

type EVM interface {
	BalanceMonitor() BalanceMonitor
	Transactions() Transactions
	OCR() OCR
}

type OCR interface {
	ContractConfirmations() uint16
	ContractTransmitterTransmitTimeout() time.Duration
	ObservationGracePeriod() time.Duration
	DatabaseTimeout() time.Duration
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

//go:generate mockery --quiet --name ChainScopedConfig --output ./mocks/ --case=underscore
type ChainScopedConfig interface {
	config.AppConfig
	ChainScopedOnlyConfig // Deprecated, to be replaced by EVM() below
	Validate() error

	EVM() EVM
}
