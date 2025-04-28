package legacyevm

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	gotoml "github.com/pelletier/go-toml/v2"
	"go.uber.org/multierr"

	common "github.com/smartcontractkit/chainlink-common/pkg/chains"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/chaintype"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas/rollups"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads"
	"github.com/smartcontractkit/chainlink-evm/pkg/keys"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/monitor"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"
	trontxm "github.com/smartcontractkit/chainlink-tron/relayer/txm"
	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/tron"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
)

type Chain interface {
	types.ChainService

	ID() *big.Int
	Client() client.Client
	Config() config.ChainScopedConfig
	LogBroadcaster() log.Broadcaster
	HeadBroadcaster() heads.Broadcaster
	TxManager() txmgr.TxManager
	HeadTracker() heads.Tracker
	Logger() logger.Logger
	BalanceMonitor() monitor.BalanceMonitor
	LogPoller() logpoller.LogPoller
	GasEstimator() gas.EvmFeeEstimator
}

// ChainTronSupport is an Chain interface extension for Tron support.
type ChainTronSupport interface {
	GetTronTXM() *trontxm.TronTxm
}

var (
	_           Chain = &chain{}
	nilBigInt   *big.Int
	emptyString string
)

// LegacyChains implements [LegacyChainContainer]
type LegacyChains struct {
	*chains.ChainsKV[Chain]

	cfgs toml.EVMConfigs
}

// LegacyChainContainer is container for EVM chains.
type LegacyChainContainer interface {
	Get(id string) (Chain, error)
	Len() int
	List(ids ...string) ([]Chain, error)
	Slice() []Chain

	// BCF-2516: this is only used for EVMORM. When we delete that
	// we can promote/move the needed funcs from it to LegacyChainContainer
	// so instead of EVMORM().XYZ() we'd have something like legacyChains.XYZ()
	ChainNodeConfigs() evmtypes.Configs
}

var _ LegacyChainContainer = &LegacyChains{}

func NewLegacyChains(m map[string]Chain, evmCfgs toml.EVMConfigs) *LegacyChains {
	return &LegacyChains{
		ChainsKV: chains.NewChainsKV[Chain](m),
		cfgs:     evmCfgs,
	}
}

func (c *LegacyChains) ChainNodeConfigs() evmtypes.Configs {
	return c.cfgs
}

// backward compatibility.
// eth keys are represented as multiple types in the code base;
// *big.Int, string, and int64.
//
// TODO BCF-2507 unify the type system
func (c *LegacyChains) Get(id string) (Chain, error) {
	if id == nilBigInt.String() || id == emptyString {
		return nil, fmt.Errorf("invalid chain id requested: %q", id)
	}
	return c.ChainsKV.Get(id)
}

type chain struct {
	services.StateMachine
	id              *big.Int
	cfg             *config.ChainScoped
	client          client.Client
	txm             txmgr.TxManager
	logger          logger.Logger
	headBroadcaster heads.Broadcaster
	headTracker     heads.Tracker
	logBroadcaster  log.Broadcaster
	logPoller       logpoller.LogPoller
	balanceMonitor  monitor.BalanceMonitor
	gasEstimator    gas.EvmFeeEstimator

	// Extends with support for the Tron TXM
	tronTxm *trontxm.TronTxm
}

type errChainDisabled struct {
	ChainID *ubig.Big
}

func (e errChainDisabled) Error() string {
	return fmt.Sprintf("cannot create new chain with ID %s, the chain is disabled", e.ChainID.String())
}

type FeatureConfig interface {
	LogPoller() bool
}

type ChainRelayOpts struct {
	Logger   logger.Logger
	KeyStore keys.ChainStore
	ChainOpts
}

type ChainOpts struct {
	ChainConfigs   toml.EVMConfigs
	DatabaseConfig txmgr.DatabaseConfig
	FeatureConfig  FeatureConfig
	ListenerConfig txmgr.ListenerConfig

	MailMon      *mailbox.Monitor
	GasEstimator gas.EvmFeeEstimator

	DS sqlutil.DataSource

	// TODO BCF-2513 remove test code from the API
	// Gen-functions are useful for dependency injection by tests
	GenChainStore     func(ks core.Keystore, i *big.Int) keys.ChainStore
	GenEthClient      func(*big.Int) client.Client
	GenLogBroadcaster func(*big.Int) log.Broadcaster
	GenLogPoller      func(*big.Int) logpoller.LogPoller
	GenHeadTracker    func(*big.Int, heads.Broadcaster) heads.Tracker
	GenTxManager      func(*big.Int) txmgr.TxManager
	GenGasEstimator   func(*big.Int) gas.EvmFeeEstimator
}

func (o ChainOpts) Validate() error {
	var err error
	if o.ChainConfigs == nil {
		err = errors.Join(err, errors.New("nil ChainConfigs"))
	}
	if o.DatabaseConfig == nil {
		err = errors.Join(err, errors.New("nil DatabaseConfig"))
	}
	if o.FeatureConfig == nil {
		err = errors.Join(err, errors.New("nil FeatureConfig"))
	}
	if o.ListenerConfig == nil {
		err = errors.Join(err, errors.New("nil ListenerConfig"))
	}
	if o.MailMon == nil {
		err = errors.Join(err, errors.New("nil MailMon"))
	}
	if o.DS == nil {
		err = errors.Join(err, errors.New("nil DS"))
	}
	if err != nil {
		err = fmt.Errorf("invalid ChainOpts: %w", err)
	}
	return err
}

func NewTOMLChain(chain *toml.EVMConfig, opts ChainRelayOpts, clientsByChainID map[string]rollups.DAClient) (Chain, error) {
	err := opts.Validate()
	if err != nil {
		return nil, err
	}
	chainID := chain.ChainID
	if !chain.IsEnabled() {
		return nil, errChainDisabled{ChainID: chainID}
	}
	cfg := config.NewTOMLChainScopedConfig(chain)
	// note: per-chain validation is not necessary at this point since everything is checked earlier on boot.
	return newChain(cfg, chain.Nodes, opts, clientsByChainID)
}

func newChain(cfg *config.ChainScoped, nodes []*toml.Node, opts ChainRelayOpts, clientsByChainID map[string]rollups.DAClient) (*chain, error) {
	chainID := cfg.EVM().ChainID()
	l := opts.Logger
	var cl client.Client
	var err error
	if !opts.ChainConfigs.RPCEnabled() {
		cl = client.NewNullClient(chainID, l)
	} else if opts.GenEthClient == nil {
		cl, err = client.NewEvmClient(cfg.EVM().NodePool(), cfg.EVM(), cfg.EVM().NodePool().Errors(), l, chainID, nodes, cfg.EVM().ChainType())
		if err != nil {
			return nil, err
		}
	} else {
		cl = opts.GenEthClient(chainID)
	}

	headBroadcaster := heads.NewBroadcaster(l)
	headSaver := heads.NullSaver
	var headTracker heads.Tracker
	if !opts.ChainConfigs.RPCEnabled() {
		headTracker = heads.NullTracker
	} else if opts.GenHeadTracker == nil {
		var orm heads.ORM
		if cfg.EVM().HeadTracker().PersistenceEnabled() {
			orm = heads.NewORM(*chainID, opts.DS)
		} else {
			orm = heads.NewNullORM()
		}
		headSaver = heads.NewSaver(l, orm, cfg.EVM(), cfg.EVM().HeadTracker())
		headTracker = heads.NewTracker(l, cl, cfg.EVM(), cfg.EVM().HeadTracker(), headBroadcaster, headSaver, opts.MailMon)
	} else {
		headTracker = opts.GenHeadTracker(chainID, headBroadcaster)
	}

	logPoller := logpoller.LogPollerDisabled
	if opts.FeatureConfig.LogPoller() {
		if opts.GenLogPoller != nil {
			logPoller = opts.GenLogPoller(chainID)
		} else {
			lpOpts := logpoller.Opts{
				PollPeriod:               cfg.EVM().LogPollInterval(),
				UseFinalityTag:           cfg.EVM().FinalityTagEnabled(),
				FinalityDepth:            int64(cfg.EVM().FinalityDepth()),
				BackfillBatchSize:        int64(cfg.EVM().LogBackfillBatchSize()),
				RPCBatchSize:             int64(cfg.EVM().RPCDefaultBatchSize()),
				KeepFinalizedBlocksDepth: int64(cfg.EVM().LogKeepBlocksDepth()),
				LogPrunePageSize:         int64(cfg.EVM().LogPrunePageSize()),
				BackupPollerBlockDelay:   int64(cfg.EVM().BackupLogPollerBlockDelay()),
				ClientErrors:             cfg.EVM().NodePool().Errors(),
			}

			lpORM, err := logpoller.NewObservedORM(chainID, opts.DS, l)
			if err != nil {
				return nil, fmt.Errorf("failed to create logpoller observed ORM: %w", err)
			}

			logPoller = logpoller.NewLogPoller(lpORM, cl, l, headTracker, lpOpts)
		}
	}

	// initialize gas estimator
	gasEstimator, err := newGasEstimator(cfg.EVM(), cl, l, opts, clientsByChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate gas estimator for chain with ID %s: %w", chainID, err)
	}

	// note: gas estimator is started as a part of the txm
	var txm txmgr.TxManager
	//nolint:gocritic // ignoring suggestion to convert to switch statement
	if !opts.ChainConfigs.RPCEnabled() {
		txm = &txmgr.NullTxManager{ErrMsg: fmt.Sprintf("Ethereum is disabled for chain %d", chainID)}
	} else if cfg.EVM().ChainType() == chaintype.ChainTron {
		txm = &txmgr.NullTxManager{ErrMsg: fmt.Sprintf("EVM TXM disabled for tron based chains %d, using Tron TXM instead", chainID)}
	} else if !cfg.EVM().Transactions().Enabled() {
		txm = &txmgr.NullTxManager{ErrMsg: fmt.Sprintf("TXM disabled for chain %d", chainID)}
	} else {
		txm, err = newEvmTxm(opts.DS, cfg.EVM(), opts.DatabaseConfig, opts.ListenerConfig, cl, l, logPoller, opts, headTracker, gasEstimator)
		if err != nil {
			return nil, fmt.Errorf("failed to instantiate EvmTxm for chain with ID %s: %w", chainID, err)
		}
	}

	headBroadcaster.Subscribe(txm)

	var balanceMonitor monitor.BalanceMonitor
	if opts.ChainConfigs.RPCEnabled() && cfg.EVM().BalanceMonitor().Enabled() {
		balanceMonitor = monitor.NewBalanceMonitor(cl, opts.KeyStore, l)
		headBroadcaster.Subscribe(balanceMonitor)
	}

	var logBroadcaster log.Broadcaster
	if !opts.ChainConfigs.RPCEnabled() {
		logBroadcaster = &log.NullBroadcaster{ErrMsg: fmt.Sprintf("Ethereum is disabled for chain %d", chainID)}
	} else if !cfg.EVM().LogBroadcasterEnabled() {
		logBroadcaster = &log.NullBroadcaster{ErrMsg: fmt.Sprintf("LogBroadcaster disabled for chain %d", chainID)}
	} else if opts.GenLogBroadcaster == nil {
		logORM := log.NewORM(opts.DS, *chainID)
		logBroadcaster = log.NewBroadcaster(logORM, cl, cfg.EVM(), l, headSaver.LatestHeadFromDB, opts.MailMon)
	} else {
		logBroadcaster = opts.GenLogBroadcaster(chainID)
	}

	// AddDependent for this chain
	// log broadcaster will not start until dependent ready is called by a
	// subsequent routine (job spawner)
	logBroadcaster.AddDependents(1)

	headBroadcaster.Subscribe(logBroadcaster)

	// Construct the Tron TXM, will be nil if the chaintype is not tron
	var tronTxm *trontxm.TronTxm
	if cfg.EVM().ChainType() == chaintype.ChainTron {
		tronTxm, err = tron.ConstructTronTxm(l, cfg, nodes, opts.KeyStore)
		if err != nil {
			return nil, fmt.Errorf("failed to construct tron txm: %w", err)
		}
	}

	return &chain{
		id:              chainID,
		cfg:             cfg,
		client:          cl,
		txm:             txm,
		logger:          l,
		headBroadcaster: headBroadcaster,
		headTracker:     headTracker,
		logBroadcaster:  logBroadcaster,
		logPoller:       logPoller,
		balanceMonitor:  balanceMonitor,
		gasEstimator:    gasEstimator,

		// Extends with support for the Tron TXM
		tronTxm: tronTxm,
	}, nil
}

func (c *chain) Start(ctx context.Context) error {
	return c.StartOnce("Chain", func() error {
		c.logger.Debugf("Chain: starting with ID %s", c.ID().String())
		// Must ensure that EthClient is dialed first because subsequent
		// services may make eth calls on startup
		if err := c.client.Dial(ctx); err != nil {
			return fmt.Errorf("failed to dial ethclient: %w", err)
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

		if c.cfg.EVM().ChainType() == chaintype.ChainTron {
			c.gasEstimator.Start(ctx) // Still need gas estimator to be working for the OCR2 plugin
			c.tronTxm.Start(ctx)
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

		// Tron doesn't use the EVM TXM but still uses the gas estimator, we'll close it here
		if c.cfg.EVM().ChainType() == chaintype.ChainTron {
			merr = multierr.Combine(merr, c.gasEstimator.Close())
			merr = multierr.Combine(merr, c.tronTxm.Close())
		}

		c.logger.Debug("Chain: stopping client")
		c.client.Close()
		c.logger.Debug("Chain: stopped")
		return merr
	})
}

func (c *chain) Ready() (merr error) {
	merr = multierr.Combine(
		c.StateMachine.Ready(),
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
	report := map[string]error{c.Name(): c.Healthy()}
	services.CopyHealth(report, c.txm.HealthReport())
	services.CopyHealth(report, c.headBroadcaster.HealthReport())
	services.CopyHealth(report, c.headTracker.HealthReport())
	services.CopyHealth(report, c.logBroadcaster.HealthReport())

	if c.balanceMonitor != nil {
		services.CopyHealth(report, c.balanceMonitor.HealthReport())
	}

	return report
}

func (c *chain) Transact(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return errors.New("LOOPP not yet supported")
}

func (c *chain) Replay(ctx context.Context, fromBlock string, args map[string]any) error {
	block, err := strconv.ParseInt(fromBlock, 10, 64)
	if err != nil {
		return err
	}
	return c.logPoller.Replay(ctx, block)
}

func (c *chain) SendTx(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return c.Transact(ctx, from, to, amount, balanceCheck)
}

func (c *chain) LatestHead(_ context.Context) (types.Head, error) {
	latestChain := c.headTracker.LatestChain()
	if latestChain == nil {
		return types.Head{}, errors.New("latest chain not found")
	}

	return types.Head{
		Height:    strconv.FormatInt(latestChain.BlockNumber(), 10),
		Hash:      latestChain.Hash.Bytes(),
		Timestamp: uint64(latestChain.Timestamp.Unix()),
	}, nil
}

func (c *chain) GetChainStatus(ctx context.Context) (types.ChainStatus, error) {
	toml, err := c.cfg.EVM().TOMLString()
	if err != nil {
		return types.ChainStatus{}, err
	}
	return types.ChainStatus{
		ID:      c.ID().String(),
		Enabled: c.cfg.EVM().IsEnabled(),
		Config:  toml,
	}, nil
}

// TODO BCF-2602 statuses are static for non-evm chain and should be dynamic
func (c *chain) listNodeStatuses(start, end int) ([]types.NodeStatus, int, error) {
	nodes := c.cfg.Nodes()
	total := len(nodes)
	if start >= total {
		return nil, total, common.ErrOutOfRange
	}
	if end > total {
		end = total
	}
	stats := make([]types.NodeStatus, 0)

	states := c.Client().NodeStates()
	for _, n := range nodes[start:end] {
		var (
			nodeState string
		)
		toml, err := gotoml.Marshal(n)
		if err != nil {
			return nil, -1, err
		}
		if states == nil {
			nodeState = "Unknown"
		} else {
			// The node is in the DB and the chain is enabled but it's not running
			nodeState = "NotLoaded"
			s, exists := states[*n.Name]
			if exists {
				nodeState = s
			}
		}
		stats = append(stats, types.NodeStatus{
			ChainID: c.ID().String(),
			Name:    *n.Name,
			Config:  string(toml),
			State:   nodeState,
		})
	}
	return stats, total, nil
}

func (c *chain) ListNodeStatuses(ctx context.Context, pageSize int32, pageToken string) (stats []types.NodeStatus, nextPageToken string, total int, err error) {
	return common.ListNodeStatuses(int(pageSize), pageToken, c.listNodeStatuses)
}

func (c *chain) ID() *big.Int                           { return c.id }
func (c *chain) Client() client.Client                  { return c.client }
func (c *chain) Config() config.ChainScopedConfig       { return c.cfg }
func (c *chain) LogBroadcaster() log.Broadcaster        { return c.logBroadcaster }
func (c *chain) LogPoller() logpoller.LogPoller         { return c.logPoller }
func (c *chain) HeadBroadcaster() heads.Broadcaster     { return c.headBroadcaster }
func (c *chain) TxManager() txmgr.TxManager             { return c.txm }
func (c *chain) HeadTracker() heads.Tracker             { return c.headTracker }
func (c *chain) Logger() logger.Logger                  { return c.logger }
func (c *chain) BalanceMonitor() monitor.BalanceMonitor { return c.balanceMonitor }
func (c *chain) GasEstimator() gas.EvmFeeEstimator      { return c.gasEstimator }

// Add ChainTronSupport
func (c *chain) GetTronTXM() *trontxm.TronTxm { return c.tronTxm }
