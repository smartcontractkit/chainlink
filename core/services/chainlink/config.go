package chainlink

import (
	solanadb "github.com/smartcontractkit/chainlink-solana/pkg/solana/db"
	terradb "github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/chains/solana"
	"github.com/smartcontractkit/chainlink/core/chains/terra"
	terratypes "github.com/smartcontractkit/chainlink/core/chains/terra/types"
	"github.com/smartcontractkit/chainlink/core/config/toml"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// Config is TODO doc
//
// When adding a new field:
// 	- consider including a unit suffix with the field name
// 	- TOML is limited to int64/float64, so fields requiring greater range/precision must use non-standard types
//	implementing encoding.TextMarshaler/TextUnmarshaler
type Config struct {
	toml.CoreConfig

	EVM []EVMConfig `toml:",omitempty"`

	Solana []SolanaConfig `toml:",omitempty"`

	Terra []TerraConfig `toml:",omitempty"`
}

type EVMConfig struct {
	ChainID *utils.Big
	Enabled *bool // TODO or find a way to output toml comments....
	evmtypes.ChainTOMLCfg
	Nodes []evmtypes.TOMLNode
}

func newEVMConfigFromDB(ch evmtypes.DBChain, nodes []evmtypes.Node) (EVMConfig, error) {
	c := EVMConfig{
		ChainID: &ch.ID,
		Enabled: &ch.Enabled,
	}
	if err := c.ChainTOMLCfg.SetFromDB(ch.Cfg); err != nil {
		return EVMConfig{}, err
	}
	for _, db := range nodes {
		n, err := evmtypes.NewTOMLNodeFromDB(db)
		if err != nil {
			return EVMConfig{}, err
		}
		c.Nodes = append(c.Nodes, n)
	}
	return c, nil
}

type SolanaConfig struct {
	ChainID string
	Enabled *bool
	solana.TOMLChain
	Nodes []solana.TOMLNode
}

func newSolanaConfigFromDB(ch solana.DBChain, nodes []solanadb.Node) (SolanaConfig, error) {
	c := SolanaConfig{
		ChainID: ch.ID,
		Enabled: &ch.Enabled,
	}
	if err := c.TOMLChain.SetFromDB(ch.Cfg); err != nil {
		return SolanaConfig{}, err
	}
	for _, db := range nodes {
		n, err := solana.NewTOMLNodeFromDB(db)
		if err != nil {
			return SolanaConfig{}, err
		}
		c.Nodes = append(c.Nodes, n)
	}
	return c, nil
}

type TerraConfig struct {
	ChainID string
	Enabled *bool
	terra.TOMLChain
	Nodes []terra.TOMLNode
}

func newTerraConfigFromDB(ch terratypes.DBChain, nodes []terradb.Node) (TerraConfig, error) {
	c := TerraConfig{
		ChainID: ch.ID,
		Enabled: &ch.Enabled,
	}
	if err := c.TOMLChain.SetFromDB(ch.Cfg); err != nil {
		return TerraConfig{}, err
	}
	for _, db := range nodes {
		n, err := terra.NewTOMLNodeFromDB(db)
		if err != nil {
			return TerraConfig{}, err
		}
		c.Nodes = append(c.Nodes, n)
	}
	return c, nil
}
