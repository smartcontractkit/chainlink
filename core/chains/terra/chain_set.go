package terra

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra"
	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
	coreconfig "github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

var (
	// ErrChainIDEmpty is returned when chain is required but was empty.
	ErrChainIDEmpty = errors.New("chain id empty")
	// ErrChainIDInvalid is returned when a chain id does not match any configured chains.
	ErrChainIDInvalid = errors.New("chain id does not match any local chains")
)

// ChainSetOpts holds options for configuring a ChainSet.
type ChainSetOpts struct {
	Config           coreconfig.BasicConfig
	Logger           logger.Logger
	DB               *sqlx.DB
	KeyStore         keystore.Terra
	EventBroadcaster pg.EventBroadcaster
	ORM              types.ORM
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
	if o.KeyStore == nil {
		err = multierr.Append(err, required("KeyStore"))
	}
	if o.EventBroadcaster == nil {
		err = multierr.Append(err, required("EventBroadcaster"))
	}
	if o.ORM == nil {
		err = multierr.Append(err, required("ORM"))
	}
	return
}

func (o *ChainSetOpts) ORMAndLogger() (chains.ORM[string, *db.ChainCfg, db.Node], logger.Logger) {
	return o.ORM, o.Logger
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func (o *ChainSetOpts) NewChain(dbchain types.DBChain) (terra.Chain, error) {
	if !dbchain.Enabled {
		return nil, errors.Errorf("cannot create new chain with ID %s, the chain is disabled", dbchain.ID)
	}
	id := dbchain.ID
	cfg := terra.NewConfig(*dbchain.Cfg, o.Logger)
	return newChain(id, cfg, o.DB, o.KeyStore, o.Config, o.EventBroadcaster, o.ORM, o.Logger)
}

func (o *ChainSetOpts) NewTOMLChain(cfg *TerraConfig) (terra.Chain, error) {
	if !cfg.IsEnabled() {
		return nil, errors.Errorf("cannot create new chain with ID %s, the chain is disabled", *cfg.ChainID)
	}
	c, err := newChain(*cfg.ChainID, cfg, o.DB, o.KeyStore, o.Config, o.EventBroadcaster, o.ORM, o.Logger)
	if err != nil {
		return nil, err
	}
	c.cfgImmutable = true
	return c, nil
}

//go:generate mockery --quiet --name ChainSet --srcpkg github.com/smartcontractkit/chainlink-terra/pkg/terra --output ./mocks/ --case=underscore

// ChainSet extends terra.ChainSet with mutability and exposes the underlying ORM.
type ChainSet interface {
	terra.ChainSet

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
	return chains.NewChainSet[string, *db.ChainCfg, db.Node, terra.Chain](&opts, func(s string) string { return s })
}

func NewChainSetImmut(opts ChainSetOpts, cfgs TerraConfigs) (ChainSet, error) {
	solChains := map[string]terra.Chain{}
	var err error
	for _, chain := range cfgs {
		if !chain.IsEnabled() {
			continue
		}
		var err2 error
		solChains[*chain.ChainID], err2 = opts.NewTOMLChain(chain)
		if err2 != nil {
			err = multierr.Combine(err, err2)
			continue
		}
	}
	if err != nil {
		return nil, errors.Wrap(err, "failed to load some Solana chains")
	}
	return chains.NewChainSetImmut[string, *db.ChainCfg, db.Node, terra.Chain](solChains, &opts, func(s string) string { return s })
}
