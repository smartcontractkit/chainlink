package starknet

import (
	"fmt"
	"net/url"
	"time"

	"go.uber.org/multierr"
	"golang.org/x/exp/slices"
	"gopkg.in/guregu/null.v4"

	stkcfg "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"
	starknetdb "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/core/chains/starknet/types"
	v2 "github.com/smartcontractkit/chainlink/core/config/v2"
)

type StarknetConfigs []*StarknetConfig

func (cs StarknetConfigs) ValidateConfig() (err error) {
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

func (cs *StarknetConfigs) SetFrom(fs *StarknetConfigs) {
	for _, f := range *fs {
		if f.ChainID == nil {
			*cs = append(*cs, f)
		} else if i := slices.IndexFunc(*cs, func(c *StarknetConfig) bool {
			return c.ChainID != nil && *c.ChainID == *f.ChainID
		}); i == -1 {
			*cs = append(*cs, f)
		} else {
			(*cs)[i].SetFrom(f)
		}
	}
}

func (cs StarknetConfigs) Chains(ids ...string) (chains []types.DBChain) {
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

func (cs StarknetConfigs) Node(name string) (n db.Node, err error) {
	for i := range cs {
		for _, n := range cs[i].Nodes {
			if n.Name != nil && *n.Name == name {
				return legacyNode(n, *cs[i].ChainID), nil
			}
		}
	}
	return
}

func (cs StarknetConfigs) Nodes() (ns []db.Node) {
	for i := range cs {
		for _, n := range cs[i].Nodes {
			if n == nil {
				continue
			}
			ns = append(ns, legacyNode(n, *cs[i].ChainID))
		}
	}
	return
}

func (cs StarknetConfigs) NodesByID(chainIDs ...string) (ns []db.Node) {
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
			ns = append(ns, legacyNode(n, *cs[i].ChainID))
		}
	}
	return
}

type StarknetConfig struct {
	ChainID *string
	Enabled *bool
	stkcfg.Chain
	Nodes StarknetNodes
}

func (c *StarknetConfig) IsEnabled() bool {
	return c.Enabled == nil || *c.Enabled
}

func (c *StarknetConfig) SetFrom(f *StarknetConfig) {
	if f.ChainID != nil {
		c.ChainID = f.ChainID
	}
	if f.Enabled != nil {
		c.Enabled = f.Enabled
	}
	setFromChain(&c.Chain, &f.Chain)
	c.Nodes.SetFrom(&f.Nodes)
}

func setFromChain(c, f *stkcfg.Chain) {
	if f.OCR2CachePollPeriod != nil {
		c.OCR2CachePollPeriod = f.OCR2CachePollPeriod
	}
	if f.OCR2CacheTTL != nil {
		c.OCR2CacheTTL = f.OCR2CacheTTL
	}
	if f.RequestTimeout != nil {
		c.RequestTimeout = f.RequestTimeout
	}
	if f.TxTimeout != nil {
		c.TxTimeout = f.TxTimeout
	}
	if f.TxSendFrequency != nil {
		c.TxSendFrequency = f.TxSendFrequency
	}
	if f.TxMaxBatchSize != nil {
		c.TxMaxBatchSize = f.TxMaxBatchSize
	}
}

func (c *StarknetConfig) SetFromDB(ch types.DBChain, nodes []db.Node) error {
	c.ChainID = &ch.ID
	c.Enabled = &ch.Enabled

	if err := c.Chain.SetFromDB(ch.Cfg); err != nil {
		return err
	}
	for _, db := range nodes {
		var n stkcfg.Node
		if err := n.SetFromDB(db); err != nil {
			return err
		}
		c.Nodes = append(c.Nodes, &n)
	}

	return nil
}

func (c *StarknetConfig) ValidateConfig() (err error) {
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

func (c *StarknetConfig) AsV1() types.DBChain {
	return types.DBChain{
		ID:      *c.ChainID,
		Enabled: c.IsEnabled(),
		Cfg: &starknetdb.ChainCfg{
			OCR2CachePollPeriod: c.Chain.OCR2CachePollPeriod,
			OCR2CacheTTL:        c.Chain.OCR2CacheTTL,
			RequestTimeout:      c.Chain.RequestTimeout,
			TxTimeout:           c.Chain.TxTimeout,
			TxSendFrequency:     c.Chain.TxSendFrequency,
			TxMaxBatchSize:      null.IntFromPtr(c.Chain.TxMaxBatchSize),
		},
	}
}

type StarknetNodes []*stkcfg.Node

func (ns *StarknetNodes) SetFrom(fs *StarknetNodes) {
	for _, f := range *fs {
		if f.Name == nil {
			*ns = append(*ns, f)
		} else if i := slices.IndexFunc(*ns, func(n *stkcfg.Node) bool {
			return n.Name != nil && *n.Name == *f.Name
		}); i == -1 {
			*ns = append(*ns, f)
		} else {
			setFromNode((*ns)[i], f)
		}
	}
}

func setFromNode(n, f *stkcfg.Node) {
	if f.Name != nil {
		n.Name = f.Name
	}
	if f.URL != nil {
		n.URL = f.URL
	}
}

func legacyNode(n *stkcfg.Node, id string) db.Node {
	return db.Node{
		Name:    *n.Name,
		ChainID: id,
		URL:     (*url.URL)(n.URL).String(),
	}
}

var _ stkcfg.Config = &StarknetConfig{}

func (c *StarknetConfig) TxTimeout() time.Duration {
	return c.Chain.TxTimeout.Duration()
}

func (c *StarknetConfig) TxSendFrequency() time.Duration {
	return c.Chain.TxSendFrequency.Duration()
}

func (c *StarknetConfig) TxMaxBatchSize() int {
	return int(*c.Chain.TxMaxBatchSize)
}

func (c *StarknetConfig) OCR2CachePollPeriod() time.Duration {
	return c.Chain.OCR2CachePollPeriod.Duration()
}

func (c *StarknetConfig) OCR2CacheTTL() time.Duration {
	return c.Chain.OCR2CacheTTL.Duration()
}

func (c *StarknetConfig) RequestTimeout() time.Duration {
	return c.Chain.RequestTimeout.Duration()
}

func (c *StarknetConfig) Update(cfg db.ChainCfg) {
	panic(fmt.Errorf("cannot update: %v", v2.ErrUnsupported))
}
