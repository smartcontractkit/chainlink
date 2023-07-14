package cosmos

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/adapters"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos/types"
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

var (
	// ErrChainIDEmpty is returned when chain is required but was empty.
	ErrChainIDEmpty = errors.New("chain id empty")
	// ErrChainIDInvalid is returned when a chain id does not match any configured chains.
	ErrChainIDInvalid = errors.New("chain id does not match any local chains")
)

// ChainSetOpts holds options for configuring a ChainSet.
type ChainSetOpts struct {
	Config           coreconfig.AppConfig
	Logger           logger.Logger
	DB               *sqlx.DB
	KeyStore         keystore.Cosmos
	EventBroadcaster pg.EventBroadcaster
	Configs          types.Configs
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
	if o.Configs == nil {
		err = multierr.Append(err, required("Configs"))
	}
	return
}

func (o *ChainSetOpts) ConfigsAndLogger() (chains.Configs[string, db.Node], logger.Logger) {
	return o.Configs, o.Logger
}

func (o *ChainSetOpts) NewTOMLChain(cfg *CosmosConfig) (adapters.Chain, error) {
	if !cfg.IsEnabled() {
		return nil, errors.Errorf("cannot create new chain with ID %s, the chain is disabled", *cfg.ChainID)
	}
	c, err := newChain(*cfg.ChainID, cfg, o.DB, o.KeyStore, o.Config.Database(), o.EventBroadcaster, o.Configs, o.Logger)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// SingleChainSet is a chainset with 1 chain. TODO remove when relayer interface is updated
type SingleChainSet struct {
	adapters.ChainSet
	// TODO what type for ID?
	ID string
}

func (s SingleChainSet) GetChain(ctx context.Context) adapters.Chain {
	c, err := s.Chain(ctx, s.ID)
	if err != nil {
		panic("inconsistent single chain set")
	}
	return c
}

type ChainSetX struct {
	adapters.ChainSet
	chains map[relay.Identifier]SingleChainSet
}

func (cs *ChainSetX) Chains() map[relay.Identifier]SingleChainSet {
	return cs.chains
}

func NewChainSet(opts ChainSetOpts, cfgs CosmosConfigs) (*ChainSetX, error) {
	solChains := map[string]adapters.Chain{}
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
		return nil, errors.Wrap(err, "failed to load some Cosmos chains")
	}
	chainMap := make(map[relay.Identifier]SingleChainSet)
	for id, chain := range solChains {
		temp := map[string]adapters.Chain{id: chain}
		cs, err := chains.NewChainSet[db.Node, adapters.Chain](temp, &opts)
		if err != nil {
			return nil, err
		}

		relayID := relay.Identifier{Network: relay.Cosmos, ChainID: relay.ChainID(id)}
		chainMap[relayID] = SingleChainSet{
			ChainSet: cs,
			ID:       id,
		}

	}
	cs, err := chains.NewChainSet[db.Node, adapters.Chain](solChains, &opts)
	if err != nil {
		return nil, err
	}
	return &ChainSetX{
		ChainSet: cs,
		chains:   chainMap,
	}, nil
}
