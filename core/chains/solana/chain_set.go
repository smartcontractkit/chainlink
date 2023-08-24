package solana

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana/db"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
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

func NewChain(cfg *SolanaConfig, opts ChainSetOpts) (solana.Chain, error) {
	if !cfg.IsEnabled() {
		return nil, fmt.Errorf("cannot create new chain with ID %s: %w", *cfg.ChainID, chains.ErrChainDisabled)
	}
	c, err := newChain(*cfg.ChainID, cfg, opts.KeyStore, opts.Configs, opts.Logger)
	if err != nil {
		return nil, err
	}
	return c, nil
}
