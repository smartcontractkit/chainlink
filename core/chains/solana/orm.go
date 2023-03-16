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
	chains.ChainsORM[string, *soldb.ChainCfg, DBChain]
	chains.NodesORM[string, soldb.Node]

	EnsureChains([]string, ...pg.QOpt) error
}

var _ chains.ORM[string, *soldb.ChainCfg, soldb.Node] = (ORM)(nil)

// NewORM returns an ORM backed by db.
// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) ORM {
	q := pg.NewQ(db, lggr.Named("ORM"), cfg)
	return chains.NewORM[string, *soldb.ChainCfg, soldb.Node](q, "solana", "solana_url")
}

func NewORMImmut(cfgs chains.ChainConfig[string, *soldb.ChainCfg, soldb.Node]) ORM {
	return chains.NewORMImmut(cfgs)
}
