package solana

import (
	"fmt"
	"net/url"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
	"go.uber.org/multierr"
	"golang.org/x/exp/slices"
	"gopkg.in/guregu/null.v4"

	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	soldb "github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	v2 "github.com/smartcontractkit/chainlink/core/config/v2"
)

type SolanaConfigs []*SolanaConfig

func (cs SolanaConfigs) ValidateConfig() (err error) {
	// Unique chain IDs
	chainIDs := v2.UniqueStrings{}
	for i, c := range cs {
		if chainIDs.IsDupe(c.ChainID) {
			err = multierr.Append(err, v2.NewErrDuplicate(fmt.Sprintf("%d.ChainID", i), *c.ChainID))
		}
	}

	// Unique node names
	names := v2.UniqueStrings{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			if names.IsDupe(n.Name) {
				err = multierr.Append(err, v2.NewErrDuplicate(fmt.Sprintf("%d.Nodes.%d.Name", i, j), *n.Name))
			}
		}
	}

	// Unique URLs
	urls := v2.UniqueStrings{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			u := (*url.URL)(n.URL)
			if urls.IsDupeFmt(u) {
				err = multierr.Append(err, v2.NewErrDuplicate(fmt.Sprintf("%d.Nodes.%d.URL", i, j), u.String()))
			}
		}
	}
	return
}

func (cs *SolanaConfigs) SetFrom(fs *SolanaConfigs) {
	for _, f := range *fs {
		if f.ChainID == nil {
			*cs = append(*cs, f)
		} else if i := slices.IndexFunc(*cs, func(c *SolanaConfig) bool {
			return c.ChainID != nil && *c.ChainID == *f.ChainID
		}); i == -1 {
			*cs = append(*cs, f)
		} else {
			(*cs)[i].SetFrom(f)
		}
	}
}

func (cs SolanaConfigs) Chains(ids ...string) (chains []DBChain) {
	for _, ch := range cs {
		if ch == nil {
			continue
		}
		if len(ids) > 0 {
			var match bool
			for _, id := range ids {
				if id == *ch.ChainID {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		chains = append(chains, ch.AsV1())
	}
	return
}

func (cs SolanaConfigs) Node(name string) (soldb.Node, error) {
	for i := range cs {
		for _, n := range cs[i].Nodes {
			if n.Name != nil && *n.Name == name {
				return legacySolNode(n, *cs[i].ChainID), nil
			}
		}
	}
	return soldb.Node{}, nil
}

func (cs SolanaConfigs) Nodes() (ns []soldb.Node) {
	for i := range cs {
		for _, n := range cs[i].Nodes {
			if n == nil {
				continue
			}
			ns = append(ns, legacySolNode(n, *cs[i].ChainID))
		}
	}
	return
}

func (cs SolanaConfigs) NodesByID(chainIDs ...string) (ns []soldb.Node) {
	for i := range cs {
		var match bool
		for _, id := range chainIDs {
			if id == *cs[i].ChainID {
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
			ns = append(ns, legacySolNode(n, *cs[i].ChainID))
		}
	}
	return
}

type SolanaNodes []*solcfg.Node

func (ns *SolanaNodes) SetFrom(fs *SolanaNodes) {
	for _, f := range *fs {
		if f.Name == nil {
			*ns = append(*ns, f)
		} else if i := slices.IndexFunc(*ns, func(n *solcfg.Node) bool {
			return n.Name != nil && *n.Name == *f.Name
		}); i == -1 {
			*ns = append(*ns, f)
		} else {
			setFromNode((*ns)[i], f)
		}
	}
}

func setFromNode(n, f *solcfg.Node) {
	if f.Name != nil {
		n.Name = f.Name
	}
	if f.URL != nil {
		n.URL = f.URL
	}
}

func legacySolNode(n *solcfg.Node, chainID string) soldb.Node {
	return soldb.Node{
		Name:          *n.Name,
		SolanaChainID: chainID,
		SolanaURL:     (*url.URL)(n.URL).String(),
	}
}

type SolanaConfig struct {
	ChainID *string
	Enabled *bool
	solcfg.Chain
	Nodes SolanaNodes
}

func (c *SolanaConfig) IsEnabled() bool {
	return c.Enabled == nil || *c.Enabled
}

func (c *SolanaConfig) SetFrom(f *SolanaConfig) {
	if f.ChainID != nil {
		c.ChainID = f.ChainID
	}
	if f.Enabled != nil {
		c.Enabled = f.Enabled
	}
	setFromChain(&c.Chain, &f.Chain)
	c.Nodes.SetFrom(&f.Nodes)
}

func setFromChain(c, f *solcfg.Chain) {
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
}

func (c *SolanaConfig) SetFromDB(ch DBChain, nodes []soldb.Node) error {
	c.ChainID = &ch.ID
	c.Enabled = &ch.Enabled

	if err := c.Chain.SetFromDB(ch.Cfg); err != nil {
		return err
	}
	for _, db := range nodes {
		var n solcfg.Node
		if err := n.SetFromDB(db); err != nil {
			return err
		}
		c.Nodes = append(c.Nodes, &n)
	}
	return nil
}

func (c *SolanaConfig) ValidateConfig() (err error) {
	if c.ChainID == nil {
		err = multierr.Append(err, v2.ErrMissing{Name: "ChainID", Msg: "required for all chains"})
	} else if *c.ChainID == "" {
		err = multierr.Append(err, v2.ErrEmpty{Name: "ChainID", Msg: "required for all chains"})
	}

	if len(c.Nodes) == 0 {
		err = multierr.Append(err, v2.ErrMissing{Name: "Nodes", Msg: "must have at least one node"})
	}
	return
}

func (c *SolanaConfig) AsV1() DBChain {
	return DBChain{
		ID:      *c.ChainID,
		Enabled: c.IsEnabled(),
		Cfg: &soldb.ChainCfg{
			BalancePollPeriod:       c.Chain.BalancePollPeriod,
			ConfirmPollPeriod:       c.Chain.ConfirmPollPeriod,
			OCR2CachePollPeriod:     c.Chain.OCR2CachePollPeriod,
			OCR2CacheTTL:            c.Chain.OCR2CacheTTL,
			TxTimeout:               c.Chain.TxTimeout,
			TxRetryTimeout:          c.Chain.TxRetryTimeout,
			TxConfirmTimeout:        c.Chain.TxConfirmTimeout,
			SkipPreflight:           null.BoolFromPtr(c.Chain.SkipPreflight),
			Commitment:              null.StringFromPtr(c.Chain.Commitment),
			MaxRetries:              null.IntFromPtr(c.Chain.MaxRetries),
			FeeEstimatorMode:        null.StringFromPtr(c.Chain.FeeEstimatorMode),
			ComputeUnitPriceMax:     null.IntFrom(int64(*c.Chain.ComputeUnitPriceMax)),
			ComputeUnitPriceMin:     null.IntFrom(int64(*c.Chain.ComputeUnitPriceMin)),
			ComputeUnitPriceDefault: null.IntFrom(int64(*c.Chain.ComputeUnitPriceDefault)),
			FeeBumpPeriod:           c.Chain.FeeBumpPeriod,
		},
	}
}

var _ solcfg.Config = &SolanaConfig{}

func (c *SolanaConfig) BalancePollPeriod() time.Duration {
	return c.Chain.BalancePollPeriod.Duration()
}

func (c *SolanaConfig) ConfirmPollPeriod() time.Duration {
	return c.Chain.ConfirmPollPeriod.Duration()
}

func (c *SolanaConfig) OCR2CachePollPeriod() time.Duration {
	return c.Chain.OCR2CachePollPeriod.Duration()
}

func (c *SolanaConfig) OCR2CacheTTL() time.Duration {
	return c.Chain.OCR2CacheTTL.Duration()
}

func (c *SolanaConfig) TxTimeout() time.Duration {
	return c.Chain.TxTimeout.Duration()
}

func (c *SolanaConfig) TxRetryTimeout() time.Duration {
	return c.Chain.TxRetryTimeout.Duration()
}

func (c *SolanaConfig) TxConfirmTimeout() time.Duration {
	return c.Chain.TxConfirmTimeout.Duration()
}

func (c *SolanaConfig) SkipPreflight() bool {
	return *c.Chain.SkipPreflight
}

func (c *SolanaConfig) Commitment() rpc.CommitmentType {
	return rpc.CommitmentType(*c.Chain.Commitment)
}

func (c *SolanaConfig) MaxRetries() *uint {
	if c.Chain.MaxRetries == nil {
		return nil
	}
	mr := uint(*c.Chain.MaxRetries)
	return &mr
}

func (c *SolanaConfig) FeeEstimatorMode() string {
	return *c.Chain.FeeEstimatorMode
}

func (c *SolanaConfig) ComputeUnitPriceMax() uint64 {
	return *c.Chain.ComputeUnitPriceMax
}

func (c *SolanaConfig) ComputeUnitPriceMin() uint64 {
	return *c.Chain.ComputeUnitPriceMin
}

func (c *SolanaConfig) ComputeUnitPriceDefault() uint64 {
	return *c.Chain.ComputeUnitPriceDefault
}

func (c *SolanaConfig) FeeBumpPeriod() time.Duration {
	return c.Chain.FeeBumpPeriod.Duration()
}

func (c *SolanaConfig) Update(cfg soldb.ChainCfg) {
	panic(fmt.Errorf("cannot update: %v", v2.ErrUnsupported))
}
