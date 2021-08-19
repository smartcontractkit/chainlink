package evm

import (
	"context"
	"fmt"
	"math/big"
	"net/url"

	"github.com/pkg/errors"
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
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/utils"
	"go.uber.org/multierr"
)

type ChainIdentification interface {
	IsL2() bool
	IsArbitrum() bool
	IsOptimism() bool
}

//go:generate mockery --name Chain --output ./mocks/ --case=underscore
type Chain interface {
	service.Service
	ChainIdentification
	ID() *big.Int
	Client() eth.Client
	Config() evmconfig.ChainScopedConfig
	LogBroadcaster() log.Broadcaster
	HeadBroadcaster() httypes.HeadBroadcaster
	TxManager() bulletprooftxmanager.TxManager
	HeadTracker() httypes.Tracker
	Logger() *logger.Logger
}

var _ Chain = &chain{}

type chain struct {
	utils.StartStopOnce
	id              *big.Int
	cfg             evmconfig.ChainScopedConfig
	client          eth.Client
	txm             bulletprooftxmanager.TxManager
	logger          *logger.Logger
	headBroadcaster httypes.HeadBroadcaster
	headTracker     httypes.Tracker
	logBroadcaster  log.Broadcaster
	balanceMonitor  services.BalanceMonitor
	keyStore        keystore.EthKeyStoreInterface
}

func newChain(dbchain types.Chain, opts ChainCollectionOpts) (*chain, error) {
	// TODO: Pass this logger into all subservices
	chainID := dbchain.ID.ToInt()
	l := opts.Logger.With("chainID", chainID.String())
	cfg := evmconfig.NewChainScopedConfig(opts.DB, l, opts.Config, dbchain)
	if cfg.EVMDisabled() {
		return nil, errors.Errorf("cannot create new chain with ID %d, EVM is disabled", dbchain.ID.ToInt())
	}
	if err := cfg.Validate(); err != nil {
		return nil, errors.Wrapf(err, "cannot create new chain with ID %d, config validation failed", dbchain.ID.ToInt())
	}
	db := opts.DB
	serviceLogLevels, err := opts.Logger.GetServiceLogLevels()
	if err != nil {
		return nil, err
	}
	var client eth.Client
	if cfg.EthereumDisabled() {
		client = &eth.NullClient{CID: chainID}
	} else if opts.GenEthClient == nil {
		var err error
		client, err = newEthClientFromChain(l, dbchain)
		if err != nil {
			return nil, err
		}
	} else {
		client = opts.GenEthClient(dbchain)
	}
	headTrackerLogger, err := opts.Logger.InitServiceLevelLogger(logger.HeadTracker, serviceLogLevels[logger.HeadTracker])
	if err != nil {
		return nil, err
	}
	headBroadcaster := headtracker.NewHeadBroadcaster(l)
	orm := headtracker.NewORM(db, *chainID)
	headTracker := headtracker.NewHeadTracker(headTrackerLogger, client, cfg, orm, headBroadcaster)
	txm := bulletprooftxmanager.NewBulletproofTxManager(db, client, cfg, opts.KeyStore, opts.AdvisoryLocker, opts.EventBroadcaster, l)
	headBroadcaster.Subscribe(txm)

	// Highest seen head height is used as part of the start of LogBroadcaster backfill range
	highestSeenHead, err2 := headTracker.HighestSeenHeadFromDB()
	if err2 != nil {
		return nil, err2
	}

	var balanceMonitor services.BalanceMonitor
	if cfg.BalanceMonitorEnabled() {
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
		if err := c.client.Dial(context.TODO()); err != nil {
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

		if c.cfg.Dev() {
			fundingKeys, err := c.keyStore.FundingKeys()
			if err != nil {
				c.logger.Errorw("Chain: failed to get funding keys")
			} else {
				for _, key := range fundingKeys {
					balance, ethErr := c.client.BalanceAt(context.TODO(), key.Address.Address(), nil)
					if ethErr != nil {
						c.logger.Errorw("Chain: failed to fetch balance for funding key", "address", key.Address, "err", ethErr)
						continue
					}
					if balance.Cmp(big.NewInt(0)) == 0 {
						logger.Infow("The backup funding address does not have sufficient funds", "evmChainID", c.ID(), "address", key.Address.Hex(), "balance", balance)
					} else {
						logger.Infow("Funding address ready", "evmChainID", c.ID(), "address", key.Address.Hex(), "current-balance", balance)
					}
				}
			}
		}

		return nil
	})
}

func (c *chain) Close() error {
	return c.StopOnce("Chain", func() (merr error) {
		if c.balanceMonitor != nil {
			merr = c.balanceMonitor.Close()
		}
		merr = multierr.Combine(
			c.logBroadcaster.Close(),
			c.headTracker.Stop(),
			c.headBroadcaster.Close(),
			c.txm.Close(),
		)
		c.client.Close()
		return merr
	})
}
func (c *chain) ID() *big.Int                              { return c.id }
func (c *chain) Client() eth.Client                        { return c.client }
func (c *chain) Config() evmconfig.ChainScopedConfig       { return c.cfg }
func (c *chain) LogBroadcaster() log.Broadcaster           { return c.logBroadcaster }
func (c *chain) HeadBroadcaster() httypes.HeadBroadcaster  { return c.headBroadcaster }
func (c *chain) TxManager() bulletprooftxmanager.TxManager { return c.txm }
func (c *chain) HeadTracker() httypes.Tracker              { return c.headTracker }
func (c *chain) Logger() *logger.Logger                    { return c.logger }

func (c *chain) IsL2() bool       { return types.IsL2(c.id) }
func (c *chain) IsArbitrum() bool { return types.IsArbitrum(c.id) }
func (c *chain) IsOptimism() bool { return types.IsOptimism(c.id) }

func newEthClientFromChain(lggr *logger.Logger, chain types.Chain) (eth.Client, error) {
	nodes := chain.Nodes
	chainID := big.Int(chain.ID)
	var primary *eth.Node
	var sendonlys []*eth.SecondaryNode
	for _, node := range nodes {
		if node.SendOnly {
			sendonly, err := newSendOnly(node)
			if err != nil {
				return nil, err
			}
			sendonlys = append(sendonlys, sendonly)
		} else {
			if primary != nil {
				return nil, errors.Errorf("Got multiple primaries for chain %d, only one primary is currently supported", chain.ID.ToInt())
			}
			var err error
			primary, err = newPrimary(node)
			if err != nil {
				return nil, err
			}
		}
	}
	if primary == nil {
		return nil, errors.New("no primary node found")
	}
	return eth.NewClientWithNodes(lggr, primary, sendonlys, &chainID)
}

func newPrimary(n types.Node) (*eth.Node, error) {
	if n.SendOnly {
		return nil, errors.New("cannot cast send-only node to primary")
	}
	wsuri, err := url.Parse(n.WSURL)
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

	return eth.NewNode(*wsuri, httpuri, n.Name), nil
}

func newSendOnly(n types.Node) (*eth.SecondaryNode, error) {
	if !n.SendOnly {
		return nil, errors.New("cannot cast non send-only node to secondarynode")
	}
	if !n.HTTPURL.Valid {
		return nil, errors.New("send only node was missing HTTP url")
	}
	httpuri, err := url.Parse(n.HTTPURL.String)
	if err != nil {
		return nil, errors.Wrap(err, "invalid http uri")
	}

	return eth.NewSecondaryNode(*httpuri, n.Name), nil
}
