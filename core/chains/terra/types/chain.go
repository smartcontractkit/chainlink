package types

import (
	"fmt"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/sqlx"

	"github.com/smartcontractkit/chainlink-terra/pkg/terra"
	terraclient "github.com/smartcontractkit/chainlink-terra/pkg/terra/client"
	terraconfig "github.com/smartcontractkit/chainlink-terra/pkg/terra/config"

	"github.com/smartcontractkit/chainlink/core/chains/terra/terratxm"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
)

// ChainSetOpts holds options for configuring a ChainSet.
type ChainSetOpts struct {
	Config           config.GeneralConfig
	Logger           logger.Logger
	DB               *sqlx.DB
	KeyStore         keystore.Terra
	EventBroadcaster pg.EventBroadcaster
	ORM              ORM
}

var _ terra.ChainSet = (*chainSet)(nil)

type chainSet struct {
	opts   ChainSetOpts
	mu     sync.RWMutex
	chains map[string]terra.Chain
	lggr   logger.Logger
}

// NewChainSet returns a new chain set for opts.
func NewChainSet(opts ChainSetOpts) (terra.ChainSet, error) {
	dbchains, err := opts.ORM.EnabledChainsWithNodes()
	if err != nil {
		return nil, errors.Wrap(err, "error loading chains")
	}
	cs := &chainSet{
		opts:   opts,
		chains: make(map[string]terra.Chain),
		lggr:   opts.Logger.Named("Terra"),
	}
	for _, c := range dbchains {
		n := c.Nodes[0] //TODO client pool
		var err2 error
		cs.chains[c.ID], err2 = NewChain(opts.DB, opts.KeyStore, n, opts.Config, opts.EventBroadcaster, c.Cfg, opts.Logger)
		if err2 != nil {
			err = multierr.Combine(err, err2)
			continue
		}
	}
	return cs, err
}

func (c *chainSet) Get(id string) (terra.Chain, error) {
	c.mu.RLock()
	ch := c.chains[id]
	c.mu.RUnlock()
	if ch != nil {
		return ch, nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	ch = c.chains[id]
	if ch != nil {
		return ch, nil
	}

	opts := c.opts
	chain, err := opts.ORM.Chain(id)
	if err != nil {
		return nil, err
	}
	if len(chain.Nodes) == 0 {
		return nil, fmt.Errorf("no nodes for terra chain: %s", id)
	}

	node := chain.Nodes[0] // TODO client pool
	tChain, err := NewChain(opts.DB, opts.KeyStore, node, opts.Config, opts.EventBroadcaster, chain.Cfg, opts.Logger)
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
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, ch := range c.chains {
		if err := ch.Start(); err != nil {
			c.lggr.Errorw(fmt.Sprintf("EVM: Chain with ID %s failed to start. You will need to fix this issue and restart the Chainlink node before any services that use this chain will work properly. Got error: %v", ch.ID(), err), "evmChainID", ch.ID(), "err", err)
			continue
		}
	}
	return nil
}

func (c *chainSet) Close() (err error) {
	c.lggr.Debug("Stopping")
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, c := range c.chains {
		err = multierr.Combine(err, c.Close())
	}
	return
}

func (c *chainSet) Ready() (err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, c := range c.chains {
		err = multierr.Combine(err, c.Ready())
	}
	return
}

func (c *chainSet) Healthy() (err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, c := range c.chains {
		err = multierr.Combine(err, c.Healthy())
	}
	return
}

var _ terra.Chain = (*chain)(nil)

type chain struct {
	utils.StartStopOnce
	id     string
	cfg    terraconfig.ChainCfg
	client *terraclient.Client
	txm    *terratxm.Txm
	lggr   logger.Logger
}

// NewChain returns a new chain backed by node.
func NewChain(db *sqlx.DB, ks keystore.Terra, node Node, logCfg pg.LogConfig, eb pg.EventBroadcaster, cfg terraconfig.ChainCfg, lggr logger.Logger) (terra.Chain, error) {
	id := node.TerraChainID
	client, err := terraclient.NewClient(id,
		node.TendermintURL, node.FCDURL, 10, lggr)
	if err != nil {
		return nil, err
	}
	txm, err := terratxm.NewTxm(db, client, cfg.FallbackGasPriceULuna, cfg.GasLimitMultiplier, ks, lggr.Named(id), logCfg, eb, 5*time.Second)
	if err != nil {
		return nil, err
	}
	return &chain{
		id:     id,
		cfg:    cfg,
		client: client,
		txm:    txm,
		lggr:   lggr,
	}, nil
}

func (c *chain) ID() string {
	return c.id
}

func (c *chain) Config() terraconfig.ChainCfg {
	return c.cfg
}

func (c *chain) MsgEnqueuer() terra.MsgEnqueuer {
	return c.txm
}

func (c *chain) Reader() terraclient.Reader {
	return c.client
}

func (c *chain) Start() error {
	return c.StartOnce("Chain", func() error {
		//TODO dial client?

		return c.txm.Start()
	})
}

func (c *chain) Close() error {
	return c.StopOnce("Chain", func() error {
		return c.txm.Close()
	})
}

func (c *chain) Ready() error {
	return c.txm.Ready()
}

func (c *chain) Healthy() error {
	return c.txm.Healthy()
}
