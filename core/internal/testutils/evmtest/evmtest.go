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

	"github.com/smartcontractkit/chainlink-evm/pkg/chains"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	evmclient "github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/client/clienttest"
	evmconfig "github.com/smartcontractkit/chainlink-evm/pkg/config"
	configtoml "github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/gas"
	evmheads "github.com/smartcontractkit/chainlink-evm/pkg/heads"
	"github.com/smartcontractkit/chainlink-evm/pkg/log"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	"github.com/smartcontractkit/chainlink-evm/pkg/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	evmrelay "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

func NewChainScopedConfig(t testing.TB, cfg configtoml.HasEVMConfigs) evmconfig.ChainScopedConfig {
	var evmCfg *configtoml.EVMConfig
	if len(cfg.EVMConfigs()) > 0 {
		evmCfg = cfg.EVMConfigs()[0]
	} else {
		var chainID = (*ubig.Big)(testutils.FixtureChainID)
		evmCfg = &configtoml.EVMConfig{
			ChainID: chainID,
			Chain:   configtoml.Defaults(chainID),
		}
	}

	return evmconfig.NewTOMLChainScopedConfig(evmCfg)
}

type TestChainOpts struct {
	Client         evmclient.Client
	LogBroadcaster log.Broadcaster
	LogPoller      logpoller.LogPoller
	ChainConfigs   configtoml.EVMConfigs
	DatabaseConfig txmgr.DatabaseConfig
	FeatureConfig  legacyevm.FeatureConfig
	ListenerConfig txmgr.ListenerConfig
	HeadTracker    evmheads.Tracker
	DB             sqlutil.DataSource
	TxManager      txmgr.TxManager
	KeyStore       keystore.Eth
	MailMon        *mailbox.Monitor
	GasEstimator   gas.EvmFeeEstimator
}

// NewLegacyChains returns a simple chain collection with one chain and
// allows to mock client/config on that chain
func NewLegacyChains(t testing.TB, testopts TestChainOpts) *legacyevm.LegacyChains {
	lggr, ks, opts := NewChainOpts(t, testopts)
	cc, err := evmrelay.NewLegacyChainsAndConfig(lggr, ks, opts)
	require.NoError(t, err)
	return cc.NewLegacyChains()
}

func NewChainOpts(t testing.TB, testopts TestChainOpts) (logger.Logger, keystore.Eth, legacyevm.ChainOpts) {
	require.NotNil(t, testopts.KeyStore)
	lggr := logger.TestLogger(t)
	opts := legacyevm.ChainOpts{
		ChainConfigs:   testopts.ChainConfigs,
		DatabaseConfig: testopts.DatabaseConfig,
		ListenerConfig: testopts.ListenerConfig,
		FeatureConfig:  testopts.FeatureConfig,
		MailMon:        testopts.MailMon,
		GasEstimator:   testopts.GasEstimator,
		DS:             testopts.DB,
	}
	opts.GenEthClient = func(*big.Int) evmclient.Client {
		if testopts.Client != nil {
			return testopts.Client
		}
		return evmclient.NewNullClient(MustGetDefaultChainID(t, testopts.ChainConfigs), logger.TestLogger(t))
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
		opts.GenHeadTracker = func(*big.Int, evmheads.Broadcaster) evmheads.Tracker {
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

	return lggr, testopts.KeyStore, opts
}

// Deprecated, this is a replacement function for tests for now removed default evmChainID logic
func MustGetDefaultChainID(t testing.TB, evmCfgs configtoml.EVMConfigs) *big.Int {
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
	configtoml.EVMConfigs
}

var _ evmtypes.Configs = &TestConfigs{}

func NewTestConfigs(cs ...*configtoml.EVMConfig) *TestConfigs {
	return &TestConfigs{EVMConfigs: configtoml.EVMConfigs(cs)}
}

func (mo *TestConfigs) PutChains(cs ...configtoml.EVMConfig) {
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

func legacyNode(n *configtoml.Node, chainID *ubig.Big) (v2 evmtypes.Node) {
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

func nodeStatus(n *configtoml.Node, chainID string) (types.NodeStatus, error) {
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

// Deprecated: use clienttest.NewClient
func NewEthClientMock(t *testing.T) *clienttest.Client {
	return clienttest.NewClient(t)
}

// Deprecated: use clienttest.NewClientWithDefaultChainID
func NewEthClientMockWithDefaultChain(t *testing.T) *clienttest.Client {
	c := NewEthClientMock(t)
	c.On("ConfiguredChainID").Return(testutils.FixtureChainID).Maybe()
	c.On("IsL2").Return(false).Maybe()
	return c
}
