package starknet

import (
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

func (o *ChainSetOpts) NewChain(cc types.ChainConfig) (starkchain.Chain, error) {
	if !cc.Enabled {
		return nil, errors.Errorf("cannot create new chain with ID %s, the chain is disabled", cc.ID)
	}
	cfg := config.NewConfig(*cc.Cfg, o.Logger)
	return newChain(cc.ID, cfg, o.KeyStore, o.ORM, o.Logger)
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
	chains.ChainsConfig[string, *db.ChainCfg]
	chains.NodesConfig[string, db.Node]
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
