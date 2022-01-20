package evm

import (
	"fmt"
	"math"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/sqlx"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/chains/evm/bulletprooftxmanager"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var ErrNoChains = errors.New("no chains loaded, are you running with EVM_DISABLED=true?")

var _ ChainSet = &chainSet{}

type ChainConfigUpdater func(*types.ChainCfg) error

//go:generate mockery --name ChainSet --output ./mocks/ --case=underscore
type ChainSet interface {
	services.Service
	Get(id *big.Int) (Chain, error)
	Add(id *big.Int, config types.ChainCfg) (types.Chain, error)
	Remove(id *big.Int) error
	Default() (Chain, error)
	Configure(id *big.Int, enabled bool, config types.ChainCfg) (types.Chain, error)
	UpdateConfig(id *big.Int, updaters ...ChainConfigUpdater) error
	Chains() []Chain
	ChainCount() int
	ORM() types.ORM
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

func (cll *chainSet) Start() error {
	if cll.opts.Config.EVMDisabled() {
		cll.logger.Warn("EVM is disabled, no EVM-based chains will be started")
		return nil
	}
	for _, c := range cll.Chains() {
		if err := c.Start(); err != nil {
			cll.logger.Errorw(fmt.Sprintf("EVM: Chain with ID %s failed to start. You will need to fix this issue and restart the Chainlink node before any services that use this chain will work properly. Got error: %v", c.ID(), err), "evmChainID", c.ID(), "err", err)
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

func (cll *chainSet) Default() (Chain, error) {
	cll.chainsMu.RLock()
	len := len(cll.chains)
	cll.chainsMu.RUnlock()
	if len == 0 {
		return nil, errors.Wrap(ErrNoChains, "cannot get default chain")
	}
	if cll.defaultID == nil {
		return nil, errors.New("no default chain ID specified")
	}

	return cll.Get(cll.defaultID)
}

// Requires a lock on chainsMu
func (cll *chainSet) initializeChain(dbchain *types.Chain) error {
	// preload nodes
	nodes, _, err := cll.orm.NodesForChain(dbchain.ID, 0, math.MaxInt)
	if err != nil {
		return err
	}
	dbchain.Nodes = nodes

	cid := dbchain.ID.String()
	chain, err := newChain(*dbchain, cll.opts)
	if err != nil {
		return errors.Wrapf(err, "initializeChain: failed to instantiate chain %s", dbchain.ID.String())
	}
	if err = chain.Start(); err != nil {
		return errors.Wrapf(err, "initializeChain: failed to start chain %s", dbchain.ID.String())
	}
	cll.chains[cid] = chain
	return nil
}

func (cll *chainSet) Add(id *big.Int, config types.ChainCfg) (types.Chain, error) {
	cll.chainsMu.Lock()
	defer cll.chainsMu.Unlock()

	cid := id.String()
	if _, exists := cll.chains[cid]; exists {
		return types.Chain{}, errors.Errorf("chain already exists with id %s", id.String())
	}

	bid := utils.NewBig(id)
	dbchain, err := cll.orm.CreateChain(*bid, config)
	if err != nil {
		return types.Chain{}, err
	}
	return dbchain, cll.initializeChain(&dbchain)
}

func (cll *chainSet) Remove(id *big.Int) error {
	cll.chainsMu.Lock()
	defer cll.chainsMu.Unlock()

	if err := cll.orm.DeleteChain(*utils.NewBig(id)); err != nil {
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

func (cll *chainSet) Configure(id *big.Int, enabled bool, config types.ChainCfg) (types.Chain, error) {
	cll.chainsMu.Lock()
	defer cll.chainsMu.Unlock()

	// Update configuration stored in the database
	bid := utils.NewBig(id)
	dbchain, err := cll.orm.UpdateChain(*bid, enabled, config)
	if err != nil {
		return types.Chain{}, err
	}

	cid := id.String()

	chain, exists := cll.chains[cid]

	switch {
	case exists && !enabled:
		// Chain was toggled to disabled
		delete(cll.chains, cid)
		return types.Chain{}, chain.Close()
	case !exists && enabled:
		// Chain was toggled to enabled
		return dbchain, cll.initializeChain(&dbchain)
	case exists:
		// Exists in memory, no toggling: Update in-memory chain
		if err = chain.Config().Configure(config); err != nil {
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

	_, err = cll.orm.UpdateChain(*bid, dbchain.Enabled, updatedConfig)
	if err == nil {
		return chain.Config().Configure(updatedConfig)
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

type ChainSetOpts struct {
	Config           config.GeneralConfig
	Logger           logger.Logger
	DB               *sqlx.DB
	KeyStore         keystore.Eth
	EventBroadcaster pg.EventBroadcaster
	ORM              types.ORM

	// Gen-functions are useful for dependency injection by tests
	GenEthClient      func(types.Chain) evmclient.Client
	GenLogBroadcaster func(types.Chain) log.Broadcaster
	GenHeadTracker    func(types.Chain) httypes.HeadTracker
	GenTxManager      func(types.Chain) bulletprooftxmanager.TxManager
}

func LoadChainSet(opts ChainSetOpts) (ChainSet, error) {
	if err := checkOpts(&opts); err != nil {
		return nil, err
	}
	dbchains, err := opts.ORM.EnabledChainsWithNodes()
	if err != nil {
		return nil, errors.Wrap(err, "error loading chains")
	}
	return NewChainSet(opts, dbchains)
}

func NewChainSet(opts ChainSetOpts, dbchains []types.Chain) (ChainSet, error) {
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
		chain, err2 := newChain(dbchains[i], opts)
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
		opts.ORM = NewORM(opts.DB)
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
