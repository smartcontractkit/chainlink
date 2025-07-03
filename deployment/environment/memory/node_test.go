package memory

import (
	"maps"
	"slices"
	"testing"

	"github.com/smartcontractkit/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

func TestNode(t *testing.T) {
	blockchains := cldf_chain.NewBlockChainsFromSlice(NewMemoryChainsEVM(t, 3, 5))
	ports := freeport.GetN(t, 1)
	c := NewNodeConfig{
		Port:           ports[0],
		BlockChains:    blockchains,
		LogLevel:       zapcore.DebugLevel,
		Bootstrap:      false,
		RegistryConfig: deployment.CapabilityRegistryConfig{},
		CustomDBSetup:  nil,
	}
	node := NewNode(t, c)
	// We expect 3 transmitter keys
	keys, err := node.App.GetKeyStore().Eth().GetAll(t.Context())
	require.NoError(t, err)
	require.Len(t, keys, 3)
	// We expect 3 chains supported
	evmChains := node.App.GetRelayers().List(chainlink.FilterRelayersByType(relay.NetworkEVM)).Slice()
	require.Len(t, evmChains, 3)

	t.Run("DeploymentNode", func(t *testing.T) {
		dn, err := node.DeploymentNode()
		require.NoError(t, err)
		assert.Equal(t, node.Keys.PeerID, dn.PeerID)
		assert.Equal(t, node.Keys.CSA.ID(), dn.CSAKey)
		assert.Len(t, dn.SelToOCRConfig, 3)
		gotChains := make([]uint64, len(dn.SelToOCRConfig))
		i := 0
		for k := range dn.SelToOCRConfig {
			gotChains[i] = k.ChainSelector
			i++
		}
		assert.ElementsMatch(t, slices.Collect(maps.Keys(blockchains.EVMChains())), gotChains)
	})

	t.Run("JDChainConfigs", func(t *testing.T) {
		jdChainConfigs, err := node.JDChainConfigs()
		require.NoError(t, err)
		assert.Len(t, jdChainConfigs, 3)
		for i, cc := range jdChainConfigs {
			assert.Equal(t, node.Keys.PeerID.String(), cc.Ocr2Config.P2PKeyBundle.PeerId, "chain %d", i)
		}
	})
}
