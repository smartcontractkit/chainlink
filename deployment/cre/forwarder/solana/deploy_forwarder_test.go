package solana

import (
	"crypto/ecdsa"
	"os"
	"testing"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/wsrpc/logger"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	cldftesthelpers "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils/testhelpers"

	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/helpers"
	"github.com/smartcontractkit/chainlink/deployment/internal/soltestutils"
)

// Tests with transfer upgrade authority require downloading and building artifacts
// from chainlink-solana
// so we disable them in CI since it will take too long to run
func TestDeployForwarder(t *testing.T) {
	skipInCI(t)

	t.Parallel()

	selector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithSolanaContainer(t, []uint64{selector}, t.TempDir(), map[string]string{}),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	chain := rt.Environment().BlockChains.SolanaChains()[selector]

	t.Run("should deploy forwarder", func(t *testing.T) {
		err := rt.Exec(
			runtime.ChangesetTask(DeployForwarder{},
				&DeployForwarderRequest{
					ChainSel:  selector,
					Qualifier: testQualifier,
					Version:   testVersion,
					BuildConfig: &helpers.BuildSolanaConfig{
						GitCommitSha:   "3305b4d55b5469e110133e5a36e5600aadf436fb",
						DestinationDir: chain.ProgramsPath,
						LocalBuild:     helpers.LocalBuildConfig{BuildLocally: true, CreateDestinationDir: true},
					},
				},
			),
		)
		require.NoError(t, err)
	})

	// Rebuilding is what makes this the realistic upgrade path: the artifacts are built again and
	// have to declare the deployed program ID for the upgraded program to keep working.
	t.Run("should upgrade forwarder", func(t *testing.T) {
		err := rt.Exec(
			runtime.ChangesetTask(UpgradeForwarder{},
				&UpgradeForwarderRequest{
					ChainSel:   selector,
					Qualifier:  testQualifier,
					Version:    testVersion,
					NewVersion: testUpgradedVersion,
					BuildConfig: &helpers.BuildSolanaConfig{
						GitCommitSha:   "3305b4d55b5469e110133e5a36e5600aadf436fb",
						DestinationDir: chain.ProgramsPath,
						LocalBuild:     helpers.LocalBuildConfig{BuildLocally: true, CreateDestinationDir: true},
					},
				},
			),
		)
		require.NoError(t, err)

		// The upgrade happens in place, so the new version points at the program deployed above.
		ds := rt.Environment().DataStore.Addresses()
		before, err := ds.Get(datastore.NewAddressRefKey(selector, ForwarderContract, semver.MustParse(testVersion), testQualifier))
		require.NoError(t, err)
		after, err := ds.Get(datastore.NewAddressRefKey(selector, ForwarderContract, semver.MustParse(testUpgradedVersion), testQualifier))
		require.NoError(t, err)
		require.Equal(t, before.Address, after.Address)
	})

	t.Run("should pass upgrade authority", func(t *testing.T) {
		err := rt.Exec(
			runtime.ChangesetTask(SetForwarderUpgradeAuthority{},
				&SetForwarderUpgradeAuthorityRequest{
					ChainSel:            selector,
					Qualifier:           testQualifier,
					Version:             testVersion,
					NewUpgradeAuthority: chain.DeployerKey.PublicKey(),
				},
			),
		)
		require.NoError(t, err)
	})
}

func TestConfigureForwarder(t *testing.T) {
	t.Parallel()

	// Setup the solana programs
	programsPath := t.TempDir()
	programIDs := soltestutils.LoadKeystonePrograms(t, programsPath)

	t.Run("set config without mcms", func(t *testing.T) {
		te := setupForwarderTestEnv(t, programsPath, programIDs, false)

		err := te.Runtime.Exec(
			runtime.ChangesetTask(DeployForwarder{},
				&DeployForwarderRequest{
					ChainSel:  te.Selector,
					Qualifier: testQualifier,
					Version:   testVersion,
				},
			),
			runtime.ChangesetTask(ConfigureForwarders{},
				&ConfigureForwarderRequest{
					DON:       te.DON,
					Version:   testVersion,
					Qualifier: testQualifier,
				},
			),
		)
		require.NoError(t, err)

		requireOraclesConfig(t, te.Runtime.Environment(), te.Selector, te.DON.ID, te.DON.Version, true)
	})

	t.Run("set config with mcms", func(t *testing.T) {
		te := setupForwarderTestEnv(t, programsPath, programIDs, true)

		// Deploy the forwarder and transfer ownership to the MCMS
		err := te.Runtime.Exec(
			runtime.ChangesetTask(DeployForwarder{},
				&DeployForwarderRequest{
					ChainSel:  te.Selector,
					Qualifier: testQualifier,
					Version:   testVersion,
				},
			),
			runtime.ChangesetTask(TransferOwnershipForwarder{},
				&TransferOwnershipForwarderRequest{
					ChainSel:  te.Selector,
					MCMSCfg:   cldfproposalutils.TimelockConfig{MinDelay: 1 * time.Second},
					Qualifier: testQualifier,
					Version:   testVersion,
				},
			),
			runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{cldftesthelpers.TestXXXMCMSSigner}),
		)
		require.NoError(t, err)

		// Configure the forwarder using MCMS
		err = te.Runtime.Exec(
			runtime.ChangesetTask(ConfigureForwarders{},
				&ConfigureForwarderRequest{
					DON:       te.DON,
					Version:   testVersion,
					Qualifier: testQualifier,
					MCMS: &cldfproposalutils.TimelockConfig{
						MinDelay: time.Second,
					},
				},
			),
			runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{cldftesthelpers.TestXXXMCMSSigner}),
		)
		require.NoError(t, err)

		requireOraclesConfig(t, te.Runtime.Environment(), te.Selector, te.DON.ID, te.DON.Version, true)
	})
}

func TestUpgradeForwarder(t *testing.T) {
	t.Parallel()
	skipInCI(t)
	// Setup the solana programs
	programsPath := t.TempDir()
	programIDs := soltestutils.LoadKeystonePrograms(t, programsPath)

	// The upgrade rewrites the same artifact, which exercises the buffer write and the in place
	// upgrade of the deployed program. The binary does not grow, so the program data account does
	// not have to be extended.
	t.Run("upgrade without mcms", func(t *testing.T) {
		te := setupForwarderTestEnv(t, programsPath, programIDs, false)

		err := te.Runtime.Exec(
			runtime.ChangesetTask(DeployForwarder{},
				&DeployForwarderRequest{
					ChainSel:  te.Selector,
					Qualifier: testQualifier,
					Version:   testVersion,
				},
			),
			runtime.ChangesetTask(ConfigureForwarders{},
				&ConfigureForwarderRequest{
					DON:       te.DON,
					Version:   testVersion,
					Qualifier: testQualifier,
				},
			),
		)
		require.NoError(t, err)

		err = te.Runtime.Exec(
			runtime.ChangesetTask(UpgradeForwarder{},
				&UpgradeForwarderRequest{
					ChainSel:   te.Selector,
					Qualifier:  testQualifier,
					Version:    testVersion,
					NewVersion: testUpgradedVersion,
				},
			),
		)
		require.NoError(t, err)

		// The upgrade happens in place, so the new version points at the same program and state.
		env := te.Runtime.Environment()
		for _, contractType := range []datastore.ContractType{ForwarderContract, ForwarderState} {
			before, err := env.DataStore.Addresses().Get(
				datastore.NewAddressRefKey(te.Selector, contractType, semver.MustParse(testVersion), testQualifier))
			require.NoError(t, err)
			after, err := env.DataStore.Addresses().Get(
				datastore.NewAddressRefKey(te.Selector, contractType, semver.MustParse(testUpgradedVersion), testQualifier))
			require.NoError(t, err)
			require.Equal(t, before.Address, after.Address, "%s address must survive the upgrade", contractType)
		}

		// The config written before the upgrade is untouched, and the upgraded program still serves
		// instructions addressed to it.
		requireOraclesConfig(t, env, te.Selector, te.DON.ID, te.DON.Version, true)
		err = te.Runtime.Exec(
			runtime.ChangesetTask(ConfigureForwarders{},
				&ConfigureForwarderRequest{
					DON:       te.DON,
					Version:   testUpgradedVersion,
					Qualifier: testQualifier,
				},
			),
		)
		require.NoError(t, err)
	})

	t.Run("upgrade with mcms", func(t *testing.T) {
		te := setupForwarderTestEnv(t, programsPath, programIDs, true)

		timelockSigner, err := helpers.FetchTimelockSigner(
			te.Runtime.Environment().DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(te.Selector)))
		require.NoError(t, err)

		// Deploy the forwarder and hand the upgrade authority to the timelock
		err = te.Runtime.Exec(
			runtime.ChangesetTask(DeployForwarder{},
				&DeployForwarderRequest{
					ChainSel:  te.Selector,
					Qualifier: testQualifier,
					Version:   testVersion,
				},
			),
			runtime.ChangesetTask(SetForwarderUpgradeAuthority{},
				&SetForwarderUpgradeAuthorityRequest{
					ChainSel:            te.Selector,
					Qualifier:           testQualifier,
					Version:             testVersion,
					NewUpgradeAuthority: timelockSigner,
				},
			),
		)
		require.NoError(t, err)

		// Upgrade through a timelock proposal
		err = te.Runtime.Exec(
			runtime.ChangesetTask(UpgradeForwarder{},
				&UpgradeForwarderRequest{
					ChainSel:   te.Selector,
					Qualifier:  testQualifier,
					Version:    testVersion,
					NewVersion: testUpgradedVersion,
					MCMS: &cldfproposalutils.TimelockConfig{
						MinDelay: time.Second,
					},
				},
			),
			runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{cldftesthelpers.TestXXXMCMSSigner}),
		)
		require.NoError(t, err)

		// The upgraded program still serves instructions addressed to it. Here the upgrade lands
		// when the proposal executes rather than inside the changeset, so the test has to wait out
		// the slot in which the program became invisible itself.
		require.Eventually(t, func() bool {
			return te.Runtime.Exec(
				runtime.ChangesetTask(ConfigureForwarders{},
					&ConfigureForwarderRequest{
						DON:       te.DON,
						Version:   testUpgradedVersion,
						Qualifier: testQualifier,
					},
				),
			) == nil
		}, time.Minute, 2*time.Second, "upgraded forwarder must accept instructions")

		requireOraclesConfig(t, te.Runtime.Environment(), te.Selector, te.DON.ID, te.DON.Version, true)
	})
}

func TestWithPinnedUpgradeKey(t *testing.T) {
	t.Parallel()

	programID := solana.MustPublicKeyFromBase58("whV7Q5pi17hPPyaPksToDw1nMx6Lh8qmNWKFaLRQ4wz")
	otherKey := "3NUQ6J1XaaApjNYWqGG5HQzbCuooQobyfsar45cSQB92"

	t.Run("pins the deployed program id", func(t *testing.T) {
		got := withPinnedUpgradeKey(helpers.BuildSolanaConfig{
			LocalBuild: helpers.LocalBuildConfig{BuildLocally: true},
		}, programID)

		require.Equal(t, programID.String(), got.LocalBuild.UpgradeKeys[keystoneForwarder])
	})

	t.Run("keeps a key the caller pinned", func(t *testing.T) {
		got := withPinnedUpgradeKey(helpers.BuildSolanaConfig{
			LocalBuild: helpers.LocalBuildConfig{
				BuildLocally: true,
				UpgradeKeys:  map[cldf.ContractType]string{keystoneForwarder: otherKey},
			},
		}, programID)

		require.Equal(t, otherKey, got.LocalBuild.UpgradeKeys[keystoneForwarder])
	})

	t.Run("leaves the callers config untouched", func(t *testing.T) {
		upgradeKeys := map[cldf.ContractType]string{"other_program": otherKey}
		in := helpers.BuildSolanaConfig{
			LocalBuild: helpers.LocalBuildConfig{BuildLocally: true, UpgradeKeys: upgradeKeys},
		}

		got := withPinnedUpgradeKey(in, programID)

		require.Equal(t, map[cldf.ContractType]string{"other_program": otherKey}, upgradeKeys)
		require.Equal(t, programID.String(), got.LocalBuild.UpgradeKeys[keystoneForwarder])
	})

	t.Run("does not pin keys of a downloaded build", func(t *testing.T) {
		got := withPinnedUpgradeKey(helpers.BuildSolanaConfig{GitCommitSha: "deadbeef"}, programID)

		require.Empty(t, got.LocalBuild.UpgradeKeys)
	})
}

func skipInCI(t *testing.T) {
	ci := os.Getenv("CI") == "true"
	if ci {
		t.Skip("Skipping in CI")
	}
}
