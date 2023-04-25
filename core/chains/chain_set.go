package chains

import (
	"context"
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	// ErrChainIDEmpty is returned when chain is required but was empty.
	ErrChainIDEmpty = errors.New("chain id empty")
	ErrNotFound     = errors.New("not found")
)

// Chains is a generic interface for chain configuration.
type Chains interface {
	ChainStatus(ctx context.Context, id string) (types.ChainStatus, error)
	ChainStatuses(ctx context.Context, offset, limit int) ([]types.ChainStatus, int, error)
}

// Nodes is an interface for node configuration and state.
type Nodes interface {
	NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []types.NodeStatus, count int, err error)
}

// ChainService is a live, runtime chain instance, with supporting services.
type ChainService interface {
	services.ServiceCtx
	SendTx(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error
}

// ChainSetOpts holds options for configuring a ChainSet via NewChainSet.
type ChainSetOpts[I ID, N Node] interface {
	Validate() error
	ConfigsAndLogger() (Configs[I, N], logger.Logger)
}

type chainSet[N Node, S ChainService] struct {
	utils.StartStopOnce
	opts    ChainSetOpts[string, N]
	configs Configs[string, N]
	lggr    logger.Logger
	chains  map[string]S
}

// NewChainSet returns a new immutable ChainSet for the given ChainSetOpts.
func NewChainSet[N Node, S ChainService](
	chains map[string]S,
	opts ChainSetOpts[string, N],
) (types.ChainSet[string, S], error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	cfgs, lggr := opts.ConfigsAndLogger()
	cs := chainSet[N, S]{
		opts:    opts,
		configs: cfgs,
		lggr:    lggr.Named("ChainSet"),
		chains:  chains,
	}

	return &cs, nil
}

func (c *chainSet[N, S]) Chain(ctx context.Context, id string) (s S, err error) {
	if err = c.StartStopOnce.Ready(); err != nil {
		return
	}
	ch, ok := c.chains[id]
	if !ok {
		err = ErrNotFound
		return
	}
	return ch, nil
}

func (c *chainSet[N, S]) ChainStatus(ctx context.Context, id string) (cfg types.ChainStatus, err error) {
	var cs []types.ChainStatus
	cs, _, err = c.configs.Chains(0, -1, id)
	if err != nil {
		return
	}
	l := len(cs)
	if l == 0 {
		err = ErrNotFound
		return
	}
	if l > 1 {
		err = fmt.Errorf("multiple chains found: %d", len(cs))
		return
	}
	cfg = cs[0]
	return
}

func (c *chainSet[N, S]) ChainStatuses(ctx context.Context, offset, limit int) ([]types.ChainStatus, int, error) {
	return c.configs.Chains(offset, limit)
}

func (c *chainSet[N, S]) NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []types.NodeStatus, count int, err error) {
	return c.configs.NodeStatusesPaged(offset, limit, chainIDs...)
}

func (c *chainSet[N, S]) SendTx(ctx context.Context, chainID, from, to string, amount *big.Int, balanceCheck bool) error {
	chain, err := c.Chain(ctx, chainID)
	if err != nil {
		return err
	}

	return chain.SendTx(ctx, from, to, amount, balanceCheck)
}

func (c *chainSet[N, S]) Start(ctx context.Context) error {
	return c.StartOnce("ChainSet", func() error {
		c.lggr.Debug("Starting")

		var ms services.MultiStart
		for id, ch := range c.chains {
			if err := ms.Start(ctx, ch); err != nil {
				return errors.Wrapf(err, "failed to start chain %q", id)
			}
		}
		c.lggr.Info(fmt.Sprintf("Started %d chains", len(c.chains)))
		return nil
	})
}

func (c *chainSet[N, S]) Close() error {
	return c.StopOnce("ChainSet", func() (err error) {
		c.lggr.Debug("Stopping")

		for _, c := range c.chains {
			err = multierr.Combine(err, c.Close())
		}
		return
	})
}

func (c *chainSet[N, S]) Ready() (err error) {
	err = c.StartStopOnce.Ready()
	for _, c := range c.chains {
		err = multierr.Combine(err, c.Ready())
	}
	return
}

func (c *chainSet[N, S]) Name() string {
	return c.lggr.Name()
}

func (c *chainSet[N, S]) HealthReport() map[string]error {
	report := map[string]error{c.Name(): c.StartStopOnce.Healthy()}
	for _, c := range c.chains {
		maps.Copy(report, c.HealthReport())
	}
	return report
}
