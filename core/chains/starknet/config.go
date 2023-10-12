package starknet

import (
	"fmt"
	"net/url"
	"slices"
	"time"

	"github.com/pelletier/go-toml/v2"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	stkcfg "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/utils/config"
)

type StarknetConfigs []*StarknetConfig

func (cs StarknetConfigs) ValidateConfig() (err error) {
	return cs.validateKeys()
}

func (cs StarknetConfigs) validateKeys() (err error) {
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

	// Unique URLs
	urls := config.UniqueStrings{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			u := (*url.URL)(n.URL)
			if urls.IsDupeFmt(u) {
				err = multierr.Append(err, config.NewErrDuplicate(fmt.Sprintf("%d.Nodes.%d.URL", i, j), u.String()))
			}
		}
	}
	return
}

func (cs *StarknetConfigs) SetFrom(fs *StarknetConfigs) (err error) {
	if err1 := fs.validateKeys(); err1 != nil {
		return err1
	}
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
	return
}

func nodeStatus(n *stkcfg.Node, id relay.ChainID) (types.NodeStatus, error) {
	var s types.NodeStatus
	s.ChainID = id
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
	// Do not access directly. Use [IsEnabled]
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
	if f.ConfirmationPoll != nil {
		c.ConfirmationPoll = f.ConfirmationPoll
	}
}

func (c *StarknetConfig) ValidateConfig() (err error) {
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

func legacyNode(n *stkcfg.Node, id relay.ChainID) db.Node {
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

func (c *StarknetConfig) ConfirmationPoll() time.Duration {
	return c.Chain.ConfirmationPoll.Duration()
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

func (c *StarknetConfig) ListNodes() ([]db.Node, error) {
	var allNodes []db.Node
	for _, n := range c.Nodes {
		allNodes = append(allNodes, legacyNode(n, *c.ChainID))
	}
	return allNodes, nil
}
