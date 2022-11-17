package terra

import (
	"fmt"
	"net/url"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"
	"golang.org/x/exp/slices"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra"
	tercfg "github.com/smartcontractkit/chainlink-terra/pkg/terra/config"
	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"
	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
	v2 "github.com/smartcontractkit/chainlink/core/config/v2"
)

type TerraConfigs []*TerraConfig

func (cs TerraConfigs) ValidateConfig() (err error) {
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

func (cs *TerraConfigs) SetFrom(fs *TerraConfigs) {
	for _, f := range *fs {
		if f.ChainID == nil {
			*cs = append(*cs, f)
		} else if i := slices.IndexFunc(*cs, func(c *TerraConfig) bool {
			return c.ChainID != nil && *c.ChainID == *f.ChainID
		}); i == -1 {
			*cs = append(*cs, f)
		} else {
			(*cs)[i].SetFrom(f)
		}
	}
}

func (cs TerraConfigs) Chains(ids ...string) (chains []types.DBChain) {
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

func (cs TerraConfigs) Node(name string) (n db.Node, err error) {
	for i := range cs {
		for _, n := range cs[i].Nodes {
			if n.Name != nil && *n.Name == name {
				return legacyNode(n, *cs[i].ChainID), nil
			}
		}
	}
	return
}

func (cs TerraConfigs) Nodes() (ns []db.Node) {
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

func (cs TerraConfigs) NodesByID(chainIDs ...string) (ns []db.Node) {
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

type TerraNodes []*tercfg.Node

func (ns *TerraNodes) SetFrom(fs *TerraNodes) {
	for _, f := range *fs {
		if f.Name == nil {
			*ns = append(*ns, f)
		} else if i := slices.IndexFunc(*ns, func(n *tercfg.Node) bool {
			return n.Name != nil && *n.Name == *f.Name
		}); i == -1 {
			*ns = append(*ns, f)
		} else {
			setFromNode((*ns)[i], f)
		}
	}
}

func setFromNode(n, f *tercfg.Node) {
	if f.Name != nil {
		n.Name = f.Name
	}
	if f.TendermintURL != nil {
		n.TendermintURL = f.TendermintURL
	}
}

func legacyNode(n *tercfg.Node, id string) db.Node {
	return db.Node{
		Name:          *n.Name,
		TerraChainID:  id,
		TendermintURL: (*url.URL)(n.TendermintURL).String(),
	}
}

type TerraConfig struct {
	ChainID *string
	Enabled *bool
	tercfg.Chain
	Nodes TerraNodes
}

func (c *TerraConfig) IsEnabled() bool {
	return c.Enabled == nil || *c.Enabled
}

func (c *TerraConfig) SetFrom(f *TerraConfig) {
	if f.ChainID != nil {
		c.ChainID = f.ChainID
	}
	if f.Enabled != nil {
		c.Enabled = f.Enabled
	}
	setFromChain(&c.Chain, &f.Chain)
	c.Nodes.SetFrom(&f.Nodes)
}

func setFromChain(c, f *tercfg.Chain) {
	if f.BlockRate != nil {
		c.BlockRate = f.BlockRate
	}
	if f.BlocksUntilTxTimeout != nil {
		c.BlocksUntilTxTimeout = f.BlocksUntilTxTimeout
	}
	if f.ConfirmPollPeriod != nil {
		c.ConfirmPollPeriod = f.ConfirmPollPeriod
	}
	if f.FallbackGasPriceULuna != nil {
		c.FallbackGasPriceULuna = f.FallbackGasPriceULuna
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

func (c *TerraConfig) SetFromDB(ch types.DBChain, nodes []db.Node) error {
	c.ChainID = &ch.ID
	c.Enabled = &ch.Enabled

	if err := c.Chain.SetFromDB(ch.Cfg); err != nil {
		return err
	}
	for _, db := range nodes {
		var n tercfg.Node
		if err := n.SetFromDB(db); err != nil {
			return err
		}
		c.Nodes = append(c.Nodes, &n)
	}
	return nil
}

func (c *TerraConfig) ValidateConfig() (err error) {
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

func (c *TerraConfig) AsV1() types.DBChain {
	return types.DBChain{
		ID:      *c.ChainID,
		Enabled: c.IsEnabled(),
		Cfg: &db.ChainCfg{
			BlockRate:             c.Chain.BlockRate,
			BlocksUntilTxTimeout:  null.IntFromPtr(c.Chain.BlocksUntilTxTimeout),
			ConfirmPollPeriod:     c.Chain.ConfirmPollPeriod,
			FallbackGasPriceULuna: nullString(c.Chain.FallbackGasPriceULuna),
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

var _ terra.Config = &TerraConfig{}

func (c *TerraConfig) BlockRate() time.Duration {
	return c.Chain.BlockRate.Duration()
}

func (c *TerraConfig) BlocksUntilTxTimeout() int64 {
	return *c.Chain.BlocksUntilTxTimeout
}

func (c *TerraConfig) ConfirmPollPeriod() time.Duration {
	return c.Chain.ConfirmPollPeriod.Duration()
}

func (c *TerraConfig) FallbackGasPriceULuna() sdk.Dec {
	return sdkDecFromDecimal(c.Chain.FallbackGasPriceULuna)
}

func (c *TerraConfig) FCDURL() url.URL {
	return (url.URL)(*c.Chain.FCDURL)
}

func (c *TerraConfig) GasLimitMultiplier() float64 {
	return c.Chain.GasLimitMultiplier.InexactFloat64()
}

func (c *TerraConfig) MaxMsgsPerBatch() int64 {
	return *c.Chain.MaxMsgsPerBatch
}

func (c *TerraConfig) OCR2CachePollPeriod() time.Duration {
	return c.Chain.OCR2CachePollPeriod.Duration()
}

func (c *TerraConfig) OCR2CacheTTL() time.Duration {
	return c.Chain.OCR2CacheTTL.Duration()
}

func (c *TerraConfig) TxMsgTimeout() time.Duration {
	return c.Chain.TxMsgTimeout.Duration()
}

func (c *TerraConfig) Update(cfg db.ChainCfg) {
	panic(fmt.Errorf("cannot update: %v", v2.ErrUnsupported))
}

func sdkDecFromDecimal(d *decimal.Decimal) sdk.Dec {
	i := d.Shift(sdk.Precision)
	return sdk.NewDecFromBigIntWithPrec(i.BigInt(), sdk.Precision)
}
