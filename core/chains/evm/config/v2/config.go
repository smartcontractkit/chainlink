package v2

import (
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	v2 "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type Chain struct {
	BlockBackfillDepth       *uint32
	BlockBackfillSkip        *bool
	ChainType                *string
	FinalityDepth            *uint32
	FlagsContractAddress     *ethkey.EIP55Address
	LinkContractAddress      *ethkey.EIP55Address
	LogBackfillBatchSize     *uint32
	LogPollInterval          *models.Duration
	MaxInFlightTransactions  *uint32
	MaxQueuedTransactions    *uint32
	MinIncomingConfirmations *uint32
	MinimumContractPayment   *assets.Link
	NonceAutoSync            *bool
	OperatorFactoryAddress   *ethkey.EIP55Address
	RPCDefaultBatchSize      *uint32
	TxReaperInterval         *models.Duration
	TxReaperThreshold        *models.Duration
	TxResendAfterThreshold   *models.Duration

	UseForwarders *bool

	BalanceMonitor *BalanceMonitor

	GasEstimator *GasEstimator

	HeadTracker *HeadTracker

	KeySpecific KeySpecificConfig `toml:",omitempty"`

	NodePool *NodePool

	OCR *OCR
}

type BalanceMonitor struct {
	Enabled    *bool
	BlockDelay *uint16
}

type GasEstimator struct {
	Mode *string

	PriceDefault *utils.Wei
	PriceMax     *utils.Wei
	PriceMin     *utils.Wei

	LimitDefault    *uint32
	LimitMultiplier *decimal.Decimal
	LimitTransfer   *uint32

	LimitOCRJobType    *uint32
	LimitDRJobType     *uint32
	LimitVRFJobType    *uint32
	LimitFMJobType     *uint32
	LimitKeeperJobType *uint32

	BumpMin       *utils.Wei
	BumpPercent   *uint16
	BumpThreshold *uint32
	BumpTxDepth   *uint16

	EIP1559DynamicFees *bool

	FeeCapDefault *utils.Wei
	TipCapDefault *utils.Wei
	TipCapMinimum *utils.Wei

	BlockHistory *BlockHistoryEstimator
}

type BlockHistoryEstimator struct {
	BatchSize                 *uint32
	BlockDelay                *uint16
	BlockHistorySize          *uint16
	EIP1559FeeCapBufferBlocks *uint16
	TransactionPercentile     *uint16
}

type KeySpecificConfig []KeySpecific

func (ks KeySpecificConfig) ValidateConfig() (err error) {
	addrs := map[string]struct{}{}
	for _, k := range ks {
		addr := k.Key.String()
		if _, ok := addrs[addr]; ok {
			err = multierr.Append(err, fmt.Errorf("duplicate address: %s", addr))
		} else {
			addrs[addr] = struct{}{}
		}
	}
	return
}

type KeySpecific struct {
	Key          *ethkey.EIP55Address
	GasEstimator *KeySpecificGasEstimator
}

type KeySpecificGasEstimator struct {
	PriceMax *utils.Wei
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

type OCR struct {
	ContractConfirmations              *uint16
	ContractTransmitterTransmitTimeout *models.Duration
	DatabaseTimeout                    *models.Duration
	ObservationTimeout                 *models.Duration
	ObservationGracePeriod             *models.Duration
}

func (c *Chain) SetFromDB(cfg *types.ChainCfg) error {
	if cfg == nil {
		return nil
	}
	if cfg.ChainType.Valid {
		c.ChainType = &cfg.ChainType.String
	}
	c.TxReaperThreshold = cfg.EthTxReaperThreshold
	c.TxResendAfterThreshold = cfg.EthTxResendAfterThreshold
	if cfg.EvmFinalityDepth.Valid {
		v := uint32(cfg.EvmFinalityDepth.Int64)
		c.FinalityDepth = &v
	}
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
			return errors.Errorf("invalid FlagsContractAddress: %s", s)
		}
		a := common.HexToAddress(s)
		v := ethkey.EIP55AddressFromAddress(a)
		c.FlagsContractAddress = &v
	}
	if cfg.GasEstimatorMode.Valid {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		c.GasEstimator.Mode = &cfg.GasEstimatorMode.String
	}
	if cfg.EvmMaxGasPriceWei != nil {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		c.GasEstimator.PriceMax = cfg.EvmMaxGasPriceWei.Wei()
	}
	if cfg.EvmEIP1559DynamicFees.Valid {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		c.GasEstimator.EIP1559DynamicFees = &cfg.EvmEIP1559DynamicFees.Bool
	}
	if cfg.EvmGasPriceDefault != nil {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		c.GasEstimator.PriceDefault = cfg.EvmGasPriceDefault.Wei()
	}
	if cfg.EvmGasLimitMultiplier.Valid {
		v := decimal.NewFromFloat(cfg.EvmGasLimitMultiplier.Float64)
		c.GasEstimator.LimitMultiplier = &v
	}
	if cfg.EvmGasTipCapDefault != nil {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		c.GasEstimator.TipCapDefault = cfg.EvmGasTipCapDefault.Wei()
	}
	if cfg.EvmGasTipCapMinimum != nil {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		c.GasEstimator.TipCapMinimum = cfg.EvmGasTipCapMinimum.Wei()
	}
	if cfg.EvmGasBumpPercent.Valid {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		v := uint16(cfg.EvmGasBumpPercent.Int64)
		c.GasEstimator.BumpPercent = &v
	}
	if cfg.EvmGasBumpTxDepth.Valid {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		v := uint16(cfg.EvmGasBumpTxDepth.Int64)
		c.GasEstimator.BumpTxDepth = &v
	}
	if cfg.EvmGasBumpWei != nil {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		c.GasEstimator.BumpMin = cfg.EvmGasBumpWei.Wei()
	}
	if cfg.EvmGasFeeCapDefault != nil {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		c.GasEstimator.FeeCapDefault = cfg.EvmGasFeeCapDefault.Wei()
	}
	if cfg.EvmGasLimitDefault.Valid {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		v := uint32(cfg.EvmGasLimitDefault.Int64)
		c.GasEstimator.LimitDefault = &v
	}
	if cfg.EvmGasLimitOCRJobType.Valid {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		v := uint32(cfg.EvmGasLimitOCRJobType.Int64)
		c.GasEstimator.LimitOCRJobType = &v
	}
	if cfg.EvmGasLimitDRJobType.Valid {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		v := uint32(cfg.EvmGasLimitDRJobType.Int64)
		c.GasEstimator.LimitDRJobType = &v
	}
	if cfg.EvmGasLimitVRFJobType.Valid {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		v := uint32(cfg.EvmGasLimitVRFJobType.Int64)
		c.GasEstimator.LimitVRFJobType = &v
	}
	if cfg.EvmGasLimitFMJobType.Valid {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		v := uint32(cfg.EvmGasLimitFMJobType.Int64)
		c.GasEstimator.LimitFMJobType = &v
	}
	if cfg.EvmGasLimitKeeperJobType.Valid {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		v := uint32(cfg.EvmGasLimitKeeperJobType.Int64)
		c.GasEstimator.LimitKeeperJobType = &v
	}

	if cfg.BlockHistoryEstimatorBlockDelay.Valid || cfg.BlockHistoryEstimatorBlockHistorySize.Valid || cfg.BlockHistoryEstimatorEIP1559FeeCapBufferBlocks.Valid {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		c.GasEstimator.BlockHistory = &BlockHistoryEstimator{}
		if cfg.BlockHistoryEstimatorBlockDelay.Valid {
			v := uint16(cfg.BlockHistoryEstimatorBlockDelay.Int64)
			c.GasEstimator.BlockHistory.BlockDelay = &v
		}
		if cfg.BlockHistoryEstimatorBlockHistorySize.Valid {
			v := uint16(cfg.BlockHistoryEstimatorBlockHistorySize.Int64)
			c.GasEstimator.BlockHistory.BlockHistorySize = &v
		}
		if cfg.BlockHistoryEstimatorEIP1559FeeCapBufferBlocks.Valid {
			v := uint16(cfg.BlockHistoryEstimatorEIP1559FeeCapBufferBlocks.Int64)
			c.GasEstimator.BlockHistory.EIP1559FeeCapBufferBlocks = &v
		}
	}
	for s, kcfg := range cfg.KeySpecific {
		if !common.IsHexAddress(s) {
			return errors.Errorf("invalid address KeySpecific: %s", s)
		}
		a := common.HexToAddress(s)
		v := ethkey.EIP55AddressFromAddress(a)
		c.KeySpecific = append(c.KeySpecific, KeySpecific{
			Key: &v,
			GasEstimator: &KeySpecificGasEstimator{
				PriceMax: kcfg.EvmMaxGasPriceWei.Wei(),
			},
		})
	}
	if cfg.LinkContractAddress.Valid {
		s := cfg.LinkContractAddress.String
		if !common.IsHexAddress(s) {
			return errors.Errorf("invalid LinkContractAddress: %s", s)
		}
		a := common.HexToAddress(s)
		v := ethkey.EIP55AddressFromAddress(a)
		c.LinkContractAddress = &v
	}
	if cfg.OperatorFactoryAddress.Valid {
		s := cfg.OperatorFactoryAddress.String
		if !common.IsHexAddress(s) {
			return errors.Errorf("invalid OperatorFactoryAddress: %s", s)
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
	if cfg.OCRObservationTimeout != nil {
		c.OCR = &OCR{ObservationTimeout: cfg.OCRObservationTimeout}
	}
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

func (n *Node) ValidateConfig() (err error) {
	if n.Name == nil {
		err = multierr.Append(err, v2.ErrMissing{Name: "Name", Msg: "required for all nodes"})
	} else if *n.Name == "" {
		err = multierr.Append(err, v2.ErrEmpty{Name: "Name", Msg: "required for all nodes"})
	}
	if s := n.SendOnly; s != nil && *s {
		if n.WSURL == nil {
			err = multierr.Append(err, v2.ErrMissing{Name: "WSURL", Msg: "required for SendOnly nodes"})
		}
	}
	if n.HTTPURL == nil {
		err = multierr.Append(err, v2.ErrMissing{Name: "HTTPURL", Msg: "required for all nodes"})
	}
	return
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
