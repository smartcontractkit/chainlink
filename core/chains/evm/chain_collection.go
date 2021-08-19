package evm

import (
	"math/big"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"go.uber.org/multierr"
	"gorm.io/gorm"
)

var ErrNoChains = errors.New("no chains loaded, are you running with EVM_DISABLED=true ?")

var _ ChainCollection = &chainCollection{}

//go:generate mockery --name ChainCollection --output ./mocks/ --case=underscore
type ChainCollection interface {
	service.Service
	Get(id *big.Int) (Chain, error)
	Default() (Chain, error)
	Chains() []Chain
	ChainCount() int
}

type chainCollection struct {
	defaultID *big.Int
	chains    map[string]*chain
}

func (cll *chainCollection) Start() (err error) {
	for _, c := range cll.Chains() {
		err = multierr.Combine(err, c.Start())
	}
	logger.Infof("ChainCollection: Started %d chains, default chain ID is %d", len(cll.chains), cll.defaultID)
	return
}
func (cll *chainCollection) Close() (err error) {
	for _, c := range cll.Chains() {
		err = multierr.Combine(err, c.Close())
	}
	return
}
func (cll *chainCollection) Healthy() (err error) {
	for _, c := range cll.Chains() {
		err = multierr.Combine(err, c.Healthy())
	}
	return
}
func (cll *chainCollection) Ready() (err error) {
	for _, c := range cll.Chains() {
		err = multierr.Combine(err, c.Ready())
	}
	return
}

func (cll *chainCollection) Get(id *big.Int) (Chain, error) {
	if id == nil {
		return nil, errors.New("cannot lookup chain by nil ID")
	}
	c, exists := cll.chains[id.String()]
	if exists {
		return c, nil
	}
	return nil, errors.Errorf("chain not found with id %d", id)
}

func (cll *chainCollection) Default() (Chain, error) {
	if len(cll.chains) == 0 {
		return nil, ErrNoChains
	}
	if cll.defaultID == nil {
		return nil, errors.New("no default chain ID specified")
	}
	return cll.Get(cll.defaultID)
}

func (cll *chainCollection) Chains() (c []Chain) {
	for _, chain := range cll.chains {
		c = append(c, chain)
	}
	return c
}

func (cll *chainCollection) ChainCount() int {
	return len(cll.chains)
}

type ChainCollectionOpts struct {
	Config           config.GeneralConfig
	Logger           *logger.Logger
	DB               *gorm.DB
	KeyStore         keystore.EthKeyStoreInterface
	AdvisoryLocker   postgres.AdvisoryLocker
	EventBroadcaster postgres.EventBroadcaster

	// Gen-functions are useful for dependency injection by tests
	GenEthClient              func(types.Chain) eth.Client
	GenLogBroadcaster         func(types.Chain) log.Broadcaster
	GenHeadTrackerBroadcaster func(types.Chain) (httypes.Tracker, httypes.HeadBroadcaster)
	GenTxManager              func(types.Chain) bulletprooftxmanager.TxManager
}

func LoadChainCollection(opts ChainCollectionOpts) (ChainCollection, error) {
	if opts.Config.EVMDisabled() {
		logger.Info("ChainCollection: EVM is disabled, no chains will be loaded")
		return &chainCollection{}, nil
	}
	var dbchains []types.Chain
	var nodes []types.Node
	// TODO: Can we use sqlx instead of gorm here
	// See: https://app.clubhouse.io/chainlinklabs/story/8781/remove-dependency-on-gorm
	if err := opts.DB.Find(&dbchains).Error; err != nil {
		return nil, err
	}
	if err := opts.DB.Find(&nodes).Error; err != nil {
		return nil, err
	}
	// HACK: gorm can't handle non-comparable foreign keys (utils.Big cannot be
	// used with ==), so preloading is not possible. Just manually assign here
	// instead
	for i, c := range dbchains {
		for _, n := range nodes {
			if n.EVMChainID.ToInt().Cmp(c.ID.ToInt()) == 0 {
				// Performance note: quadratic
				dbchains[i].Nodes = append(dbchains[i].Nodes, n)
			}
		}
	}
	return NewChainCollection(opts, dbchains)
}

func NewChainCollection(opts ChainCollectionOpts, dbchains []types.Chain) (ChainCollection, error) {
	if opts.Config == nil {
		panic("config must be non-nil")
	}
	var err error
	cll := &chainCollection{opts.Config.DefaultChainID(), make(map[string]*chain)}
	for i := range dbchains {
		chain, err2 := newChain(dbchains[i], opts)
		if err2 != nil {
			err = multierr.Combine(err, err2)
			continue
		}
		cll.chains[chain.ID().String()] = chain
	}
	return cll, err
}
