package evm

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"

	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name Chain --output ./mocks/ --case=underscore
type Chain interface {
	service.Service
	ID() *big.Int
	Client() eth.Client
	Config() evmconfig.ChainScopedConfig
	LogBroadcaster() log.Broadcaster
	HeadBroadcaster() httypes.HeadBroadcaster
	TxManager() bulletprooftxmanager.TxManager
	HeadTracker() httypes.Tracker
	Logger() logger.Logger
}

var _ Chain = &chain{}

type chain struct {
	utils.StartStopOnce
	id              *big.Int
	cfg             evmconfig.ChainScopedConfig
	client          eth.Client
	txm             bulletprooftxmanager.TxManager
	logger          logger.Logger
	headBroadcaster httypes.HeadBroadcaster
	headTracker     httypes.Tracker
	logBroadcaster  log.Broadcaster
	balanceMonitor  services.BalanceMonitor
	keyStore        keystore.Eth
}

func newChain(dbchain types.Chain, opts ChainSetOpts) (*chain, error) {
	chainID := dbchain.ID.ToInt()
	l := opts.Logger.With("evmChainID", chainID.String())
	cfg := evmconfig.NewChainScopedConfig(chainID, dbchain.Cfg, opts.ORM, l, opts.Config)
	if cfg.EVMDisabled() {
		return nil, errors.Errorf("cannot create new chain with ID %s, EVM is disabled", dbchain.ID.String())
	}
	if !dbchain.Enabled {
		return nil, errors.Errorf("cannot create new chain with ID %s, the chain is disabled", dbchain.ID.String())
	}
	if err := cfg.Validate(); err != nil {
		return nil, errors.Wrapf(err, "cannot create new chain with ID %s, config validation failed", dbchain.ID.String())
	}
	db := opts.GormDB
	headTrackerLL := opts.Config.LogLevel().String()
	if db != nil {
		if ll, ok := logger.NewORM(db).GetServiceLogLevel(logger.HeadTracker); ok {
			headTrackerLL = ll
		}
	}
	var client eth.Client
	if cfg.EthereumDisabled() {
		client = &eth.NullClient{CID: chainID}
	} else if opts.GenEthClient == nil {
		var err2 error
		client, err2 = newEthClientFromChain(l, dbchain)
		if err2 != nil {
			return nil, errors.Wrapf(err2, "failed to instantiate eth client for chain with ID %s", dbchain.ID.String())
		}
	} else {
		client = opts.GenEthClient(dbchain)
	}

	headBroadcaster := headtracker.NewHeadBroadcaster(l)
	var headTracker httypes.Tracker
	if cfg.EthereumDisabled() {
		headTracker = &headtracker.NullTracker{}
	} else if opts.GenHeadTracker == nil {
		var ll zapcore.Level
		if err2 := ll.UnmarshalText([]byte(headTrackerLL)); err2 != nil {
			return nil, err2
		}
		headTrackerLogger, err2 := l.NewRootLogger(ll)
		if err2 != nil {
			return nil, errors.Wrapf(err2, "failed to instantiate head tracker for chain with ID %s", dbchain.ID.String())
		}
		orm := headtracker.NewORM(db, *chainID)
		headTracker = headtracker.NewHeadTracker(headTrackerLogger, client, cfg, orm, headBroadcaster)
	} else {
		headTracker = opts.GenHeadTracker(dbchain)
	}

	var txm bulletprooftxmanager.TxManager
	if cfg.EthereumDisabled() {
		txm = &bulletprooftxmanager.NullTxManager{ErrMsg: fmt.Sprintf("Ethereum is disabled for chain %d", chainID)}
	} else if opts.GenTxManager == nil {
		txm = bulletprooftxmanager.NewBulletproofTxManager(db, client, cfg, opts.KeyStore, opts.EventBroadcaster, l)
	} else {
		txm = opts.GenTxManager(dbchain)
	}

	headBroadcaster.Subscribe(txm)

	// Highest seen head height is used as part of the start of LogBroadcaster backfill range
	highestSeenHead, err := headTracker.HighestSeenHeadFromDB(context.Background())
	if err != nil {
		return nil, err
	}

	var balanceMonitor services.BalanceMonitor
	if !cfg.EthereumDisabled() && cfg.BalanceMonitorEnabled() {
		balanceMonitor = services.NewBalanceMonitor(db, client, opts.KeyStore, l)
		headBroadcaster.Subscribe(balanceMonitor)
	}

	var logBroadcaster log.Broadcaster
	if cfg.EthereumDisabled() {
		logBroadcaster = &log.NullBroadcaster{ErrMsg: fmt.Sprintf("Ethereum is disabled for chain %d", chainID)}
	} else if opts.GenLogBroadcaster == nil {
		logBroadcaster = log.NewBroadcaster(log.NewORM(db, *chainID), client, cfg, l, highestSeenHead)
	} else {
		logBroadcaster = opts.GenLogBroadcaster(dbchain)
	}

	// Log Broadcaster waits for other services' registrations
	// until app.LogBroadcaster.DependentReady() call (see below)
	logBroadcaster.AddDependents(1)

	headBroadcaster.Subscribe(logBroadcaster)

	c := chain{
		utils.StartStopOnce{},
		chainID,
		cfg,
		client,
		txm,
		l,
		headBroadcaster,
		headTracker,
		logBroadcaster,
		balanceMonitor,
		opts.KeyStore,
	}
	return &c, nil
}

func (c *chain) Start() error {
	return c.StartOnce("Chain", func() (merr error) {
		c.logger.Debugf("Chain: starting with ID %s", c.ID().String())
		// Must ensure that EthClient is dialed first because subsequent
		// services may make eth calls on startup
		ctx, cancel := eth.DefaultQueryCtx()
		defer cancel()
		if err := c.client.Dial(ctx); err != nil {
			return errors.Wrap(err, "failed to Dial ethclient")
		}
		merr = multierr.Combine(
			c.txm.Start(),
			c.headBroadcaster.Start(),
			c.headTracker.Start(),
			c.logBroadcaster.Start(),
		)
		if c.balanceMonitor != nil {
			merr = multierr.Combine(merr, c.balanceMonitor.Start())
		}

		if merr != nil {
			return merr
		}

		// Log Broadcaster fully starts after all initial Register calls are done from other starting services
		// to make sure the initial backfill covers those subscribers.
		c.logBroadcaster.DependentReady()

		if !c.cfg.Dev() {
			return nil
		}
		return c.checkKeys()
	})
}

func (c *chain) checkKeys() error {
	fundingKeys, err := c.keyStore.FundingKeys()
	if err != nil {
		return errors.New("failed to get funding keys")
	}
	var wg sync.WaitGroup
	for _, key := range fundingKeys {
		wg.Add(1)
		go func(k ethkey.KeyV2) {
			defer wg.Done()
			ctx, cancel := eth.DefaultQueryCtx()
			defer cancel()
			balance, ethErr := c.client.BalanceAt(ctx, k.Address.Address(), nil)
			if ethErr != nil {
				c.logger.Errorw("Chain: failed to fetch balance for funding key", "address", k.Address, "err", ethErr)
				return
			}
			if balance.Cmp(big.NewInt(0)) == 0 {
				c.logger.Infow("The backup funding address does not have sufficient funds", "address", k.Address.Hex(), "balance", balance)
			} else {
				c.logger.Infow("Funding address ready", "address", k.Address.Hex(), "current-balance", balance)
			}
		}(key)
	}
	wg.Wait()

	return nil
}

func (c *chain) Close() error {
	return c.StopOnce("Chain", func() (merr error) {
		c.logger.Debug("Chain: stopping")

		if c.balanceMonitor != nil {
			c.logger.Debug("Chain: stopping balance monitor")
			merr = c.balanceMonitor.Close()
		}
		c.logger.Debug("Chain: stopping logBroadcaster")
		merr = multierr.Combine(merr, c.logBroadcaster.Close())
		c.logger.Debug("Chain: stopping headTracker")
		merr = multierr.Combine(merr, c.headTracker.Stop())
		c.logger.Debug("Chain: stopping headBroadcaster")
		merr = multierr.Combine(merr, c.headBroadcaster.Close())
		c.logger.Debug("Chain: stopping txm")
		merr = multierr.Combine(merr, c.txm.Close())
		c.logger.Debug("Chain: stopping client")
		c.client.Close()
		c.logger.Debug("Chain: stopped")
		return merr
	})
}

func (c *chain) Ready() (merr error) {
	merr = multierr.Combine(
		c.StartStopOnce.Ready(),
		c.txm.Ready(),
		c.headBroadcaster.Ready(),
		c.headTracker.Ready(),
		c.logBroadcaster.Ready(),
	)
	if c.balanceMonitor != nil {
		merr = multierr.Combine(merr, c.balanceMonitor.Ready())
	}
	return
}

func (c *chain) Healthy() (merr error) {
	merr = multierr.Combine(
		c.StartStopOnce.Healthy(),
		c.txm.Healthy(),
		c.headBroadcaster.Healthy(),
		c.headTracker.Healthy(),
		c.logBroadcaster.Healthy(),
	)
	if c.balanceMonitor != nil {
		merr = multierr.Combine(merr, c.balanceMonitor.Healthy())
	}
	return
}

func (c *chain) ID() *big.Int                              { return c.id }
func (c *chain) Client() eth.Client                        { return c.client }
func (c *chain) Config() evmconfig.ChainScopedConfig       { return c.cfg }
func (c *chain) LogBroadcaster() log.Broadcaster           { return c.logBroadcaster }
func (c *chain) HeadBroadcaster() httypes.HeadBroadcaster  { return c.headBroadcaster }
func (c *chain) TxManager() bulletprooftxmanager.TxManager { return c.txm }
func (c *chain) HeadTracker() httypes.Tracker              { return c.headTracker }
func (c *chain) Logger() logger.Logger                     { return c.logger }

var ErrNoPrimaryNode = errors.New("no primary node found")

func newEthClientFromChain(lggr logger.Logger, chain types.Chain) (eth.Client, error) {
	nodes := chain.Nodes
	chainID := big.Int(chain.ID)
	var primaries []eth.Node
	var sendonlys []eth.SendOnlyNode
	for _, node := range nodes {
		if node.SendOnly {
			sendonly, err := newSendOnly(lggr, node)
			if err != nil {
				return nil, err
			}
			sendonlys = append(sendonlys, sendonly)
		} else {
			primary, err := newPrimary(lggr, node)
			if err != nil {
				return nil, err
			}
			primaries = append(primaries, primary)
		}
	}
	if len(primaries) == 0 {
		return nil, ErrNoPrimaryNode
	}
	return eth.NewClientWithNodes(lggr, primaries, sendonlys, &chainID)
}

func newPrimary(lggr logger.Logger, n types.Node) (eth.Node, error) {
	if n.SendOnly {
		return nil, errors.New("cannot cast send-only node to primary")
	}
	if !n.WSURL.Valid {
		return nil, errors.New("primary node was missing WS url")
	}
	wsuri, err := url.Parse(n.WSURL.String)
	if err != nil {
		return nil, errors.Wrap(err, "invalid websocket uri")
	}
	var httpuri *url.URL
	if n.HTTPURL.Valid {
		u, err := url.Parse(n.HTTPURL.String)
		if err != nil {
			return nil, errors.Wrap(err, "invalid http uri")
		}
		httpuri = u
	}

	return eth.NewNode(lggr, *wsuri, httpuri, n.Name, n.EVMChainID.ToInt()), nil
}

func newSendOnly(lggr logger.Logger, n types.Node) (eth.SendOnlyNode, error) {
	if !n.SendOnly {
		return nil, errors.New("cannot cast non send-only node to send-only node")
	}
	if !n.HTTPURL.Valid {
		return nil, errors.New("send only node was missing HTTP url")
	}
	httpuri, err := url.Parse(n.HTTPURL.String)
	if err != nil {
		return nil, errors.Wrap(err, "invalid http uri")
	}

	return eth.NewSendOnlyNode(lggr, *httpuri, n.Name, n.EVMChainID.ToInt()), nil
}
