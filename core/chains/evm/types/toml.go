package types

import (
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/core/assets"
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
	OCRObservationTimeout                 *models.Duration
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
}

func (c *ChainTOMLCfg) SetFromDB(cfg *ChainCfg) error {
	if cfg == nil {
		return nil
	}
	if cfg.BlockHistoryEstimatorBlockDelay.Valid || cfg.BlockHistoryEstimatorBlockHistorySize.Valid || cfg.BlockHistoryEstimatorEIP1559FeeCapBufferBlocks.Valid {
		c.BlockHistoryEstimator = &BlockHistoryEstimatorConfig{}
		if cfg.BlockHistoryEstimatorBlockDelay.Valid {
			v := uint16(cfg.BlockHistoryEstimatorBlockDelay.Int64)
			c.BlockHistoryEstimator.BlockDelay = &v
		}
		if cfg.BlockHistoryEstimatorBlockHistorySize.Valid {
			v := uint16(cfg.BlockHistoryEstimatorBlockHistorySize.Int64)
			c.BlockHistoryEstimator.BlockHistorySize = &v
		}
		if cfg.BlockHistoryEstimatorEIP1559FeeCapBufferBlocks.Valid {
			v := uint16(cfg.BlockHistoryEstimatorEIP1559FeeCapBufferBlocks.Int64)
			c.BlockHistoryEstimator.EIP1559FeeCapBufferBlocks = &v
		}
	}
	if cfg.ChainType.Valid {
		c.ChainType = &cfg.ChainType.String
	}
	c.TxReaperThreshold = cfg.EthTxReaperThreshold
	c.TxResendAfterThreshold = cfg.EthTxResendAfterThreshold
	if cfg.EvmEIP1559DynamicFees.Valid {
		c.EIP1559DynamicFees = &cfg.EvmEIP1559DynamicFees.Bool
	}
	if cfg.EvmFinalityDepth.Valid {
		v := uint32(cfg.EvmFinalityDepth.Int64)
		c.FinalityDepth = &v
	}
	if cfg.EvmGasBumpPercent.Valid {
		v := uint16(cfg.EvmGasBumpPercent.Int64)
		c.GasBumpPercent = &v
	}
	if cfg.EvmGasBumpTxDepth.Valid {
		v := uint16(cfg.EvmGasBumpTxDepth.Int64)
		c.GasBumpTxDepth = &v
	}
	c.GasBumpWei = cfg.EvmGasBumpWei
	c.GasFeeCapDefault = cfg.EvmGasFeeCapDefault
	if cfg.EvmGasLimitDefault.Valid {
		c.GasLimitDefault = utils.NewBigI(cfg.EvmGasLimitDefault.Int64)
	}
	if cfg.EvmGasLimitMultiplier.Valid {
		v := decimal.NewFromFloat(cfg.EvmGasLimitMultiplier.Float64)
		c.GasLimitMultiplier = &v
	}
	c.GasPriceDefault = cfg.EvmGasPriceDefault
	c.GasTipCapDefault = cfg.EvmGasTipCapDefault
	c.GasTipCapMinimum = cfg.EvmGasTipCapMinimum
	if cfg.EvmHeadTrackerHistoryDepth.Valid {
		v := uint32(cfg.EvmHeadTrackerHistoryDepth.Int64)
		c.HeadTrackerHistoryDepth = &v
	}
	if cfg.EvmHeadTrackerMaxBufferSize.Valid {
		v := uint32(cfg.EvmHeadTrackerMaxBufferSize.Int64)
		c.HeadTrackerMaxBufferSize = &v
	}
	c.HeadTrackerSamplingInterval = cfg.EvmHeadTrackerSamplingInterval
	if cfg.EvmLogBackfillBatchSize.Valid {
		v := uint32(cfg.EvmLogBackfillBatchSize.Int64)
		c.LogBackfillBatchSize = &v
	}
	c.LogPollInterval = cfg.EvmLogPollInterval
	c.MaxGasPriceWei = cfg.EvmMaxGasPriceWei
	if cfg.EvmNonceAutoSync.Valid {
		c.NonceAutoSync = &cfg.EvmNonceAutoSync.Bool
	}
	if cfg.EvmUseForwarders.Valid {
		c.UseForwarders = &cfg.EvmUseForwarders.Bool
	}
	if cfg.EvmRPCDefaultBatchSize.Valid {
		v := uint32(cfg.EvmRPCDefaultBatchSize.Int64)
		c.RPCDefaultBatchSize = &v
	}
	if cfg.FlagsContractAddress.Valid {
		s := cfg.FlagsContractAddress.String
		if !common.IsHexAddress(s) {
			return fmt.Errorf("invalid FlagsContractAddress: %s", s)
		}
		a := common.HexToAddress(s)
		v := ethkey.EIP55AddressFromAddress(a)
		c.FlagsContractAddress = &v
	}
	if cfg.GasEstimatorMode.Valid {
		c.GasEstimatorMode = &cfg.GasEstimatorMode.String
	}
	for s, kcfg := range cfg.KeySpecific {
		if !common.IsHexAddress(s) {
			return fmt.Errorf("invalid address KeySpecific: %s", s)
		}
		a := common.HexToAddress(s)
		v := ethkey.EIP55AddressFromAddress(a)
		c.KeySpecific = append(c.KeySpecific, KeySpecificConfig{
			Key:            &v,
			MaxGasPriceWei: kcfg.EvmMaxGasPriceWei,
		})
	}
	if cfg.LinkContractAddress.Valid {
		s := cfg.LinkContractAddress.String
		if !common.IsHexAddress(s) {
			return fmt.Errorf("invalid LinkContractAddress: %s", s)
		}
		a := common.HexToAddress(s)
		v := ethkey.EIP55AddressFromAddress(a)
		c.LinkContractAddress = &v
	}
	if cfg.OperatorFactoryAddress.Valid {
		s := cfg.OperatorFactoryAddress.String
		if !common.IsHexAddress(s) {
			return fmt.Errorf("invalid OperatorFactoryAddress: %s", s)
		}
		a := common.HexToAddress(s)
		v := ethkey.EIP55AddressFromAddress(a)
		c.OperatorFactoryAddress = &v
	}
	if cfg.MinIncomingConfirmations.Valid {
		v := uint32(cfg.MinIncomingConfirmations.Int64)
		c.MinIncomingConfirmations = &v
	}
	c.MinimumContractPayment = cfg.MinimumContractPayment
	c.OCRObservationTimeout = cfg.OCRObservationTimeout
	c.NodeNoNewHeadsThreshold = cfg.NodeNoNewHeadsThreshold
	return nil
}

type TOMLNode struct {
	Name     *string
	WSURL    *models.URL
	HTTPURL  *models.URL
	SendOnly *bool
}

func NewTOMLNodeFromDB(db Node) (n TOMLNode, err error) {
	n.Name = &db.Name
	if db.WSURL.Valid {
		var u *url.URL
		u, err = url.Parse(db.WSURL.String)
		if err != nil {
			return
		}
		n.WSURL = (*models.URL)(u)
	}
	if db.HTTPURL.Valid {
		var u *url.URL
		u, err = url.Parse(db.HTTPURL.String)
		if err != nil {
			return
		}
		n.HTTPURL = (*models.URL)(u)
	}
	if db.SendOnly {
		// Only necessary if true
		n.SendOnly = &db.SendOnly
	}
	return
}
