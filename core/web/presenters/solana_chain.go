package presenters

import (
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
	NodeResource
}

// GetName implements the api2go EntityNamer interface
func (r SolanaNodeResource) GetName() string {
	return "solana_node"
}

// NewSolanaNodeResource returns a new SolanaNodeResource for node.
func NewSolanaNodeResource(node chains.NodeStatus) SolanaNodeResource {
	return SolanaNodeResource{NodeResource{
		JAID:    NewJAID(node.Name),
		ChainID: node.ChainID,
		Name:    node.Name,
		State:   node.State,
		Config:  node.Config,
	}}
}
