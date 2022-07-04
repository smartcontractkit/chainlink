package starknet

import (
	"github.com/smartcontractkit/sqlx"

	starknetdb "github.com/smartcontractkit/chainlink-starknet/pkg/starknet/db"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/chains/starknet/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig) types.ORM {
	q := pg.NewQ(db, lggr.Named("ORM"), cfg)
	return chains.NewORM[string, *starknetdb.ChainCfg, starknetdb.Node](q, "starknet", "starknet_url")
}
