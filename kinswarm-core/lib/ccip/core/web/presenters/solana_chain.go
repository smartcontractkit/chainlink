package presenters

import (
	"github.com/smartcontractkit/chainlink-common/pkg/types"
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
func NewSolanaChainResource(chain types.ChainStatus) SolanaChainResource {
	return SolanaChainResource{ChainResource{
		JAID:    NewJAID(chain.ID),
		Config:  chain.Config,
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
func NewSolanaNodeResource(node types.NodeStatus) SolanaNodeResource {
	return SolanaNodeResource{NodeResource{
		JAID:    NewPrefixedJAID(node.Name, node.ChainID),
		ChainID: node.ChainID,
		Name:    node.Name,
		State:   node.State,
		Config:  node.Config,
	}}
}
