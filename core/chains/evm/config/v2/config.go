package v2

import (
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type Chain struct {
	BalanceMonitorEnabled *bool
	BlockBackfillDepth    *uint32
	BlockBackfillSkip     *bool

	ChainType            *string
	EIP1559DynamicFees   *bool
	FinalityDepth        *uint32
	FlagsContractAddress *ethkey.EIP55Address

	GasBumpPercent     *uint16
	GasBumpThreshold   *utils.Big
	GasBumpTxDepth     *uint16
	GasBumpWei         *utils.Big
	GasEstimatorMode   *string
	GasFeeCapDefault   *utils.Big
	GasLimitDefault    *utils.Big
	GasLimitMultiplier *decimal.Decimal
	GasLimitTransfer   *utils.Big
	GasPriceDefault    *utils.Big
	GasTipCapDefault   *utils.Big
	GasTipCapMinimum   *utils.Big

	LinkContractAddress  *ethkey.EIP55Address
	LogBackfillBatchSize *uint32
	LogPollInterval      *models.Duration

	MaxGasPriceWei           *utils.Big
	MaxInFlightTransactions  *uint32
	MaxQueuedTransactions    *uint32
	MinGasPriceWei           *utils.Big
	MinIncomingConfirmations *uint32
	MinimumContractPayment   *assets.Link

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

	BlockHistoryEstimator *BlockHistoryEstimator

	HeadTracker *HeadTracker

	KeySpecific []KeySpecific `toml:",omitempty"`

	NodePool *NodePool
}

type BlockHistoryEstimator struct {
	BatchSize                 *uint32
	BlockDelay                *uint16
	BlockHistorySize          *uint16
	EIP1559FeeCapBufferBlocks *uint16
	TransactionPercentile     *uint16
}

type KeySpecific struct {
	Key            *ethkey.EIP55Address
	MaxGasPriceWei *utils.Big
}

type HeadTracker struct {
	BlockEmissionIdleWarningThreshold *models.Duration
	HistoryDepth                      *uint32
	MaxBufferSize                     *uint32
	SamplingInterval                  *models.Duration
}

type NodePool struct {
	NoNewHeadsThreshold  *models.Duration
	PollFailureThreshold *uint32
	PollInterval         *models.Duration
}

func (c *Chain) SetFromDB(cfg *types.ChainCfg) error {
	if cfg == nil {
		return nil
	}
	if cfg.BlockHistoryEstimatorBlockDelay.Valid || cfg.BlockHistoryEstimatorBlockHistorySize.Valid || cfg.BlockHistoryEstimatorEIP1559FeeCapBufferBlocks.Valid {
		c.BlockHistoryEstimator = &BlockHistoryEstimator{}
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
		if c.HeadTracker == nil {
			c.HeadTracker = &HeadTracker{}
		}
		v := uint32(cfg.EvmHeadTrackerHistoryDepth.Int64)
		c.HeadTracker.HistoryDepth = &v
	}
	if cfg.EvmHeadTrackerMaxBufferSize.Valid {
		if c.HeadTracker == nil {
			c.HeadTracker = &HeadTracker{}
		}
		v := uint32(cfg.EvmHeadTrackerMaxBufferSize.Int64)
		c.HeadTracker.MaxBufferSize = &v
	}
	if i := cfg.EvmHeadTrackerSamplingInterval; i != nil {
		if c.HeadTracker == nil {
			c.HeadTracker = &HeadTracker{}
		}
		c.HeadTracker.SamplingInterval = cfg.EvmHeadTrackerSamplingInterval
	}
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
		c.KeySpecific = append(c.KeySpecific, KeySpecific{
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
	if cfg.NodeNoNewHeadsThreshold != nil {
		c.NodePool = &NodePool{NoNewHeadsThreshold: cfg.NodeNoNewHeadsThreshold}
	}
	return nil
}

type Node struct {
	Name     *string
	WSURL    *models.URL
	HTTPURL  *models.URL
	SendOnly *bool
}

func (n *Node) SetFromDB(db types.Node) (err error) {
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
