package solana

import (
	"crypto/ecdsa"
	"os"
	"testing"
	"time"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/wsrpc/logger"
	"github.com/stretchr/testify/require"

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

func skipInCI(t *testing.T) {
	ci := os.Getenv("CI") == "true"
	if ci {
		t.Skip("Skipping in CI")
	}
}
