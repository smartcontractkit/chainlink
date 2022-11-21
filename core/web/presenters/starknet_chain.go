package presenters

import (
	"time"

	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	starknet "github.com/smartcontractkit/chainlink/core/chains/starknet/types"
)

// StarkNetChainResource is an StarkNet chain JSONAPI resource.
type StarkNetChainResource struct {
	chainResource[*db.ChainCfg]
}

// GetName implements the api2go EntityNamer interface
func (r StarkNetChainResource) GetName() string {
	return "starknet_chain"
}

// NewStarkNetChainResource returns a new StarkNetChainResource for chain.
func NewStarkNetChainResource(chain starknet.DBChain) StarkNetChainResource {
	return StarkNetChainResource{chainResource[*db.ChainCfg]{
		JAID:      NewJAID(chain.ID),
		Config:    chain.Cfg,
		Enabled:   chain.Enabled,
		CreatedAt: chain.CreatedAt,
		UpdatedAt: chain.UpdatedAt,
	}}
}

// StarkNetNodeResource is a StarkNet node JSONAPI resource.
type StarkNetNodeResource struct {
	JAID
	Name      string    `json:"name"`
	ChainID   string    `json:"chainID"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// GetName implements the api2go EntityNamer interface
func (r StarkNetNodeResource) GetName() string {
	return "starknet_node"
}

// NewStarkNetNodeResource returns a new StarkNetNodeResource for node.
func NewStarkNetNodeResource(node db.Node) StarkNetNodeResource {
	return StarkNetNodeResource{
		JAID:      NewJAIDInt32(node.ID),
		Name:      node.Name,
		ChainID:   node.ChainID,
		URL:       node.URL,
		CreatedAt: node.CreatedAt,
		UpdatedAt: node.UpdatedAt,
	}
}
