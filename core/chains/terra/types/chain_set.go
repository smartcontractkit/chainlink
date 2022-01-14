package types

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra"

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
	ORM              ORM
}

func (o ChainSetOpts) validate() (err error) {
	required := func(s string) error {
		return fmt.Errorf("%s is required", s)
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

func (o ChainSetOpts) newChain(dbchain Chain) (*chain, error) {
	return NewChain(o.DB, o.KeyStore, o.Config, o.EventBroadcaster, dbchain, o.Logger)
}

var _ terra.ChainSet = (*chainSet)(nil)

type chainSet struct {
	utils.StartStopOnce
	opts     ChainSetOpts
	chainsMu sync.RWMutex
	chains   map[string]terra.Chain
	lggr     logger.Logger
}

// NewChainSet returns a new chain set for opts.
func NewChainSet(opts ChainSetOpts) (terra.ChainSet, error) {
	if err := opts.validate(); err != nil {
		return nil, err
	}
	dbchains, err := opts.ORM.EnabledChainsWithNodes()
	if err != nil {
		return nil, errors.Wrap(err, "error loading chains")
	}
	cs := chainSet{
		opts:   opts,
		chains: make(map[string]terra.Chain),
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
	if len(dbchain.Nodes) == 0 {
		return nil, fmt.Errorf("no nodes for terra chain: %s", id)
	}

	// Start it
	tChain, err := opts.newChain(dbchain)
	if err != nil {
		return nil, err
	}
	err = tChain.Start()
	if err != nil {
		return nil, err
	}
	c.chains[id] = tChain
	return tChain, nil
}

func (c *chainSet) Start() error {
	//TODO if terra disabled, warn and return?
	return c.StartOnce("ChainSet", func() error {
		c.lggr.Debug("Starting")

		c.chainsMu.RLock()
		defer c.chainsMu.RUnlock()
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
		c.chainsMu.RLock()
		defer c.chainsMu.RUnlock()
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
