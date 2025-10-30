package solana_test

import (
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	chainSelectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	cldfChain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	solCommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestOnboardTokenPoolForSelfServeWithMCMs(t *testing.T) {
	t.Parallel()
	doTestOnboardTokenPoolForSelfServe(t, true)
}

func TestOnboardTokenPoolForSelfServeWithoutMCMs(t *testing.T) {
	t.Parallel()
	doTestOnboardTokenPoolForSelfServe(t, false)
}

func doTestOnboardTokenPoolForSelfServe(t *testing.T, isMCMsOwner bool) {
	ctx := testcontext.Get(t)
	tenv, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithSolChains(1), testhelpers.WithCCIPSolanaContractVersion(ccipChangesetSolana.SolanaContractV0_1_1))
	solChain := tenv.Env.BlockChains.ListChainSelectors(cldfChain.WithFamily(chainSelectors.FamilySolana))[0]
	e, tokenAddress, err := deployTokenAndMint(t, tenv.Env, solChain, []string{}, "TEST_TOKEN")
	require.NoError(t, err)
	state, err := stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	customerAdmin, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)

	var mcmsConfig *proposalutils.TimelockConfig
	if isMCMsOwner {
		_, _ = testhelpers.TransferOwnershipSolanaV0_1_1(t, &e, solChain, true,
			ccipChangesetSolana.CCIPContractsToTransfer{
				Router: true,
			})
		mcmsConfig = &proposalutils.TimelockConfig{
			MinDelay: 1 * time.Second,
		}
	}

	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.OnboardTokenPoolsForSelfServe),
			ccipChangesetSolana.OnboardTokenPoolsForSelfServeConfig{
				ChainSelector: solChain,
				RegisterTokenConfigs: []ccipChangesetSolana.OnboardTokenPoolConfig{
					{
						TokenPubKey:             tokenAddress,
						TokenAdminRegistryAdmin: customerAdmin.PublicKey(),
					},
				},
				MCMS: mcmsConfig,
			},
		),
	},
	)
	require.NoError(t, err)

	var tokenAdminRegistryAccount solCommon.TokenAdminRegistry

	// Verify that the proposed admin in the token admin registry was updated
	tokenAdminRegistryPDA, _, _ := solState.FindTokenAdminRegistryPDA(tokenAddress, state.SolChains[solChain].Router)
	err = e.BlockChains.SolanaChains()[solChain].GetAccountDataBorshInto(ctx, tokenAdminRegistryPDA, &tokenAdminRegistryAccount)
	require.NoError(t, err)
	// the actual administrator needs to accept the role
	require.Equal(t, solana.PublicKey{}, tokenAdminRegistryAccount.Administrator)
	// pending administrator should be the proposed admin key
	require.Equal(t, customerAdmin.PublicKey(), tokenAdminRegistryAccount.PendingAdministrator)
}
