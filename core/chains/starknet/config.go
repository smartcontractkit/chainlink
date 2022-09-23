package starknet

import (
	"fmt"
	"net/url"
	"time"

	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v4"

	stkcfg "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"
	starknetdb "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/core/chains/starknet/types"
	v2 "github.com/smartcontractkit/chainlink/core/config/v2"
)

type StarknetConfigs []*StarknetConfig

func (cs StarknetConfigs) ValidateConfig() (err error) {
	chainIDs := map[string]struct{}{}
	for i, c := range cs {
		if c.ChainID == nil {
			continue
		}
		chainID := *c.ChainID
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

	// Unique URLs
	urls := map[string]struct{}{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			if n.URL == nil {
				continue
			}
			us := (*url.URL)(n.URL).String()
			if _, ok := urls[us]; ok {
				err = multierr.Append(err, v2.ErrInvalid{Name: fmt.Sprintf("%d.Nodes.%d.URL", i, j), Msg: "duplicate - must be unique", Value: us})
			}
			urls[us] = struct{}{}
		}
	}
	return
}

func (cs StarknetConfigs) Chains(ids ...string) (chains []types.DBChain) {
	for _, ch := range cs {
		if ch == nil {
			continue
		}
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
			if id != *cs[i].ChainID {
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

	return
}

func (c *StarknetConfig) AsV1() types.DBChain {
	return types.DBChain{
		ID:      *c.ChainID,
		Enabled: *c.Enabled,
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
