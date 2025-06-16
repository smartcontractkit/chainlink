package solana

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/wsrpc/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	solanaMCMS "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
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

	// replace default program path since memory env sets it to ccip
	chain := env.BlockChains.SolanaChains()[solSel]
	chain.ProgramsPath = getProgramsPath()
	env.BlockChains = cldfchain.NewBlockChains(map[uint64]cldfchain.BlockChain{solSel: chain})

	t.Run("should deploy forwarder", func(t *testing.T) {
		configuredChangeset := commonchangeset.Configure(DeployForwarder{},
			&DeployForwarderRequest{
				ChainSel:  solSel,
				Qualifier: testQualifier,
				Version:   "1.0.0",
				BuildConfig: &helpers.BuildSolanaConfig{
					GitCommitSha:   "ba5a33ab378020fac73bda72b6bc2f9ae6bddb83",
					DestinationDir: getProgramsPath(),
					LocalBuild:     helpers.LocalBuildConfig{BuildLocally: true, CreateDestinationDir: true},
				},
			},
		)

		// deploy
		var err error
		env, _, err = commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{configuredChangeset})
		require.NoError(t, err)
	})

	t.Run("should pass upgrade authority", func(t *testing.T) {
		configuredChangeset := commonchangeset.Configure(SetForwarderUpgradeAuthority{},
			&SetForwarderUpgradeAuthorityRequest{
				ChainSel:            solSel,
				Qualifier:           testQualifier,
				Version:             "1.0.0",
				NewUpgradeAuthority: chain.DeployerKey.PublicKey(),
			},
		)

		// deploy
		var err error
		_, _, err = commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{configuredChangeset})
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

				// configure don for solana chain
				te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
					WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4, ChainSelectors: []uint64{solSel}},
					AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
					WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
					NumChains:       nChains,
				})

				solChain := env.BlockChains.SolanaChains()[solSel]
				blockchains := make(map[uint64]cldfchain.BlockChain)

				solChain.ProgramsPath = getProgramsPath()
				blockchains[solSel] = solChain

				for _, ch := range te.Env.BlockChains.All() {
					blockchains[ch.ChainSelector()] = ch
				}

				te.Env.BlockChains = cldfchain.NewBlockChains(blockchains)

				deployChangeset := commonchangeset.Configure(DeployForwarder{},
					&DeployForwarderRequest{
						ChainSel:  solSel,
						Qualifier: testQualifier,
						Version:   "1.0.0",
					},
				)

				var wfNodes []string
				for _, id := range te.GetP2PIDs("wfDon") {
					wfNodes = append(wfNodes, id.String())
				}

				cfg := ConfigureForwarderRequest{
					WFDonName:        "test-wf-don",
					WFNodeIDs:        wfNodes,
					RegistryChainSel: te.RegistrySelector,
					Version:          "1.0.0",
					Qualifier:        testQualifier,
				}

				configureChangeset := commonchangeset.Configure(ConfigureForwarders{},
					&cfg,
				)

				var err error
				_, _, err = commonchangeset.ApplyChangesets(t, te.Env, []commonchangeset.ConfiguredChangeSet{deployChangeset, configureChangeset})
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

				solChain.ProgramsPath = getProgramsPath()
				blockchains[solSel] = solChain

				for _, ch := range te.Env.BlockChains.All() {
					blockchains[ch.ChainSelector()] = ch
				}

				te.Env.BlockChains = cldfchain.NewBlockChains(blockchains)

				ds := datastore.NewMemoryDataStore()

				// deploy mcms
				mcmsState, err := solanaMCMS.DeployMCMSWithTimelockProgramsSolanaV2(env, ds, solChain,
					commontypes.MCMSWithTimelockConfigV2{
						Canceller:        proposalutils.SingleGroupMCMSV2(t),
						Proposer:         proposalutils.SingleGroupMCMSV2(t),
						Bypasser:         proposalutils.SingleGroupMCMSV2(t),
						TimelockMinDelay: big.NewInt(0),
					},
				)
				require.NoError(t, err)

				te.Env.DataStore = ds.Seal()
				fundSignerPDAs(t, te.Env, solSel, mcmsState)

				deployChangeset := commonchangeset.Configure(DeployForwarder{},
					&DeployForwarderRequest{
						ChainSel:  solSel,
						Qualifier: testQualifier,
						Version:   "1.0.0",
					},
				)

				var wfNodes []string
				for _, id := range te.GetP2PIDs("wfDon") {
					wfNodes = append(wfNodes, id.String())
				}

				cfg := ConfigureForwarderRequest{
					WFDonName:        "test-wf-don",
					WFNodeIDs:        wfNodes,
					RegistryChainSel: te.RegistrySelector,
					Version:          "1.0.0",
					Qualifier:        testQualifier,
					MCMS: &proposalutils.TimelockConfig{
						MinDelay: time.Second,
					},
				}

				configureChangeset := commonchangeset.Configure(ConfigureForwarders{},
					&cfg,
				)

				transferOwnershipChangeset := commonchangeset.Configure(TransferOwnershipForwarder{},
					&TransferOwnershipForwarderRequest{
						ChainSel:  solSel,
						MCMSCfg:   proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
						Qualifier: testQualifier,
						Version:   "1.0.0",
					})

				_, _, err = commonchangeset.ApplyChangesets(t, te.Env, []commonchangeset.ConfiguredChangeSet{deployChangeset, transferOwnershipChangeset,
					configureChangeset})
				require.NoError(t, err)
			})
		}
	})
}

const (
	testQualifier = "test-deploy"
)

func getProgramsPath() string {
	// Get the directory of the current file (environment.go)
	_, currentFile, _, _ := runtime.Caller(0)
	// Go up to the root of the deployment package
	rootDir := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
	// Construct the absolute path
	return filepath.Join(rootDir, "changeset/solana", "solana_contracts")
}
func fundSignerPDAs(
	t *testing.T, env cldf.Environment, chainSelector uint64, chainState *state.MCMSWithTimelockStateSolana,
) {
	t.Helper()
	solChain := env.BlockChains.SolanaChains()[chainSelector]
	timelockSignerPDA := state.GetTimelockSignerPDA(chainState.TimelockProgram, chainState.TimelockSeed)
	mcmSignerPDA := state.GetMCMSignerPDA(chainState.McmProgram, chainState.ProposerMcmSeed)
	signerPDAs := []solana.PublicKey{timelockSignerPDA, mcmSignerPDA}
	err := memory.FundSolanaAccounts(env.GetContext(), signerPDAs, 1, solChain.Client)
	require.NoError(t, err)
}

func skipInCI(t *testing.T) {
	ci := os.Getenv("CI") == "true"
	if ci {
		t.Skip("Skipping in CI")
	}
}
