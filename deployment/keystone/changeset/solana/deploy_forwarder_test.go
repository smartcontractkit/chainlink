package solana

import (
	"crypto/ecdsa"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/wsrpc/logger"
	"github.com/stretchr/testify/require"

	cldfchain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/onchain"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	solanaMCMS "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/helpers"
	"github.com/smartcontractkit/chainlink/deployment/internal/soltestutils"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
)

const (
	testQualifier = "test-deploy"
)

// Tests with transfer upgrade authority require downloading and building artifacts
// from chainlink-solana
// so we disable them in CI since it will take too long to run
func TestDeployForwarder(t *testing.T) {
	skipInCI(t)
	t.Parallel()

	selector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	env, err := environment.New(t.Context(),
		environment.WithSolanaContainer(t, []uint64{selector}, t.TempDir(), map[string]string{}),
		environment.WithLogger(logger.Test(t)),
	)
	require.NoError(t, err)

	chain := env.BlockChains.SolanaChains()[selector]

	t.Run("should deploy forwarder", func(t *testing.T) {
		configuredChangeset := commonchangeset.Configure(DeployForwarder{},
			&DeployForwarderRequest{
				ChainSel:  selector,
				Qualifier: testQualifier,
				Version:   "1.0.0",
				BuildConfig: &helpers.BuildSolanaConfig{
					GitCommitSha:   "3305b4d55b5469e110133e5a36e5600aadf436fb",
					DestinationDir: chain.ProgramsPath,
					LocalBuild:     helpers.LocalBuildConfig{BuildLocally: true, CreateDestinationDir: true},
				},
			},
		)

		// deploy
		var err error
		_, _, err = commonchangeset.ApplyChangesets(t, *env, []commonchangeset.ConfiguredChangeSet{configuredChangeset})
		require.NoError(t, err)
	})

	t.Run("should pass upgrade authority", func(t *testing.T) {
		configuredChangeset := commonchangeset.Configure(SetForwarderUpgradeAuthority{},
			&SetForwarderUpgradeAuthorityRequest{
				ChainSel:            selector,
				Qualifier:           testQualifier,
				Version:             "1.0.0",
				NewUpgradeAuthority: chain.DeployerKey.PublicKey(),
			},
		)

		// deploy
		var err error
		_, _, err = commonchangeset.ApplyChangesets(t, *env, []commonchangeset.ConfiguredChangeSet{configuredChangeset})
		require.NoError(t, err)
	})
}

func TestConfigureForwarder(t *testing.T) {
	t.Parallel()

	// Setup the solana programs
	programsPath := t.TempDir()
	solSel := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	programIDs := soltestutils.LoadKeystonePrograms(t, programsPath)

	t.Run("set config without mcms", func(t *testing.T) {
		// Initialize a solana chain
		solChain, err := onchain.
			NewSolanaContainerLoader(programsPath, programIDs).
			Load(t, []uint64{solSel})
		require.NoError(t, err)

		// Configure don without the solana chain
		te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
			WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4, ChainSelectors: []uint64{solSel}},
			AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
			WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
			NumChains:       1,
		})

		// Inject the solana chain into the environment
		blockchains := make(map[uint64]cldfchain.BlockChain)
		blockchains[solSel] = solChain[0]
		for _, ch := range te.Env.BlockChains.All() {
			blockchains[ch.ChainSelector()] = ch
		}
		te.Env.BlockChains = cldfchain.NewBlockChains(blockchains)

		// Populate the datastore with the keystone forwarder contract
		ds := datastore.NewMemoryDataStore()
		err = ds.AddressRefStore.Add(datastore.AddressRef{
			Address:       programIDs["keystone_forwarder"],
			ChainSelector: solSel,
			Type:          ForwarderContract,
			Version:       semver.MustParse("1.0.0"),
			Qualifier:     testQualifier,
		})
		require.NoError(t, err)
		te.Env.DataStore = ds.Seal()

		// We set up a new runtime to execute the changesets based on the previously set up environment
		rt := runtime.NewFromEnvironment(te.Env)
		require.NoError(t, err)

		var wfNodes []string
		for _, id := range te.GetP2PIDs("wfDon") {
			wfNodes = append(wfNodes, id.String())
		}

		err = rt.Exec(
			runtime.ChangesetTask(DeployForwarder{},
				&DeployForwarderRequest{
					ChainSel:  solSel,
					Qualifier: testQualifier,
					Version:   "1.0.0",
				},
			),
			runtime.ChangesetTask(ConfigureForwarders{},
				&ConfigureForwarderRequest{
					WFDonName:        "test-wf-don",
					WFNodeIDs:        wfNodes,
					RegistryChainSel: te.RegistrySelector,
					Version:          "1.0.0",
					Qualifier:        testQualifier,
				},
			),
		)
		require.NoError(t, err)
	})

	t.Run("set config with mcms", func(t *testing.T) {
		// Initialize a solana chain
		solChains, err := onchain.
			NewSolanaContainerLoader(programsPath, programIDs).
			Load(t, []uint64{solSel})
		require.NoError(t, err)

		te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
			WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4, ChainSelectors: []uint64{solSel}},
			AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
			WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
			NumChains:       1,
		})

		// Inject the solana chain into the environment
		blockchains := make(map[uint64]cldfchain.BlockChain)
		blockchains[solSel] = solChains[0]
		for _, ch := range te.Env.BlockChains.All() {
			blockchains[ch.ChainSelector()] = ch
		}
		te.Env.BlockChains = cldfchain.NewBlockChains(blockchains)

		ds := datastore.NewMemoryDataStore()
		soltestutils.RegisterMCMSPrograms(t, solSel, ds)

		err = ds.AddressRefStore.Add(datastore.AddressRef{
			Address:       programIDs["keystone_forwarder"],
			ChainSelector: solSel,
			Type:          ForwarderContract,
			Version:       semver.MustParse("1.0.0"),
			Qualifier:     testQualifier,
		})
		require.NoError(t, err)

		te.Env.DataStore = ds.Seal()

		rt := runtime.NewFromEnvironment(te.Env)
		require.NoError(t, err)

		mcmsState, err := solanaMCMS.DeployMCMSWithTimelockProgramsSolanaV2(
			rt.Environment(),
			ds,
			rt.Environment().BlockChains.SolanaChains()[solSel],
			commontypes.MCMSWithTimelockConfigV2{
				Canceller:        proposalutils.SingleGroupMCMSV2(t),
				Proposer:         proposalutils.SingleGroupMCMSV2(t),
				Bypasser:         proposalutils.SingleGroupMCMSV2(t),
				TimelockMinDelay: big.NewInt(0),
			},
		)
		require.NoError(t, err)

		chain := te.Env.BlockChains.SolanaChains()[solSel]
		soltestutils.FundSignerPDAs(t, chain, mcmsState)

		var wfNodes []string
		for _, id := range te.GetP2PIDs("wfDon") {
			wfNodes = append(wfNodes, id.String())
		}

		// Deploy the forwarder and transfer ownership to the MCMS
		err = rt.Exec(
			runtime.ChangesetTask(DeployForwarder{},
				&DeployForwarderRequest{
					ChainSel:  solSel,
					Qualifier: testQualifier,
					Version:   "1.0.0",
				},
			),
			runtime.ChangesetTask(TransferOwnershipForwarder{},
				&TransferOwnershipForwarderRequest{
					ChainSel:  solSel,
					MCMSCfg:   proposalutils.TimelockConfig{MinDelay: 1 * time.Second},
					Qualifier: testQualifier,
					Version:   "1.0.0",
				},
			),
			runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{proposalutils.TestXXXMCMSSigner}),
		)
		require.NoError(t, err)

		// Configure the forwarder using MCMS
		err = rt.Exec(
			runtime.ChangesetTask(ConfigureForwarders{},
				&ConfigureForwarderRequest{
					WFDonName:        "test-wf-don",
					WFNodeIDs:        wfNodes,
					RegistryChainSel: te.RegistrySelector,
					Version:          "1.0.0",
					Qualifier:        testQualifier,
					MCMS: &proposalutils.TimelockConfig{
						MinDelay: time.Second,
					},
				},
			),
			runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{proposalutils.TestXXXMCMSSigner}),
		)
		require.NoError(t, err)
	})
}

func skipInCI(t *testing.T) {
	ci := os.Getenv("CI") == "true"
	if ci {
		t.Skip("Skipping in CI")
	}
}
