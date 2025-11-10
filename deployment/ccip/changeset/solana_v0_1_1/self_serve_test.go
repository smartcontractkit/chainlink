package solana_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"

	chainSelectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	lockrelease "github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/lockrelease_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
	cldfsolana "github.com/smartcontractkit/chainlink-deployments-framework/chain/solana"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"

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

		// Print out deployer key
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
		require.NoError(t, err)
		tenv.Env = e
		mcmsConfig = &proposalutils.TimelockConfig{
			MinDelay: 1 * time.Second,
		}
		progDataAddr, err := deployment.GetProgramDataAddress(e.BlockChains.SolanaChains()[solChainSelector].Client, lockAndReleaseTokenPoolProgramID)
		require.NoError(t, err)
		upgradeAuthority, _, err := deployment.GetUpgradeAuthority(e.BlockChains.SolanaChains()[solChainSelector].Client, progDataAddr)
		require.NoError(t, err)
		require.Equal(t, timelockSignerPDA, upgradeAuthority)
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
						TokenMint:        lnrTokenMint,
						TokenProgramName: shared.SPLTokens,
						ProposedOwner:    customerAdmin.PublicKey(),
						Metadata:         customerAdmin.PublicKey().String(),
						PoolType:         shared.LockReleaseTokenPool,
					},
					{
						TokenMint:        bnmTokenMint,
						TokenProgramName: shared.SPLTokens,
						ProposedOwner:    customerAdmin.PublicKey(),
						Metadata:         customerAdmin.PublicKey().String(),
						PoolType:         shared.BurnMintTokenPool,
					},
				},
				MCMS: mcmsConfig,
			},
		),
	},
	)
	require.NoError(t, err)
	tenv.Env = e

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
	// Test with override
	e, _, err = commonchangeset.ApplyChangesets(t, e, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(
			cldf.CreateLegacyChangeSet(ccipChangesetSolana.OnboardTokenPoolsForSelfServe),
			ccipChangesetSolana.OnboardTokenPoolsForSelfServeConfig{
				ChainSelector: solChainSelector,
				RegisterTokenConfigs: []ccipChangesetSolana.OnboardTokenPoolConfig{
					{
						TokenMint:        lnrTokenMint,
						TokenProgramName: shared.SPLTokens,
						ProposedOwner:    anotherCustomerAdmin.PublicKey(),
						Metadata:         anotherCustomerAdmin.PublicKey().String(),
						PoolType:         shared.LockReleaseTokenPool,
						Override:         true,
					},
				},
				MCMS: mcmsConfig,
			},
		),
	},
	)
	require.NoError(t, err)
	tenv.Env = e

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

func modifyMintAuthority(state cldfsolana.Chain, deployerKey solana.PublicKey, mint solana.PublicKey, newAuthority solana.PublicKey) error {
	mintI, err := token.NewSetAuthorityInstruction(token.AuthorityMintTokens, newAuthority, mint, deployerKey, []solana.PublicKey{}).ValidateAndBuild()
	if err != nil {
		return err
	}
	mintWrap := &tokens.TokenInstruction{Instruction: mintI, Program: solana.TokenProgramID}
	if err := state.Confirm([]solana.Instruction{mintWrap}); err != nil {
		return fmt.Errorf("failed to confirm instructions: %w", err)
	}
	return nil
}
