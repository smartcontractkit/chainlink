package example

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func Test_ExemplarDeployLinkToken(t *testing.T) {
	t.Parallel()

	lggr := logger.TestLogger(t)
	e := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain1 := e.BlockChains.ListChainSelectors()[0]

	result, err := ExemplarDeployLinkToken{}.Apply(e, chain1)
	require.NoError(t, err)

	// Check that one address ref was created
	addresRefs, err := result.DataStore.Addresses().Fetch()
	require.NoError(t, err)
	require.Len(t, addresRefs, 1)

	// Check that one contract metadata ref was created
	contractMetadata, err := result.DataStore.ContractMetadata().Fetch()
	require.NoError(t, err)
	require.Len(t, contractMetadata, 1)

	// Check that env metadata was set correctly
	_, err = result.DataStore.EnvMetadata().Get()
	require.NoError(t, err)
}
