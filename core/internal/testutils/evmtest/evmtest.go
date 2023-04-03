package evmtest

import (
	"database/sql"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	v2 "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/v2"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	evmmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/srvctest"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func NewChainScopedConfig(t testing.TB, cfg evm.GeneralConfig) evmconfig.ChainScopedConfig {
	var evmCfg *v2.EVMConfig
	if len(cfg.EVMConfigs()) > 0 {
		evmCfg = cfg.EVMConfigs()[0]
	} else {
		chainID := utils.NewBigI(0)
		evmCfg = &v2.EVMConfig{
			ChainID: chainID,
			Chain:   v2.Defaults(chainID),
		}
	}

	return v2.NewTOMLChainScopedConfig(cfg, evmCfg, logger.TestLogger(t))

}

type TestChainOpts struct {
	Client         evmclient.Client
	LogBroadcaster log.Broadcaster
	LogPoller      logpoller.LogPoller
	GeneralConfig  evm.GeneralConfig
	HeadTracker    httypes.HeadTracker
	DB             *sqlx.DB
	TxManager      txmgr.EvmTxManager
	KeyStore       keystore.Eth
	MailMon        *utils.MailboxMonitor
	GasEstimator   gas.EvmFeeEstimator
}

// NewChainSet returns a simple chain collection with one chain and
// allows to mock client/config on that chain
func NewChainSet(t testing.TB, testopts TestChainOpts) evm.ChainSet {
	opts := NewChainSetOpts(t, testopts)
	cc, err := evm.NewTOMLChainSet(testutils.Context(t), opts)
	require.NoError(t, err)
	return cc
}

// NewMockChainSetWithChain returns a mock chainset with one chain
func NewMockChainSetWithChain(t testing.TB, ch evm.Chain) *evmmocks.ChainSet {
	cc := evmmocks.NewChainSet(t)
	cc.On("Default").Return(ch, nil)
	return cc
}

func NewChainSetOpts(t testing.TB, testopts TestChainOpts) evm.ChainSetOpts {
	require.NotNil(t, testopts.KeyStore)
	opts := evm.ChainSetOpts{
		Config:           testopts.GeneralConfig,
		Logger:           logger.TestLogger(t),
		DB:               testopts.DB,
		KeyStore:         testopts.KeyStore,
		EventBroadcaster: pg.NewNullEventBroadcaster(),
		MailMon:          testopts.MailMon,
		GasEstimator:     testopts.GasEstimator,
	}
	opts.GenEthClient = func(*big.Int) evmclient.Client {
		if testopts.Client != nil {
			return testopts.Client
		}
		return evmclient.NewNullClient(testopts.GeneralConfig.DefaultChainID(), logger.TestLogger(t))
	}
	if testopts.LogBroadcaster != nil {
		opts.GenLogBroadcaster = func(*big.Int) log.Broadcaster {
			return testopts.LogBroadcaster
		}
	}
	if testopts.LogPoller != nil {
		opts.GenLogPoller = func(*big.Int) logpoller.LogPoller {
			return testopts.LogPoller
		}
	}
	if testopts.HeadTracker != nil {
		opts.GenHeadTracker = func(*big.Int, httypes.HeadBroadcaster) httypes.HeadTracker {
			return testopts.HeadTracker
		}
	}
	if testopts.TxManager != nil {
		opts.GenTxManager = func(*big.Int) txmgr.EvmTxManager {
			return testopts.TxManager
		}
	}
	if opts.MailMon == nil {
		opts.MailMon = srvctest.Start(t, utils.NewMailboxMonitor(t.Name()))
	}
	if testopts.GasEstimator != nil {
		opts.GenGasEstimator = func(*big.Int) gas.EvmFeeEstimator {
			return testopts.GasEstimator
		}
	}

	return opts
}

func MustGetDefaultChain(t testing.TB, cc evm.ChainSet) evm.Chain {
	chain, err := cc.Default()
	require.NoError(t, err)
	return chain
}

type MockORM struct {
	mu     sync.RWMutex
	chains map[string]chains.ChainConfig
	nodes  map[string][]evmtypes.Node
}

var _ evmtypes.Configs = &MockORM{}

func NewMockORM(cs []chains.ChainConfig, nodes []evmtypes.Node) *MockORM {
	mo := &MockORM{
		chains: make(map[string]chains.ChainConfig),
		nodes:  make(map[string][]evmtypes.Node),
	}
	mo.PutChains(cs...)
	mo.AddNodes(nodes...)
	return mo
}

func (mo *MockORM) PutChains(cs ...chains.ChainConfig) {
	mo.mu.Lock()
	defer mo.mu.Unlock()
	for _, c := range cs {
		mo.chains[c.ID] = c
	}
}

func (mo *MockORM) AddNodes(ns ...evmtypes.Node) {
	for _, n := range ns {
		id := n.EVMChainID.String()
		mo.nodes[id] = append(mo.nodes[id], n)
	}
}

func (mo *MockORM) Chains(offset int, limit int, ids ...utils.Big) (cs []chains.ChainConfig, count int, err error) {
	mo.mu.RLock()
	defer mo.mu.RUnlock()
	if len(ids) == 0 {
		cs = maps.Values(mo.chains)
		count = len(cs)
		return
	}
	for _, id := range ids {
		c, ok := mo.chains[id.String()]
		if ok {
			cs = append(cs, c)
		}
	}
	count = len(cs)
	return
}

// Nodes implements evmtypes.Configs
func (mo *MockORM) Nodes(offset int, limit int) (nodes []evmtypes.Node, cnt int, err error) {
	mo.mu.RLock()
	defer mo.mu.RUnlock()
	for _, ns := range maps.Values(mo.nodes) {
		nodes = append(nodes, ns...)
	}
	cnt = len(nodes)
	return
}

func (mo *MockORM) NodeNamed(name string) (evmtypes.Node, error) {
	mo.mu.RLock()
	defer mo.mu.RUnlock()
	for _, ns := range maps.Values(mo.nodes) {
		for _, n := range ns {
			if n.Name == name {
				return n, nil
			}
		}
	}
	return evmtypes.Node{}, sql.ErrNoRows
}

// GetNodesByChainIDs implements evmtypes.Configs
func (mo *MockORM) GetNodesByChainIDs(chainIDs []utils.Big) (nodes []evmtypes.Node, err error) {
	ids := map[string]struct{}{}
	for _, chainID := range chainIDs {
		ids[chainID.String()] = struct{}{}
	}
	mo.mu.RLock()
	defer mo.mu.RUnlock()
	for _, ns := range maps.Values(mo.nodes) {
		for _, n := range ns {
			if _, ok := ids[n.EVMChainID.String()]; ok {
				nodes = append(nodes, n)
			}
		}
	}
	return
}

// NodesForChain implements evmtypes.Configs
func (mo *MockORM) NodesForChain(chainID utils.Big, offset int, limit int) ([]evmtypes.Node, int, error) {
	panic("not implemented")
}

func NewEthClientMock(t *testing.T) *evmclimocks.Client {
	return evmclimocks.NewClient(t)
}

func NewEthClientMockWithDefaultChain(t *testing.T) *evmclimocks.Client {
	c := NewEthClientMock(t)
	c.On("ChainID").Return(testutils.FixtureChainID).Maybe()
	return c
}

type MockEth struct {
	EthClient       *evmclimocks.Client
	CheckFilterLogs func(int64, int64)

	subsMu           sync.RWMutex
	subs             []*evmclimocks.Subscription
	errChs           []chan error
	subscribeCalls   atomic.Int32
	unsubscribeCalls atomic.Int32
}

func (m *MockEth) SubscribeCallCount() int32 {
	return m.subscribeCalls.Load()
}

func (m *MockEth) UnsubscribeCallCount() int32 {
	return m.unsubscribeCalls.Load()
}

func (m *MockEth) NewSub(t *testing.T) ethereum.Subscription {
	m.subscribeCalls.Add(1)
	sub := evmclimocks.NewSubscription(t)
	errCh := make(chan error)
	sub.On("Err").
		Return(func() <-chan error { return errCh }).Maybe()
	sub.On("Unsubscribe").
		Run(func(mock.Arguments) {
			m.unsubscribeCalls.Add(1)
			close(errCh)
		}).Return().Maybe()
	m.subsMu.Lock()
	m.subs = append(m.subs, sub)
	m.errChs = append(m.errChs, errCh)
	m.subsMu.Unlock()
	return sub
}

func (m *MockEth) SubsErr(err error) {
	m.subsMu.Lock()
	defer m.subsMu.Unlock()
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
