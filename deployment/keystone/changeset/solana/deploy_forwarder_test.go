package solana

import (
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/wsrpc/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestDeployForwarder(t *testing.T) {
	t.Parallel()

	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:     1, // nodes unused but required in config
		SolChains: 1,
	}
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	registrySel := env.AllChainSelectorsSolana()[0]

	t.Run("should deploy forwarder", func(t *testing.T) {
		ab := cldf.NewMemoryAddressBook()
		cfg := BuildSolanaConfig{
			GitCommitSha:   "6442d0e438ca175b1b2ce059a174ba4bf4e8afc1",
			DestinationDir: "./contracts",
			LocalBuild: LocalBuildConfig{
				BuildLocally:         true,
				CleanDestinationDir:  false,
				CreateDestinationDir: true,
				CleanGitDir:          true,
			},
		}
		// deploy forwarder
		env.ExistingAddresses = ab
		//	resp, err := changeset.DeployForwarder(env, changeset.DeployForwarderRequest{})
		resp, err := DeployForwarder(env, &DeployRequest{
			ChainSel:    registrySel,
			Qualifier:   "my-test-forwarder",
			BuildConfig: &cfg,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
		// registry, ocr3, forwarder should be deployed on registry chain
		addrs, err := resp.AddressBook.AddressesForChain(registrySel)
		require.NoError(t, err)
		require.Len(t, addrs, 1)
		fa := resp.DataStore.Addresses().Filter(datastore.AddressRefByQualifier("my-test-forwarder"))
		require.Len(t, fa, 1, "expected to find 'my-test-forwarder' qualifier")
	})
}
