//go:build terra
// +build terra

package terra

import (
	"fmt"
	"math"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra"
	"github.com/smartcontractkit/chainlink-terra/pkg/terra/db"
	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink/core/chains/terra/types"
	coreconfig "github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// ChainSetOpts holds options for configuring a ChainSet.
type ChainSetOpts struct {
	Config           coreconfig.GeneralConfig
	Logger           logger.Logger
	DB               *sqlx.DB
	KeyStore         keystore.Terra
	EventBroadcaster pg.EventBroadcaster
	ORM              types.ORM
}

func (o ChainSetOpts) validate() (err error) {
	required := func(s string) error {
		return errors.Errorf("%s is required", s)
	}
	if o.Config == nil {
		err = multierr.Append(err, required("Config"))
	}
	if o.Logger == nil {
		err = multierr.Append(err, required("Logger'"))
	}
	if o.DB == nil {
		err = multierr.Append(err, required("DB"))
	}
	if o.KeyStore == nil {
		err = multierr.Append(err, required("KeyStore"))
	}
	if o.EventBroadcaster == nil {
		err = multierr.Append(err, required("EventBroadcaster"))
	}
	if o.ORM == nil {
		err = multierr.Append(err, required("ORM"))
	}
	return
}

func (o ChainSetOpts) newChain(dbchain db.Chain) (*chain, error) {
	return NewChain(o.DB, o.KeyStore, o.Config, o.EventBroadcaster, dbchain, o.Logger)
}

// ChainSet extends terra.ChainSet with mutability and exposes the underlying ORM.
type ChainSet interface {
	terra.ChainSet

	Add(string, db.ChainCfg) (db.Chain, error)
	Remove(string) error
	Configure(id string, enabled bool, config db.ChainCfg) (db.Chain, error)

	ORM() types.ORM
}

//go:generate mockery --name ChainSet --srcpkg github.com/smartcontractkit/chainlink-terra/pkg/terra --output ./mocks/ --case=underscore
var _ ChainSet = (*chainSet)(nil)

type chainSet struct {
	utils.StartStopOnce
	opts     ChainSetOpts
	chainsMu sync.RWMutex
	chains   map[string]*chain
	lggr     logger.Logger
}

// NewChainSet returns a new chain set for opts.
func NewChainSet(opts ChainSetOpts) (*chainSet, error) {
	if err := opts.validate(); err != nil {
		return nil, err
	}
	dbchains, err := opts.ORM.EnabledChainsWithNodes()
	if err != nil {
		return nil, errors.Wrap(err, "error loading chains")
	}
	cs := chainSet{
		opts:   opts,
		chains: make(map[string]*chain),
		lggr:   opts.Logger.Named("ChainSet"),
	}
	for _, dbc := range dbchains {
		var err2 error
		cs.chains[dbc.ID], err2 = opts.newChain(dbc)
		if err2 != nil {
			err = multierr.Combine(err, err2)
			continue
		}
	}
	return &cs, err
}

func (c *chainSet) ORM() types.ORM {
	return c.opts.ORM
}

func (c *chainSet) Chain(id string) (terra.Chain, error) {
	if err := c.StartStopOnce.Ready(); err != nil {
		return nil, err
	}
	c.chainsMu.RLock()
	ch := c.chains[id]
	c.chainsMu.RUnlock()
	if ch != nil {
		// Already known/started
		return ch, nil
	}

	// Unknown/unstarted
	c.chainsMu.Lock()
	defer c.chainsMu.Unlock()

	// Double check now that we have the lock, so we don't start an orphan.
	if err := c.StartStopOnce.Ready(); err != nil {
		return nil, err
	}

	ch = c.chains[id]
	if ch != nil {
		// Someone else beat us to it
		return ch, nil
	}

	// Do we have nodes/config?
	opts := c.opts
	dbchain, err := opts.ORM.Chain(id)
	if err != nil {
		return nil, err
	}

	err = c.initializeChain(&dbchain)
	if err != nil {
		return nil, err
	}
	return c.chains[id], nil
}

func (c *chainSet) Add(id string, config db.ChainCfg) (db.Chain, error) {
	c.chainsMu.Lock()
	defer c.chainsMu.Unlock()

	if _, exists := c.chains[id]; exists {
		return db.Chain{}, errors.Errorf("chain already exists with id %s", id)
	}

	dbchain, err := c.opts.ORM.CreateChain(id, config)
	if err != nil {
		return db.Chain{}, err
	}
	return dbchain, c.initializeChain(&dbchain)
}

// Requires a lock on chainsMu
func (c *chainSet) initializeChain(dbchain *db.Chain) error {
	// preload nodes
	nodes, cnt, err := c.opts.ORM.NodesForChain(dbchain.ID, 0, math.MaxInt)
	if err != nil {
		return err
	}
	if cnt == 0 {
		// Can't start without nodes
		return nil
	}
	dbchain.Nodes = nodes

	// Start it
	cid := dbchain.ID
	chain, err := c.opts.newChain(*dbchain)
	if err != nil {
		return errors.Wrapf(err, "initializeChain: failed to instantiate chain %s", dbchain.ID)
	}
	if err = chain.Start(); err != nil {
		return errors.Wrapf(err, "initializeChain: failed to start chain %s", dbchain.ID)
	}
	c.chains[cid] = chain
	return nil
}

func (c *chainSet) Remove(id string) error {
	c.chainsMu.Lock()
	defer c.chainsMu.Unlock()

	if err := c.opts.ORM.DeleteChain(id); err != nil {
		return err
	}

	chain, exists := c.chains[id]
	if !exists {
		// If a chain was removed from the DB that wasn't loaded into the memory set we're done.
		return nil
	}
	delete(c.chains, id)
	return chain.Close()
}

func (c *chainSet) Configure(id string, enabled bool, config db.ChainCfg) (db.Chain, error) {
	c.chainsMu.Lock()
	defer c.chainsMu.Unlock()

	// Update configuration stored in the database
	dbchain, err := c.opts.ORM.UpdateChain(id, enabled, config)
	if err != nil {
		return db.Chain{}, err
	}

	chain, exists := c.chains[id]

	switch {
	case exists && !enabled:
		// Chain was toggled to disabled
		delete(c.chains, id)
		return db.Chain{}, chain.Close()
	case !exists && enabled:
		// Chain was toggled to enabled
		return dbchain, c.initializeChain(&dbchain)
	case exists:
		// Exists in memory, no toggling: Update in-memory chain
		chain.UpdateConfig(config)
	}

	return dbchain, nil
}

func (c *chainSet) Start() error {
	//TODO if terra disabled, warn and return?
	return c.StartOnce("ChainSet", func() error {
		c.lggr.Debug("Starting")

		c.chainsMu.Lock()
		defer c.chainsMu.Unlock()
		var started int
		for _, ch := range c.chains {
			if err := ch.Start(); err != nil {
				c.lggr.Errorw(fmt.Sprintf("Chain with ID %s failed to start. You will need to fix this issue and restart the Chainlink node before any services that use this chain will work properly. Got error: %v", ch.ID(), err), "err", err)
				continue
			}
			started++
		}
		c.lggr.Info(fmt.Sprintf("Started %d/%d chains", started, len(c.chains)))
		return nil
	})
}

func (c *chainSet) Close() error {
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

func (c *chainSet) Ready() (err error) {
	err = c.StartStopOnce.Ready()
	c.chainsMu.RLock()
	defer c.chainsMu.RUnlock()
	for _, c := range c.chains {
		err = multierr.Combine(err, c.Ready())
	}
	return
}

func (c *chainSet) Healthy() (err error) {
	err = c.StartStopOnce.Healthy()
	c.chainsMu.RLock()
	defer c.chainsMu.RUnlock()
	for _, c := range c.chains {
		err = multierr.Combine(err, c.Healthy())
	}
	return
}
