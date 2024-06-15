package evmtest

import (
	"fmt"
	"math/big"
	"slices"
	"sync"
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox/mailboxtest"

	"github.com/smartcontractkit/chainlink/v2/core/chains"
	evmclient "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	evmconfig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	evmtoml "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/gas"
	httypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/log"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

func NewChainScopedConfig(t testing.TB, cfg legacyevm.AppConfig) evmconfig.ChainScopedConfig {
	var evmCfg *evmtoml.EVMConfig
	if len(cfg.EVMConfigs()) > 0 {
		evmCfg = cfg.EVMConfigs()[0]
	} else {
		var chainID = (*ubig.Big)(testutils.FixtureChainID)
		evmCfg = &evmtoml.EVMConfig{
			ChainID: chainID,
			Chain:   evmtoml.Defaults(chainID),
		}
	}

	return evmconfig.NewTOMLChainScopedConfig(evmCfg, logger.TestLogger(t))
}

type TestChainOpts struct {
	Client         evmclient.Client
	LogBroadcaster log.Broadcaster
	LogPoller      logpoller.LogPoller
	GeneralConfig  legacyevm.AppConfig
	HeadTracker    httypes.HeadTracker
	DB             sqlutil.DataSource
	TxManager      txmgr.TxManager
	KeyStore       keystore.Eth
	MailMon        *mailbox.Monitor
	GasEstimator   gas.EvmFeeEstimator
}

// NewChainRelayExtenders returns a simple chain collection with one chain and
// allows to mock client/config on that chain
func NewChainRelayExtenders(t testing.TB, testopts TestChainOpts) *evmrelay.ChainRelayerExtenders {
	opts := NewChainRelayExtOpts(t, testopts)
	cc, err := evmrelay.NewChainRelayerExtenders(testutils.Context(t), opts)
	require.NoError(t, err)
	return cc
}

func NewChainRelayExtOpts(t testing.TB, testopts TestChainOpts) legacyevm.ChainRelayExtenderConfig {
	require.NotNil(t, testopts.KeyStore)
	lggr := logger.TestLogger(t)
	opts := legacyevm.ChainRelayExtenderConfig{
		Logger:   lggr,
		KeyStore: testopts.KeyStore,
		ChainOpts: legacyevm.ChainOpts{
			AppConfig:    testopts.GeneralConfig,
			MailMon:      testopts.MailMon,
			GasEstimator: testopts.GasEstimator,
			DS:           testopts.DB,
		},
	}
	opts.GenEthClient = func(*big.Int) evmclient.Client {
		if testopts.Client != nil {
			return testopts.Client
		}
		return evmclient.NewNullClient(MustGetDefaultChainID(t, testopts.GeneralConfig.EVMConfigs()), logger.TestLogger(t))
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
		opts.GenTxManager = func(*big.Int) txmgr.TxManager {
			return testopts.TxManager
		}
	}
	if opts.MailMon == nil {
		opts.MailMon = servicetest.Run(t, mailboxtest.NewMonitor(t))
	}
	if testopts.GasEstimator != nil {
		opts.GenGasEstimator = func(*big.Int) gas.EvmFeeEstimator {
			return testopts.GasEstimator
		}
	}

	return opts
}

// Deprecated, this is a replacement function for tests for now removed default evmChainID logic
func MustGetDefaultChainID(t testing.TB, evmCfgs evmtoml.EVMConfigs) *big.Int {
	if len(evmCfgs) == 0 {
		t.Fatalf("at least one evm chain config must be defined")
	}
	return evmCfgs[0].ChainID.ToInt()
}

// Deprecated, this is a replacement function for tests for now removed default chain logic
func MustGetDefaultChain(t testing.TB, cc legacyevm.LegacyChainContainer) legacyevm.Chain {
	if len(cc.Slice()) == 0 {
		t.Fatalf("at least one evm chain container must be defined")
	}

	return cc.Slice()[0]
}

type TestConfigs struct {
	mu sync.RWMutex
	evmtoml.EVMConfigs
}

var _ evmtypes.Configs = &TestConfigs{}

func NewTestConfigs(cs ...*evmtoml.EVMConfig) *TestConfigs {
	return &TestConfigs{EVMConfigs: evmtoml.EVMConfigs(cs)}
}

func (mo *TestConfigs) PutChains(cs ...evmtoml.EVMConfig) {
	mo.mu.Lock()
	defer mo.mu.Unlock()
chains:
	for i := range cs {
		id := cs[i].ChainID
		for j, c2 := range mo.EVMConfigs {
			if c2.ChainID == id {
				mo.EVMConfigs[j] = &cs[i] // replace
				continue chains
			}
		}
		mo.EVMConfigs = append(mo.EVMConfigs, &cs[i])
	}
}

func (mo *TestConfigs) Chains(chainIDs ...string) (cs []types.ChainStatus, count int, err error) {
	mo.mu.RLock()
	defer mo.mu.RUnlock()
	if len(chainIDs) == 0 {
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
		if !slices.Contains(chainIDs, chainID) {
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
func (mo *TestConfigs) Nodes(chainID string) (nodes []evmtypes.Node, err error) {
	mo.mu.RLock()
	defer mo.mu.RUnlock()

	for i := range mo.EVMConfigs {
		c := mo.EVMConfigs[i]
		if chainID == c.ChainID.String() {
			for _, n := range c.Nodes {
				nodes = append(nodes, legacyNode(n, c.ChainID))
			}
		}
	}
	err = fmt.Errorf("no nodes: chain %s: %w", chainID, chains.ErrNotFound)
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
	return evmtypes.Node{}, fmt.Errorf("node %s: %w", name, chains.ErrNotFound)
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

func legacyNode(n *evmtoml.Node, chainID *ubig.Big) (v2 evmtypes.Node) {
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

func nodeStatus(n *evmtoml.Node, chainID string) (types.NodeStatus, error) {
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
	c.On("IsL2").Return(false).Maybe()
	return c
}
