package evmtest

import (
	"database/sql"
	"math/big"
	"math/rand"
	"sync"
	"testing"
	"time"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
	evmconfig "github.com/smartcontractkit/chainlink/core/chains/evm/config"
	httypes "github.com/smartcontractkit/chainlink/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type TestChainOpts struct {
	Client         evmclient.Client
	LogBroadcaster log.Broadcaster
	GeneralConfig  config.GeneralConfig
	ChainCfg       evmtypes.ChainCfg
	HeadTracker    httypes.HeadTracker
	DB             *sqlx.DB
	TxManager      txmgr.TxManager
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
		EventBroadcaster: pg.NewNullEventBroadcaster(),
	}
	if testopts.Client != nil {
		opts.GenEthClient = func(c evmtypes.Chain) evmclient.Client {
			return testopts.Client
		}
	}
	if testopts.LogBroadcaster != nil {
		opts.GenLogBroadcaster = func(c evmtypes.Chain) log.Broadcaster {
			return testopts.LogBroadcaster
		}
	}
	if testopts.HeadTracker != nil {
		opts.GenHeadTracker = func(evmtypes.Chain, httypes.HeadBroadcaster) httypes.HeadTracker {
			return testopts.HeadTracker
		}
	}
	if testopts.TxManager != nil {
		opts.GenTxManager = func(evmtypes.Chain) txmgr.TxManager {
			return testopts.TxManager
		}

	}
	opts.Logger = logger.TestLogger(t)
	opts.Config = testopts.GeneralConfig

	chains := []evmtypes.Chain{
		{
			ID:      *utils.NewBigI(0),
			Cfg:     testopts.ChainCfg,
			Enabled: true,
		},
	}
	nodes := map[string][]evmtypes.Node{
		"0": {{
			Name:       "evm-test-only-0",
			EVMChainID: *utils.NewBigI(0),
			WSURL:      null.StringFrom("ws://example.invalid"),
		}},
	}

	cc, err := evm.NewChainSet(opts, chains, nodes)
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
	mu     sync.RWMutex
	chains map[string]evmtypes.Chain
	nodes  map[string][]evmtypes.Node
}

var _ evmtypes.ORM = &MockORM{}

func NewMockORM(chains []evmtypes.Chain, nodes []evmtypes.Node) *MockORM {
	mo := &MockORM{
		chains: make(map[string]evmtypes.Chain),
		nodes:  make(map[string][]evmtypes.Node),
	}
	mo.PutChains(chains...)
	mo.AddNodes(nodes...)
	return mo
}

func (mo *MockORM) PutChains(cs ...evmtypes.Chain) {
	for _, c := range cs {
		mo.chains[c.ID.String()] = c
	}
}

func (mo *MockORM) AddNodes(ns ...evmtypes.Node) {
	for _, n := range ns {
		id := n.EVMChainID.String()
		mo.nodes[id] = append(mo.nodes[id], n)
	}
}

func (mo *MockORM) EnabledChainsWithNodes() ([]evmtypes.Chain, map[string][]evmtypes.Node, error) {
	mo.mu.RLock()
	defer mo.mu.RUnlock()
	return maps.Values(mo.chains), mo.nodes, nil
}

func (mo *MockORM) StoreString(chainID utils.Big, key, val string) error {
	panic("not implemented")
}

func (mo *MockORM) Clear(chainID utils.Big, key string) error {
	panic("not implemented")
}

func (mo *MockORM) Chain(id utils.Big, qopts ...pg.QOpt) (evmtypes.Chain, error) {
	mo.mu.RLock()
	defer mo.mu.RUnlock()
	c, ok := mo.chains[id.String()]
	if !ok {
		return evmtypes.Chain{}, sql.ErrNoRows
	}
	return c, nil
}

func (mo *MockORM) CreateChain(id utils.Big, config evmtypes.ChainCfg, qopts ...pg.QOpt) (evmtypes.Chain, error) {
	panic("not implemented")
}

func (mo *MockORM) UpdateChain(id utils.Big, enabled bool, config evmtypes.ChainCfg, qopts ...pg.QOpt) (evmtypes.Chain, error) {
	return evmtypes.Chain{}, nil
}

func (mo *MockORM) DeleteChain(id utils.Big, qopts ...pg.QOpt) error {
	panic("not implemented")
}

func (mo *MockORM) Chains(offset int, limit int, qopts ...pg.QOpt) (chains []evmtypes.Chain, count int, err error) {
	mo.mu.RLock()
	defer mo.mu.RUnlock()
	chains = maps.Values(mo.chains)
	count = len(chains)
	return
}

func (mo *MockORM) GetChainsByIDs(ids []utils.Big) (chains []evmtypes.Chain, err error) {
	mo.mu.RLock()
	defer mo.mu.RUnlock()
	for _, id := range ids {
		c, ok := mo.chains[id.String()]
		if ok {
			chains = append(chains, c)
		}
	}
	return
}

func (mo *MockORM) CreateNode(data evmtypes.NewNode, qopts ...pg.QOpt) (n evmtypes.Node, err error) {
	mo.mu.Lock()
	defer mo.mu.Unlock()
	n.ID = rand.Int31()
	n.Name = data.Name
	n.EVMChainID = data.EVMChainID
	n.WSURL = data.WSURL
	n.HTTPURL = data.HTTPURL
	n.SendOnly = data.SendOnly
	n.CreatedAt = time.Now()
	n.UpdatedAt = n.CreatedAt
	mo.AddNodes(n)
	return n, nil
}

func (mo *MockORM) DeleteNode(id int32, qopts ...pg.QOpt) error {
	mo.mu.Lock()
	defer mo.mu.Unlock()
	for chainID, ns := range mo.nodes {
		i := slices.IndexFunc(ns, func(n evmtypes.Node) bool {
			return n.ID == id
		})
		if i < 0 {
			continue
		}
		mo.nodes[chainID] = slices.Delete(ns, i, i)
		return nil
	}
	return sql.ErrNoRows
}

// Nodes implements evmtypes.ORM
func (mo *MockORM) Nodes(offset int, limit int, qopts ...pg.QOpt) ([]evmtypes.Node, int, error) {
	panic("not implemented")
}

// Node implements evmtypes.ORM
func (mo *MockORM) Node(id int32, qopts ...pg.QOpt) (evmtypes.Node, error) {
	panic("not implemented")
}

// GetNodesByChainIDs implements evmtypes.ORM
func (mo *MockORM) GetNodesByChainIDs(chainIDs []utils.Big, qopts ...pg.QOpt) (nodes []evmtypes.Node, err error) {
	panic("not implemented")
}

// NodesForChain implements evmtypes.ORM
func (mo *MockORM) NodesForChain(chainID utils.Big, offset int, limit int, qopts ...pg.QOpt) ([]evmtypes.Node, int, error) {
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
