package chains

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	cfgv2 "github.com/smartcontractkit/chainlink/core/config/v2"
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
	Add(ctx context.Context, id I, cfg C) (DBChain[I, C], error)
	Show(id I) (DBChain[I, C], error)
	Configure(ctx context.Context, id I, enabled bool, cfg C) (DBChain[I, C], error)
	Remove(id I) error
	Index(offset, limit int) ([]DBChain[I, C], int, error)
}

// DBNodeSet is a generic interface for Node configuration.
type DBNodeSet[I ID, N Node] interface {
	GetNodes(ctx context.Context, offset, limit int) (nodes []N, count int, err error)
	GetNodesForChain(ctx context.Context, chainID I, offset, limit int) (nodes []N, count int, err error)
	CreateNode(context.Context, N) (N, error)
	DeleteNode(context.Context, int32) error
}

// ChainSet manages a live set of ChainService instances.
type ChainSet[I ID, C Config, N Node, S ChainService[C]] interface {
	services.ServiceCtx

	Name() string
	HealthReport() map[string]error

	DBChainSet[I, C]

	DBNodeSet[I, N]

	// Chain returns the ChainService for this ID (if a configuration is available), creating one if necessary.
	Chain(context.Context, I) (S, error)
}

// ChainService is a live, runtime chain instance, with supporting services.
type ChainService[C Config] interface {
	services.ServiceCtx
	UpdateConfig(C)
}

// ChainSetOpts holds options for configuring a ChainSet via NewChainSet.
type ChainSetOpts[I ID, C Config, N Node, S ChainService[C]] interface {
	Validate() error
	// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
	NewChain(DBChain[I, C]) (S, error)
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

	// immutability will be standard https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
	immutable bool // toml chain set is immutable
}

// NewChainSet returns a new ChainSet for the given ChainSetOpts.
// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func NewChainSet[I ID, C Config, N Node, S ChainService[C]](
	opts ChainSetOpts[I, C, N, S], formatID func(I) string,
) (ChainSet[I, C, N, S], error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	orm, lggr := opts.ORMAndLogger()
	dbchains, err := orm.EnabledChains()
	if err != nil {
		return nil, errors.Wrap(err, "error loading chains")
	}
	cs := chainSet[I, C, N, S]{
		opts:     opts,
		formatID: formatID,
		orm:      orm,
		lggr:     lggr.Named("ChainSet"),
		chains:   make(map[string]S),
	}
	for _, dbc := range dbchains {
		var err2 error
		cs.chains[formatID(dbc.ID)], err2 = opts.NewChain(dbc)
		if err2 != nil {
			err = multierr.Combine(err, err2)
			continue
		}
	}

	return &cs, err
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
		opts:      opts,
		formatID:  formatID,
		orm:       orm,
		lggr:      lggr.Named("ChainSet"),
		chains:    chains,
		immutable: true,
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
	if ok {
		// Already known/started
		return ch, nil
	}

	// Unknown/unstarted
	c.chainsMu.Lock()
	defer c.chainsMu.Unlock()

	// Double check now that we have the lock, so we don't start an orphan.
	if err = c.StartStopOnce.Ready(); err != nil {
		return
	}

	ch, ok = c.chains[sid]
	if ok {
		// Someone else beat us to it
		return ch, nil
	}

	// Do we have nodes/config?
	var dbchain DBChain[I, C]
	dbchain, err = c.orm.Chain(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = ErrChainIDInvalid
		}
		return
	}

	err = c.initializeChain(ctx, dbchain)
	if err != nil {
		return
	}
	return c.chains[sid], nil
}

// Requires a lock on chainsMu
func (c *chainSet[I, C, N, S]) initializeChain(ctx context.Context, dbchain DBChain[I, C]) error {
	cid := c.formatID(dbchain.ID)
	chain, err := c.opts.NewChain(dbchain)
	if err != nil {
		return errors.Wrapf(err, "initializeChain: failed to instantiate chain %s", cid)
	}
	if err = chain.Start(ctx); err != nil {
		return errors.Wrapf(err, "initializeChain: failed to start chain %s", cid)
	}
	c.chains[cid] = chain
	return nil
}

func (c *chainSet[I, C, N, S]) Add(ctx context.Context, id I, config C) (DBChain[I, C], error) {
	if c.immutable {
		return DBChain[I, C]{}, cfgv2.ErrUnsupported
	}
	c.chainsMu.Lock()
	defer c.chainsMu.Unlock()

	sid := c.formatID(id)
	if _, exists := c.chains[sid]; exists {
		return DBChain[I, C]{}, errors.Errorf("chain already exists with id %s", sid)
	}

	dbchain, err := c.orm.CreateChain(id, config)
	if err != nil {
		return DBChain[I, C]{}, err
	}
	return dbchain, c.initializeChain(ctx, dbchain)
}

func (c *chainSet[I, C, N, S]) Show(id I) (DBChain[I, C], error) {
	return c.orm.Chain(id)
}

func (c *chainSet[I, C, N, S]) Configure(ctx context.Context, id I, enabled bool, config C) (DBChain[I, C], error) {
	if c.immutable {
		return DBChain[I, C]{}, cfgv2.ErrUnsupported
	}
	c.chainsMu.Lock()
	defer c.chainsMu.Unlock()

	// Update configuration stored in the database
	dbchain, err := c.orm.UpdateChain(id, enabled, config)
	if err != nil {
		return DBChain[I, C]{}, err
	}

	sid := c.formatID(id)
	chain, exists := c.chains[sid]

	switch {
	case exists && !enabled:
		// Chain was toggled to disabled
		delete(c.chains, sid)
		return DBChain[I, C]{}, chain.Close()
	case !exists && enabled:
		// Chain was toggled to enabled
		return dbchain, c.initializeChain(ctx, dbchain)
	case exists:
		// Exists in memory, no toggling: Update in-memory chain
		chain.UpdateConfig(config)
	}

	return dbchain, nil
}

func (c *chainSet[I, C, N, S]) Remove(id I) error {
	if c.immutable {
		return cfgv2.ErrUnsupported
	}
	c.chainsMu.Lock()
	defer c.chainsMu.Unlock()

	if err := c.orm.DeleteChain(id); err != nil {
		return err
	}

	sid := c.formatID(id)
	chain, exists := c.chains[sid]
	if !exists {
		// If a chain was removed from the DB that wasn't loaded into the memory set we're done.
		return nil
	}
	delete(c.chains, sid)
	return chain.Close()
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

func (c *chainSet[I, C, N, S]) CreateNode(ctx context.Context, n N) (N, error) {
	return c.orm.CreateNode(n, pg.WithParentCtx(ctx))
}

func (c *chainSet[I, C, N, S]) DeleteNode(ctx context.Context, id int32) error {
	if c.immutable {
		return cfgv2.ErrUnsupported
	}
	return c.orm.DeleteNode(id, pg.WithParentCtx(ctx))
}

func (c *chainSet[I, C, N, S]) Start(ctx context.Context) error {
	return c.StartOnce("ChainSet", func() error {
		c.lggr.Debug("Starting")

		c.chainsMu.Lock()
		defer c.chainsMu.Unlock()
		if c.immutable {
			var ms services.MultiStart
			for id, ch := range c.chains {
				if err := ms.Start(ctx, ch); err != nil {
					return errors.Wrapf(err, "failed to start chain %q", id)
				}
			}
			c.lggr.Info(fmt.Sprintf("Started %d chains", len(c.chains)))
		} else {
			var started int
			for id, ch := range c.chains {
				if err := ch.Start(ctx); err != nil {
					c.lggr.Errorw(fmt.Sprintf("Chain with ID %s failed to start. You will need to fix this issue and restart the Chainlink node before any services that use this chain will work properly.", id), "err", err)
					continue
				}
				started++
			}
			c.lggr.Info(fmt.Sprintf("Started %d/%d chains", started, len(c.chains)))
		}
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

func (c *chainSet[I, C, N, S]) Healthy() (err error) {
	err = c.StartStopOnce.Healthy()
	c.chainsMu.RLock()
	defer c.chainsMu.RUnlock()
	for _, c := range c.chains {
		err = multierr.Combine(err, c.Healthy())
	}
	return
}

func (c *chainSet[I, C, N, S]) Name() string {
	return c.lggr.Name()
}

func (c *chainSet[I, C, N, S]) HealthReport() map[string]error {
	return map[string]error{c.Name(): c.Healthy()}
}
