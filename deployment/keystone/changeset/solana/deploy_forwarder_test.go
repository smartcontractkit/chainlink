package solana

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
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
		cfg := helpers.BuildSolanaConfig{
			GitCommitSha:   "6442d0e438ca175b1b2ce059a174ba4bf4e8afc1",
			DestinationDir: "./contracts",
			LocalBuild: helpers.LocalBuildConfig{
				BuildLocally:         true,
				CleanDestinationDir:  false,
				CreateDestinationDir: true,
				CleanGitDir:          false,
			},
		}

		// default solChain looking for contracts in ccip directory
		chain := env.SolChains[registrySel]
		chain.ProgramsPath = getProgramsPath()
		env.SolChains[registrySel] = chain
		// deploy forwarder
		env.ExistingAddresses = ab
		resp, err := DeployForwarder(env, &DeployRequest{
			ChainSel:    registrySel,
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

func getProgramsPath() string {
	// Get the directory of the current file (environment.go)
	_, currentFile, _, _ := runtime.Caller(0)
	// Go up to the root of the deployment package
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
	// Construct the absolute path
	return filepath.Join(rootDir, "changeset/solana", "contracts")
}
