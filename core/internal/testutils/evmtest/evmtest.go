package evmtest

import (
	"math/big"
	"testing"

	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type TestChainOpts struct {
	Client         eth.Client
	LogBroadcaster log.Broadcaster
	GeneralConfig  config.GeneralConfig
	ChainCfg       evmtypes.ChainCfg
	HeadTracker    services.Tracker
	DB             *sqlx.DB
	TxManager      bulletprooftxmanager.TxManager
	KeyStore       keystore.Eth
}

func NewChainScopedConfig(t testing.TB, cfg config.GeneralConfig) evmconfig.ChainScopedConfig {
	return evmconfig.NewChainScopedConfig(big.NewInt(0), evmtypes.ChainCfg{},
		nil, logger.TestLogger(t), cfg)
}

// NewChainSet returns a simple chain collection with one chain and
// allows to mock client/config on that chain
func NewChainSet(t testing.TB, testopts TestChainOpts) evm.ChainSet {
	opts := evm.ChainSetOpts{
		Config:           testopts.GeneralConfig,
		DB:               testopts.DB,
		KeyStore:         testopts.KeyStore,
		EventBroadcaster: services.NewNullEventBroadcaster(),
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
		opts.GenHeadTracker = func(evmtypes.Chain) services.Tracker {
			return testopts.HeadTracker
		}
	}
	if testopts.TxManager != nil {
		opts.GenTxManager = func(evmtypes.Chain) bulletprooftxmanager.TxManager {
			return testopts.TxManager
		}

	}
	opts.Logger = logger.TestLogger(t)
	opts.Config = testopts.GeneralConfig

	chains := []evmtypes.Chain{
		{
			ID:  *utils.NewBigI(0),
			Cfg: testopts.ChainCfg,
			Nodes: []evmtypes.Node{{
				Name:       "evm-test-only-0",
				EVMChainID: *utils.NewBigI(0),
				WSURL:      null.StringFrom("ws://example.invalid"),
			}},
			Enabled: true,
		},
	}

	cc, err := evm.NewChainSet(opts, chains)
	require.NoError(t, err)
	return cc
}

func MustGetDefaultChain(t testing.TB, cc evm.ChainSet) evm.Chain {
	chain, err := cc.Default()
	require.NoError(t, err)
	return chain
}

func MustInsertChain(t testing.TB, db *sqlx.DB, chain *evmtypes.Chain) {
	query, args, e := db.BindNamed(`
INSERT INTO evm_chains (id, cfg, enabled, created_at, updated_at) VALUES (:id, :cfg, :enabled, NOW(), NOW()) RETURNING *;`, chain)
	require.NoError(t, e)
	err := db.Get(chain, query, args...)
	require.NoError(t, err)
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

func (mo *MockORM) EnabledChainsWithNodes() ([]evmtypes.Chain, error) {
	return mo.chains, nil
}

func (mo *MockORM) StoreString(chainID *big.Int, key, val string) error {
	return nil
}

func (mo *MockORM) Clear(chainID *big.Int, key string) error {
	return nil
}

func (mo *MockORM) Chain(id utils.Big) (evmtypes.Chain, error) {
	panic("not implemented")
}

func (mo *MockORM) CreateChain(id utils.Big, config evmtypes.ChainCfg) (evmtypes.Chain, error) {
	panic("not implemented")
}

func (mo *MockORM) UpdateChain(id utils.Big, enabled bool, config evmtypes.ChainCfg) (evmtypes.Chain, error) {
	return evmtypes.Chain{}, nil
}

func (mo *MockORM) DeleteChain(id utils.Big) error {
	panic("not implemented")
}

func (mo *MockORM) Chains(offset int, limit int) ([]evmtypes.Chain, int, error) {
	panic("not implemented")
}

func (mo *MockORM) GetChainsByIDs(ids []utils.Big) (chains []evmtypes.Chain, err error) {
	panic("not implemented")
}

func (mo *MockORM) CreateNode(data evmtypes.NewNode) (evmtypes.Node, error) {
	panic("not implemented")
}

func (mo *MockORM) DeleteNode(id int64) error {
	panic("not implemented")
}

func (mo *MockORM) Nodes(offset int, limit int) ([]evmtypes.Node, int, error) {
	panic("not implemented")
}

func (mo *MockORM) Node(id int32) (evmtypes.Node, error) {
	panic("not implemented")
}

func (mo *MockORM) GetNodesByChainIDs(chainIDs []utils.Big) (nodes []evmtypes.Node, err error) {
	panic("not implemented")
}

func (mo *MockORM) NodesForChain(chainID utils.Big, offset int, limit int) ([]evmtypes.Node, int, error) {
	panic("not implemented")
}

func ChainEthMainnet(t *testing.T) evmconfig.ChainScopedConfig      { return scopedConfig(t, 1) }
func ChainOptimismMainnet(t *testing.T) evmconfig.ChainScopedConfig { return scopedConfig(t, 10) }
func ChainOptimismKovan(t *testing.T) evmconfig.ChainScopedConfig   { return scopedConfig(t, 69) }
func ChainArbitrumMainnet(t *testing.T) evmconfig.ChainScopedConfig { return scopedConfig(t, 42161) }
func ChainArbitrumRinkeby(t *testing.T) evmconfig.ChainScopedConfig { return scopedConfig(t, 421611) }

func scopedConfig(t *testing.T, chainID int64) evmconfig.ChainScopedConfig {
	return evmconfig.NewChainScopedConfig(big.NewInt(chainID), evmtypes.ChainCfg{}, nil,
		logger.TestLogger(t), configtest.NewTestGeneralConfig(t))
}
