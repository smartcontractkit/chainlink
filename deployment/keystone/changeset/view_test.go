package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestKeystoneView(t *testing.T) {
	t.Parallel()
	env := memory.NewMemoryEnvironment(t, logger.Test(t), zapcore.DebugLevel, memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 2,
	})
	registryChain := env.AllChainSelectors()[0]
	resp, err := DeployCapabilityRegistry(env, registryChain)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NoError(t, env.ExistingAddresses.Merge(resp.AddressBook))
	resp, err = DeployOCR3(env, registryChain)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NoError(t, env.ExistingAddresses.Merge(resp.AddressBook))
	resp, err = DeployForwarder(env, DeployForwarderRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NoError(t, env.ExistingAddresses.Merge(resp.AddressBook))

	a, err := ViewKeystone(env)
	require.NoError(t, err)
	b, err := a.MarshalJSON()
	require.NoError(t, err)
	require.NotEmpty(t, b)
	t.Log(string(b))
}
