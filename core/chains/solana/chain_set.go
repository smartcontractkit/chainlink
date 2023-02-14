package solana

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/logger"
)

// ChainSetOpts holds options for configuring a ChainSet.
type ChainSetOpts struct {
	Logger   logger.Logger
	KeyStore loop.Keystore
	ORM      ORM
}

func (o *ChainSetOpts) Validate() (err error) {
	required := func(s string) error {
		return errors.Errorf("%s is required", s)
	}
	if o.Logger == nil {
		err = multierr.Append(err, required("Logger"))
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

func (o *ChainSetOpts) NewTOMLChain(cfg *SolanaConfig) (solana.Chain, error) {
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

//go:generate mockery --quiet --name ChainSet --srcpkg github.com/smartcontractkit/chainlink-solana/pkg/solana --output ./mocks/ --case=underscore

// ChainSet extends solana.ChainSet with mutability.
type ChainSet interface {
	solana.ChainSet

	Show(id string) (DBChain, error)
	Index(offset, limit int) ([]DBChain, int, error)
	GetNodes(ctx context.Context, offset, limit int) (nodes []db.Node, count int, err error)
	GetNodesForChain(ctx context.Context, chainID string, offset, limit int) (nodes []db.Node, count int, err error)
}

func NewChainSetImmut(opts ChainSetOpts, cfgs SolanaConfigs) (ChainSet, error) {
	solChains := map[string]solana.Chain{}
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
	return chains.NewChainSetImmut[string, *db.ChainCfg, db.Node, solana.Chain](solChains, &opts, func(s string) string { return s })
}
