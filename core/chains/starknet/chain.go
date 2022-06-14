package starknet

import (
	"context"
	"github.com/smartcontractkit/chainlink-starknet/pkg/relay/starknet"
	"github.com/smartcontractkit/chainlink-starknet/pkg/relay/starknet/db"

	"github.com/smartcontractkit/chainlink/core/chains/starknet/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/sqlx"
)

var _ starknet.Chain = (*chain)(nil)

type chain struct {
	utils.StartStopOnce
	id   string
	cfg  starknet.Config
	orm  types.ORM
	lggr logger.Logger
}

func NewChain(db *sqlx.DB, dbchain types.DBChain, orm types.ORM, lggr logger.Logger) (*chain, error) {
	cfg := starknet.NewConfig(*dbchain.Cfg, lggr)
	lggr = lggr.With("starknetChainID", dbchain.ID)
	var ch = chain{
		id:   dbchain.ID,
		cfg:  cfg,
		orm:  orm,
		lggr: lggr.Named("Chain"),
	}

	return &ch, nil
}

func (c *chain) Config() starknet.Config {
	return c.cfg
}

func (c *chain) UpdateConfig(cfg *db.ChainCfg) {
	c.cfg.Update(*cfg)
}

func (c *chain) Start(ctx context.Context) error {
	return c.StartOnce("Chain", func() error {
		return nil
	})
}

func (c *chain) Close() error {
	return c.StopOnce("Chain", func() error {
		return nil
	})
}

func (c *chain) Ready() error {
	return c.StartStopOnce.Ready()
}

func (c *chain) Healthy() error {
	return c.StartStopOnce.Healthy()
}
