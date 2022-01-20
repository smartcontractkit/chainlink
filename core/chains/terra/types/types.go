package types

import (
	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/services/pg"
)

// ORM manages terra chains and nodes.
type ORM interface {
	Chain(string, ...pg.QOpt) (db.Chain, error)
	Chains(offset, limit int, qopts ...pg.QOpt) ([]db.Chain, int, error)
	CreateChain(id string, config db.ChainCfg, qopts ...pg.QOpt) (db.Chain, error)
	UpdateChain(id string, enabled bool, config db.ChainCfg, qopts ...pg.QOpt) (db.Chain, error)
	DeleteChain(id string, qopts ...pg.QOpt) error

	// EnabledChainsWithNodes returns enabled chains with nodes (if any) included.
	EnabledChainsWithNodes(...pg.QOpt) ([]db.Chain, error)

	CreateNode(NewNode, ...pg.QOpt) (db.Node, error)
	DeleteNode(int32, ...pg.QOpt) error
	Node(int32, ...pg.QOpt) (db.Node, error)
	Nodes(offset, limit int, qopts ...pg.QOpt) (nodes []db.Node, count int, err error)
	NodesForChain(chainID string, offset, limit int, qopts ...pg.QOpt) (nodes []db.Node, count int, err error)
}

// NewNode defines a new node to create.
type NewNode struct {
	Name          string `json:"name"`
	TerraChainID  string `json:"terraChainId"`
	TendermintURL string `json:"tendermintURL" db:"tendermint_url"`
	FCDURL        string `json:"fcdURL" db:"fcd_url"`
}
