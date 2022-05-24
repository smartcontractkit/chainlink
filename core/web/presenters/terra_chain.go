package presenters

import (
	"time"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
)

// TerraChainResource is an Terra chain JSONAPI resource.
type TerraChainResource struct {
	chainResource[*db.ChainCfg]
}

// GetName implements the api2go EntityNamer interface
func (r TerraChainResource) GetName() string {
	return "terra_chain"
}

// NewTerraChainResource returns a new TerraChainResource for chain.
func NewTerraChainResource(chain types.DBChain) TerraChainResource {
	return TerraChainResource{chainResource[*db.ChainCfg]{
		JAID:      NewJAID(chain.ID),
		Config:    chain.Cfg,
		Enabled:   chain.Enabled,
		CreatedAt: chain.CreatedAt,
		UpdatedAt: chain.UpdatedAt,
	}}
}

// TerraNodeResource is a Terra node JSONAPI resource.
type TerraNodeResource struct {
	JAID
	Name          string    `json:"name"`
	TerraChainID  string    `json:"terraChainID"`
	TendermintURL string    `json:"tendermintURL"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// GetName implements the api2go EntityNamer interface
func (r TerraNodeResource) GetName() string {
	return "terra_node"
}

// NewTerraNodeResource returns a new TerraNodeResource for node.
func NewTerraNodeResource(node db.Node) TerraNodeResource {
	return TerraNodeResource{
		JAID:          NewJAIDInt32(node.ID),
		Name:          node.Name,
		TerraChainID:  node.TerraChainID,
		TendermintURL: node.TendermintURL,
		CreatedAt:     node.CreatedAt,
		UpdatedAt:     node.UpdatedAt,
	}
}
