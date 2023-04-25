package solana

import (
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// ChainSetOpts holds options for configuring a ChainSet.
type ChainSetOpts struct {
	Logger   logger.Logger
	KeyStore loop.Keystore
	Configs  Configs
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
	if o.Configs == nil {
		err = multierr.Append(err, required("Configs"))
	}
	return
}

func (o *ChainSetOpts) ConfigsAndLogger() (chains.Configs[string, db.Node], logger.Logger) {
	return o.Configs, o.Logger
}

func (o *ChainSetOpts) NewTOMLChain(cfg *SolanaConfig) (solana.Chain, error) {
	if !cfg.IsEnabled() {
		return nil, errors.Errorf("cannot create new chain with ID %s, the chain is disabled", *cfg.ChainID)
	}
	c, err := newChain(*cfg.ChainID, cfg, o.KeyStore, o.Configs, o.Logger)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func NewChainSet(opts ChainSetOpts, cfgs SolanaConfigs) (solana.ChainSet, error) {
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
	return chains.NewChainSet[db.Node, solana.Chain](solChains, &opts)
}
