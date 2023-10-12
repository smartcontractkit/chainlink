package cosmos

import (
	"fmt"
	"net/url"
	"slices"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pelletier/go-toml/v2"
	"github.com/shopspring/decimal"
	"go.uber.org/multierr"

	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/utils/config"
)

type CosmosConfigs []*CosmosConfig

func (cs CosmosConfigs) validateKeys() (err error) {
	// Unique chain IDs
	chainIDs := config.UniqueStrings{}
	for i, c := range cs {
		if chainIDs.IsDupe(c.ChainID) {
			err = multierr.Append(err, config.NewErrDuplicate(fmt.Sprintf("%d.ChainID", i), *c.ChainID))
		}
	}

	// Unique node names
	names := config.UniqueStrings{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			if names.IsDupe(n.Name) {
				err = multierr.Append(err, config.NewErrDuplicate(fmt.Sprintf("%d.Nodes.%d.Name", i, j), *n.Name))
			}
		}
	}

	// Unique TendermintURLs
	urls := config.UniqueStrings{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			u := (*url.URL)(n.TendermintURL)
			if urls.IsDupeFmt(u) {
				err = multierr.Append(err, config.NewErrDuplicate(fmt.Sprintf("%d.Nodes.%d.TendermintURL", i, j), u.String()))
			}
		}
	}
	return

}

func (cs CosmosConfigs) ValidateConfig() (err error) {
	return cs.validateKeys()
}

func (cs *CosmosConfigs) SetFrom(fs *CosmosConfigs) (err error) {
	if err1 := fs.validateKeys(); err1 != nil {
		return err1
	}
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
	return
}

func nodeStatus(n *coscfg.Node, id relay.ChainID) (relaytypes.NodeStatus, error) {
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
	// Do not access directly. Use [IsEnabled]
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
	if f.Bech32Prefix != nil {
		c.Bech32Prefix = f.Bech32Prefix
	}
	if f.BlockRate != nil {
		c.BlockRate = f.BlockRate
	}
	if f.BlocksUntilTxTimeout != nil {
		c.BlocksUntilTxTimeout = f.BlocksUntilTxTimeout
	}
	if f.ConfirmPollPeriod != nil {
		c.ConfirmPollPeriod = f.ConfirmPollPeriod
	}
	if f.FallbackGasPrice != nil {
		c.FallbackGasPrice = f.FallbackGasPrice
	}
	if f.GasToken != nil {
		c.GasToken = f.GasToken
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
	if f.OCR2CacheTTL != nil {
		c.OCR2CacheTTL = f.OCR2CacheTTL
	}
	if f.TxMsgTimeout != nil {
		c.TxMsgTimeout = f.TxMsgTimeout
	}
}

func (c *CosmosConfig) ValidateConfig() (err error) {
	if c.ChainID == nil {
		err = multierr.Append(err, config.ErrMissing{Name: "ChainID", Msg: "required for all chains"})
	} else if *c.ChainID == "" {
		err = multierr.Append(err, config.ErrEmpty{Name: "ChainID", Msg: "required for all chains"})
	}

	if len(c.Nodes) == 0 {
		err = multierr.Append(err, config.ErrMissing{Name: "Nodes", Msg: "must have at least one node"})
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

func (c *CosmosConfig) Bech32Prefix() string {
	return *c.Chain.Bech32Prefix
}

func (c *CosmosConfig) BlockRate() time.Duration {
	return c.Chain.BlockRate.Duration()
}

func (c *CosmosConfig) BlocksUntilTxTimeout() int64 {
	return *c.Chain.BlocksUntilTxTimeout
}

func (c *CosmosConfig) ConfirmPollPeriod() time.Duration {
	return c.Chain.ConfirmPollPeriod.Duration()
}

func (c *CosmosConfig) FallbackGasPrice() sdk.Dec {
	return sdkDecFromDecimal(c.Chain.FallbackGasPrice)
}

func (c *CosmosConfig) GasToken() string {
	return *c.Chain.GasToken
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

func (c *CosmosConfig) GetNode(name string) (db.Node, error) {
	for _, n := range c.Nodes {
		if *n.Name == name {
			return legacyNode(n, *c.ChainID), nil
		}
	}
	return db.Node{}, fmt.Errorf("%w: node %q", chains.ErrNotFound, name)
}

func (c *CosmosConfig) ListNodes() ([]db.Node, error) {
	var allNodes []db.Node
	for _, n := range c.Nodes {
		allNodes = append(allNodes, legacyNode(n, *c.ChainID))
	}
	return allNodes, nil
}
