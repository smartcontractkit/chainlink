package solana

import (
	"github.com/smartcontractkit/sqlx"

	soldb "github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

type DBChain = chains.DBChain[string, *soldb.ChainCfg]

// ORM manages solana chains and nodes.
type ORM interface {
	Chain(string, ...pg.QOpt) (DBChain, error)
	Chains(offset, limit int, qopts ...pg.QOpt) ([]DBChain, int, error)
	CreateChain(id string, config *soldb.ChainCfg, qopts ...pg.QOpt) (DBChain, error)
	UpdateChain(id string, enabled bool, config *soldb.ChainCfg, qopts ...pg.QOpt) (DBChain, error)
	DeleteChain(id string, qopts ...pg.QOpt) error
	GetChainsByIDs(ids []string) (chains []DBChain, err error)
	EnabledChains(...pg.QOpt) ([]DBChain, error)

	CreateNode(soldb.Node, ...pg.QOpt) (soldb.Node, error)
	DeleteNode(int32, ...pg.QOpt) error
	GetNodesByChainIDs(chainIDs []string, qopts ...pg.QOpt) (nodes []soldb.Node, err error)
	Node(int32, ...pg.QOpt) (soldb.Node, error)
	NodeNamed(string, ...pg.QOpt) (soldb.Node, error)
	Nodes(offset, limit int, qopts ...pg.QOpt) (nodes []soldb.Node, count int, err error)
	NodesForChain(chainID string, offset, limit int, qopts ...pg.QOpt) (nodes []soldb.Node, count int, err error)

	SetupNodes([]soldb.Node, []string) error

	StoreString(chainID string, key, val string) error
	Clear(chainID string, key string) error
}

var _ chains.ORM[string, *soldb.ChainCfg, soldb.Node] = (ORM)(nil)

// NewORM returns an ORM backed by db.
func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig) ORM {
	q := pg.NewQ(db, lggr.Named("ORM"), cfg)
	return chains.NewORM[string, *soldb.ChainCfg, soldb.Node](q, "solana", "solana_url")
}
