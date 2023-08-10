package solana

import (
	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
)

// ChainOpts holds options for configuring a Chain.
type ChainOpts struct {
	Logger   logger.Logger
	KeyStore loop.Keystore
	// ChainNodeStatuser ConfigStater
}

func (o *ChainOpts) Validate() (err error) {
	required := func(s string) error {
		return errors.Errorf("%s is required", s)
	}
	if o.Logger == nil {
		err = multierr.Append(err, required("Logger"))
	}
	if o.KeyStore == nil {
		err = multierr.Append(err, required("KeyStore"))
	}
	return
}

/*
	func (o *ChainOpts) ConfigsAndLogger() (chains.Statuser[db.Node], logger.Logger) {
		return o.ChainNodeStatuser, o.Logger
	}
*/
func NewChain(cfg *SolanaConfig, o ChainOpts) (solana.Chain, error) {
	if !cfg.IsEnabled() {
		return nil, errors.Errorf("cannot create new chain with ID %s, the chain is disabled", *cfg.ChainID)
	}
	c, err := newChain(*cfg.ChainID, cfg, o)
	if err != nil {
		return nil, err
	}
	return c, nil
}
