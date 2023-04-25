package evm

import (
	"context"
	"fmt"
	"math/big"
	"net/url"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"
	"go.uber.org/multierr"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	v2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/monitor"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

//go:generate mockery --quiet --name Chain --output ./mocks/ --case=underscore
type Chain interface {
	services.ServiceCtx
	ID() *big.Int
	Client() evmclient.Client
	Config() evmconfig.ChainScopedConfig
	LogBroadcaster() log.Broadcaster
	HeadBroadcaster() httypes.HeadBroadcaster
	TxManager() txmgr.EvmTxManager
	HeadTracker() httypes.HeadTracker
	Logger() logger.Logger
	BalanceMonitor() monitor.BalanceMonitor
	LogPoller() logpoller.LogPoller
	GasEstimator() gas.EvmFeeEstimator

	SendTx(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error
}

var _ Chain = &chain{}

type chain struct {
	utils.StartStopOnce
	id              *big.Int
	cfg             evmconfig.ChainScopedConfig
	client          evmclient.Client
	txm             txmgr.EvmTxManager
	logger          logger.Logger
	headBroadcaster httypes.HeadBroadcaster
	headTracker     httypes.HeadTracker
	logBroadcaster  log.Broadcaster
	logPoller       logpoller.LogPoller
	balanceMonitor  monitor.BalanceMonitor
	keyStore        keystore.Eth
	gasEstimator    gas.EvmFeeEstimator
}

type errChainDisabled struct {
	ChainID *utils.Big
}

func (e errChainDisabled) Error() string {
	return fmt.Sprintf("cannot create new chain with ID %s, the chain is disabled", e.ChainID.String())
}

func newTOMLChain(ctx context.Context, chain *v2.EVMConfig, opts ChainSetOpts) (*chain, error) {
	chainID := chain.ChainID
	l := opts.Logger.With("evmChainID", chainID.String())
	if !chain.IsEnabled() {
		return nil, errChainDisabled{ChainID: chainID}
	}
	cfg := v2.NewTOMLChainScopedConfig(opts.Config, chain, l)
	// note: per-chain validation is not ncessary at this point since everything is checked earlier on boot.
	return newChain(ctx, cfg, chain.Nodes, opts)
}

func newChain(ctx context.Context, cfg evmconfig.ChainScopedConfig, nodes []*v2.Node, opts ChainSetOpts) (*chain, error) {
	chainID := cfg.ChainID()
	l := opts.Logger.Named(chainID.String()).With("evmChainID", chainID.String())
	var client evmclient.Client
	if !cfg.EVMRPCEnabled() {
		client = evmclient.NewNullClient(chainID, l)
	} else if opts.GenEthClient == nil {
		var err2 error
		client, err2 = newEthClientFromChain(cfg, l, cfg.ChainID(), cfg.ChainType(), nodes)
		if err2 != nil {
			return nil, errors.Wrapf(err2, "failed to instantiate eth client for chain with ID %s", cfg.ChainID().String())
		}
	} else {
		client = opts.GenEthClient(chainID)
	}

	db := opts.DB
	headBroadcaster := headtracker.NewHeadBroadcaster(l)
	headSaver := headtracker.NullSaver
	var headTracker httypes.HeadTracker
	if !cfg.EVMRPCEnabled() {
		headTracker = headtracker.NullTracker
	} else if opts.GenHeadTracker == nil {
		orm := headtracker.NewORM(db, l, cfg, *chainID)
		headSaver = headtracker.NewHeadSaver(l, orm, cfg)
		headTracker = headtracker.NewHeadTracker(l, client, cfg, headBroadcaster, headSaver, opts.MailMon)
	} else {
		headTracker = opts.GenHeadTracker(chainID, headBroadcaster)
	}

	logPoller := logpoller.LogPollerDisabled
	if cfg.FeatureLogPoller() {
		if opts.GenLogPoller != nil {
			logPoller = opts.GenLogPoller(chainID)
		} else {
			logPoller = logpoller.NewLogPoller(logpoller.NewORM(chainID, db, l, cfg), client, l, cfg.EvmLogPollInterval(), int64(cfg.EvmFinalityDepth()), int64(cfg.EvmLogBackfillBatchSize()), int64(cfg.EvmRPCDefaultBatchSize()), int64(cfg.EvmLogKeepBlocksDepth()))
		}
	}

	// note: gas estimator is started as a part of the txm
	txm, gasEstimator, err := newEvmTxm(db, cfg, client, l, logPoller, opts)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to instantiate EvmTxm for chain with ID %s", cfg.ChainID().String())
	}

	headBroadcaster.Subscribe(txm)

	// Highest seen head height is used as part of the start of LogBroadcaster backfill range
	highestSeenHead, err := headSaver.LatestHeadFromDB(ctx)
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
		logBroadcaster = log.NewBroadcaster(logORM, client, cfg, l, highestSeenHead, opts.MailMon)
	} else {
		logBroadcaster = opts.GenLogBroadcaster(chainID)
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
		gasEstimator:    gasEstimator,
	}, nil
}

func (c *chain) Start(ctx context.Context) error {
	return c.StartOnce("Chain", func() error {
		c.logger.Debugf("Chain: starting with ID %s", c.ID().String())
		// Must ensure that EthClient is dialed first because subsequent
		// services may make eth calls on startup
		if err := c.client.Dial(ctx); err != nil {
			return errors.Wrap(err, "failed to dial ethclient")
		}
		// Services should be able to handle a non-functional eth client and
		// not block start in this case, instead retrying in a background loop
		// until it becomes available.
		//
		// We do not start the log poller here, it gets
		// started after the jobs so they have a chance to apply their filters.
		var ms services.MultiStart
		if err := ms.Start(ctx, c.txm, c.headBroadcaster, c.headTracker, c.logBroadcaster); err != nil {
			return err
		}
		if c.balanceMonitor != nil {
			if err := ms.Start(ctx, c.balanceMonitor); err != nil {
				return err
			}
		}

		return nil
	})
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
		c.logger.Debug("Chain: stopping evmTxm")
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

func (c *chain) Name() string {
	return c.logger.Name()
}

func (c *chain) HealthReport() map[string]error {
	report := map[string]error{
		c.Name(): c.StartStopOnce.Healthy(),
	}
	maps.Copy(report, c.txm.HealthReport())
	maps.Copy(report, c.headBroadcaster.HealthReport())
	maps.Copy(report, c.headTracker.HealthReport())
	maps.Copy(report, c.logBroadcaster.HealthReport())

	if c.balanceMonitor != nil {
		maps.Copy(report, c.balanceMonitor.HealthReport())
	}

	return report
}

func (c *chain) SendTx(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return errors.New("unsupported") //TODO
}

func (c *chain) ID() *big.Int                             { return c.id }
func (c *chain) Client() evmclient.Client                 { return c.client }
func (c *chain) Config() evmconfig.ChainScopedConfig      { return c.cfg }
func (c *chain) LogBroadcaster() log.Broadcaster          { return c.logBroadcaster }
func (c *chain) LogPoller() logpoller.LogPoller           { return c.logPoller }
func (c *chain) HeadBroadcaster() httypes.HeadBroadcaster { return c.headBroadcaster }
func (c *chain) TxManager() txmgr.EvmTxManager            { return c.txm }
func (c *chain) HeadTracker() httypes.HeadTracker         { return c.headTracker }
func (c *chain) Logger() logger.Logger                    { return c.logger }
func (c *chain) BalanceMonitor() monitor.BalanceMonitor   { return c.balanceMonitor }
func (c *chain) GasEstimator() gas.EvmFeeEstimator        { return c.gasEstimator }

func newEthClientFromChain(cfg evmclient.NodeConfig, lggr logger.Logger, chainID *big.Int, chainType config.ChainType, nodes []*v2.Node) (evmclient.Client, error) {
	var primaries []evmclient.Node
	var sendonlys []evmclient.SendOnlyNode
	for i, node := range nodes {
		if node.SendOnly != nil && *node.SendOnly {
			sendonly := evmclient.NewSendOnlyNode(lggr, (url.URL)(*node.HTTPURL), *node.Name, chainID)
			sendonlys = append(sendonlys, sendonly)
		} else {
			primary, err := newPrimary(cfg, lggr, node, int32(i), chainID)
			if err != nil {
				return nil, err
			}
			primaries = append(primaries, primary)
		}
	}
	return evmclient.NewClientWithNodes(lggr, cfg, primaries, sendonlys, chainID, chainType)
}

func newPrimary(cfg evmclient.NodeConfig, lggr logger.Logger, n *v2.Node, id int32, chainID *big.Int) (evmclient.Node, error) {
	if n.SendOnly != nil && *n.SendOnly {
		return nil, errors.New("cannot cast send-only node to primary")
	}

	return evmclient.NewNode(cfg, lggr, (url.URL)(*n.WSURL), (*url.URL)(n.HTTPURL), *n.Name, id, chainID), nil
}

func EnsureChains(db *sqlx.DB, lggr logger.Logger, cfg pg.QConfig, ids []utils.Big) error {
	q := pg.NewQ(db, lggr.Named("Ensure"), cfg)
	return chains.EnsureChains[utils.Big](q, "evm", ids)
}
