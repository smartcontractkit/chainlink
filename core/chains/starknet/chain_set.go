package starknet

import (
	"context"
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink-starknet/pkg/starknet"
	"github.com/smartcontractkit/chainlink-starknet/pkg/starknet/db"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/chains/starknet/types"
	coreconfig "github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
)

type ChainSetOpts struct {
	Config coreconfig.GeneralConfig
	Logger logger.Logger
	DB     *sqlx.DB
	ORM    types.ORM
}

func (o *ChainSetOpts) Validate() (err error) {
	required := func(s string) error {
		return errors.Errorf("%s is required", s)
	}
	if o.Config == nil {
		err = multierr.Append(err, required("Config"))
	}
	if o.Logger == nil {
		err = multierr.Append(err, required("Logger'"))
	}
	if o.DB == nil {
		err = multierr.Append(err, required("DB"))
	}
	if o.ORM == nil {
		err = multierr.Append(err, required("ORM"))
	}
	return
}

func (o *ChainSetOpts) ORMAndLogger() (chains.ORM[string, *db.ChainCfg, db.Node], logger.Logger) {
	return o.ORM, o.Logger
}

func (o *ChainSetOpts) NewChain(dbchain types.DBChain) (starknet.Chain, error) {
	if !dbchain.Enabled {
		return nil, errors.Errorf("cannot create new chain with ID %s, the chain is disabled", dbchain.ID)
	}
	return NewChain(o.DB, dbchain, o.ORM, o.Logger)
}

type ChainSet interface {
	starknet.ChainSet

	Add(context.Context, string, *db.ChainCfg) (types.DBChain, error)
	Remove(string) error
	Configure(ctx context.Context, id string, enabled bool, config *db.ChainCfg) (types.DBChain, error)
	Show(id string) (types.DBChain, error)
	Index(offset, limit int) ([]types.DBChain, int, error)
	GetNodes(ctx context.Context, offset, limit int) (nodes []db.Node, count int, err error)
	GetNodesForChain(ctx context.Context, chainID string, offset, limit int) (nodes []db.Node, count int, err error)
	CreateNode(ctx context.Context, data db.Node) (db.Node, error)
	DeleteNode(ctx context.Context, id int32) error
}

// NewChainSet returns a new chain set for opts.
func NewChainSet(opts ChainSetOpts) (ChainSet, error) {
	return chains.NewChainSet[string, *db.ChainCfg, db.Node, starknet.Chain](&opts, func(s string) string { return s })
}
