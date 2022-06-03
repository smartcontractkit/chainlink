package chainlink

import (
	soldb "github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	terdb "github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	evmcfg "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	evmtyp "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/chains/solana"
	solcfg "github.com/smartcontractkit/chainlink/core/chains/solana/config"
	tercfg "github.com/smartcontractkit/chainlink/core/chains/terra/config"
	tertyp "github.com/smartcontractkit/chainlink/core/chains/terra/types"
	config "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Config is TODO doc
//
// When adding a new field:
// 	- consider including a unit suffix with the field name
// 	- TOML is limited to int64/float64, so fields requiring greater range/precision must use non-standard types
//	implementing encoding.TextMarshaler/TextUnmarshaler
type Config struct {
	config.Core

	EVM []EVMConfig `toml:",omitempty"`

	Solana []SolanaConfig `toml:",omitempty"`

	Terra []TerraConfig `toml:",omitempty"`
}

type EVMConfig struct {
	ChainID *utils.Big
	Enabled *bool // TODO or find a way to output toml comments....
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
