package toml

import (
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strconv"

	"github.com/ethereum/go-ethereum/core/txpool/legacypool"
	"github.com/pelletier/go-toml/v2"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v4"

	commonassets "github.com/smartcontractkit/chainlink-common/pkg/assets"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

var ErrNotFound = errors.New("not found")

type HasEVMConfigs interface {
	EVMConfigs() EVMConfigs
}

type EVMConfigs []*EVMConfig

func (cs EVMConfigs) ValidateConfig() (err error) {
	return cs.validateKeys()
}

func (cs EVMConfigs) validateKeys() (err error) {
	// Unique chain IDs
	chainIDs := commonconfig.UniqueStrings{}
	for i, c := range cs {
		if chainIDs.IsDupeFmt(c.ChainID) {
			err = multierr.Append(err, commonconfig.NewErrDuplicate(fmt.Sprintf("%d.ChainID", i), c.ChainID.String()))
		}
	}

	// Unique node names
	names := commonconfig.UniqueStrings{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			if names.IsDupe(n.Name) {
				err = multierr.Append(err, commonconfig.NewErrDuplicate(fmt.Sprintf("%d.Nodes.%d.Name", i, j), *n.Name))
			}
		}
	}

	// Unique node WSURLs
	wsURLs := commonconfig.UniqueStrings{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			u := (*url.URL)(n.WSURL)
			if wsURLs.IsDupeFmt(u) {
				err = multierr.Append(err, commonconfig.NewErrDuplicate(fmt.Sprintf("%d.Nodes.%d.WSURL", i, j), u.String()))
			}
		}
	}

	// Unique node HTTPURLs
	httpURLs := commonconfig.UniqueStrings{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			u := (*url.URL)(n.HTTPURL)
			if httpURLs.IsDupeFmt(u) {
				err = multierr.Append(err, commonconfig.NewErrDuplicate(fmt.Sprintf("%d.Nodes.%d.HTTPURL", i, j), u.String()))
			}
		}
	}
	return
}

func (cs *EVMConfigs) SetFrom(fs *EVMConfigs) (err error) {
	if err1 := fs.validateKeys(); err1 != nil {
		return err1
	}
	for _, f := range *fs {
		if f.ChainID == nil {
			*cs = append(*cs, f)
		} else if i := slices.IndexFunc(*cs, func(c *EVMConfig) bool {
			return c.ChainID != nil && c.ChainID.Cmp(f.ChainID) == 0
		}); i == -1 {
			*cs = append(*cs, f)
		} else {
			(*cs)[i].SetFrom(f)
		}
	}
	return
}

func (cs EVMConfigs) totalChains() int {
	total := 0
	for _, ch := range cs {
		if ch == nil {
			continue
		}
		total++
	}
	return total
}
func (cs EVMConfigs) Chains(ids ...string) (r []commontypes.ChainStatus, total int, err error) {
	total = cs.totalChains()
	for _, ch := range cs {
		if ch == nil {
			continue
		}
		chainID := ch.ChainID.String()
		if len(ids) > 0 {
			var match bool
			for _, id := range ids {
				if id == chainID {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		ch2 := commontypes.ChainStatus{
			ID:      ch.ChainID.String(),
			Enabled: ch.IsEnabled(),
		}
		ch2.Config, err = ch.TOMLString()
		if err != nil {
			return
		}
		r = append(r, ch2)
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
	return types.Node{}, fmt.Errorf("node %s: %w", name, ErrNotFound)
}

func (cs EVMConfigs) NodeStatus(name string) (commontypes.NodeStatus, error) {
	for i := range cs {
		for _, n := range cs[i].Nodes {
			if n.Name != nil && *n.Name == name {
				return nodeStatus(n, cs[i].ChainID.String())
			}
		}
	}
	return commontypes.NodeStatus{}, fmt.Errorf("node %s: %w", name, ErrNotFound)
}

func legacyNode(n *Node, chainID *big.Big) (v2 types.Node) {
	v2.Name = *n.Name
	v2.EVMChainID = *chainID
	if n.HTTPURL != nil {
		v2.HTTPURL = null.StringFrom(n.HTTPURL.String())
	}
	if n.WSURL != nil {
		v2.WSURL = null.StringFrom(n.WSURL.String())
	}
	if n.SendOnly != nil {
		v2.SendOnly = *n.SendOnly
	}
	if n.Order != nil {
		v2.Order = *n.Order
	}
	return
}

func nodeStatus(n *Node, chainID string) (commontypes.NodeStatus, error) {
	var s commontypes.NodeStatus
	s.ChainID = chainID
	s.Name = *n.Name
	b, err := toml.Marshal(n)
	if err != nil {
		return commontypes.NodeStatus{}, err
	}
	s.Config = string(b)
	return s, nil
}

func (cs EVMConfigs) nodes(id string) (ns EVMNodes) {
	for _, c := range cs {
		if c.ChainID.String() == id {
			return c.Nodes
		}
	}
	return nil
}

func (cs EVMConfigs) Nodes(chainID string) (ns []types.Node, err error) {
	evmID, err := ChainIDInt64(chainID)
	if err != nil {
		return nil, fmt.Errorf("invalid evm chain id %q : %w", chainID, err)
	}
	nodes := cs.nodes(chainID)
	if nodes == nil {
		err = fmt.Errorf("no nodes: chain %q: %w", chainID, ErrNotFound)
		return
	}
	for _, n := range nodes {
		if n == nil {
			continue
		}

		ns = append(ns, legacyNode(n, big.NewI(evmID)))
	}
	return
}

func (cs EVMConfigs) NodeStatuses(chainIDs ...string) (ns []commontypes.NodeStatus, err error) {
	if len(chainIDs) == 0 {
		for i := range cs {
			for _, n := range cs[i].Nodes {
				if n == nil {
					continue
				}
				n2, err := nodeStatus(n, cs[i].ChainID.String())
				if err != nil {
					return nil, err
				}
				ns = append(ns, n2)
			}
		}
		return
	}
	for _, id := range chainIDs {
		for _, n := range cs.nodes(id) {
			if n == nil {
				continue
			}
			n2, err := nodeStatus(n, id)
			if err != nil {
				return nil, err
			}
			ns = append(ns, n2)
		}
	}
	return
}

type EVMNodes []*Node

func (ns *EVMNodes) SetFrom(fs *EVMNodes) {
	for _, f := range *fs {
		if f.Name == nil {
			*ns = append(*ns, f)
		} else if i := slices.IndexFunc(*ns, func(n *Node) bool {
			return n.Name != nil && *n.Name == *f.Name
		}); i == -1 {
			*ns = append(*ns, f)
		} else {
			(*ns)[i].SetFrom(f)
		}
	}
}

type EVMConfig struct {
	ChainID *big.Big
	Enabled *bool
	Chain
	Nodes EVMNodes
}

func (c *EVMConfig) IsEnabled() bool {
	return c.Enabled == nil || *c.Enabled
}

func (c *EVMConfig) SetFrom(f *EVMConfig) {
	if f.ChainID != nil {
		c.ChainID = f.ChainID
	}
	if f.Enabled != nil {
		c.Enabled = f.Enabled
	}
	c.Chain.SetFrom(&f.Chain)
	c.Nodes.SetFrom(&f.Nodes)
}

func (c *EVMConfig) ValidateConfig() (err error) {
	if c.ChainID == nil {
		err = multierr.Append(err, commonconfig.ErrMissing{Name: "ChainID", Msg: "required for all chains"})
	} else if c.ChainID.String() == "" {
		err = multierr.Append(err, commonconfig.ErrEmpty{Name: "ChainID", Msg: "required for all chains"})
	} else if must, ok := ChainTypeForID(c.ChainID); ok { // known chain id
		// Check if the parsed value matched the expected value
		is := c.ChainType.ChainType()
		if is != must {
			if must == "" {
				err = multierr.Append(err, commonconfig.ErrInvalid{Name: "ChainType", Value: c.ChainType.ChainType(),
					Msg: "must not be set with this chain id"})
			} else {
				err = multierr.Append(err, commonconfig.ErrInvalid{Name: "ChainType", Value: c.ChainType.ChainType(),
					Msg: fmt.Sprintf("only %q can be used with this chain id", must)})
			}
		}
	}

	if len(c.Nodes) == 0 {
		err = multierr.Append(err, commonconfig.ErrMissing{Name: "Nodes", Msg: "must have at least one node"})
	} else {
		var hasPrimary bool
		for _, n := range c.Nodes {
			if n.SendOnly != nil && *n.SendOnly {
				continue
			}
			hasPrimary = true
			break
		}
		if !hasPrimary {
			err = multierr.Append(err, commonconfig.ErrMissing{Name: "Nodes",
				Msg: "must have at least one primary node with WSURL"})
		}
	}

	err = multierr.Append(err, c.Chain.ValidateConfig())

	return
}

func (c *EVMConfig) TOMLString() (string, error) {
	b, err := toml.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

type Chain struct {
	AutoCreateKey             *bool
	BlockBackfillDepth        *uint32
	BlockBackfillSkip         *bool
	ChainType                 *chaintype.ChainTypeConfig
	FinalityDepth             *uint32
	FinalityTagEnabled        *bool
	FlagsContractAddress      *types.EIP55Address
	LinkContractAddress       *types.EIP55Address
	LogBackfillBatchSize      *uint32
	LogPollInterval           *commonconfig.Duration
	LogKeepBlocksDepth        *uint32
	LogPrunePageSize          *uint32
	BackupLogPollerBlockDelay *uint64
	MinIncomingConfirmations  *uint32
	MinContractPayment        *commonassets.Link
	NonceAutoSync             *bool
	NoNewHeadsThreshold       *commonconfig.Duration
	OperatorFactoryAddress    *types.EIP55Address
	RPCDefaultBatchSize       *uint32
	RPCBlockQueryDelay        *uint16

	Transactions   Transactions      `toml:",omitempty"`
	BalanceMonitor BalanceMonitor    `toml:",omitempty"`
	GasEstimator   GasEstimator      `toml:",omitempty"`
	HeadTracker    HeadTracker       `toml:",omitempty"`
	KeySpecific    KeySpecificConfig `toml:",omitempty"`
	NodePool       NodePool          `toml:",omitempty"`
	OCR            OCR               `toml:",omitempty"`
	OCR2           OCR2              `toml:",omitempty"`
	Workflow       Workflow          `toml:",omitempty"`
}

func (c *Chain) ValidateConfig() (err error) {
	if !c.ChainType.ChainType().IsValid() {
		err = multierr.Append(err, commonconfig.ErrInvalid{Name: "ChainType", Value: c.ChainType.ChainType(),
			Msg: chaintype.ErrInvalidChainType.Error()})
	}

	if c.GasEstimator.BumpTxDepth != nil && *c.GasEstimator.BumpTxDepth > *c.Transactions.MaxInFlight {
		err = multierr.Append(err, commonconfig.ErrInvalid{Name: "GasEstimator.BumpTxDepth", Value: *c.GasEstimator.BumpTxDepth,
			Msg: "must be less than or equal to Transactions.MaxInFlight"})
	}
	if *c.HeadTracker.HistoryDepth < *c.FinalityDepth {
		err = multierr.Append(err, commonconfig.ErrInvalid{Name: "HeadTracker.HistoryDepth", Value: *c.HeadTracker.HistoryDepth,
			Msg: "must be equal to or greater than FinalityDepth"})
	}
	if *c.FinalityDepth < 1 {
		err = multierr.Append(err, commonconfig.ErrInvalid{Name: "FinalityDepth", Value: *c.FinalityDepth,
			Msg: "must be greater than or equal to 1"})
	}
	if *c.MinIncomingConfirmations < 1 {
		err = multierr.Append(err, commonconfig.ErrInvalid{Name: "MinIncomingConfirmations", Value: *c.MinIncomingConfirmations,
			Msg: "must be greater than or equal to 1"})
	}

	// AutoPurge configs depend on ChainType so handling validation on per chain basis
	if c.Transactions.AutoPurge.Enabled != nil && *c.Transactions.AutoPurge.Enabled {
		chainType := c.ChainType.ChainType()
		switch chainType {
		case chaintype.ChainScroll:
			if c.Transactions.AutoPurge.DetectionApiUrl == nil {
				err = multierr.Append(err, commonconfig.ErrMissing{Name: "Transactions.AutoPurge.DetectionApiUrl", Msg: fmt.Sprintf("must be set for %s", chainType)})
			} else if c.Transactions.AutoPurge.DetectionApiUrl.IsZero() {
				err = multierr.Append(err, commonconfig.ErrInvalid{Name: "Transactions.AutoPurge.DetectionApiUrl", Value: c.Transactions.AutoPurge.DetectionApiUrl, Msg: fmt.Sprintf("must be set for %s", chainType)})
			} else {
				switch c.Transactions.AutoPurge.DetectionApiUrl.Scheme {
				case "http", "https":
				default:
					err = multierr.Append(err, commonconfig.ErrInvalid{Name: "Transactions.AutoPurge.DetectionApiUrl", Value: c.Transactions.AutoPurge.DetectionApiUrl.Scheme, Msg: "must be http or https"})
				}
			}
		case chaintype.ChainZkEvm:
			// No other configs are needed
		default:
			// Bump Threshold is required because the stuck tx heuristic relies on a minimum number of bump attempts to exist
			if c.GasEstimator.BumpThreshold == nil {
				err = multierr.Append(err, commonconfig.ErrMissing{Name: "GasEstimator.BumpThreshold", Msg: fmt.Sprintf("must be set if auto-purge feature is enabled for %s", chainType)})
			} else if *c.GasEstimator.BumpThreshold == 0 {
				err = multierr.Append(err, commonconfig.ErrInvalid{Name: "GasEstimator.BumpThreshold", Value: 0, Msg: fmt.Sprintf("cannot be 0 if auto-purge feature is enabled for %s", chainType)})
			}
			if c.Transactions.AutoPurge.Threshold == nil {
				err = multierr.Append(err, commonconfig.ErrMissing{Name: "Transactions.AutoPurge.Threshold", Msg: fmt.Sprintf("needs to be set if auto-purge feature is enabled for %s", chainType)})
			} else if *c.Transactions.AutoPurge.Threshold == 0 {
				err = multierr.Append(err, commonconfig.ErrInvalid{Name: "Transactions.AutoPurge.Threshold", Value: 0, Msg: fmt.Sprintf("cannot be 0 if auto-purge feature is enabled for %s", chainType)})
			}
			if c.Transactions.AutoPurge.MinAttempts == nil {
				err = multierr.Append(err, commonconfig.ErrMissing{Name: "Transactions.AutoPurge.MinAttempts", Msg: fmt.Sprintf("needs to be set if auto-purge feature is enabled for %s", chainType)})
			} else if *c.Transactions.AutoPurge.MinAttempts == 0 {
				err = multierr.Append(err, commonconfig.ErrInvalid{Name: "Transactions.AutoPurge.MinAttempts", Value: 0, Msg: fmt.Sprintf("cannot be 0 if auto-purge feature is enabled for %s", chainType)})
			}
		}
	}

	return
}

type Transactions struct {
	ForwardersEnabled    *bool
	MaxInFlight          *uint32
	MaxQueued            *uint32
	ReaperInterval       *commonconfig.Duration
	ReaperThreshold      *commonconfig.Duration
	ResendAfterThreshold *commonconfig.Duration

	AutoPurge AutoPurgeConfig `toml:",omitempty"`
}

func (t *Transactions) setFrom(f *Transactions) {
	if v := f.ForwardersEnabled; v != nil {
		t.ForwardersEnabled = v
	}
	if v := f.MaxInFlight; v != nil {
		t.MaxInFlight = v
	}
	if v := f.MaxQueued; v != nil {
		t.MaxQueued = v
	}
	if v := f.ReaperInterval; v != nil {
		t.ReaperInterval = v
	}
	if v := f.ReaperThreshold; v != nil {
		t.ReaperThreshold = v
	}
	if v := f.ResendAfterThreshold; v != nil {
		t.ResendAfterThreshold = v
	}
	t.AutoPurge.setFrom(&f.AutoPurge)
}

type AutoPurgeConfig struct {
	Enabled         *bool
	Threshold       *uint32
	MinAttempts     *uint32
	DetectionApiUrl *commonconfig.URL
}

func (a *AutoPurgeConfig) setFrom(f *AutoPurgeConfig) {
	if v := f.Enabled; v != nil {
		a.Enabled = v
	}
	if v := f.Threshold; v != nil {
		a.Threshold = v
	}
	if v := f.MinAttempts; v != nil {
		a.MinAttempts = v
	}
	if v := f.DetectionApiUrl; v != nil {
		a.DetectionApiUrl = v
	}
}

type OCR2 struct {
	Automation Automation `toml:",omitempty"`
}

func (o *OCR2) setFrom(f *OCR2) {
	o.Automation.setFrom(&f.Automation)
}

type Automation struct {
	GasLimit *uint32
}

func (a *Automation) setFrom(f *Automation) {
	if v := f.GasLimit; v != nil {
		a.GasLimit = v
	}
}

type Workflow struct {
	FromAddress      *types.EIP55Address `toml:",omitempty"`
	ForwarderAddress *types.EIP55Address `toml:",omitempty"`
}

func (m *Workflow) setFrom(f *Workflow) {
	if v := f.FromAddress; v != nil {
		m.FromAddress = v
	}
	if v := f.ForwarderAddress; v != nil {
		m.ForwarderAddress = v
	}
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

	PriceDefault *assets.Wei
	PriceMax     *assets.Wei
	PriceMin     *assets.Wei

	LimitDefault    *uint64
	LimitMax        *uint64
	LimitMultiplier *decimal.Decimal
	LimitTransfer   *uint64
	LimitJobType    GasLimitJobType `toml:",omitempty"`

	BumpMin       *assets.Wei
	BumpPercent   *uint16
	BumpThreshold *uint32
	BumpTxDepth   *uint32

	EIP1559DynamicFees *bool

	FeeCapDefault *assets.Wei
	TipCapDefault *assets.Wei
	TipCapMin     *assets.Wei

	BlockHistory BlockHistoryEstimator `toml:",omitempty"`
}

func (e *GasEstimator) ValidateConfig() (err error) {
	if uint64(*e.BumpPercent) < legacypool.DefaultConfig.PriceBump {
		err = multierr.Append(err, commonconfig.ErrInvalid{Name: "BumpPercent", Value: *e.BumpPercent,
			Msg: fmt.Sprintf("may not be less than Geth's default of %d", legacypool.DefaultConfig.PriceBump)})
	}
	if e.TipCapDefault.Cmp(e.TipCapMin) < 0 {
		err = multierr.Append(err, commonconfig.ErrInvalid{Name: "TipCapDefault", Value: e.TipCapDefault,
			Msg: "must be greater than or equal to TipCapMinimum"})
	}
	if e.FeeCapDefault.Cmp(e.TipCapDefault) < 0 {
		err = multierr.Append(err, commonconfig.ErrInvalid{Name: "FeeCapDefault", Value: e.TipCapDefault,
			Msg: "must be greater than or equal to TipCapDefault"})
	}
	if *e.Mode == "FixedPrice" && *e.BumpThreshold == 0 && *e.EIP1559DynamicFees && e.FeeCapDefault.Cmp(e.PriceMax) != 0 {
		err = multierr.Append(err, commonconfig.ErrInvalid{Name: "FeeCapDefault", Value: e.FeeCapDefault,
			Msg: fmt.Sprintf("must be equal to PriceMax (%s) since you are using FixedPrice estimation with gas bumping disabled in "+
				"EIP1559 mode - PriceMax will be used as the FeeCap for transactions instead of FeeCapDefault", e.PriceMax)})
	} else if e.FeeCapDefault.Cmp(e.PriceMax) > 0 {
		err = multierr.Append(err, commonconfig.ErrInvalid{Name: "FeeCapDefault", Value: e.FeeCapDefault,
			Msg: fmt.Sprintf("must be less than or equal to PriceMax (%s)", e.PriceMax)})
	}

	if e.PriceMin.Cmp(e.PriceDefault) > 0 {
		err = multierr.Append(err, commonconfig.ErrInvalid{Name: "PriceMin", Value: e.PriceMin,
			Msg: "must be less than or equal to PriceDefault"})
	}
	if e.PriceMax.Cmp(e.PriceDefault) < 0 {
		err = multierr.Append(err, commonconfig.ErrInvalid{Name: "PriceMax", Value: e.PriceMin,
			Msg: "must be greater than or equal to PriceDefault"})
	}
	if *e.Mode == "BlockHistory" && *e.BlockHistory.BlockHistorySize <= 0 {
		err = multierr.Append(err, commonconfig.ErrInvalid{Name: "BlockHistory.BlockHistorySize", Value: *e.BlockHistory.BlockHistorySize,
			Msg: "must be greater than or equal to 1 with BlockHistory Mode"})
	}

	return
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
	if v := f.PriceDefault; v != nil {
		e.PriceDefault = v
	}
	if v := f.TipCapDefault; v != nil {
		e.TipCapDefault = v
	}
	if v := f.TipCapMin; v != nil {
		e.TipCapMin = v
	}
	if v := f.PriceMax; v != nil {
		e.PriceMax = v
	}
	if v := f.PriceMin; v != nil {
		e.PriceMin = v
	}
	e.LimitJobType.setFrom(&f.LimitJobType)
	e.BlockHistory.setFrom(&f.BlockHistory)
}

type GasLimitJobType struct {
	OCR    *uint32 `toml:",inline"`
	OCR2   *uint32 `toml:",inline"`
	DR     *uint32 `toml:",inline"`
	VRF    *uint32 `toml:",inline"`
	FM     *uint32 `toml:",inline"`
	Keeper *uint32 `toml:",inline"`
}

func (t *GasLimitJobType) setFrom(f *GasLimitJobType) {
	if f.OCR != nil {
		t.OCR = f.OCR
	}
	if f.OCR2 != nil {
		t.OCR2 = f.OCR2
	}
	if f.DR != nil {
		t.DR = f.DR
	}
	if f.VRF != nil {
		t.VRF = f.VRF
	}
	if f.FM != nil {
		t.FM = f.FM
	}
	if f.Keeper != nil {
		t.Keeper = f.Keeper
	}
}

type BlockHistoryEstimator struct {
	BatchSize                 *uint32
	BlockHistorySize          *uint16
	CheckInclusionBlocks      *uint16
	CheckInclusionPercentile  *uint16
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
	if v := f.CheckInclusionBlocks; v != nil {
		e.CheckInclusionBlocks = v
	}
	if v := f.CheckInclusionPercentile; v != nil {
		e.CheckInclusionPercentile = v
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
			err = multierr.Append(err, commonconfig.NewErrDuplicate("Key", addr))
		} else {
			addrs[addr] = struct{}{}
		}
	}
	return
}

type KeySpecific struct {
	Key          *types.EIP55Address
	GasEstimator KeySpecificGasEstimator `toml:",omitempty"`
}

type KeySpecificGasEstimator struct {
	PriceMax *assets.Wei
}

func (e *KeySpecificGasEstimator) setFrom(f *KeySpecificGasEstimator) {
	if v := f.PriceMax; v != nil {
		e.PriceMax = v
	}
}

type HeadTracker struct {
	HistoryDepth            *uint32
	MaxBufferSize           *uint32
	SamplingInterval        *commonconfig.Duration
	MaxAllowedFinalityDepth *uint32
	FinalityTagBypass       *bool
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
	if v := f.MaxAllowedFinalityDepth; v != nil {
		t.MaxAllowedFinalityDepth = v
	}
	if v := f.FinalityTagBypass; v != nil {
		t.FinalityTagBypass = v
	}
}

func (t *HeadTracker) ValidateConfig() (err error) {
	if *t.MaxAllowedFinalityDepth < 1 {
		err = multierr.Append(err, commonconfig.ErrInvalid{Name: "MaxAllowedFinalityDepth", Value: *t.MaxAllowedFinalityDepth,
			Msg: "must be greater than or equal to 1"})
	}

	return
}

type ClientErrors struct {
	NonceTooLow                       *string `toml:",omitempty"`
	NonceTooHigh                      *string `toml:",omitempty"`
	ReplacementTransactionUnderpriced *string `toml:",omitempty"`
	LimitReached                      *string `toml:",omitempty"`
	TransactionAlreadyInMempool       *string `toml:",omitempty"`
	TerminallyUnderpriced             *string `toml:",omitempty"`
	InsufficientEth                   *string `toml:",omitempty"`
	TxFeeExceedsCap                   *string `toml:",omitempty"`
	L2FeeTooLow                       *string `toml:",omitempty"`
	L2FeeTooHigh                      *string `toml:",omitempty"`
	L2Full                            *string `toml:",omitempty"`
	TransactionAlreadyMined           *string `toml:",omitempty"`
	Fatal                             *string `toml:",omitempty"`
	ServiceUnavailable                *string `toml:",omitempty"`
}

func (r *ClientErrors) setFrom(f *ClientErrors) bool {
	if v := f.NonceTooLow; v != nil {
		r.NonceTooLow = v
	}
	if v := f.NonceTooHigh; v != nil {
		r.NonceTooHigh = v
	}
	if v := f.ReplacementTransactionUnderpriced; v != nil {
		r.ReplacementTransactionUnderpriced = v
	}
	if v := f.LimitReached; v != nil {
		r.LimitReached = v
	}
	if v := f.TransactionAlreadyInMempool; v != nil {
		r.TransactionAlreadyInMempool = v
	}
	if v := f.TerminallyUnderpriced; v != nil {
		r.TerminallyUnderpriced = v
	}
	if v := f.InsufficientEth; v != nil {
		r.InsufficientEth = v
	}
	if v := f.TxFeeExceedsCap; v != nil {
		r.TxFeeExceedsCap = v
	}
	if v := f.L2FeeTooLow; v != nil {
		r.L2FeeTooLow = v
	}
	if v := f.L2FeeTooHigh; v != nil {
		r.L2FeeTooHigh = v
	}
	if v := f.L2Full; v != nil {
		r.L2Full = v
	}
	if v := f.TransactionAlreadyMined; v != nil {
		r.TransactionAlreadyMined = v
	}
	if v := f.Fatal; v != nil {
		r.Fatal = v
	}
	if v := f.ServiceUnavailable; v != nil {
		r.ServiceUnavailable = v
	}
	return true
}

type NodePool struct {
	PollFailureThreshold       *uint32
	PollInterval               *commonconfig.Duration
	SelectionMode              *string
	SyncThreshold              *uint32
	LeaseDuration              *commonconfig.Duration
	NodeIsSyncingEnabled       *bool
	FinalizedBlockPollInterval *commonconfig.Duration
	Errors                     ClientErrors `toml:",omitempty"`
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
	if v := f.SyncThreshold; v != nil {
		p.SyncThreshold = v
	}
	if v := f.LeaseDuration; v != nil {
		p.LeaseDuration = v
	}
	if v := f.NodeIsSyncingEnabled; v != nil {
		p.NodeIsSyncingEnabled = v
	}
	if v := f.FinalizedBlockPollInterval; v != nil {
		p.FinalizedBlockPollInterval = v
	}
	p.Errors.setFrom(&f.Errors)
}

type OCR struct {
	ContractConfirmations              *uint16
	ContractTransmitterTransmitTimeout *commonconfig.Duration
	DatabaseTimeout                    *commonconfig.Duration
	DeltaCOverride                     *commonconfig.Duration
	DeltaCJitterOverride               *commonconfig.Duration
	ObservationGracePeriod             *commonconfig.Duration
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
	if v := f.DeltaCOverride; v != nil {
		o.DeltaCOverride = v
	}
	if v := f.DeltaCJitterOverride; v != nil {
		o.DeltaCJitterOverride = v
	}
	if v := f.ObservationGracePeriod; v != nil {
		o.ObservationGracePeriod = v
	}
}

type Node struct {
	Name     *string
	WSURL    *commonconfig.URL
	HTTPURL  *commonconfig.URL
	SendOnly *bool
	Order    *int32
}

func (n *Node) ValidateConfig() (err error) {
	if n.Name == nil {
		err = multierr.Append(err, commonconfig.ErrMissing{Name: "Name", Msg: "required for all nodes"})
	} else if *n.Name == "" {
		err = multierr.Append(err, commonconfig.ErrEmpty{Name: "Name", Msg: "required for all nodes"})
	}

	var sendOnly bool
	if n.SendOnly != nil {
		sendOnly = *n.SendOnly
	}
	if n.WSURL == nil {
		if !sendOnly {
			err = multierr.Append(err, commonconfig.ErrMissing{Name: "WSURL", Msg: "required for primary nodes"})
		}
	} else if n.WSURL.IsZero() {
		if !sendOnly {
			err = multierr.Append(err, commonconfig.ErrEmpty{Name: "WSURL", Msg: "required for primary nodes"})
		}
	} else {
		switch n.WSURL.Scheme {
		case "ws", "wss":
		default:
			err = multierr.Append(err, commonconfig.ErrInvalid{Name: "WSURL", Value: n.WSURL.Scheme, Msg: "must be ws or wss"})
		}
	}

	if n.HTTPURL == nil {
		err = multierr.Append(err, commonconfig.ErrMissing{Name: "HTTPURL", Msg: "required for all nodes"})
	} else if n.HTTPURL.IsZero() {
		err = multierr.Append(err, commonconfig.ErrEmpty{Name: "HTTPURL", Msg: "required for all nodes"})
	} else {
		switch n.HTTPURL.Scheme {
		case "http", "https":
		default:
			err = multierr.Append(err, commonconfig.ErrInvalid{Name: "HTTPURL", Value: n.HTTPURL.Scheme, Msg: "must be http or https"})
		}
	}

	if n.Order != nil && (*n.Order < 1 || *n.Order > 100) {
		err = multierr.Append(err, commonconfig.ErrInvalid{Name: "Order", Value: *n.Order, Msg: "must be between 1 and 100"})
	} else if n.Order == nil {
		z := int32(100)
		n.Order = &z
	}

	return
}

func (n *Node) SetFrom(f *Node) {
	if f.Name != nil {
		n.Name = f.Name
	}
	if f.WSURL != nil {
		n.WSURL = f.WSURL
	}
	if f.HTTPURL != nil {
		n.HTTPURL = f.HTTPURL
	}
	if f.SendOnly != nil {
		n.SendOnly = f.SendOnly
	}
	if f.Order != nil {
		n.Order = f.Order
	}
}

func ChainIDInt64(cid string) (int64, error) {
	return strconv.ParseInt(cid, 10, 64)
}
