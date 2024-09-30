package memory

import (
	"testing"

	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
)

func TestNode(t *testing.T) {
	chains := GenerateChains(t, 3)
	ports := freeport.GetN(t, 1)
	node := NewNode(t, ports[0], chains, zapcore.DebugLevel, false, deployment.CapabilityRegistryConfig{})
	// We expect 3 transmitter keys
	keys, err := node.App.GetKeyStore().Eth().GetAll(Context(t))
	require.NoError(t, err)
	require.Len(t, keys, 3)
	// We expect 3 chains supported
	evmChains := node.App.GetRelayers().LegacyEVMChains().Slice()
	require.NoError(t, err)
	require.Len(t, evmChains, 3)
}
