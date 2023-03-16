package evm

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	v2 "github.com/smartcontractkit/chainlink/core/chains/evm/config/v2"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/config"
	cfgv2 "github.com/smartcontractkit/chainlink/core/config/v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

// ErrNoChains indicates that no EVM chains have been started
var ErrNoChains = errors.New("no EVM chains loaded")

var _ ChainSet = &chainSet{}

type ChainConfigUpdater func(*types.ChainCfg) error

//go:generate mockery --quiet --name ChainSet --output ./mocks/ --case=underscore
type ChainSet interface {
	services.ServiceCtx
	Get(id *big.Int) (Chain, error)

	Show(id utils.Big) (types.DBChain, error)

	Default() (Chain, error)
	Chains() []Chain
	ChainCount() int

	ORM() types.ORM

	Add(ctx context.Context, id utils.Big, config *types.ChainCfg) (types.DBChain, error)
	Remove(id utils.Big) error
	Index(offset, limit int) ([]types.DBChain, int, error)
	UpdateConfig(id *big.Int, updaters ...ChainConfigUpdater) error
	Configure(ctx context.Context, id utils.Big, enabled bool, config *types.ChainCfg) (types.DBChain, error)

	// GetNodes et al retrieves Nodes from the ORM and adds additional state info
	GetNodes(ctx context.Context, offset, limit int) (nodes []types.Node, count int, err error)
	GetNodesForChain(ctx context.Context, chainID utils.Big, offset, limit int) (nodes []types.Node, count int, err error)
	GetNodesByChainIDs(ctx context.Context, chainIDs []utils.Big) (nodes []types.Node, err error)

	CreateNode(ctx context.Context, data types.Node) (types.Node, error)
	DeleteNode(ctx context.Context, id int32) error
}

type chainSet struct {
	defaultID     *big.Int
	chains        map[string]*chain
	startedChains []Chain
	chainsMu      sync.RWMutex
	logger        logger.Logger
	opts          ChainSetOpts

	immutable bool // toml config is immutable
}

func (cll *chainSet) Start(ctx context.Context) error {
	if !cll.opts.Config.EVMEnabled() {
		cll.logger.Warn("EVM is disabled, no EVM-based chains will be started")
		return nil
	}
	if !cll.opts.Config.EVMRPCEnabled() {
		cll.logger.Warn("EVM RPC connections are disabled. Chainlink will not connect to any EVM RPC node.")
	}
	if cll.immutable {
		var ms services.MultiStart
		for _, c := range cll.Chains() {
			if err := ms.Start(ctx, c); err != nil {
				return errors.Wrapf(err, "failed to start chain %s", c.ID().String())
			}
			cll.startedChains = append(cll.startedChains, c)
		}
	} else {
		for _, c := range cll.Chains() {
			if err := c.Start(ctx); err != nil {
				id := c.ID().String()
				cll.logger.Criticalw(fmt.Sprintf("EVM: Chain with ID %s failed to start. You will need to fix this issue and restart the Chainlink node before any services that use this chain will work properly. Got error: %v", id, err), "evmChainID", id, "err", err)
				continue
			}
			cll.startedChains = append(cll.startedChains, c)
		}
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
func (cll *chainSet) Healthy() (err error) {
	for _, c := range cll.Chains() {
		err = multierr.Combine(err, c.Healthy())
	}
	return
}

func (cll *chainSet) Name() string {
	return cll.logger.Name()
}

func (cll *chainSet) HealthReport() map[string]error {
	return map[string]error{cll.Name(): cll.Healthy()}
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
	cll.chainsMu.RLock()
	defer cll.chainsMu.RUnlock()
	c, exists := cll.chains[id.String()]
	if exists {
		return c, nil
	}
	return nil, errors.Errorf("chain not found with id %v", id.String())
}

func (cll *chainSet) Show(id utils.Big) (types.DBChain, error) {
	return cll.opts.ORM.Chain(id)
}

func (cll *chainSet) Index(offset, limit int) ([]types.DBChain, int, error) {
	return cll.opts.ORM.Chains(offset, limit)
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

// Requires a lock on chainsMu
// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func (cll *chainSet) initializeChain(ctx context.Context, dbchain *types.DBChain) error {
	// preload nodes
	nodes, _, err := cll.opts.ORM.NodesForChain(dbchain.ID, 0, math.MaxInt)
	if err != nil {
		return err
	}

	cid := dbchain.ID.String()
	chain, err := newDBChain(ctx, *dbchain, nodes, cll.opts)
	if err != nil {
		return errors.Wrapf(err, "initializeChain: failed to instantiate chain %s", dbchain.ID.String())
	}
	if err = chain.Start(ctx); err != nil {
		return errors.Wrapf(err, "initializeChain: failed to start chain %s", dbchain.ID.String())
	}
	cll.startedChains = append(cll.startedChains, chain)
	cll.chains[cid] = chain
	return nil
}

func (cll *chainSet) Add(ctx context.Context, id utils.Big, config *types.ChainCfg) (types.DBChain, error) {
	if cll.immutable {
		return types.DBChain{}, cfgv2.ErrUnsupported
	}
	cll.chainsMu.Lock()
	defer cll.chainsMu.Unlock()

	cid := id.String()
	if _, exists := cll.chains[cid]; exists {
		return types.DBChain{}, errors.Errorf("chain already exists with id %s", id.String())
	}

	dbchain, err := cll.opts.ORM.CreateChain(id, config)
	if err != nil {
		return types.DBChain{}, err
	}
	return dbchain, cll.initializeChain(ctx, &dbchain)
}

func (cll *chainSet) Remove(id utils.Big) error {
	if cll.immutable {
		return cfgv2.ErrUnsupported
	}
	cll.chainsMu.Lock()
	defer cll.chainsMu.Unlock()

	if err := cll.opts.ORM.DeleteChain(id); err != nil {
		return err
	}

	cid := id.String()
	chain, exists := cll.chains[cid]
	if !exists {
		// If a chain was removed from the DB that wasn't loaded into the memory set we're done.
		return nil
	}
	delete(cll.chains, cid)
	return chain.Close()
}

func (cll *chainSet) Configure(ctx context.Context, id utils.Big, enabled bool, config *types.ChainCfg) (types.DBChain, error) {
	if cll.immutable {
		return types.DBChain{}, cfgv2.ErrUnsupported
	}
	cll.chainsMu.Lock()
	defer cll.chainsMu.Unlock()

	// Update configuration stored in the database
	dbchain, err := cll.opts.ORM.UpdateChain(id, enabled, config)
	if err != nil {
		return types.DBChain{}, err
	}

	cid := id.String()

	chain, exists := cll.chains[cid]

	switch {
	case exists && !enabled:
		// Chain was toggled to disabled
		delete(cll.chains, cid)
		return types.DBChain{}, chain.Close()
	case !exists && enabled:
		// Chain was toggled to enabled
		return dbchain, cll.initializeChain(ctx, &dbchain)
	case exists:
		// Exists in memory, no toggling: Update in-memory chain
		if chain.Config().Configure(*config); err != nil {
			return dbchain, err
		}
		// TODO: recreate ethClient etc if node set changed
		// https://app.shortcut.com/chainlinklabs/story/17044/chainset-should-update-chains-when-nodes-are-changed
	}

	return dbchain, nil
}

func (cll *chainSet) UpdateConfig(id *big.Int, updaters ...ChainConfigUpdater) error {
	if cll.immutable {
		return cfgv2.ErrUnsupported
	}
	bid := utils.NewBig(id)
	dbchain, err := cll.opts.ORM.Chain(*bid)
	if err != nil {
		return err
	}

	cll.chainsMu.RLock()
	chain, exists := cll.chains[id.String()]
	cll.chainsMu.RUnlock()
	if !exists {
		return errors.New("chain does not exist")
	}

	updatedConfig := chain.Config().PersistedConfig()
	for _, updater := range updaters {
		if err = updater(&updatedConfig); err != nil {
			cll.chainsMu.RUnlock()
			return err
		}
	}

	_, err = cll.opts.ORM.UpdateChain(*bid, dbchain.Enabled, &updatedConfig)
	if err == nil {
		chain.Config().Configure(updatedConfig)
	}

	return err
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

func (cll *chainSet) ORM() types.ORM {
	return cll.opts.ORM
}

func (cll *chainSet) GetNodes(ctx context.Context, offset, limit int) (nodes []types.Node, count int, err error) {
	nodes, count, err = cll.opts.ORM.Nodes(offset, limit, pg.WithParentCtx(ctx))
	if err != nil {
		err = errors.Wrap(err, "GetNodes failed to load nodes from DB")
		return
	}
	for i := range nodes {
		cll.addStateToNode(&nodes[i])
	}
	return
}

func (cll *chainSet) GetNodesForChain(ctx context.Context, chainID utils.Big, offset, limit int) (nodes []types.Node, count int, err error) {
	nodes, count, err = cll.opts.ORM.NodesForChain(chainID, offset, limit, pg.WithParentCtx(ctx))
	if err != nil {
		err = errors.Wrap(err, "GetNodesForChain failed to load nodes from DB")
		return
	}
	for i := range nodes {
		cll.addStateToNode(&nodes[i])
	}
	return
}

func (cll *chainSet) GetNodesByChainIDs(ctx context.Context, chainIDs []utils.Big) (nodes []types.Node, err error) {
	nodes, err = cll.opts.ORM.GetNodesByChainIDs(chainIDs, pg.WithParentCtx(ctx))
	if err != nil {
		err = errors.Wrap(err, "GetNodesForChain failed to load nodes from DB")
		return
	}
	for i := range nodes {
		cll.addStateToNode(&nodes[i])
	}
	return
}

func (cll *chainSet) CreateNode(ctx context.Context, data types.Node) (types.Node, error) {
	if cll.immutable {
		return types.Node{}, cfgv2.ErrUnsupported
	}
	return cll.opts.ORM.CreateNode(data, pg.WithParentCtx(ctx))
}

func (cll *chainSet) DeleteNode(ctx context.Context, id int32) error {
	if cll.immutable {
		return cfgv2.ErrUnsupported
	}
	return cll.opts.ORM.DeleteNode(id, pg.WithParentCtx(ctx))
}

func (cll *chainSet) addStateToNode(n *types.Node) {
	cll.chainsMu.RLock()
	chain, exists := cll.chains[n.EVMChainID.String()]
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

type ChainSetOpts struct {
	Config           config.GeneralConfig
	Logger           logger.Logger
	DB               *sqlx.DB
	KeyStore         keystore.Eth
	EventBroadcaster pg.EventBroadcaster
	ORM              types.ORM
	MailMon          *utils.MailboxMonitor

	// Gen-functions are useful for dependency injection by tests
	GenEthClient      func(*big.Int) client.Client
	GenLogBroadcaster func(*big.Int) log.Broadcaster
	GenLogPoller      func(*big.Int) logpoller.LogPoller
	GenHeadTracker    func(*big.Int, httypes.HeadBroadcaster) httypes.HeadTracker
	GenTxManager      func(*big.Int) txmgr.TxManager
}

// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func LoadChainSet(ctx context.Context, opts ChainSetOpts) (ChainSet, error) {
	if err := opts.check(); err != nil {
		return nil, err
	}
	if h, ok := opts.Config.(v2.HasEVMConfigs); ok {
		return NewTOMLChainSet(ctx, opts, h.EVMConfigs())
	}

	chains, err := opts.ORM.EnabledChains()
	if err != nil {
		return nil, errors.Wrap(err, "error loading chains")
	}
	nodesSlice, _, err := opts.ORM.Nodes(0, -1)
	if err != nil {
		return nil, errors.Wrap(err, "error loading nodes")
	}
	nodes := make(map[string][]types.Node)
	for _, n := range nodesSlice {
		id := n.EVMChainID.String()
		nodes[id] = append(nodes[id], n)
	}
	return NewDBChainSet(ctx, opts, chains, nodes)
}

// NewDBChainSet returns a new ChainSet from legacy configuration.
// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
func NewDBChainSet(ctx context.Context, opts ChainSetOpts, dbchains []types.DBChain, nodes map[string][]types.Node) (ChainSet, error) {
	if err := opts.check(); err != nil {
		return nil, err
	}
	opts.Logger = opts.Logger.Named("EVM")
	defaultChainID := opts.Config.DefaultChainID()
	if defaultChainID == nil && len(dbchains) >= 1 {
		defaultChainID = dbchains[0].ID.ToInt()
		if len(dbchains) > 1 {
			opts.Logger.Debugf("Multiple chains present but ETH_CHAIN_ID was not specified, falling back to default chain: %s", defaultChainID.String())
		}
	}
	var err error
	cll := newChainSet(opts)
	cll.defaultID = defaultChainID
	for i := range dbchains {
		cid := dbchains[i].ID.String()
		cll.logger.Infow(fmt.Sprintf("Loading chain %s", cid), "evmChainID", cid)
		chain, err2 := newDBChain(ctx, dbchains[i], nodes[cid], opts)
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

// NewTOMLChainSet returns a new ChainSet from TOML configuration.
func NewTOMLChainSet(ctx context.Context, opts ChainSetOpts, chains []*v2.EVMConfig) (ChainSet, error) {
	if err := opts.check(); err != nil {
		return nil, err
	}
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
	cll.immutable = true
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

	if tomlConfig, ok := opts.Config.(v2.HasEVMConfigs); ok {
		opts.ORM = chains.NewORMImmut[utils.Big, *types.ChainCfg, types.Node](tomlConfig.EVMConfigs())
	} else if opts.ORM == nil {
		// legacy config only
		opts.ORM = NewORM(opts.DB, opts.Logger, opts.Config)
	}
	return nil
}

func UpdateKeySpecificMaxGasPrice(addr common.Address, maxGasPriceWei *assets.Wei) ChainConfigUpdater {
	return func(config *types.ChainCfg) error {
		keyChainConfig, ok := config.KeySpecific[addr.Hex()]
		if !ok {
			keyChainConfig = types.ChainCfg{}
		}
		keyChainConfig.EvmMaxGasPriceWei = maxGasPriceWei
		if config.KeySpecific == nil {
			config.KeySpecific = map[string]types.ChainCfg{}
		}
		config.KeySpecific[addr.Hex()] = keyChainConfig
		return nil
	}
}
