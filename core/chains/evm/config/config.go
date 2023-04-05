package config

import (
	"math/big"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/config"
)

type ChainScopedOnlyConfig interface {
	evmclient.NodeConfig

	AutoCreateKey() bool
	BalanceMonitorEnabled() bool
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
	EthTxReaperInterval() time.Duration
	EthTxReaperThreshold() time.Duration
	EthTxResendAfterThreshold() time.Duration
	EvmFinalityDepth() uint32
	EvmGasBumpPercent() uint16
	EvmGasBumpThreshold() uint64
	EvmGasBumpTxDepth() uint16
	EvmGasBumpWei() *assets.Wei
	EvmGasFeeCapDefault() *assets.Wei
	EvmGasLimitDefault() uint32
	EvmGasLimitMax() uint32
	EvmGasLimitMultiplier() float32
	EvmGasLimitTransfer() uint32
	EvmGasLimitOCRJobType() *uint32
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
	EvmMaxInFlightTransactions() uint32
	EvmMaxQueuedTransactions() uint64
	EvmMinGasPriceWei() *assets.Wei
	EvmNonceAutoSync() bool
	EvmUseForwarders() bool
	EvmRPCDefaultBatchSize() uint32
	FlagsContractAddress() string
	GasEstimatorMode() string
	ChainType() config.ChainType
	KeySpecificMaxGasPriceWei(addr gethcommon.Address) *assets.Wei
	LinkContractAddress() string
	OperatorFactoryAddress() string
	MinIncomingConfirmations() uint32
	MinimumContractPayment() *assets.Link

	// OCR1 chain specific config
	OCRContractConfirmations() uint16
	OCRContractTransmitterTransmitTimeout() time.Duration
	OCRObservationGracePeriod() time.Duration
	OCRDatabaseTimeout() time.Duration

	// OCR2 chain specific config
	OCR2AutomationGasLimit() uint32
}

//go:generate mockery --quiet --name ChainScopedConfig --output ./mocks/ --case=underscore
type ChainScopedConfig interface {
	config.BasicConfig
	ChainScopedOnlyConfig
	Validate() error
}
