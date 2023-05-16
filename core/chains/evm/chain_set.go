package evm

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/sqlx"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	v2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// ErrNoChains indicates that no EVM chains have been started
var ErrNoChains = errors.New("no EVM chains loaded")

var _ ChainSet = &chainSet{}

//go:generate mockery --quiet --name ChainSet --output ./mocks/ --case=underscore
type ChainSet interface {
	services.ServiceCtx
	chains.Chains
	chains.Nodes

	Get(id *big.Int) (Chain, error)

	Default() (Chain, error)
	Chains() []Chain
	ChainCount() int

	Configs() types.Configs

	SendTx(ctx context.Context, chainID, from, to string, amount *big.Int, balanceCheck bool) error
}

type chainSet struct {
	defaultID     *big.Int
	chains        map[string]*chain
	startedChains []Chain
	chainsMu      sync.RWMutex
	logger        logger.Logger
	opts          ChainSetOpts
}

func (cll *chainSet) Start(ctx context.Context) error {
	if !cll.opts.Config.EVMEnabled() {
		cll.logger.Warn("EVM is disabled, no EVM-based chains will be started")
		return nil
	}
	if !cll.opts.Config.EVMRPCEnabled() {
		cll.logger.Warn("EVM RPC connections are disabled. Chainlink will not connect to any EVM RPC node.")
	}
	var ms services.MultiStart
	for _, c := range cll.Chains() {
		if err := ms.Start(ctx, c); err != nil {
			return errors.Wrapf(err, "failed to start chain %q", c.ID())
		}
		cll.startedChains = append(cll.startedChains, c)
	}
	evmChainIDs := make([]*big.Int, len(cll.startedChains))
	for i, c := range cll.startedChains {
		evmChainIDs[i] = c.ID()
	}
	defChainID := "unspecified"
	if cll.defaultID != nil {
		defChainID = fmt.Sprintf("%q", cll.defaultID.String())
	}
	cll.logger.Infow(fmt.Sprintf("EVM: Started %d/%d chains, default chain ID is %s", len(cll.startedChains), len(cll.Chains()), defChainID), "startedEvmChainIDs", evmChainIDs)
	return nil
}

func (cll *chainSet) Close() (err error) {
	cll.logger.Debug("EVM: stopping")
	for _, c := range cll.startedChains {
		err = multierr.Combine(err, c.Close())
	}
	return
}

func (cll *chainSet) Name() string {
	return cll.logger.Name()
}

func (cll *chainSet) HealthReport() map[string]error {
	report := map[string]error{}
	for _, c := range cll.Chains() {
		maps.Copy(report, c.HealthReport())
	}
	return report
}

func (cll *chainSet) Ready() (err error) {
	for _, c := range cll.Chains() {
		err = multierr.Combine(err, c.Ready())
	}
	return
}

func (cll *chainSet) Get(id *big.Int) (Chain, error) {
	if id == nil {
		if cll.defaultID == nil {
			cll.logger.Debug("Chain ID not specified, and default is nil")
			return nil, errors.New("chain ID not specified, and default is nil")
		}
		cll.logger.Debugf("Chain ID not specified, using default: %s", cll.defaultID.String())
		return cll.Default()
	}
	return cll.get(id.String())
}

func (cll *chainSet) get(id string) (Chain, error) {
	cll.chainsMu.RLock()
	defer cll.chainsMu.RUnlock()
	c, exists := cll.chains[id]
	if exists {
		return c, nil
	}
	return nil, errors.Wrap(chains.ErrNotFound, fmt.Sprintf("failed to get chain with id %s", id))
}

func (cll *chainSet) ChainStatus(ctx context.Context, id string) (cfg relaytypes.ChainStatus, err error) {
	var cs []relaytypes.ChainStatus
	cs, _, err = cll.opts.Configs.Chains(0, -1, id)
	if err != nil {
		return
	}
	l := len(cs)
	if l == 0 {
		err = chains.ErrNotFound
		return
	}
	if l > 1 {
		err = fmt.Errorf("multiple chains found: %d", len(cs))
		return
	}
	cfg = cs[0]
	return
}

func (cll *chainSet) ChainStatuses(ctx context.Context, offset, limit int) ([]relaytypes.ChainStatus, int, error) {
	return cll.opts.Configs.Chains(offset, limit)
}

func (cll *chainSet) Default() (Chain, error) {
	cll.chainsMu.RLock()
	length := len(cll.chains)
	cll.chainsMu.RUnlock()
	if length == 0 {
		return nil, errors.Wrap(ErrNoChains, "cannot get default EVM chain; no EVM chains are available")
	}
	if cll.defaultID == nil {
		// This is an invariant violation; if any chains are available then a
		// default should _always_ have been set in the constructor
		return nil, errors.New("no default chain ID specified")
	}

	return cll.Get(cll.defaultID)
}

func (cll *chainSet) Chains() (c []Chain) {
	cll.chainsMu.RLock()
	defer cll.chainsMu.RUnlock()
	for _, chain := range cll.chains {
		c = append(c, chain)
	}
	return c
}

func (cll *chainSet) ChainCount() int {
	cll.chainsMu.RLock()
	defer cll.chainsMu.RUnlock()
	return len(cll.chains)
}

func (cll *chainSet) Configs() types.Configs {
	return cll.opts.Configs
}

func (cll *chainSet) NodeStatuses(ctx context.Context, offset, limit int, chainIDs ...string) (nodes []relaytypes.NodeStatus, count int, err error) {
	nodes, count, err = cll.opts.Configs.NodeStatusesPaged(offset, limit, chainIDs...)
	if err != nil {
		err = errors.Wrap(err, "GetNodesForChain failed to load nodes from DB")
		return
	}
	for i := range nodes {
		cll.addStateToNode(&nodes[i])
	}
	return
}

func (cll *chainSet) addStateToNode(n *relaytypes.NodeStatus) {
	cll.chainsMu.RLock()
	chain, exists := cll.chains[n.ChainID]
	cll.chainsMu.RUnlock()
	if !exists {
		// The EVM chain is disabled
		n.State = "Disabled"
		return
	}
	states := chain.Client().NodeStates()
	if states == nil {
		n.State = "Unknown"
		return
	}
	state, exists := states[n.Name]
	if exists {
		n.State = state
		return
	}
	// The node is in the DB and the chain is enabled but it's not running
	n.State = "NotLoaded"
}

func (cll *chainSet) SendTx(ctx context.Context, chainID, from, to string, amount *big.Int, balanceCheck bool) error {
	chain, err := cll.get(chainID)
	if err != nil {
		return err
	}

	return chain.SendTx(ctx, from, to, amount, balanceCheck)
}

type GeneralConfig interface {
	config.GeneralConfig
	v2.HasEVMConfigs
}

type ChainSetOpts struct {
	Config           GeneralConfig
	Logger           logger.Logger
	DB               *sqlx.DB
	KeyStore         keystore.Eth
	EventBroadcaster pg.EventBroadcaster
	Configs          types.Configs
	MailMon          *utils.MailboxMonitor
	GasEstimator     gas.EvmFeeEstimator

	// Gen-functions are useful for dependency injection by tests
	GenEthClient      func(*big.Int) client.Client
	GenLogBroadcaster func(*big.Int) log.Broadcaster
	GenLogPoller      func(*big.Int) logpoller.LogPoller
	GenHeadTracker    func(*big.Int, httypes.HeadBroadcaster) httypes.HeadTracker
	GenTxManager      func(*big.Int) txmgr.EvmTxManager
	GenGasEstimator   func(*big.Int) gas.EvmFeeEstimator
}

// NewTOMLChainSet returns a new ChainSet from TOML configuration.
func NewTOMLChainSet(ctx context.Context, opts ChainSetOpts) (ChainSet, error) {
	if err := opts.check(); err != nil {
		return nil, err
	}
	chains := opts.Config.EVMConfigs()
	var enabled []*v2.EVMConfig
	for i := range chains {
		if chains[i].IsEnabled() {
			enabled = append(enabled, chains[i])
		}
	}
	opts.Logger = opts.Logger.Named("EVM")
	defaultChainID := opts.Config.DefaultChainID()
	if defaultChainID == nil && len(enabled) >= 1 {
		defaultChainID = enabled[0].ChainID.ToInt()
		if len(enabled) > 1 {
			opts.Logger.Debugf("Multiple chains present, default chain: %s", defaultChainID.String())
		}
	}
	var err error
	cll := newChainSet(opts)
	cll.defaultID = defaultChainID
	for i := range enabled {
		cid := enabled[i].ChainID.String()
		cll.logger.Infow(fmt.Sprintf("Loading chain %s", cid), "evmChainID", cid)
		chain, err2 := newTOMLChain(ctx, enabled[i], opts)
		if err2 != nil {
			err = multierr.Combine(err, err2)
			continue
		}
		if _, exists := cll.chains[cid]; exists {
			return nil, errors.Errorf("duplicate chain with ID %s", cid)
		}
		cll.chains[cid] = chain
	}
	return cll, err
}

func newChainSet(opts ChainSetOpts) *chainSet {
	return &chainSet{
		chains:        make(map[string]*chain),
		startedChains: make([]Chain, 0),
		logger:        opts.Logger.Named("ChainSet"),
		opts:          opts,
	}
}

func (opts *ChainSetOpts) check() error {
	if opts.Logger == nil {
		return errors.New("logger must be non-nil")
	}
	if opts.Config == nil {
		return errors.New("config must be non-nil")
	}

	opts.Configs = chains.NewConfigs[utils.Big, types.Node](opts.Config.EVMConfigs())
	return nil
}
