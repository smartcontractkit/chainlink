package changeset_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

func TestDeployLinkToken(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1,
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	chainSelector := env.AllChainSelectors()[0]

	resp, err := changeset.DeployLinkToken(env, chainSelector)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// LinkToken should be deployed on chain 0
	addrs, err := resp.AddressBook.AddressesForChain(chainSelector)
	require.NoError(t, err)
	require.Len(t, addrs, 1)

	// nothing on chain 1
	require.NotEqual(t, chainSelector, env.AllChainSelectors()[1])
	oaddrs, _ := resp.AddressBook.AddressesForChain(env.AllChainSelectors()[1])
	assert.Len(t, oaddrs, 0)
}
