package solana

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/wsrpc/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	solanaMCMS "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
)

// Tests require downloading and building artifacts
// from chainlink-solana and chainlink-ccip
// so we disable them in CI since it will take too long to run
func TestDeployForwarder(t *testing.T) {
	skipInCI(t)
	t.Parallel()

	lggr := logger.Test(t)
	cfg := memory.MemoryEnvironmentConfig{
		Nodes:     1, // nodes unused but required in config
		SolChains: 1,
	}

	env := memory.NewMemoryEnvironment(t, lggr, zapcore.DebugLevel, cfg)
	solSel := env.BlockChains.ListChainSelectors(cldfchain.WithFamily(chain_selectors.FamilySolana))[0]
	ab := cldf.NewMemoryAddressBook()

	t.Run("should deploy forwarder", func(t *testing.T) {
		env = shouldDeployForwarder(t, env, solSel, ab)
	})
	t.Run("should pass upgrade authority", func(t *testing.T) {
		_, err := SetForwarderUpgradeAuthority(env, &SetForwarderUpgradeAuthorityRequest{
			ChainSel:            solSel,
			NewUpgradeAuthority: env.BlockChains.SolanaChains()[solSel].DeployerKey.PublicKey(),
		})
		require.NoError(t, err)
	})
}

func TestConfigureForwarder(t *testing.T) {
	skipInCI(t)
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

				solSel := env.BlockChains.ListChainSelectors(cldfchain.WithFamily(chain_selectors.FamilySolana))[0]
				te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
					WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4, ChainSelectors: []uint64{solSel}},
					AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
					WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
					NumChains:       nChains,
				})
				solChain := env.BlockChains.SolanaChains()[solSel]
				blockchains := make(map[uint64]cldfchain.BlockChain)
				blockchains[solSel] = solChain
				for _, ch := range te.Env.BlockChains.All() {
					blockchains[ch.ChainSelector()] = ch
				}

				te.Env.BlockChains = cldfchain.NewBlockChains(blockchains)
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

				solSel := env.BlockChains.ListChainSelectors(cldfchain.WithFamily(chain_selectors.FamilySolana))[0]
				te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
					WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4, ChainSelectors: []uint64{solSel}},
					AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
					WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
					NumChains:       nChains,
				})

				solChain := env.BlockChains.SolanaChains()[solSel]
				blockchains := make(map[uint64]cldfchain.BlockChain)
				blockchains[solSel] = solChain
				for _, ch := range te.Env.BlockChains.All() {
					blockchains[ch.ChainSelector()] = ch
				}

				te.Env.BlockChains = cldfchain.NewBlockChains(blockchains)

				env = shouldDeployForwarder(t, te.Env, solSel, te.Env.ExistingAddresses)
				_, err := solanaMCMS.DeployMCMSWithTimelockProgramsSolana(env, env.BlockChains.SolanaChains()[solSel], env.ExistingAddresses,
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
				require.Len(t, out.MCMSTimelockProposals, 1)
			})
		}
	})
}

func shouldDeployForwarder(t *testing.T, env cldf.Environment, solSel uint64, _ cldf.AddressBook) cldf.Environment {
	cfg := helpers.BuildSolanaConfig{
		GitCommitSha:   "d047073ea230f965626716029f8d902729ddffed",
		DestinationDir: "./solana_contracts",
		LocalBuild: helpers.LocalBuildConfig{
			BuildLocally:         true,
			CreateDestinationDir: true,
		},
	}
	// defult
	fmt.Println(cfg.GitCommitSha)
	// default solChain looking for contracts in ccip directory
	chain := env.BlockChains.SolanaChains()[solSel]
	chain.ProgramsPath = getProgramsPath()
	env.BlockChains = cldfchain.NewBlockChains(map[uint64]cldfchain.BlockChain{solSel: chain})

	// deploy forwarder
	resp, err := DeployForwarder(env, &DeployRequest{
		ChainSel:    solSel,
		BuildConfig: nil,
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	// forwarder should be deployed on registry chain
	addrs, err := resp.AddressBook.AddressesForChain(solSel) //nolint:staticcheck migrate to datastore
	require.NoError(t, err)
	require.Len(t, addrs, 2)                            // forwarder programID, forwarder state
	err = env.ExistingAddresses.Merge(resp.AddressBook) //nolint:staticcheck migrate to datastore

	require.NoError(t, err)
	return env
}

func getProgramsPath() string {
	// Get the directory of the current file (environment.go)
	_, currentFile, _, _ := runtime.Caller(0)
	// Go up to the root of the deployment package
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
	// Construct the absolute path
	return filepath.Join(rootDir, "changeset/solana", "solana_contracts")
}

func skipInCI(t *testing.T) {
	ci := os.Getenv("CI") == "true"
	if ci {
		t.Skip("Skipping in CI")
	}
}
