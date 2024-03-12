package client_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestNewEvmClient(t *testing.T) {
	t.Parallel()

	noNewHeadsThreshold := 3 * time.Minute
	selectionMode := ptr("HighestHead")
	leaseDuration := 0 * time.Second
	pollFailureThreshold := ptr(uint32(5))
	pollInterval := 10 * time.Second
	syncThreshold := ptr(uint32(5))
	nodeIsSyncingEnabled := ptr(false)
	chainTypeStr := ""
	nodeConfigs := []client.TestNodeConfig{
		{
			Name:    ptr("foo"),
			WSURL:   ptr("ws://foo.test"),
			HTTPURL: ptr("http://foo.test"),
		},
	}
	nodePool, nodes, chainType, err := client.NewClientConfigs(selectionMode, leaseDuration, chainTypeStr, nodeConfigs, pollFailureThreshold, pollInterval, syncThreshold, nodeIsSyncingEnabled)
	require.NoError(t, err)

	client := client.NewEvmClient(nodePool, noNewHeadsThreshold, logger.TestLogger(t), testutils.FixtureChainID, chainType, nodes)
	require.NotNil(t, client)
}
