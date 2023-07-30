package cosmos

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/multierr"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/adapters"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"

	pkgcosmos "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos"
	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos/types"
	"github.com/smartcontractkit/chainlink/v2/core/services"
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
	Config           pg.QConfig //coreconfig.AppConfig
	Logger           logger.Logger
	DB               *sqlx.DB
	KeyStore         keystore.Cosmos
	EventBroadcaster pg.EventBroadcaster
	Configs          types.Configs
}

func (o *ChainSetOpts) Validate() (err error) {
	required := func(s string) error {
		return fmt.Errorf("%s is required", s)
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
		return nil, fmt.Errorf("cannot create new chain with ID %s, the chain is disabled", *cfg.ChainID)
	}
	c, err := newChain(*cfg.ChainID, cfg, o.DB, o.KeyStore, o.Config, o.EventBroadcaster, o.Configs, o.Logger)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// LegacyChainContainer is container interface for Cosmos chains
type LegacyChainContainer interface {
	Get(id string) (adapters.Chain, error)
	Len() int
	List(ids ...string) ([]adapters.Chain, error)
	Put(id string, chain adapters.Chain)
	Slice() []adapters.Chain
}

type LegacyChains = chains.ChainsKV[adapters.Chain]

var _ LegacyChainContainer = &LegacyChains{}

func NewLegacyChains() *LegacyChains {
	return chains.NewChainsKV[adapters.Chain]()
}

type LoopRelayerChainer interface {
	loop.Relayer
	Chain() adapters.Chain
}

type LoopRelayerSingleChain struct {
	loop.Relayer
	singleChain *SingleChainSet
}

func NewLoopRelayerSingleChain(r *pkgcosmos.Relayer, s *SingleChainSet) *LoopRelayerSingleChain {
	ra := relay.NewRelayerAdapter(r, s)
	return &LoopRelayerSingleChain{
		Relayer:     ra,
		singleChain: s,
	}
}
func (r *LoopRelayerSingleChain) Chain() adapters.Chain {
	return r.singleChain.chain
}

// implement service interface
func (r *LoopRelayerSingleChain) Start(ctx context.Context) error {
	var ms services.MultiStart
	return ms.Start(ctx, r.singleChain, r.Relayer)
}

func (r *LoopRelayerSingleChain) Close() error {
	return services.CloseAll(r.Relayer, r.singleChain)
}

func (r *LoopRelayerSingleChain) Name() string {
	return fmt.Sprintf("%s-%s", r.Relayer.Name(), r.singleChain.Name())
}

func (r *LoopRelayerSingleChain) Ready() (err error) {
	return errors.Join(r.Relayer.Ready(), r.singleChain.Ready())
}

func (r *LoopRelayerSingleChain) HealthReport() map[string]error {
	hr := make(map[string]error)
	maps.Copy(r.Relayer.HealthReport(), hr)
	maps.Copy(r.singleChain.HealthReport(), hr)
	return hr
}

var _ LoopRelayerChainer = &LoopRelayerSingleChain{}

func newChainSet(opts ChainSetOpts, cfgs CosmosConfigs) (adapters.ChainSet, map[string]adapters.Chain, error) {
	cosmosChains := map[string]adapters.Chain{}
	var err error
	for _, chain := range cfgs {
		if !chain.IsEnabled() {
			continue
		}
		var err2 error
		cosmosChains[*chain.ChainID], err2 = opts.NewTOMLChain(chain)
		if err2 != nil {
			err = multierr.Combine(err, err2)
			continue
		}
	}
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load some Cosmos chains: %w", err)
	}

	cs, err := chains.NewChainSet[db.Node, adapters.Chain](cosmosChains, &opts)
	if err != nil {
		return nil, nil, err
	}

	return cs, cosmosChains, nil
}

// SingleChainSet is a chainset with 1 chain. TODO remove when relayer interface is updated
type SingleChainSet struct {
	adapters.ChainSet
	ID    string
	chain adapters.Chain
}

/*
func (s *SingleChainSet) getChain(ctx context.Context) adapters.Chain {
	c, err := s.Chain(ctx, s.ID)
	if err != nil {
		panic(fmt.Errorf("inconsistent single chain set: %s", err))
	}
	return c
}
*/

func (s *SingleChainSet) Chain(ctx context.Context, id string) (adapters.Chain, error) {
	return s.chain, nil
}

func NewSingleChainSet(opts ChainSetOpts, cfg *CosmosConfig) (*SingleChainSet, error) {
	cs, m, err := newChainSet(opts, CosmosConfigs{cfg})
	if err != nil {
		return nil, err
	}
	if len(m) != 1 {
		return nil, fmt.Errorf("invalid Single chain: more than one chain %d", len(m))
	}
	var chain adapters.Chain
	for _, ch := range m {
		chain = ch
	}
	return &SingleChainSet{
		ChainSet: cs,
		ID:       *cfg.ChainID,
		chain:    chain,
	}, nil
}
