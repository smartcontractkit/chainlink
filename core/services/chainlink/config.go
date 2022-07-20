package chainlink

import (
	"strings"

	"github.com/pelletier/go-toml/v2"

	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	soldb "github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	tercfg "github.com/smartcontractkit/chainlink-terra/pkg/terra/config"
	terdb "github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	evmcfg "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	evmtyp "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/chains/solana"
	tertyp "github.com/smartcontractkit/chainlink/core/chains/terra/types"
	config "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Config is the root type used for TOML configuration.
//
// When adding a new field:
// 	- consider including a unit suffix with the field name
// 	- TOML is limited to int64/float64, so fields requiring greater range/precision must use non-standard types
//	implementing encoding.TextMarshaler/TextUnmarshaler, like utils.Big and decimal.Decimal
//  - std lib types that don't implement encoding.TextMarshaler/TextUnmarshaler (time.Duration, url.URL, big.Int) won't
//   work as expected, and require wrapper types. See models.Duration, models.URL, utils.Big.
type Config struct {
	config.Core

	EVM []EVMConfig `toml:",omitempty"`

	Solana []SolanaConfig `toml:",omitempty"`

	Terra []TerraConfig `toml:",omitempty"`
}

// TOMLString returns a pretty-printed TOML encoded string, with extra line breaks removed.
func (c *Config) TOMLString() (string, error) {
	b, err := toml.Marshal(c)
	if err != nil {
		return "", err
	}
	// remove runs of line breaks
	s := multiLineBreak.ReplaceAllLiteralString(string(b), "\n")
	// restore them preceding keys
	s = strings.Replace(s, "\n[", "\n\n[", -1)
	s = strings.TrimPrefix(s, "\n")
	return s, nil
}

type EVMConfig struct {
	ChainID *utils.Big
	Enabled *bool
	evmcfg.Chain
	Nodes []evmcfg.Node
}

func (c *EVMConfig) setFromDB(ch evmtyp.DBChain, nodes []evmtyp.Node) error {
	c.ChainID = &ch.ID
	c.Enabled = &ch.Enabled

	if err := c.Chain.SetFromDB(ch.Cfg); err != nil {
		return err
	}
	for _, db := range nodes {
		var n evmcfg.Node
		if err := n.SetFromDB(db); err != nil {
			return err
		}
		c.Nodes = append(c.Nodes, n)
	}
	return nil
}

type SolanaConfig struct {
	ChainID string
	Enabled *bool
	solcfg.Chain
	Nodes []solcfg.Node
}

func (c *SolanaConfig) setFromDB(ch solana.DBChain, nodes []soldb.Node) error {
	c.ChainID = ch.ID
	c.Enabled = &ch.Enabled

	if err := c.Chain.SetFromDB(ch.Cfg); err != nil {
		return err
	}
	for _, db := range nodes {
		var n solcfg.Node
		if err := n.SetFromDB(db); err != nil {
			return err
		}
		c.Nodes = append(c.Nodes, n)
	}
	return nil
}

type TerraConfig struct {
	ChainID string
	Enabled *bool
	tercfg.Chain
	Nodes []tercfg.Node
}

func (c *TerraConfig) setFromDB(ch tertyp.DBChain, nodes []terdb.Node) error {
	c.ChainID = ch.ID
	c.Enabled = &ch.Enabled

	if err := c.Chain.SetFromDB(ch.Cfg); err != nil {
		return err
	}
	for _, db := range nodes {
		var n tercfg.Node
		if err := n.SetFromDB(db); err != nil {
			return err
		}
		c.Nodes = append(c.Nodes, n)
	}
	return nil
}
