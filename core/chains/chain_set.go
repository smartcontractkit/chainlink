package chains

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	// ErrChainIDEmpty is returned when chain is required but was empty.
	ErrChainIDEmpty = errors.New("chain id empty")
	ErrNotFound     = errors.New("not found")
)

// Chains is a generic interface for ChainConfig[I, C] configuration.
type Chains[I ID] interface {
	Show(id I) (ChainConfig, error)
	Index(offset, limit int) ([]ChainConfig, int, error)
}

// Nodes is an interface for node configuration and state.
type Nodes interface {
	NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []NodeStatus, count int, err error)
}

// ChainSet manages a live set of ChainService instances.
type ChainSet[I ID, S ChainService] interface {
	services.ServiceCtx
	Chains[I]
	Nodes

	Name() string
	HealthReport() map[string]error

	// Chain returns the ChainService for this ID (if a configuration is available), creating one if necessary.
	Chain(context.Context, I) (S, error)
}

// ChainService is a live, runtime chain instance, with supporting services.
type ChainService interface {
	services.ServiceCtx
}

// ChainSetOpts holds options for configuring a ChainSet via NewChainSet.
type ChainSetOpts[I ID, N Node] interface {
	Validate() error
	ConfigsAndLogger() (Configs[I, N], logger.Logger)
}

type chainSet[N Node, S ChainService] struct {
	utils.StartStopOnce
	opts   ChainSetOpts[string, N]
	orm    Configs[string, N]
	lggr   logger.Logger
	chains map[string]S
}

// NewChainSet returns a new immutable ChainSet for the given ChainSetOpts.
func NewChainSet[N Node, S ChainService](
	chains map[string]S,
	opts ChainSetOpts[string, N],
) (ChainSet[string, S], error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	cfgs, lggr := opts.ConfigsAndLogger()
	cs := chainSet[N, S]{
		opts:   opts,
		orm:    cfgs,
		lggr:   lggr.Named("ChainSet"),
		chains: chains,
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

func (c *chainSet[N, S]) Show(id string) (cfg ChainConfig, err error) {
	var cs []ChainConfig
	cs, _, err = c.orm.Chains(0, -1, id)
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

func (c *chainSet[N, S]) Index(offset, limit int) ([]ChainConfig, int, error) {
	return c.orm.Chains(offset, limit)
}

func (c *chainSet[N, S]) NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []NodeStatus, count int, err error) {
	return c.orm.NodeStatusesPaged(offset, limit, chainIDs...)
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
