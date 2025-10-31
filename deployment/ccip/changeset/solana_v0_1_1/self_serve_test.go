package solana_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	chainSelectors "github.com/smartcontractkit/chain-selectors"
	lockrelease "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/lockrelease_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	cldf_solana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/stretchr/testify/require"

	cldfChain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	solCommon "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_common"
	solState "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/state"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/testcontext"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	ccipChangesetSolana "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/solana_v0_1_1"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
)

func TestOnboardTokenPoolForSelfServeWithMCMs(t *testing.T) {
	// TODO: Fix this test
	skipInCI(t)
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
	solChainSelector := tenv.Env.BlockChains.ListChainSelectors(cldfChain.WithFamily(chainSelectors.FamilySolana))[0]
	e, lnrTokenMint, err := deployTokenAndMint(t, tenv.Env, solChainSelector, []string{}, "TEST_TOKEN")
	require.NoError(t, err)
	tenv.Env = e
	e, bnmTokenMint, err := deployTokenAndMint(t, tenv.Env, solChainSelector, []string{}, "TEST_TOKEN_2")
	require.NoError(t, err)
	tenv.Env = e
	customerAdmin, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	state, err := stateview.LoadOnchainStateSolana(e)
	require.NoError(t, err)
	err = modifyMintAuthority(e.BlockChains.SolanaChains()[solChainSelector], tenv.Env.BlockChains.SolanaChains()[solChainSelector].DeployerKey.PublicKey(), lnrTokenMint, customerAdmin.PublicKey())
	require.NoError(t, err)
	err = modifyMintAuthority(e.BlockChains.SolanaChains()[solChainSelector], tenv.Env.BlockChains.SolanaChains()[solChainSelector].DeployerKey.PublicKey(), bnmTokenMint, customerAdmin.PublicKey())
	require.NoError(t, err)
	lockAndReleaseTokenPoolProgramID := state.SolChains[solChainSelector].LockReleaseTokenPools[shared.CLLMetadata]
	burnAndMintTokenPoolProgramID := state.SolChains[solChainSelector].BurnMintTokenPools[shared.CLLMetadata]
	var mcmsConfig *proposalutils.TimelockConfig
	if isMCMsOwner {
		timelockSignerPDA, _ := testhelpers.TransferOwnershipSolanaV0_1_1(t, &e, solChainSelector, true,
			ccipChangesetSolana.CCIPContractsToTransfer{
				Router: true,
			})
		e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
			commonchangeset.Configure(
				cldf.CreateLegacyChangeSet(ccipChangesetSolana.SetUpgradeAuthorityChangeset),
				ccipChangesetSolana.SetUpgradeAuthorityConfig{
					ChainSelector:         solChainSelector,
					NewUpgradeAuthority:   timelockSignerPDA,
					SetAfterInitialDeploy: true,
					SetOffRamp:            true,
					SetMCMSPrograms:       true,
					TransferKeys: []solana.PublicKey{
						lockAndReleaseTokenPoolProgramID,
						burnAndMintTokenPoolProgramID,
					},
				},
			),
		})
		mcmsConfig = &proposalutils.TimelockConfig{
			MinDelay: 1 * time.Second,
		}
	}
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			// Setup needed for the token pool program
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.InitGlobalConfigTokenPoolProgram),
			ccipChangesetSolana.TokenPoolConfigWithMCM{
				ChainSelector: solChainSelector,
				MCMS:          mcmsConfig,
				TokenPoolConfigs: []ccipChangesetSolana.TokenPoolConfig{
					{
						PoolType: shared.LockReleaseTokenPool,
						Metadata: shared.CLLMetadata,
					},
					{
						PoolType: shared.BurnMintTokenPool,
						Metadata: shared.CLLMetadata,
					},
				},
			},
		),
		commonchangeset.Configure(
			// Actual changeset to test
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.OnboardTokenPoolsForSelfServe),
			ccipChangesetSolana.OnboardTokenPoolsForSelfServeConfig{
				ChainSelector: solChainSelector,
				RegisterTokenConfigs: []ccipChangesetSolana.OnboardTokenPoolConfig{
					{
						TokenMint:     lnrTokenMint,
						ProposedOwner: customerAdmin.PublicKey(),
						PoolType:      shared.LockReleaseTokenPool,
					},
					{
						TokenMint:     bnmTokenMint,
						ProposedOwner: customerAdmin.PublicKey(),
						PoolType:      shared.BurnMintTokenPool,
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
	tokenAdminRegistryPDA, _, err := solState.FindTokenAdminRegistryPDA(lnrTokenMint, state.SolChains[solChainSelector].Router)
	require.NoError(t, err)
	err = e.BlockChains.SolanaChains()[solChainSelector].GetAccountDataBorshInto(ctx, tokenAdminRegistryPDA, &tokenAdminRegistryAccount)
	require.NoError(t, err)
	// the actual administrator needs to accept the role
	require.Equal(t, solana.PublicKey{}, tokenAdminRegistryAccount.Administrator)
	// pending administrator should be the proposed admin key
	require.Equal(t, customerAdmin.PublicKey(), tokenAdminRegistryAccount.PendingAdministrator)

	var tokenPoolAccount lockrelease.State
	// Verify that the token pool PDA was initialized
	tokenPoolPDA, err := tokens.TokenPoolConfigAddress(lnrTokenMint, lockAndReleaseTokenPoolProgramID)
	require.NoError(t, err)
	err = e.BlockChains.SolanaChains()[solChainSelector].GetAccountDataBorshInto(ctx, tokenPoolPDA, &tokenPoolAccount)
	require.NoError(t, err)
	// Verify the mint address is correct
	require.Equal(t, lnrTokenMint, tokenPoolAccount.Config.Mint)
	// Verify the proposed owner is correct
	require.Equal(t, customerAdmin.PublicKey(), tokenPoolAccount.Config.ProposedOwner)

	anotherCustomerAdmin, err := solana.NewRandomPrivateKey()
	require.NoError(t, err)
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			// Actual changeset to test
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.OnboardTokenPoolsForSelfServe),
			ccipChangesetSolana.OnboardTokenPoolsForSelfServeConfig{
				ChainSelector: solChainSelector,
				RegisterTokenConfigs: []ccipChangesetSolana.OnboardTokenPoolConfig{
					{
						TokenMint:     lnrTokenMint,
						ProposedOwner: anotherCustomerAdmin.PublicKey(),
						PoolType:      shared.LockReleaseTokenPool,
						Override:      true,
					},
				},
				MCMS: mcmsConfig,
			},
		),
	},
	)
	require.NoError(t, err)

	var tokenAdminRegistryAccount2 solCommon.TokenAdminRegistry
	// Verify that the proposed admin in the token admin registry was updated
	err = e.BlockChains.SolanaChains()[solChainSelector].GetAccountDataBorshInto(ctx, tokenAdminRegistryPDA, &tokenAdminRegistryAccount2)
	require.NoError(t, err)
	require.Equal(t, anotherCustomerAdmin.PublicKey(), tokenAdminRegistryAccount2.PendingAdministrator)

	var tokenPoolAccount2 lockrelease.State
	// Verify the proposed owner is updated
	err = e.BlockChains.SolanaChains()[solChainSelector].GetAccountDataBorshInto(ctx, tokenPoolPDA, &tokenPoolAccount2)
	require.NoError(t, err)
	require.Equal(t, anotherCustomerAdmin.PublicKey(), tokenPoolAccount2.Config.ProposedOwner)
}

func modifyMintAuthority(state cldf_solana.Chain, deployerKey solana.PublicKey, mint solana.PublicKey, newAuthority solana.PublicKey) error {
	mintI, err := token.NewSetAuthorityInstruction(token.AuthorityMintTokens, newAuthority, mint, deployerKey, []solana.PublicKey{}).ValidateAndBuild()
	if err != nil {
		return err
	}
	mintWrap := &tokens.TokenInstruction{mintI, solana.TokenProgramID}
	if err := state.Confirm([]solana.Instruction{mintWrap}); err != nil {
		return fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return nil
}
