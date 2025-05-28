package solana

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

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
	ab := cldf.NewMemoryAddressBook()

	t.Run("should deploy forwarder", func(t *testing.T) {
		env = shouldDeployForwarder(t, env, registrySel, ab)
	})
	t.Run("should pass upgrade authority", func(t *testing.T) {
		_, err := SetForwarderUpgradeAuthority(env, &SetForwarderUpgradeAuthorityRequest{
			ChainSel:            registrySel,
			NewUpgradeAuthority: env.SolChains[registrySel].DeployerKey.PublicKey(),
		})
		require.NoError(t, err)
	})
}

func shouldDeployForwarder(t *testing.T, env cldf.Environment, registrySel uint64, ab cldf.AddressBook) cldf.Environment {
	cfg := helpers.BuildSolanaConfig{
		GitCommitSha:   "d047073ea230f965626716029f8d902729ddffed",
		DestinationDir: "./contracts",
		LocalBuild: helpers.LocalBuildConfig{
			BuildLocally:         true,
			CleanDestinationDir:  true,
			CreateDestinationDir: true,
			CleanGitDir:          true,
		},
	}
	// defult
	fmt.Println(cfg.GitCommitSha)
	// default solChain looking for contracts in ccip directory
	chain := env.SolChains[registrySel]
	chain.ProgramsPath = getProgramsPath()
	env.SolChains[registrySel] = chain
	// deploy forwarder
	env.ExistingAddresses = ab
	resp, err := DeployForwarder(env, &DeployRequest{
		ChainSel:    registrySel,
		BuildConfig: nil,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	// forwarder should be deployed on registry chain
	addrs, err := resp.AddressBook.AddressesForChain(registrySel)
	require.NoError(t, err)
	require.Len(t, addrs, 2) // forwarder programID, forwarder state
	env.ExistingAddresses = resp.AddressBook
	return env
}

func getProgramsPath() string {
	// Get the directory of the current file (environment.go)
	_, currentFile, _, _ := runtime.Caller(0)
	// Go up to the root of the deployment package
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
	// Construct the absolute path
	return filepath.Join(rootDir, "changeset/solana", "contracts")
}
