package metrics

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
)

func TestNewMetrics(t *testing.T) {
	t.Parallel()

	donConfig := &config.DONConfig{
		Members: []config.NodeConfig{
			{Address: "0xnode1", Name: "node1"},
			{Address: "0xnode2", Name: "node2"},
		},
	}
	metrics, err := NewMetrics(donConfig)
	require.NoError(t, err)
	require.NotNil(t, metrics)
	require.NotNil(t, metrics.action)
	require.NotNil(t, metrics.trigger)
	require.NotNil(t, metrics.common)
	require.Equal(t, "node1", metrics.nodeAddressToNodeName["0xnode1"])
	require.Equal(t, "node2", metrics.nodeAddressToNodeName["0xnode2"])
}
