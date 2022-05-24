package presenters

import (
	"time"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/chains/solana"
)

// SolanaChainResource is an Solana chain JSONAPI resource.
type SolanaChainResource struct {
	chainResource[*db.ChainCfg]
}

// GetName implements the api2go EntityNamer interface
func (r SolanaChainResource) GetName() string {
	return "solana_chain"
}

// NewSolanaChainResource returns a new SolanaChainResource for chain.
func NewSolanaChainResource(chain solana.DBChain) SolanaChainResource {
	return SolanaChainResource{chainResource[*db.ChainCfg]{
		JAID:      NewJAID(chain.ID),
		Config:    chain.Cfg,
		Enabled:   chain.Enabled,
		CreatedAt: chain.CreatedAt,
		UpdatedAt: chain.UpdatedAt,
	}}
}

// SolanaNodeResource is a Solana node JSONAPI resource.
type SolanaNodeResource struct {
	JAID
	Name          string    `json:"name"`
	SolanaChainID string    `json:"solanaChainID"`
	SolanaURL     string    `json:"solanaURL"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// GetName implements the api2go EntityNamer interface
func (r SolanaNodeResource) GetName() string {
	return "solana_node"
}

// NewSolanaNodeResource returns a new SolanaNodeResource for node.
func NewSolanaNodeResource(node db.Node) SolanaNodeResource {
	return SolanaNodeResource{
		JAID:          NewJAIDInt32(node.ID),
		Name:          node.Name,
		SolanaChainID: node.SolanaChainID,
		SolanaURL:     node.SolanaURL,
		CreatedAt:     node.CreatedAt,
		UpdatedAt:     node.UpdatedAt,
	}
}
