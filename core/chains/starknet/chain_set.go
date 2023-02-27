package starknet

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	starkchain "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/chain"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/chains/starknet/types"
	coreconfig "github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
)

type ChainSetOpts struct {
	Config   coreconfig.BasicConfig
	Logger   logger.Logger
	KeyStore keystore.StarkNet
	ORM      types.ORM
}

func (o *ChainSetOpts) Name() string {
	return o.Logger.Name()
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
	if o.KeyStore == nil {
		err = multierr.Append(err, required("KeyStore"))
	}
	if o.ORM == nil {
		err = multierr.Append(err, required("ORM"))
	}
	return
}

func (o *ChainSetOpts) ORMAndLogger() (chains.ORM[string, *db.ChainCfg, db.Node], logger.Logger) {
	return o.ORM, o.Logger
}

func (o *ChainSetOpts) NewChain(dbchain types.DBChain) (starkchain.Chain, error) {
	if !dbchain.Enabled {
		return nil, errors.Errorf("cannot create new chain with ID %s, the chain is disabled", dbchain.ID)
	}
	cfg := config.NewConfig(*dbchain.Cfg, o.Logger)
	return newChain(dbchain.ID, cfg, o.KeyStore, o.ORM, o.Logger)
}

func (o *ChainSetOpts) NewTOMLChain(cfg *StarknetConfig) (starkchain.Chain, error) {
	if !cfg.IsEnabled() {
		return nil, errors.Errorf("cannot create new chain with ID %s, the chain is disabled", *cfg.ChainID)
	}
	c, err := newChain(*cfg.ChainID, cfg, o.KeyStore, o.ORM, o.Logger)
	if err != nil {
		return nil, err
	}
	c.cfgImmutable = true
	return c, nil
}

type ChainSet interface {
	starkchain.ChainSet

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
// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func NewChainSet(opts ChainSetOpts) (ChainSet, error) {
	return chains.NewChainSet[string, *db.ChainCfg, db.Node, starkchain.Chain](&opts, func(s string) string { return s })
}

func NewChainSetImmut(opts ChainSetOpts, cfgs StarknetConfigs) (ChainSet, error) {
	stkChains := map[string]starkchain.Chain{}
	var err error
	for _, chain := range cfgs {
		if !chain.IsEnabled() {
			continue
		}
		var err2 error
		stkChains[*chain.ChainID], err2 = opts.NewTOMLChain(chain)
		if err2 != nil {
			err = multierr.Combine(err, err2)
			continue
		}
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to load some Solana chains")
	}
	return chains.NewChainSetImmut[string, *db.ChainCfg, db.Node, starkchain.Chain](stkChains, &opts, func(s string) string { return s })
}
