package cosmos

import (
	"fmt"
	"net/url"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"
	"golang.org/x/exp/slices"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-terra/pkg/cosmos"
	coscfg "github.com/smartcontractkit/chainlink-terra/pkg/cosmos/config"
	"github.com/smartcontractkit/chainlink-terra/pkg/cosmos/db"
	"github.com/smartcontractkit/chainlink/core/chains/cosmos/types"
	v2 "github.com/smartcontractkit/chainlink/core/config/v2"
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

func (cs CosmosConfigs) Chains(ids ...string) (chains []types.DBChain) {
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

func (cs CosmosConfigs) Node(name string) (n db.Node, err error) {
	for i := range cs {
		for _, n := range cs[i].Nodes {
			if n.Name != nil && *n.Name == name {
				return legacyNode(n, *cs[i].ChainID), nil
			}
		}
	}
	return
}

func (cs CosmosConfigs) Nodes() (ns []db.Node) {
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

func (cs CosmosConfigs) NodesByID(chainIDs ...string) (ns []db.Node) {
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

func (c *CosmosConfig) SetFromDB(ch types.DBChain, nodes []db.Node) error {
	c.ChainID = &ch.ID
	c.Enabled = &ch.Enabled

	if err := c.Chain.SetFromDB(ch.Cfg); err != nil {
		return err
	}
	for _, db := range nodes {
		var n coscfg.Node
		if err := n.SetFromDB(db); err != nil {
			return err
		}
		c.Nodes = append(c.Nodes, &n)
	}
	return nil
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

func (c *CosmosConfig) AsV1() types.DBChain {
	return types.DBChain{
		ID:      *c.ChainID,
		Enabled: c.IsEnabled(),
		Cfg: &db.ChainCfg{
			BlockRate:             c.Chain.BlockRate,
			BlocksUntilTxTimeout:  null.IntFromPtr(c.Chain.BlocksUntilTxTimeout),
			ConfirmPollPeriod:     c.Chain.ConfirmPollPeriod,
			FallbackGasPriceUAtom: nullString(c.Chain.FallbackGasPriceUAtom),
			FCDURL:                nullString((*url.URL)(c.Chain.FCDURL)),
			GasLimitMultiplier:    null.FloatFrom(c.Chain.GasLimitMultiplier.InexactFloat64()),
			MaxMsgsPerBatch:       null.IntFromPtr(c.Chain.MaxMsgsPerBatch),
			OCR2CachePollPeriod:   c.Chain.OCR2CachePollPeriod,
			OCR2CacheTTL:          c.Chain.OCR2CacheTTL,
			TxMsgTimeout:          c.Chain.TxMsgTimeout,
		},
	}
}

func nullString(s fmt.Stringer) null.String {
	if s == nil {
		return null.String{}
	}
	return null.StringFrom(s.String())
}

var _ cosmos.Config = &CosmosConfig{}

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

func (c *CosmosConfig) Update(cfg db.ChainCfg) {
	panic(fmt.Errorf("cannot update: %v", v2.ErrUnsupported))
}

func sdkDecFromDecimal(d *decimal.Decimal) sdk.Dec {
	i := d.Shift(sdk.Precision)
	return sdk.NewDecFromBigIntWithPrec(i.BigInt(), sdk.Precision)
}
