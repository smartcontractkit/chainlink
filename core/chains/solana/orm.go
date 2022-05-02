package solana

import (
	"github.com/smartcontractkit/sqlx"

	soldb "github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type Chain = chains.Chain[string, soldb.ChainCfg]

// ORM manages solana chains and nodes.
type ORM interface {
	Chain(string, ...pg.QOpt) (Chain, error)
	Chains(offset, limit int, qopts ...pg.QOpt) ([]Chain, int, error)
	CreateChain(id string, config soldb.ChainCfg, qopts ...pg.QOpt) (Chain, error)
	UpdateChain(id string, enabled bool, config soldb.ChainCfg, qopts ...pg.QOpt) (Chain, error)
	DeleteChain(id string, qopts ...pg.QOpt) error
	EnabledChains(...pg.QOpt) ([]Chain, error)

	CreateNode(soldb.Node, ...pg.QOpt) (soldb.Node, error)
	DeleteNode(int32, ...pg.QOpt) error
	Node(int32, ...pg.QOpt) (soldb.Node, error)
	NodeNamed(string, ...pg.QOpt) (soldb.Node, error)
	Nodes(offset, limit int, qopts ...pg.QOpt) (nodes []soldb.Node, count int, err error)
	NodesForChain(chainID string, offset, limit int, qopts ...pg.QOpt) (nodes []soldb.Node, count int, err error)

	SetupNodes([]soldb.Node, []string) error
}

// NewORM returns an ORM backed by db.
func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig) ORM {
	q := pg.NewQ(db, lggr.Named("ORM"), cfg)
	return chains.NewORM[string, soldb.ChainCfg, soldb.Node](q, "solana", "solana_url")
}
