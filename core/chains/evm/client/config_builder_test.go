package client_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
)

func TestClientConfigBuilder(t *testing.T) {
	t.Parallel()

	selectionMode := ptr("HighestHead")
	leaseDuration := 0 * time.Second
	pollFailureThreshold := ptr(uint32(5))
	pollInterval := 10 * time.Second
	syncThreshold := ptr(uint32(5))
	nodeIsSyncingEnabled := ptr(false)
	chainTypeStr := ""
	nodeConfigs := []client.NodeConfig{
		{
			Name:    ptr("foo"),
			WSURL:   ptr("ws://foo.test"),
			HTTPURL: ptr("http://foo.test"),
		},
	}
	finalityDepth := ptr(uint32(10))
	finalityTagEnabled := ptr(true)
	noNewHeadsThreshold := time.Second
	chainCfg, nodePool, nodes, err := client.NewClientConfigs(selectionMode, leaseDuration, chainTypeStr, nodeConfigs,
		pollFailureThreshold, pollInterval, syncThreshold, nodeIsSyncingEnabled, noNewHeadsThreshold, finalityDepth, finalityTagEnabled)
	require.NoError(t, err)

	// Validate node pool configs
	require.Equal(t, *selectionMode, nodePool.SelectionMode())
	require.Equal(t, leaseDuration, nodePool.LeaseDuration())
	require.Equal(t, *pollFailureThreshold, nodePool.PollFailureThreshold())
	require.Equal(t, pollInterval, nodePool.PollInterval())
	require.Equal(t, *syncThreshold, nodePool.SyncThreshold())
	require.Equal(t, *nodeIsSyncingEnabled, nodePool.NodeIsSyncingEnabled())

	// Validate node configs
	require.Equal(t, *nodeConfigs[0].Name, *nodes[0].Name)
	require.Equal(t, *nodeConfigs[0].WSURL, (*nodes[0].WSURL).String())
	require.Equal(t, *nodeConfigs[0].HTTPURL, (*nodes[0].HTTPURL).String())

	// Validate chain config
	require.Equal(t, chainTypeStr, string(chainCfg.ChainType()))
	require.Equal(t, noNewHeadsThreshold, chainCfg.NodeNoNewHeadsThreshold())
	require.Equal(t, *finalityDepth, chainCfg.FinalityDepth())
	require.Equal(t, *finalityTagEnabled, chainCfg.FinalityTagEnabled())

	// let combiler tell us, when we do not have sufficient data to create evm client
	_ = client.NewEvmClient(nodePool, chainCfg, nil, logger.Test(t), big.NewInt(10), nodes)
}

func TestNodeConfigs(t *testing.T) {
	t.Parallel()

	t.Run("parsing unique node configs succeeds", func(t *testing.T) {
		nodeConfigs := []client.NodeConfig{
			{
				Name:    ptr("foo1"),
				WSURL:   ptr("ws://foo1.test"),
				HTTPURL: ptr("http://foo1.test"),
			},
			{
				Name:    ptr("foo2"),
				WSURL:   ptr("ws://foo2.test"),
				HTTPURL: ptr("http://foo2.test"),
			},
		}
		tomlNodes, err := client.ParseTestNodeConfigs(nodeConfigs)
		require.NoError(t, err)
		require.Len(t, tomlNodes, len(nodeConfigs))
	})

	t.Run("parsing missing ws url fails", func(t *testing.T) {
		nodeConfigs := []client.NodeConfig{
			{
				Name:    ptr("foo1"),
				HTTPURL: ptr("http://foo1.test"),
			},
		}
		_, err := client.ParseTestNodeConfigs(nodeConfigs)
		require.Error(t, err)
	})

	t.Run("parsing missing http url fails", func(t *testing.T) {
		nodeConfigs := []client.NodeConfig{
			{
				Name:  ptr("foo1"),
				WSURL: ptr("ws://foo1.test"),
			},
		}
		_, err := client.ParseTestNodeConfigs(nodeConfigs)
		require.Error(t, err)
	})

	t.Run("parsing invalid ws url fails", func(t *testing.T) {
		nodeConfigs := []client.NodeConfig{
			{
				Name:    ptr("foo1"),
				WSURL:   ptr("http://foo1.test"),
				HTTPURL: ptr("http://foo1.test"),
			},
		}
		_, err := client.ParseTestNodeConfigs(nodeConfigs)
		require.Error(t, err)
	})

	t.Run("parsing duplicate http url fails", func(t *testing.T) {
		nodeConfigs := []client.NodeConfig{
			{
				Name:    ptr("foo1"),
				WSURL:   ptr("ws://foo1.test"),
				HTTPURL: ptr("ws://foo1.test"),
			},
		}
		_, err := client.ParseTestNodeConfigs(nodeConfigs)
		require.Error(t, err)
	})

	t.Run("parsing duplicate node names fails", func(t *testing.T) {
		nodeConfigs := []client.NodeConfig{
			{
				Name:    ptr("foo1"),
				WSURL:   ptr("ws://foo1.test"),
				HTTPURL: ptr("http://foo1.test"),
			},
			{
				Name:    ptr("foo1"),
				WSURL:   ptr("ws://foo2.test"),
				HTTPURL: ptr("http://foo2.test"),
			},
		}
		_, err := client.ParseTestNodeConfigs(nodeConfigs)
		require.Error(t, err)
	})

	t.Run("parsing duplicate node ws urls fails", func(t *testing.T) {
		nodeConfigs := []client.NodeConfig{
			{
				Name:    ptr("foo1"),
				WSURL:   ptr("ws://foo1.test"),
				HTTPURL: ptr("http://foo1.test"),
			},
			{
				Name:    ptr("foo2"),
				WSURL:   ptr("ws://foo2.test"),
				HTTPURL: ptr("http://foo1.test"),
			},
		}
		_, err := client.ParseTestNodeConfigs(nodeConfigs)
		require.Error(t, err)
	})

	t.Run("parsing duplicate node http urls fails", func(t *testing.T) {
		nodeConfigs := []client.NodeConfig{
			{
				Name:    ptr("foo1"),
				WSURL:   ptr("ws://foo1.test"),
				HTTPURL: ptr("http://foo1.test"),
			},
			{
				Name:    ptr("foo2"),
				WSURL:   ptr("ws://foo1.test"),
				HTTPURL: ptr("http://foo2.test"),
			},
		}
		_, err := client.ParseTestNodeConfigs(nodeConfigs)
		require.Error(t, err)
	})

	t.Run("parsing order too large fails", func(t *testing.T) {
		nodeConfigs := []client.NodeConfig{
			{
				Name:    ptr("foo1"),
				WSURL:   ptr("ws://foo1.test"),
				HTTPURL: ptr("http://foo1.test"),
				Order:   ptr(int32(101)),
			},
		}
		_, err := client.ParseTestNodeConfigs(nodeConfigs)
		require.Error(t, err)
	})
}

func ptr[T any](t T) *T { return &t }
