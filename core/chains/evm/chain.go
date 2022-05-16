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

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/core/chains/evm/headtracker"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/chains/evm/monitor"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/utils"
)

//go:generate mockery --name Chain --output ./mocks/ --case=underscore
type Chain interface {
	services.ServiceCtx
	ID() *big.Int
	Client() evmclient.Client
	Config() evmconfig.ChainScopedConfig
	UpdateConfig(*types.ChainCfg)
	LogBroadcaster() log.Broadcaster
	HeadBroadcaster() httypes.HeadBroadcaster
	TxManager() txmgr.TxManager
	HeadTracker() httypes.HeadTracker
	Logger() logger.Logger
	BalanceMonitor() monitor.BalanceMonitor
	LogPoller() *logpoller.LogPoller
}

var _ Chain = &chain{}

type chain struct {
	utils.StartStopOnce
	id              *big.Int
	cfg             evmconfig.ChainScopedConfig
	client          evmclient.Client
	txm             txmgr.TxManager
	logger          logger.Logger
	headBroadcaster httypes.HeadBroadcaster
	headTracker     httypes.HeadTracker
	logBroadcaster  log.Broadcaster
	logPoller       *logpoller.LogPoller
	balanceMonitor  monitor.BalanceMonitor
	keyStore        keystore.Eth
}

func newChain(dbchain types.DBChain, nodes []types.Node, opts ChainSetOpts) (*chain, error) {
	chainID := dbchain.ID.ToInt()
	l := opts.Logger.With("evmChainID", chainID.String())
	if !dbchain.Enabled {
		return nil, errors.Errorf("cannot create new chain with ID %s, the chain is disabled", dbchain.ID.String())
	}
	cfg := evmconfig.NewChainScopedConfig(chainID, *dbchain.Cfg, opts.ORM, l, opts.Config)
	if err := cfg.Validate(); err != nil {
		return nil, errors.Wrapf(err, "cannot create new chain with ID %s, config validation failed", dbchain.ID.String())
	}
	headTrackerLL := opts.Config.LogLevel().String()
	db := opts.DB
	if db != nil {
		if ll, ok := logger.NewORM(db, l).GetServiceLogLevel(logger.HeadTracker); ok {
			headTrackerLL = ll
		}
	}
	var client evmclient.Client
	if !cfg.EVMRPCEnabled() {
		client = evmclient.NewNullClient(chainID, l)
	} else if opts.GenEthClient == nil {
		var err2 error
		client, err2 = newEthClientFromChain(cfg, l, dbchain, nodes)
		if err2 != nil {
			return nil, errors.Wrapf(err2, "failed to instantiate eth client for chain with ID %s", dbchain.ID.String())
		}
	} else {
		client = opts.GenEthClient(dbchain)
	}

	headBroadcaster := headtracker.NewHeadBroadcaster(l)
	headSaver := headtracker.NullSaver
	var headTracker httypes.HeadTracker
	if !cfg.EVMRPCEnabled() {
		headTracker = headtracker.NullTracker
	} else if opts.GenHeadTracker == nil {
		var ll zapcore.Level
		if err2 := ll.UnmarshalText([]byte(headTrackerLL)); err2 != nil {
			return nil, err2
		}
		headTrackerLogger, err2 := l.NewRootLogger(ll)
		if err2 != nil {
			return nil, errors.Wrapf(err2, "failed to instantiate head tracker for chain with ID %s", dbchain.ID.String())
		}
		orm := headtracker.NewORM(db, l, cfg, *chainID)
		headSaver = headtracker.NewHeadSaver(headTrackerLogger, orm, cfg)
		headTracker = headtracker.NewHeadTracker(headTrackerLogger, client, cfg, headBroadcaster, headSaver)
	} else {
		headTracker = opts.GenHeadTracker(dbchain, headBroadcaster)
	}

	logPoller := logpoller.NewLogPoller(logpoller.NewORM(chainID, db, l, cfg), client, l, cfg.EvmLogPollInterval(), int64(cfg.EvmFinalityDepth()), int64(cfg.EvmLogBackfillBatchSize()))
	if opts.GenLogPoller != nil {
		logPoller = opts.GenLogPoller(dbchain)
	}

	var txm txmgr.TxManager
	if !cfg.EVMRPCEnabled() {
		txm = &txmgr.NullTxManager{ErrMsg: fmt.Sprintf("Ethereum is disabled for chain %d", chainID)}
	} else if opts.GenTxManager == nil {
		checker := &txmgr.CheckerFactory{Client: client}
		txm = txmgr.NewTxm(db, client, cfg, opts.KeyStore, opts.EventBroadcaster, l, checker, logPoller)
	} else {
		txm = opts.GenTxManager(dbchain)
	}

	headBroadcaster.Subscribe(txm)

	// Highest seen head height is used as part of the start of LogBroadcaster backfill range
	highestSeenHead, err := headSaver.LatestHeadFromDB(context.Background())
	if err != nil {
		return nil, err
	}

	var balanceMonitor monitor.BalanceMonitor
	if cfg.EVMRPCEnabled() && cfg.BalanceMonitorEnabled() {
		balanceMonitor = monitor.NewBalanceMonitor(client, opts.KeyStore, l)
		headBroadcaster.Subscribe(balanceMonitor)
	}

	var logBroadcaster log.Broadcaster
	if !cfg.EVMRPCEnabled() {
		logBroadcaster = &log.NullBroadcaster{ErrMsg: fmt.Sprintf("Ethereum is disabled for chain %d", chainID)}
	} else if opts.GenLogBroadcaster == nil {
		logORM := log.NewORM(db, l, cfg, *chainID)
		logBroadcaster = log.NewBroadcaster(logORM, client, cfg, l, highestSeenHead)
	} else {
		logBroadcaster = opts.GenLogBroadcaster(dbchain)
	}

	// AddDependent for this chain
	// log broadcaster will not start until dependent ready is called by a
	// subsequent routine (job spawner)
	logBroadcaster.AddDependents(1)

	headBroadcaster.Subscribe(logBroadcaster)

	return &chain{
		id:              chainID,
		cfg:             cfg,
		client:          client,
		txm:             txm,
		logger:          l,
		headBroadcaster: headBroadcaster,
		headTracker:     headTracker,
		logBroadcaster:  logBroadcaster,
		logPoller:       logPoller,
		balanceMonitor:  balanceMonitor,
		keyStore:        opts.KeyStore,
	}, nil
}

func (c *chain) Start(ctx context.Context) error {
	return c.StartOnce("Chain", func() (merr error) {
		c.logger.Debugf("Chain: starting with ID %s", c.ID().String())
		// Must ensure that EthClient is dialed first because subsequent
		// services may make eth calls on startup
		if err := c.client.Dial(ctx); err != nil {
			return errors.Wrap(err, "failed to dial ethclient")
		}
		// We do not start the log poller here, it gets
		// started after the jobs so they have a chance to apply their filters.
		merr = multierr.Combine(
			c.txm.Start(ctx),
			c.headBroadcaster.Start(ctx),
			c.headTracker.Start(ctx),
			c.logBroadcaster.Start(ctx),
		)
		if c.balanceMonitor != nil {
			merr = multierr.Combine(merr, c.balanceMonitor.Start(ctx))
		}

		if merr != nil {
			return merr
		}

		if !c.cfg.Dev() {
			return nil
		}
		return c.checkKeys(ctx)
	})
}

func (c *chain) checkKeys(ctx context.Context) error {
	fundingKeys, err := c.keyStore.FundingKeys()
	if err != nil {
		return errors.New("failed to get funding keys")
	}
	var wg sync.WaitGroup
	for _, key := range fundingKeys {
		wg.Add(1)
		go func(k ethkey.KeyV2) {
			defer wg.Done()
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
		merr = multierr.Combine(merr, c.headTracker.Close())
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

func (c *chain) ID() *big.Int                             { return c.id }
func (c *chain) Client() evmclient.Client                 { return c.client }
func (c *chain) Config() evmconfig.ChainScopedConfig      { return c.cfg }
func (c *chain) UpdateConfig(cfg *types.ChainCfg)         { c.cfg.Configure(*cfg) }
func (c *chain) LogBroadcaster() log.Broadcaster          { return c.logBroadcaster }
func (c *chain) LogPoller() *logpoller.LogPoller          { return c.logPoller }
func (c *chain) HeadBroadcaster() httypes.HeadBroadcaster { return c.headBroadcaster }
func (c *chain) TxManager() txmgr.TxManager               { return c.txm }
func (c *chain) HeadTracker() httypes.HeadTracker         { return c.headTracker }
func (c *chain) Logger() logger.Logger                    { return c.logger }
func (c *chain) BalanceMonitor() monitor.BalanceMonitor   { return c.balanceMonitor }

func newEthClientFromChain(cfg evmclient.NodeConfig, lggr logger.Logger, chain types.DBChain, nodes []types.Node) (evmclient.Client, error) {
	chainID := big.Int(chain.ID)
	var primaries []evmclient.Node
	var sendonlys []evmclient.SendOnlyNode
	for _, node := range nodes {
		if node.SendOnly {
			sendonly, err := newSendOnly(lggr, node)
			if err != nil {
				return nil, err
			}
			sendonlys = append(sendonlys, sendonly)
		} else {
			primary, err := newPrimary(cfg, lggr, node)
			if err != nil {
				return nil, err
			}
			primaries = append(primaries, primary)
		}
	}
	return evmclient.NewClientWithNodes(lggr, primaries, sendonlys, &chainID)
}

func newPrimary(cfg evmclient.NodeConfig, lggr logger.Logger, n types.Node) (evmclient.Node, error) {
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

	return evmclient.NewNode(cfg, lggr, *wsuri, httpuri, n.Name, n.ID, (*big.Int)(&n.EVMChainID)), nil
}

func newSendOnly(lggr logger.Logger, n types.Node) (evmclient.SendOnlyNode, error) {
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

	return evmclient.NewSendOnlyNode(lggr, *httpuri, n.Name, (*big.Int)(&n.EVMChainID)), nil
}
