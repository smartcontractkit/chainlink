package workflowregistry

import (
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func Test_Deploy(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1, // nodes unused but required in config
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	registrySel := env.AllChainSelectors()[0]

	resp, err := Deploy(env, registrySel)
	require.NoError(t, err)
	require.NotNil(t, resp)
	// workflow registry should be deployed on chain 0
	addrs, err := resp.AddressBook.AddressesForChain(registrySel)
	require.NoError(t, err)
	require.Len(t, addrs, 1)

	// assert nothing on chain 1
	require.NotEqual(t, registrySel, env.AllChainSelectors()[1])
	oaddrs, _ := resp.AddressBook.AddressesForChain(env.AllChainSelectors()[1])
	assert.Empty(t, oaddrs)
}
