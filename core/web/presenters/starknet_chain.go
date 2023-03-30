package presenters

import (
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
)

// StarkNetChainResource is an StarkNet chain JSONAPI resource.
type StarkNetChainResource struct {
	ChainResource
}

// GetName implements the api2go EntityNamer interface
func (r StarkNetChainResource) GetName() string {
	return "starknet_chain"
}

// NewStarkNetChainResource returns a new StarkNetChainResource for chain.
func NewStarkNetChainResource(chain chains.ChainConfig) StarkNetChainResource {
	return StarkNetChainResource{ChainResource{
		JAID:    NewJAID(chain.ID),
		Config:  chain.Cfg,
		Enabled: chain.Enabled,
	}}
}

// StarkNetNodeResource is a StarkNet node JSONAPI resource.
type StarkNetNodeResource struct {
	JAID
	Name    string `json:"name"`
	ChainID string `json:"chainID"`
	URL     string `json:"url"`
}

// GetName implements the api2go EntityNamer interface
func (r StarkNetNodeResource) GetName() string {
	return "starknet_node"
}

// NewStarkNetNodeResource returns a new StarkNetNodeResource for node.
func NewStarkNetNodeResource(node db.Node) StarkNetNodeResource {
	return StarkNetNodeResource{
		JAID:    NewJAID(node.Name),
		Name:    node.Name,
		ChainID: node.ChainID,
		URL:     node.URL,
	}
}
