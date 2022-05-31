package chainlink

import (
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/chains/solana"
	"github.com/smartcontractkit/chainlink/core/chains/terra"
	"github.com/smartcontractkit/chainlink/core/config/toml"
)

// Config is TODO doc
//
// When adding a new field:
// 	- consider including a unit suffix with the field name
// 	- TOML is limited to int64/float64, so fields requiring greating range/precision must use non-standard types
//	implementing encoding.TextMarshaler/TextUnmarshaler
type Config struct {
	toml.CoreConfig

	EVM []EVMConfig `toml:",omitempty"`

	Solana []SolanaConfig `toml:",omitempty"`

	Terra []TerraConfig `toml:",omitempty"`
}

type EVMConfig struct {
	ChainID int64 `toml:",omitempty"` //TODO big.Int?
	//TODO Enabled bool?
	evmtypes.ChainTOMLCfg
	Nodes []evmtypes.TOMLNode
}

type SolanaConfig struct {
	ChainID string `toml:",omitempty"`
	solana.TOMLChain
	Nodes []solana.TOMLNode
}

type TerraConfig struct {
	ChainID string `toml:",omitempty"`
	terra.TOMLChain
	Nodes []terra.TOMLNode
}
