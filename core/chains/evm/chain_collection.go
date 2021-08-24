package evm

import (
	"math/big"

	"github.com/pkg/errors"
	"go.uber.org/multierr"
	"gorm.io/gorm"

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
	logger    *logger.Logger
}

func (cll *chainCollection) Start() (err error) {
	for _, c := range cll.Chains() {
		err = multierr.Combine(err, c.Start())
	}
	cll.logger.Infof("EVM: Started %d chains, default chain ID is %d", len(cll.chains), cll.defaultID)
	return
}
func (cll *chainCollection) Close() (err error) {
	cll.logger.Debug("EVM: stopping")
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
		cll.logger.Debugf("Chain ID not specified, using default: %s", cll.defaultID.String())
		return cll.Default()
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
	ORM              types.ORM

	// Gen-functions are useful for dependency injection by tests
	GenEthClient      func(types.Chain) eth.Client
	GenLogBroadcaster func(types.Chain) log.Broadcaster
	GenHeadTracker    func(types.Chain) httypes.Tracker
	GenTxManager      func(types.Chain) bulletprooftxmanager.TxManager
}

func LoadChainCollection(opts ChainCollectionOpts) (ChainCollection, error) {
	if err := checkOpts(&opts); err != nil {
		return nil, err
	}
	if opts.Config.EVMDisabled() {
		opts.Logger.Info("EVM is disabled, no chains will be loaded")
		return &chainCollection{logger: opts.Logger}, nil
	}
	if opts.ORM == nil {
		opts.ORM = NewORM(opts.DB)
	}
	dbchains, err := opts.ORM.LoadChains()
	if err != nil {
		return nil, errors.Wrap(err, "error loading chains")
	}
	return NewChainCollection(opts, dbchains)
}

func NewChainCollection(opts ChainCollectionOpts, dbchains []types.Chain) (ChainCollection, error) {
	if err := checkOpts(&opts); err != nil {
		return nil, err
	}
	opts.Logger.Infof("Creating ChainCollection with default chain id: %v and number of chains: %v", opts.Config.DefaultChainID(), len(dbchains))
	var err error
	cll := &chainCollection{opts.Config.DefaultChainID(), make(map[string]*chain), opts.Logger}
	for i := range dbchains {
		opts.Logger.Infof("EVM: Loading chain %s", dbchains[i].ID.String())
		chain, err2 := newChain(dbchains[i], opts)
		if err2 != nil {
			if errors.Cause(err2) == ErrNoPrimaryNode {
				opts.Logger.Warnf("EVM: No primary node found for chain %s; this chain will be ignored", dbchains[i].ID.String())
			} else {
				err = multierr.Combine(err, err2)
			}
			continue
		}
		cll.chains[chain.ID().String()] = chain
	}
	return cll, err
}

func checkOpts(opts *ChainCollectionOpts) error {
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
