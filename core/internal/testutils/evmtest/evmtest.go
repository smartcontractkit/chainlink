package evmtest

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	evmmocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/postgres"
	"github.com/smartcontractkit/chainlink/core/store/config"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func verifyMatchingChainIDs(t testing.TB, n *big.Int, m *big.Int) {
	require.Equal(t, n.Cmp(m), 0, "expected chain IDs to match")
}

type TestChainOpts struct {
	Client         eth.Client
	LogBroadcaster log.Broadcaster
	GeneralConfig  config.GeneralConfig
	ChainCfg       evmtypes.ChainCfg
	HeadTracker    httypes.Tracker
	DB             *gorm.DB
	TxManager      bulletprooftxmanager.TxManager
	KeyStore       keystore.EthKeyStoreInterface
	Logger         *logger.Logger
}

func NewChainScopedConfig(t testing.TB, cfg config.GeneralConfig) evmconfig.ChainScopedConfig {
	return evmconfig.NewChainScopedConfig(nil, logger.Default, cfg, evmtypes.Chain{ID: *utils.NewBigI(0)})
}

// NewChainCollection returns a simple chain collection with one chain and
// allows to mock client/config on that chain
func NewChainCollection(t testing.TB, testopts TestChainOpts) evm.ChainCollection {
	opts := evm.ChainCollectionOpts{
		Config:           testopts.GeneralConfig,
		DB:               testopts.DB,
		KeyStore:         testopts.KeyStore,
		EventBroadcaster: postgres.NewNullEventBroadcaster(),
		AdvisoryLocker:   postgres.NewNullAdvisoryLocker(),
	}
	if testopts.Client != nil {
		opts.GenEthClient = func(c evmtypes.Chain) eth.Client {
			return testopts.Client
		}
	}
	if testopts.LogBroadcaster != nil {
		opts.GenLogBroadcaster = func(c evmtypes.Chain) log.Broadcaster {
			return testopts.LogBroadcaster
		}
	}
	if testopts.HeadTracker != nil {
		opts.GenHeadTracker = func(evmtypes.Chain) httypes.Tracker {
			return testopts.HeadTracker
		}
	}
	if testopts.TxManager != nil {
		opts.GenTxManager = func(evmtypes.Chain) bulletprooftxmanager.TxManager {
			return testopts.TxManager
		}

	}
	if testopts.Logger != nil {
		opts.Logger = testopts.Logger
	} else if opts.Config != nil {
		opts.Logger = opts.Config.CreateProductionLogger()
	} else {
		opts.Logger = logger.Default
	}

	opts.Config = testopts.GeneralConfig

	chains := []evmtypes.Chain{
		{
			ID:  *utils.NewBigI(0),
			Cfg: testopts.ChainCfg,
			Nodes: []evmtypes.Node{{
				Name:       "evm-test-only-0",
				EVMChainID: *utils.NewBigI(0),
				WSURL:      "ws://example.invalid",
			}},
			Enabled: true,
		},
	}

	cc, err := evm.NewChainCollection(opts, chains)
	require.NoError(t, err)
	return cc
}

func MustGetDefaultChain(t testing.TB, cc evm.ChainCollection) evm.Chain {
	chain, err := cc.Default()
	require.NoError(t, err)
	return chain
}

func MustInsertChainWithNode(t testing.TB, db *gorm.DB, chain evmtypes.Chain) evmtypes.Chain {
	err := db.Create(&chain).Error
	require.NoError(t, err)
	return chain
}

type MockORM struct {
	chains []evmtypes.Chain
}

var _ evmtypes.ORM = &MockORM{}

func NewMockORM(chains []evmtypes.Chain) *MockORM {
	mo := &MockORM{
		chains: chains,
	}
	return mo
}

func (mo *MockORM) LoadChains() ([]evmtypes.Chain, error) {
	return mo.chains, nil
}

func (mo *MockORM) StoreString(chainID *big.Int, key, val string) error {
	return nil
}

func (mo *MockORM) Clear(chainID *big.Int, key string) error {
	return nil
}

func ChainEthMainnet() evmtypes.Chain      { return evmtypes.Chain{ID: *utils.NewBigI(1)} }
func ChainOptimismMainnet() evmtypes.Chain { return evmtypes.Chain{ID: *utils.NewBigI(10)} }
func ChainOptimismKovan() evmtypes.Chain   { return evmtypes.Chain{ID: *utils.NewBigI(69)} }
func ChainArbitrumMainnet() evmtypes.Chain { return evmtypes.Chain{ID: *utils.NewBigI(42161)} }
func ChainArbitrumRinkeby() evmtypes.Chain { return evmtypes.Chain{ID: *utils.NewBigI(421611)} }

func NewMockChainCollection(t testing.TB) *evmmocks.ChainCollection {
	cc := new(evmmocks.ChainCollection)
	cc.Test(t)
	return cc
}
