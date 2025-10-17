package solana

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
)

// setupTest sets up a test runtime with a single solana chain with deployed the MCMS and Timelock
// contracts
func setupTest(t *testing.T) (*runtime.Runtime, uint64) {
	memory.DownloadSolanaProgramArtifactsForTest(t)

	// Setup the runtime with preloaded programs. The address book is updated with the preloaded programs.
	selector := chainselectors.TEST_22222222222222222222222222222222222222222222.Selector
	ab := getPreloadedAddressBook(t, selector)
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithSolanaContainer(t, []uint64{selector}, memory.ProgramsPath, memory.SolanaProgramIDs),
		environment.WithAddressBook(ab),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	// Deploy MCMS and Timelock
	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployMCMSWithTimelockV2), map[uint64]commontypes.MCMSWithTimelockConfigV2{
			selector: proposalutils.SingleGroupTimelockConfigV2(t),
		}),
	)
	require.NoError(t, err)

	return rt, selector
}

// getPreloadedAddressBook returns an address book with the preloaded MCMS Solana addresses for
// the given selector.
func getPreloadedAddressBook(t *testing.T, selector uint64) *cldf.AddressBookMap {
	t.Helper()

	ab := cldf.NewMemoryAddressBook()

	tv := cldf.NewTypeAndVersion(commontypes.ManyChainMultisigProgram, deployment.Version1_0_0)
	err := ab.Save(selector, memory.SolanaProgramIDs["mcm"], tv)
	require.NoError(t, err)

	tv = cldf.NewTypeAndVersion(commontypes.AccessControllerProgram, deployment.Version1_0_0)
	err = ab.Save(selector, memory.SolanaProgramIDs["access_controller"], tv)
	require.NoError(t, err)

	tv = cldf.NewTypeAndVersion(commontypes.RBACTimelockProgram, deployment.Version1_0_0)
	err = ab.Save(selector, memory.SolanaProgramIDs["timelock"], tv)
	require.NoError(t, err)

	return ab
}

// fundSignerPDAs funds the timelock signer and MCMS signer PDAs with 1 SOL
func fundSignerPDAs(
	t *testing.T, chain cldf_solana.Chain, mcmsState *state.MCMSWithTimelockStateSolana,
) {
	t.Helper()

	timelockSignerPDA := state.GetTimelockSignerPDA(mcmsState.TimelockProgram, mcmsState.TimelockSeed)
	mcmSignerPDA := state.GetMCMSignerPDA(mcmsState.McmProgram, mcmsState.ProposerMcmSeed)
	signerPDAs := []solana.PublicKey{timelockSignerPDA, mcmSignerPDA}
	err := memory.FundSolanaAccounts(t.Context(), signerPDAs, 1, chain.Client)
	require.NoError(t, err)
}
