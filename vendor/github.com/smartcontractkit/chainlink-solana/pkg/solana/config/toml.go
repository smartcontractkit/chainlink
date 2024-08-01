package config

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	relaytypes "github.com/smartcontractkit/chainlink-common/pkg/types"
)

type TOMLConfigs []*TOMLConfig

func (cs TOMLConfigs) ValidateConfig() (err error) {
	return cs.validateKeys()
}

func (cs TOMLConfigs) validateKeys() (err error) {
	// Unique chain IDs
	chainIDs := config.UniqueStrings{}
	for i, c := range cs {
		if chainIDs.IsDupe(c.ChainID) {
			err = errors.Join(err, config.NewErrDuplicate(fmt.Sprintf("%d.ChainID", i), *c.ChainID))
		}
	}

	// Unique node names
	names := config.UniqueStrings{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			if names.IsDupe(n.Name) {
				err = errors.Join(err, config.NewErrDuplicate(fmt.Sprintf("%d.Nodes.%d.Name", i, j), *n.Name))
			}
		}
	}

	// Unique URLs
	urls := config.UniqueStrings{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			u := (*url.URL)(n.URL)
			if urls.IsDupeFmt(u) {
				err = errors.Join(err, config.NewErrDuplicate(fmt.Sprintf("%d.Nodes.%d.URL", i, j), u.String()))
			}
		}
	}
	return
}

func (cs *TOMLConfigs) SetFrom(fs *TOMLConfigs) (err error) {
	if err1 := fs.validateKeys(); err1 != nil {
		return err1
	}
	for _, f := range *fs {
		if f.ChainID == nil {
			*cs = append(*cs, f)
		} else if i := slices.IndexFunc(*cs, func(c *TOMLConfig) bool {
			return c.ChainID != nil && *c.ChainID == *f.ChainID
		}); i == -1 {
			*cs = append(*cs, f)
		} else {
			(*cs)[i].SetFrom(f)
		}
	}
	return
}

func NodeStatus(n *Node, id string) (relaytypes.NodeStatus, error) {
	var s relaytypes.NodeStatus
	s.ChainID = id
	s.Name = *n.Name
	b, err := toml.Marshal(n)
	if err != nil {
		return relaytypes.NodeStatus{}, err
	}
	s.Config = string(b)
	return s, nil
}

type Nodes []*Node

func (ns *Nodes) SetFrom(fs *Nodes) {
	for _, f := range *fs {
		if f.Name == nil {
			*ns = append(*ns, f)
		} else if i := slices.IndexFunc(*ns, func(n *Node) bool {
			return n.Name != nil && *n.Name == *f.Name
		}); i == -1 {
			*ns = append(*ns, f)
		} else {
			setFromNode((*ns)[i], f)
		}
	}
}

func setFromNode(n, f *Node) {
	if f.Name != nil {
		n.Name = f.Name
	}
	if f.URL != nil {
		n.URL = f.URL
	}
}

type TOMLConfig struct {
	ChainID *string
	// Do not access directly, use [IsEnabled]
	Enabled *bool
	Chain
	Nodes Nodes
}

func (c *TOMLConfig) IsEnabled() bool {
	return c.Enabled == nil || *c.Enabled
}

func (c *TOMLConfig) SetFrom(f *TOMLConfig) {
	if f.ChainID != nil {
		c.ChainID = f.ChainID
	}
	if f.Enabled != nil {
		c.Enabled = f.Enabled
	}
	setFromChain(&c.Chain, &f.Chain)
	c.Nodes.SetFrom(&f.Nodes)
}

func setFromChain(c, f *Chain) {
	if f.BalancePollPeriod != nil {
		c.BalancePollPeriod = f.BalancePollPeriod
	}
	if f.ConfirmPollPeriod != nil {
		c.ConfirmPollPeriod = f.ConfirmPollPeriod
	}
	if f.OCR2CachePollPeriod != nil {
		c.OCR2CachePollPeriod = f.OCR2CachePollPeriod
	}
	if f.OCR2CacheTTL != nil {
		c.OCR2CacheTTL = f.OCR2CacheTTL
	}
	if f.TxTimeout != nil {
		c.TxTimeout = f.TxTimeout
	}
	if f.TxRetryTimeout != nil {
		c.TxRetryTimeout = f.TxRetryTimeout
	}
	if f.TxConfirmTimeout != nil {
		c.TxConfirmTimeout = f.TxConfirmTimeout
	}
	if f.SkipPreflight != nil {
		c.SkipPreflight = f.SkipPreflight
	}
	if f.Commitment != nil {
		c.Commitment = f.Commitment
	}
	if f.MaxRetries != nil {
		c.MaxRetries = f.MaxRetries
	}
	if f.FeeEstimatorMode != nil {
		c.FeeEstimatorMode = f.FeeEstimatorMode
	}
	if f.ComputeUnitPriceMax != nil {
		c.ComputeUnitPriceMax = f.ComputeUnitPriceMax
	}
	if f.ComputeUnitPriceMin != nil {
		c.ComputeUnitPriceMin = f.ComputeUnitPriceMin
	}
	if f.ComputeUnitPriceDefault != nil {
		c.ComputeUnitPriceDefault = f.ComputeUnitPriceDefault
	}
	if f.FeeBumpPeriod != nil {
		c.FeeBumpPeriod = f.FeeBumpPeriod
	}
	if f.BlockHistoryPollPeriod != nil {
		c.BlockHistoryPollPeriod = f.BlockHistoryPollPeriod
	}
}

func (c *TOMLConfig) ValidateConfig() (err error) {
	if c.ChainID == nil {
		err = errors.Join(err, config.ErrMissing{Name: "ChainID", Msg: "required for all chains"})
	} else if *c.ChainID == "" {
		err = errors.Join(err, config.ErrEmpty{Name: "ChainID", Msg: "required for all chains"})
	}

	if len(c.Nodes) == 0 {
		err = errors.Join(err, config.ErrMissing{Name: "Nodes", Msg: "must have at least one node"})
	}
	return
}

func (c *TOMLConfig) TOMLString() (string, error) {
	b, err := toml.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

var _ Config = &TOMLConfig{}

func (c *TOMLConfig) BalancePollPeriod() time.Duration {
	return c.Chain.BalancePollPeriod.Duration()
}

func (c *TOMLConfig) ConfirmPollPeriod() time.Duration {
	return c.Chain.ConfirmPollPeriod.Duration()
}

func (c *TOMLConfig) OCR2CachePollPeriod() time.Duration {
	return c.Chain.OCR2CachePollPeriod.Duration()
}

func (c *TOMLConfig) OCR2CacheTTL() time.Duration {
	return c.Chain.OCR2CacheTTL.Duration()
}

func (c *TOMLConfig) TxTimeout() time.Duration {
	return c.Chain.TxTimeout.Duration()
}

func (c *TOMLConfig) TxRetryTimeout() time.Duration {
	return c.Chain.TxRetryTimeout.Duration()
}

func (c *TOMLConfig) TxConfirmTimeout() time.Duration {
	return c.Chain.TxConfirmTimeout.Duration()
}

func (c *TOMLConfig) SkipPreflight() bool {
	return *c.Chain.SkipPreflight
}

func (c *TOMLConfig) Commitment() rpc.CommitmentType {
	return rpc.CommitmentType(*c.Chain.Commitment)
}

func (c *TOMLConfig) MaxRetries() *uint {
	if c.Chain.MaxRetries == nil {
		return nil
	}
	if *c.Chain.MaxRetries < 0 {
		return nil // interpret negative numbers as nil (prevents unlikely case of overflow)
	}
	mr := uint(*c.Chain.MaxRetries)
	return &mr
}

func (c *TOMLConfig) FeeEstimatorMode() string {
	return *c.Chain.FeeEstimatorMode
}

func (c *TOMLConfig) ComputeUnitPriceMax() uint64 {
	return *c.Chain.ComputeUnitPriceMax
}

func (c *TOMLConfig) ComputeUnitPriceMin() uint64 {
	return *c.Chain.ComputeUnitPriceMin
}

func (c *TOMLConfig) ComputeUnitPriceDefault() uint64 {
	return *c.Chain.ComputeUnitPriceDefault
}

func (c *TOMLConfig) FeeBumpPeriod() time.Duration {
	return c.Chain.FeeBumpPeriod.Duration()
}

func (c *TOMLConfig) BlockHistoryPollPeriod() time.Duration {
	return c.Chain.BlockHistoryPollPeriod.Duration()
}

func (c *TOMLConfig) ListNodes() Nodes {
	return c.Nodes
}

func NewDefault() *TOMLConfig {
	cfg := &TOMLConfig{}
	cfg.SetDefaults()
	return cfg
}
