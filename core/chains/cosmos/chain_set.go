package cosmos

import (
	"errors"
	"fmt"

	"go.uber.org/multierr"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/adapters"
	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/db"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"

	pkgcosmos "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos"
	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos/types"
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

// Chain is a wrap for easy use in other places in the core node
type Chain = adapters.Chain

// ChainSetOpts holds options for configuring a ChainSet.
type ChainSetOpts struct {
	QueryConfig      pg.QConfig
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
	if o.QueryConfig == nil {
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
	c, err := newChain(*cfg.ChainID, cfg, o.DB, o.KeyStore, o.QueryConfig, o.EventBroadcaster, o.Configs, o.Logger)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func NewChain(cfg *CosmosConfig, opts ChainSetOpts) (adapters.Chain, error) {
	if !cfg.IsEnabled() {
		return nil, fmt.Errorf("cannot create new chain with ID %s, the chain is disabled", *cfg.ChainID)
	}
	c, err := newChain(*cfg.ChainID, cfg, opts.DB, opts.KeyStore, opts.QueryConfig, opts.EventBroadcaster, opts.Configs, opts.Logger)
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
	Slice() []adapters.Chain
}

type LegacyChains = chains.ChainsKV[adapters.Chain]

var _ LegacyChainContainer = &LegacyChains{}

func NewLegacyChains(m map[string]adapters.Chain) *LegacyChains {
	return chains.NewChainsKV[adapters.Chain](m)
}

type LoopRelayerChainer interface {
	loop.Relayer
	Chain() adapters.Chain
}

type LoopRelayerChain struct {
	loop.Relayer
	chain adapters.Chain
}

func NewLoopRelayerChain(r *pkgcosmos.Relayer, s adapters.Chain) *LoopRelayerChain {
	ra := relay.NewRelayerAdapter(r, s)
	return &LoopRelayerChain{
		Relayer: ra,
		chain:   s,
	}
}
func (r *LoopRelayerChain) Chain() adapters.Chain {
	return r.chain
}

var _ LoopRelayerChainer = &LoopRelayerChain{}
