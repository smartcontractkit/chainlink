package starknet

import (
	"context"
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-relay/pkg/logger"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
	starkchain "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/chain"
	"github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/db"

	"github.com/smartcontractkit/chainlink/v2/core/chains"

	"github.com/smartcontractkit/chainlink/v2/core/chains/starknet/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
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

// TODO remove these wrappers after BCF-2441
type RelayExtender struct {
	starkchain.Chain
	chainImpl *chain
}

var _ relay.RelayerExt = &RelayExtender{}

func NewRelayExtender(cfg *StarknetConfig, opts ChainSetOpts) (*RelayExtender, error) {
	c, err := NewChain(cfg, opts)
	if err != nil {
		return nil, err
	}
	chainImpl, ok := (c).(*chain)
	if !ok {
		return nil, fmt.Errorf("internal error: starkent relay extender got wrong type %t", c)
	}
	return &RelayExtender{Chain: chainImpl, chainImpl: chainImpl}, nil
}
func (r *RelayExtender) GetChainStatus(ctx context.Context) (relaytypes.ChainStatus, error) {
	return r.chainImpl.GetChainStatus(ctx)
}
func (r *RelayExtender) ListNodeStatuses(ctx context.Context, page_size int32, page_token string) (stats []relaytypes.NodeStatus, next_page_token string, total int, err error) {
	return r.chainImpl.ListNodeStatuses(ctx, page_size, page_token)
}
func (r *RelayExtender) Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return chains.ErrLOOPPUnsupported
}
