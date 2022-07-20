package evm

import (
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// NewORM returns a new EVM ORM
func NewORM(db *sqlx.DB, lggr logger.Logger, cfg pg.LogConfig) types.ORM {
	q := pg.NewQ(db, lggr.Named("EVMORM"), cfg)
	return chains.NewORM[utils.Big, *types.ChainCfg, types.Node](q, "evm", "ws_url", "http_url", "send_only")
}
