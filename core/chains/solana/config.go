package solana

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pelletier/go-toml/v2"
	"go.uber.org/multierr"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink-relay/pkg/utils"

	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	soldb "github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/utils/config"
)

type SolanaConfigs []*SolanaConfig

func (cs SolanaConfigs) ValidateConfig() (err error) {
	return cs.validateKeys()
}

func (cs SolanaConfigs) validateKeys() (err error) {
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
		for j, n := range c.nodes {
			if names.IsDupe(n.Name) {
				err = multierr.Append(err, config.NewErrDuplicate(fmt.Sprintf("%d.Nodes.%d.Name", i, j), *n.Name))
			}
		}
	}

	// Unique URLs
	urls := config.UniqueStrings{}
	for i, c := range cs {
		for j, n := range c.nodes {
			u := (*url.URL)(n.URL)
			if urls.IsDupeFmt(u) {
				err = multierr.Append(err, config.NewErrDuplicate(fmt.Sprintf("%d.Nodes.%d.URL", i, j), u.String()))
			}
		}
	}
	return
}

func (cs *SolanaConfigs) SetFrom(fs *SolanaConfigs) (err error) {
	if err1 := fs.validateKeys(); err1 != nil {
		return err1
	}
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
	return
}
func nodeStatus(n *solcfg.Node, chainID string) (types.NodeStatus, error) {
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

func nodeFromLegacy(l soldb.Node) (*solcfg.Node, error) {
	url, err := utils.ParseURL(l.SolanaURL)
	if err != nil {
		return nil, err
	}
	return &solcfg.Node{
		Name: &l.Name,
		URL:  url,
	}, nil
}

type SolanaConfig struct {
	ChainID *string
	Enabled *bool
	solcfg.Chain
	nodes SolanaNodes
}

func (c *SolanaConfig) ChainStatus() (stat types.ChainStatus, err error) {

	stat = types.ChainStatus{
		ID:      *c.ChainID,
		Enabled: c.IsEnabled(),
	}
	stat.Config, err = c.TOMLString()
	if err != nil {
		return stat, err
	}

	return stat, nil
}

func (c *SolanaConfig) Nodes(names ...string) (nodes []soldb.Node, err error) {
	filter := func(filterFn func(node *solcfg.Node) bool) {
		for _, n := range c.nodes {
			if n.Name != nil {
				if filterFn(n) {
					nodes = append(nodes, legacySolNode(n, *c.ChainID))
				}
			}
		}
	}
	if len(names) == 0 {
		allNodes := func(ignored *solcfg.Node) bool { return true }
		filter(allNodes)
	} else {
		for _, name := range names {
			matchName := func(node *solcfg.Node) bool { return name == *node.Name }
			filter(matchName)
		}
	}

	return nodes, nil
}

func (c SolanaConfig) NodeStatuses(names ...string) ([]types.NodeStatus, error) {
	legacyNodes, err := c.Nodes(names...)
	if err != nil {
		return nil, err
	}
	stats := make([]types.NodeStatus, 0)
	for _, legacyNodes := range legacyNodes {
		n, err2 := nodeFromLegacy(legacyNodes)
		if err2 != nil {
			err = errors.Join(err, err2)
			continue
		}

		stat, err2 := nodeStatus(n, *c.ChainID)
		if err2 != nil {
			err = errors.Join(err, err2)
			continue
		}
		stats = append(stats, stat)
	}
	return stats, err

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
	c.nodes.SetFrom(&f.nodes)
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

func (c *SolanaConfig) ValidateConfig() (err error) {
	if c.ChainID == nil {
		err = multierr.Append(err, config.ErrMissing{Name: "ChainID", Msg: "required for all chains"})
	} else if *c.ChainID == "" {
		err = multierr.Append(err, config.ErrEmpty{Name: "ChainID", Msg: "required for all chains"})
	}

	if len(c.nodes) == 0 {
		err = multierr.Append(err, config.ErrMissing{Name: "Nodes", Msg: "must have at least one node"})
	}
	return
}

func (c *SolanaConfig) TOMLString() (string, error) {
	b, err := toml.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
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

// ConfigStater manages solana chains and nodes.
type ConfigStater interface {
	chains.ChainConfig
	chains.NodeConfigs[soldb.Node]
}

var _ chains.Statuser[soldb.Node] = (ConfigStater)(nil)

func NewConfigStater(cfgs chains.ConfigsV2[soldb.Node]) ConfigStater {
	return chains.NewConfigs(cfgs)
}
