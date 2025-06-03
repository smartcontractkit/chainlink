package solana

import (
	"fmt"
	"math/big"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	solanaMCMS "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
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

func TestConfigureForwarder(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		nChains      int
		ExcludeChain bool // if true, configuration should be applied to all except one chain
	}{
		{
			nChains: 1,
		},
		{
			nChains:      3,
			ExcludeChain: true,
		},
	}

	t.Run("set config without mcms", func(t *testing.T) {
		for _, tcase := range testCases {
			nChains := tcase.nChains
			name := fmt.Sprintf("nChains=%d", nChains)

			t.Run(name, func(t *testing.T) {
				lggr := logger.Test(t)
				env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, memory.MemoryEnvironmentConfig{
					Nodes:     1, // nodes unused but required in config
					SolChains: 1,
				})

				solSel := env.AllChainSelectorsSolana()[0]
				te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
					WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4, ChainSelectors: []uint64{solSel}},
					AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
					WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
					NumChains:       nChains,
				})

				te.Env.SolChains = env.SolChains

				env = shouldDeployForwarder(t, te.Env, solSel, te.Env.ExistingAddresses)

				var wfNodes []string
				for _, id := range te.GetP2PIDs("wfDon") {
					wfNodes = append(wfNodes, id.String())
				}

				cfg := ConfigureForwarderRequest{
					WFDonName:        "test-wf-don",
					WFNodeIDs:        wfNodes,
					RegistryChainSel: te.RegistrySelector,
				}
				_, err := ConfigureForwarders(te.Env, cfg)
				require.NoError(t, err)
			})
		}
	})

	t.Run("set config with mcms", func(t *testing.T) {
		for _, tcase := range testCases {
			nChains := tcase.nChains
			name := fmt.Sprintf("nChains=%d", nChains)

			t.Run(name, func(t *testing.T) {
				lggr := logger.Test(t)
				env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, memory.MemoryEnvironmentConfig{
					Nodes:     1, // nodes unused but required in config
					SolChains: 1,
				})

				solSel := env.AllChainSelectorsSolana()[0]
				te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
					WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4, ChainSelectors: []uint64{solSel}},
					AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
					WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
					NumChains:       nChains,
				})

				te.Env.SolChains = env.SolChains

				env = shouldDeployForwarder(t, te.Env, solSel, te.Env.ExistingAddresses)
				_, err := solanaMCMS.DeployMCMSWithTimelockProgramsSolana(env, env.SolChains[solSel], env.ExistingAddresses,
					commontypes.MCMSWithTimelockConfigV2{
						Canceller:        proposalutils.SingleGroupMCMSV2(t),
						Proposer:         proposalutils.SingleGroupMCMSV2(t),
						Bypasser:         proposalutils.SingleGroupMCMSV2(t),
						TimelockMinDelay: big.NewInt(0),
					},
				)

				require.NoError(t, err)

				var wfNodes []string
				for _, id := range te.GetP2PIDs("wfDon") {
					wfNodes = append(wfNodes, id.String())
				}

				cfg := ConfigureForwarderRequest{
					WFDonName:        "test-wf-don",
					WFNodeIDs:        wfNodes,
					RegistryChainSel: te.RegistrySelector,
					MCMS: &proposalutils.TimelockConfig{
						MinDelay: time.Millisecond,
					}}

				out, err := ConfigureForwarders(te.Env, cfg)
				require.NoError(t, err)
				require.Equal(t, 1, len(out.MCMSTimelockProposals))
			})
		}
	})
}

func shouldDeployForwarder(t *testing.T, env cldf.Environment, registrySel uint64, _ cldf.AddressBook) cldf.Environment {
	cfg := helpers.BuildSolanaConfig{
		GitCommitSha:   "d047073ea230f965626716029f8d902729ddffed",
		DestinationDir: "./solana_contracts",
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
	env.ExistingAddresses.Merge(resp.AddressBook)
	return env
}

// prepare tests
// 1. clone chainlink-solana, chainlink-ccip
// 2. move mcp + timelock from ccip to chainlink-solana
func prepareTests() {

}

func getProgramsPath() string {
	// Get the directory of the current file (environment.go)
	_, currentFile, _, _ := runtime.Caller(0)
	// Go up to the root of the deployment package
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
	// Construct the absolute path
	return filepath.Join(rootDir, "changeset/solana", "solana_contracts")
}
