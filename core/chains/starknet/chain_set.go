package starknet

import (
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	starkchain "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/chain"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/chains/starknet/types"
)

// TODO rename to ChainOpts
type ChainSetOpts struct {
	Logger logger.Logger
	// the implementation used here needs to be co-ordinated with the starknet transaction manager keystore adapter
	KeyStore loop.Keystore
	Configs  types.Configs
}

func (o *ChainSetOpts) Name() string {
	return o.Logger.Name()
}

func (o *ChainSetOpts) Validate() (err error) {
	required := func(s string) error {
		return errors.Errorf("%s is required", s)
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

func NewChain(cfg *StarknetConfig, opts ChainSetOpts) (starkchain.Chain, error) {
	if !cfg.IsEnabled() {
		return nil, fmt.Errorf("cannot create new chain with ID %s: %w", *cfg.ChainID, chains.ErrChainDisabled)
	}
	c, err := newChain(*cfg.ChainID, cfg, opts.KeyStore, opts.Configs, opts.Logger)
	if err != nil {
		return nil, err
	}
	return c, nil
}
