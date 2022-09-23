package v2

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"
	"golang.org/x/exp/constraints"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/config"
	v2 "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type HasEVMConfigs interface {
	EVMConfigs() EVMConfigs
}

type EVMConfigs []*EVMConfig

func (cs EVMConfigs) ValidateConfig() (err error) {
	// Unique chain IDs
	chainIDs := map[string]struct{}{}
	for i, c := range cs {
		if c.ChainID == nil {
			continue
		}
		chainID := c.ChainID.String()
		if chainID == "" {
			continue
		}
		if _, ok := chainIDs[chainID]; ok {
			err = multierr.Append(err, v2.ErrInvalid{Name: fmt.Sprintf("%d.ChainID", i), Msg: "duplicate - must be unique", Value: chainID})
		} else {
			chainIDs[chainID] = struct{}{}
		}
	}

	// Unique node names
	names := map[string]struct{}{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			if n.Name == nil || *n.Name == "" {
				continue
			}
			if _, ok := names[*n.Name]; ok {
				err = multierr.Append(err, v2.ErrInvalid{Name: fmt.Sprintf("%d.Nodes.%d.Name", i, j), Msg: "duplicate - must be unique", Value: *n.Name})
			}
			names[*n.Name] = struct{}{}
		}
	}

	// Unique node WSURLs
	wsURLs := map[string]struct{}{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			if n.WSURL == nil {
				continue
			}
			us := (*url.URL)(n.WSURL).String()
			if _, ok := wsURLs[us]; ok {
				err = multierr.Append(err, v2.ErrInvalid{Name: fmt.Sprintf("%d.Nodes.%d.WSURL", i, j), Msg: "duplicate - must be unique", Value: us})
			}
			wsURLs[us] = struct{}{}
		}
	}

	// Unique node HTTPURLs
	httpURLs := map[string]struct{}{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			if n.HTTPURL == nil {
				continue
			}
			us := (*url.URL)(n.HTTPURL).String()
			if _, ok := httpURLs[us]; ok {
				err = multierr.Append(err, v2.ErrInvalid{Name: fmt.Sprintf("%d.Nodes.%d.HTTPURL", i, j), Msg: "duplicate - must be unique", Value: us})
			}
			httpURLs[us] = struct{}{}
		}
	}
	return
}

func (cs EVMConfigs) Chains(ids ...utils.Big) (chains []types.DBChain) {
	for _, ch := range cs {
		if ch == nil {
			continue
		}
		var match bool
		for _, id := range ids {
			if id.Cmp(ch.ChainID) == 0 {
				match = true
				break
			}
		}
		if !match {
			continue
		}
		chains = append(chains, types.DBChain{
			ID:      *ch.ChainID,
			Enabled: *ch.Enabled,
			Cfg:     ch.asV1(),
		})
	}
	return
}

func (cs EVMConfigs) Node(name string) (types.Node, error) {
	for i := range cs {
		for _, n := range cs[i].Nodes {
			if n.Name != nil && *n.Name == name {
				return legacyNode(n, cs[i].ChainID), nil
			}
		}
	}
	return types.Node{}, sql.ErrNoRows
}

func legacyNode(n *Node, chainID *utils.Big) types.Node {
	return types.Node{
		Name:       *n.Name,
		EVMChainID: *chainID,
		WSURL:      null.StringFrom((*n).WSURL.String()),
		HTTPURL:    null.StringFrom((*n).HTTPURL.String()),
		SendOnly:   *n.SendOnly,
	}
}

func (cs EVMConfigs) Nodes() (ns []types.Node) {
	for i := range cs {
		for _, n := range cs[i].Nodes {
			if n == nil {
				continue
			}
			ns = append(ns, legacyNode(n, cs[i].ChainID))
		}
	}
	return
}

func (cs EVMConfigs) NodesByID(chainIDs ...utils.Big) (ns []types.Node) {
	for i := range cs {
		var match bool
		for _, chainID := range chainIDs {
			if chainID.Cmp(cs[i].ChainID) != 0 {
				match = true
				break
			}
		}
		if !match {
			continue
		}
		for _, n := range cs[i].Nodes {
			if n == nil {
				continue
			}
			ns = append(ns, legacyNode(n, cs[i].ChainID))
		}
	}
	return
}

type EVMNodes []*Node

type EVMConfig struct {
	ChainID *utils.Big
	Enabled *bool
	Chain
	Nodes EVMNodes
}

// Ensure that the embedded struct will be validated (w/o requiring a pointer receiver).
var _ v2.Validated = Chain{}

func (c *EVMConfig) SetFromDB(ch types.DBChain, nodes []types.Node) error {
	c.ChainID = &ch.ID
	c.Enabled = &ch.Enabled

	if err := c.Chain.SetFromDB(ch.Cfg); err != nil {
		return err
	}
	for _, db := range nodes {
		var n Node
		if err := n.SetFromDB(db); err != nil {
			return err
		}
		c.Nodes = append(c.Nodes, &n)
	}
	return nil
}

func (c *EVMConfig) ValidateConfig() (err error) {
	if c.ChainID == nil {
		err = multierr.Append(err, v2.ErrMissing{Name: "ChainID", Msg: "required for all chains"})
	} else if c.ChainID.String() == "" {
		err = multierr.Append(err, v2.ErrEmpty{Name: "ChainID", Msg: "required for all chains"})
	}
	//TODO more from chain scoped?

	return
}

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
	NoNewHeadsThreshold      *models.Duration
	OperatorFactoryAddress   *ethkey.EIP55Address
	RPCDefaultBatchSize      *uint32
	RPCBlockQueryDelay       *uint16
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

func (c Chain) ValidateConfig() (err error) {
	if c.ChainType != nil && !config.ChainType(*c.ChainType).IsValid() {
		err = multierr.Append(err, v2.ErrInvalid{Name: "ChainType", Value: *c.ChainType, Msg: config.ErrInvalidChainType.Error()})
	}
	//TODO more from chain scoped?
	return
}

func (c *Chain) asV1() *types.ChainCfg {
	cfg := types.ChainCfg{
		BlockHistoryEstimatorBlockDelay:                null.Int{},
		BlockHistoryEstimatorBlockHistorySize:          null.Int{},
		BlockHistoryEstimatorEIP1559FeeCapBufferBlocks: null.Int{},
		ChainType:                      null.StringFromPtr(c.ChainType),
		EthTxReaperThreshold:           c.TxReaperThreshold,
		EthTxResendAfterThreshold:      c.TxResendAfterThreshold,
		EvmEIP1559DynamicFees:          null.BoolFromPtr(c.GasEstimator.EIP1559DynamicFees),
		EvmFinalityDepth:               nullInt(c.FinalityDepth),
		EvmGasBumpPercent:              nullInt(c.GasEstimator.BumpPercent),
		EvmGasBumpTxDepth:              nullInt(c.GasEstimator.BumpTxDepth),
		EvmGasBumpWei:                  (*utils.Big)(c.GasEstimator.BumpMin),
		EvmGasFeeCapDefault:            (*utils.Big)(c.GasEstimator.FeeCapDefault),
		EvmGasLimitDefault:             nullInt(c.GasEstimator.LimitDefault),
		EvmGasLimitMax:                 nullInt(c.GasEstimator.LimitMax),
		EvmGasLimitMultiplier:          nullFloat(c.GasEstimator.LimitMultiplier),
		EvmGasLimitOCRJobType:          nullInt(c.GasEstimator.LimitOCRJobType),
		EvmGasLimitDRJobType:           nullInt(c.GasEstimator.LimitDRJobType),
		EvmGasLimitVRFJobType:          nullInt(c.GasEstimator.LimitVRFJobType),
		EvmGasLimitFMJobType:           nullInt(c.GasEstimator.LimitFMJobType),
		EvmGasLimitKeeperJobType:       nullInt(c.GasEstimator.LimitKeeperJobType),
		EvmGasPriceDefault:             (*utils.Big)(c.GasEstimator.PriceDefault),
		EvmGasTipCapDefault:            (*utils.Big)(c.GasEstimator.TipCapDefault),
		EvmGasTipCapMinimum:            (*utils.Big)(c.GasEstimator.TipCapMinimum),
		EvmHeadTrackerHistoryDepth:     nullInt(c.HeadTracker.HistoryDepth),
		EvmHeadTrackerMaxBufferSize:    nullInt(c.HeadTracker.MaxBufferSize),
		EvmHeadTrackerSamplingInterval: c.HeadTracker.SamplingInterval,
		EvmLogBackfillBatchSize:        nullInt(c.LogBackfillBatchSize),
		EvmLogPollInterval:             c.LogPollInterval,
		EvmMaxGasPriceWei:              (*utils.Big)(c.GasEstimator.PriceMax),
		EvmNonceAutoSync:               null.BoolFromPtr(c.NonceAutoSync),
		EvmUseForwarders:               null.BoolFromPtr(c.UseForwarders),
		EvmRPCDefaultBatchSize:         nullInt(c.RPCDefaultBatchSize),
		FlagsContractAddress:           nullString(c.FlagsContractAddress),
		GasEstimatorMode:               null.StringFromPtr(c.GasEstimator.Mode),
		LinkContractAddress:            nullString(c.LinkContractAddress),
		OperatorFactoryAddress:         nullString(c.OperatorFactoryAddress),
		MinIncomingConfirmations:       nullInt(c.MinIncomingConfirmations),
		MinimumContractPayment:         c.MinimumContractPayment,
		NodeNoNewHeadsThreshold:        c.NoNewHeadsThreshold,
	}
	for _, ks := range c.KeySpecific {
		if cfg.KeySpecific == nil {
			cfg.KeySpecific = map[string]types.ChainCfg{}
		}
		cfg.KeySpecific[ks.Key.String()] = types.ChainCfg{
			EvmMaxGasPriceWei: (*utils.Big)(ks.GasEstimator.PriceMax),
		}
	}
	return &cfg
}

func nullInt[I constraints.Integer](i *I) null.Int {
	if i == nil {
		return null.Int{}
	}
	return null.IntFrom(int64(*i))
}

func nullFloat(d *decimal.Decimal) null.Float {
	if d == nil {
		return null.Float{}
	}
	return null.FloatFrom(d.InexactFloat64())
}

func nullString[S fmt.Stringer](s *S) null.String {
	if s == nil {
		return null.String{}
	}
	return null.StringFrom((*s).String())
}

type BalanceMonitor struct {
	Enabled *bool
}

func (m *BalanceMonitor) setFrom(f *BalanceMonitor) {
	if v := f.Enabled; v != nil {
		m.Enabled = v
	}
}

type GasEstimator struct {
	Mode *string

	PriceDefault *utils.Wei
	PriceMax     *utils.Wei
	PriceMin     *utils.Wei

	LimitDefault    *uint32
	LimitMax        *uint32
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

func (e *GasEstimator) setFrom(f *GasEstimator) {
	if v := f.Mode; v != nil {
		e.Mode = v
	}
	if v := f.EIP1559DynamicFees; v != nil {
		e.EIP1559DynamicFees = v
	}
	if v := f.BumpPercent; v != nil {
		e.BumpPercent = v
	}
	if v := f.BumpThreshold; v != nil {
		e.BumpThreshold = v
	}
	if v := f.BumpTxDepth; v != nil {
		e.BumpTxDepth = v
	}
	if v := f.BumpMin; v != nil {
		e.BumpMin = v
	}
	if v := f.FeeCapDefault; v != nil {
		e.FeeCapDefault = v
	}
	if v := f.LimitDefault; v != nil {
		e.LimitDefault = v
	}
	if v := f.LimitMax; v != nil {
		e.LimitMax = v
	}
	if v := f.LimitMultiplier; v != nil {
		e.LimitMultiplier = v
	}
	if v := f.LimitTransfer; v != nil {
		e.LimitTransfer = v
	}
	if v := f.LimitOCRJobType; v != nil {
		e.LimitOCRJobType = v
	}
	if v := f.LimitDRJobType; v != nil {
		e.LimitDRJobType = v
	}
	if v := f.LimitVRFJobType; v != nil {
		e.LimitVRFJobType = v
	}
	if v := f.LimitFMJobType; v != nil {
		e.LimitFMJobType = v
	}
	if v := f.LimitKeeperJobType; v != nil {
		e.LimitKeeperJobType = v
	}
	if v := f.PriceDefault; v != nil {
		e.PriceDefault = v
	}
	if v := f.TipCapDefault; v != nil {
		e.TipCapDefault = v
	}
	if v := f.TipCapMinimum; v != nil {
		e.TipCapMinimum = v
	}
	if v := f.PriceMax; v != nil {
		e.PriceMax = v
	}
	if v := f.PriceMin; v != nil {
		e.PriceMin = v
	}
	if f.BlockHistory != nil {
		if e.BlockHistory == nil {
			e.BlockHistory = &BlockHistoryEstimator{}
		}
		e.BlockHistory.setFrom(f.BlockHistory)
	}
}

type BlockHistoryEstimator struct {
	BatchSize                 *uint32
	BlockHistorySize          *uint16
	EIP1559FeeCapBufferBlocks *uint16
	TransactionPercentile     *uint16
}

func (e *BlockHistoryEstimator) setFrom(f *BlockHistoryEstimator) {
	if v := f.BatchSize; v != nil {
		e.BatchSize = v
	}
	if v := f.BlockHistorySize; v != nil {
		e.BlockHistorySize = v
	}
	if v := f.EIP1559FeeCapBufferBlocks; v != nil {
		e.EIP1559FeeCapBufferBlocks = v
	}
	if v := f.TransactionPercentile; v != nil {
		e.TransactionPercentile = v
	}
}

type KeySpecificConfig []KeySpecific

func (ks KeySpecificConfig) ValidateConfig() (err error) {
	addrs := map[string]struct{}{}
	for _, k := range ks {
		addr := k.Key.String()
		if _, ok := addrs[addr]; ok {
			err = multierr.Append(err, v2.ErrInvalid{Name: "Key", Msg: "duplicate - must be unique", Value: addr})
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
	HistoryDepth     *uint32
	MaxBufferSize    *uint32
	SamplingInterval *models.Duration
}

func (t *HeadTracker) setFrom(f *HeadTracker) {
	if v := f.HistoryDepth; v != nil {
		t.HistoryDepth = v
	}
	if v := f.MaxBufferSize; v != nil {
		t.MaxBufferSize = v
	}
	if v := f.SamplingInterval; v != nil {
		t.SamplingInterval = v
	}
}

type NodePool struct {
	PollFailureThreshold *uint32
	PollInterval         *models.Duration
	SelectionMode        *string
}

func (p *NodePool) setFrom(f *NodePool) {
	if v := f.PollFailureThreshold; v != nil {
		p.PollFailureThreshold = v
	}
	if v := f.PollInterval; v != nil {
		p.PollInterval = v
	}
	if v := f.SelectionMode; v != nil {
		p.SelectionMode = v
	}
}

type OCR struct {
	ContractConfirmations              *uint16
	ContractTransmitterTransmitTimeout *models.Duration
	DatabaseTimeout                    *models.Duration
	ObservationGracePeriod             *models.Duration
}

func (o *OCR) setFrom(f *OCR) {
	if v := f.ContractConfirmations; v != nil {
		o.ContractConfirmations = v
	}
	if v := f.ContractTransmitterTransmitTimeout; v != nil {
		o.ContractTransmitterTransmitTimeout = v
	}
	if v := f.DatabaseTimeout; v != nil {
		o.DatabaseTimeout = v
	}
	if v := f.ObservationGracePeriod; v != nil {
		o.ObservationGracePeriod = v
	}
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
	if cfg.BlockHistoryEstimatorBlockDelay.Valid {
		v := uint16(cfg.BlockHistoryEstimatorBlockDelay.Int64)
		c.RPCBlockQueryDelay = &v
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
	if cfg.EvmGasLimitMax.Valid {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		v := uint32(cfg.EvmGasLimitMax.Int64)
		c.GasEstimator.LimitMax = &v
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

	if cfg.BlockHistoryEstimatorBlockHistorySize.Valid || cfg.BlockHistoryEstimatorEIP1559FeeCapBufferBlocks.Valid {
		if c.GasEstimator == nil {
			c.GasEstimator = &GasEstimator{}
		}
		c.GasEstimator.BlockHistory = &BlockHistoryEstimator{}
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
	if cfg.NodeNoNewHeadsThreshold != nil {
		c.NoNewHeadsThreshold = cfg.NodeNoNewHeadsThreshold
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
