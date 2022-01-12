package presenters

import (
	"time"

	terratypes "github.com/smartcontractkit/chainlink/core/chains/terra/types"
)

// TerraNodeResource is a Terra node JSONAPI resource.
type TerraNodeResource struct {
	JAID
	Name          string    `json:"name"`
	TerraChainID  string    `json:"terraChainID"`
	TendermintURL string    `json:"tendermintURL"`
	FCDURL        string    `json:"fcdURL"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// GetName implements the api2go EntityNamer interface
func (r TerraNodeResource) GetName() string {
	return "terra_node"
}

// NewTerraNodeResource returns a new TerraNodeResource for node.
func NewTerraNodeResource(node terratypes.Node) TerraNodeResource {
	return TerraNodeResource{
		JAID:          NewJAIDInt32(node.ID),
		Name:          node.Name,
		TerraChainID:  node.TerraChainID,
		TendermintURL: node.TendermintURL,
		FCDURL:        node.FCDURL,
		CreatedAt:     node.CreatedAt,
		UpdatedAt:     node.UpdatedAt,
	}
}
