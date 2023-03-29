package presenters

import (
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
)

// StarknetChainResource is an Starknet chain JSONAPI resource.
type StarknetChainResource struct {
	ChainResource
}

// GetName implements the api2go EntityNamer interface
func (r StarknetChainResource) GetName() string {
	return "starknet_chain"
}

// NewStarknetChainResource returns a new StarknetChainResource for chain.
func NewStarknetChainResource(chain chains.ChainConfig) StarknetChainResource {
	return StarknetChainResource{ChainResource{
		JAID:    NewJAID(chain.ID),
		Config:  chain.Cfg,
		Enabled: chain.Enabled,
	}}
}

// StarknetNodeResource is a Starknet node JSONAPI resource.
type StarknetNodeResource struct {
	JAID
	Name    string `json:"name"`
	ChainID string `json:"chainID"`
	URL     string `json:"url"`
}

// GetName implements the api2go EntityNamer interface
func (r StarknetNodeResource) GetName() string {
	return "starknet_node"
}

// NewStarknetNodeResource returns a new StarknetNodeResource for node.
func NewStarknetNodeResource(node db.Node) StarknetNodeResource {
	return StarknetNodeResource{
		JAID:    NewJAID(node.Name),
		Name:    node.Name,
		ChainID: node.ChainID,
		URL:     node.URL,
	}
}
