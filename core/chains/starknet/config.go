package starknet

import (
	"fmt"
	"net/url"
	"time"

	"github.com/pelletier/go-toml/v2"
	"go.uber.org/multierr"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	stkcfg "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
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

func (cs StarknetConfigs) Chains(ids ...string) (r []types.ChainStatus, err error) {
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
		ch2 := types.ChainStatus{
			ID:      *ch.ChainID,
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

func (cs StarknetConfigs) Node(name string) (n db.Node, err error) {
	for i := range cs {
		for _, n := range cs[i].Nodes {
			if n.Name != nil && *n.Name == name {
				return legacyNode(n, *cs[i].ChainID), nil
			}
		}
	}
	err = chains.ErrNotFound
	return
}

func (cs StarknetConfigs) nodes(chainID string) (ns StarknetNodes) {
	for _, c := range cs {
		if *c.ChainID == chainID {
			return c.Nodes
		}
	}
	return nil
}

func (cs StarknetConfigs) Nodes(chainID string) (ns []db.Node, err error) {
	nodes := cs.nodes(chainID)
	if nodes == nil {
		err = chains.ErrNotFound
		return
	}
	for _, n := range nodes {
		if n == nil {
			continue
		}
		ns = append(ns, legacyNode(n, chainID))
	}
	return
}

func (cs StarknetConfigs) NodeStatus(name string) (n types.NodeStatus, err error) {
	for i := range cs {
		for _, n := range cs[i].Nodes {
			if n.Name != nil && *n.Name == name {
				return nodeStatus(n, *cs[i].ChainID)
			}
		}
	}
	err = chains.ErrNotFound
	return
}

func (cs StarknetConfigs) NodeStatuses(chainIDs ...string) (ns []types.NodeStatus, err error) {
	if len(chainIDs) == 0 {
		for i := range cs {
			for _, n := range cs[i].Nodes {
				if n == nil {
					continue
				}
				n2, err := nodeStatus(n, *cs[i].ChainID)
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

func nodeStatus(n *stkcfg.Node, chainID string) (types.NodeStatus, error) {
	var s types.NodeStatus
	s.ChainID = chainID
	s.Name = *n.Name
	b, err := toml.Marshal(n)
	if err != nil {
		return types.NodeStatus{}, err
	}
	s.Config = string(b)
	return s, nil
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

func (c *StarknetConfig) TOMLString() (string, error) {
	b, err := toml.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
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
