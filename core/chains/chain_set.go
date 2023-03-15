package chains

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	// ErrChainIDEmpty is returned when chain is required but was empty.
	ErrChainIDEmpty = errors.New("chain id empty")
	// ErrChainIDInvalid is returned when a chain id does not match any configured chains.
	ErrChainIDInvalid = errors.New("chain id does not match any local chains")
)

// DBChainSet is a generic interface for DBChain[I, C] configuration.
type DBChainSet[I ID, C Config] interface {
	Show(id I) (DBChain[I, C], error)
	Index(offset, limit int) ([]DBChain[I, C], int, error)
}

// DBNodeSet is a generic interface for Node configuration.
type DBNodeSet[I ID, N Node] interface {
	GetNodes(ctx context.Context, offset, limit int) (nodes []N, count int, err error)
	GetNodesForChain(ctx context.Context, chainID I, offset, limit int) (nodes []N, count int, err error)
}

// ChainSet manages a live set of ChainService instances.
type ChainSet[I ID, C Config, N Node, S ChainService[C]] interface {
	services.ServiceCtx

	Name() string
	HealthReport() map[string]error

	// FIXME: for backward compat we will leave this until relayer libs remove Healthy refs
	// https://smartcontract-it.atlassian.net/browse/BCF-2140
	Healthy() error

	DBChainSet[I, C]

	DBNodeSet[I, N]

	// Chain returns the ChainService for this ID (if a configuration is available), creating one if necessary.
	Chain(context.Context, I) (S, error)
}

// ChainService is a live, runtime chain instance, with supporting services.
type ChainService[C Config] interface {
	services.ServiceCtx
}

// ChainSetOpts holds options for configuring a ChainSet via NewChainSet.
type ChainSetOpts[I ID, C Config, N Node, S ChainService[C]] interface {
	Validate() error
	ORMAndLogger() (ORM[I, C, N], logger.Logger)
}

type chainSet[I ID, C Config, N Node, S ChainService[C]] struct {
	utils.StartStopOnce
	opts     ChainSetOpts[I, C, N, S]
	formatID func(I) string
	orm      ORM[I, C, N]
	lggr     logger.Logger

	chainsMu sync.RWMutex
	chains   map[string]S
}

// NewChainSetImmut returns a new immutable ChainSet for the given ChainSetOpts.
func NewChainSetImmut[I ID, C Config, N Node, S ChainService[C]](chains map[string]S,
	opts ChainSetOpts[I, C, N, S], formatID func(I) string,
) (ChainSet[I, C, N, S], error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	orm, lggr := opts.ORMAndLogger()
	cs := chainSet[I, C, N, S]{
		opts:     opts,
		formatID: formatID,
		orm:      orm,
		lggr:     lggr.Named("ChainSet"),
		chains:   chains,
	}

	return &cs, nil
}

func (c *chainSet[I, C, N, S]) Chain(ctx context.Context, id I) (s S, err error) {
	sid := c.formatID(id)
	if sid == "" {
		err = ErrChainIDEmpty
		return
	}
	if err = c.StartStopOnce.Ready(); err != nil {
		return
	}
	c.chainsMu.RLock()
	ch, ok := c.chains[sid]
	c.chainsMu.RUnlock()
	if !ok {
		err = ErrChainIDInvalid
		return
	}
	return ch, nil
}

func (c *chainSet[I, C, N, S]) Show(id I) (DBChain[I, C], error) {
	return c.orm.Chain(id)
}

func (c *chainSet[I, C, N, S]) Index(offset, limit int) ([]DBChain[I, C], int, error) {
	return c.orm.Chains(offset, limit)
}

func (c *chainSet[I, C, N, S]) GetNodes(ctx context.Context, offset, limit int) (nodes []N, count int, err error) {
	return c.orm.Nodes(offset, limit, pg.WithParentCtx(ctx))
}

func (c *chainSet[I, C, N, S]) GetNodesForChain(ctx context.Context, chainID I, offset, limit int) (nodes []N, count int, err error) {
	return c.orm.NodesForChain(chainID, offset, limit, pg.WithParentCtx(ctx))
}

func (c *chainSet[I, C, N, S]) Start(ctx context.Context) error {
	return c.StartOnce("ChainSet", func() error {
		c.lggr.Debug("Starting")

		c.chainsMu.Lock()
		defer c.chainsMu.Unlock()
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

func (c *chainSet[I, C, N, S]) Close() error {
	return c.StopOnce("ChainSet", func() (err error) {
		c.lggr.Debug("Stopping")

		c.chainsMu.Lock()
		defer c.chainsMu.Unlock()
		for _, c := range c.chains {
			err = multierr.Combine(err, c.Close())
		}
		return
	})
}

func (c *chainSet[I, C, N, S]) Ready() (err error) {
	err = c.StartStopOnce.Ready()
	c.chainsMu.RLock()
	defer c.chainsMu.RUnlock()
	for _, c := range c.chains {
		err = multierr.Combine(err, c.Ready())
	}
	return
}

// FIXME: for backward compat we will leave this until relayer libs remove Healthy refs
// https://smartcontract-it.atlassian.net/browse/BCF-2140
func (c *chainSet[I, C, N, S]) Healthy() error {
	return nil
}

func (c *chainSet[I, C, N, S]) Name() string {
	return c.lggr.Name()
}

func (c *chainSet[I, C, N, S]) HealthReport() map[string]error {
	report := map[string]error{c.Name(): c.StartStopOnce.Healthy()}
	c.chainsMu.RLock()
	defer c.chainsMu.RUnlock()
	for _, c := range c.chains {
		utils.MergeMaps(report, c.HealthReport())
	}
	return report
}
