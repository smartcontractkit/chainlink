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

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/config"
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

//go:generate mockery --name ChainSet --output ./mocks/ --case=underscore
type ChainSet interface {
	services.ServiceCtx
	Get(id *big.Int) (Chain, error)
	Show(id utils.Big) (types.DBChain, error)
	Add(ctx context.Context, id utils.Big, config *types.ChainCfg) (types.DBChain, error)
	Remove(id utils.Big) error
	Default() (Chain, error)
	Configure(ctx context.Context, id utils.Big, enabled bool, config *types.ChainCfg) (types.DBChain, error)
	UpdateConfig(id *big.Int, updaters ...ChainConfigUpdater) error
	Chains() []Chain
	Index(offset, limit int) ([]types.DBChain, int, error)
	ChainCount() int
	ORM() types.ORM
	// GetNode et al retrieves Nodes from the ORM and adds additional state info
	GetNode(ctx context.Context, id int32) (node types.Node, err error)
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
	orm           types.ORM
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
	for _, c := range cll.Chains() {
		if err := c.Start(ctx); err != nil {
			cll.logger.Criticalw(fmt.Sprintf("EVM: Chain with ID %s failed to start. You will need to fix this issue and restart the Chainlink node before any services that use this chain will work properly. Got error: %v", c.ID(), err), "evmChainID", c.ID(), "err", err)
			continue
		}
		cll.startedChains = append(cll.startedChains, c)
	}
	evmChainIDs := make([]*big.Int, len(cll.startedChains))
	for i, c := range cll.startedChains {
		evmChainIDs[i] = c.ID()
	}
	cll.logger.Infow(fmt.Sprintf("EVM: Started %d/%d chains, default chain ID is %s", len(cll.startedChains), len(cll.Chains()), cll.defaultID.String()), "startedEvmChainIDs", evmChainIDs)
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
func (cll *chainSet) Ready() (err error) {
	for _, c := range cll.Chains() {
		err = multierr.Combine(err, c.Ready())
	}
	return
}

func (cll *chainSet) Get(id *big.Int) (Chain, error) {
	if id == nil {
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
	return cll.orm.Chain(id)
}

func (cll *chainSet) Index(offset, limit int) ([]types.DBChain, int, error) {
	return cll.orm.Chains(offset, limit)
}

func (cll *chainSet) Default() (Chain, error) {
	cll.chainsMu.RLock()
	len := len(cll.chains)
	cll.chainsMu.RUnlock()
	if len == 0 {
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
func (cll *chainSet) initializeChain(ctx context.Context, dbchain *types.DBChain) error {
	// preload nodes
	nodes, _, err := cll.orm.NodesForChain(dbchain.ID, 0, math.MaxInt)
	if err != nil {
		return err
	}

	cid := dbchain.ID.String()
	chain, err := newChain(*dbchain, nodes, cll.opts)
	if err != nil {
		return errors.Wrapf(err, "initializeChain: failed to instantiate chain %s", dbchain.ID.String())
	}
	if err = chain.Start(ctx); err != nil {
		return errors.Wrapf(err, "initializeChain: failed to start chain %s", dbchain.ID.String())
	}
	cll.chains[cid] = chain
	return nil
}

func (cll *chainSet) Add(ctx context.Context, id utils.Big, config *types.ChainCfg) (types.DBChain, error) {
	cll.chainsMu.Lock()
	defer cll.chainsMu.Unlock()

	cid := id.String()
	if _, exists := cll.chains[cid]; exists {
		return types.DBChain{}, errors.Errorf("chain already exists with id %s", id.String())
	}

	dbchain, err := cll.orm.CreateChain(id, config)
	if err != nil {
		return types.DBChain{}, err
	}
	return dbchain, cll.initializeChain(ctx, &dbchain)
}

func (cll *chainSet) Remove(id utils.Big) error {
	cll.chainsMu.Lock()
	defer cll.chainsMu.Unlock()

	if err := cll.orm.DeleteChain(id); err != nil {
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
	cll.chainsMu.Lock()
	defer cll.chainsMu.Unlock()

	// Update configuration stored in the database
	dbchain, err := cll.orm.UpdateChain(id, enabled, config)
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
	bid := utils.NewBig(id)
	dbchain, err := cll.orm.Chain(*bid)
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

	_, err = cll.orm.UpdateChain(*bid, dbchain.Enabled, &updatedConfig)
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
	return cll.orm
}

func (cll *chainSet) GetNode(ctx context.Context, id int32) (n evmtypes.Node, err error) {
	n, err = cll.orm.Node(id, pg.WithParentCtx(ctx))
	if err != nil {
		err = errors.Wrap(err, "GetNode failed to load node from DB")
		return
	}
	cll.addStateToNode(&n)
	return
}

func (cll *chainSet) GetNodes(ctx context.Context, offset, limit int) (nodes []evmtypes.Node, count int, err error) {
	nodes, count, err = cll.orm.Nodes(offset, limit, pg.WithParentCtx(ctx))
	if err != nil {
		err = errors.Wrap(err, "GetNodes failed to load nodes from DB")
		return
	}
	for i := range nodes {
		cll.addStateToNode(&nodes[i])
	}
	return
}

func (cll *chainSet) GetNodesForChain(ctx context.Context, chainID utils.Big, offset, limit int) (nodes []evmtypes.Node, count int, err error) {
	nodes, count, err = cll.orm.NodesForChain(chainID, offset, limit, pg.WithParentCtx(ctx))
	if err != nil {
		err = errors.Wrap(err, "GetNodesForChain failed to load nodes from DB")
		return
	}
	for i := range nodes {
		cll.addStateToNode(&nodes[i])
	}
	return
}

func (cll *chainSet) GetNodesByChainIDs(ctx context.Context, chainIDs []utils.Big) (nodes []evmtypes.Node, err error) {
	nodes, err = cll.orm.GetNodesByChainIDs(chainIDs, pg.WithParentCtx(ctx))
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
	return cll.opts.ORM.CreateNode(data, pg.WithParentCtx(ctx))
}

func (cll *chainSet) DeleteNode(ctx context.Context, id int32) error {
	return cll.opts.ORM.DeleteNode(id, pg.WithParentCtx(ctx))
}

func (cll *chainSet) addStateToNode(n *evmtypes.Node) {
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
	state, exists := states[n.ID]
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

	// Gen-functions are useful for dependency injection by tests
	GenEthClient      func(types.DBChain) evmclient.Client
	GenLogBroadcaster func(types.DBChain) log.Broadcaster
	GenLogPoller      func(types.DBChain) *logpoller.LogPoller
	GenHeadTracker    func(types.DBChain, httypes.HeadBroadcaster) httypes.HeadTracker
	GenTxManager      func(types.DBChain) txmgr.TxManager
}

func LoadChainSet(opts ChainSetOpts) (ChainSet, error) {
	if err := checkOpts(&opts); err != nil {
		return nil, err
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
	return NewChainSet(opts, chains, nodes)
}

func NewChainSet(opts ChainSetOpts, dbchains []types.DBChain, nodes map[string][]types.Node) (ChainSet, error) {
	if err := checkOpts(&opts); err != nil {
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
	cll := &chainSet{defaultChainID, make(map[string]*chain), make([]Chain, 0), sync.RWMutex{}, opts.Logger.Named("ChainSet"), opts.ORM, opts}
	for i := range dbchains {
		cid := dbchains[i].ID.String()
		cll.logger.Infow(fmt.Sprintf("Loading chain %s", cid), "evmChainID", cid)
		chain, err2 := newChain(dbchains[i], nodes[cid], opts)
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

func checkOpts(opts *ChainSetOpts) error {
	if opts.Logger == nil {
		return errors.New("logger must be non-nil")
	}
	if opts.Config == nil {
		return errors.New("config must be non-nil")
	}
	if opts.ORM == nil {
		opts.ORM = NewORM(opts.DB, opts.Logger, opts.Config)
	}
	return nil
}

func UpdateKeySpecificMaxGasPrice(addr common.Address, maxGasPriceWei *big.Int) ChainConfigUpdater {
	return func(config *types.ChainCfg) error {
		keyChainConfig, ok := config.KeySpecific[addr.Hex()]
		if !ok {
			keyChainConfig = types.ChainCfg{}
		}
		keyChainConfig.EvmMaxGasPriceWei = (*utils.Big)(maxGasPriceWei)
		if config.KeySpecific == nil {
			config.KeySpecific = map[string]types.ChainCfg{}
		}
		config.KeySpecific[addr.Hex()] = keyChainConfig
		return nil
	}
}
