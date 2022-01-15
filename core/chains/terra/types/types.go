package types

import (
	"time"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra/config"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

// ORM manages terra chains and nodes.
type ORM interface {
	EnabledChainsWithNodes(...pg.QOpt) ([]Chain, error)
	Chain(string, ...pg.QOpt) (Chain, error)
	CreateNode(NewNode, ...pg.QOpt) (Node, error)
	DeleteNode(int32, ...pg.QOpt) error
	Node(int32, ...pg.QOpt) (Node, error)
	Nodes(offset, limit int, qopts ...pg.QOpt) (nodes []Node, count int, err error)
	NodesForChain(chainID string, offset, limit int, qopts ...pg.QOpt) (nodes []Node, count int, err error)
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

// Chain is a an existing chain.
type Chain struct {
	ID    string
	Nodes []Node
	Cfg   config.ChainCfg
}
