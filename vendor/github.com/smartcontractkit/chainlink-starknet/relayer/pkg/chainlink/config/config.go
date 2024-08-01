package config

import (
	"errors"
	"fmt"
	"net/url"
	"slices"
	"time"

	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/ocr2"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/txm"
)

var DefaultConfigSet = ConfigSet{
	OCR2CachePollPeriod: 5 * time.Second,
	OCR2CacheTTL:        time.Minute,
	RequestTimeout:      10 * time.Second,
	TxTimeout:           10 * time.Second,
	ConfirmationPoll:    5 * time.Second,
}

type ConfigSet struct { //nolint:revive
	OCR2CachePollPeriod time.Duration
	OCR2CacheTTL        time.Duration

	// client config
	RequestTimeout time.Duration

	// txm config
	TxTimeout        time.Duration
	ConfirmationPoll time.Duration
}

type Config interface {
	txm.Config // txm config

	// ocr2 config
	ocr2.Config

	// client config
	RequestTimeout() time.Duration
}

type Chain struct {
	OCR2CachePollPeriod *config.Duration
	OCR2CacheTTL        *config.Duration
	RequestTimeout      *config.Duration
	TxTimeout           *config.Duration
	ConfirmationPoll    *config.Duration
}

func (c *Chain) SetDefaults() {
	if c.OCR2CachePollPeriod == nil {
		c.OCR2CachePollPeriod = config.MustNewDuration(DefaultConfigSet.OCR2CachePollPeriod)
	}
	if c.OCR2CacheTTL == nil {
		c.OCR2CacheTTL = config.MustNewDuration(DefaultConfigSet.OCR2CacheTTL)
	}
	if c.RequestTimeout == nil {
		c.RequestTimeout = config.MustNewDuration(DefaultConfigSet.RequestTimeout)
	}
	if c.TxTimeout == nil {
		c.TxTimeout = config.MustNewDuration(DefaultConfigSet.TxTimeout)
	}
	if c.ConfirmationPoll == nil {
		c.ConfirmationPoll = config.MustNewDuration(DefaultConfigSet.ConfirmationPoll)
	}
}

type Node struct {
	Name *string
	URL  *config.URL
	// optional, only if rpc url needs api key passed in header
	APIKey *string
}

type TOMLConfigs []*TOMLConfig

func (cs TOMLConfigs) ValidateConfig() (err error) {
	return cs.validateKeys()
}

func (cs TOMLConfigs) validateKeys() (err error) {
	// Unique chain IDs
	chainIDs := config.UniqueStrings{}
	for i, c := range cs {
		if chainIDs.IsDupe(c.ChainID) {
			err = errors.Join(err, config.NewErrDuplicate(fmt.Sprintf("%d.ChainID", i), *c.ChainID))
		}
	}

	// Unique node names
	names := config.UniqueStrings{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			if names.IsDupe(n.Name) {
				err = errors.Join(err, config.NewErrDuplicate(fmt.Sprintf("%d.Nodes.%d.Name", i, j), *n.Name))
			}
		}
	}

	// Unique URLs
	urls := config.UniqueStrings{}
	for i, c := range cs {
		for j, n := range c.Nodes {
			u := (*url.URL)(n.URL)
			if urls.IsDupeFmt(u) {
				err = errors.Join(err, config.NewErrDuplicate(fmt.Sprintf("%d.Nodes.%d.URL", i, j), u.String()))
			}
		}
	}
	return
}

func (cs *TOMLConfigs) SetFrom(fs *TOMLConfigs) (err error) {
	if err1 := fs.validateKeys(); err1 != nil {
		return err1
	}
	for _, f := range *fs {
		if f.ChainID == nil {
			*cs = append(*cs, f)
		} else if i := slices.IndexFunc(*cs, func(c *TOMLConfig) bool {
			return c.ChainID != nil && *c.ChainID == *f.ChainID
		}); i == -1 {
			*cs = append(*cs, f)
		} else {
			(*cs)[i].SetFrom(f)
		}
	}
	return
}

type TOMLConfig struct {
	ChainID   *string
	FeederURL *config.URL
	// Do not access directly. Use [IsEnabled]
	Enabled *bool
	Chain
	Nodes Nodes
}

func (c *TOMLConfig) IsEnabled() bool {
	return c.Enabled == nil || *c.Enabled
}

func (c *TOMLConfig) SetFrom(f *TOMLConfig) {
	if f.ChainID != nil {
		c.ChainID = f.ChainID
	}
	if f.Enabled != nil {
		c.Enabled = f.Enabled
	}
	if f.FeederURL != nil {
		c.FeederURL = f.FeederURL
	}
	setFromChain(&c.Chain, &f.Chain)
	c.Nodes.SetFrom(&f.Nodes)
}

func setFromChain(c, f *Chain) {
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

func (c *TOMLConfig) ValidateConfig() (err error) {
	if c.ChainID == nil {
		err = errors.Join(err, config.ErrMissing{Name: "ChainID", Msg: "required for all chains"})
	} else if *c.ChainID == "" {
		err = errors.Join(err, config.ErrEmpty{Name: "ChainID", Msg: "required for all chains"})
	}

	if len(c.Nodes) == 0 {
		err = errors.Join(err, config.ErrMissing{Name: "Nodes", Msg: "must have at least one node"})
	}

	return
}

func (c *TOMLConfig) TOMLString() (string, error) {
	b, err := toml.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

type Nodes []*Node

func (ns *Nodes) SetFrom(fs *Nodes) {
	for _, f := range *fs {
		if f.Name == nil {
			*ns = append(*ns, f)
		} else if i := slices.IndexFunc(*ns, func(n *Node) bool {
			return n.Name != nil && *n.Name == *f.Name
		}); i == -1 {
			*ns = append(*ns, f)
		} else {
			setFromNode((*ns)[i], f)
		}
	}
}

func setFromNode(n, f *Node) {
	if f.Name != nil {
		n.Name = f.Name
	}
	if f.URL != nil {
		n.URL = f.URL
	}
}

func legacyNode(n *Node, id string) db.Node {
	var apiKey string
	if n.APIKey == nil {
		apiKey = ""
	} else {
		apiKey = *n.APIKey
	}
	return db.Node{
		Name:    *n.Name,
		ChainID: id,
		URL:     (*url.URL)(n.URL).String(),
		APIKey:  apiKey,
	}
}

var _ Config = &TOMLConfig{}

func (c *TOMLConfig) TxTimeout() time.Duration {
	return c.Chain.TxTimeout.Duration()
}

func (c *TOMLConfig) ConfirmationPoll() time.Duration {
	return c.Chain.ConfirmationPoll.Duration()
}

func (c *TOMLConfig) OCR2CachePollPeriod() time.Duration {
	return c.Chain.OCR2CachePollPeriod.Duration()
}

func (c *TOMLConfig) OCR2CacheTTL() time.Duration {
	return c.Chain.OCR2CacheTTL.Duration()
}

func (c *TOMLConfig) RequestTimeout() time.Duration {
	return c.Chain.RequestTimeout.Duration()
}

func (c *TOMLConfig) ListNodes() ([]db.Node, error) {
	var allNodes []db.Node
	for _, n := range c.Nodes {
		allNodes = append(allNodes, legacyNode(n, *c.ChainID))
	}
	return allNodes, nil
}
