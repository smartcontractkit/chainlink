package changeset_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

func TestDeployBalanceReader(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:  1, // nodes unused but required in config
		Chains: 2,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)

	registrySel := env.AllChainSelectors()[0]

	t.Run("should deploy balancereader", func(t *testing.T) {
		ab := deployment.NewMemoryAddressBook()

		// deploy balancereader
		env.ExistingAddresses = ab
		resp, err := changeset.DeployBalanceReader(env, changeset.DeployBalanceReaderRequest{})
		require.NoError(t, err)
		require.NotNil(t, resp)
		// registry, ocr3, balancereader should be deployed on registry chain
		addrs, err := resp.AddressBook.AddressesForChain(registrySel)
		require.NoError(t, err)
		require.Len(t, addrs, 1)

		// only balancereader on chain 1
		require.NotEqual(t, registrySel, env.AllChainSelectors()[1])
		oaddrs, err := resp.AddressBook.AddressesForChain(env.AllChainSelectors()[1])
		require.NoError(t, err)
		require.Len(t, oaddrs, 1)
	})
}
