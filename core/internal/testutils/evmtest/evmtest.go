package evmtest

import (
	"database/sql"
	"math/big"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/stretchr/testify/mock"
	"go.uber.org/atomic"
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
	evmMocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
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
		opts.GenEthClient = func(c evmtypes.DBChain) evmclient.Client {
			return testopts.Client
		}
	}
	if testopts.LogBroadcaster != nil {
		opts.GenLogBroadcaster = func(c evmtypes.DBChain) log.Broadcaster {
			return testopts.LogBroadcaster
		}
	}
	if testopts.HeadTracker != nil {
		opts.GenHeadTracker = func(evmtypes.DBChain, httypes.HeadBroadcaster) httypes.HeadTracker {
			return testopts.HeadTracker
		}
	}
	if testopts.TxManager != nil {
		opts.GenTxManager = func(evmtypes.DBChain) txmgr.TxManager {
			return testopts.TxManager
		}

	}
	opts.Logger = logger.TestLogger(t)
	opts.Config = testopts.GeneralConfig

	chains := []evmtypes.DBChain{
		{
			ID:      *utils.NewBigI(0),
			Cfg:     &testopts.ChainCfg,
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

func MustInsertChain(t testing.TB, db *sqlx.DB, chain *evmtypes.DBChain) {
	query, args, e := db.BindNamed(`
INSERT INTO evm_chains (id, cfg, enabled, created_at, updated_at) VALUES (:id, :cfg, :enabled, NOW(), NOW()) RETURNING *;`, chain)
	require.NoError(t, e)
	err := db.Get(chain, query, args...)
	require.NoError(t, err)
}

type MockORM struct {
	mu     sync.RWMutex
	chains map[string]evmtypes.DBChain
	nodes  map[string][]evmtypes.Node
}

var _ evmtypes.ORM = &MockORM{}

func NewMockORM(chains []evmtypes.DBChain, nodes []evmtypes.Node) *MockORM {
	mo := &MockORM{
		chains: make(map[string]evmtypes.DBChain),
		nodes:  make(map[string][]evmtypes.Node),
	}
	mo.PutChains(chains...)
	mo.AddNodes(nodes...)
	return mo
}

func (mo *MockORM) PutChains(cs ...evmtypes.DBChain) {
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

func (mo *MockORM) EnabledChains(qopts ...pg.QOpt) ([]evmtypes.DBChain, error) {
	mo.mu.RLock()
	defer mo.mu.RUnlock()
	return maps.Values(mo.chains), nil
}

func (mo *MockORM) StoreString(chainID utils.Big, key, val string) error {
	panic("not implemented")
}

func (mo *MockORM) Clear(chainID utils.Big, key string) error {
	panic("not implemented")
}

func (mo *MockORM) Chain(id utils.Big, qopts ...pg.QOpt) (evmtypes.DBChain, error) {
	mo.mu.RLock()
	defer mo.mu.RUnlock()
	c, ok := mo.chains[id.String()]
	if !ok {
		return evmtypes.DBChain{}, sql.ErrNoRows
	}
	return c, nil
}

func (mo *MockORM) CreateChain(id utils.Big, config *evmtypes.ChainCfg, qopts ...pg.QOpt) (evmtypes.DBChain, error) {
	panic("not implemented")
}

func (mo *MockORM) UpdateChain(id utils.Big, enabled bool, config *evmtypes.ChainCfg, qopts ...pg.QOpt) (evmtypes.DBChain, error) {
	return evmtypes.DBChain{}, nil
}

func (mo *MockORM) DeleteChain(id utils.Big, qopts ...pg.QOpt) error {
	panic("not implemented")
}

func (mo *MockORM) Chains(offset int, limit int, qopts ...pg.QOpt) (chains []evmtypes.DBChain, count int, err error) {
	mo.mu.RLock()
	defer mo.mu.RUnlock()
	chains = maps.Values(mo.chains)
	count = len(chains)
	return
}

func (mo *MockORM) GetChainsByIDs(ids []utils.Big) (chains []evmtypes.DBChain, err error) {
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

func (mo *MockORM) CreateNode(data evmtypes.Node, qopts ...pg.QOpt) (n evmtypes.Node, err error) {
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
func (mo *MockORM) Nodes(offset int, limit int, qopts ...pg.QOpt) (nodes []evmtypes.Node, cnt int, err error) {
	mo.mu.RLock()
	defer mo.mu.RUnlock()
	for _, ns := range maps.Values(mo.nodes) {
		nodes = append(nodes, ns...)
	}
	cnt = len(nodes)
	return
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

// NodesForChain implements evmtypes.ORM
func (mo *MockORM) SetupNodes([]evmtypes.Node, []utils.Big) error {
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

func NewEthClientMock(t mock.TestingT) *evmMocks.Client {
	mockEth := new(evmMocks.Client)
	mockEth.Test(t)
	return mockEth
}

func NewEthClientMockWithDefaultChain(t testing.TB) *evmMocks.Client {
	c := NewEthClientMock(t)
	c.On("ChainID").Return(testutils.FixtureChainID).Maybe()
	return c
}

type MockEth struct {
	EthClient       *evmMocks.Client
	CheckFilterLogs func(int64, int64)

	subs             []*evmMocks.Subscription
	errChs           []chan error
	subscribeCalls   atomic.Int32
	unsubscribeCalls atomic.Int32
}

func (m *MockEth) AssertExpectations(t *testing.T) {
	m.EthClient.AssertExpectations(t)
	for _, sub := range m.subs {
		sub.AssertExpectations(t)
	}
}

func (m *MockEth) SubscribeCallCount() int32 {
	return m.subscribeCalls.Load()
}

func (m *MockEth) UnsubscribeCallCount() int32 {
	return m.unsubscribeCalls.Load()
}

func (m *MockEth) NewSub(t *testing.T) ethereum.Subscription {
	m.subscribeCalls.Inc()
	sub := new(evmMocks.Subscription)
	sub.Test(t)
	errCh := make(chan error)
	sub.On("Err").
		Return(func() <-chan error { return errCh })
	sub.On("Unsubscribe").
		Run(func(mock.Arguments) {
			m.unsubscribeCalls.Inc()
			close(errCh)
		}).Return().Maybe()
	m.subs = append(m.subs, sub)
	m.errChs = append(m.errChs, errCh)
	return sub
}

func (m *MockEth) SubsErr(err error) {
	for _, errCh := range m.errChs {
		errCh <- err
	}
}

type RawSub[T any] struct {
	ch  chan<- T
	err <-chan error
}

func NewRawSub[T any](ch chan<- T, err <-chan error) RawSub[T] {
	return RawSub[T]{ch: ch, err: err}
}

func (r *RawSub[T]) CloseCh() {
	close(r.ch)
}

func (r *RawSub[T]) TrySend(t T) {
	select {
	case <-r.err:
	case r.ch <- t:
	}
}
