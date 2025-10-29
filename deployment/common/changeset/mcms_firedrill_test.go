package changeset_test

import (
	"testing"

	"github.com/gagliardetto/solana-go"
	mcmsTypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/internal/soltestutils"
	"github.com/smartcontractkit/chainlink/deployment/utils/solutils"
)

func TestMCMSSignFireDrillChangeset(t *testing.T) {
	t.Parallel()

	evmSelector1 := chain_selectors.TEST_90000001.Selector
	evmSelector2 := chain_selectors.TEST_90000002.Selector
	solSelector := chain_selectors.TEST_22222222222222222222222222222222222222222222.Selector
	programsPath, programIDs, ab := soltestutils.PreloadMCMS(t, solSelector)

	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{evmSelector1, evmSelector2}),
		environment.WithSolanaContainer(t, []uint64{solSelector}, programsPath, programIDs),
		environment.WithAddressBook(ab),
		environment.WithLogger(logger.Test(t)),
	))
	require.NoError(t, err)

	solChain := rt.Environment().BlockChains.SolanaChains()[solSelector]

	// Deploy MCMS and Timelock
	config := proposalutils.SingleGroupTimelockConfigV2(t)

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.DeployMCMSWithTimelockV2), map[uint64]commontypes.MCMSWithTimelockConfigV2{
			evmSelector1: config,
			evmSelector2: config,
			solSelector:  config,
		}),
	)
	require.NoError(t, err)

	// Fund the signer PDAs for the MCMS contracts
	mcmsState := soltestutils.GetMCMSStateFromAddressBook(t, rt.State().AddressBook, solChain)

	timelockSigner := state.GetTimelockSignerPDA(mcmsState.TimelockProgram, mcmsState.TimelockSeed)
	mcmSigner := state.GetMCMSignerPDA(mcmsState.McmProgram, mcmsState.ProposerMcmSeed)
	mcmSignerBypasser := state.GetMCMSignerPDA(mcmsState.McmProgram, mcmsState.BypasserMcmSeed)

	// Note we cannot use FundSignerPDAs here because we also have to fund the bypasser signer PDA.
	err = solutils.FundAccounts(t.Context(),
		solChain.Client,
		[]solana.PublicKey{timelockSigner, mcmSigner, mcmSignerBypasser},
		150,
	)
	require.NoError(t, err)

	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(commonchangeset.MCMSSignFireDrillChangeset), commonchangeset.FireDrillConfig{
			Selectors: []uint64{evmSelector1, evmSelector2, solSelector},
			TimelockCfg: proposalutils.TimelockConfig{
				MCMSAction: mcmsTypes.TimelockActionBypass,
			},
		}),
	)
	require.NoError(t, err)
}
