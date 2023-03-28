package presenters

import (
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
)

// SolanaChainResource is an Solana chain JSONAPI resource.
type SolanaChainResource struct {
	ChainResource
}

// GetName implements the api2go EntityNamer interface
func (r SolanaChainResource) GetName() string {
	return "solana_chain"
}

// NewSolanaChainResource returns a new SolanaChainResource for chain.
func NewSolanaChainResource(chain chains.ChainConfig) SolanaChainResource {
	return SolanaChainResource{ChainResource{
		JAID:    NewJAID(chain.ID),
		Config:  chain.Cfg,
		Enabled: chain.Enabled,
	}}
}

// SolanaNodeResource is a Solana node JSONAPI resource.
type SolanaNodeResource struct {
	JAID
	Name          string `json:"name"`
	SolanaChainID string `json:"solanaChainID"`
	SolanaURL     string `json:"solanaURL"`
}

// GetName implements the api2go EntityNamer interface
func (r SolanaNodeResource) GetName() string {
	return "solana_node"
}

// NewSolanaNodeResource returns a new SolanaNodeResource for node.
func NewSolanaNodeResource(node db.Node) SolanaNodeResource {
	return SolanaNodeResource{
		JAID:          NewJAID(node.Name),
		Name:          node.Name,
		SolanaChainID: node.SolanaChainID,
		SolanaURL:     node.SolanaURL,
	}
}
