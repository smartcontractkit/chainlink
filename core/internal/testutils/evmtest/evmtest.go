package evmtest

import (
	"math/big"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/pelletier/go-toml/v2"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slices"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-relay/pkg/types"

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

type TestConfigs struct {
	mu sync.RWMutex
	v2.EVMConfigs
}

var _ evmtypes.Configs = &TestConfigs{}

func NewTestConfigs(cs ...*v2.EVMConfig) *TestConfigs {
	return &TestConfigs{EVMConfigs: v2.EVMConfigs(cs)}
}

func (mo *TestConfigs) PutChains(cs ...v2.EVMConfig) {
	mo.mu.Lock()
	defer mo.mu.Unlock()
chains:
	for i := range cs {
		id := cs[i].ChainID.String()
		for j, c2 := range mo.EVMConfigs {
			if c2.ChainID.String() == id {
				mo.EVMConfigs[j] = &cs[i] // replace
				continue chains
			}
		}
		mo.EVMConfigs = append(mo.EVMConfigs, &cs[i])
	}
}

func (mo *TestConfigs) Chains(offset int, limit int, ids ...string) (cs []types.ChainStatus, count int, err error) {
	mo.mu.RLock()
	defer mo.mu.RUnlock()
	if len(ids) == 0 {
		for _, c := range mo.EVMConfigs {
			c2 := types.ChainStatus{
				ID:      c.ChainID.String(),
				Enabled: c.IsEnabled(),
			}
			c2.Config, err = c.TOMLString()
			if err != nil {
				return
			}
			cs = append(cs, c2)
		}
		count = len(cs)
		return
	}
	for i := range mo.EVMConfigs {
		c := mo.EVMConfigs[i]
		chainID := c.ChainID.String()
		if !slices.Contains(ids, chainID) {
			continue
		}
		c2 := types.ChainStatus{
			ID:      chainID,
			Enabled: c.IsEnabled(),
		}
		c2.Config, err = c.TOMLString()
		if err != nil {
			return
		}
		cs = append(cs, c2)
	}
	count = len(cs)
	return
}

// Nodes implements evmtypes.Configs
func (mo *TestConfigs) Nodes(chainID utils.Big) (nodes []evmtypes.Node, err error) {
	mo.mu.RLock()
	defer mo.mu.RUnlock()

	for i := range mo.EVMConfigs {
		c := mo.EVMConfigs[i]
		if chainID.Cmp(c.ChainID) == 0 {
			for _, n := range c.Nodes {
				nodes = append(nodes, legacyNode(n, c.ChainID))
			}
		}
	}
	err = chains.ErrNotFound
	return
}

func (mo *TestConfigs) Node(name string) (evmtypes.Node, error) {
	mo.mu.RLock()
	defer mo.mu.RUnlock()

	for i := range mo.EVMConfigs {
		c := mo.EVMConfigs[i]
		for _, n := range c.Nodes {
			if *n.Name == name {
				return legacyNode(n, c.ChainID), nil
			}
		}
	}
	return evmtypes.Node{}, chains.ErrNotFound
}

func (mo *TestConfigs) NodeStatusesPaged(offset int, limit int, chainIDs ...string) (nodes []types.NodeStatus, cnt int, err error) {
	mo.mu.RLock()
	defer mo.mu.RUnlock()

	for i := range mo.EVMConfigs {
		c := mo.EVMConfigs[i]
		id := c.ChainID.String()
		if !slices.Contains(chainIDs, id) {
			continue
		}
		for _, n := range c.Nodes {
			var n2 types.NodeStatus
			n2, err = nodeStatus(n, id)
			if err != nil {
				return
			}
			nodes = append(nodes, n2)
		}
	}
	cnt = len(nodes)
	return
}

func legacyNode(n *v2.Node, chainID *utils.Big) (v2 evmtypes.Node) {
	v2.Name = *n.Name
	v2.EVMChainID = *chainID
	if n.HTTPURL != nil {
		v2.HTTPURL = null.StringFrom(n.HTTPURL.String())
	}
	if n.WSURL != nil {
		v2.WSURL = null.StringFrom(n.WSURL.String())
	}
	if n.SendOnly != nil {
		v2.SendOnly = *n.SendOnly
	}
	return
}

func nodeStatus(n *v2.Node, chainID string) (types.NodeStatus, error) {
	var s types.NodeStatus
	s.ChainID = chainID
	s.Name = *n.Name
	b, err := toml.Marshal(n)
	if err != nil {
		return types.NodeStatus{}, err
	}
	s.Config = string(b)
	return s, nil
}

func NewEthClientMock(t *testing.T) *evmclimocks.Client {
	return evmclimocks.NewClient(t)
}

func NewEthClientMockWithDefaultChain(t *testing.T) *evmclimocks.Client {
	c := NewEthClientMock(t)
	c.On("ConfiguredChainID").Return(testutils.FixtureChainID).Maybe()
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
