package types

import (
	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

// ORM manages terra chains and nodes.
type ORM interface {
	Chain(string, ...pg.QOpt) (Chain, error)
	Chains(offset, limit int, qopts ...pg.QOpt) ([]Chain, int, error)
	CreateChain(id string, config db.ChainCfg, qopts ...pg.QOpt) (Chain, error)
	UpdateChain(id string, enabled bool, config db.ChainCfg, qopts ...pg.QOpt) (Chain, error)
	DeleteChain(id string, qopts ...pg.QOpt) error
	EnabledChains(...pg.QOpt) ([]Chain, error)

	CreateNode(db.Node, ...pg.QOpt) (db.Node, error)
	DeleteNode(int32, ...pg.QOpt) error
	Node(int32, ...pg.QOpt) (db.Node, error)
	NodeNamed(string, ...pg.QOpt) (db.Node, error)
	Nodes(offset, limit int, qopts ...pg.QOpt) (nodes []db.Node, count int, err error)
	NodesForChain(chainID string, offset, limit int, qopts ...pg.QOpt) (nodes []db.Node, count int, err error)
}

type Chain = chains.Chain[string, db.ChainCfg]

// NewNode defines a new node to create.
type NewNode struct {
	Name          string `json:"name"`
	TerraChainID  string `json:"terraChainId"`
	TendermintURL string `json:"tendermintURL" db:"tendermint_url"`
}
