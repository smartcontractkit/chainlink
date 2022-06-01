package types

import (
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/config/toml"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type ChainTOMLCfg struct {
	BalanceMonitorEnabled             *bool
	BlockBackfillDepth                *uint32
	BlockBackfillSkip                 *bool
	BlockEmissionIdleWarningThreshold *models.Duration

	BlockHistoryEstimator *BlockHistoryEstimatorConfig

	ChainType            *string
	EIP1559DynamicFees   *bool
	FinalityDepth        *uint32
	FlagsContractAddress *ethkey.EIP55Address

	GasBumpPercent     *uint16
	GasBumpThreshold   *utils.Big
	GasBumpTxDepth     *uint16
	GasBumpWei         *utils.Big
	GasEstimatorMode   *string //TODO strict type?
	GasFeeCapDefault   *utils.Big
	GasLimitDefault    *utils.Big
	GasLimitMultiplier *decimal.Decimal
	GasLimitTransfer   *utils.Big
	GasPriceDefault    *utils.Big
	GasTipCapDefault   *utils.Big
	GasTipCapMinimum   *utils.Big

	HeadTrackerHistoryDepth     *uint32
	HeadTrackerMaxBufferSize    *uint32
	HeadTrackerSamplingInterval *models.Duration

	KeySpecific []KeySpecificConfig `toml:",omitempty"`

	LinkContractAddress  *ethkey.EIP55Address
	LogBackfillBatchSize *uint32
	LogPollInterval      *models.Duration

	MaxGasPriceWei           *utils.Big
	MaxInFlightTransactions  *uint32
	MaxQueuedTransactions    *uint32
	MinGasPriceWei           *utils.Big
	MinIncomingConfirmations *uint32
	MinimumContractPayment   *assets.Link

	NodeNoNewHeadsThreshold  *models.Duration
	NodePollFailureThreshold *uint32
	NodePollInterval         *models.Duration

	NonceAutoSync *bool

	OCRContractConfirmations              *uint16
	OCRContractTransmitterTransmitTimeout *models.Duration
	OCRDatabaseTimeout                    *models.Duration
	OCRObservationGracePeriod             *models.Duration
	OCR2ContractConfirmations             *uint16

	OperatorFactoryAddress *ethkey.EIP55Address
	RPCDefaultBatchSize    *uint32

	TxReaperInterval       *models.Duration
	TxReaperThreshold      *models.Duration
	TxResendAfterThreshold *models.Duration

	UseForwarders *bool
}

type BlockHistoryEstimatorConfig struct {
	BatchSize                 *uint32
	BlockDelay                *uint16
	BlockHistorySize          *uint16
	EIP1559FeeCapBufferBlocks *uint16
	TransactionPercentile     *uint16
}

type KeySpecificConfig struct {
	Key            *ethkey.EIP55Address
	MaxGasPriceWei *utils.Big
	//TODO more?
}

type TOMLNode struct {
	Name     *string
	WSURL    *toml.URL
	HTTPURL  *toml.URL
	SendOnly *bool
}
