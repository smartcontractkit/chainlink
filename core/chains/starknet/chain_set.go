package starknet

import (
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	starkchain "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/chain"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/chains/starknet/types"
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
)

type ChainSetOpts struct {
	Config   coreconfig.BasicConfig
	Logger   logger.Logger
	KeyStore keystore.StarkNet
	Configs  types.Configs
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
	if o.Configs == nil {
		err = multierr.Append(err, required("Configs"))
	}
	return
}

func (o *ChainSetOpts) ConfigsAndLogger() (chains.Configs[string, db.Node], logger.Logger) {
	return o.Configs, o.Logger
}

func (o *ChainSetOpts) NewTOMLChain(cfg *StarknetConfig) (starkchain.Chain, error) {
	if !cfg.IsEnabled() {
		return nil, errors.Errorf("cannot create new chain with ID %s, the chain is disabled", *cfg.ChainID)
	}
	c, err := newChain(*cfg.ChainID, cfg, o.KeyStore, o.Configs, o.Logger)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func NewChainSet(opts ChainSetOpts, cfgs StarknetConfigs) (starkchain.ChainSet, error) {
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
	return chains.NewChainSet[db.Node, starkchain.Chain](stkChains, &opts)
}
