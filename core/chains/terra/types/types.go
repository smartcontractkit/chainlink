package types

import (
	"time"
)

// ORM manages terra chains and nodes.
type ORM interface {
	CreateNode(NewNode) (Node, error)
	DeleteNode(int32) error
	Node(int32) (Node, error)
	Nodes(offset, limit int) (nodes []Node, count int, err error)
	NodesForChain(chainID string, offset, limit int) (nodes []Node, count int, err error)
}

// ChainCfg is configuration parameters for a terra chain.
type ChainCfg struct {
	FallbackGasPriceULuna string
	GasLimitMultiplier    string
}

// NewNode defines a new node to create.
type NewNode struct {
	Name          string `json:"name"`
	TerraChainID  string `json:"terraChainId"`
	TendermintURL string `json:"tendermintURL" db:"tendermint_url"`
	FCDURL        string `json:"fcdURL" db:"fcd_url"`
}

// Node is an existing node.
type Node struct {
	ID            int32
	Name          string
	TerraChainID  string
	TendermintURL string `db:"tendermint_url"`
	FCDURL        string `db:"fcd_url"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
