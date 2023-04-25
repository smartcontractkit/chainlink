package cosmos

import (
	"fmt"
	"net/url"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pelletier/go-toml/v2"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"
	"golang.org/x/exp/slices"

	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos/types"
	v2 "github.com/smartcontractkit/chainlink/v2/core/config/v2"
)

type CosmosConfigs []*CosmosConfig

func (cs CosmosConfigs) ValidateConfig() (err error) {
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

	// Unique TendermintURLs
	urls := v2.UniqueStrings{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			u := (*url.URL)(n.TendermintURL)
			if urls.IsDupeFmt(u) {
				err = multierr.Append(err, v2.NewErrDuplicate(fmt.Sprintf("%d.Nodes.%d.TendermintURL", i, j), u.String()))
			}
		}
	}
	return
}

func (cs *CosmosConfigs) SetFrom(fs *CosmosConfigs) {
	for _, f := range *fs {
		if f.ChainID == nil {
			*cs = append(*cs, f)
		} else if i := slices.IndexFunc(*cs, func(c *CosmosConfig) bool {
			return c.ChainID != nil && *c.ChainID == *f.ChainID
		}); i == -1 {
			*cs = append(*cs, f)
		} else {
			(*cs)[i].SetFrom(f)
		}
	}
}

func (cs CosmosConfigs) Chains(ids ...string) (r []relaytypes.ChainStatus, err error) {
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
		ch2 := relaytypes.ChainStatus{
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

func (cs CosmosConfigs) Node(name string) (n db.Node, err error) {
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

func (cs CosmosConfigs) nodes(chainID string) (ns CosmosNodes) {
	for _, c := range cs {
		if *c.ChainID == chainID {
			return c.Nodes
		}
	}
	return nil
}

func (cs CosmosConfigs) Nodes(chainID string) (ns []db.Node, err error) {
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

func (cs CosmosConfigs) NodeStatus(name string) (n relaytypes.NodeStatus, err error) {
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

func (cs CosmosConfigs) NodeStatuses(chainIDs ...string) (ns []relaytypes.NodeStatus, err error) {
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

func nodeStatus(n *coscfg.Node, chainID string) (relaytypes.NodeStatus, error) {
	var s relaytypes.NodeStatus
	s.ChainID = chainID
	s.Name = *n.Name
	b, err := toml.Marshal(n)
	if err != nil {
		return relaytypes.NodeStatus{}, err
	}
	s.Config = string(b)
	return s, nil
}

type CosmosNodes []*coscfg.Node

func (ns *CosmosNodes) SetFrom(fs *CosmosNodes) {
	for _, f := range *fs {
		if f.Name == nil {
			*ns = append(*ns, f)
		} else if i := slices.IndexFunc(*ns, func(n *coscfg.Node) bool {
			return n.Name != nil && *n.Name == *f.Name
		}); i == -1 {
			*ns = append(*ns, f)
		} else {
			setFromNode((*ns)[i], f)
		}
	}
}

func setFromNode(n, f *coscfg.Node) {
	if f.Name != nil {
		n.Name = f.Name
	}
	if f.TendermintURL != nil {
		n.TendermintURL = f.TendermintURL
	}
}

func legacyNode(n *coscfg.Node, id string) db.Node {
	return db.Node{
		Name:          *n.Name,
		CosmosChainID: id,
		TendermintURL: (*url.URL)(n.TendermintURL).String(),
	}
}

type CosmosConfig struct {
	ChainID *string
	Enabled *bool
	coscfg.Chain
	Nodes CosmosNodes
}

func (c *CosmosConfig) IsEnabled() bool {
	return c.Enabled == nil || *c.Enabled
}

func (c *CosmosConfig) SetFrom(f *CosmosConfig) {
	if f.ChainID != nil {
		c.ChainID = f.ChainID
	}
	if f.Enabled != nil {
		c.Enabled = f.Enabled
	}
	setFromChain(&c.Chain, &f.Chain)
	c.Nodes.SetFrom(&f.Nodes)
}

func setFromChain(c, f *coscfg.Chain) {
	if f.BlockRate != nil {
		c.BlockRate = f.BlockRate
	}
	if f.BlocksUntilTxTimeout != nil {
		c.BlocksUntilTxTimeout = f.BlocksUntilTxTimeout
	}
	if f.ConfirmPollPeriod != nil {
		c.ConfirmPollPeriod = f.ConfirmPollPeriod
	}
	if f.FallbackGasPriceUAtom != nil {
		c.FallbackGasPriceUAtom = f.FallbackGasPriceUAtom
	}
	if f.FCDURL != nil {
		c.FCDURL = f.FCDURL
	}
	if f.GasLimitMultiplier != nil {
		c.GasLimitMultiplier = f.GasLimitMultiplier
	}
	if f.MaxMsgsPerBatch != nil {
		c.MaxMsgsPerBatch = f.MaxMsgsPerBatch
	}
	if f.OCR2CachePollPeriod != nil {
		c.OCR2CachePollPeriod = f.OCR2CachePollPeriod
	}
	if f.BlockRate != nil {
		c.BlockRate = f.BlockRate
	}
	if f.BlockRate != nil {
		c.BlockRate = f.BlockRate
	}
}

func (c *CosmosConfig) ValidateConfig() (err error) {
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

func (c *CosmosConfig) TOMLString() (string, error) {
	b, err := toml.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

var _ coscfg.Config = &CosmosConfig{}

func (c *CosmosConfig) BlockRate() time.Duration {
	return c.Chain.BlockRate.Duration()
}

func (c *CosmosConfig) BlocksUntilTxTimeout() int64 {
	return *c.Chain.BlocksUntilTxTimeout
}

func (c *CosmosConfig) ConfirmPollPeriod() time.Duration {
	return c.Chain.ConfirmPollPeriod.Duration()
}

func (c *CosmosConfig) FallbackGasPriceUAtom() sdk.Dec {
	return sdkDecFromDecimal(c.Chain.FallbackGasPriceUAtom)
}

func (c *CosmosConfig) FCDURL() url.URL {
	return (url.URL)(*c.Chain.FCDURL)
}

func (c *CosmosConfig) GasLimitMultiplier() float64 {
	return c.Chain.GasLimitMultiplier.InexactFloat64()
}

func (c *CosmosConfig) MaxMsgsPerBatch() int64 {
	return *c.Chain.MaxMsgsPerBatch
}

func (c *CosmosConfig) OCR2CachePollPeriod() time.Duration {
	return c.Chain.OCR2CachePollPeriod.Duration()
}

func (c *CosmosConfig) OCR2CacheTTL() time.Duration {
	return c.Chain.OCR2CacheTTL.Duration()
}

func (c *CosmosConfig) TxMsgTimeout() time.Duration {
	return c.Chain.TxMsgTimeout.Duration()
}

func sdkDecFromDecimal(d *decimal.Decimal) sdk.Dec {
	i := d.Shift(sdk.Precision)
	return sdk.NewDecFromBigIntWithPrec(i.BigInt(), sdk.Precision)
}

func NewConfigs(cfgs chains.ConfigsV2[string, db.Node]) types.Configs {
	return chains.NewConfigs(cfgs)
}
