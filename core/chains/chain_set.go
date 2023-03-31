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

// Nodes is a generic interface for Node configuration.
type Nodes[I ID, N Node] interface {
	GetNodes(ctx context.Context, offset, limit int) (nodes []N, count int, err error)
	GetNodesForChain(ctx context.Context, chainID I, offset, limit int) (nodes []N, count int, err error)
}

// ChainSet manages a live set of ChainService instances.
type ChainSet[I ID, N Node, S ChainService] interface {
	services.ServiceCtx
	Chains[I]
	Nodes[I, N]

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

type chainSet[I ID, N Node, S ChainService] struct {
	utils.StartStopOnce
	opts     ChainSetOpts[I, N]
	formatID func(I) string
	orm      Configs[I, N]
	lggr     logger.Logger
	chains   map[string]S
}

// NewChainSet returns a new immutable ChainSet for the given ChainSetOpts.
func NewChainSet[I ID, N Node, S ChainService](chains map[string]S,
	opts ChainSetOpts[I, N], formatID func(I) string,
) (ChainSet[I, N, S], error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	cfgs, lggr := opts.ConfigsAndLogger()
	cs := chainSet[I, N, S]{
		opts:     opts,
		formatID: formatID,
		orm:      cfgs,
		lggr:     lggr.Named("ChainSet"),
		chains:   chains,
	}

	return &cs, nil
}

func (c *chainSet[I, N, S]) Chain(ctx context.Context, id I) (s S, err error) {
	sid := c.formatID(id)
	if sid == "" {
		err = ErrChainIDEmpty
		return
	}
	if err = c.StartStopOnce.Ready(); err != nil {
		return
	}
	ch, ok := c.chains[sid]
	if !ok {
		err = ErrNotFound
		return
	}
	return ch, nil
}

func (c *chainSet[I, N, S]) Show(id I) (cfg ChainConfig, err error) {
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

func (c *chainSet[I, N, S]) Index(offset, limit int) ([]ChainConfig, int, error) {
	return c.orm.Chains(offset, limit)
}

func (c *chainSet[I, N, S]) GetNodes(ctx context.Context, offset, limit int) (nodes []N, count int, err error) {
	return c.orm.Nodes(offset, limit)
}

func (c *chainSet[I, N, S]) GetNodesForChain(ctx context.Context, chainID I, offset, limit int) (nodes []N, count int, err error) {
	return c.orm.NodesForChain(chainID, offset, limit)
}

func (c *chainSet[I, N, S]) Start(ctx context.Context) error {
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

func (c *chainSet[I, N, S]) Close() error {
	return c.StopOnce("ChainSet", func() (err error) {
		c.lggr.Debug("Stopping")

		for _, c := range c.chains {
			err = multierr.Combine(err, c.Close())
		}
		return
	})
}

func (c *chainSet[I, N, S]) Ready() (err error) {
	err = c.StartStopOnce.Ready()
	for _, c := range c.chains {
		err = multierr.Combine(err, c.Ready())
	}
	return
}

func (c *chainSet[I, N, S]) Name() string {
	return c.lggr.Name()
}

func (c *chainSet[I, N, S]) HealthReport() map[string]error {
	report := map[string]error{c.Name(): c.StartStopOnce.Healthy()}
	for _, c := range c.chains {
		maps.Copy(report, c.HealthReport())
	}
	return report
}
