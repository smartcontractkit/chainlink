package terra

import (
	"github.com/smartcontractkit/sqlx"

	terradb "github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

// NewORM returns an ORM backed by db.
// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig) types.ORM {
	q := pg.NewQ(db, lggr.Named("ORM"), cfg)
	return chains.NewORM[string, *terradb.ChainCfg, terradb.Node](q, "terra", "tendermint_url")
}

func NewORMImmut(cfgs chains.ChainConfig[string, *terradb.ChainCfg, terradb.Node]) types.ORM {
	return chains.NewORMImmut(cfgs)
}
