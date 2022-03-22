package types

import (
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/services/pg"
)

// ORM manages solana chains and nodes.
type ORM interface {
	Chain(string, ...pg.QOpt) (db.Chain, error)
	Chains(offset, limit int, qopts ...pg.QOpt) ([]db.Chain, int, error)
	CreateChain(id string, config db.ChainCfg, qopts ...pg.QOpt) (db.Chain, error)
	UpdateChain(id string, enabled bool, config db.ChainCfg, qopts ...pg.QOpt) (db.Chain, error)
	DeleteChain(id string, qopts ...pg.QOpt) error
	EnabledChains(...pg.QOpt) ([]db.Chain, error)

	CreateNode(db.NewNode, ...pg.QOpt) (db.Node, error)
	DeleteNode(int32, ...pg.QOpt) error
	Node(int32, ...pg.QOpt) (db.Node, error)
	NodeNamed(string, ...pg.QOpt) (db.Node, error)
	Nodes(offset, limit int, qopts ...pg.QOpt) (nodes []db.Node, count int, err error)
	NodesForChain(chainID string, offset, limit int, qopts ...pg.QOpt) (nodes []db.Node, count int, err error)
}
