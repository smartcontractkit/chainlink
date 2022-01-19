package types

import (
	"github.com/smartcontractkit/chainlink/core/services/pg"

	. "github.com/smartcontractkit/chainlink-terra/pkg/terra/db"
)

// ORM manages terra chains and nodes.
type ORM interface {
	Chain(string, ...pg.QOpt) (Chain, error)
	Chains(offset, limit int, qopts ...pg.QOpt) ([]Chain, int, error)
	CreateChain(id string, config ChainCfg, qopts ...pg.QOpt) (Chain, error)
	UpdateChain(id string, enabled bool, config ChainCfg, qopts ...pg.QOpt) (Chain, error)
	DeleteChain(id string, qopts ...pg.QOpt) error

	// EnabledChainsWithNodes returns enabled chains with nodes (if any) included.
	EnabledChainsWithNodes(...pg.QOpt) ([]Chain, error)

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
