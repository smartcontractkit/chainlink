package starknet

import (
	"github.com/smartcontractkit/sqlx"

	starknetdb "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/chains/starknet/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

func EnsureChains(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig, ids []string) error {
	q := pg.NewQ(db, lggr.Named("Ensure"), cfg)
	return chains.EnsureChains[string](q, "starknet", ids)
}

func NewConfigs(cfgs chains.ConfigsV2[string, starknetdb.Node]) types.Configs {
	return chains.NewConfigs(cfgs)
}
