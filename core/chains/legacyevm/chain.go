package legacyevm

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	gotoml "github.com/pelletier/go-toml/v2"
	"go.uber.org/multierr"

	common "github.com/smartcontractkit/chainlink-common/pkg/chains"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/monitor"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

//go:generate mockery --quiet --name Chain --output ./mocks/ --case=underscore
type Chain interface {
	types.ChainService

	ID() *big.Int
	Client() evmclient.Client
	Config() evmconfig.ChainScopedConfig
	LogBroadcaster() log.Broadcaster
	HeadBroadcaster() httypes.HeadBroadcaster
	TxManager() txmgr.TxManager
	HeadTracker() httypes.HeadTracker
	Logger() logger.Logger
	BalanceMonitor() monitor.BalanceMonitor
	LogPoller() logpoller.LogPoller
	GasEstimator() gas.EvmFeeEstimator
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
//
//go:generate mockery --quiet --name LegacyChainContainer --output ./mocks/ --case=underscore
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
	cfg             *evmconfig.ChainScoped
	client          evmclient.Client
	txm             txmgr.TxManager
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
	ChainID *ubig.Big
}

func (e errChainDisabled) Error() string {
	return fmt.Sprintf("cannot create new chain with ID %s, the chain is disabled", e.ChainID.String())
}

// TODO BCF-2509 what is this and does it need the entire app config?
type AppConfig interface {
	config.AppConfig
	toml.HasEVMConfigs
}

type ChainRelayExtenderConfig struct {
	Logger   logger.Logger
	KeyStore keystore.Eth
	ChainOpts
}

func (c ChainRelayExtenderConfig) Validate() error {
	err := c.ChainOpts.Validate()
	if c.Logger == nil {
		err = errors.Join(err, errors.New("nil Logger"))
	}
	if c.KeyStore == nil {
		err = errors.Join(err, errors.New("nil Keystore"))
	}

	if err != nil {
		err = fmt.Errorf("invalid ChainRelayerExtenderConfig: %w", err)
	}
	return err
}

type ChainOpts struct {
	AppConfig AppConfig

	MailMon      *mailbox.Monitor
	GasEstimator gas.EvmFeeEstimator

	DS sqlutil.DataSource

	// TODO BCF-2513 remove test code from the API
	// Gen-functions are useful for dependency injection by tests
	GenEthClient      func(*big.Int) evmclient.Client
	GenLogBroadcaster func(*big.Int) log.Broadcaster
	GenLogPoller      func(*big.Int) logpoller.LogPoller
	GenHeadTracker    func(*big.Int, httypes.HeadBroadcaster) httypes.HeadTracker
	GenTxManager      func(*big.Int) txmgr.TxManager
	GenGasEstimator   func(*big.Int) gas.EvmFeeEstimator
}

func (o ChainOpts) Validate() error {
	var err error
	if o.AppConfig == nil {
		err = errors.Join(err, errors.New("nil AppConfig"))
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

func NewTOMLChain(ctx context.Context, chain *toml.EVMConfig, opts ChainRelayExtenderConfig) (Chain, error) {
	err := opts.Validate()
	if err != nil {
		return nil, err
	}
	chainID := chain.ChainID
	l := opts.Logger.With("evmChainID", chainID.String())
	if !chain.IsEnabled() {
		return nil, errChainDisabled{ChainID: chainID}
	}
	cfg := evmconfig.NewTOMLChainScopedConfig(chain, l)
	// note: per-chain validation is not necessary at this point since everything is checked earlier on boot.
	return newChain(ctx, cfg, chain.Nodes, opts)
}

func newChain(ctx context.Context, cfg *evmconfig.ChainScoped, nodes []*toml.Node, opts ChainRelayExtenderConfig) (*chain, error) {
	chainID := cfg.EVM().ChainID()
	l := opts.Logger
	var client evmclient.Client
	if !opts.AppConfig.EVMRPCEnabled() {
		client = evmclient.NewNullClient(chainID, l)
	} else if opts.GenEthClient == nil {
		client = evmclient.NewEvmClient(cfg.EVM().NodePool(), cfg.EVM(), cfg.EVM().NodePool().Errors(), l, chainID, nodes)
	} else {
		client = opts.GenEthClient(chainID)
	}

	headBroadcaster := headtracker.NewHeadBroadcaster(l)
	headSaver := headtracker.NullSaver
	var headTracker httypes.HeadTracker
	if !opts.AppConfig.EVMRPCEnabled() {
		headTracker = headtracker.NullTracker
	} else if opts.GenHeadTracker == nil {
		orm := headtracker.NewORM(*chainID, opts.DS)
		headSaver = headtracker.NewHeadSaver(l, orm, cfg.EVM(), cfg.EVM().HeadTracker())
		headTracker = headtracker.NewHeadTracker(l, client, cfg.EVM(), cfg.EVM().HeadTracker(), headBroadcaster, headSaver, opts.MailMon)
	} else {
		headTracker = opts.GenHeadTracker(chainID, headBroadcaster)
	}

	logPoller := logpoller.LogPollerDisabled
	if opts.AppConfig.Feature().LogPoller() {
		if opts.GenLogPoller != nil {
			logPoller = opts.GenLogPoller(chainID)
		} else {
			lpOpts := logpoller.Opts{
				PollPeriod:               cfg.EVM().LogPollInterval(),
				UseFinalityTag:           cfg.EVM().FinalityTagEnabled(),
				FinalityDepth:            int64(cfg.EVM().FinalityDepth()),
				BackfillBatchSize:        int64(cfg.EVM().LogBackfillBatchSize()),
				RpcBatchSize:             int64(cfg.EVM().RPCDefaultBatchSize()),
				KeepFinalizedBlocksDepth: int64(cfg.EVM().LogKeepBlocksDepth()),
				LogPrunePageSize:         int64(cfg.EVM().LogPrunePageSize()),
				BackupPollerBlockDelay:   int64(cfg.EVM().BackupLogPollerBlockDelay()),
			}
			logPoller = logpoller.NewLogPoller(logpoller.NewObservedORM(chainID, opts.DS, l), client, l, lpOpts)
		}
	}

	// note: gas estimator is started as a part of the txm
	txm, gasEstimator, err := newEvmTxm(opts.DS, cfg.EVM(), opts.AppConfig.EVMRPCEnabled(), opts.AppConfig.Database(), opts.AppConfig.Database().Listener(), client, l, logPoller, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate EvmTxm for chain with ID %s: %w", chainID.String(), err)
	}

	headBroadcaster.Subscribe(txm)

	// Highest seen head height is used as part of the start of LogBroadcaster backfill range
	highestSeenHead, err := headSaver.LatestHeadFromDB(ctx)
	if err != nil {
		return nil, err
	}

	var balanceMonitor monitor.BalanceMonitor
	if opts.AppConfig.EVMRPCEnabled() && cfg.EVM().BalanceMonitor().Enabled() {
		balanceMonitor = monitor.NewBalanceMonitor(client, opts.KeyStore, l)
		headBroadcaster.Subscribe(balanceMonitor)
	}

	var logBroadcaster log.Broadcaster
	if !opts.AppConfig.EVMRPCEnabled() {
		logBroadcaster = &log.NullBroadcaster{ErrMsg: fmt.Sprintf("Ethereum is disabled for chain %d", chainID)}
	} else if opts.GenLogBroadcaster == nil {
		logORM := log.NewORM(opts.DS, *chainID)
		logBroadcaster = log.NewBroadcaster(logORM, client, cfg.EVM(), l, highestSeenHead, opts.MailMon)
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
	return chains.ErrLOOPPUnsupported
}

func (c *chain) SendTx(ctx context.Context, from, to string, amount *big.Int, balanceCheck bool) error {
	return c.Transact(ctx, from, to, amount, balanceCheck)
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
			exists    bool
		)
		toml, err := gotoml.Marshal(n)
		if err != nil {
			return nil, -1, err
		}
		if states == nil {
			nodeState = "Unknown"
		} else {
			nodeState, exists = states[*n.Name]
			if !exists {
				// The node is in the DB and the chain is enabled but it's not running
				nodeState = "NotLoaded"
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

func (c *chain) ID() *big.Int                             { return c.id }
func (c *chain) Client() evmclient.Client                 { return c.client }
func (c *chain) Config() evmconfig.ChainScopedConfig      { return c.cfg }
func (c *chain) LogBroadcaster() log.Broadcaster          { return c.logBroadcaster }
func (c *chain) LogPoller() logpoller.LogPoller           { return c.logPoller }
func (c *chain) HeadBroadcaster() httypes.HeadBroadcaster { return c.headBroadcaster }
func (c *chain) TxManager() txmgr.TxManager               { return c.txm }
func (c *chain) HeadTracker() httypes.HeadTracker         { return c.headTracker }
func (c *chain) Logger() logger.Logger                    { return c.logger }
func (c *chain) BalanceMonitor() monitor.BalanceMonitor   { return c.balanceMonitor }
func (c *chain) GasEstimator() gas.EvmFeeEstimator        { return c.gasEstimator }
